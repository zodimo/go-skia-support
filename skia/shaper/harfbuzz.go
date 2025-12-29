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

// shapedRunData holds the result of shaping a single run.
// This allows us to collect all runs before emitting callbacks in the correct order.
type shapedRunData struct {
	info      RunInfo
	glyphs    []uint16
	positions []models.Point
	clusters  []uint32
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
// This implementation follows the C++ callback ordering:
// 1. BeginLine()
// 2. RunInfo() for ALL runs (in visual order)
// 3. CommitRunInfo() once
// 4. RunBuffer() + CommitRunBuffer() for ALL runs (in visual order)
// 5. CommitLine()
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

	// Phase 1: Collect all shaped runs (in logical order)
	var shapedRuns []shapedRunData

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

		// Split at feature boundaries
		for _, f := range features {
			if f.Start > currentOffset && f.Start < end {
				end = f.Start
			}
			if f.End > currentOffset && f.End < end {
				end = f.End
			}
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

		// Shape the run and collect the data
		if runData := s.shapeRunCollect(text, currentOffset, end, currentFont, currentBidiLevel, currentScript, currentLang, features); runData != nil {
			shapedRuns = append(shapedRuns, *runData)
		}

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

	// Phase 2: Compute visual ordering for BiDi text
	numRuns := len(shapedRuns)
	visualOrder := make([]int, numRuns)

	if numRuns > 0 {
		// Collect bidi levels
		levels := make([]uint8, numRuns)
		for i, run := range shapedRuns {
			levels[i] = run.info.BidiLevel
		}
		// Compute visual-to-logical mapping
		visualOrder = reorderVisual(levels)
	}

	// Phase 3: Emit callbacks in the correct order (visual order, matching C++ behavior)
	runHandler.BeginLine()

	// 3a. Call RunInfo for ALL runs in visual order
	for i := 0; i < numRuns; i++ {
		logicalIndex := visualOrder[i]
		runHandler.RunInfo(shapedRuns[logicalIndex].info)
	}

	// 3b. CommitRunInfo once (after all RunInfo calls)
	runHandler.CommitRunInfo()

	// 3c. Call RunBuffer + CommitRunBuffer for ALL runs in visual order
	for i := 0; i < numRuns; i++ {
		logicalIndex := visualOrder[i]
		run := shapedRuns[logicalIndex]

		buffer := runHandler.RunBuffer(run.info)

		copy(buffer.Glyphs, run.glyphs)
		copy(buffer.Positions, run.positions)
		copy(buffer.Clusters, run.clusters)

		runHandler.CommitRunBuffer(run.info)
	}

	// 3d. CommitLine
	runHandler.CommitLine()
}

// reorderVisual computes the visual order of runs based on their BiDi levels.
// It returns a slice where visualOrder[visualIndex] = logicalIndex.
// This implements the Unicode Bidirectional Algorithm L2 rule for reordering.
//
// The algorithm:
// 1. Find the highest level among all runs
// 2. For each level from highest down to the lowest odd level:
//   - Reverse any contiguous sequence of runs at that level or higher
func reorderVisual(levels []uint8) []int {
	n := len(levels)
	if n == 0 {
		return nil
	}

	// Initialize visual order as identity (logical order)
	order := make([]int, n)
	for i := range order {
		order[i] = i
	}

	// Find highest and lowest odd levels
	var highestLevel uint8 = 0
	var lowestOddLevel uint8 = 255

	for _, level := range levels {
		if level > highestLevel {
			highestLevel = level
		}
		if level%2 == 1 && level < lowestOddLevel {
			lowestOddLevel = level
		}
	}

	// If no odd levels, all text is LTR, no reordering needed
	if lowestOddLevel == 255 {
		return order
	}

	// Apply L2: reverse runs at each level from highest to lowestOddLevel
	for level := highestLevel; level >= lowestOddLevel; level-- {
		// Find and reverse contiguous sequences at this level or higher
		start := -1
		for i := 0; i <= n; i++ {
			if i < n && levels[order[i]] >= level {
				if start == -1 {
					start = i
				}
			} else {
				if start != -1 {
					// Reverse the sequence from start to i-1
					reverseSlice(order, start, i-1)
					start = -1
				}
			}
		}
	}

	return order
}

// reverseSlice reverses elements in slice from index start to end (inclusive).
func reverseSlice(slice []int, start, end int) {
	for start < end {
		slice[start], slice[end] = slice[end], slice[start]
		start++
		end--
	}
}

// shapeRunCollect shapes a run and returns the data without calling RunHandler callbacks.
// Returns nil if the run produces no glyphs.
func (s *HarfbuzzShaper) shapeRunCollect(text string, start, end int,
	skFont interfaces.SkFont, bidiLevel uint8, script uint32, lang string,
	features []Feature) *shapedRunData {

	// 1. Resolve Face
	face := resolveFace(skFont)
	if face == nil {
		// Cannot shape without a face that supports go-text/typesetting
		log.Println("HarfbuzzShaper: typeface does not implement UseGoTextFace or returns nil")
		return nil
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
		// Since we segmented by feature boundaries, a feature applies if it covers the entire segment.
		if f.Start <= start && f.End >= end {
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
		Language:     language.NewLanguage(lang),
	}

	// 3. Shape
	output := s.hb.Shape(input)

	// 4. Map to output data
	count := len(output.Glyphs)
	if count == 0 {
		return nil
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

	return &shapedRunData{
		info: RunInfo{
			Font:       skFont,
			BidiLevel:  bidiLevel,
			Script:     script,
			Language:   lang,
			Advance:    models.Point{X: models.Scalar(currentX), Y: models.Scalar(currentY)},
			GlyphCount: uint64(count),
			Utf8Range:  Range{Begin: start, End: end},
		},
		glyphs:    glyphs,
		positions: positions,
		clusters:  clusters,
	}
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
