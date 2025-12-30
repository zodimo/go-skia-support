package models

// ISize holds a 32-bit integer width and height
// Matches C++ SkISize
type ISize struct {
	Width, Height int32
}

func NewISize(width, height int) ISize {
	return ISize{
		Width:  int32(width),
		Height: int32(height),
	}
}
