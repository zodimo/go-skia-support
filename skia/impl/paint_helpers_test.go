package impl

import (
	"math"
	"testing"

	"github.com/zodimo/go-skia-support/skia/base"
	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/models"
)

// TestColor4fFromColor tests the Color4fFromColor helper function
// Ported from: skia-source/tests/SkColor4fTest.cpp:DEF_TEST(SkColor4f_FromColor, reporter)
func TestColor4fFromColor(t *testing.T) {
	tests := []struct {
		name     string
		color    uint32
		expected models.Color4f
	}{
		{
			name:  "SK_ColorBLACK",
			color: 0xFF000000, // ARGB: alpha=FF, r=00, g=00, b=00
			expected: models.Color4f{
				R: 0.0,
				G: 0.0,
				B: 0.0,
				A: 1.0,
			},
		},
		{
			name:  "SK_ColorWHITE",
			color: 0xFFFFFFFF, // ARGB: alpha=FF, r=FF, g=FF, b=FF
			expected: models.Color4f{
				R: 1.0,
				G: 1.0,
				B: 1.0,
				A: 1.0,
			},
		},
		{
			name:  "SK_ColorRED",
			color: 0xFFFF0000, // ARGB: alpha=FF, r=FF, g=00, b=00
			expected: models.Color4f{
				R: 1.0,
				G: 0.0,
				B: 0.0,
				A: 1.0,
			},
		},
		{
			name:  "SK_ColorGREEN",
			color: 0xFF00FF00, // ARGB: alpha=FF, r=00, g=FF, b=00
			expected: models.Color4f{
				R: 0.0,
				G: 1.0,
				B: 0.0,
				A: 1.0,
			},
		},
		{
			name:  "SK_ColorBLUE",
			color: 0xFF0000FF, // ARGB: alpha=FF, r=00, g=00, b=FF
			expected: models.Color4f{
				R: 0.0,
				G: 0.0,
				B: 1.0,
				A: 1.0,
			},
		},
		{
			name:  "transparent (zero alpha)",
			color: 0x00000000, // ARGB: alpha=00, r=00, g=00, b=00
			expected: models.Color4f{
				R: 0.0,
				G: 0.0,
				B: 0.0,
				A: 0.0,
			},
		},
		{
			name:  "semi-transparent red",
			color: 0x80FF0000, // ARGB: alpha=80 (128/255 ≈ 0.502), r=FF, g=00, b=00
			expected: models.Color4f{
				R: 1.0,
				G: 0.0,
				B: 0.0,
				A: 128.0 / 255.0, // ≈ 0.502
			},
		},
		{
			name:  "gray (128,128,128)",
			color: 0xFF808080, // ARGB: alpha=FF, r=80 (128), g=80, b=80
			expected: models.Color4f{
				R: 128.0 / 255.0, // ≈ 0.502
				G: 128.0 / 255.0,
				B: 128.0 / 255.0,
				A: 1.0,
			},
		},
		{
			name:  "yellow",
			color: 0xFFFFFF00, // ARGB: alpha=FF, r=FF, g=FF, b=00
			expected: models.Color4f{
				R: 1.0,
				G: 1.0,
				B: 0.0,
				A: 1.0,
			},
		},
		{
			name:  "cyan",
			color: 0xFF00FFFF, // ARGB: alpha=FF, r=00, g=FF, b=FF
			expected: models.Color4f{
				R: 0.0,
				G: 1.0,
				B: 1.0,
				A: 1.0,
			},
		},
		{
			name:  "magenta",
			color: 0xFFFF00FF, // ARGB: alpha=FF, r=FF, g=00, b=FF
			expected: models.Color4f{
				R: 1.0,
				G: 0.0,
				B: 1.0,
				A: 1.0,
			},
		},
		{
			name:  "partial alpha with color",
			color: 0x7F3F5F9F, // ARGB: alpha=7F (127), r=3F (63), g=5F (95), b=9F (159)
			expected: models.Color4f{
				R: 63.0 / 255.0,  // ≈ 0.247
				G: 95.0 / 255.0,  // ≈ 0.373
				B: 159.0 / 255.0, // ≈ 0.624
				A: 127.0 / 255.0, // ≈ 0.498
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Color4fFromColor(tt.color)

			// Compare each component with tolerance
			if !NearlyEqualScalar(result.R, tt.expected.R) {
				t.Errorf("R component mismatch: got %f, expected %f", result.R, tt.expected.R)
			}
			if !NearlyEqualScalar(result.G, tt.expected.G) {
				t.Errorf("G component mismatch: got %f, expected %f", result.G, tt.expected.G)
			}
			if !NearlyEqualScalar(result.B, tt.expected.B) {
				t.Errorf("B component mismatch: got %f, expected %f", result.B, tt.expected.B)
			}
			if !NearlyEqualScalar(result.A, tt.expected.A) {
				t.Errorf("A component mismatch: got %f, expected %f", result.A, tt.expected.A)
			}
		})
	}
}

// TestColor4fFromColor_RoundTrip tests round-trip conversion
// Converts Color4f -> SkColor -> Color4f and verifies equality
func TestColor4fFromColor_RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		color models.Color4f
	}{
		{
			name: "opaque red",
			color: models.Color4f{
				R: 1.0,
				G: 0.0,
				B: 0.0,
				A: 1.0,
			},
		},
		{
			name: "semi-transparent blue",
			color: models.Color4f{
				R: 0.0,
				G: 0.0,
				B: 1.0,
				A: 0.5,
			},
		},
		{
			name: "gray",
			color: models.Color4f{
				R: 0.5,
				G: 0.5,
				B: 0.5,
				A: 1.0,
			},
		},
		{
			name: "transparent",
			color: models.Color4f{
				R: 0.0,
				G: 0.0,
				B: 0.0,
				A: 0.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert Color4f -> SkColor
			skColor := tt.color.ToSkColor()

			// Convert SkColor -> Color4f
			result := Color4fFromColor(skColor)

			// Note: Round-trip may have slight precision loss due to uint8 quantization
			// Use slightly larger tolerance for round-trip tests
			tolerance := base.Scalar(1.0 / 255.0) // One quantized step

			diffR := result.R - tt.color.R
			if diffR < 0 {
				diffR = -diffR
			}
			if diffR > tolerance {
				t.Errorf("R component round-trip mismatch: got %f, expected %f (diff: %f)", result.R, tt.color.R, diffR)
			}

			diffG := result.G - tt.color.G
			if diffG < 0 {
				diffG = -diffG
			}
			if diffG > tolerance {
				t.Errorf("G component round-trip mismatch: got %f, expected %f (diff: %f)", result.G, tt.color.G, diffG)
			}

			diffB := result.B - tt.color.B
			if diffB < 0 {
				diffB = -diffB
			}
			if diffB > tolerance {
				t.Errorf("B component round-trip mismatch: got %f, expected %f (diff: %f)", result.B, tt.color.B, diffB)
			}

			diffA := result.A - tt.color.A
			if diffA < 0 {
				diffA = -diffA
			}
			if diffA > tolerance {
				t.Errorf("A component round-trip mismatch: got %f, expected %f (diff: %f)", result.A, tt.color.A, diffA)
			}
		})
	}
}

// TestGetInflationRadiusForStroke tests the GetInflationRadiusForStroke helper function
// Ported from: skia-source/src/core/SkStrokeRec.cpp:GetInflationRadius()
func TestGetInflationRadiusForStroke(t *testing.T) {
	tests := []struct {
		name        string
		join        enums.PaintJoin
		miterLimit  base.Scalar
		cap         enums.PaintCap
		strokeWidth base.Scalar
		matrixScale []base.Scalar // optional
		expected    base.Scalar
	}{
		// Fill cases (negative stroke width)
		{
			name:        "fill (negative width)",
			join:        enums.PaintJoinMiter,
			miterLimit:   4.0,
			cap:         enums.PaintCapRound,
			strokeWidth: -1.0,
			expected:    0.0,
		},
		{
			name:        "fill (negative width, any join/cap)",
			join:        enums.PaintJoinRound,
			miterLimit:   10.0,
			cap:         enums.PaintCapSquare,
			strokeWidth: -5.0,
			expected:    0.0,
		},

		// Hairline cases (zero stroke width)
		{
			name:        "hairline (zero width, no matrixScale)",
			join:        enums.PaintJoinMiter,
			miterLimit:   4.0,
			cap:         enums.PaintCapRound,
			strokeWidth: 0.0,
			expected:    1.0, // Default hairline radius
		},
		{
			name:        "hairline (zero width, with matrixScale)",
			join:        enums.PaintJoinMiter,
			miterLimit:   4.0,
			cap:         enums.PaintCapRound,
			strokeWidth: 0.0,
			matrixScale: []base.Scalar{2.5},
			expected:    2.5, // Uses matrixScale
		},
		{
			name:        "hairline (zero width, zero matrixScale)",
			join:        enums.PaintJoinMiter,
			miterLimit:   4.0,
			cap:         enums.PaintCapRound,
			strokeWidth: 0.0,
			matrixScale: []base.Scalar{0.0}, // Zero or negative uses default
			expected:    1.0,
		},
		{
			name:        "hairline (zero width, negative matrixScale)",
			join:        enums.PaintJoinMiter,
			miterLimit:   4.0,
			cap:         enums.PaintCapRound,
			strokeWidth: 0.0,
			matrixScale: []base.Scalar{-1.0}, // Negative uses default
			expected:    1.0,
		},

		// Normal stroke cases (positive stroke width)
		{
			name:        "normal stroke, round join, round cap",
			join:        enums.PaintJoinRound,
			miterLimit:   4.0,
			cap:         enums.PaintCapRound,
			strokeWidth: 10.0,
			expected:    5.0, // width/2 * 1.0 = 5.0
		},
		{
			name:        "normal stroke, bevel join, round cap",
			join:        enums.PaintJoinBevel,
			miterLimit:   4.0,
			cap:         enums.PaintCapRound,
			strokeWidth: 8.0,
			expected:    4.0, // width/2 * 1.0 = 4.0
		},
		{
			name:        "normal stroke, miter join (miterLimit > 1.0), round cap",
			join:        enums.PaintJoinMiter,
			miterLimit:   4.0,
			cap:         enums.PaintCapRound,
			strokeWidth: 6.0,
			expected:    12.0, // width/2 * 4.0 = 3.0 * 4.0 = 12.0
		},
		{
			name:        "normal stroke, miter join (miterLimit < 1.0), round cap",
			join:        enums.PaintJoinMiter,
			miterLimit:   0.5,
			cap:         enums.PaintCapRound,
			strokeWidth: 6.0,
			expected:    3.0, // width/2 * 1.0 = 3.0 (miterLimit < 1.0, so multiplier stays 1.0)
		},
		{
			name:        "normal stroke, round join, square cap",
			join:        enums.PaintJoinRound,
			miterLimit:  4.0,
			cap:         enums.PaintCapSquare,
			strokeWidth: 10.0,
			expected:    base.Scalar(math.Sqrt2) * 5.0, // width/2 * sqrt(2) ≈ 7.071
		},
		{
			name:        "normal stroke, bevel join, square cap",
			join:        enums.PaintJoinBevel,
			miterLimit:  4.0,
			cap:         enums.PaintCapSquare,
			strokeWidth: 8.0,
			expected:    base.Scalar(math.Sqrt2) * 4.0, // width/2 * sqrt(2) ≈ 5.657
		},
		{
			name:        "normal stroke, miter join (miterLimit > sqrt2), square cap",
			join:        enums.PaintJoinMiter,
			miterLimit:  4.0,
			cap:         enums.PaintCapSquare,
			strokeWidth: 6.0,
			expected:    12.0, // width/2 * max(1.0, 4.0, sqrt2) = 3.0 * 4.0 = 12.0
		},
		{
			name:        "normal stroke, miter join (miterLimit < sqrt2), square cap",
			join:        enums.PaintJoinMiter,
			miterLimit:  1.0,
			cap:         enums.PaintCapSquare,
			strokeWidth: 6.0,
			expected:    base.Scalar(math.Sqrt2) * 3.0, // width/2 * max(1.0, 1.0, sqrt2) = 3.0 * sqrt2 ≈ 4.243
		},
		{
			name:        "normal stroke, miter join (miterLimit = sqrt2), square cap",
			join:        enums.PaintJoinMiter,
			miterLimit:  base.Scalar(math.Sqrt2),
			cap:         enums.PaintCapSquare,
			strokeWidth: 6.0,
			expected:    base.Scalar(math.Sqrt2) * 3.0, // width/2 * max(1.0, sqrt2, sqrt2) = 3.0 * sqrt2 ≈ 4.243
		},

		// Edge cases
		{
			name:        "very small stroke width",
			join:        enums.PaintJoinRound,
			miterLimit:  4.0,
			cap:         enums.PaintCapRound,
			strokeWidth: 0.001,
			expected:    0.0005, // width/2 * 1.0 = 0.0005
		},
		{
			name:        "very large stroke width",
			join:        enums.PaintJoinRound,
			miterLimit:  4.0,
			cap:         enums.PaintCapRound,
			strokeWidth: 1000.0,
			expected:    500.0, // width/2 * 1.0 = 500.0
		},
		{
			name:        "very large miter limit",
			join:        enums.PaintJoinMiter,
			miterLimit:  100.0,
			cap:         enums.PaintCapRound,
			strokeWidth: 10.0,
			expected:    500.0, // width/2 * 100.0 = 5.0 * 100.0 = 500.0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result base.Scalar
			if len(tt.matrixScale) > 0 {
				result = GetInflationRadiusForStroke(tt.join, tt.miterLimit, tt.cap, tt.strokeWidth, tt.matrixScale...)
			} else {
				result = GetInflationRadiusForStroke(tt.join, tt.miterLimit, tt.cap, tt.strokeWidth)
			}

			if !NearlyEqualScalar(result, tt.expected) {
				t.Errorf("GetInflationRadiusForStroke() = %f, expected %f", result, tt.expected)
			}
		})
	}
}

// TestGetInflationRadiusForStroke_EdgeCases tests edge cases and boundary conditions
func TestGetInflationRadiusForStroke_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		join        enums.PaintJoin
		miterLimit  base.Scalar
		cap         enums.PaintCap
		strokeWidth base.Scalar
		matrixScale []base.Scalar
		expected    base.Scalar
	}{
		// Boundary between fill and hairline
		{
			name:        "stroke width exactly zero",
			join:        enums.PaintJoinRound,
			miterLimit:  4.0,
			cap:         enums.PaintCapRound,
			strokeWidth: 0.0,
			expected:    1.0,
		},
		{
			name:        "stroke width just negative",
			join:        enums.PaintJoinRound,
			miterLimit:  4.0,
			cap:         enums.PaintCapRound,
			strokeWidth: -0.0001,
			expected:    0.0,
		},
		{
			name:        "stroke width just positive",
			join:        enums.PaintJoinRound,
			miterLimit:  4.0,
			cap:         enums.PaintCapRound,
			strokeWidth: 0.0001,
			expected:    0.00005, // width/2 * 1.0
		},

		// Miter limit boundary conditions
		{
			name:        "miter limit exactly 1.0",
			join:        enums.PaintJoinMiter,
			miterLimit:  1.0,
			cap:         enums.PaintCapRound,
			strokeWidth: 10.0,
			expected:    5.0, // width/2 * 1.0 (miterLimit == 1.0, so multiplier stays 1.0)
		},
		{
			name:        "miter limit just above 1.0",
			join:        enums.PaintJoinMiter,
			miterLimit:  1.0001,
			cap:         enums.PaintCapRound,
			strokeWidth: 10.0,
			expected:    5.0005, // width/2 * 1.0001 ≈ 5.0005
		},
		{
			name:        "miter limit just below 1.0",
			join:        enums.PaintJoinMiter,
			miterLimit:  0.9999,
			cap:         enums.PaintCapRound,
			strokeWidth: 10.0,
			expected:    5.0, // width/2 * 1.0 (miterLimit < 1.0, so multiplier stays 1.0)
		},

		// Matrix scale boundary conditions
		{
			name:        "hairline with matrixScale exactly zero",
			join:        enums.PaintJoinRound,
			miterLimit:  4.0,
			cap:         enums.PaintCapRound,
			strokeWidth: 0.0,
			matrixScale: []base.Scalar{0.0},
			expected:    1.0, // Uses default (matrixScale <= 0)
		},
		{
			name:        "hairline with matrixScale just positive",
			join:        enums.PaintJoinRound,
			miterLimit:  4.0,
			cap:         enums.PaintCapRound,
			strokeWidth: 0.0,
			matrixScale: []base.Scalar{0.0001},
			expected:    0.0001, // Uses matrixScale
		},
		{
			name:        "hairline with matrixScale just negative",
			join:        enums.PaintJoinRound,
			miterLimit:  4.0,
			cap:         enums.PaintCapRound,
			strokeWidth: 0.0,
			matrixScale: []base.Scalar{-0.0001},
			expected:    1.0, // Uses default (matrixScale <= 0)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result base.Scalar
			if len(tt.matrixScale) > 0 {
				result = GetInflationRadiusForStroke(tt.join, tt.miterLimit, tt.cap, tt.strokeWidth, tt.matrixScale...)
			} else {
				result = GetInflationRadiusForStroke(tt.join, tt.miterLimit, tt.cap, tt.strokeWidth)
			}

			if !NearlyEqualScalar(result, tt.expected) {
				t.Errorf("GetInflationRadiusForStroke() = %f, expected %f", result, tt.expected)
			}
		})
	}
}

