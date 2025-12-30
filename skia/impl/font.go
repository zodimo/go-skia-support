package impl

import (
	"bytes"
	"encoding/binary"
	"unicode/utf16"

	"github.com/go-text/typesetting/font"
	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/models"
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
// It decodes the input text according to the specified encoding and measures
// the cumulative advance width of the corresponding glyphs.
func (f *Font) MeasureText(text []byte, encoding enums.TextEncoding, bounds *Rect) Scalar {
	if len(text) == 0 {
		if bounds != nil {
			*bounds = Rect{}
		}
		return 0
	}

	var totalWidth Scalar
	var runes []rune
	switch encoding {
	case enums.TextEncodingUTF8:
		runes = []rune(string(text))
	case enums.TextEncodingUTF16:
		if len(text)%2 != 0 {
			// standard behavior for invalid length? often 0 or best effort
			return 0
		}
		u16s := make([]uint16, len(text)/2)
		// Assume Little Endian for Skia compatibility unless BOM says otherwise,
		// but Skia usually defaults to native or specific setup.
		// Go's text/encoding/unicode might be overkill here.
		// We'll stick to LE as a safe default for modern systems.
		err := binary.Read(bytes.NewReader(text), binary.LittleEndian, &u16s)
		if err != nil {
			return 0
		}
		runes = utf16.Decode(u16s)
	case enums.TextEncodingUTF32:
		if len(text)%4 != 0 {
			return 0
		}
		runes = make([]rune, len(text)/4)
		// Assume Little Endian
		err := binary.Read(bytes.NewReader(text), binary.LittleEndian, &runes)
		if err != nil {
			// If binary.Read fails on []rune (int32 alias), do manual loop
			// binary.Read works for fixed-size values.
			// Let's degrade to manual loop to be safe if []rune is platform dependent (usually int32)
			rdr := bytes.NewReader(text)
			for i := 0; i < len(runes); i++ {
				var u32 uint32
				if err := binary.Read(rdr, binary.LittleEndian, &u32); err != nil {
					break
				}
				runes[i] = rune(u32)
			}
		}
	case enums.TextEncodingGlyphID:
		// Input text is actually a slice of GlyphIDs (uint16)
		if len(text)%2 != 0 {
			return 0
		}
		glyphs := make([]uint16, len(text)/2)
		binary.Read(bytes.NewReader(text), binary.LittleEndian, &glyphs)

		var totalWidth Scalar
		widths := f.GetWidths(glyphs)
		for _, w := range widths {
			totalWidth += w
		}

		if bounds != nil {
			metrics := f.GetMetrics()
			*bounds = Rect{
				Left:   0,
				Top:    metrics.Ascent,
				Right:  totalWidth,
				Bottom: metrics.Descent,
			}
		}
		return totalWidth

	default:
		runes = []rune(string(text))
	}

	glyphs := make([]uint16, len(runes))
	for i, r := range runes {
		glyphs[i] = f.UnicharToGlyph(r)
	}

	widths := f.GetWidths(glyphs)
	for _, w := range widths {
		totalWidth += w
	}

	if bounds != nil {
		metrics := f.GetMetrics()
		*bounds = Rect{
			Left:   0,
			Top:    metrics.Ascent, // Ascent is typically negative
			Right:  totalWidth,
			Bottom: metrics.Descent,
		}
	}

	return totalWidth
}

// UnicharToGlyph returns the glyph ID for the given Unicode character.
// Delegates to the typeface, matching C++ SkFont::unicharToGlyph.
// Ported from: SkFont::unicharToGlyph
func (f *Font) UnicharToGlyph(unichar rune) uint16 {
	return f.typeface.UnicharToGlyph(unichar)
}

// GetWidths returns the advance widths for a slice of glyph IDs.
func (f *Font) GetWidths(glyphs []uint16) []Scalar {
	if len(glyphs) == 0 {
		return nil
	}
	widths := make([]Scalar, len(glyphs))

	tf, ok := f.typeface.(*Typeface)
	if !ok || tf.goTextFace == nil {
		// Fallback to heuristic
		charWidth := f.size * 0.6 * f.scaleX
		for i := range glyphs {
			widths[i] = charWidth
		}
		return widths
	}

	face := tf.goTextFace
	upem := Scalar(face.Upem())
	scale := f.size / upem

	for i, gid := range glyphs {
		// HorizontalAdvance returns the advance width in font units
		adv := Scalar(face.HorizontalAdvance(font.GID(gid)))
		widths[i] = adv * scale * f.scaleX
	}
	return widths
}

// GetMetrics returns the font metrics for this font.
func (f *Font) GetMetrics() models.FontMetrics {
	tf, ok := f.typeface.(*Typeface)
	if !ok || tf.goTextFace == nil {
		// Fallback
		return models.FontMetrics{
			Ascent:  -f.size * 0.8,
			Descent: f.size * 0.2,
			Leading: f.size * 0.05,
		}
	}

	face := tf.goTextFace
	scale := f.size / Scalar(face.Upem())

	extents, ok := face.FontHExtents()
	if !ok {
		// Fallback if no extents
		return models.FontMetrics{
			Ascent:  -f.size * 0.8,
			Descent: f.size * 0.2,
			Leading: f.size * 0.05,
		}
	}

	// SkFontMetrics conventions (SkFontMetrics.h):
	// fAscent: distance to reserve above baseline, typically negative.
	// fDescent: distance to reserve below baseline, typically positive.
	//
	// go-text/typesetting/font conventions (standard OpenType):
	// Ascender: typically positive (up).
	// Descender: typically negative (down).
	//
	// Conversion:
	// Skia Ascent = -Ascender
	// Skia Descent = -Descender (since Descender is negative, result is positive)

	return models.FontMetrics{
		Ascent:  Scalar(-extents.Ascender) * scale,
		Descent: Scalar(-extents.Descender) * scale,
		Leading: Scalar(extents.LineGap) * scale,
	}
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
