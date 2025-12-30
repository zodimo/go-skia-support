package models

// ColorSpace describes the range of colors.
// Currently a placeholder to match C++ SkImage flexibility.
// In the future, this will hold detailed gamut and transfer function data.
type ColorSpace struct {
	// TODO: Add SkColorSpace fields (gamma, matrix, etc.)
}

// NewColorSpaceSrgb creates a standard sRGB color space.
func NewColorSpaceSrgb() *ColorSpace {
	return &ColorSpace{}
}
