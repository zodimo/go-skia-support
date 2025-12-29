package shaper

import (
	"github.com/zodimo/go-skia-support/skia/interfaces"
)

// Shaper is the interface for text shaping.
// It maps to SkShaper in C++.
type Shaper interface {
	// Shape shapes the text using the font and runHandler.
	// leftToRight indicates the base direction of the text.
	// width is the width of the shape.
	Shape(text string, font interfaces.SkFont, leftToRight bool, width float32, runHandler RunHandler)
}

// RunIterator is the base interface for iterators over runs of text.
// It maps to SkShaper::RunIterator in C++.
type RunIterator interface {
	// Consume consumes the next n characters.
	Consume()
	// EndOfCurrentRun returns the end index of the current run.
	EndOfCurrentRun() int
	// AtEnd returns true if the iterator is at the end of the text.
	AtEnd() bool
}

// FontRunIterator iterates over runs of fonts.
// It maps to SkShaper::FontRunIterator in C++.
type FontRunIterator interface {
	RunIterator
	// CurrentFont returns the font for the current run.
	CurrentFont() interfaces.SkFont
}

// BiDiRunIterator iterates over runs of bidirectional levels.
// It maps to SkShaper::BiDiRunIterator in C++.
type BiDiRunIterator interface {
	RunIterator
	// CurrentLevel returns the bidi level for the current run.
	CurrentLevel() uint8
}

// ScriptRunIterator iterates over runs of scripts.
// It maps to SkShaper::ScriptRunIterator in C++.
type ScriptRunIterator interface {
	RunIterator
	// CurrentScript returns the script code for the current run.
	CurrentScript() uint32
}

// LanguageRunIterator iterates over runs of languages.
// It maps to SkShaper::LanguageRunIterator in C++.
type LanguageRunIterator interface {
	RunIterator
	// CurrentLanguage returns the language string for the current run.
	CurrentLanguage() string
}
