package impl

import (
	"github.com/zodimo/go-skia-support/skia/enums"
)

// PathIterRec represents a single iteration result from PathIter
// Ported from: skia-source/include/core/SkPathIter.h:Rec
type PathIterRec struct {
	Points      []Point
	ConicWeight Scalar
	Verb        enums.PathVerb
}

// PathIter iterates through path verbs and points, handling implicit lines for Close
// Ported from: skia-source/src/core/SkPathIter.cpp
type PathIter struct {
	pIndex                int
	vIndex                int
	cIndex                int
	points                []Point
	verbs                 []enums.PathVerb
	conics                []Scalar
	closePointStorage     [2]Point
	firstPointFromMove    Point
	hasFirstPointFromMove bool
}

// NewPathIter creates a new PathIter for the given path data
// Ported from: skia-source/src/core/SkPathIter.cpp:SkPathIter constructor
func NewPathIter(points []Point, verbs []enums.PathVerb, conics []Scalar) *PathIter {
	// Trim trailing Move verb for compatibility (as C++ does)
	trimmedVerbs := verbs
	if len(verbs) > 0 && verbs[len(verbs)-1] == enums.PathVerbMove {
		trimmedVerbs = verbs[:len(verbs)-1]
	}

	return &PathIter{
		pIndex: 0,
		vIndex: 0,
		cIndex: 0,
		points: points,
		verbs:  trimmedVerbs,
		conics: conics,
	}
}

// Next returns the next iteration result, or nil if done
// Ported from: skia-source/src/core/SkPathIter.cpp:next()
// Close is funny -- it has no explicit point data, but we return 2 points,
// the logical 2 points that would make up the line connecting the end of
// the contour, and its beginning.
func (iter *PathIter) Next() *PathIterRec {
	if iter.vIndex >= len(iter.verbs) {
		return nil
	}

	var n int
	var w Scalar = -1
	v := iter.verbs[iter.vIndex]
	iter.vIndex++

	switch v {
	case enums.PathVerbMove:
		if iter.pIndex < len(iter.points) {
			iter.closePointStorage[1] = iter.points[iter.pIndex] // remember for close
			iter.firstPointFromMove = iter.points[iter.pIndex]
			iter.hasFirstPointFromMove = true
			iter.pIndex++
		}
		return &PathIterRec{
			Points:      []Point{iter.closePointStorage[1]},
			ConicWeight: w,
			Verb:        v,
		}
	case enums.PathVerbLine:
		n = 1
	case enums.PathVerbQuad:
		n = 2
	case enums.PathVerbConic:
		n = 2
		if iter.cIndex < len(iter.conics) {
			w = iter.conics[iter.cIndex]
			iter.cIndex++
		}
	case enums.PathVerbCubic:
		n = 3
	case enums.PathVerbClose:
		// Close has no explicit point data, but we return 2 points:
		// [last point we saw, first point from Move]
		if iter.pIndex > 0 {
			iter.closePointStorage[0] = iter.points[iter.pIndex-1] // the last point we saw
		}
		if iter.hasFirstPointFromMove {
			iter.closePointStorage[1] = iter.firstPointFromMove
		}
		return &PathIterRec{
			Points:      iter.closePointStorage[:],
			ConicWeight: w,
			Verb:        v,
		}
	default:
		return nil
	}

	// For Line, Quad, Conic, Cubic: return points starting from last point
	if iter.pIndex == 0 {
		return nil // Invalid state
	}
	start := iter.pIndex - 1
	if start+n >= len(iter.points) {
		return nil // Not enough points
	}

	// Return n+1 points: [start point, ...n additional points]
	resultPoints := make([]Point, n+1)
	resultPoints[0] = iter.points[start]
	for i := 0; i < n; i++ {
		resultPoints[i+1] = iter.points[start+1+i]
	}

	iter.pIndex += n

	return &PathIterRec{
		Points:      resultPoints,
		ConicWeight: w,
		Verb:        v,
	}
}
