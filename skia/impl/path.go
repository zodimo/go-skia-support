package impl

import (
	"github.com/zodimo/go-skia-support/skia/base"
	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/interfaces"
)

var _ interfaces.SkPath = &pathImpl{}

// pathImpl implements SkPath interface
// This is a verbatim port from the C++ SkPath implementation
type pathImpl struct {
	points          []Point
	verbs           []enums.PathVerb
	conicWeights    []Scalar
	fillType        enums.PathFillType
	isVolatile      bool
	convexity       enums.PathConvexity
	lastMoveToIndex int
	bounds          Rect
	boundsDirty     bool
}

const initialLastMoveToIndexValue = ^0

// NewSkPath creates a new empty SkPath with the specified fill type
func NewSkPath(fillType enums.PathFillType) interfaces.SkPath {
	return &pathImpl{
		fillType:        fillType,
		lastMoveToIndex: initialLastMoveToIndexValue,
		convexity:       enums.PathConvexityUnknown,
		boundsDirty:     true,
	}
}

// FillType returns the fill type used to determine which parts are inside.
func (p *pathImpl) FillType() enums.PathFillType {
	return p.fillType
}

// SetFillType sets the fill type used to determine which parts are inside.
func (p *pathImpl) SetFillType(fillType enums.PathFillType) {
	p.fillType = fillType
}

// IsInverseFillType returns true if the fill type is inverse.
func (p *pathImpl) IsInverseFillType() bool {
	return PathFillTypeIsInverse(p.fillType)
}

// ToggleInverseFillType toggles between inverse and non-inverse fill types.
func (p *pathImpl) ToggleInverseFillType() {
	p.fillType = PathFillTypeToggleInverse(p.fillType)
}

// Convexity returns the convexity type of the path.
func (p *pathImpl) Convexity() enums.PathConvexity {
	convexity := p.getConvexityOrUnknown()
	if convexity == enums.PathConvexityUnknown {
		convexity = p.computeConvexity()
	}
	return convexity
}

// IsConvex returns true if the path is convex.
func (p *pathImpl) IsConvex() bool {
	return PathConvexityIsConvex(p.Convexity())
}

// Reset clears the path, removing all verbs, points, and conic weights.
func (p *pathImpl) Reset() {
	p.points = nil
	p.verbs = nil
	p.conicWeights = nil
	p.lastMoveToIndex = initialLastMoveToIndexValue
	p.fillType = enums.PathFillTypeDefault
	p.setConvexity(enums.PathConvexityUnknown)
	p.boundsDirty = true
}

// IsEmpty returns true if the path has no verbs.
func (p *pathImpl) IsEmpty() bool {
	return len(p.verbs) == 0
}

// IsFinite returns true if all points in the path are finite.
func (p *pathImpl) IsFinite() bool {
	for _, pt := range p.points {
		if !IsFinite(pt.X) || !IsFinite(pt.Y) {
			return false
		}
	}
	return true
}

// IsLine returns true if the path contains only one line.
func (p *pathImpl) IsLine() bool {
	if len(p.verbs) == 2 && p.verbs[1] == enums.PathVerbLine {
		return p.verbs[0] == enums.PathVerbMove
	}
	return false
}

// CountPoints returns the number of points in the path.
func (p *pathImpl) CountPoints() int {
	return len(p.points)
}

// Point returns the point at the specified index.
func (p *pathImpl) Point(index int) Point {
	if index >= 0 && index < len(p.points) {
		return p.points[index]
	}
	return Point{X: 0, Y: 0}
}

// GetPoints copies all points from the path into the provided slice.
func (p *pathImpl) GetPoints(points []Point) int {
	n := len(points)
	if n > len(p.points) {
		n = len(p.points)
	}
	copy(points, p.points[:n])
	return len(p.points)
}

// CountVerbs returns the number of verbs in the path.
func (p *pathImpl) CountVerbs() int {
	return len(p.verbs)
}

// GetVerbs copies all verbs from the path into the provided slice.
func (p *pathImpl) GetVerbs(verbs []enums.PathVerb) int {
	n := len(verbs)
	if n > len(p.verbs) {
		n = len(p.verbs)
	}
	copy(verbs, p.verbs[:n])
	return len(p.verbs)
}

// ConicWeights returns a read-only view of the path's conic weights.
// Returns a copy of the conic weights slice.
func (p *pathImpl) ConicWeights() []Scalar {
	if len(p.conicWeights) == 0 {
		return nil
	}
	weights := make([]Scalar, len(p.conicWeights))
	copy(weights, p.conicWeights)
	return weights
}

// GetLastPoint returns the last point in the path.
// Returns the point and true if the path contains one or more points,
// otherwise returns a zero point and false.
func (p *pathImpl) GetLastPoint() (Point, bool) {
	if len(p.points) == 0 {
		return Point{}, false
	}
	return p.points[len(p.points)-1], true
}

// Bounds returns the bounding box of the path.
func (p *pathImpl) Bounds() Rect {
	if p.boundsDirty {
		p.updateBounds()
	}
	return p.bounds
}

// UpdateBoundsCache updates the cached bounds of the path.
func (p *pathImpl) UpdateBoundsCache() {
	p.Bounds() // This will update the cache
}

// ComputeTightBounds returns a tight bounding box of the path.
func (p *pathImpl) ComputeTightBounds() Rect {
	// If we're only lines, then our (quick) bounds is also tight.
	if p.getSegmentMasks() == base.SegmentMaskLine {
		return p.Bounds()
	}
	return p.computeTightBounds()
}

// MoveTo starts a new contour at the specified point.
func (p *pathImpl) MoveTo(x, y Scalar) {
	if len(p.verbs) > 0 && p.verbs[len(p.verbs)-1] == enums.PathVerbMove {
		// Replace the last move point
		p.points[len(p.points)-1] = Point{X: x, Y: y}
	} else {
		// Remember our index
		p.lastMoveToIndex = len(p.points)
		p.verbs = append(p.verbs, enums.PathVerbMove)
		p.points = append(p.points, Point{X: x, Y: y})
	}
	p.dirtyAfterEdit()
}

// MoveToPoint starts a new contour at the specified point.
func (p *pathImpl) MoveToPoint(pt Point) {
	p.MoveTo(pt.X, pt.Y)
}

// LineTo adds a line from the last point to the specified point.
func (p *pathImpl) LineTo(x, y Scalar) {
	p.injectMoveToIfNeeded()
	p.verbs = append(p.verbs, enums.PathVerbLine)
	p.points = append(p.points, Point{X: x, Y: y})
	p.dirtyAfterEdit()
}

// LineToPoint adds a line from the last point to the specified point.
func (p *pathImpl) LineToPoint(pt Point) {
	p.LineTo(pt.X, pt.Y)
}

// QuadTo adds a quadratic bezier from the last point to the specified point.
func (p *pathImpl) QuadTo(cx, cy, x, y Scalar) {
	p.injectMoveToIfNeeded()
	p.verbs = append(p.verbs, enums.PathVerbQuad)
	p.points = append(p.points, Point{X: cx, Y: cy}, Point{X: x, Y: y})
	p.dirtyAfterEdit()
}

// QuadToPoint adds a quadratic bezier from the last point to the specified point.
func (p *pathImpl) QuadToPoint(c, pt Point) {
	p.QuadTo(c.X, c.Y, pt.X, pt.Y)
}

// ConicTo adds a conic bezier from the last point to the specified point.
func (p *pathImpl) ConicTo(cx, cy, x, y Scalar, w Scalar) {
	// check for <= 0 or NaN with this test
	if !(w > 0) {
		p.LineTo(x, y)
	} else if !IsFinite(w) {
		p.LineTo(cx, cy)
		p.LineTo(x, y)
	} else if w == 1.0 {
		p.QuadTo(cx, cy, x, y)
	} else {
		p.injectMoveToIfNeeded()
		p.verbs = append(p.verbs, enums.PathVerbConic)
		p.points = append(p.points, Point{X: cx, Y: cy}, Point{X: x, Y: y})
		p.conicWeights = append(p.conicWeights, w)
		p.dirtyAfterEdit()
	}
}

// ConicToPoint adds a conic bezier from the last point to the specified point.
func (p *pathImpl) ConicToPoint(c, pt Point, w Scalar) {
	p.ConicTo(c.X, c.Y, pt.X, pt.Y, w)
}

// CubicTo adds a cubic bezier from the last point to the specified point.
func (p *pathImpl) CubicTo(cx1, cy1, cx2, cy2, x, y Scalar) {
	p.injectMoveToIfNeeded()
	p.verbs = append(p.verbs, enums.PathVerbCubic)
	p.points = append(p.points, Point{X: cx1, Y: cy1}, Point{X: cx2, Y: cy2}, Point{X: x, Y: y})
	p.dirtyAfterEdit()
}

// CubicToPoint adds a cubic bezier from the last point to the specified point.
func (p *pathImpl) CubicToPoint(c1, c2, pt Point) {
	p.CubicTo(c1.X, c1.Y, c2.X, c2.Y, pt.X, pt.Y)
}

// Close closes the current contour.
func (p *pathImpl) Close() {
	if len(p.verbs) > 0 {
		switch p.verbs[len(p.verbs)-1] {
		case enums.PathVerbLine, enums.PathVerbQuad, enums.PathVerbConic, enums.PathVerbCubic, enums.PathVerbMove:
			p.verbs = append(p.verbs, enums.PathVerbClose)
		case enums.PathVerbClose:
			// don't add a close if it's a repeat
		}
	}
	// signal that we need a moveTo to follow us (unless we're done)
	if p.lastMoveToIndex >= 0 {
		p.lastMoveToIndex = ^p.lastMoveToIndex
	}
}

// AddRect adds a rectangle to the path.
func (p *pathImpl) AddRect(rect Rect, dir enums.PathDirection, startIndex uint) {
	p.addRaw(RectPathRaw(rect, dir, startIndex))
}

// AddOval adds an oval to the path.
func (p *pathImpl) AddOval(rect Rect, dir enums.PathDirection) {
	// legacy start index: 1
	p.addRaw(OvalPathRaw(rect, dir, 1))
}

// AddCircle adds a circle to the path.
func (p *pathImpl) AddCircle(cx, cy, radius Scalar, dir enums.PathDirection) {
	if radius > 0 {
		p.AddOval(Rect{
			Left:   cx - radius,
			Top:    cy - radius,
			Right:  cx + radius,
			Bottom: cy + radius,
		}, dir)
	}
}

// AddRRect adds a rounded rectangle to the path.
func (p *pathImpl) AddRRect(rrect RRect, dir enums.PathDirection) {
	// legacy start indices: 6 (CW) and 7 (CCW)
	startIndex := uint(6)
	if dir == enums.PathDirectionCCW {
		startIndex = 7
	}
	p.addRRectWithStart(rrect, dir, startIndex)
}

// AddRRect adds a rounded rectangle to the path with a specific start index.
func (p *pathImpl) addRRectWithStart(rrect RRect, dir enums.PathDirection, startIndex uint) {
	if rrect.IsRect() || rrect.IsEmpty() {
		// degenerate(rect) => radii points are collapsing
		bounds := rrect.Bounds()
		p.AddRect(bounds, dir, (startIndex+1)/2)
	} else if rrect.IsOval() {
		// degenerate(oval) => line points are collapsing
		bounds := rrect.Bounds()
		p.addRaw(OvalPathRaw(bounds, dir, startIndex/2))
	} else {
		p.addRaw(RRectPathRaw(rrect, dir, startIndex))
	}
}

// AddPath adds another path to this path with offset.
func (p *pathImpl) AddPath(path interfaces.SkPath, dx, dy Scalar, addMode enums.AddPathMode) {
	// Create a translation matrix for the offset
	matrix := NewMatrixTranslate(dx, dy)
	p.addPathWithMatrix(path, matrix, addMode)
}

// AddPathNoOffset adds another path to this path without offset.
func (p *pathImpl) AddPathNoOffset(path SkPath, addMode enums.AddPathMode) {
	// Use identity matrix (no transformation)
	matrix := NewMatrixIdentity()
	p.addPathWithMatrix(path, matrix, addMode)
}

// AddPathMatrix adds another path to this path with matrix transformation.
func (p *pathImpl) AddPathMatrix(path SkPath, matrix SkMatrix, addMode enums.AddPathMode) {
	p.addPathWithMatrix(path, matrix, addMode)
}

// addPathWithMatrix adds another path to this path with a matrix transformation.
func (p *pathImpl) addPathWithMatrix(srcPath SkPath, matrix SkMatrix, mode enums.AddPathMode) {

	// If source path is empty, nothing to do
	if srcPath == nil || srcPath.IsEmpty() {
		return
	}

	// Check if we can replace this path entirely
	canReplaceThis := (mode == enums.AddPathModeAppend && p.isEffectivelyEmpty()) || p.IsEmpty()
	if canReplaceThis && matrix.IsIdentity() {
		// Save fill type and copy the source path
		fillType := p.fillType
		if srcImpl, ok := srcPath.(*pathImpl); ok {
			// Copy the path data
			p.points = make([]Point, len(srcImpl.points))
			copy(p.points, srcImpl.points)
			p.verbs = make([]enums.PathVerb, len(srcImpl.verbs))
			copy(p.verbs, srcImpl.verbs)
			p.conicWeights = make([]Scalar, len(srcImpl.conicWeights))
			copy(p.conicWeights, srcImpl.conicWeights)
			p.lastMoveToIndex = srcImpl.lastMoveToIndex
			p.fillType = fillType // Restore original fill type
			p.dirtyAfterEdit()
		}
		return
	}

	// Handle self-addition: if we're adding ourselves, we need to copy first
	var src *pathImpl
	var tmpPath *pathImpl
	if srcImpl, ok := srcPath.(*pathImpl); ok {
		if p == srcImpl {
			// Copy the path to avoid modifying while iterating
			tmpPath = &pathImpl{
				points:          make([]Point, len(srcImpl.points)),
				verbs:           make([]enums.PathVerb, len(srcImpl.verbs)),
				conicWeights:    make([]Scalar, len(srcImpl.conicWeights)),
				fillType:        srcImpl.fillType,
				lastMoveToIndex: srcImpl.lastMoveToIndex,
			}
			copy(tmpPath.points, srcImpl.points)
			copy(tmpPath.verbs, srcImpl.verbs)
			copy(tmpPath.conicWeights, srcImpl.conicWeights)
			src = tmpPath
		} else {
			src = srcImpl
		}
	} else {
		// If it's not a pathImpl, we'll need to iterate through it
		// For now, we'll handle it in the iteration loop below
		src = nil
	}

	// Optimized path for Append mode with simple translation (no perspective)
	if mode == enums.AddPathModeAppend && !matrix.HasPerspective() && src != nil {
		// Update lastMoveToIndex
		if src.lastMoveToIndex >= 0 {
			p.lastMoveToIndex = src.lastMoveToIndex + p.CountPoints()
		} else if src.lastMoveToIndex != initialLastMoveToIndexValue {
			p.lastMoveToIndex = src.lastMoveToIndex - p.CountPoints()
		}

		// Reserve space
		p.incReserve(len(src.points), len(src.verbs), len(src.conicWeights))

		// Copy and transform points
		for i := range src.points {
			p.points = append(p.points, matrix.MapPoint(src.points[i]))
		}

		// Copy verbs
		p.verbs = append(p.verbs, src.verbs...)

		// Copy conic weights
		p.conicWeights = append(p.conicWeights, src.conicWeights...)

		p.dirtyAfterEdit()
		return
	}

	// General case: iterate through the source path and add each verb
	// We need to iterate through verbs and track points correctly
	if src == nil {
		// For non-pathImpl, get all verbs and points upfront
		verbCount := srcPath.CountVerbs()
		verbs := make([]enums.PathVerb, verbCount)
		srcPath.GetVerbs(verbs)
		pointCount := srcPath.CountPoints()
		points := make([]Point, pointCount)
		srcPath.GetPoints(points)

		firstVerb := true
		pointIdx := 0

		for _, verb := range verbs {
			switch verb {
			case enums.PathVerbMove:
				mappedPt := matrix.MapPoint(points[pointIdx])
				pointIdx++

				if firstVerb && mode == enums.AddPathModeExtend && !p.IsEmpty() {
					p.injectMoveToIfNeeded() // In case last contour is closed
					lastPt := p.getLastPt()
					// Don't add lineTo if it's degenerate (same point)
					if !lastPt.hasValue || lastPt.value != mappedPt {
						p.LineToPoint(mappedPt)
					}
				} else {
					p.MoveToPoint(mappedPt)
				}
				firstVerb = false

			case enums.PathVerbLine:
				mappedPt := matrix.MapPoint(points[pointIdx])
				pointIdx++
				p.LineToPoint(mappedPt)
				firstVerb = false

			case enums.PathVerbQuad:
				mappedCtrl := matrix.MapPoint(points[pointIdx])
				mappedEnd := matrix.MapPoint(points[pointIdx+1])
				pointIdx += 2
				p.QuadToPoint(mappedCtrl, mappedEnd)
				firstVerb = false

			case enums.PathVerbConic:
				mappedCtrl := matrix.MapPoint(points[pointIdx])
				mappedEnd := matrix.MapPoint(points[pointIdx+1])
				pointIdx += 2
				// For non-pathImpl, we can't get conic weights, so use default
				weight := Scalar(1.0)
				p.ConicToPoint(mappedCtrl, mappedEnd, weight)
				firstVerb = false

			case enums.PathVerbCubic:
				mappedCtrl1 := matrix.MapPoint(points[pointIdx])
				mappedCtrl2 := matrix.MapPoint(points[pointIdx+1])
				mappedEnd := matrix.MapPoint(points[pointIdx+2])
				pointIdx += 3
				p.CubicToPoint(mappedCtrl1, mappedCtrl2, mappedEnd)
				firstVerb = false

			case enums.PathVerbClose:
				p.Close()
				firstVerb = false
			}
		}
	} else {
		// For pathImpl, iterate through verbs and track points
		firstVerb := true
		pointIdx := 0
		conicWeightIdx := 0

		for _, verb := range src.verbs {
			switch verb {
			case enums.PathVerbMove:
				mappedPt := matrix.MapPoint(src.points[pointIdx])
				pointIdx++

				if firstVerb && mode == enums.AddPathModeExtend && !p.IsEmpty() {
					p.injectMoveToIfNeeded() // In case last contour is closed
					lastPt := p.getLastPt()
					// Don't add lineTo if it's degenerate (same point)
					if !lastPt.hasValue || lastPt.value != mappedPt {
						p.LineToPoint(mappedPt)
					}
				} else {
					p.MoveToPoint(mappedPt)
				}
				firstVerb = false

			case enums.PathVerbLine:
				mappedPt := matrix.MapPoint(src.points[pointIdx])
				pointIdx++
				p.LineToPoint(mappedPt)
				firstVerb = false

			case enums.PathVerbQuad:
				mappedCtrl := matrix.MapPoint(src.points[pointIdx])
				mappedEnd := matrix.MapPoint(src.points[pointIdx+1])
				pointIdx += 2
				p.QuadToPoint(mappedCtrl, mappedEnd)
				firstVerb = false

			case enums.PathVerbConic:
				mappedCtrl := matrix.MapPoint(src.points[pointIdx])
				mappedEnd := matrix.MapPoint(src.points[pointIdx+1])
				weight := src.conicWeights[conicWeightIdx]
				pointIdx += 2
				conicWeightIdx++
				p.ConicToPoint(mappedCtrl, mappedEnd, weight)
				firstVerb = false

			case enums.PathVerbCubic:
				mappedCtrl1 := matrix.MapPoint(src.points[pointIdx])
				mappedCtrl2 := matrix.MapPoint(src.points[pointIdx+1])
				mappedEnd := matrix.MapPoint(src.points[pointIdx+2])
				pointIdx += 3
				p.CubicToPoint(mappedCtrl1, mappedCtrl2, mappedEnd)
				firstVerb = false

			case enums.PathVerbClose:
				p.Close()
				firstVerb = false
			}
		}
	}
}

// Transform applies a matrix transformation to the path.
func (p *pathImpl) Transform(matrix SkMatrix) {
	// Transform all points
	for i := range p.points {
		p.points[i] = matrix.MapPoint(p.points[i])
	}
	p.boundsDirty = true
	p.setConvexity(enums.PathConvexityUnknown)
}

// Offset translates the path by the specified offset.
func (p *pathImpl) Offset(dx, dy Scalar) {
	for i := range p.points {
		p.points[i].X += dx
		p.points[i].Y += dy
	}
	p.boundsDirty = true
}

func (p *pathImpl) trimTrailingMoves() ([]Point, []enums.PathVerb) {
	points := p.points
	verbs := p.verbs

	// Trim trailing moves
	for len(verbs) > 0 && verbs[len(verbs)-1] == enums.PathVerbMove {
		verbs = verbs[:len(verbs)-1]
		if len(points) > 0 {
			points = points[:len(points)-1]
		}
	}

	return points, verbs
}

// Private helper methods

func (p *pathImpl) getConvexityOrUnknown() enums.PathConvexity {
	return p.convexity
}

func (p *pathImpl) setConvexity(c enums.PathConvexity) {
	p.convexity = c
}

func (p *pathImpl) computeConvexity() enums.PathConvexity {
	if !p.IsFinite() {
		return enums.PathConvexityConcave
	}

	// Trim trailing moves
	points, verbs := p.trimTrailingMoves()

	if len(verbs) == 0 {
		return enums.PathConvexityConvexDegenerate
	}

	// Quick concave test: check if path changes direction more than three times
	if isConcaveBySign(points) {
		return enums.PathConvexityConcave
	}

	// Iterate through the path and check convexity
	contourCount := 0
	needsClose := false
	state := newConvexicator()

	pointIdx := 0
	conicWeightIdx := 0

	for _, verb := range verbs {
		// Looking for the last moveTo before non-move verbs start
		if contourCount == 0 {
			if verb == enums.PathVerbMove {
				if pointIdx < len(points) {
					state.setMovePt(points[pointIdx])
					pointIdx++
				}
			} else {
				// Starting the actual contour, fall through to add the points
				// Note: This assumes there was a MoveTo (which should always be the case)
				contourCount++
				needsClose = true
			}
		}

		// Accumulating points into the Convexicator until we hit a close or another move
		if contourCount == 1 {
			if verb == enums.PathVerbClose || verb == enums.PathVerbMove {
				if !state.close() {
					return enums.PathConvexityConcave
				}
				needsClose = false
				contourCount++
				if verb == enums.PathVerbMove {
					if pointIdx < len(points) {
						state.setMovePt(points[pointIdx])
						pointIdx++
					}
				}
			} else {
				// Lines add 1 point, cubics add 3, conics and quads add 2
				// These are the points AFTER the start point (which is tracked in state.lastPt)
				count := ptsInVerb(verb)
				if count > 0 && pointIdx+count-1 < len(points) {
					for i := 0; i < count; i++ {
						if !state.addPt(points[pointIdx+i]) {
							return enums.PathConvexityConcave
						}
					}
					pointIdx += count
					if verb == enums.PathVerbConic {
						conicWeightIdx++
					}
				}
			}
		} else {
			// The first contour has closed and anything other than spurious trailing moves means
			// there's multiple contours and the path can't be convex
			if verb != enums.PathVerbMove {
				return enums.PathConvexityConcave
			}
			if pointIdx < len(points) {
				pointIdx++
			}
		}
	}

	// If the path isn't explicitly closed, do so implicitly
	if needsClose && !state.close() {
		return enums.PathConvexityConcave
	}

	firstDir := state.getFirstDirection()
	if firstDir == enums.PathFirstDirectionUnknown && state.reversals >= 3 {
		return enums.PathConvexityConcave
	}

	return pathFirstDirectionToConvexity(firstDir)
}

func (p *pathImpl) getSegmentMasks() uint32 {
	var mask uint32
	for _, verb := range p.verbs {
		switch verb {
		case enums.PathVerbLine:
			mask |= base.SegmentMaskLine
		case enums.PathVerbQuad:
			mask |= base.SegmentMaskQuad
		case enums.PathVerbConic:
			mask |= base.SegmentMaskConic
		case enums.PathVerbCubic:
			mask |= base.SegmentMaskCubic
		}
	}
	return mask
}

func (p *pathImpl) injectMoveToIfNeeded() {
	if p.lastMoveToIndex < 0 {
		var x, y Scalar
		if len(p.verbs) == 0 {
			x, y = 0, 0
		} else {
			pt := p.points[^p.lastMoveToIndex]
			x, y = pt.X, pt.Y
		}
		p.MoveTo(x, y)
	}
}

func (p *pathImpl) dirtyAfterEdit() {
	p.boundsDirty = true
	p.setConvexity(enums.PathConvexityUnknown)
}

// isEffectivelyEmpty returns true if the path has at most one verb (effectively empty)
func (p *pathImpl) isEffectivelyEmpty() bool {
	return len(p.verbs) <= 1
}

// getLastPt returns the last point in the path, if any
func (p *pathImpl) getLastPt() lastPointResult {
	if len(p.points) > 0 {
		return lastPointResult{
			hasValue: true,
			value:    p.points[len(p.points)-1],
		}
	}
	return lastPointResult{hasValue: false}
}

func (p *pathImpl) updateBounds() {
	if len(p.points) == 0 {
		p.bounds = Rect{Left: 0, Top: 0, Right: 0, Bottom: 0}
		p.boundsDirty = false
		return
	}

	left := p.points[0].X
	top := p.points[0].Y
	right := p.points[0].X
	bottom := p.points[0].Y

	for _, pt := range p.points[1:] {
		if pt.X < left {
			left = pt.X
		}
		if pt.X > right {
			right = pt.X
		}
		if pt.Y < top {
			top = pt.Y
		}
		if pt.Y > bottom {
			bottom = pt.Y
		}
	}

	p.bounds = Rect{Left: left, Top: top, Right: right, Bottom: bottom}
	p.boundsDirty = false
}

func (p *pathImpl) computeTightBounds() Rect {
	if len(p.verbs) == 0 {
		return Rect{Left: 0, Top: 0, Right: 0, Bottom: 0}
	}

	// Initialize with the first MoveTo point
	if len(p.points) == 0 {
		return Rect{Left: 0, Top: 0, Right: 0, Bottom: 0}
	}

	left := p.points[0].X
	top := p.points[0].Y
	right := p.points[0].X
	bottom := p.points[0].Y

	// Iterate through verbs and compute extrema for curves
	pointIdx := 0
	conicWeightIdx := 0
	lastPointIdx := -1 // Track the last point index for curve start points

	for _, verb := range p.verbs {
		var extremas []Point
		var count int

		switch verb {
		case enums.PathVerbMove:
			if pointIdx < len(p.points) {
				extremas = []Point{p.points[pointIdx]}
				count = 1
				lastPointIdx = pointIdx
				pointIdx++
			}

		case enums.PathVerbLine:
			if pointIdx < len(p.points) {
				extremas = []Point{p.points[pointIdx]}
				count = 1
				lastPointIdx = pointIdx
				pointIdx++
			}

		case enums.PathVerbQuad:
			if pointIdx+1 < len(p.points) && lastPointIdx >= 0 {
				// Quad needs: start point (lastPointIdx), control (pointIdx), end (pointIdx+1)
				quadPts := []Point{p.points[lastPointIdx], p.points[pointIdx], p.points[pointIdx+1]}
				extremas, count = computeQuadExtremas(quadPts)
				lastPointIdx = pointIdx + 1
				pointIdx += 2
			}

		case enums.PathVerbConic:
			if pointIdx+1 < len(p.points) && conicWeightIdx < len(p.conicWeights) && lastPointIdx >= 0 {
				// Conic needs: start point (lastPointIdx), control (pointIdx), end (pointIdx+1)
				conicPts := []Point{p.points[lastPointIdx], p.points[pointIdx], p.points[pointIdx+1]}
				extremas, count = computeConicExtremas(conicPts, p.conicWeights[conicWeightIdx])
				lastPointIdx = pointIdx + 1
				pointIdx += 2
				conicWeightIdx++
			}

		case enums.PathVerbCubic:
			if pointIdx+2 < len(p.points) && lastPointIdx >= 0 {
				// Cubic needs: start point (lastPointIdx), control1 (pointIdx), control2 (pointIdx+1), end (pointIdx+2)
				cubicPts := []Point{p.points[lastPointIdx], p.points[pointIdx], p.points[pointIdx+1], p.points[pointIdx+2]}
				extremas, count = computeCubicExtremas(cubicPts)
				lastPointIdx = pointIdx + 2
				pointIdx += 3
			}

		case enums.PathVerbClose:
			// No points to add for close, but lastPointIdx might need to be reset
			// Actually, we keep it for potential next curve
		}

		// Update bounds with extrema points
		for i := 0; i < count; i++ {
			pt := extremas[i]
			if pt.X < left {
				left = pt.X
			}
			if pt.X > right {
				right = pt.X
			}
			if pt.Y < top {
				top = pt.Y
			}
			if pt.Y > bottom {
				bottom = pt.Y
			}
		}
	}

	return Rect{Left: left, Top: top, Right: right, Bottom: bottom}
}

func (p *pathImpl) addRaw(raw PathRaw) {
	// Reserve space
	// Count conics in raw path to reserve space for weights
	conicCount := 0
	for _, verb := range raw.Verbs {
		if verb == enums.PathVerbConic {
			conicCount++
		}
	}
	p.incReserve(len(raw.Points), len(raw.Verbs), conicCount)

	// Iterate through raw path and add elements
	for i, verb := range raw.Verbs {
		switch verb {
		case enums.PathVerbMove:
			p.MoveToPoint(raw.Points[raw.PointIndices[i]])
		case enums.PathVerbLine:
			p.LineToPoint(raw.Points[raw.PointIndices[i]+1])
		case enums.PathVerbQuad:
			p.QuadToPoint(raw.Points[raw.PointIndices[i]], raw.Points[raw.PointIndices[i]+1])
		case enums.PathVerbConic:
			weight := raw.ConicWeights[raw.ConicIndex[i]]
			p.ConicToPoint(raw.Points[raw.PointIndices[i]], raw.Points[raw.PointIndices[i]+1], weight)
		case enums.PathVerbCubic:
			p.CubicToPoint(raw.Points[raw.PointIndices[i]], raw.Points[raw.PointIndices[i]+1], raw.Points[raw.PointIndices[i]+2])
		case enums.PathVerbClose:
			p.Close()
		}
	}
}

func (p *pathImpl) incReserve(extraPtCount, extraVerbCount, extraConicCount int) {
	// Reserve capacity for points, verbs, and conic weights
	if cap(p.points) < len(p.points)+extraPtCount {
		newPoints := make([]Point, len(p.points), len(p.points)+extraPtCount)
		copy(newPoints, p.points)
		p.points = newPoints
	}
	if cap(p.verbs) < len(p.verbs)+extraVerbCount {
		newVerbs := make([]enums.PathVerb, len(p.verbs), len(p.verbs)+extraVerbCount)
		copy(newVerbs, p.verbs)
		p.verbs = newVerbs
	}
	if cap(p.conicWeights) < len(p.conicWeights)+extraConicCount {
		newConicWeights := make([]Scalar, len(p.conicWeights), len(p.conicWeights)+extraConicCount)
		copy(newConicWeights, p.conicWeights)
		p.conicWeights = newConicWeights
	}
}
