package impl

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/models"
)

// TestPath_Bounds tests path bounds calculation
// Ported from: skia-source/tests/PathTest.cpp:test_bounds()
func TestPath_Bounds(t *testing.T) {
	// Test empty path bounds
	path := NewSkPath(enums.PathFillTypeDefault)
	bounds := path.Bounds()
	expected := models.Rect{Left: 0, Top: 0, Right: 0, Bottom: 0}
	if bounds != expected {
		t.Errorf("Empty path bounds: got %v, want %v", bounds, expected)
	}

	// Test single point bounds
	path = NewSkPath(enums.PathFillTypeDefault)
	path.MoveTo(10, 20)
	bounds = path.Bounds()
	expected = models.Rect{Left: 10, Top: 20, Right: 10, Bottom: 20}
	if bounds != expected {
		t.Errorf("Single point bounds: got %v, want %v", bounds, expected)
	}

	// Test line bounds
	path = NewSkPath(enums.PathFillTypeDefault)
	path.MoveTo(10, 20)
	path.LineTo(30, 40)
	bounds = path.Bounds()
	expected = models.Rect{Left: 10, Top: 20, Right: 30, Bottom: 40}
	if bounds.Left != expected.Left || bounds.Top != expected.Top ||
		bounds.Right != expected.Right || bounds.Bottom != expected.Bottom {
		t.Errorf("Line bounds: got %v, want %v", bounds, expected)
	}

	// Test curve bounds (quadratic)
	path = NewSkPath(enums.PathFillTypeDefault)
	path.MoveTo(0, 0)
	path.QuadTo(10, 20, 20, 0)
	bounds = path.Bounds()
	// Bounds should include control point (10, 20)
	if bounds.Top > 20 || bounds.Top < 0 {
		t.Errorf("Quadratic curve bounds top: got %v, should include control point y=20", bounds.Top)
	}
	if bounds.Left != 0 || bounds.Right < 20 {
		t.Errorf("Quadratic curve bounds: got %v", bounds)
	}

	// Test curve bounds (cubic)
	path = NewSkPath(enums.PathFillTypeDefault)
	path.MoveTo(0, 0)
	path.CubicTo(10, 20, 30, 20, 40, 0)
	bounds = path.Bounds()
	// Bounds should include control points
	if bounds.Top > 20 || bounds.Top < 0 {
		t.Errorf("Cubic curve bounds top: got %v, should include control points", bounds.Top)
	}
	if bounds.Left != 0 || bounds.Right < 40 {
		t.Errorf("Cubic curve bounds: got %v", bounds)
	}

	// Test multiple contour bounds
	path = NewSkPath(enums.PathFillTypeDefault)
	path.MoveTo(10, 10)
	path.LineTo(20, 10)
	path.MoveTo(30, 30)
	path.LineTo(40, 40)
	bounds = path.Bounds()
	expected = models.Rect{Left: 10, Top: 10, Right: 40, Bottom: 40}
	if bounds.Left != expected.Left || bounds.Top != expected.Top ||
		bounds.Right != expected.Right || bounds.Bottom != expected.Bottom {
		t.Errorf("Multiple contour bounds: got %v, want %v", bounds, expected)
	}
}

// TestPath_BoundsInvalidation tests bounds invalidation after edits
func TestPath_BoundsInvalidation(t *testing.T) {
	path := NewSkPath(enums.PathFillTypeDefault)
	path.MoveTo(10, 20)
	path.LineTo(30, 40)

	// Get initial bounds
	bounds1 := path.Bounds()
	expected1 := models.Rect{Left: 10, Top: 20, Right: 30, Bottom: 40}
	if bounds1 != expected1 {
		t.Errorf("Initial bounds: got %v, want %v", bounds1, expected1)
	}

	// Add a point that extends bounds
	path.LineTo(50, 60)
	bounds2 := path.Bounds()
	expected2 := models.Rect{Left: 10, Top: 20, Right: 50, Bottom: 60}
	if bounds2 != expected2 {
		t.Errorf("Bounds after edit: got %v, want %v", bounds2, expected2)
	}

	// Add a new contour
	path.MoveTo(100, 100)
	path.LineTo(110, 110)
	bounds3 := path.Bounds()
	expected3 := models.Rect{Left: 10, Top: 20, Right: 110, Bottom: 110}
	if bounds3.Left != expected3.Left || bounds3.Top != expected3.Top ||
		bounds3.Right != expected3.Right || bounds3.Bottom != expected3.Bottom {
		t.Errorf("Bounds after new contour: got %v, want %v", bounds3, expected3)
	}
}

// TestPath_UpdateBoundsCache tests UpdateBoundsCache() behavior
func TestPath_UpdateBoundsCache(t *testing.T) {
	path := NewSkPath(enums.PathFillTypeDefault)
	path.MoveTo(10, 20)
	path.LineTo(30, 40)

	// Update bounds cache explicitly
	path.UpdateBoundsCache()

	// Bounds should be computed
	bounds := path.Bounds()
	expected := models.Rect{Left: 10, Top: 20, Right: 30, Bottom: 40}
	if bounds != expected {
		t.Errorf("Bounds after UpdateBoundsCache: got %v, want %v", bounds, expected)
	}

	// Calling UpdateBoundsCache multiple times should be safe
	path.UpdateBoundsCache()
	bounds2 := path.Bounds()
	if bounds2 != bounds {
		t.Error("Multiple UpdateBoundsCache calls should produce same bounds")
	}
}

// TestPath_BoundsWithRects tests bounds calculation with rect addition
// Based on: skia-source/tests/PathTest.cpp:test_bounds()
func TestPath_BoundsWithRects(t *testing.T) {
	// Test single rect bounds
	rect := models.Rect{Left: 10, Top: 20, Right: 30, Bottom: 40}
	path := NewSkPath(enums.PathFillTypeDefault)
	path.AddRect(rect, enums.PathDirectionCW, 0)

	bounds := path.Bounds()
	if bounds.Left != rect.Left || bounds.Top != rect.Top ||
		bounds.Right != rect.Right || bounds.Bottom != rect.Bottom {
		t.Errorf("Single rect bounds: got %v, want %v", bounds, rect)
	}
}

