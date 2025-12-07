package interfaces

type SkMatrix interface {
	GetPerspX() Scalar
	GetPerspY() Scalar
	GetScaleX() Scalar
	GetScaleY() Scalar
	GetSkewX() Scalar
	GetSkewY() Scalar
	GetTranslateX() Scalar
	GetTranslateY() Scalar
	GetType() MatrixType
	HasPerspective() bool
	Invert() (SkMatrix, bool)
	IsIdentity() bool
	IsScaleTranslate() bool
	MapPoint(pt Point) Point
	MapPoints(dst []Point, src []Point) int
	MapRect(rect Rect) Rect
	MapRectToRect(src Rect, dst Rect) bool
	PostConcat(other SkMatrix)
	PostRotate(degrees Scalar, px Scalar, py Scalar)
	PostScale(sx Scalar, sy Scalar)
	PostSkew(kx Scalar, ky Scalar)
	PostTranslate(dx Scalar, dy Scalar)
	PreConcat(other SkMatrix)
	PreRotate(degrees Scalar, px Scalar, py Scalar)
	PreScale(sx Scalar, sy Scalar)
	PreSkew(kx Scalar, ky Scalar)
	PreTranslate(dx Scalar, dy Scalar)
	PreservesRightAngles() bool
	RectStaysRect() bool
	Reset()
	SetConcat(a SkMatrix, b SkMatrix)
	SetIdentity()
	SetRotate(degrees Scalar, px Scalar, py Scalar)
	SetScale(sx Scalar, sy Scalar)
	SetSkew(kx Scalar, ky Scalar)
	SetTranslate(dx Scalar, dy Scalar)

	// computeDeterminant(isPerspective bool) float64
	// computeInv(dst *[9]Scalar, src [9]Scalar, invDet float64, isPersp bool)
	// computeInvDeterminant(isPerspective bool) float64
	// hasPerspective() bool
	// isFinite() bool
	// mapPointAffine(pt Point) Point
	// mapPointPerspective(pt Point) Point
}
