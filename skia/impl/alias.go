package impl

import (
	"github.com/zodimo/go-skia-support/skia/base"
	"github.com/zodimo/go-skia-support/skia/helpers"
	"github.com/zodimo/go-skia-support/skia/models"
)

// Matrix indices matching C++ SkMatrix constants
const (
	kMScaleX = base.KMScaleX // horizontal scale factor
	kMSkewX  = base.KMSkewX  // horizontal skew factor
	kMTransX = base.KMTransX // horizontal translation
	kMSkewY  = base.KMSkewY  // vertical skew factor
	kMScaleY = base.KMScaleY // vertical scale factor
	kMTransY = base.KMTransY // vertical translation
	kMPersp0 = base.KMPersp0 // input x perspective factor
	kMPersp1 = base.KMPersp1 // input y perspective factor
	kMPersp2 = base.KMPersp2 // perspective bias
)

const skScalarNearlyZero = base.SkScalarNearlyZero

func sign(x base.Scalar) int {
	return helpers.Sign(x)
}

func crossProduct(a, b models.Point) base.Scalar {
	return helpers.CrossProduct(a, b)
}

func dotProduct(a, b models.Point) base.Scalar {
	return helpers.DotProduct(a, b)
}

func scalarPin(x, lo, hi base.Scalar) base.Scalar {
	return helpers.ScalarPin(x, lo, hi)
}

type RSXform = models.RSXform
