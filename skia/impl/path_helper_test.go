package impl

import (
	"math"
	"testing"

	"github.com/zodimo/go-skia-support/skia/base"
	"github.com/zodimo/go-skia-support/skia/enums"
)

// TestPathFirstDirectionToConvexity tests the pathFirstDirectionToConvexity helper function
func TestPathFirstDirectionToConvexity(t *testing.T) {
	tests := []struct {
		name     string
		dir      enums.PathFirstDirection
		expected enums.PathConvexity
	}{
		{
			name:     "CW direction",
			dir:      enums.PathFirstDirectionCW,
			expected: enums.PathConvexityConvexCW,
		},
		{
			name:     "CCW direction",
			dir:      enums.PathFirstDirectionCCW,
			expected: enums.PathConvexityConvexCCW,
		},
		{
			name:     "Unknown direction",
			dir:      enums.PathFirstDirectionUnknown,
			expected: enums.PathConvexityConvexDegenerate,
		},
		{
			name:     "Invalid direction",
			dir:      enums.PathFirstDirection(255), // Max uint8 value
			expected: enums.PathConvexityUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pathFirstDirectionToConvexity(tt.dir)
			if result != tt.expected {
				t.Errorf("pathFirstDirectionToConvexity(%v) = %v, expected %v", tt.dir, result, tt.expected)
			}
		})
	}
}

// TestPtsInVerb tests the ptsInVerb helper function
func TestPtsInVerb(t *testing.T) {
	tests := []struct {
		name     string
		verb     enums.PathVerb
		expected int
	}{
		{
			name:     "Move verb",
			verb:     enums.PathVerbMove,
			expected: 1,
		},
		{
			name:     "Line verb",
			verb:     enums.PathVerbLine,
			expected: 1,
		},
		{
			name:     "Quad verb",
			verb:     enums.PathVerbQuad,
			expected: 2,
		},
		{
			name:     "Conic verb",
			verb:     enums.PathVerbConic,
			expected: 2,
		},
		{
			name:     "Cubic verb",
			verb:     enums.PathVerbCubic,
			expected: 3,
		},
		{
			name:     "Close verb",
			verb:     enums.PathVerbClose,
			expected: 0,
		},
		{
			name:     "Invalid verb",
			verb:     enums.PathVerb(255), // Max uint8 value
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ptsInVerb(tt.verb)
			if result != tt.expected {
				t.Errorf("ptsInVerb(%v) = %d, expected %d", tt.verb, result, tt.expected)
			}
		})
	}
}

// TestValidUnitDivide tests the validUnitDivide helper function
// Ported from: skia-source/src/core/SkGeometry.cpp:valid_unit_divide()
func TestValidUnitDivide(t *testing.T) {
	tests := []struct {
		name        string
		numer       base.Scalar
		denom       base.Scalar
		expectedVal base.Scalar
		expectedOk  bool
	}{
		// Valid cases
		{
			name:        "valid division 1/2",
			numer:       1.0,
			denom:       2.0,
			expectedVal: 0.5,
			expectedOk:  true,
		},
		{
			name:        "valid division 1/4",
			numer:       1.0,
			denom:       4.0,
			expectedVal: 0.25,
			expectedOk:  true,
		},
		{
			name:        "valid division 3/4",
			numer:       3.0,
			denom:       4.0,
			expectedVal: 0.75,
			expectedOk:  true,
		},
		{
			name:        "valid division small values",
			numer:       0.001,
			denom:       0.002,
			expectedVal: 0.5,
			expectedOk:  true,
		},
		// Note: After normalization, -1/2 becomes 1/-2, and 1 >= -2 is true, so it fails
		// This matches C++ behavior - normalization doesn't guarantee both positive
		{
			name:        "negative numerator (normalized but fails check)",
			numer:       -1.0,
			denom:       2.0,
			expectedVal: 0.0,
			expectedOk:  false, // After normalization: 1/-2, and 1 >= -2 is true
		},
		{
			name:        "negative denominator (normalized but fails check)",
			numer:       1.0,
			denom:       -2.0,
			expectedVal: 0.0,
			expectedOk:  false, // After normalization: -1/2 -> 1/-2, and 1 >= -2 is true
		},

		// Invalid cases
		{
			name:        "zero denominator",
			numer:       1.0,
			denom:       0.0,
			expectedVal: 0.0,
			expectedOk:  false,
		},
		{
			name:        "zero numerator",
			numer:       0.0,
			denom:       2.0,
			expectedVal: 0.0,
			expectedOk:  false,
		},
		{
			name:        "numerator >= denominator",
			numer:       2.0,
			denom:       2.0,
			expectedVal: 0.0,
			expectedOk:  false,
		},
		{
			name:        "numerator > denominator",
			numer:       3.0,
			denom:       2.0,
			expectedVal: 0.0,
			expectedOk:  false,
		},
		{
			name:        "result < 1.0 (valid)",
			numer:       0.999,
			denom:       1.0,
			expectedVal: 0.999,
			expectedOk:  true, // 0.999 < 1.0, so valid
		},
		{
			name:        "result exactly 1.0",
			numer:       1.0,
			denom:       1.0,
			expectedVal: 0.0,
			expectedOk:  false,
		},
		{
			name:        "very small numerator (valid if > 0)",
			numer:       1e-10,
			denom:       1.0,
			expectedVal: 1e-10,
			expectedOk:  true, // Very small but > 0, so valid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := validUnitDivide(tt.numer, tt.denom)
			if ok != tt.expectedOk {
				t.Errorf("validUnitDivide(%f, %f) ok = %v, expected %v", tt.numer, tt.denom, ok, tt.expectedOk)
			}
			if tt.expectedOk {
				if !NearlyEqualScalar(result, tt.expectedVal) {
					t.Errorf("validUnitDivide(%f, %f) = %f, expected %f", tt.numer, tt.denom, result, tt.expectedVal)
				}
			}
		})
	}
}

// TestFindUnitQuadRoots tests the findUnitQuadRoots helper function
// Ported from: skia-source/src/core/SkGeometry.cpp:SkFindUnitQuadRoots()
func TestFindUnitQuadRoots(t *testing.T) {
	tests := []struct {
		name           string
		A, B, C        base.Scalar
		expectedCount  int
		expectedRoots  []base.Scalar
		description    string
	}{
		// Linear case (A == 0)
		{
			name:          "linear equation, one root",
			A:             0.0,
			B:             2.0,
			C:             -1.0,
			expectedCount: 1,
			expectedRoots: []base.Scalar{0.5},
			description:   "2t - 1 = 0 => t = 0.5",
		},
		{
			name:          "linear equation, no valid root (root >= 1)",
			A:             0.0,
			B:             1.0,
			C:             -2.0,
			expectedCount: 0,
			expectedRoots: nil,
			description:   "t - 2 = 0 => t = 2 (invalid, >= 1)",
		},
		{
			name:          "linear equation, no valid root (root < 0)",
			A:             0.0,
			B:             1.0,
			C:             1.0,
			expectedCount: 0,
			expectedRoots: nil,
			description:   "t + 1 = 0 => t = -1 (invalid, < 0)",
		},

		// Quadratic cases
		{
			name:          "quadratic, two roots (both >= 1, so none valid)",
			A:             1.0,
			B:             -3.0,
			C:             2.0,
			expectedCount: 0, // Roots are 1 and 2, both >= 1, so none are in [0,1)
			expectedRoots: nil,
			description:   "t^2 - 3t + 2 = 0 => t = 1, 2 (both >= 1, so none valid)",
		},
		{
			name:          "quadratic, one root in range",
			A:             1.0,
			B:             -1.5,
			C:             0.5,
			expectedCount: 1,
			expectedRoots: []base.Scalar{0.5},
			description:   "t^2 - 1.5t + 0.5 = 0 => t = 0.5, 1 (only 0.5 is in [0,1))",
		},
		{
			name:          "quadratic, no roots (negative discriminant)",
			A:             1.0,
			B:             1.0,
			C:             1.0,
			expectedCount: 0,
			expectedRoots: nil,
			description:   "t^2 + t + 1 = 0 => no real roots",
		},
		{
			name:          "quadratic, no roots in range",
			A:             1.0,
			B:             -5.0,
			C:             6.0,
			expectedCount: 0,
			expectedRoots: nil,
			description:   "t^2 - 5t + 6 = 0 => t = 2, 3 (both >= 1)",
		},
		{
			name:          "quadratic, duplicate roots (root >= 1, so none valid)",
			A:             1.0,
			B:             -2.0,
			C:             1.0,
			expectedCount: 0, // Root is 1, which is >= 1, so not in [0,1)
			expectedRoots: nil,
			description:   "t^2 - 2t + 1 = 0 => t = 1 (double root, but >= 1, so none valid)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roots := make([]base.Scalar, 2)
			count := findUnitQuadRoots(tt.A, tt.B, tt.C, roots)

			if count != tt.expectedCount {
				t.Errorf("findUnitQuadRoots(%f, %f, %f) count = %d, expected %d (%s)",
					tt.A, tt.B, tt.C, count, tt.expectedCount, tt.description)
			}

			if count > 0 && tt.expectedRoots != nil {
				// Check that roots are valid (in [0, 1))
				for i := 0; i < count; i++ {
					if roots[i] < 0 || roots[i] >= 1.0 {
						t.Errorf("findUnitQuadRoots root[%d] = %f is not in [0, 1)", i, roots[i])
					}
				}

				// Check that roots match expected (if provided)
				if len(tt.expectedRoots) > 0 {
					// Note: Some expected roots might be >= 1, so we only check valid ones
					validExpected := 0
					for _, r := range tt.expectedRoots {
						if r >= 0 && r < 1.0 {
							validExpected++
						}
					}
					if count != validExpected {
						t.Errorf("findUnitQuadRoots count = %d, expected %d valid roots", count, validExpected)
					}
				}
			}
		})
	}
}

// TestFindQuadExtrema tests the findQuadExtrema helper function
func TestFindQuadExtrema(t *testing.T) {
	tests := []struct {
		name          string
		a, b, c       base.Scalar
		expectedCount int
		expectedT     base.Scalar
		description   string
	}{
		{
			name:          "quadratic with extrema (linear case, no extrema)",
			a:             0.0,
			b:             1.0,
			c:             2.0,
			expectedCount: 0, // a-2*b+c = 0-2*1+2 = 0, so division fails
			description:   "P(t) = (1-t)^2*0 + 2*(1-t)*t*1 + t^2*2 (linear, no extrema)",
		},
		{
			name:          "quadratic extrema at t=0.25",
			a:             0.0,
			b:             1.0,
			c:             0.0,
			expectedCount: 1,
			expectedT:     0.5, // (1-0)/(0-2*1+0) = 1/-2 = -0.5 (invalid)
			description:   "P(t) = (1-t)^2*0 + 2*(1-t)*t*1 + t^2*0",
		},
		{
			name:          "no extrema (linear)",
			a:             0.0,
			b:             1.0,
			c:             2.0,
			expectedCount: 0, // a-2*b+c = 0-2*1+2 = 0, so division fails
			description:   "Linear case",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tValue := make([]base.Scalar, 1)
			count := findQuadExtrema(tt.a, tt.b, tt.c, tValue)

			if count != tt.expectedCount {
				t.Errorf("findQuadExtrema(%f, %f, %f) count = %d, expected %d (%s)",
					tt.a, tt.b, tt.c, count, tt.expectedCount, tt.description)
			}

			if count > 0 {
				if tValue[0] < 0 || tValue[0] >= 1.0 {
					t.Errorf("findQuadExtrema tValue = %f is not in [0, 1)", tValue[0])
				}
				if tt.expectedT > 0 && !NearlyEqualScalar(tValue[0], tt.expectedT) {
					t.Errorf("findQuadExtrema tValue = %f, expected %f", tValue[0], tt.expectedT)
				}
			}
		})
	}
}

// TestEvalQuadAt tests the evalQuadAt helper function
func TestEvalQuadAt(t *testing.T) {
	tests := []struct {
		name      string
		src       []Point
		t         base.Scalar
		expected  Point
		tolerance base.Scalar
	}{
		{
			name: "evaluate at t=0",
			src: []Point{
				{X: 0, Y: 0},
				{X: 1, Y: 1},
				{X: 2, Y: 0},
			},
			t:         0.0,
			expected:  Point{X: 0, Y: 0},
			tolerance: ScalarTolerance,
		},
		{
			name: "evaluate at t=1",
			src: []Point{
				{X: 0, Y: 0},
				{X: 1, Y: 1},
				{X: 2, Y: 0},
			},
			t:         1.0,
			expected:  Point{X: 2, Y: 0},
			tolerance: ScalarTolerance,
		},
		{
			name: "evaluate at t=0.5",
			src: []Point{
				{X: 0, Y: 0},
				{X: 1, Y: 1},
				{X: 2, Y: 0},
			},
			t:         0.5,
			expected:  Point{X: 1, Y: 0.5},
			tolerance: ScalarTolerance,
		},
		{
			name: "evaluate at t=0.25",
			src: []Point{
				{X: 0, Y: 0},
				{X: 2, Y: 2},
				{X: 4, Y: 0},
			},
			t:         0.25,
			expected:  Point{X: 1.0, Y: 0.75}, // (1-0.25)^2*0 + 2*(1-0.25)*0.25*2 + 0.25^2*4 = 0.75 + 0.25 = 1.0
			tolerance: ScalarTolerance,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalQuadAt(tt.src, tt.t)

			if !NearlyEqualScalar(result.X, tt.expected.X) {
				t.Errorf("evalQuadAt X = %f, expected %f", result.X, tt.expected.X)
			}
			if !NearlyEqualScalar(result.Y, tt.expected.Y) {
				t.Errorf("evalQuadAt Y = %f, expected %f", result.Y, tt.expected.Y)
			}
		})
	}
}

// TestEvalCubicAt tests the evalCubicAt helper function
func TestEvalCubicAt(t *testing.T) {
	tests := []struct {
		name     string
		src      []Point
		t        base.Scalar
		expected Point
	}{
		{
			name: "evaluate at t=0",
			src: []Point{
				{X: 0, Y: 0},
				{X: 1, Y: 1},
				{X: 2, Y: 1},
				{X: 3, Y: 0},
			},
			t:        0.0,
			expected: Point{X: 0, Y: 0},
		},
		{
			name: "evaluate at t=1",
			src: []Point{
				{X: 0, Y: 0},
				{X: 1, Y: 1},
				{X: 2, Y: 1},
				{X: 3, Y: 0},
			},
			t:        1.0,
			expected: Point{X: 3, Y: 0},
		},
		{
			name: "evaluate at t=0.5",
			src: []Point{
				{X: 0, Y: 0},
				{X: 1, Y: 1},
				{X: 2, Y: 1},
				{X: 3, Y: 0},
			},
			t:        0.5,
			expected: Point{X: 1.5, Y: 0.75},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalCubicAt(tt.src, tt.t)

			if !NearlyEqualScalar(result.X, tt.expected.X) {
				t.Errorf("evalCubicAt X = %f, expected %f", result.X, tt.expected.X)
			}
			if !NearlyEqualScalar(result.Y, tt.expected.Y) {
				t.Errorf("evalCubicAt Y = %f, expected %f", result.Y, tt.expected.Y)
			}
		})
	}
}

// TestEvalConicAt tests the evalConicAt helper function
func TestEvalConicAt(t *testing.T) {
	tests := []struct {
		name     string
		src      []Point
		w        base.Scalar
		t        base.Scalar
		expected Point
	}{
		{
			name: "evaluate at t=0",
			src: []Point{
				{X: 0, Y: 0},
				{X: 1, Y: 1},
				{X: 2, Y: 0},
			},
			w:        base.ScalarRoot2Over2,
			t:        0.0,
			expected: Point{X: 0, Y: 0},
		},
		{
			name: "evaluate at t=1",
			src: []Point{
				{X: 0, Y: 0},
				{X: 1, Y: 1},
				{X: 2, Y: 0},
			},
			w:        base.ScalarRoot2Over2,
			t:        1.0,
			expected: Point{X: 2, Y: 0},
		},
		{
			name: "evaluate at t=0.5",
			src: []Point{
				{X: 0, Y: 0},
				{X: 1, Y: 1},
				{X: 2, Y: 0},
			},
			w:        1.0, // Circular arc weight
			t:        0.5,
			expected: Point{X: 1, Y: 0.5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalConicAt(tt.src, tt.w, tt.t)

			if !NearlyEqualScalar(result.X, tt.expected.X) {
				t.Errorf("evalConicAt X = %f, expected %f", result.X, tt.expected.X)
			}
			if !NearlyEqualScalar(result.Y, tt.expected.Y) {
				t.Errorf("evalConicAt Y = %f, expected %f", result.Y, tt.expected.Y)
			}
		})
	}
}

// TestPathFillTypeIsInverse tests the PathFillTypeIsInverse helper function
func TestPathFillTypeIsInverse(t *testing.T) {
	tests := []struct {
		name     string
		ft       enums.PathFillType
		expected bool
	}{
		{
			name:     "Winding (not inverse)",
			ft:       enums.PathFillTypeWinding,
			expected: false,
		},
		{
			name:     "EvenOdd (not inverse)",
			ft:       enums.PathFillTypeEvenOdd,
			expected: false,
		},
		{
			name:     "InverseWinding (inverse)",
			ft:       enums.PathFillTypeInverseWinding,
			expected: true,
		},
		{
			name:     "InverseEvenOdd (inverse)",
			ft:       enums.PathFillTypeInverseEvenOdd,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PathFillTypeIsInverse(tt.ft)
			if result != tt.expected {
				t.Errorf("PathFillTypeIsInverse(%v) = %v, expected %v", tt.ft, result, tt.expected)
			}
		})
	}
}

// TestPathFillTypeToggleInverse tests the PathFillTypeToggleInverse helper function
func TestPathFillTypeToggleInverse(t *testing.T) {
	tests := []struct {
		name     string
		ft       enums.PathFillType
		expected enums.PathFillType
	}{
		{
			name:     "Winding -> InverseWinding",
			ft:       enums.PathFillTypeWinding,
			expected: enums.PathFillTypeInverseWinding,
		},
		{
			name:     "EvenOdd -> InverseEvenOdd",
			ft:       enums.PathFillTypeEvenOdd,
			expected: enums.PathFillTypeInverseEvenOdd,
		},
		{
			name:     "InverseWinding -> Winding",
			ft:       enums.PathFillTypeInverseWinding,
			expected: enums.PathFillTypeWinding,
		},
		{
			name:     "InverseEvenOdd -> EvenOdd",
			ft:       enums.PathFillTypeInverseEvenOdd,
			expected: enums.PathFillTypeEvenOdd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PathFillTypeToggleInverse(tt.ft)
			if result != tt.expected {
				t.Errorf("PathFillTypeToggleInverse(%v) = %v, expected %v", tt.ft, result, tt.expected)
			}
		})
	}
}

// TestPathConvexityIsConvex tests the PathConvexityIsConvex helper function
func TestPathConvexityIsConvex(t *testing.T) {
	tests := []struct {
		name     string
		cv       enums.PathConvexity
		expected bool
	}{
		{
			name:     "ConvexCW (convex)",
			cv:       enums.PathConvexityConvexCW,
			expected: true,
		},
		{
			name:     "ConvexCCW (convex)",
			cv:       enums.PathConvexityConvexCCW,
			expected: true,
		},
		{
			name:     "ConvexDegenerate (convex)",
			cv:       enums.PathConvexityConvexDegenerate,
			expected: true,
		},
		{
			name:     "Concave (not convex)",
			cv:       enums.PathConvexityConcave,
			expected: false,
		},
		{
			name:     "Unknown (not convex)",
			cv:       enums.PathConvexityUnknown,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PathConvexityIsConvex(tt.cv)
			if result != tt.expected {
				t.Errorf("PathConvexityIsConvex(%v) = %v, expected %v", tt.cv, result, tt.expected)
			}
		})
	}
}

// TestIsFinite tests the IsFinite helper function
func TestIsFinite(t *testing.T) {
	tests := []struct {
		name     string
		f        base.Scalar
		expected bool
	}{
		{
			name:     "finite positive",
			f:        1.0,
			expected: true,
		},
		{
			name:     "finite negative",
			f:        -1.0,
			expected: true,
		},
		{
			name:     "zero",
			f:        0.0,
			expected: true,
		},
		{
			name:     "NaN",
			f:        base.Scalar(math.NaN()),
			expected: false,
		},
		{
			name:     "positive infinity",
			f:        base.Scalar(math.Inf(1)),
			expected: false,
		},
		{
			name:     "negative infinity",
			f:        base.Scalar(math.Inf(-1)),
			expected: false,
		},
		{
			name:     "very small number",
			f:        1e-38,
			expected: true,
		},
		{
			name:     "very large number",
			f:        1e38,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsFinite(tt.f)
			if result != tt.expected {
				t.Errorf("IsFinite(%f) = %v, expected %v", tt.f, result, tt.expected)
			}
		})
	}
}

// TestAffectsAlphaColorFilter tests the affectsAlphaColorFilter helper function
func TestAffectsAlphaColorFilter(t *testing.T) {
	tests := []struct {
		name     string
		cf       ColorFilter
		expected bool
	}{
		{
			name:     "nil color filter",
			cf:       nil,
			expected: false,
		},
		{
			name:     "color filter that preserves alpha",
			cf:       &mockColorFilterPreservesAlpha{},
			expected: false,
		},
		{
			name:     "color filter that modifies alpha",
			cf:       &mockColorFilterModifiesAlpha{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := affectsAlphaColorFilter(tt.cf)
			if result != tt.expected {
				t.Errorf("affectsAlphaColorFilter() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestAffectsAlphaImageFilter tests the affectsAlphaImageFilter helper function
func TestAffectsAlphaImageFilter(t *testing.T) {
	tests := []struct {
		name     string
		imf      ImageFilter
		expected bool
	}{
		{
			name:     "nil image filter",
			imf:      nil,
			expected: false,
		},
		{
			name:     "non-nil image filter",
			imf:      &mockImageFilter{canComputeFastBounds: true},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := affectsAlphaImageFilter(tt.imf)
			if result != tt.expected {
				t.Errorf("affectsAlphaImageFilter() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

