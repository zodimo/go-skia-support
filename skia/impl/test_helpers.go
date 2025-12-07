package impl

import (
	"github.com/zodimo/go-skia-support/skia/base"
)

// Tolerance constants for floating-point comparisons.
// These values match Skia C++ test tolerance thresholds.
const (
	// ScalarTolerance is the tolerance used for scalar (float32) comparisons.
	// Matches C++ test: SK_Scalar1 / 200000 = 1.0 / 200000 = 0.000005
	// This tolerance is used in nearly_equal_scalar() from MatrixTest.cpp
	ScalarTolerance base.Scalar = 1.0 / 200000

	// Epsilon is the machine epsilon for float32 (FLT_EPSILON).
	// Approximately 1.19209290E-07 = 1 / (2^23)
	// Used for very precise comparisons when needed.
	Epsilon base.Scalar = 1.19209290e-07
)

// NearlyEqualScalar compares two scalar values with tolerance.
// Returns true if |a - b| <= ScalarTolerance.
//
// Ported from: skia-source/tests/MatrixTest.cpp:nearly_equal_scalar()
// Tolerance: SK_Scalar1 / 200000 = 0.000005
func NearlyEqualScalar(a, b base.Scalar) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= ScalarTolerance
}

// NearlyEqual compares two matrices element-by-element using NearlyEqualScalar.
// Returns true if all 9 matrix elements are nearly equal within tolerance.
//
// Ported from: skia-source/tests/MatrixTest.cpp:nearly_equal()
func NearlyEqual(a, b SkMatrix) bool {
	if a == nil || b == nil {
		return a == b
	}

	// Compare all 9 matrix elements
	for i := 0; i < 9; i++ {
		if !NearlyEqualScalar(a.Get(i), b.Get(i)) {
			return false
		}
	}
	return true
}

// IsIdentity checks if a matrix is the identity matrix within tolerance.
// Identity matrix:
//
//	| 1 0 0 |
//	| 0 1 0 |
//	| 0 0 1 |
//
// Ported from: skia-source/tests/MatrixTest.cpp:is_identity()
func IsIdentity(m SkMatrix) bool {
	if m == nil {
		return false
	}

	identity := NewMatrixIdentity()
	return NearlyEqual(m, identity)
}

// IsFiniteMatrix checks if all matrix elements are finite.
// Uses the existing IsFinite() function from path_helper.go for scalar checks.
func IsFiniteMatrix(m SkMatrix) bool {
	if m == nil {
		return false
	}

	for i := 0; i < 9; i++ {
		if !IsFinite(m.Get(i)) {
			return false
		}
	}
	return true
}

