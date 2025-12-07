package helpers

func CrossProduct(a, b Point) Scalar {
	return a.X*b.Y - a.Y*b.X
}

func DotProduct(a, b Point) Scalar {
	return a.X*b.X + a.Y*b.Y
}

func Sign(x Scalar) int {
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
func ScalarPin(x, lo, hi Scalar) Scalar {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}
