package shaper

import (
	"github.com/zodimo/go-skia-support/skia/interfaces"
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
func (ps *PrimitiveShaper) Shape(text string, font interfaces.SkFont, leftToRight bool, width float32, runHandler RunHandler) {
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

	ps.ShapeWithIterators(text, fontIter, bidiIter, scriptIter, langIter, width, runHandler)
}

// ShapeWithIterators allows shaping with custom iterators.
// This is effectively the "Shape" method from the contract that takes iterators.
func (ps *PrimitiveShaper) ShapeWithIterators(text string,
	fontIter FontRunIterator,
	bidiIter BiDiRunIterator,
	scriptIter ScriptRunIterator,
	langIter LanguageRunIterator,
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

	runHandler.CommitLine()
}

// shapeRun is a placeholder for the actual shaping logic (Story 2+).
func (ps *PrimitiveShaper) shapeRun(text string, start, end int,
	font interfaces.SkFont, bidiLevel uint8, script uint32, lang string,
	width float32, runHandler RunHandler) {

	// Placeholder for Story 1: report the run range so we can test the loop logic.
	info := RunInfo{
		Font:      font,
		BidiLevel: bidiLevel,
		Utf8Range: Range{Begin: start, End: end},
	}
	runHandler.RunInfo(info)
}

// Helper utility for min (not strictly needed since we expanded it above for clarity)
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
