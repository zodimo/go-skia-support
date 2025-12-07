package impl

import "github.com/zodimo/go-skia-support/skia/enums"

// lastPointResult holds the result of getLastPt
type lastPointResult struct {
	hasValue bool
	value    Point
}

// PathRaw represents raw path data
type PathRaw struct {
	Points       []Point
	Verbs        []enums.PathVerb
	ConicWeights []Scalar
	PointIndices []int
	ConicIndex   []int
}

// ovalPointIterator iterates through oval points (ellipse midpoints) based on direction and start index
type ovalPointIterator struct {
	points     [4]Point
	currentIdx int
	advance    int
}

func newOvalPointIterator(rect Rect, dir enums.PathDirection, startIndex uint) *ovalPointIterator {
	cx := (rect.Left + rect.Right) / 2
	cy := (rect.Top + rect.Bottom) / 2

	// Oval points in order (midpoints of ellipse):
	// [0] = (centerX, top) - top of ellipse
	// [1] = (right, centerY) - right of ellipse
	// [2] = (centerX, bottom) - bottom of ellipse
	// [3] = (left, centerY) - left of ellipse
	points := [4]Point{
		{X: cx, Y: rect.Top},
		{X: rect.Right, Y: cy},
		{X: cx, Y: rect.Bottom},
		{X: rect.Left, Y: cy},
	}

	currentIdx := int(startIndex % 4)
	advance := 1
	if dir == enums.PathDirectionCCW {
		advance = 3 // equivalent to going backwards (N-1 where N=4)
	}

	return &ovalPointIterator{
		points:     points,
		currentIdx: currentIdx,
		advance:    advance,
	}
}

func (o *ovalPointIterator) current() Point {
	return o.points[o.currentIdx]
}

func (o *ovalPointIterator) next() Point {
	o.currentIdx = (o.currentIdx + o.advance) % 4
	return o.current()
}

// rectPointIterator iterates through rectangle corners based on direction and start index
type rectPointIterator struct {
	points     [4]Point
	currentIdx int
	advance    int
}

func newRectPointIterator(rect Rect, dir enums.PathDirection, startIndex uint) *rectPointIterator {
	// Rectangle corners in order:
	// [0] = (Left, Top)
	// [1] = (Right, Top)
	// [2] = (Right, Bottom)
	// [3] = (Left, Bottom)
	points := [4]Point{
		{X: rect.Left, Y: rect.Top},
		{X: rect.Right, Y: rect.Top},
		{X: rect.Right, Y: rect.Bottom},
		{X: rect.Left, Y: rect.Bottom},
	}

	currentIdx := int(startIndex % 4)
	advance := 1
	if dir == enums.PathDirectionCCW {
		advance = 3 // equivalent to going backwards (N-1 where N=4)
	}

	return &rectPointIterator{
		points:     points,
		currentIdx: currentIdx,
		advance:    advance,
	}
}

func (r *rectPointIterator) current() Point {
	return r.points[r.currentIdx]
}

func (r *rectPointIterator) next() Point {
	r.currentIdx = (r.currentIdx + r.advance) % 4
	return r.current()
}

// rrectPointIterator iterates through rounded rectangle points (8 points total)
// Points are arranged around the rounded rectangle perimeter
type rrectPointIterator struct {
	points     [8]Point
	currentIdx int
	advance    int
}

func newRRectPointIterator(rrect RRect, dir enums.PathDirection, startIndex uint) *rrectPointIterator {
	bounds := rrect.Bounds()
	L := bounds.Left
	T := bounds.Top
	R := bounds.Right
	B := bounds.Bottom

	// RRect has 8 points around the perimeter:
	// [0] = (L + UL_radii.X, T) - top-left corner start
	// [1] = (R - UR_radii.X, T) - top-right corner start
	// [2] = (R, T + UR_radii.Y) - top-right corner end
	// [3] = (R, B - LR_radii.Y) - bottom-right corner start
	// [4] = (R - LR_radii.X, B) - bottom-right corner end
	// [5] = (L + LL_radii.X, B) - bottom-left corner start
	// [6] = (L, B - LL_radii.Y) - bottom-left corner end
	// [7] = (L, T + UL_radii.Y) - top-left corner end
	points := [8]Point{
		{X: L + rrect.Radii[0].X, Y: T}, // UL start
		{X: R - rrect.Radii[1].X, Y: T}, // UR start
		{X: R, Y: T + rrect.Radii[1].Y}, // UR end
		{X: R, Y: B - rrect.Radii[2].Y}, // LR start
		{X: R - rrect.Radii[2].X, Y: B}, // LR end
		{X: L + rrect.Radii[3].X, Y: B}, // LL start
		{X: L, Y: B - rrect.Radii[3].Y}, // LL end
		{X: L, Y: T + rrect.Radii[0].Y}, // UL end
	}

	currentIdx := int(startIndex % 8)
	advance := 1
	if dir == enums.PathDirectionCCW {
		advance = 7 // equivalent to going backwards (N-1 where N=8)
	}

	return &rrectPointIterator{
		points:     points,
		currentIdx: currentIdx,
		advance:    advance,
	}
}

func (r *rrectPointIterator) current() Point {
	return r.points[r.currentIdx]
}

func (r *rrectPointIterator) next() Point {
	r.currentIdx = (r.currentIdx + r.advance) % 8
	return r.current()
}

// Convexicator tracks convexity state while iterating through a path
type convexicator struct {
	firstPt        Point
	firstVec       Point // direction leaving firstPt
	lastPt         Point
	lastVec        Point // direction that brought path to lastPt
	expectedDir    enums.DirChange
	firstDirection enums.PathFirstDirection
	reversals      int
	isFinite       bool
}

func newConvexicator() *convexicator {
	return &convexicator{
		expectedDir:    enums.DirChangeInvalid,
		firstDirection: enums.PathFirstDirectionUnknown,
		isFinite:       true,
	}
}

func (c *convexicator) setMovePt(pt Point) {
	c.firstPt = pt
	c.lastPt = pt
	c.expectedDir = enums.DirChangeInvalid
	// Reset vectors to zero to match C++ initialization
	// Ported from: skia-source/src/core/SkPathPriv.cpp:setMovePt() (lines 424-427)
	c.lastVec = Point{X: 0, Y: 0}
	c.firstVec = Point{X: 0, Y: 0}
	c.reversals = 0
	c.isFinite = true
}

func (c *convexicator) addPt(pt Point) bool {
	if c.lastPt == pt {
		return true
	}
	// Should only be true for first non-zero vector after setMovePt was called.
	// It is possible we doubled back at the start so need to check if lastVec is zero or not.
	// Ported from: skia-source/src/core/SkPathPriv.cpp:addPt() (lines 429-443)
	vec := Point{X: pt.X - c.lastPt.X, Y: pt.Y - c.lastPt.Y}
	if c.firstPt == c.lastPt && c.expectedDir == enums.DirChangeInvalid && c.lastVec.X == 0 && c.lastVec.Y == 0 {
		c.lastVec = vec
		c.firstVec = vec
	} else if !c.addVec(vec) {
		return false
	}
	c.lastPt = pt
	return true
}

func (c *convexicator) close() bool {
	// If this was an explicit close, there was already a lineTo to firstPt, so this
	// addPt() is a no-op. Otherwise, the addPt implicitly closes the contour.
	return c.addPt(c.firstPt) && c.addVec(c.firstVec)
}

func (c *convexicator) getFirstDirection() enums.PathFirstDirection {
	return c.firstDirection
}

// func (c *convexicator) reversals() int {
// 	return c.reversals
// }

func (c *convexicator) directionChange(curVec Point) enums.DirChange {
	cross := crossProduct(c.lastVec, curVec)
	if !IsFinite(cross) {
		return enums.DirChangeUnknown
	}
	if cross == 0 {
		dot := dotProduct(c.lastVec, curVec)
		if dot < 0 {
			return enums.DirChangeBackwards
		}
		return enums.DirChangeStraight
	}
	if cross > 0 {
		return enums.DirChangeRight
	}
	return enums.DirChangeLeft
}

func (c *convexicator) addVec(curVec Point) bool {
	dir := c.directionChange(curVec)
	switch dir {
	case enums.DirChangeLeft, enums.DirChangeRight:
		if c.expectedDir == enums.DirChangeInvalid {
			c.expectedDir = dir
			if dir == enums.DirChangeRight {
				c.firstDirection = enums.PathFirstDirectionCW
			} else {
				c.firstDirection = enums.PathFirstDirectionCCW
			}
		} else if dir != c.expectedDir {
			c.firstDirection = enums.PathFirstDirectionUnknown
			return false
		}
		c.lastVec = curVec
	case enums.DirChangeStraight:
		// Continue with same direction
	case enums.DirChangeBackwards:
		// Allow path to reverse direction twice
		c.lastVec = curVec
		c.reversals++
		if c.reversals >= 3 {
			return false
		}
	case enums.DirChangeUnknown:
		c.isFinite = false
		return false
	case enums.DirChangeInvalid:
		// Should not happen
		return false
	}
	return true
}
