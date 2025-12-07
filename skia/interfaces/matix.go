package interfaces

type SkMatrix interface {
	// Getters
	Get(index int) Scalar
	Get9() [9]Scalar
	GetRC(row, col int) Scalar
	GetPerspX() Scalar
	GetPerspY() Scalar
	GetScaleX() Scalar
	GetScaleY() Scalar
	GetSkewX() Scalar
	GetSkewY() Scalar
	GetTranslateX() Scalar
	GetTranslateY() Scalar
	GetType() MatrixType

	// Setters
	Set(index int, value Scalar)
	Set9(values [9]Scalar)
	SetAll(scaleX, skewX, transX, skewY, scaleY, transY, persp0, persp1, persp2 Scalar)
	SetScaleX(v Scalar)
	SetScaleY(v Scalar)
	SetSkewX(v Scalar)
	SetSkewY(v Scalar)
	SetTranslateX(v Scalar)
	SetTranslateY(v Scalar)
	SetPerspX(v Scalar)
	SetPerspY(v Scalar)
	SetScale(sx Scalar, sy Scalar)
	SetSkew(kx Scalar, ky Scalar)
	SetTranslate(dx Scalar, dy Scalar)
	SetRotate(degrees Scalar, px Scalar, py Scalar)
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
	PreTranslate(dx Scalar, dy Scalar)
	PreScale(sx Scalar, sy Scalar)
	PreSkew(kx Scalar, ky Scalar)
	PreRotate(degrees Scalar, px Scalar, py Scalar)
	PreConcat(other SkMatrix)
	PostTranslate(dx Scalar, dy Scalar)
	PostScale(sx Scalar, sy Scalar)
	PostSkew(kx Scalar, ky Scalar)
	PostRotate(degrees Scalar, px Scalar, py Scalar)
	PostConcat(other SkMatrix)

	// Mapping
	MapPoint(pt Point) Point
	MapXY(x, y Scalar) (Scalar, Scalar)
	MapPoints(dst []Point, src []Point) int
	MapRect(rect Rect) Rect
	MapRectToRect(src Rect, dst Rect) bool

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
