package impl

import (
	"math"

	"github.com/zodimo/go-skia-support/skia/enums"
)

// NewMatrixIdentity creates an identity matrix:
//
//	| 1 0 0 |
//	| 0 1 0 |
//	| 0 0 1 |
func NewMatrixIdentity() SkMatrix {
	return &Matrix{
		mat: [9]Scalar{1, 0, 0, 0, 1, 0, 0, 0, 1},
	}
}

// NewMatrixTranslate creates a translation matrix:
//
//	| 1 0 dx |
//	| 0 1 dy |
//	| 0 0  1 |
func NewMatrixTranslate(dx, dy Scalar) SkMatrix {
	return &Matrix{
		mat: [9]Scalar{1, 0, dx, 0, 1, dy, 0, 0, 1},
	}
}

// NewMatrixScale creates a scale matrix:
//
//	| sx  0  0 |
//	|  0 sy  0 |
//	|  0  0  1 |
func NewMatrixScale(sx, sy Scalar) SkMatrix {
	return &Matrix{
		mat: [9]Scalar{sx, 0, 0, 0, sy, 0, 0, 0, 1},
	}
}

// NewMatrixRotate creates a rotation matrix (rotation in degrees, positive rotates clockwise):
//
//	| cos(deg) -sin(deg) 0 |
//	| sin(deg)  cos(deg) 0 |
//	|    0         0     1 |
func NewMatrixRotate(deg Scalar) SkMatrix {
	rad := deg * math.Pi / 180.0
	cos := Scalar(math.Cos(float64(rad)))
	sin := Scalar(math.Sin(float64(rad)))
	return &Matrix{
		mat: [9]Scalar{cos, -sin, 0, sin, cos, 0, 0, 0, 1},
	}
}

// NewMatrixSkew creates a skew matrix
func NewMatrixSkew(kx, ky Scalar) SkMatrix {
	m := &Matrix{}
	m.SetSkew(kx, ky)
	return m
}

// NewMatrixRotateRad creates a rotation matrix from radians.
// Rotation in radians, positive rotates clockwise.
func NewMatrixRotateRad(rad Scalar) SkMatrix {
	deg := rad * 180.0 / math.Pi
	return NewMatrixRotate(deg)
}

// NewMatrixRotateWithPivot creates a rotation matrix about a pivot point.
// Rotation in degrees, positive rotates clockwise.
func NewMatrixRotateWithPivot(deg Scalar, px, py Scalar) SkMatrix {
	m := &Matrix{}
	m.SetRotate(deg, px, py)
	return m
}

// NewMatrixAll creates a matrix from all nine values:
//
//	| scaleX  skewX transX |
//	|  skewY scaleY transY |
//	| persp0 persp1 persp2 |
func NewMatrixAll(scaleX, skewX, transX, skewY, scaleY, transY, persp0, persp1, persp2 Scalar) SkMatrix {
	m := &Matrix{}
	m.SetAll(scaleX, skewX, transX, skewY, scaleY, transY, persp0, persp1, persp2)
	return m
}

// NewMatrixScaleTranslate creates a matrix that scales and then translates.
// Equivalent to Scale(sx, sy) * Translate(tx, ty)
func NewMatrixScaleTranslate(sx, sy, tx, ty Scalar) SkMatrix {
	m := NewMatrixScale(sx, sy)
	m.PostTranslate(tx, ty)
	return m
}

var _ SkMatrix = (*Matrix)(nil)

// Matrix represents a 3x3 transformation matrix.
// The matrix is stored in row-major order:
//
//	[0] = scaleX,  [1] = skewX,   [2] = transX
//	[3] = skewY,   [4] = scaleY,  [5] = transY
//	[6] = persp0,  [7] = persp1,  [8] = persp2
type Matrix struct {
	mat [9]Scalar
}

// hasPerspective returns true if the matrix contains perspective elements.
// A matrix has perspective if persp0 or persp1 is non-zero, or persp2 is not 1.
func (m Matrix) hasPerspective() bool {
	return m.mat[kMPersp0] != 0 || m.mat[kMPersp1] != 0 || m.mat[kMPersp2] != 1
}

// isIdentity returns true if the matrix is the identity matrix.
func (m Matrix) isIdentity() bool {
	return m.mat[kMScaleX] == 1 && m.mat[kMSkewX] == 0 && m.mat[kMTransX] == 0 &&
		m.mat[kMSkewY] == 0 && m.mat[kMScaleY] == 1 && m.mat[kMTransY] == 0 &&
		m.mat[kMPersp0] == 0 && m.mat[kMPersp1] == 0 && m.mat[kMPersp2] == 1
}

// MapPoint transforms a point using the matrix.
// For affine matrices: x' = x*scaleX + y*skewX + transX, y' = x*skewY + y*scaleY + transY
// For perspective matrices: applies perspective division
func (m Matrix) MapPoint(pt Point) Point {
	if m.hasPerspective() {
		return m.mapPointPerspective(pt)
	}
	return m.mapPointAffine(pt)
}

// MapXY transforms a single x,y coordinate pair using the matrix.
// Returns the transformed (x, y) coordinates.
func (m Matrix) MapXY(x, y Scalar) (Scalar, Scalar) {
	pt := m.MapPoint(Point{X: x, Y: y})
	return pt.X, pt.Y
}

// mapPointAffine transforms a point assuming the matrix has no perspective.
func (m Matrix) mapPointAffine(pt Point) Point {
	return Point{
		X: pt.X*m.mat[kMScaleX] + pt.Y*m.mat[kMSkewX] + m.mat[kMTransX],
		Y: pt.X*m.mat[kMSkewY] + pt.Y*m.mat[kMScaleY] + m.mat[kMTransY],
	}
}

// mapPointPerspective transforms a point with perspective division.
func (m Matrix) mapPointPerspective(pt Point) Point {
	x := pt.X*m.mat[kMScaleX] + pt.Y*m.mat[kMSkewX] + m.mat[kMTransX]
	y := pt.X*m.mat[kMSkewY] + pt.Y*m.mat[kMScaleY] + m.mat[kMTransY]
	z := pt.X*m.mat[kMPersp0] + pt.Y*m.mat[kMPersp1] + m.mat[kMPersp2]

	if z != 0 {
		z = 1 / z
	}

	return Point{
		X: x * z,
		Y: y * z,
	}
}

// Reset sets the matrix to the identity matrix.
func (m *Matrix) Reset() {
	id := NewMatrixIdentity().(*Matrix)
	*m = *id
}

// SetIdentity sets the matrix to the identity matrix.
func (m *Matrix) SetIdentity() {
	m.Reset()
}

// SetScale sets the matrix to scale by (sx, sy).
func (m *Matrix) SetScale(sx, sy Scalar) {
	m.mat[kMScaleX] = sx
	m.mat[kMSkewX] = 0
	m.mat[kMTransX] = 0
	m.mat[kMSkewY] = 0
	m.mat[kMScaleY] = sy
	m.mat[kMTransY] = 0
	m.mat[kMPersp0] = 0
	m.mat[kMPersp1] = 0
	m.mat[kMPersp2] = 1
}

// SetTranslate sets the matrix to translate by (dx, dy).
func (m *Matrix) SetTranslate(dx, dy Scalar) {
	m.mat[kMScaleX] = 1
	m.mat[kMSkewX] = 0
	m.mat[kMTransX] = dx
	m.mat[kMSkewY] = 0
	m.mat[kMScaleY] = 1
	m.mat[kMTransY] = dy
	m.mat[kMPersp0] = 0
	m.mat[kMPersp1] = 0
	m.mat[kMPersp2] = 1
}

// SetSkew sets the matrix to skew by (kx, ky).
func (m *Matrix) SetSkew(kx, ky Scalar) {
	m.mat[kMScaleX] = 1
	m.mat[kMSkewX] = kx
	m.mat[kMTransX] = 0
	m.mat[kMSkewY] = ky
	m.mat[kMScaleY] = 1
	m.mat[kMTransY] = 0
	m.mat[kMPersp0] = 0
	m.mat[kMPersp1] = 0
	m.mat[kMPersp2] = 1
}

// SetRotate sets the matrix to rotate by degrees about a pivot point.
func (m *Matrix) SetRotate(degrees Scalar, px, py Scalar) {
	if px == 0 && py == 0 {
		rad := degrees * math.Pi / 180.0
		cos := Scalar(math.Cos(float64(rad)))
		sin := Scalar(math.Sin(float64(rad)))
		m.mat[kMScaleX] = cos
		m.mat[kMSkewX] = -sin
		m.mat[kMTransX] = 0
		m.mat[kMSkewY] = sin
		m.mat[kMScaleY] = cos
		m.mat[kMTransY] = 0
		m.mat[kMPersp0] = 0
		m.mat[kMPersp1] = 0
		m.mat[kMPersp2] = 1
	} else {
		rad := degrees * math.Pi / 180.0
		cos := Scalar(math.Cos(float64(rad)))
		sin := Scalar(math.Sin(float64(rad)))
		dx := sin*py + (1-cos)*px
		dy := -sin*px + (1-cos)*py
		m.mat[kMScaleX] = cos
		m.mat[kMSkewX] = -sin
		m.mat[kMTransX] = dx
		m.mat[kMSkewY] = sin
		m.mat[kMScaleY] = cos
		m.mat[kMTransY] = dy
		m.mat[kMPersp0] = 0
		m.mat[kMPersp1] = 0
		m.mat[kMPersp2] = 1
	}
}

// SetConcat sets the matrix to the concatenation of a and b.
func (m *Matrix) SetConcat(a, b SkMatrix) {

	aMat := a.(*Matrix)
	bMat := b.(*Matrix)

	// Check for identity matrices
	if aMat.IsIdentity() {
		*m = *bMat
		return
	}
	if b.IsIdentity() {
		*m = *aMat
		return
	}

	// Check if both are scale+translate only
	aType := a.GetType()
	bType := b.GetType()
	if (aType&(enums.MatrixTypeAffine|enums.MatrixTypePerspective)) == 0 &&
		(bType&(enums.MatrixTypeAffine|enums.MatrixTypePerspective)) == 0 {
		// Both are scale+translate only
		m.mat[kMScaleX] = aMat.mat[kMScaleX] * bMat.mat[kMScaleX]
		m.mat[kMScaleY] = aMat.mat[kMScaleY] * bMat.mat[kMScaleY]
		m.mat[kMTransX] = aMat.mat[kMScaleX]*bMat.mat[kMTransX] + aMat.mat[kMTransX]
		m.mat[kMTransY] = aMat.mat[kMScaleY]*bMat.mat[kMTransY] + aMat.mat[kMTransY]
		m.mat[kMSkewX] = 0
		m.mat[kMSkewY] = 0
		m.mat[kMPersp0] = 0
		m.mat[kMPersp1] = 0
		m.mat[kMPersp2] = 1
		return
	}

	// General matrix multiplication
	if (aType|bType)&enums.MatrixTypePerspective != 0 {
		// Perspective case
		m.mat[kMScaleX] = rowcol3(aMat.mat[0:], bMat.mat[0:])
		m.mat[kMSkewX] = rowcol3(aMat.mat[0:], bMat.mat[1:])
		m.mat[kMTransX] = rowcol3(aMat.mat[0:], bMat.mat[2:])
		m.mat[kMSkewY] = rowcol3(aMat.mat[3:], bMat.mat[0:])
		m.mat[kMScaleY] = rowcol3(aMat.mat[3:], bMat.mat[1:])
		m.mat[kMTransY] = rowcol3(aMat.mat[3:], bMat.mat[2:])
		m.mat[kMPersp0] = rowcol3(aMat.mat[6:], bMat.mat[0:])
		m.mat[kMPersp1] = rowcol3(aMat.mat[6:], bMat.mat[1:])
		m.mat[kMPersp2] = rowcol3(aMat.mat[6:], bMat.mat[2:])
	} else {
		// Affine case
		m.mat[kMScaleX] = muladdmul(aMat.mat[kMScaleX], bMat.mat[kMScaleX], aMat.mat[kMSkewX], bMat.mat[kMSkewY])
		m.mat[kMSkewX] = muladdmul(aMat.mat[kMScaleX], bMat.mat[kMSkewX], aMat.mat[kMSkewX], bMat.mat[kMScaleY])
		m.mat[kMTransX] = muladdmul(aMat.mat[kMScaleX], bMat.mat[kMTransX], aMat.mat[kMSkewX], bMat.mat[kMTransY]) + aMat.mat[kMTransX]
		m.mat[kMSkewY] = muladdmul(aMat.mat[kMSkewY], bMat.mat[kMScaleX], aMat.mat[kMScaleY], bMat.mat[kMSkewY])
		m.mat[kMScaleY] = muladdmul(aMat.mat[kMSkewY], bMat.mat[kMSkewX], aMat.mat[kMScaleY], bMat.mat[kMScaleY])
		m.mat[kMTransY] = muladdmul(aMat.mat[kMSkewY], bMat.mat[kMTransX], aMat.mat[kMScaleY], bMat.mat[kMTransY]) + aMat.mat[kMTransY]
		m.mat[kMPersp0] = 0
		m.mat[kMPersp1] = 0
		m.mat[kMPersp2] = 1
	}
}

// PreTranslate premultiplies the matrix with a translation.
func (m *Matrix) PreTranslate(dx, dy Scalar) {
	if m.HasPerspective() {
		t := NewMatrixTranslate(dx, dy)
		m.SetConcat(m, t)
	} else {
		m.mat[kMTransX] += m.mat[kMScaleX]*dx + m.mat[kMSkewX]*dy
		m.mat[kMTransY] += m.mat[kMSkewY]*dx + m.mat[kMScaleY]*dy
	}
}

// PreScale premultiplies the matrix with a scale.
func (m *Matrix) PreScale(sx, sy Scalar) {
	if sx == 1 && sy == 1 {
		return
	}
	m.mat[kMScaleX] *= sx
	m.mat[kMSkewY] *= sx
	m.mat[kMPersp0] *= sx
	m.mat[kMSkewX] *= sy
	m.mat[kMScaleY] *= sy
	m.mat[kMPersp1] *= sy
}

// PreSkew premultiplies the matrix with a skew.
func (m *Matrix) PreSkew(kx, ky Scalar) {
	s := NewMatrixSkew(kx, ky)
	m.SetConcat(m, s)
}

// PreRotate premultiplies the matrix with a rotation.
func (m *Matrix) PreRotate(degrees Scalar, px, py Scalar) {
	r := &Matrix{}
	r.SetRotate(degrees, px, py)
	m.SetConcat(m, r)
}

// PreConcat premultiplies the matrix with another matrix.
func (m *Matrix) PreConcat(other SkMatrix) {
	if !other.IsIdentity() {
		m.SetConcat(m, other)
	}
}

// PostTranslate postmultiplies the matrix with a translation.
func (m *Matrix) PostTranslate(dx, dy Scalar) {
	if m.HasPerspective() {
		t := NewMatrixTranslate(dx, dy)
		m.SetConcat(t, m)
	} else {
		m.mat[kMTransX] += dx
		m.mat[kMTransY] += dy
	}
}

// PostScale postmultiplies the matrix with a scale.
func (m *Matrix) PostScale(sx, sy Scalar) {
	if sx == 1 && sy == 1 {
		return
	}
	m.mat[kMScaleX] *= sx
	m.mat[kMSkewX] *= sx
	m.mat[kMTransX] *= sx
	m.mat[kMSkewY] *= sy
	m.mat[kMScaleY] *= sy
	m.mat[kMTransY] *= sy
}

// PostSkew postmultiplies the matrix with a skew.
func (m *Matrix) PostSkew(kx, ky Scalar) {
	s := NewMatrixSkew(kx, ky)
	m.SetConcat(s, m)
}

// PostRotate postmultiplies the matrix with a rotation.
func (m *Matrix) PostRotate(degrees Scalar, px, py Scalar) {
	r := &Matrix{}
	r.SetRotate(degrees, px, py)
	m.SetConcat(r, m)
}

// PostConcat postmultiplies the matrix with another matrix.
func (m *Matrix) PostConcat(other SkMatrix) {
	if !other.IsIdentity() {
		m.SetConcat(other, m)
	}
}

// GetScaleX returns the x-axis scale factor.
func (m Matrix) GetScaleX() Scalar {
	return m.mat[kMScaleX]
}

// GetScaleY returns the y-axis scale factor.
func (m Matrix) GetScaleY() Scalar {
	return m.mat[kMScaleY]
}

// GetSkewX returns the x-axis skew factor.
func (m Matrix) GetSkewX() Scalar {
	return m.mat[kMSkewX]
}

// GetSkewY returns the y-axis skew factor.
func (m Matrix) GetSkewY() Scalar {
	return m.mat[kMSkewY]
}

// GetTranslateX returns the x-axis translation.
func (m Matrix) GetTranslateX() Scalar {
	return m.mat[kMTransX]
}

// GetTranslateY returns the y-axis translation.
func (m Matrix) GetTranslateY() Scalar {
	return m.mat[kMTransY]
}

// GetPerspX returns the x-axis perspective factor.
func (m Matrix) GetPerspX() Scalar {
	return m.mat[kMPersp0]
}

// GetPerspY returns the y-axis perspective factor.
func (m Matrix) GetPerspY() Scalar {
	return m.mat[kMPersp1]
}

// Get returns one matrix value by index.
// Index must be one of: kMScaleX (0), kMSkewX (1), kMTransX (2), kMSkewY (3),
// kMScaleY (4), kMTransY (5), kMPersp0 (6), kMPersp1 (7), kMPersp2 (8)
func (m Matrix) Get(index int) Scalar {
	if index >= 0 && index < 9 {
		return m.mat[index]
	}
	return 0
}

// Get9 copies all nine matrix values into a buffer.
// Values are in member value ascending order: kMScaleX, kMSkewX, kMTransX,
// kMSkewY, kMScaleY, kMTransY, kMPersp0, kMPersp1, kMPersp2
func (m Matrix) Get9() [9]Scalar {
	return m.mat
}

// GetRC returns one matrix value from a particular row/column.
// Row and column must be in range [0, 2]
func (m Matrix) GetRC(row, col int) Scalar {
	if row >= 0 && row <= 2 && col >= 0 && col <= 2 {
		return m.mat[row*3+col]
	}
	return 0
}

// Set sets one matrix value by index and invalidates the type cache.
// Index must be one of: kMScaleX (0), kMSkewX (1), kMTransX (2), kMSkewY (3),
// kMScaleY (4), kMTransY (5), kMPersp0 (6), kMPersp1 (7), kMPersp2 (8)
func (m *Matrix) Set(index int, value Scalar) {
	if index >= 0 && index < 9 {
		m.mat[index] = value
		// Type mask will be recomputed on next GetType() call
	}
}

// Set9 sets all nine matrix values from a buffer.
// Values are in member value ascending order: kMScaleX, kMSkewX, kMTransX,
// kMSkewY, kMScaleY, kMTransY, kMPersp0, kMPersp1, kMPersp2
func (m *Matrix) Set9(values [9]Scalar) {
	m.mat = values
	// Type mask will be recomputed on next GetType() call
}

// SetAll sets all nine matrix values from parameters.
// Sets matrix to:
//
//	| scaleX  skewX transX |
//	|  skewY scaleY transY |
//	| persp0 persp1 persp2 |
func (m *Matrix) SetAll(scaleX, skewX, transX, skewY, scaleY, transY, persp0, persp1, persp2 Scalar) {
	m.mat[kMScaleX] = scaleX
	m.mat[kMSkewX] = skewX
	m.mat[kMTransX] = transX
	m.mat[kMSkewY] = skewY
	m.mat[kMScaleY] = scaleY
	m.mat[kMTransY] = transY
	m.mat[kMPersp0] = persp0
	m.mat[kMPersp1] = persp1
	m.mat[kMPersp2] = persp2
	// Type mask will be recomputed on next GetType() call
}

// SetScaleX sets the horizontal scale factor.
func (m *Matrix) SetScaleX(v Scalar) {
	m.Set(kMScaleX, v)
}

// SetScaleY sets the vertical scale factor.
func (m *Matrix) SetScaleY(v Scalar) {
	m.Set(kMScaleY, v)
}

// SetSkewX sets the horizontal skew factor.
func (m *Matrix) SetSkewX(v Scalar) {
	m.Set(kMSkewX, v)
}

// SetSkewY sets the vertical skew factor.
func (m *Matrix) SetSkewY(v Scalar) {
	m.Set(kMSkewY, v)
}

// SetTranslateX sets the horizontal translation.
func (m *Matrix) SetTranslateX(v Scalar) {
	m.Set(kMTransX, v)
}

// SetTranslateY sets the vertical translation.
func (m *Matrix) SetTranslateY(v Scalar) {
	m.Set(kMTransY, v)
}

// SetPerspX sets the input x-axis perspective factor.
func (m *Matrix) SetPerspX(v Scalar) {
	m.Set(kMPersp0, v)
}

// SetPerspY sets the input y-axis perspective factor.
func (m *Matrix) SetPerspY(v Scalar) {
	m.Set(kMPersp1, v)
}

// IsIdentity returns true if the matrix is the identity matrix.
func (m Matrix) IsIdentity() bool {
	return m.isIdentity()
}

// IsScaleTranslate returns true if the matrix only scales and translates.
func (m Matrix) IsScaleTranslate() bool {
	return m.mat[kMSkewX] == 0 && m.mat[kMSkewY] == 0 &&
		m.mat[kMPersp0] == 0 && m.mat[kMPersp1] == 0 && m.mat[kMPersp2] == 1
}

// IsTranslate returns true if the matrix is identity or only translates.
// Matrix form is:
//
//	| 1 0 translate-x |
//	| 0 1 translate-y |
//	| 0 0      1      |
func (m Matrix) IsTranslate() bool {
	mask := m.GetType()
	return (mask &^ enums.MatrixTypeTranslate) == 0
}

// HasPerspective returns true if the matrix has perspective.
func (m Matrix) HasPerspective() bool {
	return m.hasPerspective()
}

// PreservesRightAngles returns true if the matrix preserves right angles.
func (m Matrix) PreservesRightAngles() bool {
	mask := m.GetType()
	if mask <= enums.MatrixTypeTranslate {
		return true
	}
	if mask&enums.MatrixTypePerspective != 0 {
		return false
	}

	mx := m.mat[kMScaleX]
	my := m.mat[kMScaleY]
	sx := m.mat[kMSkewX]
	sy := m.mat[kMSkewY]

	// Check if upper 2x2 is degenerate
	if mx*my-sx*sy == 0 {
		return false
	}

	// Check if basis vectors are orthogonal
	dot := mx*sx + sy*my
	return scalarNearlyZero(dot)
}

// RectStaysRect returns true if the matrix maps rectangles to rectangles.
func (m Matrix) RectStaysRect() bool {
	// A matrix maps rectangles to rectangles if it's identity, scale-only,
	// or rotates by multiples of 90 degrees
	mask := m.GetType()
	if mask <= enums.MatrixTypeTranslate {
		return true
	}
	if mask&enums.MatrixTypePerspective != 0 {
		return false
	}

	// Check if it's a 90-degree rotation (or multiple)
	mx := m.mat[kMScaleX]
	my := m.mat[kMScaleY]
	sx := m.mat[kMSkewX]
	sy := m.mat[kMSkewY]

	// For 90-degree rotations, one of scale components is 0 and the other is non-zero
	// Or both are non-zero with opposite signs
	return (scalarNearlyZero(mx) && !scalarNearlyZero(my) && scalarNearlyZero(sy) && !scalarNearlyZero(sx)) ||
		(!scalarNearlyZero(mx) && scalarNearlyZero(my) && !scalarNearlyZero(sy) && scalarNearlyZero(sx)) ||
		(scalarNearlyZero(sx) && scalarNearlyZero(sy))
}

// GetType returns the type of the matrix.
func (m Matrix) GetType() enums.MatrixType {
	var mask enums.MatrixType

	if m.mat[kMTransX] != 0 || m.mat[kMTransY] != 0 {
		mask |= enums.MatrixTypeTranslate
	}

	if m.mat[kMScaleX] != 1 || m.mat[kMScaleY] != 1 {
		mask |= enums.MatrixTypeScale
	}

	if m.mat[kMSkewX] != 0 || m.mat[kMSkewY] != 0 {
		mask |= enums.MatrixTypeAffine
	}

	if m.HasPerspective() {
		mask |= enums.MatrixTypePerspective
		// Perspective implies all other types
		mask |= enums.MatrixTypeTranslate | enums.MatrixTypeScale | enums.MatrixTypeAffine
	}

	return mask
}

// MapPoints applies the matrix transformation to the points.
func (m Matrix) MapPoints(dst, src []Point) int {
	count := minInt(len(dst), len(src))
	if count == 0 {
		return 0
	}

	if m.IsIdentity() {
		copy(dst[:count], src[:count])
		return count
	}

	if m.HasPerspective() {
		for i := 0; i < count; i++ {
			dst[i] = m.mapPointPerspective(src[i])
		}
	} else {
		for i := 0; i < count; i++ {
			dst[i] = m.mapPointAffine(src[i])
		}
	}

	return count
}

// MapRect applies the matrix transformation to a rectangle.
func (m Matrix) MapRect(rect Rect) Rect {
	if m.GetType() <= enums.MatrixTypeTranslate {
		// Translation only
		tx := m.mat[kMTransX]
		ty := m.mat[kMTransY]
		return Rect{
			Left:   rect.Left + tx,
			Top:    rect.Top + ty,
			Right:  rect.Right + tx,
			Bottom: rect.Bottom + ty,
		}
	}

	if m.IsScaleTranslate() {
		// Scale and translate
		sx := m.mat[kMScaleX]
		sy := m.mat[kMScaleY]
		tx := m.mat[kMTransX]
		ty := m.mat[kMTransY]

		left := rect.Left*sx + tx
		right := rect.Right*sx + tx
		top := rect.Top*sy + ty
		bottom := rect.Bottom*sy + ty

		if left > right {
			left, right = right, left
		}
		if top > bottom {
			top, bottom = bottom, top
		}

		return Rect{Left: left, Top: top, Right: right, Bottom: bottom}
	}

	// General case: map all four corners
	corners := [4]Point{
		{X: rect.Left, Y: rect.Top},
		{X: rect.Right, Y: rect.Top},
		{X: rect.Right, Y: rect.Bottom},
		{X: rect.Left, Y: rect.Bottom},
	}

	mapped := make([]Point, 4)
	m.MapPoints(mapped, corners[:])

	// Find bounding box
	minX := mapped[0].X
	maxX := mapped[0].X
	minY := mapped[0].Y
	maxY := mapped[0].Y

	for i := 1; i < 4; i++ {
		if mapped[i].X < minX {
			minX = mapped[i].X
		}
		if mapped[i].X > maxX {
			maxX = mapped[i].X
		}
		if mapped[i].Y < minY {
			minY = mapped[i].Y
		}
		if mapped[i].Y > maxY {
			maxY = mapped[i].Y
		}
	}

	return Rect{Left: minX, Top: minY, Right: maxX, Bottom: maxY}
}

// MapRectToRect applies the matrix transformation mapping src to dst.
func (m *Matrix) MapRectToRect(src, dst Rect) bool {
	// Compute scale factors
	sx := (dst.Right - dst.Left) / (src.Right - src.Left)
	sy := (dst.Bottom - dst.Top) / (src.Bottom - src.Top)

	if !isFinite(Scalar(sx)) || !isFinite(Scalar(sy)) {
		m.Reset()
		return false
	}

	// Compute translation
	tx := dst.Left - src.Left*sx
	ty := dst.Top - src.Top*sy

	m.SetScale(sx, sy)
	m.PostTranslate(tx, ty)
	return true
}

// Invert inverts the matrix if possible.
func (m *Matrix) Invert() (SkMatrix, bool) {
	mask := m.GetType()

	if mask == enums.MatrixTypeIdentity {
		return m, true
	}

	// Optimized invert for scale+translate only
	if (mask &^ (enums.MatrixTypeScale | enums.MatrixTypeTranslate)) == 0 {
		if mask&enums.MatrixTypeScale != 0 {
			// Scale + (optional) Translate
			invSX := 1.0 / m.mat[kMScaleX]
			invSY := 1.0 / m.mat[kMScaleY]

			if !isFinite(invSX) || !isFinite(invSY) {
				return &Matrix{}, false
			}

			invTX := -m.mat[kMTransX] * invSX
			invTY := -m.mat[kMTransY] * invSY

			if !isFinite(invTX) || !isFinite(invTY) {
				return &Matrix{}, false
			}

			inv := &Matrix{}
			inv.mat[kMScaleX] = invSX
			inv.mat[kMScaleY] = invSY
			inv.mat[kMTransX] = invTX
			inv.mat[kMTransY] = invTY
			inv.mat[kMSkewX] = 0
			inv.mat[kMSkewY] = 0
			inv.mat[kMPersp0] = 0
			inv.mat[kMPersp1] = 0
			inv.mat[kMPersp2] = 1
			return inv, true
		}

		// Translate-only
		if !isFinite(m.mat[kMTransX]) || !isFinite(m.mat[kMTransY]) {
			return &Matrix{}, false
		}

		return NewMatrixTranslate(-m.mat[kMTransX], -m.mat[kMTransY]), true
	}

	// General case: compute determinant and inverse
	isPersp := (mask & enums.MatrixTypePerspective) != 0
	invDet := m.computeInvDeterminant(isPersp)

	if invDet == 0 {
		return &Matrix{}, false
	}

	inv := &Matrix{}
	m.computeInv(&inv.mat, m.mat, invDet, isPersp)

	if !inv.isFinite() {
		return &Matrix{}, false
	}

	return inv, true
}

// Equals compares two matrices for equality.
// Returns true if all nine matrix values are equal.
func (m Matrix) Equals(other SkMatrix) bool {
	if other == nil {
		return false
	}
	otherMat := other.(*Matrix)
	for i := 0; i < 9; i++ {
		if m.mat[i] != otherMat.mat[i] {
			return false
		}
	}
	return true
}

func (m Matrix) isFinite() bool {
	for i := 0; i < 9; i++ {
		if !isFinite(m.mat[i]) {
			return false
		}
	}
	return true
}

func (m Matrix) computeInvDeterminant(isPerspective bool) float64 {
	det := m.computeDeterminant(isPerspective)

	// Check if determinant is nearly zero
	if scalarNearlyZero(Scalar(det)) {
		return 0
	}

	return 1.0 / det
}

func (m Matrix) computeDeterminant(isPerspective bool) float64 {
	if isPerspective {
		return float64(m.mat[kMScaleX])*
			dcross(float64(m.mat[kMScaleY]), float64(m.mat[kMPersp2]),
				float64(m.mat[kMTransY]), float64(m.mat[kMPersp1])) +
			float64(m.mat[kMSkewX])*
				dcross(float64(m.mat[kMTransY]), float64(m.mat[kMPersp0]),
					float64(m.mat[kMSkewY]), float64(m.mat[kMPersp2])) +
			float64(m.mat[kMTransX])*
				dcross(float64(m.mat[kMSkewY]), float64(m.mat[kMPersp1]),
					float64(m.mat[kMScaleY]), float64(m.mat[kMPersp0]))
	} else {
		return dcross(float64(m.mat[kMScaleX]), float64(m.mat[kMScaleY]),
			float64(m.mat[kMSkewX]), float64(m.mat[kMSkewY]))
	}
}

func (m Matrix) computeInv(dst *[9]Scalar, src [9]Scalar, invDet float64, isPersp bool) {
	if isPersp {
		dst[kMScaleX] = Scalar(scross_dscale(src[kMScaleY], src[kMPersp2], src[kMTransY], src[kMPersp1], invDet))
		dst[kMSkewX] = Scalar(scross_dscale(src[kMTransX], src[kMPersp1], src[kMSkewX], src[kMPersp2], invDet))
		dst[kMTransX] = Scalar(scross_dscale(src[kMSkewX], src[kMTransY], src[kMTransX], src[kMScaleY], invDet))
		dst[kMSkewY] = Scalar(scross_dscale(src[kMTransY], src[kMPersp0], src[kMSkewY], src[kMPersp2], invDet))
		dst[kMScaleY] = Scalar(scross_dscale(src[kMScaleX], src[kMPersp2], src[kMTransX], src[kMPersp0], invDet))
		dst[kMTransY] = Scalar(scross_dscale(src[kMTransX], src[kMSkewY], src[kMScaleX], src[kMTransY], invDet))
		dst[kMPersp0] = Scalar(scross_dscale(src[kMSkewY], src[kMPersp1], src[kMScaleY], src[kMPersp0], invDet))
		dst[kMPersp1] = Scalar(scross_dscale(src[kMSkewX], src[kMPersp0], src[kMScaleX], src[kMPersp1], invDet))
		dst[kMPersp2] = Scalar(scross_dscale(src[kMScaleX], src[kMScaleY], src[kMSkewX], src[kMSkewY], invDet))
	} else {
		dst[kMScaleX] = Scalar(float64(src[kMScaleY]) * invDet)
		dst[kMSkewX] = Scalar(-float64(src[kMSkewX]) * invDet)
		dst[kMTransX] = Scalar(dcross_dscale(float64(src[kMSkewX]), float64(src[kMTransY]), float64(src[kMScaleY]), float64(src[kMTransX]), invDet))
		dst[kMSkewY] = Scalar(-float64(src[kMSkewY]) * invDet)
		dst[kMScaleY] = Scalar(float64(src[kMScaleX]) * invDet)
		dst[kMTransY] = Scalar(dcross_dscale(float64(src[kMSkewY]), float64(src[kMTransX]), float64(src[kMScaleX]), float64(src[kMTransY]), invDet))
		dst[kMPersp0] = 0
		dst[kMPersp1] = 0
		dst[kMPersp2] = 1
	}
}
