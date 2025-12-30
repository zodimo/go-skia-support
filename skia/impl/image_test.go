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
