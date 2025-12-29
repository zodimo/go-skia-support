package models

// FontWeight represents the weight (boldness) of a font.
// Values correspond to CSS font-weight values.
//
// Ported from: skia-source/include/core/SkFontStyle.h
type FontWeight int

const (
	FontWeightInvisible  FontWeight = 0
	FontWeightThin       FontWeight = 100
	FontWeightExtraLight FontWeight = 200
	FontWeightLight      FontWeight = 300
	FontWeightNormal     FontWeight = 400
	FontWeightMedium     FontWeight = 500
	FontWeightSemiBold   FontWeight = 600
	FontWeightBold       FontWeight = 700
	FontWeightExtraBold  FontWeight = 800
	FontWeightBlack      FontWeight = 900
	FontWeightExtraBlack FontWeight = 1000
)

// FontWidth represents the width of a font.
//
// Ported from: skia-source/include/core/SkFontStyle.h
type FontWidth int

const (
	FontWidthUltraCondensed FontWidth = 1
	FontWidthExtraCondensed FontWidth = 2
	FontWidthCondensed      FontWidth = 3
	FontWidthSemiCondensed  FontWidth = 4
	FontWidthNormal         FontWidth = 5
	FontWidthSemiExpanded   FontWidth = 6
	FontWidthExpanded       FontWidth = 7
	FontWidthExtraExpanded  FontWidth = 8
	FontWidthUltraExpanded  FontWidth = 9
)

// FontSlant represents the slant of a font.
//
// Ported from: skia-source/include/core/SkFontStyle.h
type FontSlant int

const (
	FontSlantUpright FontSlant = 0
	FontSlantItalic  FontSlant = 1
	FontSlantOblique FontSlant = 2
)

// FontStyle represents the style of a typeface (weight, width, slant).
// This matches C++ SkFontStyle from include/core/SkFontStyle.h
//
// Ported from: skia-source/include/core/SkFontStyle.h
type FontStyle struct {
	Weight FontWeight
	Width  FontWidth
	Slant  FontSlant
}

// NewFontStyle creates a new FontStyle with the given weight, width, and slant.
func NewFontStyle(weight FontWeight, width FontWidth, slant FontSlant) FontStyle {
	return FontStyle{
		Weight: weight,
		Width:  width,
		Slant:  slant,
	}
}

// FontStyleNormal returns the normal (upright, regular weight) style.
func FontStyleNormal() FontStyle {
	return FontStyle{
		Weight: FontWeightNormal,
		Width:  FontWidthNormal,
		Slant:  FontSlantUpright,
	}
}

// FontStyleBold returns a bold style.
func FontStyleBold() FontStyle {
	return FontStyle{
		Weight: FontWeightBold,
		Width:  FontWidthNormal,
		Slant:  FontSlantUpright,
	}
}

// FontStyleItalic returns an italic style.
func FontStyleItalic() FontStyle {
	return FontStyle{
		Weight: FontWeightNormal,
		Width:  FontWidthNormal,
		Slant:  FontSlantItalic,
	}
}

// FontStyleBoldItalic returns a bold italic style.
func FontStyleBoldItalic() FontStyle {
	return FontStyle{
		Weight: FontWeightBold,
		Width:  FontWidthNormal,
		Slant:  FontSlantItalic,
	}
}

// IsBold returns true if weight is greater than Medium.
func (fs FontStyle) IsBold() bool {
	return fs.Weight > FontWeightMedium
}

// IsItalic returns true if slant is not upright.
func (fs FontStyle) IsItalic() bool {
	return fs.Slant != FontSlantUpright
}

// Equals returns true if two FontStyles are equal.
func (fs FontStyle) Equals(other FontStyle) bool {
	return fs.Weight == other.Weight &&
		fs.Width == other.Width &&
		fs.Slant == other.Slant
}
