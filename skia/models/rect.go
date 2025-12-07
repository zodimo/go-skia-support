package models

// Rect represents a rectangle
type Rect struct {
	Left, Top, Right, Bottom Scalar
}

// IsSorted returns true if the rectangle is sorted (Left <= Right and Top <= Bottom)
func (r Rect) IsSorted() bool {
	return r.Left <= r.Right && r.Top <= r.Bottom
}

// MakeOutset returns a rectangle outset by (dx, dy).
// If dx is negative, the returned rectangle is narrower.
// If dx is positive, the returned rectangle is wider.
// If dy is negative, the returned rectangle is shorter.
// If dy is positive, the returned rectangle is taller.
func (r Rect) MakeOutset(dx, dy Scalar) Rect {
	return Rect{
		Left:   r.Left - dx,
		Top:    r.Top - dy,
		Right:  r.Right + dx,
		Bottom: r.Bottom + dy,
	}
}
