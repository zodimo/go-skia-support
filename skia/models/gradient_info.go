// Package models provides gradient information structures.
// Ported from SkShaderBase.h
// https://github.com/google/skia/blob/main/src/shaders/SkShaderBase.h

package models

import "github.com/zodimo/go-skia-support/skia/enums"

// GradientInfo contains detailed information about a gradient shader.
//
// If the shader subclass can be represented as a gradient, AsGradient
// returns the matching GradientType enum. Also, if info is not nil,
// AsGradient populates info with the relevant parameters for the gradient.
//
// ColorCount is both an input and output parameter. On input, it indicates how
// many entries in Colors and ColorOffsets can be used, if they are non-nil.
// After AsGradient has run, ColorCount indicates how many color-offset pairs
// there are in the gradient. If there is insufficient space to store all of
// the color-offset pairs, Colors and ColorOffsets will not be altered.
//
// ColorOffsets specifies where on the range of 0 to 1 to transition to the
// given color.
//
// The meaning of Point and Radius is dependent on the type of gradient:
//   - None: info is ignored.
//   - Color: ColorOffsets[0] is meaningless.
//   - Linear: Point[0] and Point[1] are the end-points of the gradient.
//   - Radial: Point[0] and Radius[0] are the center and radius.
//   - Conical: Point[0] and Radius[0] are the center and radius of the 1st circle;
//     Point[1] and Radius[1] are the center and radius of the 2nd circle.
//   - Sweep: Point[0] is the center of the sweep.
//
// Matches C++ SkShaderBase::GradientInfo struct.
type GradientInfo struct {
	// ColorCount specifies passed size of Colors/ColorOffsets on input,
	// and actual number of colors/offsets on output.
	ColorCount int

	// Colors contains the colors in the gradient.
	Colors []Color4f

	// ColorOffsets contains the unit offset for color transitions.
	ColorOffsets []Scalar

	// Point contains type-specific point data (see struct documentation).
	Point [2]Point

	// Radius contains type-specific radius data (see struct documentation).
	Radius [2]Scalar

	// TileMode specifies the tiling behavior.
	TileMode enums.TileMode

	// GradientFlags contains optional gradient flags.
	GradientFlags enums.GradientFlags
}
