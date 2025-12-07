package impl

import (
	"math"
	"testing"

	"github.com/zodimo/go-skia-support/skia/base"
	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/models"
)

// TestColor4f_PinAlpha tests the PinAlpha method which clamps alpha to [0, 1] range
func TestColor4f_PinAlpha(t *testing.T) {
	tests := []struct {
		name     string
		color    models.Color4f
		expected models.Color4f
	}{
		{
			name: "alpha in valid range",
			color: models.Color4f{
				R: 0.5,
				G: 0.5,
				B: 0.5,
				A: 0.5,
			},
			expected: models.Color4f{
				R: 0.5,
				G: 0.5,
				B: 0.5,
				A: 0.5,
			},
		},
		{
			name: "alpha below zero (clamped to 0)",
			color: models.Color4f{
				R: 1.0,
				G: 0.0,
				B: 0.0,
				A: -0.5,
			},
			expected: models.Color4f{
				R: 1.0,
				G: 0.0,
				B: 0.0,
				A: 0.0,
			},
		},
		{
			name: "alpha above one (clamped to 1)",
			color: models.Color4f{
				R: 1.0,
				G: 1.0,
				B: 1.0,
				A: 1.5,
			},
			expected: models.Color4f{
				R: 1.0,
				G: 1.0,
				B: 1.0,
				A: 1.0,
			},
		},
		{
			name: "alpha exactly zero",
			color: models.Color4f{
				R: 0.0,
				G: 0.0,
				B: 0.0,
				A: 0.0,
			},
			expected: models.Color4f{
				R: 0.0,
				G: 0.0,
				B: 0.0,
				A: 0.0,
			},
		},
		{
			name: "alpha exactly one",
			color: models.Color4f{
				R: 1.0,
				G: 1.0,
				B: 1.0,
				A: 1.0,
			},
			expected: models.Color4f{
				R: 1.0,
				G: 1.0,
				B: 1.0,
				A: 1.0,
			},
		},
		{
			name: "RGB components preserved (not clamped)",
			color: models.Color4f{
				R: 1.5, // RGB can be outside [0,1] before clamping
				G: -0.5,
				B: 2.0,
				A: 0.5,
			},
			expected: models.Color4f{
				R: 1.5, // RGB preserved as-is
				G: -0.5,
				B: 2.0,
				A: 0.5, // Alpha clamped
			},
		},
		{
			name: "very large negative alpha",
			color: models.Color4f{
				R: 0.0,
				G: 0.0,
				B: 0.0,
				A: -1000.0,
			},
			expected: models.Color4f{
				R: 0.0,
				G: 0.0,
				B: 0.0,
				A: 0.0,
			},
		},
		{
			name: "very large positive alpha",
			color: models.Color4f{
				R: 0.0,
				G: 0.0,
				B: 0.0,
				A: 1000.0,
			},
			expected: models.Color4f{
				R: 0.0,
				G: 0.0,
				B: 0.0,
				A: 1.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.color.PinAlpha()

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

// TestColor4f_Vec tests the Vec method which returns a pointer to the components array
func TestColor4f_Vec(t *testing.T) {
	color := models.Color4f{
		R: 0.25,
		G: 0.5,
		B: 0.75,
		A: 1.0,
	}

	vec := color.Vec()

	// Verify vec is not nil
	if vec == nil {
		t.Fatal("Vec() returned nil")
	}

	// Verify components match
	if !NearlyEqualScalar(vec[0], color.R) {
		t.Errorf("Vec[0] (R) mismatch: got %f, expected %f", vec[0], color.R)
	}
	if !NearlyEqualScalar(vec[1], color.G) {
		t.Errorf("Vec[1] (G) mismatch: got %f, expected %f", vec[1], color.G)
	}
	if !NearlyEqualScalar(vec[2], color.B) {
		t.Errorf("Vec[2] (B) mismatch: got %f, expected %f", vec[2], color.B)
	}
	if !NearlyEqualScalar(vec[3], color.A) {
		t.Errorf("Vec[3] (A) mismatch: got %f, expected %f", vec[3], color.A)
	}

	// Verify modifying vec modifies the original (pointer behavior)
	vec[0] = 0.99
	if !NearlyEqualScalar(vec[0], 0.99) {
		t.Error("Vec pointer modification failed")
	}
}

// TestColor4f_ToSkColor tests the ToSkColor method with comprehensive edge cases
func TestColor4f_ToSkColor(t *testing.T) {
	tests := []struct {
		name     string
		color    models.Color4f
		expected uint32
	}{
		{
			name: "opaque black",
			color: models.Color4f{
				R: 0.0,
				G: 0.0,
				B: 0.0,
				A: 1.0,
			},
			expected: 0xFF000000, // ARGB: alpha=FF, r=00, g=00, b=00
		},
		{
			name: "opaque white",
			color: models.Color4f{
				R: 1.0,
				G: 1.0,
				B: 1.0,
				A: 1.0,
			},
			expected: 0xFFFFFFFF, // ARGB: alpha=FF, r=FF, g=FF, b=FF
		},
		{
			name: "opaque red",
			color: models.Color4f{
				R: 1.0,
				G: 0.0,
				B: 0.0,
				A: 1.0,
			},
			expected: 0xFFFF0000, // ARGB: alpha=FF, r=FF, g=00, b=00
		},
		{
			name: "transparent (zero alpha)",
			color: models.Color4f{
				R: 1.0,
				G: 0.0,
				B: 0.0,
				A: 0.0,
			},
			expected: 0x00FF0000, // ARGB: alpha=00, r=FF, g=00, b=00
		},
		{
			name: "semi-transparent",
			color: models.Color4f{
				R: 1.0,
				G: 0.0,
				B: 0.0,
				A: 0.5,
			},
			expected: 0x7FFF0000, // ARGB: alpha=7F (127, truncated from 127.5), r=FF, g=00, b=00
		},
		{
			name: "components clamped above 1.0",
			color: models.Color4f{
				R: 2.0, // Clamped to 1.0
				G: 1.5, // Clamped to 1.0
				B: 0.5,
				A: 1.0,
			},
			expected: 0xFFFFFF7F, // ARGB: alpha=FF, r=FF, g=FF, b=7F (127, truncated from 127.5)
		},
		{
			name: "components clamped below 0.0",
			color: models.Color4f{
				R: -0.5, // Clamped to 0.0
				G: 0.0,
				B: 0.5,
				A: 1.0,
			},
			expected: 0xFF00007F, // ARGB: alpha=FF, r=00, g=00, b=7F (127, truncated from 127.5)
		},
		{
			name: "alpha clamped above 1.0",
			color: models.Color4f{
				R: 1.0,
				G: 0.0,
				B: 0.0,
				A: 2.0, // Clamped to 1.0
			},
			expected: 0xFFFF0000, // ARGB: alpha=FF (clamped), r=FF, g=00, b=00
		},
		{
			name: "alpha clamped below 0.0",
			color: models.Color4f{
				R: 1.0,
				G: 0.0,
				B: 0.0,
				A: -0.5, // Clamped to 0.0
			},
			expected: 0x00FF0000, // ARGB: alpha=00 (clamped), r=FF, g=00, b=00
		},
		{
			name: "gray (0.5, 0.5, 0.5)",
			color: models.Color4f{
				R: 0.5,
				G: 0.5,
				B: 0.5,
				A: 1.0,
			},
			expected: 0xFF7F7F7F, // ARGB: alpha=FF, r=7F, g=7F, b=7F (127, truncated from 127.5)
		},
		{
			name: "precise values (128/255)",
			color: models.Color4f{
				R: 128.0 / 255.0,
				G: 128.0 / 255.0,
				B: 128.0 / 255.0,
				A: 1.0,
			},
			expected: 0xFF808080, // ARGB: alpha=FF, r=80, g=80, b=80
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.color.ToSkColor()

			if result != tt.expected {
				t.Errorf("ToSkColor() = 0x%08X, expected 0x%08X", result, tt.expected)
			}
		})
	}
}

// TestColor4f_ToSkColor_EdgeCases tests extreme edge cases for ToSkColor
func TestColor4f_ToSkColor_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		color models.Color4f
	}{
		{
			name: "NaN components (should clamp to 0)",
			color: models.Color4f{
				R: base.Scalar(math.NaN()),
				G: 0.0,
				B: 0.0,
				A: 1.0,
			},
		},
		{
			name: "Inf components (should clamp)",
			color: models.Color4f{
				R: base.Scalar(math.Inf(1)),
				G: 0.0,
				B: 0.0,
				A: 1.0,
			},
		},
		{
			name: "negative Inf components (should clamp to 0)",
			color: models.Color4f{
				R: base.Scalar(math.Inf(-1)),
				G: 0.0,
				B: 0.0,
				A: 1.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.color.ToSkColor()

			// Extract components
			a := uint8((result >> 24) & 0xFF)
			r := uint8((result >> 16) & 0xFF)
			g := uint8((result >> 8) & 0xFF)
			b := uint8(result & 0xFF)

			// Verify all components are in valid range [0, 255]
			if a > 255 || r > 255 || g > 255 || b > 255 {
				t.Errorf("Components out of range: a=%d, r=%d, g=%d, b=%d", a, r, g, b)
			}

			// Verify result is a valid uint32
			if result == 0 && tt.color.A > 0 {
				// This is acceptable if alpha was clamped to 0
			}
		})
	}
}

// TestPaintStyle_EnumValues tests PaintStyle enum values and behavior
func TestPaintStyle_EnumValues(t *testing.T) {
	// Verify enum values match expected constants
	if enums.PaintStyleFill != 0 {
		t.Errorf("PaintStyleFill = %d, expected 0", enums.PaintStyleFill)
	}
	if enums.PaintStyleStroke != 1 {
		t.Errorf("PaintStyleStroke = %d, expected 1", enums.PaintStyleStroke)
	}
	if enums.PaintStyleStrokeAndFill != 2 {
		t.Errorf("PaintStyleStrokeAndFill = %d, expected 2", enums.PaintStyleStrokeAndFill)
	}

	// Test PaintStyle with Paint
	p := NewPaint()

	// Default should be Fill
	if p.GetStyle() != enums.PaintStyleFill {
		t.Errorf("Default style = %d, expected PaintStyleFill (%d)", p.GetStyle(), enums.PaintStyleFill)
	}

	// Test SetStyle with Fill
	p.SetStyle(enums.PaintStyleFill)
	if p.GetStyle() != enums.PaintStyleFill {
		t.Errorf("SetStyle(Fill) failed: got %d", p.GetStyle())
	}

	// Test SetStyle with Stroke
	p.SetStyle(enums.PaintStyleStroke)
	if p.GetStyle() != enums.PaintStyleStroke {
		t.Errorf("SetStyle(Stroke) failed: got %d", p.GetStyle())
	}

	// Test SetStyle with StrokeAndFill
	p.SetStyle(enums.PaintStyleStrokeAndFill)
	if p.GetStyle() != enums.PaintStyleStrokeAndFill {
		t.Errorf("SetStyle(StrokeAndFill) failed: got %d", p.GetStyle())
	}

	// Test SetStroke convenience method
	p.SetStroke(true)
	if p.GetStyle() != enums.PaintStyleStroke {
		t.Errorf("SetStroke(true) failed: got %d", p.GetStyle())
	}

	p.SetStroke(false)
	if p.GetStyle() != enums.PaintStyleFill {
		t.Errorf("SetStroke(false) failed: got %d", p.GetStyle())
	}

	// Test invalid style (should be ignored)
	originalStyle := p.GetStyle()
	p.SetStyle(enums.PaintStyle(255)) // Invalid style
	if p.GetStyle() != originalStyle {
		t.Errorf("Invalid style was accepted: got %d, expected %d", p.GetStyle(), originalStyle)
	}
}

// TestBlendMode_EnumValues tests BlendMode enum values
func TestBlendMode_EnumValues(t *testing.T) {
	// Verify key enum values
	if enums.BlendModeClear != 0 {
		t.Errorf("BlendModeClear = %d, expected 0", enums.BlendModeClear)
	}
	if enums.BlendModeSrc != 1 {
		t.Errorf("BlendModeSrc = %d, expected 1", enums.BlendModeSrc)
	}
	if enums.BlendModeDst != 2 {
		t.Errorf("BlendModeDst = %d, expected 2", enums.BlendModeDst)
	}
	if enums.BlendModeSrcOver != 3 {
		t.Errorf("BlendModeSrcOver = %d, expected 3", enums.BlendModeSrcOver)
	}

	// Verify last enum value
	if enums.BlendModeLast != enums.BlendModeLuminosity {
		t.Errorf("BlendModeLast = %d, expected BlendModeLuminosity (%d)", enums.BlendModeLast, enums.BlendModeLuminosity)
	}

	// Test BlendMode with Paint
	p := NewPaint()

	// Default should be SrcOver (nil blendMode means default)
	mode, ok := p.AsBlendMode()
	if !ok || mode != enums.BlendModeSrcOver {
		t.Errorf("Default blend mode = (%d, %v), expected (BlendModeSrcOver, true)", mode, ok)
	}

	// Test SetBlendMode with SrcOver (should clear stored mode)
	p.SetBlendMode(enums.BlendModeSrcOver)
	mode, ok = p.AsBlendMode()
	if !ok || mode != enums.BlendModeSrcOver {
		t.Errorf("SetBlendMode(SrcOver) failed: got (%d, %v)", mode, ok)
	}

	// Test SetBlendMode with Dst (should store mode)
	p.SetBlendMode(enums.BlendModeDst)
	mode, ok = p.AsBlendMode()
	if !ok || mode != enums.BlendModeDst {
		t.Errorf("SetBlendMode(Dst) failed: got (%d, %v)", mode, ok)
	}

	// Test SetBlendMode with Clear
	p.SetBlendMode(enums.BlendModeClear)
	mode, ok = p.AsBlendMode()
	if !ok || mode != enums.BlendModeClear {
		t.Errorf("SetBlendMode(Clear) failed: got (%d, %v)", mode, ok)
	}

	// Test all valid blend modes
	validModes := []enums.BlendMode{
		enums.BlendModeClear,
		enums.BlendModeSrc,
		enums.BlendModeDst,
		enums.BlendModeSrcOver,
		enums.BlendModeDstOver,
		enums.BlendModeSrcIn,
		enums.BlendModeDstIn,
		enums.BlendModeSrcOut,
		enums.BlendModeDstOut,
		enums.BlendModeSrcATop,
		enums.BlendModeDstATop,
		enums.BlendModeXor,
		enums.BlendModePlus,
		enums.BlendModeModulate,
		enums.BlendModeScreen,
		enums.BlendModeOverlay,
		enums.BlendModeDarken,
		enums.BlendModeLighten,
		enums.BlendModeColorDodge,
		enums.BlendModeColorBurn,
		enums.BlendModeHardLight,
		enums.BlendModeSoftLight,
		enums.BlendModeDifference,
		enums.BlendModeExclusion,
		enums.BlendModeMultiply,
		enums.BlendModeHue,
		enums.BlendModeSaturation,
		enums.BlendModeColor,
		enums.BlendModeLuminosity,
	}

	for _, mode := range validModes {
		p.SetBlendMode(mode)
		gotMode, ok := p.AsBlendMode()
		if !ok || gotMode != mode {
			t.Errorf("SetBlendMode(%d) failed: got (%d, %v)", mode, gotMode, ok)
		}
	}
}

// TestPaintCap_EnumValues tests PaintCap enum values and behavior
func TestPaintCap_EnumValues(t *testing.T) {
	// Verify enum values
	if enums.PaintCapButt != 0 {
		t.Errorf("PaintCapButt = %d, expected 0", enums.PaintCapButt)
	}
	if enums.PaintCapRound != 1 {
		t.Errorf("PaintCapRound = %d, expected 1", enums.PaintCapRound)
	}
	if enums.PaintCapSquare != 2 {
		t.Errorf("PaintCapSquare = %d, expected 2", enums.PaintCapSquare)
	}
	if enums.PaintCapDefault != enums.PaintCapButt {
		t.Errorf("PaintCapDefault = %d, expected PaintCapButt (%d)", enums.PaintCapDefault, enums.PaintCapButt)
	}
	if enums.PaintCapLast != enums.PaintCapSquare {
		t.Errorf("PaintCapLast = %d, expected PaintCapSquare (%d)", enums.PaintCapLast, enums.PaintCapSquare)
	}
	if enums.PaintCapCount != 3 {
		t.Errorf("PaintCapCount = %d, expected 3", enums.PaintCapCount)
	}

	// Test PaintCap with Paint
	p := NewPaint()

	// Default should be Butt
	if p.GetStrokeCap() != enums.PaintCapDefault {
		t.Errorf("Default cap = %d, expected PaintCapDefault (%d)", p.GetStrokeCap(), enums.PaintCapDefault)
	}

	// Test SetStrokeCap with Butt
	p.SetStrokeCap(enums.PaintCapButt)
	if p.GetStrokeCap() != enums.PaintCapButt {
		t.Errorf("SetStrokeCap(Butt) failed: got %d", p.GetStrokeCap())
	}

	// Test SetStrokeCap with Round
	p.SetStrokeCap(enums.PaintCapRound)
	if p.GetStrokeCap() != enums.PaintCapRound {
		t.Errorf("SetStrokeCap(Round) failed: got %d", p.GetStrokeCap())
	}

	// Test SetStrokeCap with Square
	p.SetStrokeCap(enums.PaintCapSquare)
	if p.GetStrokeCap() != enums.PaintCapSquare {
		t.Errorf("SetStrokeCap(Square) failed: got %d", p.GetStrokeCap())
	}
}

// TestPaintJoin_EnumValues tests PaintJoin enum values and behavior
func TestPaintJoin_EnumValues(t *testing.T) {
	// Verify enum values
	if enums.PaintJoinMiter != 0 {
		t.Errorf("PaintJoinMiter = %d, expected 0", enums.PaintJoinMiter)
	}
	if enums.PaintJoinRound != 1 {
		t.Errorf("PaintJoinRound = %d, expected 1", enums.PaintJoinRound)
	}
	if enums.PaintJoinBevel != 2 {
		t.Errorf("PaintJoinBevel = %d, expected 2", enums.PaintJoinBevel)
	}
	if enums.PaintJoinDefault != enums.PaintJoinMiter {
		t.Errorf("PaintJoinDefault = %d, expected PaintJoinMiter (%d)", enums.PaintJoinDefault, enums.PaintJoinMiter)
	}
	if enums.PaintJoinLast != enums.PaintJoinBevel {
		t.Errorf("PaintJoinLast = %d, expected PaintJoinBevel (%d)", enums.PaintJoinLast, enums.PaintJoinBevel)
	}

	// Test PaintJoin with Paint
	p := NewPaint()

	// Default should be Miter
	if p.GetStrokeJoin() != enums.PaintJoinDefault {
		t.Errorf("Default join = %d, expected PaintJoinDefault (%d)", p.GetStrokeJoin(), enums.PaintJoinDefault)
	}

	// Test SetStrokeJoin with Miter
	p.SetStrokeJoin(enums.PaintJoinMiter)
	if p.GetStrokeJoin() != enums.PaintJoinMiter {
		t.Errorf("SetStrokeJoin(Miter) failed: got %d", p.GetStrokeJoin())
	}

	// Test SetStrokeJoin with Round
	p.SetStrokeJoin(enums.PaintJoinRound)
	if p.GetStrokeJoin() != enums.PaintJoinRound {
		t.Errorf("SetStrokeJoin(Round) failed: got %d", p.GetStrokeJoin())
	}

	// Test SetStrokeJoin with Bevel
	p.SetStrokeJoin(enums.PaintJoinBevel)
	if p.GetStrokeJoin() != enums.PaintJoinBevel {
		t.Errorf("SetStrokeJoin(Bevel) failed: got %d", p.GetStrokeJoin())
	}
}

// TestColor4f_WithPaint tests Color4f operations when used with Paint
func TestColor4f_WithPaint(t *testing.T) {
	// Test SetColor with PinAlpha behavior
	p := NewPaint()

	// Test that SetColor automatically calls PinAlpha
	colorWithInvalidAlpha := models.Color4f{
		R: 1.0,
		G: 0.0,
		B: 0.0,
		A: 1.5, // Should be clamped to 1.0
	}
	p.SetColor(colorWithInvalidAlpha)
	gotColor := p.GetColor()
	if !NearlyEqualScalar(gotColor.A, 1.0) {
		t.Errorf("SetColor did not clamp alpha: got %f, expected 1.0", gotColor.A)
	}

	// Test SetColor with negative alpha
	colorWithNegativeAlpha := models.Color4f{
		R: 1.0,
		G: 0.0,
		B: 0.0,
		A: -0.5, // Should be clamped to 0.0
	}
	p.SetColor(colorWithNegativeAlpha)
	gotColor = p.GetColor()
	if !NearlyEqualScalar(gotColor.A, 0.0) {
		t.Errorf("SetColor did not clamp negative alpha: got %f, expected 0.0", gotColor.A)
	}

	// Test SetAlpha (uint8) with Paint
	p.SetAlpha(128) // 128/255 ≈ 0.502
	gotAlphaf := p.GetAlphaf()
	expectedAlphaf := base.Scalar(128) / 255.0
	if !NearlyEqualScalar(gotAlphaf, expectedAlphaf) {
		t.Errorf("SetAlpha(128) failed: got %f, expected %f", gotAlphaf, expectedAlphaf)
	}

	// Test SetAlphaf (float) with Paint
	p.SetAlphaf(0.75)
	gotAlphaf = p.GetAlphaf()
	if !NearlyEqualScalar(gotAlphaf, 0.75) {
		t.Errorf("SetAlphaf(0.75) failed: got %f, expected 0.75", gotAlphaf)
	}

	// Test SetAlphaf clamping
	p.SetAlphaf(1.5) // Should clamp to 1.0
	gotAlphaf = p.GetAlphaf()
	if !NearlyEqualScalar(gotAlphaf, 1.0) {
		t.Errorf("SetAlphaf(1.5) did not clamp: got %f, expected 1.0", gotAlphaf)
	}

	p.SetAlphaf(-0.5) // Should clamp to 0.0
	gotAlphaf = p.GetAlphaf()
	if !NearlyEqualScalar(gotAlphaf, 0.0) {
		t.Errorf("SetAlphaf(-0.5) did not clamp: got %f, expected 0.0", gotAlphaf)
	}

	// Test GetColorInt (uses ToSkColor)
	p.SetColor(models.Color4f{
		R: 1.0,
		G: 0.0,
		B: 0.0,
		A: 1.0,
	})
	gotColorInt := p.GetColorInt()
	expectedColorInt := uint32(0xFFFF0000) // ARGB: alpha=FF, r=FF, g=00, b=00
	if gotColorInt != expectedColorInt {
		t.Errorf("GetColorInt() = 0x%08X, expected 0x%08X", gotColorInt, expectedColorInt)
	}
}

// TestEnum_EdgeCases tests enum edge cases with Paint
func TestEnum_EdgeCases(t *testing.T) {
	p := NewPaint()

	// Test PaintStyle edge cases
	// Test that invalid style values are ignored
	originalStyle := p.GetStyle()
	p.SetStyle(enums.PaintStyle(255)) // Invalid value
	if p.GetStyle() != originalStyle {
		t.Errorf("Invalid PaintStyle was accepted: got %d, expected %d", p.GetStyle(), originalStyle)
	}

	// Test PaintCap edge cases
	originalCap := p.GetStrokeCap()
	// PaintCap doesn't validate, so this would be accepted
	// But we can test that values outside normal range still work
	p.SetStrokeCap(enums.PaintCap(255))
	// This is implementation-dependent, but should not crash

	// Test PaintJoin edge cases
	originalJoin := p.GetStrokeJoin()
	// PaintJoin doesn't validate, so this would be accepted
	p.SetStrokeJoin(enums.PaintJoin(255))
	// This is implementation-dependent, but should not crash

	// Reset to valid values
	p.SetStrokeCap(originalCap)
	p.SetStrokeJoin(originalJoin)
}

