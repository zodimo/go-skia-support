package models

import (
	"math"
	"reflect"
	"sync"
)

// ColorSpace Describes a color gamut with primaries and a white point.
// Matches C++ SkColorSpace.
type ColorSpace struct {
	transferFn TransferFunction
	toXYZD50   Matrix3x3

	// Lazy computed fields
	invTransferFn TransferFunction
	fromXYZD50    Matrix3x3
	lazyOnce      sync.Once

	hash uint64
}

// NewColorSpace creates a new ColorSpace.
func NewColorSpace(transferFn TransferFunction, toXYZD50 Matrix3x3) *ColorSpace {
	cs := &ColorSpace{
		transferFn: transferFn,
		toXYZD50:   toXYZD50,
	}
	// TODO: Calculate hash if needed for equality optimization
	return cs
}

// MakeSRGB creates the sRGB color space.
func MakeSRGB() *ColorSpace {
	return srgbSingleton
}

// MakeSRGBLinear creates a linear sRGB color space.
func MakeSRGBLinear() *ColorSpace {
	return srgbLinearSingleton
}

// MakeRGB creates an SkColorSpace from a transfer function and a row-major 3x3 transformation to XYZ.
func MakeRGB(transferFn TransferFunction, toXYZ Matrix3x3) *ColorSpace {
	// TODO: Handle invalid transfer functions if necessary (check validity)

	// Check for standard spaces to return singletons
	if isAlmostSRGB(transferFn) {
		if isAlmostSRGBMatrix(toXYZ) {
			return MakeSRGB()
		}
		// If TF is sRGB but matrix is not, use sRGB TF constant to ensure exact match
		transferFn = NamedTransferFnSRGB
	} else if isAlmostLinear(transferFn) {
		if isAlmostSRGBMatrix(toXYZ) {
			return MakeSRGBLinear()
		}
		transferFn = NamedTransferFnLinear
	}

	return NewColorSpace(transferFn, toXYZ)
}

func (cs *ColorSpace) computeLazyDstFields() {
	cs.lazyOnce.Do(func() {
		// Invert 3x3 gamut, defaulting to sRGB if we can't.
		if !Matrix3x3Invert(&cs.toXYZD50, &cs.fromXYZD50) {
			Matrix3x3Invert(&NamedGamutSRGB, &cs.fromXYZD50)
		}

		// Invert transfer function
		// For now simple inversion if possible or default to sRGB inverse
		// TODO: Implement proper TF inversion
		// fall back to linear inverse for now if generic
		cs.invTransferFn = NamedTransferFnSRGB // placeholder
	})
}

// ToXYZD50 returns the toXYZD50 matrix.
func (cs *ColorSpace) ToXYZD50() Matrix3x3 {
	return cs.toXYZD50
}

// TransferFn returns the transfer function.
func (cs *ColorSpace) TransferFn() TransferFunction {
	return cs.transferFn
}

// IsSRGB returns true if the color space is sRGB.
func (cs *ColorSpace) IsSRGB() bool {
	return cs == srgbSingleton
}

// GammaCloseToSRGB returns true if the color space gamma is near enough to be approximated as sRGB.
func (cs *ColorSpace) GammaCloseToSRGB() bool {
	return transferFnEqual(cs.transferFn, NamedTransferFnSRGB)
}

// GammaIsLinear returns true if the color space gamma is linear.
func (cs *ColorSpace) GammaIsLinear() bool {
	return transferFnEqual(cs.transferFn, NamedTransferFnLinear)
}

// Helper functions and constants

var (
	srgbSingleton       = NewColorSpace(NamedTransferFnSRGB, NamedGamutSRGB)
	srgbLinearSingleton = NewColorSpace(NamedTransferFnLinear, NamedGamutSRGB)
)

var NamedTransferFnSRGB = TransferFunction{G: 2.4, A: 1 / 1.055, B: 0.055 / 1.055, C: 1 / 12.92, D: 0.04045, E: 0.0, F: 0.0}
var NamedTransferFn2Dot2 = TransferFunction{G: 2.2, A: 1.0, B: 0.0, C: 0.0, D: 0.0, E: 0.0, F: 0.0}
var NamedTransferFnLinear = TransferFunction{G: 1.0, A: 1.0, B: 0.0, C: 0.0, D: 0.0, E: 0.0, F: 0.0}

var NamedGamutSRGB = Matrix3x3{Vals: [3][3]float32{
	{0.4358, 0.3853, 0.1430}, // Approximate values from Skia (fixed point conversion)
	{0.2224, 0.7169, 0.0606},
	{0.0139, 0.0971, 0.7141},
}}

// Note: Precise values for SRGB gamut should be used.
// Skia uses:
//     { SkFixedToFloat(0x6FA2), SkFixedToFloat(0x6299), SkFixedToFloat(0x24A0) },
//     { SkFixedToFloat(0x38F5), SkFixedToFloat(0xB785), SkFixedToFloat(0x0F84) },
//     { SkFixedToFloat(0x0390), SkFixedToFloat(0x18DA), SkFixedToFloat(0xB6CF) },
// 0x6FA2 / 65536.0 = 0.4360656738
// ...

func init() {
	// Re-initialize constants with higher precision if needed.
	// Using values from SkColorSpace.h slightly interpreted
	NamedGamutSRGB = Matrix3x3{Vals: [3][3]float32{
		{0.436065674, 0.385147095, 0.143066406},
		{0.222488403, 0.716873169, 0.060607910},
		{0.013916016, 0.097076416, 0.714096069},
	}}

	// Reset singletons with precise gamut
	srgbSingleton = NewColorSpace(NamedTransferFnSRGB, NamedGamutSRGB)
	srgbLinearSingleton = NewColorSpace(NamedTransferFnLinear, NamedGamutSRGB)
}

func isAlmostSRGB(tf TransferFunction) bool {
	return transferFnAlmostEqual(tf, NamedTransferFnSRGB)
}

func isAlmostLinear(tf TransferFunction) bool {
	return transferFnAlmostEqual(tf, NamedTransferFnLinear)
}

func isAlmostSRGBMatrix(m Matrix3x3) bool {
	// Simple check, should be element-wise almost equal
	return Matrix3x3AlmostEqual(m, NamedGamutSRGB)
}

func transferFnEqual(a, b TransferFunction) bool {
	return reflect.DeepEqual(a, b)
}

func transferFnAlmostEqual(a, b TransferFunction) bool {
	const tolerance = 0.001
	return math.Abs(float64(a.G-b.G)) < tolerance &&
		math.Abs(float64(a.A-b.A)) < tolerance &&
		math.Abs(float64(a.B-b.B)) < tolerance &&
		math.Abs(float64(a.C-b.C)) < tolerance &&
		math.Abs(float64(a.D-b.D)) < tolerance &&
		math.Abs(float64(a.E-b.E)) < tolerance &&
		math.Abs(float64(a.F-b.F)) < tolerance
}

func Matrix3x3AlmostEqual(a, b Matrix3x3) bool {
	const tolerance = 0.001
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			if math.Abs(float64(a.Vals[r][c]-b.Vals[r][c])) > tolerance {
				return false
			}
		}
	}
	return true
}
