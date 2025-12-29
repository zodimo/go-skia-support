package shaper

import (
	"unicode/utf8"

	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// PrimitiveShaper is a basic implementation of the Shaper interface.
// It orchestrates the shaping process by coordinating property iterators
// and invoking a run shaper on consistent runs of text.
type PrimitiveShaper struct {
	// FontMgr interfaces.SkFontMgr // Future story: will need this for fallback
}

// NewPrimitiveShaper creates a new instance of PrimitiveShaper.
func NewPrimitiveShaper() *PrimitiveShaper {
	return &PrimitiveShaper{}
}

// Shape shapes the text using the font and runHandler.
// It implements the Shaper interface.
func (ps *PrimitiveShaper) Shape(text string, font interfaces.SkFont, leftToRight bool, width float32, runHandler RunHandler, features []Feature) {
	// 1. Create trivial iterators if necessary.
	// In C++ SkShaper::Shape (the simple one) creates a TrivialFontRunIterator,
	// TrivialBiDiRunIterator, TrivialScriptRunIterator, TrivialLanguageRunIterator.
	// We will assume for this story that we can use Trivial iterators if we only have the font.

	utf8Bytes := []byte(text)
	totalLength := len(utf8Bytes)

	// TODO: When we have the specific iterator factories or structs from Story 2 (Trivial Iterators),
	// we will instantiate them here. For now, I will assume they exist or implement ad-hoc ones
	// if they are not yet available.
	// Looking at the file list, `trivial_iterators.go` exists.

	fontIter := NewTrivialFontRunIterator(font, totalLength)
	bidiIter := NewTrivialBiDiRunIterator(0, totalLength) // 0 for LTR?
	if !leftToRight {
		bidiIter = NewTrivialBiDiRunIterator(1, totalLength) // 1 for RTL?
	}
	// Correcting signatures based on trivial_iterators.go:
	// NewTrivialScriptRunIterator(script uint32, textLength int)
	// NewTrivialLanguageRunIterator(language string, textLength int)
	scriptIter := NewTrivialScriptRunIterator(0, totalLength)    // 0 as default script?
	langIter := NewTrivialLanguageRunIterator("en", totalLength) // "en" as default lang?

	ps.ShapeWithIterators(text, fontIter, bidiIter, scriptIter, langIter, features, width, runHandler)
}

// ShapeWithIterators allows shaping with custom iterators.
// This is effectively the "Shape" method from the contract that takes iterators.
func (ps *PrimitiveShaper) ShapeWithIterators(text string,
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
		// End of the current run is the minimum of all iterator ends
		end := totalLength

		fontEnd := fontIter.EndOfCurrentRun()
		if fontEnd < end {
			end = fontEnd
		}

		bidiEnd := bidiIter.EndOfCurrentRun()
		if bidiEnd < end {
			end = bidiEnd
		}

		scriptEnd := scriptIter.EndOfCurrentRun()
		if scriptEnd < end {
			end = scriptEnd
		}

		langEnd := langIter.EndOfCurrentRun()
		if langEnd < end {
			end = langEnd
		}

		// Ensure we are making progress
		if end <= currentOffset {
			// This should panic or handle error in production code, but for now we clamp or break to avoid infinite loops
			if end == totalLength {
				break
			}
			// If iterators are broken, force advance?
			// C++ SkShaper loop logic is robust.
			// For now, let's assume iterators are well-behaved.
			// If end <= currentOffset, it implies an iterator didn't advance or we are stuck.
			// We'll just break to be safe against infinite loops during dev.
			break
		}

		// Identify current properties
		currentFont := fontIter.CurrentFont()
		currentBidiLevel := bidiIter.CurrentLevel()
		currentScript := scriptIter.CurrentScript()
		currentLang := langIter.CurrentLanguage()

		// Shape this specific run
		// We pass the slice of text relevant to this run, or the indices?
		// Typically shaper needs the whole context but acts on a range.
		// For simplicity in this story (Story 1 is just the loop), we won't implement the complex
		// context awareness yet.
		ps.shapeRun(text, currentOffset, end, currentFont, currentBidiLevel, currentScript, currentLang, width, runHandler)

		// Advance logic
		// Advance each iterator if its run has ended
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

	if currentOffset == 0 && totalLength == 0 {
		runHandler.BeginLine()
		runHandler.CommitRunInfo()
		runHandler.CommitLine()
	}

	runHandler.CommitLine()
}

// shapeRun implements the primitive shaping logic including word-wrapping.
func (ps *PrimitiveShaper) shapeRun(text string, start, end int,
	font interfaces.SkFont, bidiLevel uint8, script uint32, lang string,
	width float32, runHandler RunHandler) {

	runText := text[start:end]
	runes := []rune(runText)
	count := len(runes)

	if count == 0 {
		return
	}

	glyphs := make([]uint16, count)
	for i, r := range runes {
		glyphs[i] = font.UnicharToGlyph(r)
	}

	widths := font.GetWidths(glyphs)

	glyphOffset := 0
	utf8Offset := start

	for utf8Offset < end {
		// Line break logic
		bytesConsumed, bytesCollapsed := linebreak(runText[utf8Offset-start:], font, width, widths[glyphOffset:])
		bytesVisible := bytesConsumed - bytesCollapsed

		visibleRunes := []rune(runText[utf8Offset-start : utf8Offset-start+bytesVisible])
		numGlyphs := len(visibleRunes)

		// Calculate run advance
		var runAdvanceX float32 = 0
		for i := 0; i < numGlyphs; i++ {
			runAdvanceX += float32(widths[glyphOffset+i])
		}

		info := RunInfo{
			Font:       font,
			BidiLevel:  bidiLevel,
			Script:     script,
			Language:   lang,
			Advance:    models.Point{X: models.Scalar(runAdvanceX), Y: 0},
			GlyphCount: uint64(numGlyphs),
			Utf8Range:  Range{Begin: utf8Offset, End: utf8Offset + bytesVisible},
		}

		runHandler.BeginLine()
		if info.GlyphCount > 0 {
			runHandler.RunInfo(info)
		}
		runHandler.CommitRunInfo()

		if info.GlyphCount > 0 {
			buffer := runHandler.RunBuffer(info)
			copy(buffer.Glyphs, glyphs[glyphOffset:glyphOffset+numGlyphs])

			var currentX float32 = 0
			byteOff := utf8Offset
			for i := 0; i < numGlyphs; i++ {
				buffer.Positions[i] = models.Point{X: models.Scalar(currentX), Y: 0}
				buffer.Clusters[i] = uint32(byteOff)
				currentX += float32(widths[glyphOffset+i])
				byteOff += utf8.RuneLen(visibleRunes[i])
			}
			runHandler.CommitRunBuffer(info)
		}
		runHandler.CommitLine()

		utf8Offset += bytesConsumed
		// Advance glyphOffset by number of runes in consumed bytes
		consumedRunes := []rune(runText[utf8Offset-start-bytesConsumed : utf8Offset-start])
		glyphOffset += len(consumedRunes)
	}
}

func isBreakingWhitespace(r rune) bool {
	switch r {
	case 0x0020, 0x1680, 0x180E, 0x2000, 0x2001, 0x2002, 0x2003, 0x2004, 0x2005, 0x2006, 0x2007, 0x2008, 0x2009, 0x200A, 0x200B, 0x202F, 0x205F, 0x3000:
		return true
	default:
		return false
	}
}

func linebreak(text string, font interfaces.SkFont, width float32, advances []interfaces.Scalar) (int, int) {
	var accumulatedWidth float32 = 0
	glyphIndex := 0
	start := 0
	wordStart := 0
	prevWS := true
	textBytes := []byte(text)
	stop := len(textBytes)
	curr := 0

	for curr < stop {
		prevText := curr
		r, size := utf8.DecodeRune(textBytes[curr:])
		curr += size

		accumulatedWidth += float32(advances[glyphIndex])
		glyphIndex++
		currWS := isBreakingWhitespace(r)

		if !currWS && prevWS {
			wordStart = prevText
		}
		prevWS = currWS

		if width < accumulatedWidth {
			consumeWhitespace := false
			if currWS {
				if prevText == start {
					prevText = curr
				}
				consumeWhitespace = true
			} else if wordStart != start {
				curr = wordStart
			} else if prevText > start {
				curr = prevText
			} else {
				prevText = curr
				consumeWhitespace = true
			}

			var trailing int
			if consumeWhitespace {
				next := curr
				for next < stop {
					rn, sz := utf8.DecodeRune(textBytes[next:])
					if !isBreakingWhitespace(rn) {
						break
					}
					next += sz
				}
				trailing = next - prevText
				curr = next
			}
			return curr, trailing
		}
	}

	return len(text), 0
}

// Helper utility for min (not strictly needed since we expanded it above for clarity)
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
