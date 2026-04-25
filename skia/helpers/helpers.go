package helpers

import (
	"github.com/zodimo/go-skia-support/skia/base"
	"github.com/zodimo/go-skia-support/skia/models"
)

func CrossProduct(a, b models.Point) base.Scalar {
	return a.X*b.Y - a.Y*b.X
}

func DotProduct(a, b models.Point) base.Scalar {
	return a.X*b.X + a.Y*b.Y
}

func Sign(x base.Scalar) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

// scalarPin clamps x between lo and hi, inclusively
// Similar to SkTPin in C++
func ScalarPin(x, lo, hi base.Scalar) base.Scalar {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}
