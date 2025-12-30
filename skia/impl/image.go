package impl

import (
	"sync/atomic"

	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// Global unique ID counter for images
var nextImageID uint32 = 0

// Image is the base implementation of SkImage.
// It handles the common properties like ImageInfo and UniqueID.
type Image struct {
	Info     models.ImageInfo
	uniqueID uint32 // Renamed to lower case to avoid conflict with method
}

// NewImage creates a new base Image with a generated unique ID.
func NewImage(info models.ImageInfo) *Image {
	return &Image{
		Info:     info,
		uniqueID: atomic.AddUint32(&nextImageID, 1),
	}
}

// ImageInfo returns the ImageInfo describing the width, height, color type, alpha type, etc.
func (i *Image) ImageInfo() models.ImageInfo {
	return i.Info
}

// Width returns the width of the image in pixels.
func (i *Image) Width() int {
	return i.Info.Width()
}

// Height returns the height of the image in pixels.
func (i *Image) Height() int {
	return i.Info.Height()
}

// Dimensions returns the width and height as an ISize.
func (i *Image) Dimensions() models.ISize {
	return models.NewISize(i.Info.Width(), i.Info.Height())
}

// Bounds returns the bounds as an IRect (0, 0, width, height).
func (i *Image) Bounds() models.IRect {
	return models.NewIRect(0, 0, i.Info.Width(), i.Info.Height())
}

// UniqueID returns a value unique to the image contents.
func (i *Image) UniqueID() uint32 {
	return i.uniqueID
}

// AlphaType returns the AlphaType of the image.
func (i *Image) AlphaType() enums.AlphaType {
	return i.Info.AlphaType()
}

// ColorType returns the ColorType of the image.
func (i *Image) ColorType() enums.ColorType {
	return i.Info.ColorType()
}

// IsAlphaOnly returns true if the image pixels represent transparency only.
func (i *Image) IsAlphaOnly() bool {
	// Inline check for alpha only color types
	ct := i.Info.ColorType()
	return ct == enums.ColorTypeAlpha8
}

// IsOpaque returns true if pixels ignore their alpha value and are treated as fully opaque.
func (i *Image) IsOpaque() bool {
	return i.Info.IsOpaque()
}

// MakeShader creates a shader with the specified tiling and sampling.
func (i *Image) MakeShader(tmx, tmy enums.TileMode, sampling models.SamplingOptions, localMatrix *interfaces.SkMatrix) interfaces.Shader {
	// TODO: Implement MakeShader when SkShader support is ready.
	return nil
}

// IsTextureBacked returns true if the image is backed by a GPU texture.
func (i *Image) IsTextureBacked() bool {
	return false
}

// IsValid returns true if the image can be drawn.
func (i *Image) IsValid(context interface{}) bool {
	return true
}

// PeekPixels copies address, row bytes, and info to pixmap if pixels are accessible.
func (i *Image) PeekPixels(pixmap *models.Pixmap) bool {
	return false
}

// ReadPixels copies rect of pixels to dst.
func (i *Image) ReadPixels(dstInfo models.ImageInfo, dstPixels []byte, dstRowBytes int, srcX int, srcY int) bool {
	return false
}
