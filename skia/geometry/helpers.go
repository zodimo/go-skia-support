// Package geometry provides curve mathematics operations for Skia paths.
// Ported from: skia-source/src/core/SkGeometry.cpp
// https://github.com/google/skia/blob/main/src/core/SkGeometry.cpp
package geometry

import (
	"math"

	"github.com/zodimo/go-skia-support/skia/base"
	"github.com/zodimo/go-skia-support/skia/models"
)

// Type aliases for shared types
type Scalar = base.Scalar
type Point = models.Point

const (
	// ScalarNearlyZero is the tolerance for considering a value near zero
	ScalarNearlyZero = 1.0 / (1 << 12)
	// ScalarRoot2Over2 is √2/2, used for 45° conic weights
	ScalarRoot2Over2 = 0.7071067811865476
)

// FindUnitQuadRoots finds roots of Ax² + Bx + C = 0 that are in [0, 1].
// Ported from: SkFindUnitQuadRoots
func FindUnitQuadRoots(a, b, c Scalar) []Scalar {
	if a == 0 {
		return findLinearRoot(b, c)
	}

	// Compute discriminant
	dr := Scalar(math.Sqrt(float64(b*b - 4*a*c)))
	if math.IsNaN(float64(dr)) {
		return nil
	}

	// Quadratic formula
	q := -(b + Scalar(math.Copysign(float64(dr), float64(b)))) / 2
	roots := make([]Scalar, 0, 2)

	r := q / a
	if r >= 0 && r <= 1 {
		roots = append(roots, r)
	}

	if q != 0 {
		r = c / q
		if r >= 0 && r <= 1 {
			if len(roots) == 0 || roots[0] != r {
				roots = append(roots, r)
			}
		}
	}

	// Sort roots
	if len(roots) == 2 && roots[0] > roots[1] {
		roots[0], roots[1] = roots[1], roots[0]
	}

	return roots
}

// findLinearRoot finds root of Bx + C = 0 in [0, 1]
func findLinearRoot(b, c Scalar) []Scalar {
	if b == 0 {
		return nil
	}
	r := -c / b
	if r >= 0 && r <= 1 {
		return []Scalar{r}
	}
	return nil
}

// MeasureAngleBetweenVectors measures the angle between two vectors in [0, π].
// Ported from: SkMeasureAngleBetweenVectors
func MeasureAngleBetweenVectors(v1, v2 Point) Scalar {
	// Normalize vectors
	len1 := Scalar(math.Sqrt(float64(v1.X*v1.X + v1.Y*v1.Y)))
	len2 := Scalar(math.Sqrt(float64(v2.X*v2.X + v2.Y*v2.Y)))

	if len1 < ScalarNearlyZero || len2 < ScalarNearlyZero {
		return 0
	}

	v1.X /= len1
	v1.Y /= len1
	v2.X /= len2
	v2.Y /= len2

	// dot = cos(angle), cross = sin(angle)
	dot := v1.X*v2.X + v1.Y*v2.Y
	cross := v1.X*v2.Y - v1.Y*v2.X

	return Scalar(math.Atan2(math.Abs(float64(cross)), float64(dot)))
}

// FindBisector returns a vector that bisects the two given vectors.
// The bisector will point toward the interior of the provided vectors.
// Ported from: SkFindBisector
func FindBisector(v1, v2 Point) Point {
	// Normalize input vectors
	len1 := Scalar(math.Sqrt(float64(v1.X*v1.X + v1.Y*v1.Y)))
	len2 := Scalar(math.Sqrt(float64(v2.X*v2.X + v2.Y*v2.Y)))

	if len1 < ScalarNearlyZero || len2 < ScalarNearlyZero {
		return Point{X: 1, Y: 0}
	}

	v1.X /= len1
	v1.Y /= len1
	v2.X /= len2
	v2.Y /= len2

	// Bisector is sum of unit vectors
	bisector := Point{X: v1.X + v2.X, Y: v1.Y + v2.Y}

	// If vectors are nearly opposite, use perpendicular
	bisLen := Scalar(math.Sqrt(float64(bisector.X*bisector.X + bisector.Y*bisector.Y)))
	if bisLen < ScalarNearlyZero {
		// Perpendicular to v1
		return Point{X: -v1.Y, Y: v1.X}
	}

	return bisector
}

// NearlyZero returns true if the value is within ScalarNearlyZero of zero
func NearlyZero(x Scalar) bool {
	return x*x <= ScalarNearlyZero*ScalarNearlyZero
}

// NearlyEqual returns true if two scalars are nearly equal
func NearlyEqual(a, b Scalar) bool {
	return NearlyZero(a - b)
}
