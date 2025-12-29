package interfaces

import (
	"github.com/zodimo/go-skia-support/skia/models"
)

// SkFontStyleSet represents the set of styles available for a given font family.
//
// Ported from: skia-source/include/core/SkFontMgr.h
type SkFontStyleSet interface {
	// Count returns the number of styles available in this set.
	Count() int

	// GetStyle retrieves the style and name for the given index.
	GetStyle(index int, style *models.FontStyle, name *string)

	// CreateTypeface creates a typeface for the given index.
	CreateTypeface(index int) SkTypeface

	// MatchStyle matches the given pattern to the closest style in this set.
	MatchStyle(pattern models.FontStyle) SkTypeface
}

// SkFontMgr is the interface for font management.
//
// Ported from: skia-source/include/core/SkFontMgr.h
type SkFontMgr interface {
	// CountFamilies returns the number of font families available.
	CountFamilies() int

	// GetFamilyName retrieves the family name for the given index.
	GetFamilyName(index int) string

	// CreateStyleSet creates a style set for the given index.
	CreateStyleSet(index int) SkFontStyleSet

	// MatchFamily matches the given family name to a style set.
	// If the name is not found, may return an empty set or null depending on implementation.
	// In Go, we likely return nil or an empty set implementation.
	MatchFamily(familyName string) SkFontStyleSet

	// MatchFamilyStyle matches the given family name and style to a typeface.
	MatchFamilyStyle(familyName string, style models.FontStyle) SkTypeface

	// MatchFamilyStyleCharacter matches the given family name, style, and character to a typeface.
	// bcp47 is a slice of BCP47 language codes (e.g., "en-US").
	MatchFamilyStyleCharacter(familyName string, style models.FontStyle, bcp47 []string, character rune) SkTypeface

	// MakeFromData creates a typeface from the given data.
	MakeFromData(data SkData, ttcIndex int) SkTypeface

	// MakeFromFile creates a typeface from the given file path.
	MakeFromFile(path string, ttcIndex int) SkTypeface

	// LegacyMakeTypeface creates a typeface for the given family name and style.
	LegacyMakeTypeface(familyName string, style models.FontStyle) SkTypeface
}
