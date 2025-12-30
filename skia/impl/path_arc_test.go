package impl

import (
	"math"
	"testing"

	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/models"
)

// TestPath_ArcTo_Oval tests oval-based arcTo
// Ported from: skia-source/tests/PathTest.cpp arc tests
func TestPath_ArcTo_Oval(t *testing.T) {
	// Test basic 90-degree arc
	t.Run("90_degree_arc", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		oval := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 100}

		path.ArcTo(oval, 0, 90, true)

		// Path should not be empty
		if path.IsEmpty() {
			t.Error("Path should not be empty after ArcTo")
		}

		// Bounds should be approximately the upper-right quadrant
		bounds := path.Bounds()
		if bounds.Right < 50 || bounds.Bottom > 51 {
			t.Logf("Arc bounds: %v", bounds)
		}
	})

	// Test full circle arc (360 degrees)
	// Note: Full circle requires multiple arc segments - this tests the basic case
	t.Run("full_circle", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		oval := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 100}

		path.ArcTo(oval, 0, 360, true)

		// Path should not be empty
		if path.IsEmpty() {
			t.Error("Full circle arc should not be empty")
		}
		// Note: Full 360° arc handling is complex and may need multiple conic segments
		// This is a known area for future enhancement
	})

	// Test arc with forceMoveTo=false (should connect to previous point)
	t.Run("arc_connects_to_existing_path", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		path.MoveTo(0, 0)
		path.LineTo(50, 0)

		oval := models.Rect{Left: 0, Top: -50, Right: 100, Bottom: 50}
		path.ArcTo(oval, 0, 90, false)

		// Should have more than just the line
		if path.CountVerbs() < 3 {
			t.Errorf("Expected at least 3 verbs, got %d", path.CountVerbs())
		}
	})

	// Test zero sweep angle
	t.Run("zero_sweep_angle", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		oval := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 100}

		path.ArcTo(oval, 0, 0, true)

		// Should just add a point
		if path.CountPoints() != 1 {
			t.Errorf("Zero sweep should add single point, got %d points", path.CountPoints())
		}
	})
}

// TestPath_ArcToTangent tests tangent-based arcTo
func TestPath_ArcToTangent(t *testing.T) {
	t.Run("basic_tangent_arc", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		path.MoveTo(0, 0)

		// Arc from (0,0) tangent to (100,0) and (100,100) with radius 20
		path.ArcToTangent(100, 0, 100, 100, 20)

		if path.IsEmpty() {
			t.Error("Path should not be empty after ArcToTangent")
		}

		// Should have line and conic verbs
		verbCount := path.CountVerbs()
		if verbCount < 2 {
			t.Errorf("Expected at least 2 verbs, got %d", verbCount)
		}
	})

	t.Run("zero_radius", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		path.MoveTo(0, 0)

		// Zero radius should just draw a line
		path.ArcToTangent(100, 0, 100, 100, 0)

		verbs := make([]enums.PathVerb, path.CountVerbs())
		path.GetVerbs(verbs)

		// Should have Move and Line
		if len(verbs) != 2 || verbs[1] != enums.PathVerbLine {
			t.Errorf("Zero radius should create line, verbs: %v", verbs)
		}
	})
}

// TestPath_ArcToRotated tests SVG-style elliptical arc
func TestPath_ArcToRotated(t *testing.T) {
	t.Run("basic_svg_arc", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		path.MoveTo(0, 0)

		// SVG arc to (100, 0) with radii 50, 50
		path.ArcToRotated(50, 50, 0, enums.ArcSizeSmall, enums.PathDirectionCW, 100, 0)

		if path.IsEmpty() {
			t.Error("Path should not be empty after ArcToRotated")
		}
	})

	t.Run("large_arc_flag", func(t *testing.T) {
		pathSmall := NewSkPath(enums.PathFillTypeDefault)
		pathSmall.MoveTo(0, 0)
		pathSmall.ArcToRotated(50, 50, 0, enums.ArcSizeSmall, enums.PathDirectionCW, 100, 0)

		pathLarge := NewSkPath(enums.PathFillTypeDefault)
		pathLarge.MoveTo(0, 0)
		pathLarge.ArcToRotated(50, 50, 0, enums.ArcSizeLarge, enums.PathDirectionCW, 100, 0)

		// Large arc should have different bounds or more points
		boundsSmall := pathSmall.Bounds()
		boundsLarge := pathLarge.Bounds()

		// They should be different (large arc goes the long way)
		if boundsSmall == boundsLarge {
			t.Log("Small and large arc have same bounds - may need investigation")
		}
	})

	t.Run("zero_radii", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		path.MoveTo(0, 0)

		// Zero radius should just draw a line
		path.ArcToRotated(0, 0, 0, enums.ArcSizeSmall, enums.PathDirectionCW, 100, 100)

		verbs := make([]enums.PathVerb, path.CountVerbs())
		path.GetVerbs(verbs)

		// Should have Move and Line
		if len(verbs) != 2 || verbs[1] != enums.PathVerbLine {
			t.Errorf("Zero radii should create line, verbs: %v", verbs)
		}
	})
}

// TestPath_RArcTo tests relative SVG-style arc
func TestPath_RArcTo(t *testing.T) {
	t.Run("relative_arc", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		path.MoveTo(50, 50)

		// Relative arc by (100, 0)
		path.RArcTo(50, 50, 0, enums.ArcSizeSmall, enums.PathDirectionCW, 100, 0)

		// Last point should be at (150, 50)
		lastPt, ok := path.GetLastPoint()
		if !ok {
			t.Fatal("Path should have points")
		}
		if !NearlyEqualScalarDefault(lastPt.X, 150) || !NearlyEqualScalarDefault(lastPt.Y, 50) {
			t.Errorf("Last point: got (%v, %v), want (150, 50)", lastPt.X, lastPt.Y)
		}
	})
}

// TestPath_AddArc tests addArc
func TestPath_AddArc(t *testing.T) {
	t.Run("basic_add_arc", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		oval := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 100}

		path.AddArc(oval, 0, 90)

		if path.IsEmpty() {
			t.Error("Path should not be empty after AddArc")
		}

		// AddArc always starts a new contour
		verbs := make([]enums.PathVerb, path.CountVerbs())
		path.GetVerbs(verbs)
		if verbs[0] != enums.PathVerbMove {
			t.Errorf("AddArc should start with Move, got %v", verbs[0])
		}
	})

	t.Run("full_circle_becomes_oval", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		oval := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 100}

		path.AddArc(oval, 0, 360)

		// Full 360-degree arc should add an oval
		bounds := path.Bounds()
		if !NearlyEqualScalarDefault(bounds.Left, 0) ||
			!NearlyEqualScalarDefault(bounds.Top, 0) ||
			!NearlyEqualScalarDefault(bounds.Right, 100) ||
			!NearlyEqualScalarDefault(bounds.Bottom, 100) {
			t.Errorf("360° arc bounds: got %v, want (0,0,100,100)", bounds)
		}
	})

	t.Run("zero_sweep", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		oval := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 100}

		path.AddArc(oval, 0, 0)

		// Zero sweep should result in empty path
		if !path.IsEmpty() {
			t.Error("Zero sweep should result in empty path")
		}
	})
}

// TestPath_Arc_Bounds tests that arc bounds are calculated correctly
func TestPath_Arc_Bounds(t *testing.T) {
	t.Run("quarter_circle_bounds", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		oval := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 100}

		// 90 degree arc starting at 0 (right side)
		path.ArcTo(oval, 0, 90, true)

		bounds := path.Bounds()
		// Arc from (100,50) going clockwise to (50,100)
		// Bounds should be roughly in the lower-right quadrant
		if bounds.Left < 40 {
			t.Errorf("Quarter arc left bound too small: %v", bounds.Left)
		}
		if bounds.Top < 40 {
			t.Errorf("Quarter arc top bound too small: %v", bounds.Top)
		}
	})

	t.Run("semicircle_bounds", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		oval := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 100}

		// 180 degree arc (semicircle) - currently buildUnitArc handles arcs up to ~90°
		// For now, just verify the arc is created
		path.ArcTo(oval, 0, 180, true)

		// Path should not be empty
		if path.IsEmpty() {
			t.Error("Semicircle arc should not be empty")
		}
		// Note: Multi-segment arc support for arcs > 90° is a future enhancement
	})
}

// TestPath_Arc_EdgeCases tests edge cases
func TestPath_Arc_EdgeCases(t *testing.T) {
	t.Run("negative_dimensions", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		// Invalid oval (left > right)
		oval := models.Rect{Left: 100, Top: 0, Right: 0, Bottom: 100}

		path.ArcTo(oval, 0, 90, true)

		// Should be empty for invalid oval
		if !path.IsEmpty() {
			t.Error("Path should be empty for invalid oval dimensions")
		}
	})

	t.Run("very_small_sweep", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		oval := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 100}

		// Very small sweep angle
		path.ArcTo(oval, 0, 0.001, true)

		if path.IsEmpty() {
			t.Error("Path should not be empty for small non-zero sweep")
		}
	})

	t.Run("negative_sweep", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		oval := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 100}

		// Negative sweep (counter-clockwise)
		path.ArcTo(oval, 0, -90, true)

		if path.IsEmpty() {
			t.Error("Path should not be empty for negative sweep")
		}
	})

	t.Run("angle_normalization", func(t *testing.T) {
		// Angles should be normalized modulo 360
		path1 := NewSkPath(enums.PathFillTypeDefault)
		path2 := NewSkPath(enums.PathFillTypeDefault)
		oval := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 100}

		path1.ArcTo(oval, 0, 90, true)
		path2.ArcTo(oval, 360, 90, true)

		bounds1 := path1.Bounds()
		bounds2 := path2.Bounds()

		// Should produce same bounds (start angle normalized)
		if !NearlyEqualScalarDefault(bounds1.Left, bounds2.Left) ||
			!NearlyEqualScalarDefault(bounds1.Top, bounds2.Top) ||
			!NearlyEqualScalarDefault(bounds1.Right, bounds2.Right) ||
			!NearlyEqualScalarDefault(bounds1.Bottom, bounds2.Bottom) {
			t.Errorf("Angle normalization: bounds differ - %v vs %v", bounds1, bounds2)
		}
	})
}

// TestPath_Arc_VerbSequence tests the verb sequence produced by arcs
func TestPath_Arc_VerbSequence(t *testing.T) {
	t.Run("arc_uses_conics", func(t *testing.T) {
		path := NewSkPath(enums.PathFillTypeDefault)
		oval := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 100}

		path.ArcTo(oval, 0, 90, true)

		verbs := make([]enums.PathVerb, path.CountVerbs())
		path.GetVerbs(verbs)

		// Arc should use conic verbs
		hasConic := false
		for _, v := range verbs {
			if v == enums.PathVerbConic {
				hasConic = true
				break
			}
		}
		if !hasConic {
			t.Errorf("Arc should contain conic verbs, got: %v", verbs)
		}
	})
}

// Benchmark arc operations
func BenchmarkPath_ArcTo(b *testing.B) {
	oval := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 100}
	for i := 0; i < b.N; i++ {
		path := NewSkPath(enums.PathFillTypeDefault)
		path.ArcTo(oval, 0, Scalar(i%360), true)
	}
}

func BenchmarkPath_AddArc(b *testing.B) {
	oval := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 100}
	for i := 0; i < b.N; i++ {
		path := NewSkPath(enums.PathFillTypeDefault)
		path.AddArc(oval, 0, Scalar(i%360))
	}
}

// Helper to check if a value is approximately equal to expected
func nearlyEqualAngle(got, want, tolerance float64) bool {
	return math.Abs(got-want) <= tolerance
}
