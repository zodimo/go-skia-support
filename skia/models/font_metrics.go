package models

import (
	"github.com/zodimo/go-skia-support/skia/base"
)

// FontMetricsFlags indicate when certain metrics are valid.
// The underline or strikeout metrics may be valid and zero.
// Fonts with embedded bitmaps may not have valid underline or strikeout metrics.
//
// Ported from: skia-source/include/core/SkFontMetrics.h
type FontMetricsFlags uint32

const (
	// FontMetricsUnderlineThicknessIsValidFlag set if UnderlineThickness is valid
	FontMetricsUnderlineThicknessIsValidFlag FontMetricsFlags = 1 << 0
	// FontMetricsUnderlinePositionIsValidFlag set if UnderlinePosition is valid
	FontMetricsUnderlinePositionIsValidFlag FontMetricsFlags = 1 << 1
	// FontMetricsStrikeoutThicknessIsValidFlag set if StrikeoutThickness is valid
	FontMetricsStrikeoutThicknessIsValidFlag FontMetricsFlags = 1 << 2
	// FontMetricsStrikeoutPositionIsValidFlag set if StrikeoutPosition is valid
	FontMetricsStrikeoutPositionIsValidFlag FontMetricsFlags = 1 << 3
	// FontMetricsBoundsInvalidFlag set if Top, Bottom, XMin, XMax invalid
	FontMetricsBoundsInvalidFlag FontMetricsFlags = 1 << 4
)

// FontMetrics describes the metrics of an SkFont.
// The metric values are consistent with the Skia y-down coordinate system.
//
// Ported from: skia-source/include/core/SkFontMetrics.h
type FontMetrics struct {
	// Flags indicating which metrics are valid
	Flags FontMetricsFlags

	// Top is the greatest extent above origin of any glyph bounding box, typically negative; deprecated with variable fonts
	Top base.Scalar
	// Ascent is the distance to reserve above baseline, typically negative
	Ascent base.Scalar
	// Descent is the distance to reserve below baseline, typically positive
	Descent base.Scalar
	// Bottom is the greatest extent below origin of any glyph bounding box, typically positive; deprecated with variable fonts
	Bottom base.Scalar
	// Leading is the distance to add between lines, typically positive or zero
	Leading base.Scalar
	// AvgCharWidth is the average character width, zero if unknown
	AvgCharWidth base.Scalar
	// MaxCharWidth is the maximum character width, zero if unknown
	MaxCharWidth base.Scalar
	// XMin is the greatest extent to left of origin of any glyph bounding box, typically negative; deprecated with variable fonts
	XMin base.Scalar
	// XMax is the greatest extent to right of origin of any glyph bounding box, typically positive; deprecated with variable fonts
	XMax base.Scalar
	// XHeight is the height of lower-case 'x', zero if unknown, typically negative
	XHeight base.Scalar
	// CapHeight is the height of an upper-case letter, zero if unknown, typically negative
	CapHeight base.Scalar
	// UnderlineThickness is the underline thickness
	UnderlineThickness base.Scalar
	// UnderlinePosition is the distance from baseline to top of stroke, typically positive
	UnderlinePosition base.Scalar
	// StrikeoutThickness is the strikeout thickness
	StrikeoutThickness base.Scalar
	// StrikeoutPosition is the distance from baseline to bottom of stroke, typically negative
	StrikeoutPosition base.Scalar
}

// HasUnderlineThickness returns true if FontMetrics has a valid underline thickness.
func (f *FontMetrics) HasUnderlineThickness() (bool, base.Scalar) {
	if f.Flags&FontMetricsUnderlineThicknessIsValidFlag != 0 {
		return true, f.UnderlineThickness
	}
	return false, 0
}

// HasUnderlinePosition returns true if FontMetrics has a valid underline position.
func (f *FontMetrics) HasUnderlinePosition() (bool, base.Scalar) {
	if f.Flags&FontMetricsUnderlinePositionIsValidFlag != 0 {
		return true, f.UnderlinePosition
	}
	return false, 0
}

// HasStrikeoutThickness returns true if FontMetrics has a valid strikeout thickness.
func (f *FontMetrics) HasStrikeoutThickness() (bool, base.Scalar) {
	if f.Flags&FontMetricsStrikeoutThicknessIsValidFlag != 0 {
		return true, f.StrikeoutThickness
	}
	return false, 0
}

// HasStrikeoutPosition returns true if FontMetrics has a valid strikeout position.
func (f *FontMetrics) HasStrikeoutPosition() (bool, base.Scalar) {
	if f.Flags&FontMetricsStrikeoutPositionIsValidFlag != 0 {
		return true, f.StrikeoutPosition
	}
	return false, 0
}

// HasBounds returns true if FontMetrics has a valid Top, Bottom, XMin, and XMax.
func (f *FontMetrics) HasBounds() bool {
	return f.Flags&FontMetricsBoundsInvalidFlag == 0
}

// Equals returns true if two FontMetrics are equal.
func (f *FontMetrics) Equals(other FontMetrics) bool {
	return f.Flags == other.Flags &&
		f.Top == other.Top &&
		f.Ascent == other.Ascent &&
		f.Descent == other.Descent &&
		f.Bottom == other.Bottom &&
		f.Leading == other.Leading &&
		f.AvgCharWidth == other.AvgCharWidth &&
		f.MaxCharWidth == other.MaxCharWidth &&
		f.XMin == other.XMin &&
		f.XMax == other.XMax &&
		f.XHeight == other.XHeight &&
		f.CapHeight == other.CapHeight &&
		f.UnderlineThickness == other.UnderlineThickness &&
		f.UnderlinePosition == other.UnderlinePosition &&
		f.StrikeoutThickness == other.StrikeoutThickness &&
		f.StrikeoutPosition == other.StrikeoutPosition
}
