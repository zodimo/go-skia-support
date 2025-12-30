package models

import (
	"github.com/zodimo/go-skia-support/skia/enums"
)

// ImageInfo describes the dimensions and color type of an image.
// Matches C++ SkImageInfo from include/core/SkImageInfo.h
type ImageInfo struct {
	width     int
	height    int
	colorType enums.ColorType
	alphaType enums.AlphaType
	// ColorSpace is an optional pointer to a ColorSpace.
	// In C++, this is a smart pointer (sk_sp<SkColorSpace>).
	colorSpace *ColorSpace
}

// NewImageInfo creates a new ImageInfo with the specified properties.
func NewImageInfo(width, height int, colorType enums.ColorType, alphaType enums.AlphaType) ImageInfo {
	return ImageInfo{
		width:     width,
		height:    height,
		colorType: colorType,
		alphaType: alphaType,
	}
}

// MakeN32Premul creates an ImageInfo with N32 color type and Premul alpha type.
func MakeN32Premul(width, height int) ImageInfo {
	return NewImageInfo(width, height, enums.ColorTypeN32, enums.AlphaTypePremul)
}

func (i ImageInfo) Width() int {
	return i.width
}

func (i ImageInfo) Height() int {
	return i.height
}

func (i ImageInfo) ColorType() enums.ColorType {
	return i.colorType
}

func (i ImageInfo) AlphaType() enums.AlphaType {
	return i.alphaType
}

func (i ImageInfo) IsEmpty() bool {
	return i.width <= 0 || i.height <= 0
}

func (i ImageInfo) IsOpaque() bool {
	return i.alphaType == enums.AlphaTypeOpaque
}

func (i ImageInfo) Dimensions() ISize {
	return NewISize(i.width, i.height)
}

func (i ImageInfo) Bounds() IRect {
	return NewIRect(0, 0, i.width, i.height)
}

// ColorSpace returns the color space associated with this ImageInfo.
func (i ImageInfo) ColorSpace() *ColorSpace {
	return i.colorSpace
}

// WithColorSpace returns a new ImageInfo with the specified ColorSpace.
func (i ImageInfo) WithColorSpace(cs *ColorSpace) ImageInfo {
	i.colorSpace = cs
	return i
}

// BytesPerPixel returns the number of bytes per pixel for the color type.
func (i ImageInfo) BytesPerPixel() int {
	switch i.colorType {
	case enums.ColorTypeUnknown:
		return 0
	case enums.ColorTypeAlpha8, enums.ColorTypeGray8:
		return 1
	case enums.ColorTypeRGB565, enums.ColorTypeARGB4444:
		return 2
	case enums.ColorTypeRGBA8888, enums.ColorTypeBGRA8888, enums.ColorTypeRGB888x,
		enums.ColorTypeRGBA1010102, enums.ColorTypeRGB101010x:
		return 4
	case enums.ColorTypeRGBAF16Norm, enums.ColorTypeRGBAF16:
		return 8
	case enums.ColorTypeRGBAF32:
		return 16
	default:
		return 0
	}
}

// RowBytes returns the minimum bytes per row (width * bytesPerPixel).
func (i ImageInfo) MinRowBytes() int {
	return i.width * i.BytesPerPixel()
}

// ComputeByteSize returns the size in bytes of the pixel buffer.
// Returns 0 if height is 0 or calculation overflows (though overflow check is basic here).
func (i ImageInfo) ComputeByteSize(rowBytes int) int64 {
	if i.height == 0 {
		return 0
	}
	// Basic check: rowBytes must be at least minRowBytes
	if rowBytes < i.MinRowBytes() {
		return 0 // Invalid rowBytes
	}
	return int64(i.height) * int64(rowBytes)
}

// ValidRowBytes checks if the rowBytes is valid for this ImageInfo
func (i ImageInfo) ValidRowBytes(rowBytes int) bool {
	if rowBytes < i.MinRowBytes() {
		return false
	}
	// TODO: alignment checks if strict Skia parity is needed
	return true
}
