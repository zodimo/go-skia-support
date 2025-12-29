package interfaces

import (
	"github.com/zodimo/go-skia-support/skia/enums"
)

// SkFont controls options applied when drawing and measuring text.
//
// Ported from: skia-source/include/core/SkFont.h
type SkFont interface {
	// Typeface returns the SkTypeface. Does not alter SkTypeface SkRefCnt.
	// Returns a non-null SkTypeface.
	Typeface() SkTypeface

	// Size returns the text size in local coordinate units (EM size).
	// Default value is 12.
	Size() Scalar

	// ScaleX returns the text scale on x-axis.
	// Default value is 1.
	ScaleX() Scalar

	// SkewX returns the text skew on x-axis.
	// Default value is 0.
	SkewX() Scalar

	// Edging returns how edge pixels are drawn (aliased, anti-aliased, or subpixel).
	Edging() enums.FontEdging

	// Hinting returns the level of glyph outline adjustment.
	Hinting() enums.FontHinting

	// IsForceAutoHinting returns true if glyphs are always hinted.
	// Only meaningful on platforms using FreeType.
	IsForceAutoHinting() bool

	// IsEmbeddedBitmaps returns true if font engine may return glyphs from font bitmaps.
	IsEmbeddedBitmaps() bool

	// IsSubpixel returns true if glyphs may be drawn at sub-pixel offsets.
	IsSubpixel() bool

	// IsLinearMetrics returns true if font and glyph metrics are linearly scalable.
	IsLinearMetrics() bool

	// IsEmbolden returns true if bold is approximated by increasing stroke width.
	IsEmbolden() bool

	// IsBaselineSnap returns true if baselines will be snapped to pixel positions.
	IsBaselineSnap() bool

	// --- Setters ---

	// SetTypeface sets the SkTypeface.
	// Pass nil to clear SkTypeface and use an empty typeface.
	SetTypeface(tf SkTypeface)

	// SetSize sets the text size in local coordinate units (EM size).
	// Has no effect if size is not greater than or equal to zero.
	SetSize(size Scalar)

	// SetScaleX sets the text scale on x-axis.
	// Default value is 1.
	SetScaleX(scale Scalar)

	// SetSkewX sets the text skew on x-axis.
	// Default value is 0.
	SetSkewX(skew Scalar)

	// SetEdging sets how edge pixels are drawn.
	SetEdging(edging enums.FontEdging)

	// SetHinting sets the level of glyph outline adjustment.
	SetHinting(hinting enums.FontHinting)

	// SetForceAutoHinting sets whether to always hint glyphs.
	SetForceAutoHinting(forceAutoHinting bool)

	// SetEmbeddedBitmaps requests to use bitmaps in fonts instead of outlines.
	SetEmbeddedBitmaps(embeddedBitmaps bool)

	// SetSubpixel requests that glyphs respect sub-pixel positioning.
	SetSubpixel(subpixel bool)

	// SetLinearMetrics requests linearly scalable font and glyph metrics.
	SetLinearMetrics(linearMetrics bool)

	// SetEmbolden increases stroke width to approximate bold typeface.
	SetEmbolden(embolden bool)

	// SetBaselineSnap requests baselines be snapped to pixels.
	SetBaselineSnap(baselineSnap bool)

	// --- Glyph Operations ---

	// MeasureText returns the advance width of text.
	// The advance is the normal distance to move before drawing additional text.
	// If bounds is not nil, also returns the bounding box of text.
	MeasureText(text []byte, encoding enums.TextEncoding, bounds *Rect) Scalar
}
