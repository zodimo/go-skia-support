package impl

import (
	"unsafe"

	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// RasterImage is an implementation of SkImage backed by a raster bitmap (CPU memory).
type RasterImage struct {
	*Image
	Pixels   []byte
	RowBytes int
}

// Ensure RasterImage implements interfaces.SkImage
var _ interfaces.SkImage = &RasterImage{}

// NewRasterImage creates a new RasterImage.
// It copies the pixel data to ensure immutability semantics of SkImage.
func NewRasterImage(info models.ImageInfo, pixels []byte, rowBytes int) *RasterImage {
	// Create a copy of the pixels to enforce immutability
	// Calculate total size needed
	height := info.Height()
	if height < 1 {
		height = 1
	}
	size := height * rowBytes

	// Safety check against provided slice length
	if len(pixels) < size {
		if len(pixels) < size {
			size = len(pixels)
		}
	}

	ownedPixels := make([]byte, size)
	copy(ownedPixels, pixels)

	return &RasterImage{
		Image:    NewImage(info),
		Pixels:   ownedPixels,
		RowBytes: rowBytes,
	}
}

// PeekPixels copies address, row bytes, and info to pixmap if pixels are accessible.
// For RasterImage, this is always true.
func (r *RasterImage) PeekPixels(pixmap *models.Pixmap) bool {
	if pixmap == nil {
		return false
	}
	if len(r.Pixels) == 0 {
		return false
	}
	pixmap.Info = r.Info
	pixmap.RowBytes = r.RowBytes
	pixmap.Addr = unsafe.Pointer(&r.Pixels[0])
	return true
}

// ReadPixels copies rect of pixels to dst.
func (r *RasterImage) ReadPixels(dstInfo models.ImageInfo, dstPixels []byte, dstRowBytes int, srcX int, srcY int) bool {
	// Basic validation
	if len(dstPixels) == 0 {
		return false
	}

	// Calculate bytes per pixel
	bpp := r.Info.BytesPerPixel()
	if bpp <= 0 {
		return false // Unknown or invalid color type
	}

	width := dstInfo.Width()
	height := dstInfo.Height()

	// Copy row by row
	for y := 0; y < height; y++ {
		srcRowY := srcY + y
		if srcRowY >= r.Info.Height() {
			break
		}

		dstRowOffset := y * dstRowBytes
		srcRowOffset := srcRowY * r.RowBytes

		srcPixelOffset := srcRowOffset + srcX*bpp

		// Bytes to copy for this row
		bytesToCopy := width * bpp

		// Bounds check on src
		if srcPixelOffset+bytesToCopy > len(r.Pixels) {
			bytesToCopy = len(r.Pixels) - srcPixelOffset
		}
		if bytesToCopy <= 0 {
			continue
		}

		// Bounds check on dst
		if dstRowOffset+bytesToCopy > len(dstPixels) {
			bytesToCopy = len(dstPixels) - dstRowOffset
		}
		if bytesToCopy <= 0 {
			continue
		}

		copy(dstPixels[dstRowOffset:dstRowOffset+bytesToCopy], r.Pixels[srcPixelOffset:srcPixelOffset+bytesToCopy])
	}

	return true
}
