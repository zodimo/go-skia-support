package enums

// TextEncoding specifies the text encoding for text operations.
// This matches C++ SkTextEncoding from include/core/SkFontTypes.h
//
// Ported from: skia-source/include/core/SkFontTypes.h
type TextEncoding uint8

const (
	// TextEncodingUTF8 specifies UTF-8 character encoding
	TextEncodingUTF8 TextEncoding = 0

	// TextEncodingUTF16 specifies UTF-16 character encoding
	TextEncodingUTF16 TextEncoding = 1

	// TextEncodingUTF32 specifies UTF-32 character encoding
	TextEncodingUTF32 TextEncoding = 2

	// TextEncodingGlyphID specifies glyph index encoding
	TextEncodingGlyphID TextEncoding = 3
)

// TextEncodingDefault is the default text encoding (UTF-8)
const TextEncodingDefault = TextEncodingUTF8
