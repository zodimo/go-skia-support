package impl

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// TestPath_AddRect tests rectangle addition
// Ported from: skia-source/tests/PathTest.cpp:test_addrect()
func TestPath_AddRect(t *testing.T) {
	// Test adding rect with different directions
	t.Run("add_rect_directions", func(t *testing.T) {
		rect := models.Rect{Left: 10, Top: 20, Right: 50, Bottom: 60}

		// Test CW direction
		pathCW := NewSkPath(enums.PathFillTypeDefault)
		pathCW.AddRect(rect, enums.PathDirectionCW, 0)
		boundsCW := pathCW.Bounds()
		if !NearlyEqualScalarDefault(boundsCW.Left, rect.Left) ||
			!NearlyEqualScalarDefault(boundsCW.Top, rect.Top) ||
			!NearlyEqualScalarDefault(boundsCW.Right, rect.Right) ||
			!NearlyEqualScalarDefault(boundsCW.Bottom, rect.Bottom) {
			t.Errorf("CW rect bounds: got %v, want %v", boundsCW, rect)
		}

		// Test CCW direction
		pathCCW := NewSkPath(enums.PathFillTypeDefault)
		pathCCW.AddRect(rect, enums.PathDirectionCCW, 0)
		boundsCCW := pathCCW.Bounds()
		if !NearlyEqualScalarDefault(boundsCCW.Left, rect.Left) ||
			!NearlyEqualScalarDefault(boundsCCW.Top, rect.Top) ||
			!NearlyEqualScalarDefault(boundsCCW.Right, rect.Right) ||
			!NearlyEqualScalarDefault(boundsCCW.Bottom, rect.Bottom) {
			t.Errorf("CCW rect bounds: got %v, want %v", boundsCCW, rect)
		}
	})

	// Test adding rect with different start indices
	t.Run("add_rect_start_indices", func(t *testing.T) {
		rect := models.Rect{Left: 0, Top: 0, Right: 50, Bottom: 100}

		for startIndex := uint(0); startIndex < 4; startIndex++ {
			path := NewSkPath(enums.PathFillTypeDefault)
			path.AddRect(rect, enums.PathDirectionCW, startIndex)

			// All should have same bounds
			bounds := path.Bounds()
			if !NearlyEqualScalarDefault(bounds.Left, rect.Left) ||
				!NearlyEqualScalarDefault(bounds.Top, rect.Top) ||
				!NearlyEqualScalarDefault(bounds.Right, rect.Right) ||
				!NearlyEqualScalarDefault(bounds.Bottom, rect.Bottom) {
				t.Errorf("Start index %d bounds: got %v, want %v", startIndex, bounds, rect)
			}

			// Should have 4 points (4 corners, Close doesn't add a point)
			if path.CountPoints() != 4 {
				t.Errorf("Start index %d point count: got %d, want 4", startIndex, path.CountPoints())
			}
		}
	})

	// Test verb sequence for rect
	t.Run("add_rect_verb_sequence", func(t *testing.T) {
		rect := models.Rect{Left: 10, Top: 20, Right: 30, Bottom: 40}
		path := NewSkPath(enums.PathFillTypeDefault)
		path.AddRect(rect, enums.PathDirectionCW, 0)

		verbs := make([]enums.PathVerb, path.CountVerbs())
		path.GetVerbs(verbs)

		// Expected: Move, Line, Line, Line, Close
		expectedVerbs := []enums.PathVerb{
			enums.PathVerbMove,
			enums.PathVerbLine,
			enums.PathVerbLine,
			enums.PathVerbLine,
			enums.PathVerbClose,
		}

		if len(verbs) != len(expectedVerbs) {
			t.Errorf("Verb count: got %d, want %d", len(verbs), len(expectedVerbs))
			return
		}

		for i, expectedVerb := range expectedVerbs {
			if verbs[i] != expectedVerb {
				t.Errorf("Verb %d: got %v, want %v", i, verbs[i], expectedVerb)
			}
		}
	})

	// Test adding rect after lineTo (from C++ test)
	t.Run("add_rect_after_lineTo", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		path.LineTo(0, 0)
		path.AddRect(models.Rect{Left: 0, Top: 0, Right: 50, Bottom: 100}, enums.PathDirectionCW, 0)

		// Path should have bounds matching the rect
		bounds := path.Bounds()
		expectedBounds := models.Rect{Left: 0, Top: 0, Right: 50, Bottom: 100}
		if !NearlyEqualScalarDefault(bounds.Left, expectedBounds.Left) ||
			!NearlyEqualScalarDefault(bounds.Top, expectedBounds.Top) ||
			!NearlyEqualScalarDefault(bounds.Right, expectedBounds.Right) ||
			!NearlyEqualScalarDefault(bounds.Bottom, expectedBounds.Bottom) {
			t.Errorf("Bounds after lineTo + addRect: got %v, want %v", bounds, expectedBounds)
		}
	})

	// Test adding multiple rects
	t.Run("add_multiple_rects", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		rect1 := models.Rect{Left: 10, Top: 20, Right: 30, Bottom: 40}
		rect2 := models.Rect{Left: 50, Top: 60, Right: 70, Bottom: 80}

		path.AddRect(rect1, enums.PathDirectionCW, 0)
		path.AddRect(rect2, enums.PathDirectionCW, 0)

		// Bounds should encompass both rects
		bounds := path.Bounds()
		expectedBounds := models.Rect{Left: 10, Top: 20, Right: 70, Bottom: 80}
		if !NearlyEqualScalarDefault(bounds.Left, expectedBounds.Left) ||
			!NearlyEqualScalarDefault(bounds.Top, expectedBounds.Top) ||
			!NearlyEqualScalarDefault(bounds.Right, expectedBounds.Right) ||
			!NearlyEqualScalarDefault(bounds.Bottom, expectedBounds.Bottom) {
			t.Errorf("Bounds for multiple rects: got %v, want %v", bounds, expectedBounds)
		}
	})
}

// TestPath_AddCircle tests circle addition
// Ported from: skia-source/tests/PathTest.cpp:test_circle()
func TestPath_AddCircle(t *testing.T) {
	// Test circle with different radii
	t.Run("add_circle_radii", func(t *testing.T) {
		radii := []interfaces.Scalar{1, 5, 10, 20, 50, 100}
		for _, radius := range radii {
			path := NewSkPath(enums.PathFillTypeDefault)
			path.AddCircle(0, 0, radius, enums.PathDirectionCW)

			bounds := path.Bounds()
			expectedLeft := -radius
			expectedTop := -radius
			expectedRight := radius
			expectedBottom := radius

			if !NearlyEqualScalarDefault(bounds.Left, expectedLeft) ||
				!NearlyEqualScalarDefault(bounds.Top, expectedTop) ||
				!NearlyEqualScalarDefault(bounds.Right, expectedRight) ||
				!NearlyEqualScalarDefault(bounds.Bottom, expectedBottom) {
				t.Errorf("Circle radius %v bounds: got %v, want (%v, %v, %v, %v)",
					radius, bounds, expectedLeft, expectedTop, expectedRight, expectedBottom)
			}
		}
	})

	// Test circle with different directions
	t.Run("add_circle_directions", func(t *testing.T) {
		radius := interfaces.Scalar(10)
		cx, cy := interfaces.Scalar(0), interfaces.Scalar(0)

		// Test CW direction
		pathCW := NewSkPath(enums.PathFillTypeDefault)
		pathCW.AddCircle(cx, cy, radius, enums.PathDirectionCW)
		boundsCW := pathCW.Bounds()

		// Test CCW direction
		pathCCW := NewSkPath(enums.PathFillTypeDefault)
		pathCCW.AddCircle(cx, cy, radius, enums.PathDirectionCCW)
		boundsCCW := pathCCW.Bounds()

		// Both should have same bounds
		if !NearlyEqualScalarDefault(boundsCW.Left, boundsCCW.Left) ||
			!NearlyEqualScalarDefault(boundsCW.Top, boundsCCW.Top) ||
			!NearlyEqualScalarDefault(boundsCW.Right, boundsCCW.Right) ||
			!NearlyEqualScalarDefault(boundsCW.Bottom, boundsCCW.Bottom) {
			t.Errorf("CW and CCW circle bounds differ: CW=%v, CCW=%v", boundsCW, boundsCCW)
		}
	})

	// Test circle at different positions
	t.Run("add_circle_positions", func(t *testing.T) {
		testCases := []struct {
			cx, cy, radius interfaces.Scalar
		}{
			{0, 0, 10},
			{10, 10, 20},
			{100, 200, 50},
			{-50, -50, 25},
		}

		for _, tc := range testCases {
			path := NewSkPath(enums.PathFillTypeDefault)
			path.AddCircle(tc.cx, tc.cy, tc.radius, enums.PathDirectionCW)

			bounds := path.Bounds()
			expectedLeft := tc.cx - tc.radius
			expectedTop := tc.cy - tc.radius
			expectedRight := tc.cx + tc.radius
			expectedBottom := tc.cy + tc.radius

			if !NearlyEqualScalarDefault(bounds.Left, expectedLeft) ||
				!NearlyEqualScalarDefault(bounds.Top, expectedTop) ||
				!NearlyEqualScalarDefault(bounds.Right, expectedRight) ||
				!NearlyEqualScalarDefault(bounds.Bottom, expectedBottom) {
				t.Errorf("Circle at (%v, %v) radius %v bounds: got %v, want (%v, %v, %v, %v)",
					tc.cx, tc.cy, tc.radius, bounds, expectedLeft, expectedTop, expectedRight, expectedBottom)
			}
		}
	})

	// Test multiple circles
	t.Run("add_multiple_circles", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		path.AddCircle(0, 0, 10, enums.PathDirectionCW)
		path.AddCircle(0, 0, 20, enums.PathDirectionCW)

		// Bounds should encompass both circles
		bounds := path.Bounds()
		expectedBounds := models.Rect{Left: -20, Top: -20, Right: 20, Bottom: 20}
		if !NearlyEqualScalarDefault(bounds.Left, expectedBounds.Left) ||
			!NearlyEqualScalarDefault(bounds.Top, expectedBounds.Top) ||
			!NearlyEqualScalarDefault(bounds.Right, expectedBounds.Right) ||
			!NearlyEqualScalarDefault(bounds.Bottom, expectedBounds.Bottom) {
			t.Errorf("Multiple circles bounds: got %v, want %v", bounds, expectedBounds)
		}
	})

	// Test negative radius (should result in empty path)
	t.Run("add_circle_negative_radius", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		path.AddCircle(0, 0, -1, enums.PathDirectionCW)

		if !path.IsEmpty() {
			t.Error("Expected empty path for negative radius")
		}
	})
}

// TestPath_AddOval tests oval addition
// Ported from: skia-source/tests/PathTest.cpp:test_oval()
func TestPath_AddOval(t *testing.T) {
	// Test oval with different rects
	t.Run("add_oval_rects", func(t *testing.T) {
		testCases := []models.Rect{
			{Left: 0, Top: 0, Right: 30, Bottom: 50},
			{Left: 10, Top: 20, Right: 40, Bottom: 70},
			{Left: -50, Top: -50, Right: 50, Bottom: 50},
		}

		for _, rect := range testCases {
			path := NewSkPath(enums.PathFillTypeDefault)
			path.AddOval(rect, enums.PathDirectionCW)

			bounds := path.Bounds()
			if !NearlyEqualScalarDefault(bounds.Left, rect.Left) ||
				!NearlyEqualScalarDefault(bounds.Top, rect.Top) ||
				!NearlyEqualScalarDefault(bounds.Right, rect.Right) ||
				!NearlyEqualScalarDefault(bounds.Bottom, rect.Bottom) {
				t.Errorf("Oval rect %v bounds: got %v, want %v", rect, bounds, rect)
			}
		}
	})

	// Test oval with different directions
	t.Run("add_oval_directions", func(t *testing.T) {
		rect := models.Rect{Left: 10, Top: 20, Right: 40, Bottom: 60}

		// Test CW direction
		pathCW := NewSkPath(enums.PathFillTypeDefault)
		pathCW.AddOval(rect, enums.PathDirectionCW)
		boundsCW := pathCW.Bounds()

		// Test CCW direction
		pathCCW := NewSkPath(enums.PathFillTypeDefault)
		pathCCW.AddOval(rect, enums.PathDirectionCCW)
		boundsCCW := pathCCW.Bounds()

		// Both should have same bounds
		if !NearlyEqualScalarDefault(boundsCW.Left, boundsCCW.Left) ||
			!NearlyEqualScalarDefault(boundsCW.Top, boundsCCW.Top) ||
			!NearlyEqualScalarDefault(boundsCW.Right, boundsCCW.Right) ||
			!NearlyEqualScalarDefault(boundsCW.Bottom, boundsCCW.Bottom) {
			t.Errorf("CW and CCW oval bounds differ: CW=%v, CCW=%v", boundsCW, boundsCCW)
		}
	})

	// Test oval verb sequence (should have conic verbs)
	t.Run("add_oval_verb_sequence", func(t *testing.T) {
		rect := models.Rect{Left: 0, Top: 0, Right: 30, Bottom: 50}
		path := NewSkPath(enums.PathFillTypeDefault)
		path.AddOval(rect, enums.PathDirectionCW)

		verbs := make([]enums.PathVerb, path.CountVerbs())
		path.GetVerbs(verbs)

		// Should have Move, Conic, Conic, Conic, Conic, Close
		if len(verbs) < 6 {
			t.Errorf("Oval verb count: got %d, want at least 6", len(verbs))
			return
		}

		if verbs[0] != enums.PathVerbMove {
			t.Errorf("First verb: got %v, want Move", verbs[0])
		}
		if verbs[len(verbs)-1] != enums.PathVerbClose {
			t.Errorf("Last verb: got %v, want Close", verbs[len(verbs)-1])
		}

		// Check for conic verbs
		hasConic := false
		for _, verb := range verbs {
			if verb == enums.PathVerbConic {
				hasConic = true
				break
			}
		}
		if !hasConic {
			t.Error("Oval should contain Conic verbs")
		}
	})

	// Test multiple ovals
	t.Run("add_multiple_ovals", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		rect1 := models.Rect{Left: 0, Top: 0, Right: 30, Bottom: 50}
		rect2 := models.Rect{Left: 50, Top: 60, Right: 80, Bottom: 110}

		path.AddOval(rect1, enums.PathDirectionCW)
		path.AddOval(rect2, enums.PathDirectionCW)

		// Bounds should encompass both ovals
		bounds := path.Bounds()
		expectedBounds := models.Rect{Left: 0, Top: 0, Right: 80, Bottom: 110}
		if !NearlyEqualScalarDefault(bounds.Left, expectedBounds.Left) ||
			!NearlyEqualScalarDefault(bounds.Top, expectedBounds.Top) ||
			!NearlyEqualScalarDefault(bounds.Right, expectedBounds.Right) ||
			!NearlyEqualScalarDefault(bounds.Bottom, expectedBounds.Bottom) {
			t.Errorf("Multiple ovals bounds: got %v, want %v", bounds, expectedBounds)
		}
	})
}

// TestPath_AddRRect tests rounded rectangle addition
// Ported from: skia-source/tests/PathTest.cpp:test_rrect()
func TestPath_AddRRect(t *testing.T) {
	// Helper to create RRect
	createRRect := func(rect models.Rect, radii [4]models.Point) models.RRect {
		var rrect models.RRect
		rrect.SetRectRadii(rect, radii)
		return rrect
	}

	// Test RRect with various corner radii
	t.Run("add_rrect_corner_radii", func(t *testing.T) {
		rect := models.Rect{Left: 10, Top: 20, Right: 30, Bottom: 40}
		radii := [4]models.Point{
			{X: 1, Y: 2}, // UpperLeft
			{X: 3, Y: 4}, // UpperRight
			{X: 5, Y: 6}, // LowerRight
			{X: 7, Y: 8}, // LowerLeft
		}

		rrect := createRRect(rect, radii)
		path := NewSkPath(enums.PathFillTypeDefault)
		path.AddRRect(rrect, enums.PathDirectionCW)

		bounds := path.Bounds()
		if !NearlyEqualScalarDefault(bounds.Left, rect.Left) ||
			!NearlyEqualScalarDefault(bounds.Top, rect.Top) ||
			!NearlyEqualScalarDefault(bounds.Right, rect.Right) ||
			!NearlyEqualScalarDefault(bounds.Bottom, rect.Bottom) {
			t.Errorf("RRect bounds: got %v, want %v", bounds, rect)
		}
	})

	// Test RRect with different directions
	t.Run("add_rrect_directions", func(t *testing.T) {
		rect := models.Rect{Left: 10, Top: 20, Right: 30, Bottom: 40}
		radii := [4]models.Point{
			{X: 2, Y: 2},
			{X: 2, Y: 2},
			{X: 2, Y: 2},
			{X: 2, Y: 2},
		}

		rrect := createRRect(rect, radii)

		// Test CW direction
		pathCW := NewSkPath(enums.PathFillTypeDefault)
		pathCW.AddRRect(rrect, enums.PathDirectionCW)
		boundsCW := pathCW.Bounds()

		// Test CCW direction
		pathCCW := NewSkPath(enums.PathFillTypeDefault)
		pathCCW.AddRRect(rrect, enums.PathDirectionCCW)
		boundsCCW := pathCCW.Bounds()

		// Both should have same bounds
		if !NearlyEqualScalarDefault(boundsCW.Left, boundsCCW.Left) ||
			!NearlyEqualScalarDefault(boundsCW.Top, boundsCCW.Top) ||
			!NearlyEqualScalarDefault(boundsCW.Right, boundsCCW.Right) ||
			!NearlyEqualScalarDefault(boundsCW.Bottom, boundsCCW.Bottom) {
			t.Errorf("CW and CCW RRect bounds differ: CW=%v, CCW=%v", boundsCW, boundsCCW)
		}
	})

	// Test RRect with zero radii (should be equivalent to rect)
	t.Run("add_rrect_zero_radii", func(t *testing.T) {
		rect := models.Rect{Left: 10, Top: 20, Right: 30, Bottom: 40}
		zeroRadii := [4]models.Point{
			{X: 0, Y: 0},
			{X: 0, Y: 0},
			{X: 0, Y: 0},
			{X: 0, Y: 0},
		}

		rrect := createRRect(rect, zeroRadii)
		path := NewSkPath(enums.PathFillTypeDefault)
		path.AddRRect(rrect, enums.PathDirectionCW)

		bounds := path.Bounds()
		if !NearlyEqualScalarDefault(bounds.Left, rect.Left) ||
			!NearlyEqualScalarDefault(bounds.Top, rect.Top) ||
			!NearlyEqualScalarDefault(bounds.Right, rect.Right) ||
			!NearlyEqualScalarDefault(bounds.Bottom, rect.Bottom) {
			t.Errorf("RRect with zero radii bounds: got %v, want %v", bounds, rect)
		}
	})

	// Test RRect with some zero radii
	t.Run("add_rrect_partial_zero_radii", func(t *testing.T) {
		rect := models.Rect{Left: 10, Top: 20, Right: 30, Bottom: 40}
		radii := [4]models.Point{
			{X: 0, Y: 0}, // UpperLeft - zero
			{X: 2, Y: 2}, // UpperRight
			{X: 0, Y: 0}, // LowerRight - zero
			{X: 2, Y: 2}, // LowerLeft
		}

		rrect := createRRect(rect, radii)
		path := NewSkPath(enums.PathFillTypeDefault)
		path.AddRRect(rrect, enums.PathDirectionCW)

		bounds := path.Bounds()
		if !NearlyEqualScalarDefault(bounds.Left, rect.Left) ||
			!NearlyEqualScalarDefault(bounds.Top, rect.Top) ||
			!NearlyEqualScalarDefault(bounds.Right, rect.Right) ||
			!NearlyEqualScalarDefault(bounds.Bottom, rect.Bottom) {
			t.Errorf("RRect with partial zero radii bounds: got %v, want %v", bounds, rect)
		}
	})

	// Test multiple RRects
	t.Run("add_multiple_rrects", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		rect1 := models.Rect{Left: 10, Top: 20, Right: 30, Bottom: 40}
		rect2 := models.Rect{Left: 50, Top: 60, Right: 70, Bottom: 80}
		radii := [4]models.Point{
			{X: 2, Y: 2},
			{X: 2, Y: 2},
			{X: 2, Y: 2},
			{X: 2, Y: 2},
		}

		rrect1 := createRRect(rect1, radii)
		rrect2 := createRRect(rect2, radii)

		path.AddRRect(rrect1, enums.PathDirectionCW)
		path.AddRRect(rrect2, enums.PathDirectionCW)

		// Bounds should encompass both RRects
		bounds := path.Bounds()
		expectedBounds := models.Rect{Left: 10, Top: 20, Right: 70, Bottom: 80}
		if !NearlyEqualScalarDefault(bounds.Left, expectedBounds.Left) ||
			!NearlyEqualScalarDefault(bounds.Top, expectedBounds.Top) ||
			!NearlyEqualScalarDefault(bounds.Right, expectedBounds.Right) ||
			!NearlyEqualScalarDefault(bounds.Bottom, expectedBounds.Bottom) {
			t.Errorf("Multiple RRects bounds: got %v, want %v", bounds, expectedBounds)
		}
	})
}

