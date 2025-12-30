package models

import (
	"math"
)

// TransferFunction matches skcms_TransferFunction
// T(x) = cx + f             for 0 <= x < d
//
//	= (ax + b)^g + e     for d <= x
type TransferFunction struct {
	G, A, B, C, D, E, F float32
}

// TFType roughly matches skcms_TFType
type TFType int

const (
	TFTypeInvalid TFType = iota
	TFTypeSRGBish
	TFTypePQish
	TFTypeHLGish
	TFTypeHLGInvish
	TFTypeGig
)

// Matrix3x3 matches skcms_Matrix3x3
// Row-major 3x3 matrix
type Matrix3x3 struct {
	Vals [3][3]float32
}

// Concat multiplies two 3x3 matrices: dst = a * b
func Matrix3x3Concat(a, b *Matrix3x3) Matrix3x3 {
	var dst Matrix3x3
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			dst.Vals[r][c] = a.Vals[r][0]*b.Vals[0][c] +
				a.Vals[r][1]*b.Vals[1][c] +
				a.Vals[r][2]*b.Vals[2][c]
		}
	}
	return dst
}

// Vector3 matches a 3-component vector
type Vector3 struct {
	X, Y, Z float32
}

// Matrix3x3Apply multiplies a matrix by a vector: dst = m * v
func Matrix3x3Apply(m *Matrix3x3, v *Vector3) Vector3 {
	return Vector3{
		X: m.Vals[0][0]*v.X + m.Vals[0][1]*v.Y + m.Vals[0][2]*v.Z,
		Y: m.Vals[1][0]*v.X + m.Vals[1][1]*v.Y + m.Vals[1][2]*v.Z,
		Z: m.Vals[2][0]*v.X + m.Vals[2][1]*v.Y + m.Vals[2][2]*v.Z,
	}
}

// Invert inverts a 3x3 matrix. Returns false if singular.
// Basic implementation using cofactor/determinant
func Matrix3x3Invert(src *Matrix3x3, dst *Matrix3x3) bool {
	a00, a01, a02 := float64(src.Vals[0][0]), float64(src.Vals[0][1]), float64(src.Vals[0][2])
	a10, a11, a12 := float64(src.Vals[1][0]), float64(src.Vals[1][1]), float64(src.Vals[1][2])
	a20, a21, a22 := float64(src.Vals[2][0]), float64(src.Vals[2][1]), float64(src.Vals[2][2])

	b01 := a22*a11 - a12*a21
	b11 := -a22*a10 + a12*a20
	b21 := a21*a10 - a11*a20

	det := a00*b01 + a01*b11 + a02*b21

	if math.Abs(det) < 1e-10 { // Epsilon for singularity
		return false
	}

	invDet := 1.0 / det

	dst.Vals[0][0] = float32(b01 * invDet)
	dst.Vals[0][1] = float32((-a22*a01 + a02*a21) * invDet)
	dst.Vals[0][2] = float32((a12*a01 - a02*a11) * invDet)
	dst.Vals[1][0] = float32(b11 * invDet)
	dst.Vals[1][1] = float32((a22*a00 - a02*a20) * invDet)
	dst.Vals[1][2] = float32((-a12*a00 + a02*a10) * invDet)
	dst.Vals[2][0] = float32(b21 * invDet)
	dst.Vals[2][1] = float32((-a21*a00 + a01*a20) * invDet)
	dst.Vals[2][2] = float32((a11*a00 - a01*a10) * invDet)

	return true
}
