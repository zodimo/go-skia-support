package impl

import (
	"unicode"
	"unicode/utf8"

	"github.com/zodimo/go-skia-support/skia/interfaces"
)

// SkUnicodeImpl implements interfaces.SkUnicode.
type SkUnicodeImpl struct{}

// NewSkUnicode creates a new instance of SkUnicodeImpl.
func NewSkUnicode() interfaces.SkUnicode {
	return &SkUnicodeImpl{}
}

// FindPreviousGraphemeBoundary finds the start of the grapheme cluster containing the offset.
// This is a simplified implementation that keeps combining marks with their base character.
func (u *SkUnicodeImpl) FindPreviousGraphemeBoundary(text string, offset int) int {
	if offset <= 0 || offset > len(text) {
		return 0
	}

	// Move back to find the rune boundary if offset is mid-rune
	start := offset
	for !utf8.RuneStart(text[start]) {
		start--
		if start < 0 {
			return 0
		}
	}

	// Now we are at a rune boundary `start`.
	// We want to verify if the rune at `start` is an extension of a previous one.
	// We need to look at the rune *before* this one if we are conceptually "at" this character
	// but strictly `offset` usually points AFTER the character we just processed or AT the character we are about to?
	// C++ usage: `ci` is cluster index. `findPreviousGraphemeBoundary(ci)` typically returns `ci` if it IS a boundary,
	// or the start of the grapheme if `ci` is inside one.
	//
	// If the user asks for boundary at `offset`, check if `text[offset]` starts a new grapheme.
	// If yes, return `offset`. If no (it's a combining mark), return start of previous rune (and recurse).

	// Let's iterate backwards from `start`.
	// 1. Get current rune at `start`? No, we need to check if `start` itself is a boundary.
	//    A position is a boundary if the character *starting* there is a Base, OR if it's start of text.

	// Helper to get rune at index
	r, _ := utf8.DecodeRuneInString(text[start:])
	if r == utf8.RuneError {
		return start
	}

	if u.isGraphemeBreak(r) {
		// It is a start of a grapheme (Base char)
		// So the boundary is here.
		return start
	}

	// It is NOT a grapheme break (it is a combining mark).
	// We need to find the base.
	curr := start
	for curr > 0 {
		_, size := utf8.DecodeLastRuneInString(text[:curr])
		curr -= size

		// Check the rune we just stepped over
		r, _ := utf8.DecodeRuneInString(text[curr:])
		if u.isGraphemeBreak(r) {
			return curr
		}
	}
	return 0
}

func (u *SkUnicodeImpl) isGraphemeBreak(r rune) bool {
	// Simplified: A rune breaks a grapheme if it is NOT a combining mark.
	// Categories: Mn, Mc, Me are combining.
	// Also ZWJ (U+200D) extends?
	return !unicode.Is(unicode.Mn, r) &&
		!unicode.Is(unicode.Mc, r) &&
		!unicode.Is(unicode.Me, r) &&
		r != 0x200D // Zero Width Joiner
}

// IsEmoji returns true if the rune is an emoji.
func (u *SkUnicodeImpl) IsEmoji(r rune) bool {
	// Check against Unicode Emoji property (available in Go 1.21+ unicode package potentially)
	// Fallback to ranges
	return (r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols and Pictographs
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map
		(r >= 0x1F1E0 && r <= 0x1F1FF) || // Regional Indicators
		(r >= 0x2600 && r <= 0x26FF) || // Misc Symbols
		(r >= 0x2700 && r <= 0x27BF) // Dingbats
}

// IsEmojiComponent returns true if the rune is an emoji component.
func (u *SkUnicodeImpl) IsEmojiComponent(r rune) bool {
	// E.g. Modifier Fitzpatrick
	return r >= 0x1F3FB && r <= 0x1F3FF
}

// IsRegionalIndicator returns true if the rune is a regional indicator.
func (u *SkUnicodeImpl) IsRegionalIndicator(r rune) bool {
	return r >= 0x1F1E6 && r <= 0x1F1FF
}

// CodeUnitHasProperty verifies if a code unit has a flag.
func (u *SkUnicodeImpl) CodeUnitHasProperty(text string, offset int, property interfaces.CodeUnitFlags) bool {
	if offset < 0 || offset >= len(text) {
		return false
	}

	// Ensure we look at the rune starting at or containing offset
	// For simplicity assuming offset is rune start
	r, _ := utf8.DecodeRuneInString(text[offset:])

	if property&interfaces.CodeUnitFlagControl != 0 {
		if unicode.IsControl(r) {
			return true
		}
	}

	if property&interfaces.CodeUnitFlagGraphemeStart != 0 {
		if u.isGraphemeBreak(r) {
			return true
		}
	}

	if property&interfaces.CodeUnitFlagPartOfWhitespace != 0 {
		if unicode.IsSpace(r) {
			return true
		}
	}

	return false
}
