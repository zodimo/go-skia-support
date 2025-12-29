package shaper

import (
	"log"
	"math"

	"github.com/go-text/typesetting/di"
	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/language"
	"github.com/go-text/typesetting/segmenter"
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
type shapedRunData struct {
	info      RunInfo
	glyphs    []uint16
	positions []models.Point
	clusters  []uint32
}

// textProps holds safe-to-break point information from shaped model.
type textProps struct {
	glyphLen int     // Number of glyphs up to this point
	advance  float32 // Advance width up to this point
}

// lineBuilder accumulates runs for current line.
type lineBuilder struct {
	runs    []shapedRunData
	advance float32
}

// Shape shapes the text using the font and runHandler.
func (s *HarfbuzzShaper) Shape(text string, font interfaces.SkFont, leftToRight bool, width float32, runHandler RunHandler, features []Feature) {
	totalLength := len(text)

	fontIter := NewTrivialFontRunIterator(font, totalLength)
	bidiDir := uint8(0)
	if !leftToRight {
		bidiDir = 1
	}
	bidiIter := NewTrivialBiDiRunIterator(bidiDir, totalLength)
	scriptIter := NewTrivialScriptRunIterator(0, totalLength)
	langIter := NewTrivialLanguageRunIterator("en", totalLength)

	s.ShapeWithIterators(text, fontIter, bidiIter, scriptIter, langIter, features, width, runHandler)
}

// ShapeWithIterators shapes the text using custom iterators.
// When width > 0, implements shaper-driven line breaking following C++ ShaperDrivenWrapper.
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

	// No width constraint - shape and emit as single line
	if width <= 0 {
		s.shapeWithoutWrapping(text, fontIter, bidiIter, scriptIter, langIter, features, runHandler)
		return
	}

	// Get break opportunities for line breaking
	breakPoints := getLineBreakPoints(text)

	// Line accumulator
	var line lineBuilder

	currentOffset := 0
	for currentOffset < totalLength {
		// Find end of current item (smallest iterator boundary)
		itemEnd := totalLength
		if fontEnd := fontIter.EndOfCurrentRun(); fontEnd < itemEnd {
			itemEnd = fontEnd
		}
		if bidiEnd := bidiIter.EndOfCurrentRun(); bidiEnd < itemEnd {
			itemEnd = bidiEnd
		}
		if scriptEnd := scriptIter.EndOfCurrentRun(); scriptEnd < itemEnd {
			itemEnd = scriptEnd
		}
		if langEnd := langIter.EndOfCurrentRun(); langEnd < itemEnd {
			itemEnd = langEnd
		}

		// Split at feature boundaries
		for _, f := range features {
			if f.Start > currentOffset && f.Start < itemEnd {
				itemEnd = f.Start
			}
			if f.End > currentOffset && f.End < itemEnd {
				itemEnd = f.End
			}
		}

		if itemEnd <= currentOffset {
			break
		}

		currentFont := fontIter.CurrentFont()
		currentBidiLevel := bidiIter.CurrentLevel()
		currentScript := scriptIter.CurrentScript()
		currentLang := langIter.CurrentLanguage()

		// Shape the entire item as model
		model := s.shapeRunCollect(text, currentOffset, itemEnd, currentFont, currentBidiLevel, currentScript, currentLang, features)
		if model == nil {
			// Consume iterators and continue
			s.consumeIterators(fontIter, bidiIter, scriptIter, langIter, itemEnd)
			currentOffset = itemEnd
			continue
		}

		// Build textProps map for safe-to-break points
		propsMap := buildTextProps(model, currentOffset)

		// Process item with line breaking
		itemOffset := currentOffset
		for itemOffset < itemEnd {
			widthLeft := width - line.advance

			// Find best break point
			best, bestEnd := s.findBestBreak(text, itemOffset, itemEnd, breakPoints, propsMap,
				model, currentOffset, widthLeft,
				currentFont, currentBidiLevel, currentScript, currentLang, features)

			if best == nil {
				// No valid break found, force break at item end
				best = s.shapeRunCollect(text, itemOffset, itemEnd, currentFont, currentBidiLevel, currentScript, currentLang, features)
				bestEnd = itemEnd
			}

			// Check if best fits on current line
			if line.advance+float32(best.info.Advance.X) > width && len(line.runs) > 0 {
				// Emit current line and start new one
				s.emitLine(line.runs, runHandler)
				line = lineBuilder{}
			}

			// Add best to current line
			line.runs = append(line.runs, *best)
			line.advance += float32(best.info.Advance.X)

			// If we broke the item (didn't consume all), emit line
			if bestEnd < itemEnd {
				s.emitLine(line.runs, runHandler)
				line = lineBuilder{}
			}

			itemOffset = bestEnd
		}

		// Consume iterators
		s.consumeIterators(fontIter, bidiIter, scriptIter, langIter, itemEnd)
		currentOffset = itemEnd
	}

	// Emit final line
	if len(line.runs) > 0 {
		s.emitLine(line.runs, runHandler)
	} else if currentOffset == 0 {
		// Empty text case
		runHandler.BeginLine()
		runHandler.CommitRunInfo()
		runHandler.CommitLine()
	}
}

// findBestBreak finds the best break point using C++ scoring algorithm.
func (s *HarfbuzzShaper) findBestBreak(text string, itemOffset, itemEnd int,
	breakPoints []int, propsMap map[int]textProps,
	model *shapedRunData, modelStart int, widthLeft float32,
	currentFont interfaces.SkFont, bidiLevel uint8, script uint32, lang string,
	features []Feature) (*shapedRunData, int) {

	var best *shapedRunData
	bestEnd := itemOffset
	bestScore := float32(math.Inf(-1))

	// Get props at current offset for extraction
	startProps, hasStart := propsMap[itemOffset]
	if !hasStart {
		startProps = textProps{glyphLen: 0, advance: 0}
	}

	for _, bp := range breakPoints {
		if bp <= itemOffset || bp > itemEnd {
			continue
		}

		var candidate *shapedRunData
		var candidateAdvance float32

		// Try to extract from model using cached props
		if endProps, ok := propsMap[bp]; ok && hasStart {
			candidate = extractFromModel(model, startProps, endProps, itemOffset, bp)
			candidateAdvance = endProps.advance - startProps.advance
		} else {
			// Re-shape this segment
			candidate = s.shapeRunCollect(text, itemOffset, bp, currentFont, bidiLevel, script, lang, features)
			if candidate != nil {
				candidateAdvance = float32(candidate.info.Advance.X)
			}
		}

		if candidate == nil {
			continue
		}

		// Score: fits -> text length, doesn't fit -> negative overflow
		var score float32
		if candidateAdvance < widthLeft {
			score = float32(bp - itemOffset) // Maximize text length that fits
		} else {
			score = widthLeft - candidateAdvance // Negative means overflow
		}

		if score > bestScore {
			best = candidate
			bestEnd = bp
			bestScore = score
		}
	}

	// Also consider breaking at itemEnd
	if endProps, ok := propsMap[itemEnd]; ok && hasStart {
		candidate := extractFromModel(model, startProps, endProps, itemOffset, itemEnd)
		if candidate != nil {
			candidateAdvance := endProps.advance - startProps.advance
			var score float32
			if candidateAdvance < widthLeft {
				score = float32(itemEnd - itemOffset)
			} else {
				score = widthLeft - candidateAdvance
			}
			if score > bestScore {
				best = candidate
				bestEnd = itemEnd
			}
		}
	}

	return best, bestEnd
}

// buildTextProps builds a map of safe-to-break points from shaped model.
func buildTextProps(model *shapedRunData, modelStart int) map[int]textProps {
	props := make(map[int]textProps)
	var advance float32
	var prevCluster uint32 = 0

	for i, cluster := range model.clusters {
		// Safe to break when cluster changes
		if i == 0 || cluster != prevCluster {
			props[int(cluster)] = textProps{
				glyphLen: i,
				advance:  advance,
			}
			prevCluster = cluster
		}
		if i < len(model.positions) {
			// Accumulate advance from position differences or use info
			if i+1 < len(model.positions) {
				advance = float32(model.positions[i+1].X)
			} else {
				advance = float32(model.info.Advance.X)
			}
		}
	}

	// Always safe to break at end
	props[model.info.Utf8Range.End] = textProps{
		glyphLen: len(model.glyphs),
		advance:  float32(model.info.Advance.X),
	}

	return props
}

// extractFromModel extracts a subset of glyphs from model between two textProps.
func extractFromModel(model *shapedRunData, start, end textProps, byteStart, byteEnd int) *shapedRunData {
	if end.glyphLen <= start.glyphLen {
		return nil
	}

	glyphCount := end.glyphLen - start.glyphLen
	glyphs := make([]uint16, glyphCount)
	positions := make([]models.Point, glyphCount)
	clusters := make([]uint32, glyphCount)

	copy(glyphs, model.glyphs[start.glyphLen:end.glyphLen])
	copy(positions, model.positions[start.glyphLen:end.glyphLen])
	copy(clusters, model.clusters[start.glyphLen:end.glyphLen])

	// Adjust positions to be relative to start
	startX := float32(0)
	if start.glyphLen > 0 && start.glyphLen < len(model.positions) {
		startX = float32(model.positions[start.glyphLen].X)
	}
	for i := range positions {
		positions[i].X = models.Scalar(float32(positions[i].X) - startX)
	}

	return &shapedRunData{
		info: RunInfo{
			Font:       model.info.Font,
			BidiLevel:  model.info.BidiLevel,
			Script:     model.info.Script,
			Language:   model.info.Language,
			Advance:    models.Point{X: models.Scalar(end.advance - start.advance), Y: 0},
			GlyphCount: uint64(glyphCount),
			Utf8Range:  Range{Begin: byteStart, End: byteEnd},
		},
		glyphs:    glyphs,
		positions: positions,
		clusters:  clusters,
	}
}

// getLineBreakPoints returns byte offsets where lines can be broken.
func getLineBreakPoints(text string) []int {
	if len(text) == 0 {
		return nil
	}

	textRunes := []rune(text)
	var breaks []int
	var seg segmenter.Segmenter
	seg.Init(textRunes)
	iter := seg.LineIterator()

	// Build rune-to-byte offset mapping
	runeToByte := make([]int, len(textRunes)+1)
	byteOff := 0
	for i, r := range textRunes {
		runeToByte[i] = byteOff
		byteOff += len(string(r))
	}
	runeToByte[len(textRunes)] = byteOff

	// Collect break opportunities
	for iter.Next() {
		line := iter.Line()
		lineEnd := line.Offset + len(line.Text)
		if lineEnd <= len(textRunes) {
			breaks = append(breaks, runeToByte[lineEnd])
		}
	}

	return breaks
}

// shapeWithoutWrapping shapes text without line breaking (width <= 0).
func (s *HarfbuzzShaper) shapeWithoutWrapping(text string,
	fontIter FontRunIterator,
	bidiIter BiDiRunIterator,
	scriptIter ScriptRunIterator,
	langIter LanguageRunIterator,
	features []Feature,
	runHandler RunHandler) {

	utf8Bytes := []byte(text)
	totalLength := len(utf8Bytes)
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
		for _, f := range features {
			if f.Start > currentOffset && f.Start < end {
				end = f.Start
			}
			if f.End > currentOffset && f.End < end {
				end = f.End
			}
		}

		if end <= currentOffset {
			break
		}

		currentFont := fontIter.CurrentFont()
		currentBidiLevel := bidiIter.CurrentLevel()
		currentScript := scriptIter.CurrentScript()
		currentLang := langIter.CurrentLanguage()

		if runData := s.shapeRunCollect(text, currentOffset, end, currentFont, currentBidiLevel, currentScript, currentLang, features); runData != nil {
			shapedRuns = append(shapedRuns, *runData)
		}

		s.consumeIterators(fontIter, bidiIter, scriptIter, langIter, end)
		currentOffset = end
	}

	s.emitLine(shapedRuns, runHandler)
}

// consumeIterators advances iterators that end at the given position.
func (s *HarfbuzzShaper) consumeIterators(fontIter FontRunIterator, bidiIter BiDiRunIterator,
	scriptIter ScriptRunIterator, langIter LanguageRunIterator, end int) {
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
}

// emitLine emits callbacks for a single line of runs with visual reordering.
func (s *HarfbuzzShaper) emitLine(runs []shapedRunData, runHandler RunHandler) {
	numRuns := len(runs)
	if numRuns == 0 {
		runHandler.BeginLine()
		runHandler.CommitRunInfo()
		runHandler.CommitLine()
		return
	}

	// Compute visual order
	levels := make([]uint8, numRuns)
	for i, run := range runs {
		levels[i] = run.info.BidiLevel
	}
	visualOrder := reorderVisual(levels)

	runHandler.BeginLine()

	for i := 0; i < numRuns; i++ {
		logicalIndex := visualOrder[i]
		runHandler.RunInfo(runs[logicalIndex].info)
	}
	runHandler.CommitRunInfo()

	for i := 0; i < numRuns; i++ {
		logicalIndex := visualOrder[i]
		run := runs[logicalIndex]
		buffer := runHandler.RunBuffer(run.info)
		copy(buffer.Glyphs, run.glyphs)
		copy(buffer.Positions, run.positions)
		copy(buffer.Clusters, run.clusters)
		runHandler.CommitRunBuffer(run.info)
	}

	runHandler.CommitLine()
}

// reorderVisual computes the visual order of runs based on BiDi levels.
func reorderVisual(levels []uint8) []int {
	n := len(levels)
	if n == 0 {
		return nil
	}

	order := make([]int, n)
	for i := range order {
		order[i] = i
	}

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

	if lowestOddLevel == 255 {
		return order
	}

	for level := highestLevel; level >= lowestOddLevel; level-- {
		start := -1
		for i := 0; i <= n; i++ {
			if i < n && levels[order[i]] >= level {
				if start == -1 {
					start = i
				}
			} else {
				if start != -1 {
					reverseSlice(order, start, i-1)
					start = -1
				}
			}
		}
	}

	return order
}

func reverseSlice(slice []int, start, end int) {
	for start < end {
		slice[start], slice[end] = slice[end], slice[start]
		start++
		end--
	}
}

// shapeRunCollect shapes a run and returns the data.
func (s *HarfbuzzShaper) shapeRunCollect(text string, start, end int,
	skFont interfaces.SkFont, bidiLevel uint8, script uint32, lang string,
	features []Feature) *shapedRunData {

	face := resolveFace(skFont)
	if face == nil {
		log.Println("HarfbuzzShaper: typeface does not implement UseGoTextFace or returns nil")
		return nil
	}

	fullTextRunes := []rune(text)

	// Map byte offsets to rune indices
	byteToRuneStart := 0
	byteToRuneEnd := 0
	currentByte := 0
	for i, r := range fullTextRunes {
		if currentByte == start {
			byteToRuneStart = i
		}
		if currentByte == end {
			byteToRuneEnd = i
			break
		}
		currentByte += len(string(r))
	}
	if currentByte == end {
		byteToRuneEnd = len(fullTextRunes)
	}
	if currentByte < start {
		byteToRuneStart = len(fullTextRunes)
	}

	if byteToRuneEnd <= byteToRuneStart {
		return nil
	}

	textSize := skFont.Size()

	dir := di.DirectionLTR
	if bidiLevel%2 == 1 {
		dir = di.DirectionRTL
	}

	var runFeatures []shaping.FontFeature
	for _, f := range features {
		if f.Start <= start && f.End >= end {
			runFeatures = append(runFeatures, shaping.FontFeature{
				Tag:   font.Tag(f.Tag),
				Value: f.Value,
			})
		}
	}

	input := shaping.Input{
		Text:         fullTextRunes,
		RunStart:     byteToRuneStart,
		RunEnd:       byteToRuneEnd,
		Direction:    dir,
		Face:         face,
		Size:         floatToFixed(float32(textSize)),
		Script:       language.Script(script),
		FontFeatures: runFeatures,
		Language:     language.NewLanguage(lang),
	}

	output := s.hb.Shape(input)

	if len(output.Glyphs) == 0 {
		return nil
	}

	count := len(output.Glyphs)
	glyphs := make([]uint16, count)
	positions := make([]models.Point, count)
	clusters := make([]uint32, count)

	runeToByte := make([]int, len(fullTextRunes)+1)
	byteOff := 0
	for i, r := range fullTextRunes {
		runeToByte[i] = byteOff
		byteOff += len(string(r))
	}
	runeToByte[len(fullTextRunes)] = byteOff

	var currentX float32 = 0
	var currentY float32 = 0

	for i, g := range output.Glyphs {
		glyphs[i] = uint16(g.GlyphID)

		padX := fixedToFloat(g.XOffset)
		padY := -fixedToFloat(g.YOffset)

		positions[i] = models.Point{
			X: models.Scalar(currentX + padX),
			Y: models.Scalar(currentY + padY),
		}

		currentX += fixedToFloat(g.XAdvance)
		currentY += -fixedToFloat(g.YAdvance)

		runeIdx := g.ClusterIndex
		if runeIdx < len(runeToByte) {
			clusters[i] = uint32(runeToByte[runeIdx])
		} else {
			clusters[i] = uint32(runeToByte[len(runeToByte)-1])
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
