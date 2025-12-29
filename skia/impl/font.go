package impl

import (
	"github.com/zodimo/go-skia-support/skia/enums"
)

// Font default values matching C++ Skia
const (
	FontDefaultSize   Scalar = 12.0
	FontDefaultScaleX Scalar = 1.0
	FontDefaultSkewX  Scalar = 0.0
)

// Font private flags matching C++ SkFont
const (
	fontFlagForceAutoHinting uint8 = 1 << 0
	fontFlagEmbeddedBitmaps  uint8 = 1 << 1
	fontFlagSubpixel         uint8 = 1 << 2
	fontFlagLinearMetrics    uint8 = 1 << 3
	fontFlagEmbolden         uint8 = 1 << 4
	fontFlagBaselineSnap     uint8 = 1 << 5
)

// Font controls options applied when drawing and measuring text.
//
// Ported from: skia-source/include/core/SkFont.h
type Font struct {
	typeface SkTypeface
	size     Scalar
	scaleX   Scalar
	skewX    Scalar
	flags    uint8
	edging   enums.FontEdging
	hinting  enums.FontHinting
}

// NewFont creates a new Font with default values.
func NewFont() *Font {
	return &Font{
		typeface: NewDefaultTypeface(),
		size:     FontDefaultSize,
		scaleX:   FontDefaultScaleX,
		skewX:    FontDefaultSkewX,
		flags:    fontFlagBaselineSnap, // default has baseline snap enabled
		edging:   enums.FontEdgingDefault,
		hinting:  enums.FontHintingDefault,
	}
}

// NewFontWithTypeface creates a new Font with the given typeface.
func NewFontWithTypeface(tf SkTypeface) *Font {
	f := NewFont()
	if tf != nil {
		f.typeface = tf
	}
	return f
}

// NewFontWithTypefaceAndSize creates a new Font with typeface and size.
func NewFontWithTypefaceAndSize(tf SkTypeface, size Scalar) *Font {
	f := NewFontWithTypeface(tf)
	f.SetSize(size)
	return f
}

// NewFontWithTypefaceSizeScaleSkew creates a Font with all parameters.
func NewFontWithTypefaceSizeScaleSkew(tf SkTypeface, size, scaleX, skewX Scalar) *Font {
	f := NewFontWithTypeface(tf)
	f.SetSize(size)
	f.scaleX = scaleX
	f.skewX = skewX
	return f
}

// Typeface returns the SkTypeface.
func (f *Font) Typeface() SkTypeface {
	return f.typeface
}

// Size returns the text size in local coordinate units.
func (f *Font) Size() Scalar {
	return f.size
}

// ScaleX returns the text scale on x-axis.
func (f *Font) ScaleX() Scalar {
	return f.scaleX
}

// SkewX returns the text skew on x-axis.
func (f *Font) SkewX() Scalar {
	return f.skewX
}

// Edging returns how edge pixels are drawn.
func (f *Font) Edging() enums.FontEdging {
	return f.edging
}

// Hinting returns the level of glyph outline adjustment.
func (f *Font) Hinting() enums.FontHinting {
	return f.hinting
}

// IsForceAutoHinting returns true if glyphs are always hinted.
func (f *Font) IsForceAutoHinting() bool {
	return (f.flags & fontFlagForceAutoHinting) != 0
}

// IsEmbeddedBitmaps returns true if font engine may return glyphs from font bitmaps.
func (f *Font) IsEmbeddedBitmaps() bool {
	return (f.flags & fontFlagEmbeddedBitmaps) != 0
}

// IsSubpixel returns true if glyphs may be drawn at sub-pixel offsets.
func (f *Font) IsSubpixel() bool {
	return (f.flags & fontFlagSubpixel) != 0
}

// IsLinearMetrics returns true if font and glyph metrics are linearly scalable.
func (f *Font) IsLinearMetrics() bool {
	return (f.flags & fontFlagLinearMetrics) != 0
}

// IsEmbolden returns true if bold is approximated by increasing stroke width.
func (f *Font) IsEmbolden() bool {
	return (f.flags & fontFlagEmbolden) != 0
}

// IsBaselineSnap returns true if baselines will be snapped to pixel positions.
func (f *Font) IsBaselineSnap() bool {
	return (f.flags & fontFlagBaselineSnap) != 0
}

// SetTypeface sets the SkTypeface.
func (f *Font) SetTypeface(tf SkTypeface) {
	if tf == nil {
		f.typeface = NewDefaultTypeface()
	} else {
		f.typeface = tf
	}
}

// SetSize sets the text size in local coordinate units.
func (f *Font) SetSize(size Scalar) {
	if size >= 0 {
		f.size = size
	}
}

// SetScaleX sets the text scale on x-axis.
func (f *Font) SetScaleX(scale Scalar) {
	f.scaleX = scale
}

// SetSkewX sets the text skew on x-axis.
func (f *Font) SetSkewX(skew Scalar) {
	f.skewX = skew
}

// SetEdging sets how edge pixels are drawn.
func (f *Font) SetEdging(edging enums.FontEdging) {
	f.edging = edging
}

// SetHinting sets the level of glyph outline adjustment.
func (f *Font) SetHinting(hinting enums.FontHinting) {
	f.hinting = hinting
}

// SetForceAutoHinting sets whether to always hint glyphs.
func (f *Font) SetForceAutoHinting(forceAutoHinting bool) {
	f.setFlag(fontFlagForceAutoHinting, forceAutoHinting)
}

// SetEmbeddedBitmaps requests to use bitmaps in fonts instead of outlines.
func (f *Font) SetEmbeddedBitmaps(embeddedBitmaps bool) {
	f.setFlag(fontFlagEmbeddedBitmaps, embeddedBitmaps)
}

// SetSubpixel requests that glyphs respect sub-pixel positioning.
func (f *Font) SetSubpixel(subpixel bool) {
	f.setFlag(fontFlagSubpixel, subpixel)
}

// SetLinearMetrics requests linearly scalable font and glyph metrics.
func (f *Font) SetLinearMetrics(linearMetrics bool) {
	f.setFlag(fontFlagLinearMetrics, linearMetrics)
}

// SetEmbolden increases stroke width to approximate bold typeface.
func (f *Font) SetEmbolden(embolden bool) {
	f.setFlag(fontFlagEmbolden, embolden)
}

// SetBaselineSnap requests baselines be snapped to pixels.
func (f *Font) SetBaselineSnap(baselineSnap bool) {
	f.setFlag(fontFlagBaselineSnap, baselineSnap)
}

// setFlag is a helper to set or clear a flag bit.
func (f *Font) setFlag(flag uint8, set bool) {
	if set {
		f.flags |= flag
	} else {
		f.flags &^= flag
	}
}

// MeasureText returns the advance width of text.
// This is a simplified implementation for MVP.
// In a real implementation, this would use font metrics and glyph widths.
func (f *Font) MeasureText(text []byte, encoding enums.TextEncoding, bounds *Rect) Scalar {
	if len(text) == 0 {
		if bounds != nil {
			*bounds = Rect{}
		}
		return 0
	}

	// Simplified: estimate width based on character count and font size
	// Real implementation would use actual glyph metrics
	var charCount int
	switch encoding {
	case enums.TextEncodingUTF8:
		charCount = len(text) // Simplified - doesn't handle multi-byte chars
	case enums.TextEncodingUTF16:
		charCount = len(text) / 2
	case enums.TextEncodingUTF32:
		charCount = len(text) / 4
	case enums.TextEncodingGlyphID:
		charCount = len(text) / 2 // GlyphID is uint16
	default:
		charCount = len(text)
	}

	// Estimate advance width as 0.6 * size per character (reasonable average)
	width := Scalar(charCount) * f.size * 0.6 * f.scaleX

	if bounds != nil {
		// Estimate bounding box
		*bounds = Rect{
			Left:   0,
			Top:    -f.size * 0.8, // ascent
			Right:  width,
			Bottom: f.size * 0.2, // descent
		}
	}

	return width
}

// Equals compares two fonts for equality.
func (f *Font) Equals(other *Font) bool {
	if f == nil && other == nil {
		return true
	}
	if f == nil || other == nil {
		return false
	}
	// Compare typeface by unique ID
	var typefaceEqual bool
	if f.typeface == nil && other.typeface == nil {
		typefaceEqual = true
	} else if f.typeface != nil && other.typeface != nil {
		typefaceEqual = f.typeface.UniqueID() == other.typeface.UniqueID()
	} else {
		typefaceEqual = false
	}
	return typefaceEqual &&
		f.size == other.size &&
		f.scaleX == other.scaleX &&
		f.skewX == other.skewX &&
		f.flags == other.flags &&
		f.edging == other.edging &&
		f.hinting == other.hinting
}

// Compile-time interface check
var _ SkFont = (*Font)(nil)
