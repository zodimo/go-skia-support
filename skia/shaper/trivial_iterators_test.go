package shaper

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/impl"
)

func TestTrivialRunIterator_Lifecycle(t *testing.T) {
	textLength := 10
	iter := NewTrivialRunIterator(textLength)

	if iter.AtEnd() {
		t.Error("Iterator should not be at end initially")
	}
	if iter.EndOfCurrentRun() != textLength {
		t.Errorf("EndOfCurrentRun should match text length: expected %d, got %d", textLength, iter.EndOfCurrentRun())
	}

	iter.Consume()
	if !iter.AtEnd() {
		t.Error("Iterator should be at end after Consume")
	}
}

func TestTrivialRunIterator_Empty(t *testing.T) {
	iter := NewTrivialRunIterator(0)
	if !iter.AtEnd() {
		t.Error("Iterator should be at end for empty text")
	}
}

func TestTrivialFontRunIterator(t *testing.T) {
	textLength := 5
	font := impl.NewFont()
	iter := NewTrivialFontRunIterator(font, textLength)

	if iter.AtEnd() {
		t.Error("Expected not at end")
	}
	if iter.CurrentFont() != font {
		t.Error("Font mismatch")
	}

	iter.Consume()
	if !iter.AtEnd() {
		t.Error("Iterator should be at end after Consume")
	}
}

func TestTrivialBiDiRunIterator(t *testing.T) {
	textLength := 5
	level := uint8(1)
	iter := NewTrivialBiDiRunIterator(level, textLength)

	if iter.AtEnd() {
		t.Error("Expected not at end")
	}
	if iter.CurrentLevel() != level {
		t.Errorf("Level mismatch: expected %d, got %d", level, iter.CurrentLevel())
	}

	iter.Consume()
	if !iter.AtEnd() {
		t.Error("Iterator should be at end after Consume")
	}
}

func TestTrivialScriptRunIterator(t *testing.T) {
	textLength := 5
	script := uint32(123)
	iter := NewTrivialScriptRunIterator(script, textLength)

	if iter.AtEnd() {
		t.Error("Expected not at end")
	}
	if iter.CurrentScript() != script {
		t.Errorf("Script mismatch: expected %d, got %d", script, iter.CurrentScript())
	}

	iter.Consume()
	if !iter.AtEnd() {
		t.Error("Iterator should be at end after Consume")
	}
}

func TestTrivialLanguageRunIterator(t *testing.T) {
	textLength := 5
	lang := "en-US"
	iter := NewTrivialLanguageRunIterator(lang, textLength)

	if iter.AtEnd() {
		t.Error("Expected not at end")
	}
	if iter.CurrentLanguage() != lang {
		t.Errorf("Language mismatch: expected %s, got %s", lang, iter.CurrentLanguage())
	}

	iter.Consume()
	if !iter.AtEnd() {
		t.Error("Iterator should be at end after Consume")
	}
}
