package impl

import "math"

// Helper functions for matrix multiplication
func rowcol3(row, col []Scalar) Scalar {
	return row[0]*col[0] + row[1]*col[3] + row[2]*col[6]
}

func muladdmul(a, b, c, d Scalar) Scalar {
	return a*b + c*d
}

// Helper functions
func scalarNearlyZero(x Scalar) bool {
	return x*x <= skScalarNearlyZero*skScalarNearlyZero
}

func scalarNearlyEqual(a, b Scalar) bool {
	return scalarNearlyZero(a - b)
}

func isFinite(x Scalar) bool {
	return !math.IsNaN(float64(x)) && !math.IsInf(float64(x), 0)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func dcross(a, b, c, d float64) float64 {
	return a*b - c*d
}

func scross_dscale(a, b, c, d Scalar, scale float64) float64 {
	return float64(a*b-c*d) * scale
}

func dcross_dscale(a, b, c, d, scale float64) float64 {
	return dcross(a, b, c, d) * scale
}
