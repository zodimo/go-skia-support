package impl

import (
	"github.com/zodimo/go-skia-support/skia/base"
	"github.com/zodimo/go-skia-support/skia/helpers"
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

type SkMatrix = interfaces.SkMatrix
type SkPath = interfaces.SkPath
type Shader = interfaces.Shader
type ColorFilter = interfaces.ColorFilter
type ImageFilter = interfaces.ImageFilter
type MaskFilter = interfaces.MaskFilter
type PathEffect = interfaces.PathEffect
type Blender = interfaces.Blender

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

type Scalar = base.Scalar

type Point = models.Point

type Color4f = models.Color4f

const skScalarNearlyZero = base.SkScalarNearlyZero

type Rect = models.Rect
type RRect = models.RRect

func sign(x Scalar) int {
	return helpers.Sign(x)
}

func crossProduct(a, b Point) Scalar {
	return helpers.CrossProduct(a, b)
}

func dotProduct(a, b Point) Scalar {
	return helpers.DotProduct(a, b)
}

func scalarPin(x, lo, hi Scalar) Scalar {
	return helpers.ScalarPin(x, lo, hi)
}
