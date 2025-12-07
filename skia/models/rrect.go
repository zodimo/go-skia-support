package models

import "github.com/zodimo/go-skia-support/skia/enums"

// RRect represents a rounded rectangle with corner radii
// Radii are stored in order: UpperLeft, UpperRight, LowerRight, LowerLeft
// Each Point represents (X radius, Y radius) for that corner
type RRect struct {
	bounds Rect
	Radii  [4]Point // [0]=UL, [1]=UR, [2]=LR, [3]=LL
}

func (r RRect) IsRect() bool {
	// A rect has all corners square (at least one radius is zero for each corner)
	for i := 0; i < 4; i++ {
		if r.Radii[i].X != 0 && r.Radii[i].Y != 0 {
			// This corner has both radii non-zero, so it's rounded
			return false
		}
	}
	return true
}

func (r RRect) IsEmpty() bool {
	return r.bounds.Left >= r.bounds.Right || r.bounds.Top >= r.bounds.Bottom
}

func (r RRect) IsOval() bool {
	// An oval has all radii equal and each radius >= half the corresponding dimension
	if r.IsEmpty() {
		return false
	}

	// Check if all radii are equal
	firstRad := r.Radii[0]
	for i := 1; i < 4; i++ {
		if r.Radii[i].X != firstRad.X || r.Radii[i].Y != firstRad.Y {
			return false
		}
	}

	// Check if radii are at least half the width/height
	width := r.bounds.Right - r.bounds.Left
	height := r.bounds.Bottom - r.bounds.Top
	halfWidth := width / 2
	halfHeight := height / 2

	return firstRad.X >= halfWidth && firstRad.Y >= halfHeight
}

// RadiiAt returns the radii for a specific corner.
// Returns the x-axis and y-axis radii for the specified corner.
func (r RRect) RadiiAt(corner Corner) Point {
	return r.Radii[int(corner)]
}

// GetAllRadii returns all corner radii as a slice.
// The radii are returned in order: UpperLeft, UpperRight, LowerRight, LowerLeft.
func (r RRect) GetAllRadii() []Point {
	return r.Radii[:]
}

// Type returns the type of the rounded rectangle.
func (r RRect) Type() RRectType {
	if r.IsEmpty() {
		return enums.RRectTypeEmpty
	}
	if r.IsRect() {
		return enums.RRectTypeRect
	}
	if r.IsOval() {
		return enums.RRectTypeOval
	}
	if r.IsSimple() {
		return enums.RRectTypeSimple
	}
	if r.IsNinePatch() {
		return enums.RRectTypeNinePatch
	}
	return enums.RRectTypeComplex
}

// Rect returns the bounding rectangle.
func (r RRect) Rect() Rect {
	return r.bounds
}

// Width returns the width of the rounded rectangle.
func (r RRect) Width() Scalar {
	return r.bounds.Right - r.bounds.Left
}

// Height returns the height of the rounded rectangle.
func (r RRect) Height() Scalar {
	return r.bounds.Bottom - r.bounds.Top
}

// IsSimple returns true if the rounded rectangle is simple.
// A simple RRect has all radii equal but is not an oval.
func (r RRect) IsSimple() bool {
	if r.IsEmpty() || r.IsRect() || r.IsOval() {
		return false
	}

	// Check if all radii are equal
	firstRad := r.Radii[0]
	for i := 1; i < 4; i++ {
		if r.Radii[i].X != firstRad.X || r.Radii[i].Y != firstRad.Y {
			return false
		}
	}

	// If all radii are equal but not an oval, it's simple
	return true
}

// IsNinePatch returns true if the rounded rectangle is a nine-patch.
// A nine-patch has axis-aligned radii: UL.X == LL.X, UL.Y == UR.Y, UR.X == LR.X, LL.Y == LR.Y
func (r RRect) IsNinePatch() bool {
	if r.IsEmpty() || r.IsRect() || r.IsOval() || r.IsSimple() {
		return false
	}

	ul := r.Radii[enums.CornerUpperLeft]
	ur := r.Radii[enums.CornerUpperRight]
	lr := r.Radii[enums.CornerLowerRight]
	ll := r.Radii[enums.CornerLowerLeft]

	return ul.X == ll.X && ul.Y == ur.Y && ur.X == lr.X && ll.Y == lr.Y
}

// IsComplex returns true if the rounded rectangle is complex.
// A complex RRect has arbitrary radii that don't fit other categories.
func (r RRect) IsComplex() bool {
	return !r.IsEmpty() && !r.IsRect() && !r.IsOval() && !r.IsSimple() && !r.IsNinePatch()
}

// SetRect sets the rounded rectangle to a rectangle.
func (r *RRect) SetRect(rect Rect) {
	r.bounds = rect
	// Set all radii to zero
	for i := 0; i < 4; i++ {
		r.Radii[i] = Point{X: 0, Y: 0}
	}
}

// SetOval sets the rounded rectangle to an oval.
func (r *RRect) SetOval(rect Rect) {
	r.bounds = rect
	width := rect.Right - rect.Left
	height := rect.Bottom - rect.Top
	halfWidth := width / 2
	halfHeight := height / 2
	// Set all radii to half width/height
	for i := 0; i < 4; i++ {
		r.Radii[i] = Point{X: halfWidth, Y: halfHeight}
	}
}

// SetRectXY sets the rounded rectangle with uniform radii.
func (r *RRect) SetRectXY(rect Rect, rx, ry Scalar) {
	r.bounds = rect
	// Clamp radii to valid values
	if rx < 0 {
		rx = 0
	}
	if ry < 0 {
		ry = 0
	}
	width := rect.Right - rect.Left
	height := rect.Bottom - rect.Top
	// Scale down if radii are too large
	if rx+rx > width {
		rx = width / 2
	}
	if ry+ry > height {
		ry = height / 2
	}
	// Set all radii to the same value
	for i := 0; i < 4; i++ {
		r.Radii[i] = Point{X: rx, Y: ry}
	}
}

// SetNinePatch sets the rounded rectangle with different radii for each corner.
// Parameters: rx1, ry1 (upper-left), rx2, ry2 (upper-right),
//
//	rx3, ry3 (lower-right), rx4, ry4 (lower-left)
func (r *RRect) SetNinePatch(rect Rect, rx1, ry1, rx2, ry2, rx3, ry3, rx4, ry4 Scalar) {
	r.bounds = rect
	// Clamp radii to valid values
	radii := [8]Scalar{rx1, ry1, rx2, ry2, rx3, ry3, rx4, ry4}
	for i := 0; i < 8; i++ {
		if radii[i] < 0 {
			radii[i] = 0
		}
	}
	width := rect.Right - rect.Left
	height := rect.Bottom - rect.Top
	// Scale down if radii are too large
	if radii[0]+radii[6] > width { // rx1 + rx4 (left + right)
		scale := width / (radii[0] + radii[6])
		radii[0] *= scale
		radii[6] *= scale
	}
	if radii[2]+radii[4] > width { // rx2 + rx3 (left + right)
		scale := width / (radii[2] + radii[4])
		radii[2] *= scale
		radii[4] *= scale
	}
	if radii[1]+radii[5] > height { // ry1 + ry3 (top + bottom)
		scale := height / (radii[1] + radii[5])
		radii[1] *= scale
		radii[5] *= scale
	}
	if radii[3]+radii[7] > height { // ry2 + ry4 (top + bottom)
		scale := height / (radii[3] + radii[7])
		radii[3] *= scale
		radii[7] *= scale
	}
	// Set radii: UL, UR, LR, LL
	r.Radii[0] = Point{X: radii[0], Y: radii[1]} // UL: rx1, ry1
	r.Radii[1] = Point{X: radii[2], Y: radii[3]} // UR: rx2, ry2
	r.Radii[2] = Point{X: radii[4], Y: radii[5]} // LR: rx3, ry3
	r.Radii[3] = Point{X: radii[6], Y: radii[7]} // LL: rx4, ry4
}

// SetRectRadii sets the rounded rectangle with radii for each corner.
func (r *RRect) SetRectRadii(rect Rect, radii [4]Point) {
	r.bounds = rect
	// Clamp and validate radii
	width := rect.Right - rect.Left
	height := rect.Bottom - rect.Top
	for i := 0; i < 4; i++ {
		rx := radii[i].X
		ry := radii[i].Y
		if rx < 0 {
			rx = 0
		}
		if ry < 0 {
			ry = 0
		}
		// Scale down if radius is too large
		if i == 0 || i == 3 { // left corners
			if rx > width/2 {
				rx = width / 2
			}
		} else { // right corners
			if rx > width/2 {
				rx = width / 2
			}
		}
		if i == 0 || i == 1 { // top corners
			if ry > height/2 {
				ry = height / 2
			}
		} else { // bottom corners
			if ry > height/2 {
				ry = height / 2
			}
		}
		r.Radii[i] = Point{X: rx, Y: ry}
	}
}

// Bounds returns the bounding box of the rounded rectangle.
func (r RRect) Bounds() Rect {
	return r.bounds
}
