package interfaces

import (
	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/models"
)

// SkImage describes a two dimensional array of pixels to draw.
// The pixels may be decoded in a raster bitmap, encoded in a SkPicture or compressed data stream,
// or located in GPU memory as a GPU texture.
//
// SkImage cannot be modified after it is created.
//
// Ported from: skia-source/include/core/SkImage.h
type SkImage interface {
	// ImageInfo returns the ImageInfo describing the width, height, color type, alpha type, etc.
	ImageInfo() models.ImageInfo

	// Width returns the width of the image in pixels.
	Width() int

	// Height returns the height of the image in pixels.
	Height() int

	// Dimensions returns the width and height as an ISize.
	Dimensions() models.ISize

	// Bounds returns the bounds as an IRect (0, 0, width, height).
	Bounds() models.IRect

	// UniqueID returns a value unique to the image contents.
	UniqueID() uint32

	// AlphaType returns the AlphaType of the image.
	AlphaType() enums.AlphaType

	// ColorType returns the ColorType of the image.
	ColorType() enums.ColorType

	// IsAlphaOnly returns true if the image pixels represent transparency only (e.g. Alpha8).
	IsAlphaOnly() bool

	// IsOpaque returns true if pixels ignore their alpha value and are treated as fully opaque.
	IsOpaque() bool

	// MakeShader creates a shader with the specified tiling and sampling.
	// tmx, tmy: Tile modes for x and y axes.
	// sampling: Sampling options.
	// localMatrix: Optional local matrix (can be nil).
	MakeShader(tmx, tmy enums.TileMode, sampling models.SamplingOptions, localMatrix *SkMatrix) Shader

	// IsTextureBacked returns true if the image is backed by a GPU texture.
	IsTextureBacked() bool

	// IsValid returns true if the image can be drawn (e.g. context is valid).
	// Takes a Recorder context if applicable (nil for raster/generic check).
	// In C++ this is isValid(SkRecorder*), here we use a generic interface or nil for now.
	IsValid(context interface{}) bool

	// PeekPixels copies address, row bytes, and info to pixmap if pixels are accessible.
	// Returns true if successful.
	// Requires direct access to pixels (Raster).
	PeekPixels(pixmap *models.Pixmap) bool

	// ReadPixels copies rect of pixels to dst.
	// Returns true if pixels were copied successfully.
	// srcX, srcY: offset in source image.
	ReadPixels(dstInfo models.ImageInfo, dstPixels []byte, dstRowBytes int, srcX int, srcY int) bool
}
