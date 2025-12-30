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

// MakeSubset returns a new image that is a subset of this image.
// Returns nil if the subset is invalid or outside the image bounds.
func (r *RasterImage) MakeSubset(subset models.IRect) interfaces.SkImage {
	bounds := r.Bounds()
	if !bounds.Contains(subset) {
		return nil
	}
	if subset.Width() <= 0 || subset.Height() <= 0 {
		return nil
	}

	bpp := r.Info.BytesPerPixel()
	if bpp <= 0 {
		return nil
	}

	// Calculate offset
	offset := int(subset.Top)*r.RowBytes + int(subset.Left)*bpp

	// Safety check
	if offset < 0 || offset >= len(r.Pixels) {
		return nil // Should be covered by Contains check but safety first
	}

	// Create new info for the subset
	// We use the same ColorType/AlphaType/ColorSpace, but new dimensions
	newInfo := models.NewImageInfo(int(subset.Width()), int(subset.Height()), r.Info.ColorType(), r.Info.AlphaType())
	if cs := r.Info.ColorSpace(); cs != nil {
		newInfo = newInfo.WithColorSpace(cs)
	}

	// The source slice for NewRasterImage.
	// We need to provide enough data for (Height * RowBytes).
	// Since we are preserving the original RowBytes (stride), this works perfectly.
	// The copy in NewRasterImage will copy (Height * RowBytes) bytes starting from offset.
	// This includes the padding at the end of each row if Width < Stride, which is expected for a window into a larger buffer.
	// This might be inefficient (copying padding) but it's correct.

	// Check if we have enough bytes remaining
	needed := int(subset.Height()) * r.RowBytes
	if offset+needed > len(r.Pixels) {
		// NewRasterImage does its own size check but better to fail early here?
		// NewRasterImage handles safe copying but let's trust our calc.
		// If offset+needed is OOB, slice access will panic.
		// So we MUST check.
		return nil
	}

	return NewRasterImage(newInfo, r.Pixels[offset:], r.RowBytes)
}
