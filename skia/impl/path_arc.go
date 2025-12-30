// Ported from: skia-source/src/core/SkPathBuilder.cpp (arc methods)
// https://github.com/google/skia/blob/main/src/core/SkPathBuilder.cpp

package impl

import (
	"math"

	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/geometry"
	"github.com/zodimo/go-skia-support/skia/interfaces"
)

// Conic represents a conic curve (weighted quadratic bezier)
// Used internally for arc construction
type Conic struct {
	Pts [3]Point
	W   Scalar
}

// MaxConicsForArc is the maximum number of conics needed to represent any arc
const MaxConicsForArc = 5

// arcIsLonePoint checks if the arc degenerates to a single point
// Ported from: SkPathBuilder.cpp arc_is_lone_point
func arcIsLonePoint(oval Rect, startAngle, sweepAngle Scalar) (Point, bool) {
	if sweepAngle == 0 && (startAngle == 0 || startAngle == 360) {
		// Chrome uses this path to move into and out of ovals
		return Point{X: oval.Right, Y: (oval.Top + oval.Bottom) / 2}, true
	}
	width := oval.Right - oval.Left
	height := oval.Bottom - oval.Top
	if width == 0 && height == 0 {
		// Zero-size oval
		return Point{X: oval.Right, Y: oval.Top}, true
	}
	return Point{}, false
}

// anglesToUnitVectors converts start and sweep angles to unit vectors
// Ported from: SkPathBuilder.cpp angles_to_unit_vectors
func anglesToUnitVectors(startAngle, sweepAngle Scalar) (startV, stopV Point, dir enums.PathDirection) {
	startRad := degreesToRadians(startAngle)
	stopRad := degreesToRadians(startAngle + sweepAngle)

	startV.Y = scalarSinSnapToZero(startRad)
	startV.X = scalarCosSnapToZero(startRad)
	stopV.Y = scalarSinSnapToZero(stopRad)
	stopV.X = scalarCosSnapToZero(stopRad)

	// Handle nearly complete circles
	if startV == stopV {
		sw := Scalar(math.Abs(float64(sweepAngle)))
		if sw < 360 && sw > 359 {
			// Tweak the stop vector for nearly complete circles
			deltaRad := Scalar(1.0 / 512.0)
			if sweepAngle < 0 {
				deltaRad = -deltaRad
			}
			for startV == stopV {
				stopRad -= deltaRad
				stopV.Y = scalarSinSnapToZero(stopRad)
				stopV.X = scalarCosSnapToZero(stopRad)
			}
		}
	}

	if sweepAngle > 0 {
		dir = enums.PathDirectionCW
	} else {
		dir = enums.PathDirectionCCW
	}
	return startV, stopV, dir
}

// buildUnitArc builds unit arc conics for the given start/stop vectors
// Ported from: SkConic::BuildUnitArc in SkGeometry.cpp
func buildUnitArc(startV, stopV Point, dir enums.PathDirection, conics []Conic) int {
	// Compute the bisector
	x := startV.X + stopV.X
	y := startV.Y + stopV.Y
	absX := Scalar(math.Abs(float64(x)))
	absY := Scalar(math.Abs(float64(y)))

	if scalarNearlyZero(absX) && scalarNearlyZero(absY) {
		// Nearly opposite direction - needs 2 arcs
		return 0
	}

	// Normalize bisector
	length := Scalar(math.Sqrt(float64(x*x + y*y)))
	x /= length
	y /= length

	// The angle of each conic
	cosThetaOver2 := Scalar(math.Sqrt(float64((1 + x*startV.X + y*startV.Y) / 2)))
	if scalarNearlyZero(cosThetaOver2) {
		return 0
	}

	// Single conic case
	conics[0].Pts[0] = startV
	conics[0].Pts[2] = stopV
	conics[0].Pts[1] = Point{X: x, Y: y}
	conics[0].W = cosThetaOver2
	return 1
}

// buildArcConics builds conic curves to represent an arc on an oval
// Uses geometry.BuildUnitArc for proper multi-segment arc support
func buildArcConics(oval Rect, startV, stopV Point, dir enums.PathDirection) ([]Conic, Point, int) {
	// Use geometry.BuildUnitArc for proper multi-quadrant support
	// Types are unified: geometry.Point == models.Point == impl.Point
	geoConics := geometry.BuildUnitArc(startV, stopV, dir, nil)

	rx := (oval.Right - oval.Left) / 2
	ry := (oval.Bottom - oval.Top) / 2
	cx := (oval.Left + oval.Right) / 2
	cy := (oval.Top + oval.Bottom) / 2

	if len(geoConics) == 0 {
		// Degenerate case - return single point
		singlePt := Point{X: cx + rx*stopV.X, Y: cy + ry*stopV.Y}
		return nil, singlePt, 0
	}

	// Convert geometry.Conic to impl.Conic and transform to oval
	conics := make([]Conic, len(geoConics))
	for i, gc := range geoConics {
		conics[i] = Conic{
			Pts: [3]Point{
				{X: cx + rx*gc.Pts[0].X, Y: cy + ry*gc.Pts[0].Y},
				{X: cx + rx*gc.Pts[1].X, Y: cy + ry*gc.Pts[1].Y},
				{X: cx + rx*gc.Pts[2].X, Y: cy + ry*gc.Pts[2].Y},
			},
			W: gc.W,
		}
	}

	return conics, Point{}, len(conics)
}

// ArcTo appends arc from oval from startAngle through sweepAngle.
// Ported from: SkPathBuilder.cpp arcTo(oval, startAngle, sweepAngle, forceMoveTo)
func (p *pathImpl) ArcTo(oval interfaces.Rect, startAngle, sweepAngle Scalar, forceMoveTo bool) {
	width := oval.Right - oval.Left
	height := oval.Bottom - oval.Top
	if width < 0 || height < 0 {
		return
	}

	startAngle = Scalar(math.Mod(float64(startAngle), 360.0))

	if len(p.verbs) == 0 {
		forceMoveTo = true
	}

	// Check for lone point case
	if lonePt, isLone := arcIsLonePoint(oval, startAngle, sweepAngle); isLone {
		if forceMoveTo {
			p.MoveTo(lonePt.X, lonePt.Y)
		} else {
			p.LineTo(lonePt.X, lonePt.Y)
		}
		return
	}

	startV, stopV, dir := anglesToUnitVectors(startAngle, sweepAngle)

	// Handle case where vectors are equal (very small sweep angle)
	if startV == stopV {
		endAngle := degreesToRadians(startAngle + sweepAngle)
		rx := width / 2
		ry := height / 2
		cx := (oval.Left + oval.Right) / 2
		cy := (oval.Top + oval.Bottom) / 2
		singlePt := Point{
			X: cx + rx*Scalar(math.Cos(float64(endAngle))),
			Y: cy + ry*Scalar(math.Sin(float64(endAngle))),
		}
		p.addArcPoint(singlePt, forceMoveTo)
		return
	}

	conics, singlePt, count := buildArcConics(oval, startV, stopV, dir)
	if count > 0 {
		pt := conics[0].Pts[0]
		p.addArcPoint(pt, forceMoveTo)
		for i := 0; i < count; i++ {
			p.ConicTo(conics[i].Pts[1].X, conics[i].Pts[1].Y,
				conics[i].Pts[2].X, conics[i].Pts[2].Y, conics[i].W)
		}
	} else {
		p.addArcPoint(singlePt, forceMoveTo)
	}
}

// addArcPoint adds a point to the path, either as moveTo or lineTo
func (p *pathImpl) addArcPoint(pt Point, forceMoveTo bool) {
	if forceMoveTo {
		p.MoveTo(pt.X, pt.Y)
	} else if lastPt, ok := p.GetLastPoint(); !ok || !nearlyEqual(lastPt, pt) {
		p.LineTo(pt.X, pt.Y)
	}
}

// nearlyEqual checks if two points are nearly equal
func nearlyEqual(a, b interfaces.Point) bool {
	return NearlyEqualScalarDefault(a.X, b.X) && NearlyEqualScalarDefault(a.Y, b.Y)
}

// ArcToTangent appends arc tangent to line from last point through (x1,y1)
// to line from (x1,y1) to (x2,y2), with specified radius.
// Ported from: SkPathBuilder.cpp arcTo(p1, p2, radius)
func (p *pathImpl) ArcToTangent(x1, y1, x2, y2, radius Scalar) {
	p.ensureMove()

	if radius == 0 {
		p.LineTo(x1, y1)
		return
	}

	// Get the last point
	start, ok := p.GetLastPoint()
	if !ok {
		start = Point{X: 0, Y: 0}
	}

	p1 := Point{X: x1, Y: y1}
	p2 := Point{X: x2, Y: y2}

	// Compute normalized direction vectors
	before := normalize(Point{X: p1.X - start.X, Y: p1.Y - start.Y})
	after := normalize(Point{X: p2.X - p1.X, Y: p2.Y - p1.Y})

	// Check for degenerate cases
	if !isFinitePoint(before) || !isFinitePoint(after) {
		p.LineTo(x1, y1)
		return
	}

	// Compute cross product (sinh) and dot product (cosh)
	cosh := before.X*after.X + before.Y*after.Y
	sinh := before.X*after.Y - before.Y*after.X

	// If nearly parallel, just draw a line
	if NearlyEqualScalarDefault(sinh, 0) {
		p.LineTo(x1, y1)
		return
	}

	// Compute the distance along the tangent directions
	dist := Scalar(math.Abs(float64(radius * (1 - cosh) / sinh)))
	xx := p1.X - dist*before.X
	yy := p1.Y - dist*before.Y

	// Extend after vector to the arc endpoint
	afterScaled := Point{X: after.X * dist, Y: after.Y * dist}

	p.LineTo(xx, yy)

	// Compute conic weight
	weight := Scalar(math.Sqrt(float64(0.5 + cosh*0.5)))
	endPt := Point{X: p1.X + afterScaled.X, Y: p1.Y + afterScaled.Y}
	p.ConicTo(p1.X, p1.Y, endPt.X, endPt.Y, weight)
}

// normalize returns a normalized (unit length) version of the point/vector
func normalize(pt Point) Point {
	length := Scalar(math.Sqrt(float64(pt.X*pt.X + pt.Y*pt.Y)))
	if scalarNearlyZero(length) {
		return Point{X: 0, Y: 0}
	}
	return Point{X: pt.X / length, Y: pt.Y / length}
}

// isFinitePoint checks if both coordinates are finite
func isFinitePoint(pt Point) bool {
	return math.IsInf(float64(pt.X), 0) == false &&
		math.IsInf(float64(pt.Y), 0) == false &&
		math.IsNaN(float64(pt.X)) == false &&
		math.IsNaN(float64(pt.Y)) == false &&
		pt.X != 0 || pt.Y != 0
}

// ArcToRotated appends SVG-style elliptical arc to (x,y).
// Ported from: SkPathBuilder.cpp arcTo(rad, angle, arcLarge, sweep, endPt)
func (p *pathImpl) ArcToRotated(rx, ry, xAxisRotate Scalar, largeArc enums.ArcSize, arcSweep enums.PathDirection, x, y Scalar) {
	p.ensureMove()

	endPt := Point{X: x, Y: y}

	// If rx = 0 or ry = 0 then treat as line
	if rx == 0 || ry == 0 {
		p.LineTo(x, y)
		return
	}

	// Get current point
	srcPt, ok := p.GetLastPoint()
	if !ok {
		srcPt = Point{X: 0, Y: 0}
	}

	// If points are identical, treat as zero-length path
	if srcPt.X == endPt.X && srcPt.Y == endPt.Y {
		p.LineTo(x, y)
		return
	}

	rx = Scalar(math.Abs(float64(rx)))
	ry = Scalar(math.Abs(float64(ry)))

	// Compute midpoint
	midPointDistance := Point{
		X: (srcPt.X - endPt.X) / 2,
		Y: (srcPt.Y - endPt.Y) / 2,
	}

	// Rotate to align with axes
	angleRad := degreesToRadians(-xAxisRotate)
	cosAngle := Scalar(math.Cos(float64(angleRad)))
	sinAngle := Scalar(math.Sin(float64(angleRad)))

	transformedMidPoint := Point{
		X: cosAngle*midPointDistance.X + sinAngle*midPointDistance.Y,
		Y: -sinAngle*midPointDistance.X + cosAngle*midPointDistance.Y,
	}

	squareRx := rx * rx
	squareRy := ry * ry
	squareX := transformedMidPoint.X * transformedMidPoint.X
	squareY := transformedMidPoint.Y * transformedMidPoint.Y

	// Scale radii if necessary
	radiiScale := squareX/squareRx + squareY/squareRy
	if radiiScale > 1 {
		radiiScale = Scalar(math.Sqrt(float64(radiiScale)))
		rx *= radiiScale
		ry *= radiiScale
		squareRx = rx * rx
		squareRy = ry * ry
	}

	// Compute center point
	unitPts := [2]Point{
		{X: srcPt.X/rx*cosAngle + srcPt.Y/rx*sinAngle - endPt.X/rx*cosAngle - endPt.Y/rx*sinAngle,
			Y: -srcPt.X/ry*sinAngle + srcPt.Y/ry*cosAngle + endPt.X/ry*sinAngle - endPt.Y/ry*cosAngle},
	}
	// Simplify: use direct computation
	pointTransformScale := func(pt Point) Point {
		return Point{
			X: pt.X/rx*cosAngle + pt.Y/rx*sinAngle,
			Y: -pt.X/ry*sinAngle + pt.Y/ry*cosAngle,
		}
	}

	unitPts[0] = pointTransformScale(srcPt)
	unitPts[1] = pointTransformScale(endPt)

	delta := Point{X: unitPts[1].X - unitPts[0].X, Y: unitPts[1].Y - unitPts[0].Y}
	d := delta.X*delta.X + delta.Y*delta.Y
	scaleFactorSquared := Scalar(math.Max(float64(1/d-0.25), 0))
	scaleFactor := Scalar(math.Sqrt(float64(scaleFactorSquared)))

	// Determine direction based on large arc flag and sweep
	if (arcSweep == enums.PathDirectionCCW) != (largeArc == enums.ArcSizeLarge) {
		scaleFactor = -scaleFactor
	}

	delta.X *= scaleFactor
	delta.Y *= scaleFactor

	centerPoint := Point{
		X: (unitPts[0].X + unitPts[1].X) / 2,
		Y: (unitPts[0].Y + unitPts[1].Y) / 2,
	}
	centerPoint.X -= delta.Y
	centerPoint.Y += delta.X

	unitPts[0].X -= centerPoint.X
	unitPts[0].Y -= centerPoint.Y
	unitPts[1].X -= centerPoint.X
	unitPts[1].Y -= centerPoint.Y

	theta1 := Scalar(math.Atan2(float64(unitPts[0].Y), float64(unitPts[0].X)))
	theta2 := Scalar(math.Atan2(float64(unitPts[1].Y), float64(unitPts[1].X)))
	thetaArc := theta2 - theta1

	if thetaArc < 0 && arcSweep == enums.PathDirectionCW {
		thetaArc += 2 * math.Pi
	} else if thetaArc > 0 && arcSweep != enums.PathDirectionCW {
		thetaArc -= 2 * math.Pi
	}

	// Very tiny angles - just draw a line
	if Scalar(math.Abs(float64(thetaArc))) < Scalar(math.Pi/(1000*1000)) {
		p.LineTo(x, y)
		return
	}

	// Generate conic segments
	segments := int(math.Ceil(math.Abs(float64(thetaArc) / (2 * math.Pi / 3))))
	thetaWidth := thetaArc / Scalar(segments)
	t := Scalar(math.Tan(float64(thetaWidth / 2)))
	if !math.IsInf(float64(t), 0) && !math.IsNaN(float64(t)) {
		w := Scalar(math.Sqrt(0.5 + math.Cos(float64(thetaWidth))*0.5))

		startTheta := theta1
		for i := 0; i < segments; i++ {
			endTheta := startTheta + thetaWidth
			sinEndTheta := scalarSinSnapToZero(endTheta)
			cosEndTheta := scalarCosSnapToZero(endTheta)

			unitEnd := Point{X: cosEndTheta, Y: sinEndTheta}
			unitEnd.X += centerPoint.X
			unitEnd.Y += centerPoint.Y

			unitControl := Point{X: unitEnd.X - centerPoint.X, Y: unitEnd.Y - centerPoint.Y}
			unitControl.X += t * sinEndTheta
			unitControl.Y -= t * cosEndTheta
			unitControl.X += centerPoint.X
			unitControl.Y += centerPoint.Y

			// Transform back to world coordinates
			mappedControl := Point{
				X: rx*(cosAngle*unitControl.X-sinAngle*unitControl.Y) + (srcPt.X+endPt.X)/2 - rx*cosAngle*(unitPts[0].X+centerPoint.X+unitPts[1].X+centerPoint.X)/2 + rx*sinAngle*(unitPts[0].Y+centerPoint.Y+unitPts[1].Y+centerPoint.Y)/2,
				Y: ry*(sinAngle*unitControl.X+cosAngle*unitControl.Y) + (srcPt.Y+endPt.Y)/2 - ry*sinAngle*(unitPts[0].X+centerPoint.X+unitPts[1].X+centerPoint.X)/2 - ry*cosAngle*(unitPts[0].Y+centerPoint.Y+unitPts[1].Y+centerPoint.Y)/2,
			}
			mappedEnd := Point{
				X: rx*(cosAngle*unitEnd.X-sinAngle*unitEnd.Y) + (srcPt.X+endPt.X)/2 - rx*cosAngle*(unitPts[0].X+centerPoint.X+unitPts[1].X+centerPoint.X)/2 + rx*sinAngle*(unitPts[0].Y+centerPoint.Y+unitPts[1].Y+centerPoint.Y)/2,
				Y: ry*(sinAngle*unitEnd.X+cosAngle*unitEnd.Y) + (srcPt.Y+endPt.Y)/2 - ry*sinAngle*(unitPts[0].X+centerPoint.X+unitPts[1].X+centerPoint.X)/2 - ry*cosAngle*(unitPts[0].Y+centerPoint.Y+unitPts[1].Y+centerPoint.Y)/2,
			}

			// Emit the conic arc segment
			p.ConicTo(mappedControl.X, mappedControl.Y, mappedEnd.X, mappedEnd.Y, w)
			startTheta = endTheta
		}
	}

	// Ensure we end at the exact endpoint
	p.points[len(p.points)-1] = endPt
}

// RArcTo appends SVG-style elliptical arc relative to current point.
// Ported from: SkPathBuilder.cpp rArcTo
func (p *pathImpl) RArcTo(rx, ry, xAxisRotate Scalar, largeArc enums.ArcSize, sweep enums.PathDirection, dx, dy Scalar) {
	currentPt, ok := p.GetLastPoint()
	if !ok {
		currentPt = Point{X: 0, Y: 0}
	}
	p.ArcToRotated(rx, ry, xAxisRotate, largeArc, sweep, currentPt.X+dx, currentPt.Y+dy)
}

// AddArc adds arc as a new contour (starts with implicit MoveTo).
// Ported from: SkPathBuilder.cpp addArc
func (p *pathImpl) AddArc(oval interfaces.Rect, startAngle, sweepAngle Scalar) {
	width := oval.Right - oval.Left
	height := oval.Bottom - oval.Top
	if width == 0 || height == 0 || sweepAngle == 0 {
		return
	}

	const fullCircleAngle = 360.0

	if sweepAngle >= fullCircleAngle || sweepAngle <= -fullCircleAngle {
		// Check if we can treat this as an oval
		startOver90 := startAngle / 90.0
		startOver90I := Scalar(math.Round(float64(startOver90)))
		err := startOver90 - startOver90I
		if NearlyEqualScalarDefault(err, 0) {
			// Index 1 is at startAngle == 0
			startIndex := Scalar(math.Mod(float64(startOver90I+1), 4))
			if startIndex < 0 {
				startIndex += 4
			}
			dir := enums.PathDirectionCW
			if sweepAngle < 0 {
				dir = enums.PathDirectionCCW
			}
			_ = startIndex // startIndex would be used for AddOval with start index
			p.AddOval(oval, dir)
			return
		}
	}

	p.ArcTo(oval, startAngle, sweepAngle, true)
}

// ensureMove ensures there's a moveTo before adding geometry
func (p *pathImpl) ensureMove() {
	if len(p.verbs) == 0 || p.verbs[len(p.verbs)-1] == enums.PathVerbClose {
		p.MoveTo(0, 0)
	}
}

// Helper math functions

func degreesToRadians(degrees Scalar) Scalar {
	return degrees * math.Pi / 180.0
}

func scalarSinSnapToZero(radians Scalar) Scalar {
	result := Scalar(math.Sin(float64(radians)))
	if NearlyEqualScalarDefault(result, 0) {
		return 0
	}
	return result
}

func scalarCosSnapToZero(radians Scalar) Scalar {
	result := Scalar(math.Cos(float64(radians)))
	if NearlyEqualScalarDefault(result, 0) {
		return 0
	}
	return result
}

// Note: scalarNearlyZero is defined in matrix_helpers.go
