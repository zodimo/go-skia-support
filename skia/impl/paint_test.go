package impl

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/base"
	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/models"
)

// mockMaskFilter is a simple mock implementation of MaskFilter for testing
type mockMaskFilter struct {
	name string
}

func (m *mockMaskFilter) ComputeFastBounds(bounds models.Rect, storage *models.Rect) {
	if storage != nil {
		// Expand bounds slightly for testing (simulating blur effect)
		*storage = bounds.MakeOutset(1, 1)
	}
}

// mockColorFilterPreservesAlpha is a mock ColorFilter that preserves alpha
type mockColorFilterPreservesAlpha struct{}

func (m *mockColorFilterPreservesAlpha) IsAlphaUnchanged() bool {
	return true // preserves alpha
}

// mockColorFilterModifiesAlpha is a mock ColorFilter that modifies alpha
type mockColorFilterModifiesAlpha struct{}

func (m *mockColorFilterModifiesAlpha) IsAlphaUnchanged() bool {
	return false // modifies alpha
}

// mockPathEffect is a mock PathEffect for testing
type mockPathEffect struct {
	modifiesBounds bool
}

func (m *mockPathEffect) ComputeFastBounds(bounds *models.Rect) bool {
	if m.modifiesBounds && bounds != nil {
		// Expand bounds slightly for testing
		*bounds = bounds.MakeOutset(2, 2)
	}
	return m.modifiesBounds
}

// mockImageFilter is a mock ImageFilter for testing
type mockImageFilter struct {
	canComputeFastBounds bool
}

func (m *mockImageFilter) CanComputeFastBounds() bool {
	return m.canComputeFastBounds
}

func (m *mockImageFilter) ComputeFastBounds(bounds models.Rect) models.Rect {
	// Expand bounds slightly for testing
	return bounds.MakeOutset(3, 3)
}

// TestPaint_Copy tests paint copying and equality
// Ported from: skia-source/tests/PaintTest.cpp:DEF_TEST(Paint_copy, reporter)
func TestPaint_Copy(t *testing.T) {
	// Create a paint and set a few member variables
	paint := NewPaint()
	paint.SetStyle(enums.PaintStyleStrokeAndFill)
	paint.SetStrokeWidth(base.Scalar(2))

	// Set a mask filter (using mock since we don't have full implementation yet)
	maskFilter := &mockMaskFilter{name: "test-blur"}
	paint.SetMaskFilter(maskFilter)

	// Copy the paint using assignment (Go equivalent of copy constructor)
	copiedPaint := &Paint{}
	*copiedPaint = *paint

	// Check they are the same
	if !paint.Equals(copiedPaint) {
		t.Error("Paint copy should be equal to original")
	}

	// Copy the paint using assignment again (Go equivalent of assignment operator)
	anotherCopy := &Paint{}
	*anotherCopy = *paint

	// Check they are the same
	if !paint.Equals(anotherCopy) {
		t.Error("Paint assignment copy should be equal to original")
	}

	// Create a clean paint and reset both paints
	cleanPaint := NewPaint()
	paint.Reset()
	copiedPaint.Reset()

	// Check they are back to their initial states
	if !cleanPaint.Equals(paint) {
		t.Error("Reset paint should equal clean paint")
	}
	if !cleanPaint.Equals(copiedPaint) {
		t.Error("Reset copied paint should equal clean paint")
	}
}

// TestPaint_CopyWithVariousProperties tests paint equality with various property combinations
func TestPaint_CopyWithVariousProperties(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*Paint)
	}{
		{
			name: "style only",
			setup: func(p *Paint) {
				p.SetStyle(enums.PaintStyleStroke)
			},
		},
		{
			name: "stroke width only",
			setup: func(p *Paint) {
				p.SetStrokeWidth(base.Scalar(5))
			},
		},
		{
			name: "style and stroke width",
			setup: func(p *Paint) {
				p.SetStyle(enums.PaintStyleStrokeAndFill)
				p.SetStrokeWidth(base.Scalar(3))
			},
		},
		{
			name: "color",
			setup: func(p *Paint) {
				p.SetColor(models.Color4f{R: 1.0, G: 0.5, B: 0.25, A: 1.0})
			},
		},
		{
			name: "alpha",
			setup: func(p *Paint) {
				p.SetAlpha(128)
			},
		},
		{
			name: "blend mode",
			setup: func(p *Paint) {
				p.SetBlendMode(enums.BlendModeMultiply)
			},
		},
		{
			name: "stroke cap",
			setup: func(p *Paint) {
				p.SetStrokeCap(enums.PaintCapRound)
			},
		},
		{
			name: "stroke join",
			setup: func(p *Paint) {
				p.SetStrokeJoin(enums.PaintJoinRound)
			},
		},
		{
			name: "miter limit",
			setup: func(p *Paint) {
				p.SetStrokeMiter(base.Scalar(8))
			},
		},
		{
			name: "anti alias",
			setup: func(p *Paint) {
				p.SetAntiAlias(true)
			},
		},
		{
			name: "dither",
			setup: func(p *Paint) {
				p.SetDither(true)
			},
		},
		{
			name: "mask filter",
			setup: func(p *Paint) {
				p.SetMaskFilter(&mockMaskFilter{name: "test"})
			},
		},
		{
			name: "multiple properties",
			setup: func(p *Paint) {
				p.SetStyle(enums.PaintStyleStrokeAndFill)
				p.SetStrokeWidth(base.Scalar(2))
				p.SetColor(models.Color4f{R: 1.0, G: 0.0, B: 0.0, A: 0.5})
				p.SetBlendMode(enums.BlendModeSrcOver)
				p.SetStrokeCap(enums.PaintCapRound)
				p.SetStrokeJoin(enums.PaintJoinRound)
				p.SetAntiAlias(true)
				p.SetMaskFilter(&mockMaskFilter{name: "test-blur"})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create original paint and apply setup
			original := NewPaint()
			tt.setup(original)

			// Copy the paint
			copied := &Paint{}
			*copied = *original

			// Verify equality
			if !original.Equals(copied) {
				t.Errorf("Copied paint should equal original for %s", tt.name)
			}

			// Verify reset works
			original.Reset()
			copied.Reset()
			cleanPaint := NewPaint()

			if !original.Equals(cleanPaint) {
				t.Errorf("Reset original should equal clean paint for %s", tt.name)
			}
			if !copied.Equals(cleanPaint) {
				t.Errorf("Reset copied should equal clean paint for %s", tt.name)
			}
		})
	}
}

// TestPaint_NothingToDraw tests the NothingToDraw() method
// Ported from: skia-source/tests/PaintTest.cpp:DEF_TEST(Paint_nothingToDraw, r)
func TestPaint_NothingToDraw(t *testing.T) {
	paint := NewPaint()

	// Test default paint: NothingToDraw() = false
	if paint.NothingToDraw() {
		t.Error("Default paint should not have nothingToDraw() = true")
	}

	// Test zero alpha: NothingToDraw() = true
	paint.SetAlpha(0)
	if !paint.NothingToDraw() {
		t.Error("Paint with zero alpha should have nothingToDraw() = true")
	}

	// Test Dst blend mode: NothingToDraw() = true
	paint.SetAlpha(0xFF)
	paint.SetBlendMode(enums.BlendModeDst)
	if !paint.NothingToDraw() {
		t.Error("Paint with Dst blend mode should have nothingToDraw() = true")
	}

	// Test color filter that preserves alpha: NothingToDraw() = true (when alpha is 0)
	paint.SetAlpha(0)
	paint.SetBlendMode(enums.BlendModeSrcOver)
	paint.SetColorFilter(&mockColorFilterPreservesAlpha{})
	if !paint.NothingToDraw() {
		t.Error("Paint with zero alpha and color filter that preserves alpha should have nothingToDraw() = true")
	}

	// Test color filter that modifies alpha: NothingToDraw() = false (even when alpha is 0)
	paint.SetColorFilter(&mockColorFilterModifiesAlpha{})
	if paint.NothingToDraw() {
		t.Error("Paint with zero alpha but color filter that modifies alpha should have nothingToDraw() = false")
	}
}

// TestPaint_ComputeFastBounds tests bounds computation with various paint properties
// Based on: skia-source/src/core/SkPaint.cpp:computeFastBounds()
func TestPaint_ComputeFastBounds(t *testing.T) {
	origRect := models.Rect{Left: 10, Top: 20, Right: 30, Bottom: 40}
	storage := &models.Rect{}

	// Test 1: Fill style with no effects (fast path)
	paint := NewPaint()
	paint.SetStyle(enums.PaintStyleFill)
	result := paint.ComputeFastBounds(origRect, storage)
	if result.Left != origRect.Left || result.Top != origRect.Top ||
		result.Right != origRect.Right || result.Bottom != origRect.Bottom {
		t.Error("Fill style with no effects should return original bounds")
	}

	// Test 2: Fill style with mask filter
	paint = NewPaint()
	paint.SetStyle(enums.PaintStyleFill)
	paint.SetMaskFilter(&mockMaskFilter{name: "test"})
	result = paint.ComputeFastBounds(origRect, storage)
	// Bounds should be modified by mask filter
	if result.Left == origRect.Left && result.Top == origRect.Top &&
		result.Right == origRect.Right && result.Bottom == origRect.Bottom {
		t.Error("Fill style with mask filter should modify bounds")
	}

	// Test 3: Stroke style with stroke width
	paint = NewPaint()
	paint.SetStyle(enums.PaintStyleStroke)
	paint.SetStrokeWidth(base.Scalar(4))
	result = paint.ComputeFastBounds(origRect, storage)
	// Bounds should be expanded by stroke width
	if result.Left >= origRect.Left || result.Top >= origRect.Top ||
		result.Right <= origRect.Right || result.Bottom <= origRect.Bottom {
		t.Error("Stroke style should expand bounds")
	}

	// Test 4: Stroke style with zero width (hairline)
	paint = NewPaint()
	paint.SetStyle(enums.PaintStyleStroke)
	paint.SetStrokeWidth(0)
	result = paint.ComputeFastBounds(origRect, storage)
	// Hairline stroke should still expand bounds (default 1.0 radius)
	if result.Left >= origRect.Left || result.Top >= origRect.Top ||
		result.Right <= origRect.Right || result.Bottom <= origRect.Bottom {
		t.Error("Hairline stroke should expand bounds")
	}

	// Test 5: Stroke style with path effect
	paint = NewPaint()
	paint.SetStyle(enums.PaintStyleStroke)
	paint.SetStrokeWidth(base.Scalar(2))
	paint.SetPathEffect(&mockPathEffect{modifiesBounds: true})
	result = paint.ComputeFastBounds(origRect, storage)
	// Bounds should be modified by path effect
	if result.Left >= origRect.Left || result.Top >= origRect.Top {
		t.Error("Path effect should modify bounds")
	}

	// Test 6: Fill style with image filter
	paint = NewPaint()
	paint.SetStyle(enums.PaintStyleFill)
	paint.SetImageFilter(&mockImageFilter{canComputeFastBounds: true})
	result = paint.ComputeFastBounds(origRect, storage)
	// Bounds should be modified by image filter
	if result.Left >= origRect.Left || result.Top >= origRect.Top {
		t.Error("Image filter should modify bounds")
	}

	// Test 7: StrokeAndFill style
	paint = NewPaint()
	paint.SetStyle(enums.PaintStyleStrokeAndFill)
	paint.SetStrokeWidth(base.Scalar(3))
	result = paint.ComputeFastBounds(origRect, storage)
	// Bounds should be expanded
	if result.Left >= origRect.Left || result.Top >= origRect.Top {
		t.Error("StrokeAndFill style should expand bounds")
	}
}

// TestPaint_ComputeFastStrokeBounds tests stroke bounds computation
func TestPaint_ComputeFastStrokeBounds(t *testing.T) {
	origRect := models.Rect{Left: 10, Top: 20, Right: 30, Bottom: 40}
	storage := &models.Rect{}

	// Test stroke bounds regardless of paint style
	paint := NewPaint()
	paint.SetStyle(enums.PaintStyleFill) // Set to fill, but ComputeFastStrokeBounds should use stroke
	paint.SetStrokeWidth(base.Scalar(5))
	result := paint.ComputeFastStrokeBounds(origRect, storage)
	// Should compute stroke bounds even though style is fill
	if result.Left >= origRect.Left || result.Top >= origRect.Top {
		t.Error("ComputeFastStrokeBounds should compute stroke bounds regardless of style")
	}
}
