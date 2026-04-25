package interfaces

import (
	"github.com/zodimo/go-skia-support/skia/base"
	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/models"
)

type SkMatrix interface {
	// Getters
	Get(index int) base.Scalar
	Get9() [9]base.Scalar
	GetRC(row, col int) base.Scalar
	GetPerspX() base.Scalar
	GetPerspY() base.Scalar
	GetScaleX() base.Scalar
	GetScaleY() base.Scalar
	GetSkewX() base.Scalar
	GetSkewY() base.Scalar
	GetTranslateX() base.Scalar
	GetTranslateY() base.Scalar
	GetType() enums.MatrixType

	// Setters
	Set(index int, value base.Scalar)
	Set9(values [9]base.Scalar)
	SetAll(scaleX, skewX, transX, skewY, scaleY, transY, persp0, persp1, persp2 base.Scalar)
	SetScaleX(v base.Scalar)
	SetScaleY(v base.Scalar)
	SetSkewX(v base.Scalar)
	SetSkewY(v base.Scalar)
	SetTranslateX(v base.Scalar)
	SetTranslateY(v base.Scalar)
	SetPerspX(v base.Scalar)
	SetPerspY(v base.Scalar)
	SetScale(sx base.Scalar, sy base.Scalar)
	SetSkew(kx base.Scalar, ky base.Scalar)
	SetTranslate(dx base.Scalar, dy base.Scalar)
	SetRotate(degrees base.Scalar, px base.Scalar, py base.Scalar)
	SetConcat(a SkMatrix, b SkMatrix)
	SetIdentity()
	Reset()

	// Queries
	HasPerspective() bool
	IsIdentity() bool
	IsScaleTranslate() bool
	IsTranslate() bool
	PreservesRightAngles() bool
	RectStaysRect() bool

	// Transformations
	PreTranslate(dx base.Scalar, dy base.Scalar)
	PreScale(sx base.Scalar, sy base.Scalar)
	PreSkew(kx base.Scalar, ky base.Scalar)
	PreRotate(degrees base.Scalar, px base.Scalar, py base.Scalar)
	PreConcat(other SkMatrix)
	PostTranslate(dx base.Scalar, dy base.Scalar)
	PostScale(sx base.Scalar, sy base.Scalar)
	PostSkew(kx base.Scalar, ky base.Scalar)
	PostRotate(degrees base.Scalar, px base.Scalar, py base.Scalar)
	PostConcat(other SkMatrix)

	// Mapping
	MapPoint(pt models.Point) models.Point
	MapXY(x, y base.Scalar) (base.Scalar, base.Scalar)
	MapPoints(dst []models.Point, src []models.Point) int
	MapRect(rect models.Rect) models.Rect
	MapRectToRect(src models.Rect, dst models.Rect) bool

	// Advanced
	Invert() (SkMatrix, bool)
	Equals(other SkMatrix) bool

	// computeDeterminant(isPerspective bool) float64
	// computeInv(dst *[9]Scalar, src [9]Scalar, invDet float64, isPersp bool)
	// computeInvDeterminant(isPerspective bool) float64
	// hasPerspective() bool
	// isFinite() bool
	// mapPointAffine(pt Point) Point
	// mapPointPerspective(pt Point) Point
}
