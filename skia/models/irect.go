package models

// IRect holds four 32-bit integer coordinates for a rectangle
// Matches C++ SkIRect
type IRect struct {
	Left, Top, Right, Bottom int32
}

func NewIRect(left, top, width, height int) IRect {
	return IRect{
		Left:   int32(left),
		Top:    int32(top),
		Right:  int32(left + width),
		Bottom: int32(top + height),
	}
}

func (r IRect) Width() int32 {
	return r.Right - r.Left
}

func (r IRect) Height() int32 {
	return r.Bottom - r.Top
}

func IsEmpty(r IRect) bool {
	// A rectangle is empty if its width or height are <= 0
	return r.Right <= r.Left || r.Bottom <= r.Top
}

// Contains returns true if the rectangle contains the given rectangle.
func (r IRect) Contains(other IRect) bool {
	return r.Left <= other.Left && r.Top <= other.Top &&
		r.Right >= other.Right && r.Bottom >= other.Bottom
}
