package impl

import (
	"math"

	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/models"
)

// FromColor creates a Color4f from a uint32 SkColor (ARGB format)
// This is equivalent to SkColor4f::FromColor() in C++
// SkColor is packed as ARGB: (a << 24) | (r << 16) | (g << 8) | b
// Color4f stores RGBA as floats in [0, 1] range
func Color4fFromColor(color uint32) models.Color4f {
	// Extract ARGB components
	a := Scalar((color>>24)&0xFF) / 255.0
	r := Scalar((color>>16)&0xFF) / 255.0
	g := Scalar((color>>8)&0xFF) / 255.0
	b := Scalar(color&0xFF) / 255.0
	// Return as RGBA Color4f
	return models.Color4f{R: r, G: g, B: b, A: a}
}

// GetInflationRadiusForStroke computes the inflation radius for stroke effects
// based on join type, miter limit, cap type, and stroke width.
// This is a static helper function equivalent to SkStrokeRec::GetInflationRadius().
// If matrixScale is provided and > 0, it will be used for hairline strokes (width == 0).
// Otherwise, hairlines default to 1.0.
// For hairlines, the width is determined in device space, so matrixScale should be
// the maximum scale factor from the transformation matrix (e.g., from Matrix.GetMaxScale()).
func GetInflationRadiusForStroke(join enums.PaintJoin, miterLimit Scalar, cap enums.PaintCap, strokeWidth Scalar, matrixScale ...Scalar) Scalar {
	if strokeWidth < 0 { // fill
		return 0
	} else if strokeWidth == 0 {
		// Hairline stroke - use matrixScale if provided, otherwise default to 1.0
		// Hairlines are determined in device space, so we need the matrix scale factor
		// to properly compute the inflation radius.
		if len(matrixScale) > 0 && matrixScale[0] > 0 {
			return matrixScale[0]
		}
		return 1.0
	}

	// Since we're stroked, outset the rect by the radius (and join type, caps)
	multiplier := Scalar(1.0)
	if join == enums.PaintJoinMiter {
		if miterLimit > multiplier {
			multiplier = miterLimit
		}
	}
	if cap == enums.PaintCapSquare {
		sqrt2 := Scalar(math.Sqrt2)
		if sqrt2 > multiplier {
			multiplier = sqrt2
		}
	}
	return strokeWidth / 2 * multiplier
}
