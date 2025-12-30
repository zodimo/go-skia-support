package impl

import (
	"testing"
	"unsafe"

	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/models"
)

func TestRasterImage_Properties(t *testing.T) {
	info := models.NewImageInfo(10, 20, enums.ColorTypeRGBA8888, enums.AlphaTypePremul)
	rowBytes := 10 * 4
	pixels := make([]byte, rowBytes*20)
	img := MakeRasterData(info, pixels, rowBytes)

	if img == nil {
		t.Fatal("Failed to create image")
	}

	if img.Width() != 10 {
		t.Errorf("Width mismatch: got %d, want 10", img.Width())
	}
	if img.Height() != 20 {
		t.Errorf("Height mismatch: got %d, want 20", img.Height())
	}
	if img.AlphaType() != enums.AlphaTypePremul {
		t.Errorf("AlphaType mismatch")
	}
	if img.ColorType() != enums.ColorTypeRGBA8888 {
		t.Errorf("ColorType mismatch")
	}
	if img.UniqueID() == 0 {
		t.Errorf("UniqueID should be non-zero")
	}
}

func TestRasterImage_Immutability(t *testing.T) {
	info := models.NewImageInfo(2, 2, enums.ColorTypeAlpha8, enums.AlphaTypeOpaque) // 1 byte per pixel
	rowBytes := 2
	pixels := []byte{1, 2, 3, 4}

	img := MakeRasterData(info, pixels, rowBytes)

	// Modify source
	pixels[0] = 99

	// Read back via PeekPixels
	var pixmap models.Pixmap
	if !img.PeekPixels(&pixmap) {
		t.Fatal("PeekPixels failed")
	}

	// Access the pixel data from pixmap
	// Cast unsafe.Pointer to pixel slice
	internalPixels := (*[4]byte)(pixmap.Addr)[:]

	if internalPixels[0] == 99 {
		t.Error("Image pixels changed after source modification! Image should be immutable copy.")
	}
	if internalPixels[0] != 1 {
		t.Errorf("Expected pixel 0 to be 1, got %d", internalPixels[0])
	}
}

func TestRasterImage_ReadPixels(t *testing.T) {
	weight, height := 2, 2
	info := models.NewImageInfo(weight, height, enums.ColorTypeRGBA8888, enums.AlphaTypePremul)
	// 4 bytes per pixel.
	// RowBytes = 2 * 4 = 8.
	// 2 rows = 16 bytes.
	// Pixels:
	// R0: [0,0,0,0] [1,1,1,1]
	// R1: [2,2,2,2] [3,3,3,3]
	srcPixels := []byte{
		0, 0, 0, 0, 1, 1, 1, 1,
		2, 2, 2, 2, 3, 3, 3, 3,
	}

	img := MakeRasterData(info, srcPixels, 8)

	// Read into a buffer
	dstPixels := make([]byte, 16)
	dstInfo := info

	if !img.ReadPixels(dstInfo, dstPixels, 8, 0, 0) {
		t.Error("ReadPixels failed")
	}

	for i, b := range dstPixels {
		if b != srcPixels[i] {
			t.Errorf("Byte %d mismatch: got %d, want %d", i, b, srcPixels[i])
		}
	}

	// Test partial read
	// Read 1x1 rect at (1,1) -> should correspond to last pixel [3,3,3,3]
	dstSmallInfo := models.NewImageInfo(1, 1, enums.ColorTypeRGBA8888, enums.AlphaTypePremul)
	dstSmall := make([]byte, 4)

	if !img.ReadPixels(dstSmallInfo, dstSmall, 4, 1, 1) {
		t.Error("ReadPixels partial failed")
	}

	for i := 0; i < 4; i++ {
		if dstSmall[i] != 3 {
			t.Errorf("Small read byte %d mismatch: got %d, want 3", i, dstSmall[i])
		}
	}
}

func TestMakeRasterCopy(t *testing.T) {
	info := models.NewImageInfo(2, 2, enums.ColorTypeAlpha8, enums.AlphaTypeOpaque)
	pixels := []byte{10, 20, 30, 40}

	pixmap := models.NewPixmap(info, unsafe.Pointer(&pixels[0]), 2)

	img := MakeRasterCopy(pixmap)
	if img == nil {
		t.Fatal("MakeRasterCopy returned nil")
	}

	var readMap models.Pixmap
	if !img.PeekPixels(&readMap) {
		t.Fatal("PeekPixels failed")
	}

	readPixels := (*[4]byte)(readMap.Addr)[:]
	for i, v := range pixels {
		if readPixels[i] != v {
			t.Errorf("Pixel %d mismatch: got %d, want %d", i, readPixels[i], v)
		}
	}
}

func TestRasterImage_ColorSpace(t *testing.T) {
	cs := models.NewColorSpaceSrgb()
	info := models.NewImageInfo(10, 10, enums.ColorTypeRGBA8888, enums.AlphaTypePremul).WithColorSpace(cs)
	rowBytes := 10 * 4
	pixels := make([]byte, rowBytes*10)

	img := MakeRasterData(info, pixels, rowBytes)
	if img.ColorSpace() != cs {
		t.Error("ColorSpace mismatch")
	}
}

func TestRasterImage_Subset(t *testing.T) {
	// Create 4x4 image with known pattern
	// Indices:
	// 00 01 02 03
	// 10 11 12 13
	// 20 21 22 23
	// 30 31 32 33
	// Each pixel is 1 byte (Alpha8)
	width, height := 4, 4
	info := models.NewImageInfo(width, height, enums.ColorTypeAlpha8, enums.AlphaTypeOpaque)
	pixels := make([]byte, width*height)
	for i := range pixels {
		pixels[i] = byte(i)
	}
	rowBytes := width // 4

	img := MakeRasterData(info, pixels, rowBytes)

	// Make subset: 2x2 starting at non-zero offset (1, 1)
	// Expected pixels:
	// 11(0x0B) 12(0x0C)
	// 21(0x15) 22(0x16) Is that hex? No decimal.
	// 10 11 12 13 (Row 1 starts at index 4) -> 1*4+1 = 5.
	// Row 1 indices: 4, 5, 6, 7.  Wait.
	// Pixel (1,1) is at y=1, x=1. Offset = 1*4 + 1 = 5. Value = 5.
	// Pixel (2,1) (x=2,y=1) is at y=1, x=2. Offset = 6. Value = 6.
	// Pixel (1,2) is at y=2, x=1. Offset = 2*4 + 1 = 9. Value = 9.
	// Pixel (2,2) is at y=2, x=2. Offset = 10. Value = 10.

	subsetRect := models.NewIRect(1, 1, 2, 2)
	subset := img.MakeSubset(subsetRect)

	if subset == nil {
		t.Fatal("MakeSubset returned nil")
	}
	if subset.Width() != 2 || subset.Height() != 2 {
		t.Errorf("Subset dimensions wrong: %d x %d", subset.Width(), subset.Height())
	}

	// Verify content
	var pm models.Pixmap
	if !subset.PeekPixels(&pm) {
		t.Fatal("PeekPixels failed on subset")
	}

	// The subset image should have rowBytes = original rowBytes = 4
	if pm.RowBytes != 4 {
		t.Errorf("Expected preserved RowBytes 4, got %d", pm.RowBytes)
	}

	// Read the pixels manually from the pointer
	// We expect 2 rows of 4 bytes (stride), but valid width is 2
	// Wait, NewRasterImage copies. It copies Height * RowBytes.
	// Height=2, RowBytes=4. Total 8 bytes.
	// Row 0 of subset: [5, 6, x, x] (x is padding/garbage from original line)
	// Row 1 of subset: [9, 10, x, x]

	if pm.Addr == nil {
		t.Fatal("Pixmap Addr is nil")
	}

	// Unsafe cast to access memory. We know the size is at least 8 bytes.
	bytes := (*[8]byte)(pm.Addr)[:]

	// Check Row 0 (subset)
	if bytes[0] != 5 {
		t.Errorf("Subset(0,0) (Orig 1,1) should be 5, got %d", bytes[0])
	}
	if bytes[1] != 6 {
		t.Errorf("Subset(1,0) (Orig 2,1) should be 6, got %d", bytes[1])
	}

	// Check Row 1 (subset)
	// Stride is 4, so next row starts at index 4
	if bytes[4] != 9 {
		t.Errorf("Subset(0,1) (Orig 1,2) should be 9, got %d", bytes[4])
	}
	if bytes[5] != 10 {
		t.Errorf("Subset(1,1) (Orig 2,2) should be 10, got %d", bytes[5])
	}

	// Test invalid subset
	invalidRect := models.NewIRect(-1, 0, 2, 2)
	if img.MakeSubset(invalidRect) != nil {
		t.Error("MakeSubset should fail for negative coordinates")
	}
	invalidRect2 := models.NewIRect(3, 3, 2, 2) // Goes out of bounds (4x4)
	if img.MakeSubset(invalidRect2) != nil {
		t.Error("MakeSubset should fail for out of bounds")
	}
}
