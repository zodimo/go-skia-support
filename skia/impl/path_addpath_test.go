package impl

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/models"
)

// TestPath_AddPath tests path addition with offset
// Ported from: skia-source/tests/PathTest.cpp:test_addPath()
func TestPath_AddPath(t *testing.T) {
	// Create source path q with various segments
	q := NewSkPath(enums.PathFillTypeDefault)
	q.MoveTo(4, 4)
	q.LineTo(7, 8)
	q.ConicTo(8, 7, 6, 5, 0.5)
	q.QuadTo(6, 7, 8, 6)
	q.CubicTo(5, 6, 7, 8, 7, 5)
	q.Close()

	// Test adding path with offset
	t.Run("add_path_with_offset", func(t *testing.T) {
		p := NewSkPath(enums.PathFillTypeDefault)
		p.LineTo(1, 2)
		p.AddPath(q, -4, -4, enums.AddPathModeAppend)

		// Check bounds after adding path with offset
		bounds := p.Bounds()
		expected := models.Rect{Left: 0, Top: 0, Right: 4, Bottom: 4}
		if !NearlyEqualScalarDefault(bounds.Left, expected.Left) ||
			!NearlyEqualScalarDefault(bounds.Top, expected.Top) ||
			!NearlyEqualScalarDefault(bounds.Right, expected.Right) ||
			!NearlyEqualScalarDefault(bounds.Bottom, expected.Bottom) {
			t.Errorf("Bounds after AddPath with offset: got %v, want %v", bounds, expected)
		}
	})

	// Test adding path with matrix transformation
	t.Run("add_path_with_matrix", func(t *testing.T) {
		p := NewSkPath(enums.PathFillTypeDefault)
		matrix := NewMatrixScale(2, 2)
		p.AddPathMatrix(q, matrix, enums.AddPathModeAppend)

		// Check that points were transformed
		if p.CountPoints() != q.CountPoints() {
			t.Errorf("Point count: got %d, want %d", p.CountPoints(), q.CountPoints())
		}

		// Check first point is scaled (4, 4) -> (8, 8)
		firstPt := p.Point(0)
		expectedPt := models.Point{X: 8, Y: 8}
		if !NearlyEqualScalar(firstPt.X, expectedPt.X) || !NearlyEqualScalar(firstPt.Y, expectedPt.Y) {
			t.Errorf("First point after matrix transform: got (%v, %v), want (%v, %v)", firstPt.X, firstPt.Y, expectedPt.X, expectedPt.Y)
		}
	})

	// Test adding empty path
	t.Run("add_empty_path", func(t *testing.T) {
		emptyPath := NewSkPath(enums.PathFillTypeDefault)
		p := NewSkPath(enums.PathFillTypeDefault)
		p.MoveTo(1, 1)
		p.LineTo(2, 2)

		originalPointCount := p.CountPoints()
		p.AddPathNoOffset(emptyPath, enums.AddPathModeAppend)

		// Adding empty path should not change the path
		if p.CountPoints() != originalPointCount {
			t.Errorf("Point count after adding empty path: got %d, want %d", p.CountPoints(), originalPointCount)
		}
	})

	// Test adding path to empty path
	t.Run("add_to_empty_path", func(t *testing.T) {
		emptyPath := NewSkPath(enums.PathFillTypeDefault)
		emptyPath.AddPathNoOffset(q, enums.AddPathModeAppend)

		// Empty path should now have the same points as q
		if emptyPath.CountPoints() != q.CountPoints() {
			t.Errorf("Point count after adding to empty path: got %d, want %d", emptyPath.CountPoints(), q.CountPoints())
		}

		// Check that points match
		for i := 0; i < q.CountPoints(); i++ {
			ptQ := q.Point(i)
			ptEmpty := emptyPath.Point(i)
			if !NearlyEqualScalar(ptQ.X, ptEmpty.X) || !NearlyEqualScalar(ptQ.Y, ptEmpty.Y) {
				t.Errorf("Point %d mismatch: got (%v, %v), want (%v, %v)", i, ptEmpty.X, ptEmpty.Y, ptQ.X, ptQ.Y)
			}
		}
	})

	// Test verb and point copying
	t.Run("verb_and_point_copying", func(t *testing.T) {
		p := NewSkPath(enums.PathFillTypeDefault)
		p.AddPathNoOffset(q, enums.AddPathModeAppend)

		// Check verb count matches
		if p.CountVerbs() != q.CountVerbs() {
			t.Errorf("Verb count: got %d, want %d", p.CountVerbs(), q.CountVerbs())
		}

		// Check point count matches
		if p.CountPoints() != q.CountPoints() {
			t.Errorf("Point count: got %d, want %d", p.CountPoints(), q.CountPoints())
		}

		// Check verbs match
		verbsP := make([]enums.PathVerb, p.CountVerbs())
		verbsQ := make([]enums.PathVerb, q.CountVerbs())
		p.GetVerbs(verbsP)
		q.GetVerbs(verbsQ)
		for i := range verbsP {
			if verbsP[i] != verbsQ[i] {
				t.Errorf("Verb %d mismatch: got %v, want %v", i, verbsP[i], verbsQ[i])
			}
		}

		// Check points match (within tolerance)
		for i := 0; i < p.CountPoints(); i++ {
			ptP := p.Point(i)
			ptQ := q.Point(i)
			if !NearlyEqualScalar(ptP.X, ptQ.X) || !NearlyEqualScalar(ptP.Y, ptQ.Y) {
				t.Errorf("Point %d mismatch: got (%v, %v), want (%v, %v)", i, ptP.X, ptP.Y, ptQ.X, ptQ.Y)
			}
		}
	})
}

// TestPath_AddPathMode tests AddPathMode behavior (Append vs Extend)
// Ported from: skia-source/tests/PathTest.cpp:test_addPathMode()
func TestPath_AddPathMode(t *testing.T) {
	testCases := []struct {
		name          string
		explicitMoveTo bool
		extend        bool
		expectedVerbs []enums.PathVerb
	}{
		{
			name:          "append_without_explicit_move",
			explicitMoveTo: false,
			extend:        false,
			expectedVerbs: []enums.PathVerb{enums.PathVerbMove, enums.PathVerbLine, enums.PathVerbMove, enums.PathVerbLine},
		},
		{
			name:          "append_with_explicit_move",
			explicitMoveTo: true,
			extend:        false,
			expectedVerbs: []enums.PathVerb{enums.PathVerbMove, enums.PathVerbLine, enums.PathVerbMove, enums.PathVerbLine},
		},
		{
			name:          "extend_without_explicit_move",
			explicitMoveTo: false,
			extend:        true,
			expectedVerbs: []enums.PathVerb{enums.PathVerbMove, enums.PathVerbLine, enums.PathVerbLine, enums.PathVerbLine},
		},
		{
			name:          "extend_with_explicit_move",
			explicitMoveTo: true,
			extend:        true,
			expectedVerbs: []enums.PathVerb{enums.PathVerbMove, enums.PathVerbLine, enums.PathVerbLine, enums.PathVerbLine},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewSkPath(enums.PathFillTypeDefault)
			q := NewSkPath(enums.PathFillTypeDefault)

			if tc.explicitMoveTo {
				p.MoveTo(1, 1)
			}
			p.LineTo(1, 2)

			if tc.explicitMoveTo {
				q.MoveTo(2, 1)
			}
			q.LineTo(2, 2)

			mode := enums.AddPathModeAppend
			if tc.extend {
				mode = enums.AddPathModeExtend
			}
			p.AddPathNoOffset(q, mode)

			verbs := make([]enums.PathVerb, p.CountVerbs())
			p.GetVerbs(verbs)

			if len(verbs) != len(tc.expectedVerbs) {
				t.Errorf("Verb count: got %d, want %d", len(verbs), len(tc.expectedVerbs))
				return
			}

			for i, expectedVerb := range tc.expectedVerbs {
				if verbs[i] != expectedVerb {
					t.Errorf("Verb %d: got %v, want %v", i, verbs[i], expectedVerb)
				}
			}
		})
	}
}

// TestPath_ExtendClosedPath tests extending a closed path
// Ported from: skia-source/tests/PathTest.cpp:test_extendClosedPath()
func TestPath_ExtendClosedPath(t *testing.T) {
	q := NewSkPath(enums.PathFillTypeDefault)
	q.MoveTo(2, 1)
	q.LineTo(2, 3)

	p := NewSkPath(enums.PathFillTypeDefault)
	p.MoveTo(1, 1)
	p.LineTo(1, 2)
	p.LineTo(2, 2)
	p.Close()
	p.AddPathNoOffset(q, enums.AddPathModeExtend)

	// Check verb sequence: Move, Line, Line, Close, Move, Line, Line
	verbs := make([]enums.PathVerb, p.CountVerbs())
	p.GetVerbs(verbs)
	expectedVerbs := []enums.PathVerb{
		enums.PathVerbMove,
		enums.PathVerbLine,
		enums.PathVerbLine,
		enums.PathVerbClose,
		enums.PathVerbMove,
		enums.PathVerbLine,
		enums.PathVerbLine,
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

	// Check last point
	lastPt, ok := p.GetLastPoint()
	if !ok {
		t.Error("Expected last point to exist")
	} else {
		expectedLastPt := models.Point{X: 2, Y: 3}
		if !NearlyEqualScalar(lastPt.X, expectedLastPt.X) || !NearlyEqualScalar(lastPt.Y, expectedLastPt.Y) {
			t.Errorf("Last point: got (%v, %v), want (%v, %v)", lastPt.X, lastPt.Y, expectedLastPt.X, expectedLastPt.Y)
		}
	}

	// Check that point at index 3 is the move point (1, 1)
	movePt := p.Point(3)
	expectedMovePt := models.Point{X: 1, Y: 1}
	if !NearlyEqualScalar(movePt.X, expectedMovePt.X) || !NearlyEqualScalar(movePt.Y, expectedMovePt.Y) {
		t.Errorf("Move point at index 3: got (%v, %v), want (%v, %v)", movePt.X, movePt.Y, expectedMovePt.X, expectedMovePt.Y)
	}
}

// TestPath_AddEmptyPath tests adding empty paths
// Ported from: skia-source/tests/PathTest.cpp:test_addEmptyPath()
func TestPath_AddEmptyPath(t *testing.T) {
	testCases := []struct {
		name string
		mode enums.AddPathMode
	}{
		{"append_mode", enums.AddPathModeAppend},
		{"extend_mode", enums.AddPathModeExtend},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Case 1: dst is empty
			p := NewSkPath(enums.PathFillTypeDefault)
			p.MoveTo(2, 1)
			p.LineTo(2, 3)
			q := NewSkPath(enums.PathFillTypeDefault)
			q.AddPathNoOffset(p, tc.mode)

			// q should now equal p
			if q.CountPoints() != p.CountPoints() {
				t.Errorf("Case 1 - Point count: got %d, want %d", q.CountPoints(), p.CountPoints())
			}

			// Case 2: src is empty
			r := NewSkPath(enums.PathFillTypeDefault) // empty
			originalPointCount := p.CountPoints()
			p.AddPathNoOffset(r, tc.mode)

			// p should be unchanged
			if p.CountPoints() != originalPointCount {
				t.Errorf("Case 2 - Point count: got %d, want %d", p.CountPoints(), originalPointCount)
			}

			// Case 3: src and dst are empty
			emptyDst := NewSkPath(enums.PathFillTypeDefault)
			emptySrc := NewSkPath(enums.PathFillTypeDefault)
			emptyDst.AddPathNoOffset(emptySrc, tc.mode)

			if !emptyDst.IsEmpty() {
				t.Error("Case 3 - Expected empty path after adding empty to empty")
			}
		})
	}
}

