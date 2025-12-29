package enums

// FontEdging controls how edges of glyphs are drawn.
// This matches C++ SkFont::Edging from include/core/SkFont.h
//
// Ported from: skia-source/include/core/SkFont.h
type FontEdging uint8

const (
	// FontEdgingAlias draws glyph edges with no transparency (aliased)
	FontEdgingAlias FontEdging = 0

	// FontEdgingAntiAlias draws glyph edges with transparency for smooth edges
	FontEdgingAntiAlias FontEdging = 1

	// FontEdgingSubpixelAntiAlias uses subpixel rendering for glyph positioning
	FontEdgingSubpixelAntiAlias FontEdging = 2
)

// FontEdgingDefault is the default font edging (AntiAlias)
const FontEdgingDefault = FontEdgingAntiAlias

// FontHinting specifies the level of hinting applied to glyph outlines.
// This matches C++ SkFontHinting from include/core/SkFontTypes.h
//
// Ported from: skia-source/include/core/SkFontTypes.h
type FontHinting uint8

const (
	// FontHintingNone applies no hinting
	FontHintingNone FontHinting = 0

	// FontHintingSlight applies slight hinting
	FontHintingSlight FontHinting = 1

	// FontHintingNormal applies normal hinting
	FontHintingNormal FontHinting = 2

	// FontHintingFull applies full hinting
	FontHintingFull FontHinting = 3
)

// FontHintingDefault is the default font hinting (Normal)
const FontHintingDefault = FontHintingNormal
