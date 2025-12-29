package shaper

import (
	"unicode/utf8"

	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// FontMgrRunIterator is a FontRunIterator that uses SkFontMgr for font fallback.
// It iterates through text, using the initial font when possible, and falling back
// to fonts from the font manager when characters are not supported.
//
// Ported from: skia-source/modules/skshaper/src/SkShaper.cpp (FontMgrRunIterator class)
type FontMgrRunIterator struct {
	text         string
	current      int // current position in text (bytes)
	end          int // end of text (bytes)
	fallbackMgr  interfaces.SkFontMgr
	font         interfaces.SkFont   // initial font
	fallbackFont interfaces.SkFont   // current fallback font (may be nil)
	currentFont  interfaces.SkFont   // pointer to current font in use
	requestName  string              // optional family name for fallback requests
	requestStyle models.FontStyle    // style for fallback requests
	language     LanguageRunIterator // optional language iterator
}

// MakeFontMgrRunIterator creates a FontRunIterator that uses the given font manager
// to find fallback fonts for characters not supported by the initial font.
//
// Parameters:
//   - text: the UTF-8 text to iterate over
//   - font: the initial font to use
//   - fallbackMgr: the font manager for finding fallback fonts
//
// Returns a FontRunIterator that produces font runs based on character support.
//
// Ported from: SkShaper::MakeFontMgrRunIterator
func MakeFontMgrRunIterator(text string, font interfaces.SkFont, fallbackMgr interfaces.SkFontMgr) FontRunIterator {
	if fallbackMgr == nil {
		// If no fallback manager, use trivial iterator
		return NewTrivialFontRunIterator(font, len(text))
	}

	style := models.FontStyleNormal()
	if font.Typeface() != nil {
		style = font.Typeface().FontStyle()
	}

	return &FontMgrRunIterator{
		text:         text,
		current:      0,
		end:          len(text),
		fallbackMgr:  fallbackMgr,
		font:         font,
		fallbackFont: nil,
		currentFont:  nil,
		requestName:  "",
		requestStyle: style,
		language:     nil,
	}
}

// MakeFontMgrRunIteratorWithOptions creates a FontRunIterator with additional options
// for controlling fallback behavior.
//
// Parameters:
//   - text: the UTF-8 text to iterate over
//   - font: the initial font to use
//   - fallbackMgr: the font manager for finding fallback fonts
//   - requestName: family name to request for fallbacks (can be empty)
//   - requestStyle: font style to request for fallbacks
//   - language: optional language run iterator for locale-aware fallback
//
// Ported from: SkShaper::MakeFontMgrRunIterator (overload with options)
func MakeFontMgrRunIteratorWithOptions(
	text string,
	font interfaces.SkFont,
	fallbackMgr interfaces.SkFontMgr,
	requestName string,
	requestStyle models.FontStyle,
	language LanguageRunIterator,
) FontRunIterator {
	if fallbackMgr == nil {
		return NewTrivialFontRunIterator(font, len(text))
	}

	return &FontMgrRunIterator{
		text:         text,
		current:      0,
		end:          len(text),
		fallbackMgr:  fallbackMgr,
		font:         font,
		fallbackFont: nil,
		currentFont:  nil,
		requestName:  requestName,
		requestStyle: requestStyle,
		language:     language,
	}
}

// Consume advances the iterator to find the next font run.
// It consumes characters until a different font is needed.
//
// Ported from: FontMgrRunIterator::consume()
func (iter *FontMgrRunIterator) Consume() {
	if iter.current >= iter.end {
		return
	}

	// Get the first character
	r, size := utf8.DecodeRuneInString(iter.text[iter.current:])
	iter.current += size

	// Determine which font to use for this character
	if iter.font.UnicharToGlyph(r) != 0 {
		// Initial font can handle this character
		iter.currentFont = iter.font
	} else if iter.fallbackFont != nil && iter.fallbackFont.UnicharToGlyph(r) != 0 {
		// Current fallback font can handle it
		iter.currentFont = iter.fallbackFont
	} else {
		// Need to find a fallback font
		iter.currentFont = iter.findFallbackFont(r)
	}

	// Continue consuming characters that use the same font
	for iter.current < iter.end {
		r, size = utf8.DecodeRuneInString(iter.text[iter.current:])

		// If we're using a fallback and the initial font can handle this char, end run
		if iter.currentFont != iter.font && iter.font.UnicharToGlyph(r) != 0 {
			return
		}

		// If current font can't handle this char, check if another font could
		if iter.currentFont.UnicharToGlyph(r) == 0 {
			// Try to find a different font
			candidate := iter.tryFindFallback(r)
			if candidate != nil {
				// Found a different font, end the current run
				return
			}
		}

		iter.current += size
	}
}

// findFallbackFont attempts to find a fallback font for the given character.
// Returns the fallback font or the initial font if no fallback found.
func (iter *FontMgrRunIterator) findFallbackFont(r rune) interfaces.SkFont {
	// Get language if available
	var bcp47 []string
	if iter.language != nil && !iter.language.AtEnd() {
		bcp47 = []string{iter.language.CurrentLanguage()}
	}

	// Ask the font manager for a fallback
	fallbackTypeface := iter.fallbackMgr.MatchFamilyStyleCharacter(
		iter.requestName,
		iter.requestStyle,
		bcp47,
		r,
	)

	if fallbackTypeface != nil {
		// Create a font with the fallback typeface
		// Note: In a real implementation, we would clone the font properties
		// For now, we just use the fallback typeface
		iter.fallbackFont = iter.font // Use same font (stub behavior)
		return iter.fallbackFont
	}

	// No fallback found, use initial font
	return iter.font
}

// tryFindFallback attempts to find a fallback font for the given character.
// Returns the typeface if found, nil otherwise.
func (iter *FontMgrRunIterator) tryFindFallback(r rune) interfaces.SkTypeface {
	var bcp47 []string
	if iter.language != nil && !iter.language.AtEnd() {
		bcp47 = []string{iter.language.CurrentLanguage()}
	}

	return iter.fallbackMgr.MatchFamilyStyleCharacter(
		iter.requestName,
		iter.requestStyle,
		bcp47,
		r,
	)
}

// EndOfCurrentRun returns the byte offset to one past the last element in the current run.
func (iter *FontMgrRunIterator) EndOfCurrentRun() int {
	return iter.current
}

// AtEnd returns true if there are no more runs to consume.
func (iter *FontMgrRunIterator) AtEnd() bool {
	return iter.current >= iter.end
}

// CurrentFont returns the font for the current run.
func (iter *FontMgrRunIterator) CurrentFont() interfaces.SkFont {
	if iter.currentFont == nil {
		return iter.font
	}
	return iter.currentFont
}

// Compile-time interface check
var _ FontRunIterator = (*FontMgrRunIterator)(nil)
