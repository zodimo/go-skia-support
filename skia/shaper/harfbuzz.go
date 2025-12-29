package shaper

import (
	"log"

	"github.com/go-text/typesetting/di"
	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/language"
	"github.com/go-text/typesetting/shaping"
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
	"golang.org/x/image/math/fixed"
)

// UseGoTextFace is an interface that a Typeface can implement
// to provide the underlying go-text/typesetting Face.
type UseGoTextFace interface {
	GoTextFace() *font.Face
}

// HarfbuzzShaper implements the Shaper interface using HarfBuzz shaper from go-text/typesetting.
type HarfbuzzShaper struct {
	hb shaping.HarfbuzzShaper
}

// NewHarfbuzzShaper creates a new instance of HarfbuzzShaper.
func NewHarfbuzzShaper() *HarfbuzzShaper {
	return &HarfbuzzShaper{}
}

// Shape shapes the text using the font and runHandler.
func (s *HarfbuzzShaper) Shape(text string, font interfaces.SkFont, leftToRight bool, width float32, runHandler RunHandler, features []Feature) {
	// Create trivial iterators
	totalLength := len(text) // Byte length (approximation for now, iterators handle mapping)

	fontIter := NewTrivialFontRunIterator(font, totalLength)
	bidiDir := uint8(0) // LTR
	if !leftToRight {
		bidiDir = 1 // RTL
	}
	bidiIter := NewTrivialBiDiRunIterator(bidiDir, totalLength)
	scriptIter := NewTrivialScriptRunIterator(0, totalLength) // Common/Unknown
	langIter := NewTrivialLanguageRunIterator("en", totalLength)

	s.ShapeWithIterators(text, fontIter, bidiIter, scriptIter, langIter, features, width, runHandler)
}

// ShapeWithIterators shapes the text using custom iterators.
func (s *HarfbuzzShaper) ShapeWithIterators(text string,
	fontIter FontRunIterator,
	bidiIter BiDiRunIterator,
	scriptIter ScriptRunIterator,
	langIter LanguageRunIterator,
	features []Feature,
	width float32,
	runHandler RunHandler) {

	utf8Bytes := []byte(text)
	totalLength := len(utf8Bytes)

	runHandler.BeginLine()

	currentOffset := 0
	for currentOffset < totalLength {
		end := totalLength

		if fontEnd := fontIter.EndOfCurrentRun(); fontEnd < end {
			end = fontEnd
		}
		if bidiEnd := bidiIter.EndOfCurrentRun(); bidiEnd < end {
			end = bidiEnd
		}
		if scriptEnd := scriptIter.EndOfCurrentRun(); scriptEnd < end {
			end = scriptEnd
		}
		if langEnd := langIter.EndOfCurrentRun(); langEnd < end {
			end = langEnd
		}

		if end <= currentOffset {
			if end == totalLength {
				break
			}
			// Prevent infinite loop if iterators are broken
			log.Printf("Shaper iterator stuck at %d", currentOffset)
			break
		}

		currentFont := fontIter.CurrentFont()
		currentBidiLevel := bidiIter.CurrentLevel()
		currentScript := scriptIter.CurrentScript()
		currentLang := langIter.CurrentLanguage()

		s.shapeRun(text, currentOffset, end, currentFont, currentBidiLevel, currentScript, currentLang, features, width, runHandler)

		if fontIter.EndOfCurrentRun() == end {
			fontIter.Consume()
		}
		if bidiIter.EndOfCurrentRun() == end {
			bidiIter.Consume()
		}
		if scriptIter.EndOfCurrentRun() == end {
			scriptIter.Consume()
		}
		if langIter.EndOfCurrentRun() == end {
			langIter.Consume()
		}

		currentOffset = end
	}

	runHandler.CommitLine()
}

func (s *HarfbuzzShaper) shapeRun(text string, start, end int,
	skFont interfaces.SkFont, bidiLevel uint8, script uint32, lang string,
	features []Feature,
	width float32, runHandler RunHandler) {

	// 1. Resolve Face
	face := resolveFace(skFont)
	if face == nil {
		// Cannot shape without a face that supports go-text/typesetting
		log.Println("HarfbuzzShaper: typeface does not implement UseGoTextFace or returns nil")
		return
	}

	// 2. Prepare Input
	// Convert the run text to runes for HarfBuzz.
	// We assume text[start:end] are valid byte boundaries from the iterator.
	// If the text contains invalid UTF-8, string->[]rune will insert utf8.RuneError, which is acceptable.
	runText := []rune(text[start:end])

	textSize := skFont.Size()

	dir := di.DirectionLTR
	if bidiLevel%2 == 1 {
		dir = di.DirectionRTL
	}

	// Filter features used in this run
	var runFeatures []shaping.FontFeature
	for _, f := range features {
		// Intersection of [start, end) and [f.Start, f.End)
		// max(start, f.Start) < min(end, f.End)
		fStart := f.Start
		fEnd := f.End
		if fStart < start {
			fStart = start
		}
		if fEnd > end {
			fEnd = end
		}
		if fStart < fEnd {
			runFeatures = append(runFeatures, shaping.FontFeature{
				Tag:   font.Tag(f.Tag),
				Value: f.Value,
			})
		}
	}

	input := shaping.Input{
		Text:         runText,
		RunStart:     0,
		RunEnd:       len(runText),
		Direction:    dir,
		Face:         face,
		Size:         floatToFixed(float32(textSize)),
		Script:       language.Script(script),
		FontFeatures: runFeatures,
		// Language: language.NewLanguage(lang), (Future optimization: map language string)
	}

	// 3. Shape
	output := s.hb.Shape(input)

	// 4. Map to RunHandler
	count := len(output.Glyphs)
	if count == 0 {
		return
	}

	glyphs := make([]uint16, count)
	positions := make([]models.Point, count)
	clusters := make([]uint32, count)

	// We need to map clusters back to the original byte offset.
	// `output.Glyphs[i].ClusterIndex` is index in `runText` (rune index).
	// We need to convert rune index -> byte offset in `text`.
	// Need a helper to map rune index to byte offset.

	// Create map: runeIndex -> byteOffset relative to start
	runeToByte := make([]int, len(runText)+1)
	byteOff := 0
	for i, r := range runText {
		runeToByte[i] = byteOff
		byteOff += len(string(r))
	}
	runeToByte[len(runText)] = byteOff

	var currentX float32 = 0
	var currentY float32 = 0

	for i, g := range output.Glyphs {
		glyphs[i] = uint16(g.GlyphID)

		// Positions: accumulated advance + offset
		// Skia expects positions to be absolute coordinates of the glyph origin?
		// Or relative to the run?
		// RunBuffer documentation says: "Positions of the glyphs".
		// Usually (x,y).
		// Harfbuzz returns Advance (delta) and Offset (adjustment).

		// Skia standard:
		// pos[i] = (currentX, currentY) + offset
		// currentX += advance

		// Skia is Y-down, HarfBuzz is Y-up.
		// We must flip the Y components of offset and advance.
		// Effectively: skiaY = -hbY
		padX := fixedToFloat(g.XOffset)
		padY := -fixedToFloat(g.YOffset)

		positions[i] = models.Point{
			X: models.Scalar(currentX + padX),
			Y: models.Scalar(currentY + padY),
		}

		currentX += fixedToFloat(g.XAdvance)
		currentY += -fixedToFloat(g.YAdvance)

		// Clusters
		runeIdx := g.ClusterIndex
		if runeIdx < len(runeToByte) {
			clusters[i] = uint32(start + runeToByte[runeIdx])
		} else {
			clusters[i] = uint32(start + runeToByte[len(runeToByte)-1]) // End
		}
	}

	info := RunInfo{
		Font:       skFont,
		BidiLevel:  bidiLevel,
		Advance:    models.Point{X: models.Scalar(currentX), Y: models.Scalar(currentY)},
		GlyphCount: uint64(count),
		Utf8Range:  Range{Begin: start, End: end},
	}

	runHandler.RunInfo(info)
	buffer := runHandler.RunBuffer(info)

	copy(buffer.Glyphs, glyphs)
	copy(buffer.Positions, positions)
	copy(buffer.Clusters, clusters)

	runHandler.CommitRunBuffer(info)
}

func resolveFace(skFont interfaces.SkFont) *font.Face {
	tf := skFont.Typeface()
	if tf == nil {
		return nil
	}
	if exposing, ok := tf.(UseGoTextFace); ok {
		return exposing.GoTextFace()
	}
	return nil
}

func floatToFixed(f float32) fixed.Int26_6 {
	return fixed.Int26_6(f * 64)
}

func fixedToFloat(i fixed.Int26_6) float32 {
	return float32(i) / 64.0
}
