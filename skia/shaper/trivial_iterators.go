package shaper

import (
	"github.com/zodimo/go-skia-support/skia/interfaces"
)

// TrivialRunIterator is a base implementation for trivial iterators.
// It assumes a single run covering the entire text.
type TrivialRunIterator struct {
	textLength int
	atEnd      bool
}

// NewTrivialRunIterator creates a new TrivialRunIterator.
func NewTrivialRunIterator(textLength int) *TrivialRunIterator {
	return &TrivialRunIterator{
		textLength: textLength,
		atEnd:      textLength == 0,
	}
}

// Consume consumes the next n characters.
func (t *TrivialRunIterator) Consume() {
	t.atEnd = true
}

// EndOfCurrentRun returns the end index of the current run.
func (t *TrivialRunIterator) EndOfCurrentRun() int {
	return t.textLength
}

// AtEnd returns true if the iterator is at the end of the text.
func (t *TrivialRunIterator) AtEnd() bool {
	return t.atEnd
}

// TrivialFontRunIterator is a trivial implementation of FontRunIterator.
type TrivialFontRunIterator struct {
	*TrivialRunIterator
	font interfaces.SkFont
}

// NewTrivialFontRunIterator creates a new TrivialFontRunIterator.
func NewTrivialFontRunIterator(font interfaces.SkFont, textLength int) *TrivialFontRunIterator {
	return &TrivialFontRunIterator{
		TrivialRunIterator: NewTrivialRunIterator(textLength),
		font:               font,
	}
}

// CurrentFont returns the font for the current run.
func (t *TrivialFontRunIterator) CurrentFont() interfaces.SkFont {
	return t.font
}

// TrivialBiDiRunIterator is a trivial implementation of BiDiRunIterator.
type TrivialBiDiRunIterator struct {
	*TrivialRunIterator
	level uint8
}

// NewTrivialBiDiRunIterator creates a new TrivialBiDiRunIterator.
func NewTrivialBiDiRunIterator(bidiLevel uint8, textLength int) *TrivialBiDiRunIterator {
	return &TrivialBiDiRunIterator{
		TrivialRunIterator: NewTrivialRunIterator(textLength),
		level:              bidiLevel,
	}
}

// CurrentLevel returns the bidi level for the current run.
func (t *TrivialBiDiRunIterator) CurrentLevel() uint8 {
	return t.level
}

// TrivialScriptRunIterator is a trivial implementation of ScriptRunIterator.
type TrivialScriptRunIterator struct {
	*TrivialRunIterator
	script uint32
}

// NewTrivialScriptRunIterator creates a new TrivialScriptRunIterator.
func NewTrivialScriptRunIterator(script uint32, textLength int) *TrivialScriptRunIterator {
	return &TrivialScriptRunIterator{
		TrivialRunIterator: NewTrivialRunIterator(textLength),
		script:             script,
	}
}

// CurrentScript returns the script code for the current run.
func (t *TrivialScriptRunIterator) CurrentScript() uint32 {
	return t.script
}

// TrivialLanguageRunIterator is a trivial implementation of LanguageRunIterator.
type TrivialLanguageRunIterator struct {
	*TrivialRunIterator
	language string
}

// NewTrivialLanguageRunIterator creates a new TrivialLanguageRunIterator.
func NewTrivialLanguageRunIterator(language string, textLength int) *TrivialLanguageRunIterator {
	return &TrivialLanguageRunIterator{
		TrivialRunIterator: NewTrivialRunIterator(textLength),
		language:           language,
	}
}

// CurrentLanguage returns the language string for the current run.
func (t *TrivialLanguageRunIterator) CurrentLanguage() string {
	return t.language
}
