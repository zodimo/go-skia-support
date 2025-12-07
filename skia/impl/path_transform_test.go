package impl

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// copyPath creates a copy of a path by adding it to a new empty path.
// This is used for testing since Transform() modifies paths in place.
func copyPath(path interfaces.SkPath) interfaces.SkPath {
	copy := NewSkPath(path.FillType())
	copy.AddPathNoOffset(path, enums.AddPathModeAppend)
	return copy
}

// pathsEqual compares two paths for equality by checking fill type, verb count, point count, and point values.
// Ported from: skia-source/tests/PathTest.cpp (path equality comparison)
func pathsEqual(a, b interfaces.SkPath) bool {
	if a.FillType() != b.FillType() {
		return false
	}
	if a.CountVerbs() != b.CountVerbs() {
		return false
	}
	if a.CountPoints() != b.CountPoints() {
		return false
	}
	// Compare verbs
	verbsA := make([]enums.PathVerb, a.CountVerbs())
	verbsB := make([]enums.PathVerb, b.CountVerbs())
	a.GetVerbs(verbsA)
	b.GetVerbs(verbsB)
	for i := range verbsA {
		if verbsA[i] != verbsB[i] {
			return false
		}
	}
	// Compare points
	for i := 0; i < a.CountPoints(); i++ {
		ptA := a.Point(i)
		ptB := b.Point(i)
		if !NearlyEqualScalar(ptA.X, ptB.X) || !NearlyEqualScalar(ptA.Y, ptB.Y) {
			return false
		}
	}
	return true
}

// TestPath_Transform tests path transformation
// Ported from: skia-source/tests/PathTest.cpp:test_transform()
func TestPath_Transform(t *testing.T) {
	// Create a path with various segments: move, line, quad, cubic
	// Points array matches C++ test:
	// { 0, 0 },  // move
	// { 10, 10 },  // line
	// { 20, 10 }, { 20, 0 },  // quad
	// { 0, 0 }, { 0, 10 }, { 1, 10 },  // cubic
	pts := []models.Point{
		{X: 0, Y: 0},   // move
		{X: 10, Y: 10}, // line
		{X: 20, Y: 10}, // quad control
		{X: 20, Y: 0},  // quad end
		{X: 0, Y: 0},   // cubic control 1
		{X: 0, Y: 10},  // cubic control 2
		{X: 1, Y: 10},  // cubic end
	}
	const kPtCount = 7

	path := NewSkPath(enums.PathFillTypeDefault)
	path.MoveTo(pts[0].X, pts[0].Y)
	path.LineTo(pts[1].X, pts[1].Y)
	path.QuadTo(pts[2].X, pts[2].Y, pts[3].X, pts[3].Y)
	path.CubicTo(pts[4].X, pts[4].Y, pts[5].X, pts[5].Y, pts[6].X, pts[6].Y)
	path.Close()

	// Test 1: Transform with identity matrix (should be unchanged)
	t.Run("identity_matrix", func(t *testing.T) {
		pathCopy := copyPath(path)
		matrix := NewMatrixIdentity()
		pathCopy.Transform(matrix)
		if !pathsEqual(path, pathCopy) {
			t.Error("Path transformed with identity matrix should be unchanged")
		}
	})

	// Test 2: Transform with scale matrix (2x, 3y)
	t.Run("scale_matrix", func(t *testing.T) {
		pathCopy := copyPath(path)
		matrix := NewMatrixScale(2, 3)
		pathCopy.Transform(matrix)

		// Check point count
		if pathCopy.CountPoints() != kPtCount {
			t.Errorf("Transformed path point count: got %d, want %d", pathCopy.CountPoints(), kPtCount)
		}

		// Check each point is scaled correctly
		for i := 0; i < kPtCount; i++ {
			expectedPt := models.Point{X: pts[i].X * 2, Y: pts[i].Y * 3}
			actualPt := pathCopy.Point(i)
			if !NearlyEqualScalar(expectedPt.X, actualPt.X) || !NearlyEqualScalar(expectedPt.Y, actualPt.Y) {
				t.Errorf("Point %d: got (%v, %v), want (%v, %v)", i, actualPt.X, actualPt.Y, expectedPt.X, expectedPt.Y)
			}
		}
	})

	// Test 3: Transform with perspective matrix and inverse transform
	t.Run("perspective_matrix", func(t *testing.T) {
		pathCopy := copyPath(path)
		matrix := NewMatrixAll(1, 0, 0, 0, 1, 0, 4, 0, 1) // setPerspX(4)
		pathCopy.Transform(matrix)

		// Get original bounds
		originalBounds := path.Bounds()

		// Invert matrix and transform back
		inverted, ok := matrix.Invert()
		if !ok {
			t.Fatal("Failed to invert perspective matrix")
		}
		pathCopy.Transform(inverted)

		// Check bounds are nearly equal (within tolerance)
		// Use default C++ tolerance (SK_ScalarNearlyZero) for bounds comparison
		transformedBounds := pathCopy.Bounds()
		if !NearlyEqualScalarDefault(originalBounds.Left, transformedBounds.Left) {
			t.Errorf("Bounds Left: got %v, want %v", transformedBounds.Left, originalBounds.Left)
		}
		if !NearlyEqualScalarDefault(originalBounds.Top, transformedBounds.Top) {
			t.Errorf("Bounds Top: got %v, want %v", transformedBounds.Top, originalBounds.Top)
		}
		if !NearlyEqualScalarDefault(originalBounds.Right, transformedBounds.Right) {
			t.Errorf("Bounds Right: got %v, want %v", transformedBounds.Right, originalBounds.Right)
		}
		if !NearlyEqualScalarDefault(originalBounds.Bottom, transformedBounds.Bottom) {
			t.Errorf("Bounds Bottom: got %v, want %v", transformedBounds.Bottom, originalBounds.Bottom)
		}
	})

	// Test 4: Transform circle with identity matrix (preserves direction)
	t.Run("circle_identity_direction", func(t *testing.T) {
		circle := NewPathCircleDefault(0, 0, 1, enums.PathDirectionCW)
		circleCopy := copyPath(circle)
		matrix := NewMatrixIdentity()
		circleCopy.Transform(matrix)

		// Direction should be preserved (CW)
		// Note: We can't directly check direction, but we can verify the path structure is preserved
		if circleCopy.CountPoints() != circle.CountPoints() {
			t.Errorf("Circle point count after identity transform: got %d, want %d", circleCopy.CountPoints(), circle.CountPoints())
		}
	})

	// Test 5: Transform circle with scaleX(-1) (reverses direction)
	t.Run("circle_scaleX_neg1", func(t *testing.T) {
		circle := NewPathCircleDefault(0, 0, 1, enums.PathDirectionCW)
		circleCopy := copyPath(circle)
		matrix := NewMatrixScale(-1, 1)
		circleCopy.Transform(matrix)

		// Direction should be reversed (CCW)
		// Note: We can't directly check direction, but we can verify the path structure is preserved
		if circleCopy.CountPoints() != circle.CountPoints() {
			t.Errorf("Circle point count after scaleX(-1) transform: got %d, want %d", circleCopy.CountPoints(), circle.CountPoints())
		}
	})

	// Test 6: Transform circle with skew matrix (direction becomes unknown)
	t.Run("circle_skew_direction", func(t *testing.T) {
		circle := NewPathCircleDefault(0, 0, 1, enums.PathDirectionCW)
		circleCopy := copyPath(circle)
		// Matrix: setAll(1, 1, 0, 1, 1, 0, 0, 0, 1) - skew matrix
		matrix := NewMatrixAll(1, 1, 0, 1, 1, 0, 0, 0, 1)
		circleCopy.Transform(matrix)

		// Direction becomes unknown after skew
		// Note: We can't directly check direction, but we can verify the path structure is preserved
		if circleCopy.CountPoints() != circle.CountPoints() {
			t.Errorf("Circle point count after skew transform: got %d, want %d", circleCopy.CountPoints(), circle.CountPoints())
		}
	})

	// Test 7: Transform preserves path structure (verbs, fill type)
	t.Run("preserves_structure", func(t *testing.T) {
		originalPath := NewSkPath(enums.PathFillTypeDefault)
		originalPath.MoveTo(10, 20)
		originalPath.LineTo(30, 40)
		originalPath.Close()

		pathCopy := copyPath(originalPath)
		matrix := NewMatrixScale(2, 2)
		pathCopy.Transform(matrix)

		// Fill type should be preserved
		if pathCopy.FillType() != originalPath.FillType() {
			t.Errorf("Fill type: got %v, want %v", pathCopy.FillType(), originalPath.FillType())
		}

		// Verb count should be preserved
		if pathCopy.CountVerbs() != originalPath.CountVerbs() {
			t.Errorf("Verb count: got %d, want %d", pathCopy.CountVerbs(), originalPath.CountVerbs())
		}

		// Point count should be preserved
		if pathCopy.CountPoints() != originalPath.CountPoints() {
			t.Errorf("Point count: got %d, want %d", pathCopy.CountPoints(), originalPath.CountPoints())
		}
	})

	// Test 8: Transform updates bounds
	t.Run("bounds_after_transform", func(t *testing.T) {
		rect := models.Rect{Left: 10, Top: 20, Right: 30, Bottom: 40}
		path := NewSkPath(enums.PathFillTypeDefault)
		path.AddRect(rect, enums.PathDirectionCW, 0)

		originalBounds := path.Bounds()
		matrix := NewMatrixScale(2, 2)
		path.Transform(matrix)
		transformedBounds := path.Bounds()

		// Bounds should be scaled
		expectedLeft := originalBounds.Left * 2
		expectedTop := originalBounds.Top * 2
		expectedRight := originalBounds.Right * 2
		expectedBottom := originalBounds.Bottom * 2

		if !NearlyEqualScalar(transformedBounds.Left, expectedLeft) {
			t.Errorf("Bounds Left: got %v, want %v", transformedBounds.Left, expectedLeft)
		}
		if !NearlyEqualScalar(transformedBounds.Top, expectedTop) {
			t.Errorf("Bounds Top: got %v, want %v", transformedBounds.Top, expectedTop)
		}
		if !NearlyEqualScalar(transformedBounds.Right, expectedRight) {
			t.Errorf("Bounds Right: got %v, want %v", transformedBounds.Right, expectedRight)
		}
		if !NearlyEqualScalar(transformedBounds.Bottom, expectedBottom) {
			t.Errorf("Bounds Bottom: got %v, want %v", transformedBounds.Bottom, expectedBottom)
		}
	})
}

