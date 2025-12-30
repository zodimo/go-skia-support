package impl

import (
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// MakeRasterCopy creates a new RasterImage by copying the pixels from the provided pixmap.
// Returns nil if the pixmap is invalid or pixels cannot be read.
func MakeRasterCopy(pixmap models.Pixmap) interfaces.SkImage {
	if pixmap.Addr == nil {
		return nil
	}

	width := pixmap.Info.Width()
	height := pixmap.Info.Height()
	if width <= 0 || height <= 0 {
		return nil
	}

	// Validate RowBytes
	minRowBytes := pixmap.Info.MinRowBytes()
	if pixmap.RowBytes < minRowBytes {
		return nil
	}

	// Calculate total size to read
	size := height * pixmap.RowBytes

	// Create a slice view of the unsafe memory
	// This does NOT copy, it just creates a view so we can pass it to NewRasterImage (which WILL copy)
	// We use a large capacity slice to avoid bounds checks during copy?
	// Or better, use unsafe.Slice if available (Go 1.17+).
	// Assuming standard Go, we can cast.

	// Note: We can't safely know the capacity of the memory pointed to by Addr.
	// We assume it's at least 'size' bytes.

	// Safer approach: copy row by row using pointer arithmetic?
	// But NewRasterImage expects []byte.

	// Let's create a huge slice view and slice it down?
	// Or just alloc a new slice and copy into it right here, then pass to a constructor that takes ownership?
	// NewRasterImage does a copy. So we would copy twice.

	// Let's refactor NewRasterImage to NOT copy? Or have n internal one?
	// No, keep NewRasterImage safe.

	// Let's construct a temporary byte slice view.
	const maxInt = int(^uint(0) >> 1)

	// This is a common Go pattern for converting pointer to slice
	// but it is unsafe.
	sourceBytes := (*[1 << 30]byte)(pixmap.Addr)[:size:size]

	return NewRasterImage(pixmap.Info, sourceBytes, pixmap.RowBytes)
}

// MakeRasterData creates a new RasterImage from the provided data.
// It copies the data.
func MakeRasterData(info models.ImageInfo, pixels []byte, rowBytes int) interfaces.SkImage {
	if len(pixels) == 0 {
		return nil
	}
	if info.IsEmpty() {
		return nil
	}
	if rowBytes < info.MinRowBytes() {
		return nil
	}
	return NewRasterImage(info, pixels, rowBytes)
}
