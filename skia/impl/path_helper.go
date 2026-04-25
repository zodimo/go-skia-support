package impl

import (
	"math"

	"github.com/zodimo/go-skia-support/skia/base"
	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

func pathFirstDirectionToConvexity(dir enums.PathFirstDirection) enums.PathConvexity {
	switch dir {
	case enums.PathFirstDirectionCW:
		return enums.PathConvexityConvexCW
	case enums.PathFirstDirectionCCW:
		return enums.PathConvexityConvexCCW
	case enums.PathFirstDirectionUnknown:
		return enums.PathConvexityConvexDegenerate
	default:
		return enums.PathConvexityUnknown
	}
}

// ptsInVerb returns the number of models.Points (excluding start models.Point) for each verb
func ptsInVerb(verb enums.PathVerb) int {
	switch verb {
	case enums.PathVerbMove:
		return 1
	case enums.PathVerbLine:
		return 1
	case enums.PathVerbQuad:
		return 2
	case enums.PathVerbConic:
		return 2
	case enums.PathVerbCubic:
		return 3
	case enums.PathVerbClose:
		return 0
	default:
		return 0
	}
}

// validUnitDivide performs division and ensures result is in [0, 1)
func validUnitDivide(numer, denom base.Scalar) (base.Scalar, bool) {
	if numer < 0 {
		numer = -numer
		denom = -denom
	}

	if denom == 0 || numer == 0 || numer >= denom {
		return 0, false
	}

	r := numer / denom
	if math.IsNaN(float64(r)) || r == 0 {
		return 0, false
	}

	if r >= 1.0 {
		return 0, false
	}

	return r, true
}

// findUnitQuadRoots finds roots of quadratic equation At^2 + Bt + C = 0 in [0, 1)
// Returns the number of valid roots found
func findUnitQuadRoots(A, B, C base.Scalar, roots []base.Scalar) int {
	if A == 0 {
		if t, ok := validUnitDivide(-C, B); ok {
			roots[0] = t
			return 1
		}
		return 0
	}

	// Use doubles to avoid overflow
	dr := float64(B)*float64(B) - 4*float64(A)*float64(C)
	if dr < 0 {
		return 0
	}
	dr = math.Sqrt(dr)
	R := base.Scalar(dr)
	if math.IsInf(float64(R), 0) || math.IsNaN(float64(R)) {
		return 0
	}

	var Q base.Scalar
	if B < 0 {
		Q = -(B - R) / 2
	} else {
		Q = -(B + R) / 2
	}

	count := 0
	if t, ok := validUnitDivide(Q, A); ok {
		roots[count] = t
		count++
	}
	if t, ok := validUnitDivide(C, Q); ok {
		roots[count] = t
		count++
	}

	// Sort roots
	if count == 2 {
		if roots[0] > roots[1] {
			roots[0], roots[1] = roots[1], roots[0]
		} else if roots[0] == roots[1] {
			// Nearly equal, skip duplicate
			count = 1
		}
	}

	return count
}

// findQuadExtrema finds t values where quadratic curve has extrema
// Quadratic: P(t) = (1-t)^2*P0 + 2*(1-t)*t*P1 + t^2*P2
// Derivative: 2*(P1-P0) + 2*(P2-2*P1+P0)*t = 0
// Solving: t = (P0-P1) / (P0-2*P1+P2)
func findQuadExtrema(a, b, c base.Scalar, tValue []base.Scalar) int {
	// At + B == 0, where A = a - b - b + c, B = a - b
	// t = -B / A = (b - a) / (a - 2*b + c)
	if t, ok := validUnitDivide(a-b, a-b-b+c); ok {
		tValue[0] = t
		return 1
	}
	return 0
}

// evalQuadAt evaluates a quadratic curve at parameter t
// P(t) = (1-t)^2*P0 + 2*(1-t)*t*P1 + t^2*P2
func evalQuadAt(src []models.Point, t base.Scalar) models.Point {
	// Using Bernstein basis: (1-t)^2, 2*(1-t)*t, t^2
	mt := 1 - t
	mt2 := mt * mt
	t2 := t * t
	mt2t := 2 * mt * t

	return models.Point{
		X: mt2*src[0].X + mt2t*src[1].X + t2*src[2].X,
		Y: mt2*src[0].Y + mt2t*src[1].Y + t2*src[2].Y,
	}
}

// computeQuadExtremas computes extrema points for a quadratic curve
func computeQuadExtremas(src []models.Point) ([]models.Point, int) {
	if len(src) < 3 {
		return nil, 0
	}

	ts := make([]base.Scalar, 2)
	n := findQuadExtrema(src[0].X, src[1].X, src[2].X, ts)
	n += findQuadExtrema(src[0].Y, src[1].Y, src[2].Y, ts[n:])

	if n > 2 {
		n = 2
	}

	extremas := make([]models.Point, n+1)
	for i := 0; i < n; i++ {
		extremas[i] = evalQuadAt(src, ts[i])
	}
	extremas[n] = src[2] // Always include end point

	return extremas, n + 1
}

// findCubicExtrema finds t values where cubic curve has extrema
// Cubic: P(t) = (1-t)^3*P0 + 3*(1-t)^2*t*P1 + 3*(1-t)*t^2*P2 + t^3*P3
// Derivative: 3*(P1-P0) + 6*(P2-2*P1+P0)*t + 3*(P3-3*P2+3*P1-P0)*t^2 = 0
func findCubicExtrema(a, b, c, d base.Scalar, tValues []base.Scalar) int {
	// A = d - a + 3*(b - c)
	// B = 2*(a - b - b + c) = 2*(a - 2*b + c)
	// C = b - a
	A := d - a + 3*(b-c)
	B := 2 * (a - 2*b + c)
	C := b - a

	return findUnitQuadRoots(A, B, C, tValues)
}

// evalCubicAt evaluates a cubic curve at parameter t
// P(t) = (1-t)^3*P0 + 3*(1-t)^2*t*P1 + 3*(1-t)*t^2*P2 + t^3*P3
func evalCubicAt(src []models.Point, t base.Scalar) models.Point {
	// Using Bernstein basis
	mt := 1 - t
	mt2 := mt * mt
	mt3 := mt2 * mt
	t2 := t * t
	t3 := t2 * t

	return models.Point{
		X: mt3*src[0].X + 3*mt2*t*src[1].X + 3*mt*t2*src[2].X + t3*src[3].X,
		Y: mt3*src[0].Y + 3*mt2*t*src[1].Y + 3*mt*t2*src[2].Y + t3*src[3].Y,
	}
}

// computeCubicExtremas computes extrema points for a cubic curve
func computeCubicExtremas(src []models.Point) ([]models.Point, int) {
	if len(src) < 4 {
		return nil, 0
	}

	ts := make([]base.Scalar, 4)
	n := findCubicExtrema(src[0].X, src[1].X, src[2].X, src[3].X, ts)
	n += findCubicExtrema(src[0].Y, src[1].Y, src[2].Y, src[3].Y, ts[n:])

	if n > 4 {
		n = 4
	}

	extremas := make([]models.Point, n+1)
	for i := 0; i < n; i++ {
		extremas[i] = evalCubicAt(src, ts[i])
	}
	extremas[n] = src[3] // Always include end point

	return extremas, n + 1
}

// conicDerivCoeff computes derivative coefficients for conic curve
// Conic: P(t) = [(1-t)^2*P0 + 2*w*(1-t)*t*P1 + t^2*P2] / [(1-t)^2 + 2*w*(1-t)*t + t^2]
// This computes the coefficients for the derivative numerator: coeff[0]*t^2 + coeff[1]*t + coeff[2]
// src is a 3-element array representing [P0_coord, P1_coord, P2_coord] for a single coordinate (X or Y)
func conicDerivCoeff(src [3]base.Scalar, w base.Scalar) [3]base.Scalar {
	P20 := src[2] - src[0]
	P10 := src[1] - src[0]
	wP10 := w * P10
	return [3]base.Scalar{
		w*P20 - P20,  // coeff[0] for t^2
		P20 - 2*wP10, // coeff[1] for t^1
		wP10,         // coeff[2] for t^0
	}
}

// conicFindExtrema finds extrema for conic curve (X or Y component)
func conicFindExtrema(src []models.Point, w base.Scalar, isX bool) (base.Scalar, bool) {
	if len(src) < 3 {
		return 0, false
	}

	// Extract the coordinate values for the three points
	var coordSrc [3]base.Scalar
	if isX {
		coordSrc = [3]base.Scalar{src[0].X, src[1].X, src[2].X}
	} else {
		coordSrc = [3]base.Scalar{src[0].Y, src[1].Y, src[2].Y}
	}

	coeff := conicDerivCoeff(coordSrc, w)

	tValues := make([]base.Scalar, 2)
	roots := findUnitQuadRoots(coeff[0], coeff[1], coeff[2], tValues)
	if roots == 1 {
		return tValues[0], true
	}
	return 0, false
}

// evalConicAt evaluates a conic curve at parameter t
// P(t) = [(1-t)^2*P0 + 2*w*(1-t)*t*P1 + t^2*P2] / [(1-t)^2 + 2*w*(1-t)*t + t^2]
func evalConicAt(src []models.Point, w base.Scalar, t base.Scalar) models.Point {
	mt := 1 - t
	mt2 := mt * mt
	t2 := t * t
	mt2t := 2 * w * mt * t

	numerX := mt2*src[0].X + mt2t*src[1].X + t2*src[2].X
	numerY := mt2*src[0].Y + mt2t*src[1].Y + t2*src[2].Y
	denom := mt2 + mt2t + t2

	if denom == 0 {
		return src[2] // Fallback to end point
	}

	return models.Point{
		X: numerX / denom,
		Y: numerY / denom,
	}
}

// computeConicExtremas computes extrema points for a conic curve
func computeConicExtremas(src []models.Point, w base.Scalar) ([]models.Point, int) {
	if len(src) < 3 {
		return nil, 0
	}

	ts := make([]base.Scalar, 2)
	n := 0

	if t, ok := conicFindExtrema(src, w, true); ok {
		ts[n] = t
		n++
	}
	if t, ok := conicFindExtrema(src, w, false); ok {
		// Check if this t is different from the X extrema
		isNew := true
		for i := 0; i < n; i++ {
			if ts[i] == t {
				isNew = false
				break
			}
		}
		if isNew {
			ts[n] = t
			n++
		}
	}

	if n > 2 {
		n = 2
	}

	extremas := make([]models.Point, n+1)
	for i := 0; i < n; i++ {
		extremas[i] = evalConicAt(src, w, ts[i])
	}
	extremas[n] = src[2] // Always include end point

	return extremas, n + 1
}

////////////////////////////////////////

func PathFillTypeIsInverse(ft enums.PathFillType) bool {
	return (int(ft) & 2) != 0
}

func PathFillTypeToggleInverse(ft enums.PathFillType) enums.PathFillType {
	return enums.PathFillType(int(ft) ^ 2)
}

func PathConvexityIsConvex(cv enums.PathConvexity) bool {
	return cv == enums.PathConvexityConvexCW || cv == enums.PathConvexityConvexCCW || cv == enums.PathConvexityConvexDegenerate
}

func IsFinite(f base.Scalar) bool {
	return !math.IsNaN(float64(f)) && !math.IsInf(float64(f), 0)
}

var PathVerbs = []enums.PathVerb{
	enums.PathVerbMove,
	enums.PathVerbLine,
	enums.PathVerbQuad,
	enums.PathVerbConic,
	enums.PathVerbCubic,
	enums.PathVerbClose,
}

func RectPathRaw(rect models.Rect, dir enums.PathDirection, startIndex uint) PathRaw {
	// Keep startIndex legal (0-3)
	startIndex = startIndex % 4

	// Create iterator for rectangle points
	iter := newRectPointIterator(rect, dir, startIndex)

	// Rectangle path: Move, Line, Line, Line, Close
	// 4 points total (one for each corner)
	points := make([]models.Point, 4)
	points[0] = iter.current()
	points[1] = iter.next()
	points[2] = iter.next()
	points[3] = iter.next()

	// Verbs: Move, Line, Line, Line, Close
	verbs := []enums.PathVerb{
		enums.PathVerbMove,
		enums.PathVerbLine,
		enums.PathVerbLine,
		enums.PathVerbLine,
		enums.PathVerbClose,
	}
	// Point indices: Move uses models.Point 0, each Line uses the next models.Point
	// For Move: uses models.Points[0]
	// For Line 1: uses models.Points[1] (which is models.PointIndices[1] = 0, then +1 = 1)
	// For Line 2: uses models.Points[2] (which is models.PointIndices[2] = 1, then +1 = 2)
	// For Line 3: uses models.Points[3] (which is models.PointIndices[3] = 2, then +1 = 3)
	// For Close: no models.Point needed
	pointIndices := []int{0, 0, 1, 2, -1} // -1 for Close (not used)

	return PathRaw{
		Points:       points,
		Verbs:        verbs,
		ConicWeights: nil, // Rectangles don't use conic weights
		PointIndices: pointIndices,
		ConicIndex:   nil, // Rectangles don't use conic indices
	}
}

func OvalPathRaw(rect models.Rect, dir enums.PathDirection, startIndex uint) PathRaw {
	// Keep startIndex legal (0-3)
	startIndex = startIndex % 4

	// Create iterators for oval and rectangle points
	ovalIter := newOvalPointIterator(rect, dir, startIndex)

	// Rect iterator starts at index + (dir == CW ? 0 : 1) to align properly
	rectStartIndex := startIndex
	if dir == enums.PathDirectionCCW {
		rectStartIndex = (startIndex + 1) % 4
	}
	rectIter := newRectPointIterator(rect, dir, rectStartIndex)

	// Oval path: Move, Conic, Conic, Conic, Conic, Close
	// 9 models.Points total: 1 start models.Point + 4 conics (each needs 2 models.Points: control + end)
	points := make([]models.Point, 9)
	points[0] = ovalIter.current()
	for i := 0; i < 4; i++ {
		points[i*2+1] = rectIter.next() // control point (rectangle corner)
		points[i*2+2] = ovalIter.next() // end point (oval midpoint)
	}

	// Verbs: Move, Conic, Conic, Conic, Conic, Close
	verbs := []enums.PathVerb{
		enums.PathVerbMove,
		enums.PathVerbConic,
		enums.PathVerbConic,
		enums.PathVerbConic,
		enums.PathVerbConic,
		enums.PathVerbClose,
	}

	// Conic weights: all 4 use sqrt(2)/2 for quarter-circle approximation
	conicWeights := []base.Scalar{
		base.ScalarRoot2Over2,
		base.ScalarRoot2Over2,
		base.ScalarRoot2Over2,
		base.ScalarRoot2Over2,
	}

	// Point indices: Move uses models.Point 0, each Conic uses control models.Point index
	// For Move: uses models.Points[0]
	// For Conic 0: uses models.Points[1] (control) and models.Points[2] (end)
	// For Conic 1: uses models.Points[3] (control) and models.Points[4] (end)
	// For Conic 2: uses models.Points[5] (control) and models.Points[6] (end)
	// For Conic 3: uses models.Points[7] (control) and models.Points[8] (end)
	// For Close: no models.Point needed
	pointIndices := []int{0, 1, 3, 5, 7, -1} // -1 for Close (not used)

	// Conic index: maps verb index to conic weight index
	// Verbs: [Move, Conic, Conic, Conic, Conic, Close]
	// ConicIndex: [-, 0, 1, 2, 3, -]
	conicIndex := []int{-1, 0, 1, 2, 3, -1} // -1 for non-conic verbs

	return PathRaw{
		Points:       points,
		Verbs:        verbs,
		ConicWeights: conicWeights,
		PointIndices: pointIndices,
		ConicIndex:   conicIndex,
	}
}

func RRectPathRaw(rrect models.RRect, dir enums.PathDirection, startIndex uint) PathRaw {
	// Keep startIndex legal (0-7)
	startIndex = startIndex % 8

	// Determine if we start with a conic or a line
	// We start with a conic on odd indices when moving CW vs. even indices when moving CCW
	startsWithConic := ((startIndex & 1) == 1) == (dir == enums.PathDirectionCW)

	// If we start with a conic, we end with a line, which we can skip (relying on close())
	npoints := 13
	if startsWithConic {
		npoints = 12 // skip the last line point
	}

	// Create iterators
	rrectIter := newRRectPointIterator(rrect, dir, startIndex)
	// Corner iterator indices follow the collapsed radii model,
	// adjusted such that the start pt is "behind" the radii start pt.
	rectStartIndex := startIndex / 2
	if dir == enums.PathDirectionCCW {
		rectStartIndex = (startIndex/2 + 1) % 4
	}
	rectIter := newRectPointIterator(rrect.Bounds(), dir, uint(rectStartIndex))

	// Build points array
	points := make([]models.Point, npoints)
	points[0] = rrectIter.current()

	if startsWithConic {
		// Pattern: Conic, Line, Conic, Line, Conic, Line, Conic, Close
		// points: start, (conic_ctrl, conic_end, line), (conic_ctrl, conic_end, line), ...
		for i := 0; i < 3; i++ {
			// Conic points
			points[i*3+1] = rectIter.next()  // control point (rectangle corner)
			points[i*3+2] = rrectIter.next() // end point (rounded corner end)
			// Line point
			points[i*3+3] = rrectIter.next() // line end point
		}
		// Last conic points
		points[10] = rectIter.next()  // control point
		points[11] = rrectIter.next() // end point
		// The final line is accomplished by close()
	} else {
		// Pattern: Line, Conic, Line, Conic, Line, Conic, Line, Conic, Close
		// points: start, (line, conic_ctrl, conic_end), (line, conic_ctrl, conic_end), ...
		for i := 0; i < 4; i++ {
			// Line point
			points[i*3+1] = rrectIter.next() // line end point
			// Conic points
			points[i*3+2] = rectIter.next()  // control point (rectangle corner)
			points[i*3+3] = rrectIter.next() // end point (rounded corner end)
		}
	}

	// Build verbs array
	var verbs []enums.PathVerb
	if startsWithConic {
		// Conic, Line, Conic, Line, Conic, Line, Conic, Close
		verbs = []enums.PathVerb{
			enums.PathVerbMove,
			enums.PathVerbConic, enums.PathVerbLine,
			enums.PathVerbConic, enums.PathVerbLine,
			enums.PathVerbConic, enums.PathVerbLine,
			enums.PathVerbConic, // last line skipped
			enums.PathVerbClose,
		}
	} else {
		// Line, Conic, Line, Conic, Line, Conic, Line, Conic, Close
		verbs = []enums.PathVerb{
			enums.PathVerbMove,
			enums.PathVerbLine, enums.PathVerbConic,
			enums.PathVerbLine, enums.PathVerbConic,
			enums.PathVerbLine, enums.PathVerbConic,
			enums.PathVerbLine, enums.PathVerbConic,
			enums.PathVerbClose,
		}
	}

	// Conic weights: all use sqrt(2)/2 for quarter-circle approximation
	numConics := 4
	conicWeights := make([]base.Scalar, numConics)
	for i := 0; i < numConics; i++ {
		conicWeights[i] = base.ScalarRoot2Over2
	}

	// Build point indices array
	// Point indices map verb index to the starting models.Point index in the models.Points array
	// For Line verbs, addRaw uses models.PointIndices[i]+1, so models.PointIndices[i] should point to the models.Point BEFORE the line end
	var pointIndices []int
	if startsWithConic {
		// Verbs: [Move(0), Conic(1), Line(2), Conic(3), Line(4), Conic(5), Line(6), Conic(7), Close(8)]
		// models.Points: [0=start, 1=conic0_ctrl, 2=conic0_end, 3=line0_end, 4=conic1_ctrl, 5=conic1_end, 6=line1_end, 7=conic2_ctrl, 8=conic2_end, 9=line2_end, 10=conic3_ctrl, 11=conic3_end]
		// Move(0): uses models.Points[0] -> models.PointIndices[0] = 0
		// Conic(1): uses models.Points[1] (ctrl) and models.Points[2] (end) -> models.PointIndices[1] = 1
		// Line(2): uses models.Points[3] (end) -> models.PointIndices[2] = 2 (point before line end)
		// Conic(3): uses models.Points[4] (ctrl) and models.Points[5] (end) -> models.PointIndices[3] = 4
		// Line(4): uses models.Points[6] (end) -> models.PointIndices[4] = 5 (point before line end)
		// Conic(5): uses models.Points[7] (ctrl) and models.Points[8] (end) -> models.PointIndices[5] = 7
		// Line(6): uses models.Points[9] (end) -> models.PointIndices[6] = 8 (point before line end)
		// Conic(7): uses models.Points[10] (ctrl) and models.Points[11] (end) -> models.PointIndices[7] = 10
		// Close(8): no models.Point -> models.PointIndices[8] = -1
		pointIndices = []int{0, 1, 2, 4, 5, 7, 8, 10, -1}
	} else {
		// Verbs: [Move(0), Line(1), Conic(2), Line(3), Conic(4), Line(5), Conic(6), Line(7), Conic(8), Close(9)]
		// models.Points: [0=start, 1=line0_end, 2=conic0_ctrl, 3=conic0_end, 4=line1_end, 5=conic1_ctrl, 6=conic1_end, 7=line2_end, 8=conic2_ctrl, 9=conic2_end, 10=line3_end, 11=conic3_ctrl, 12=conic3_end]
		// Move(0): uses models.Points[0] -> models.PointIndices[0] = 0
		// Line(1): uses models.Points[1] (end) -> models.PointIndices[1] = 0 (point before line end, which is start)
		// Conic(2): uses models.Points[2] (ctrl) and models.Points[3] (end) -> models.PointIndices[2] = 2
		// Line(3): uses models.Points[4] (end) -> models.PointIndices[3] = 3 (point before line end)
		// Conic(4): uses models.Points[5] (ctrl) and models.Points[6] (end) -> models.PointIndices[4] = 5
		// Line(5): uses models.Points[7] (end) -> models.PointIndices[5] = 6 (point before line end)
		// Conic(6): uses models.Points[8] (ctrl) and models.Points[9] (end) -> models.PointIndices[6] = 8
		// Line(7): uses models.Points[10] (end) -> models.PointIndices[7] = 9 (point before line end)
		// Conic(8): uses models.Points[11] (ctrl) and models.Points[12] (end) -> models.PointIndices[8] = 11
		// Close(9): no models.Point -> models.PointIndices[9] = -1
		pointIndices = []int{0, 0, 2, 3, 5, 6, 8, 9, 11, -1}
	}

	// Build conic index array
	// Maps verb index to conic weight index
	var conicIndex []int
	if startsWithConic {
		// Verbs: [Move, Conic, Line, Conic, Line, Conic, Line, Conic, Close]
		// ConicIndex: [-, 0, -, 1, -, 2, -, 3, -]
		conicIndex = []int{-1, 0, -1, 1, -1, 2, -1, 3, -1}
	} else {
		// Verbs: [Move, Line, Conic, Line, Conic, Line, Conic, Line, Conic, Close]
		// ConicIndex: [-, -, 0, -, 1, -, 2, -, 3, -]
		conicIndex = []int{-1, -1, 0, -1, 1, -1, 2, -1, 3, -1}
	}

	return PathRaw{
		Points:       points,
		Verbs:        verbs,
		ConicWeights: conicWeights,
		PointIndices: pointIndices,
		ConicIndex:   conicIndex,
	}
}

func isConcaveBySign(points []models.Point) bool {
	if len(points) <= 3 {
		// Point, line, or triangle are always convex
		return false
	}

	dxes := 0
	dyes := 0
	lastSx := 999 // kValueNeverReturnedBySign
	lastSy := 999

	// Check twice: first pass from models.Points[1] to end, second pass processes first edge only
	// This matches C++ implementation: counters and lastSx/lastSy accumulate across both passes
	currPt := points[0]
	firstPt := currPt
	pointIdx := 1 // Start from second models.Point (points[1])

	for outerLoop := 0; outerLoop < 2; outerLoop++ {
		for pointIdx < len(points) {
			vec := models.Point{X: points[pointIdx].X - currPt.X, Y: points[pointIdx].Y - currPt.Y}
			if vec.X != 0 || vec.Y != 0 {
				// Give up if vector construction failed
				if !IsFinite(vec.X) || !IsFinite(vec.Y) {
					return true // treat as concave
				}
				sx := sign(vec.X)
				sy := sign(vec.Y)
				if sx != lastSx {
					dxes++
					if dxes > 3 {
						return true
					}
				}
				if sy != lastSy {
					dyes++
					if dyes > 3 {
						return true
					}
				}
				lastSx = sx
				lastSy = sy
			}
			currPt = points[pointIdx]
			pointIdx++

			// In C++, the second pass breaks after first iteration
			if outerLoop == 1 {
				break
			}
		}
		// Second pass: reset point index to 0 (start from first point)
		if outerLoop == 0 {
			currPt = firstPt
			pointIdx = 0
		}
	}
	return false // may be convex, don't know yet
}

// affectsAlphaColorFilter returns true if the color filter exists and may affect alpha
// It checks if the color filter implements IsAlphaUnchanged() method
func affectsAlphaColorFilter(cf interfaces.ColorFilter) bool {
	if cf == nil {
		return false
	}
	// If alpha is unchanged, it doesn't affect alpha
	return !cf.IsAlphaUnchanged()
}

// affectsAlphaImageFilter returns true if the image filter exists and may affect alpha
// For now, if an image filter exists, it affects alpha (as per C++ TODO comment)
func affectsAlphaImageFilter(imf interfaces.ImageFilter) bool {
	return imf != nil
}
