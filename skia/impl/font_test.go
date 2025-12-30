package impl

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/zodimo/go-skia-support/skia/enums"
)

func TestNewFont(t *testing.T) {
	f := NewFont()
	if f == nil {
		t.Fatal("NewFont returned nil")
	}
	if f.Size() != FontDefaultSize {
		t.Errorf("Expected size %v, got %v", FontDefaultSize, f.Size())
	}
	if f.ScaleX() != FontDefaultScaleX {
		t.Errorf("Expected scaleX %v, got %v", FontDefaultScaleX, f.ScaleX())
	}
	if f.SkewX() != FontDefaultSkewX {
		t.Errorf("Expected skewX %v, got %v", FontDefaultSkewX, f.SkewX())
	}
	if f.Typeface() == nil {
		t.Error("Default font should have a typeface")
	}
	if f.Edging() != enums.FontEdgingDefault {
		t.Errorf("Expected edging %v, got %v", enums.FontEdgingDefault, f.Edging())
	}
	if f.Hinting() != enums.FontHintingDefault {
		t.Errorf("Expected hinting %v, got %v", enums.FontHintingDefault, f.Hinting())
	}
}

func TestFontWithTypeface(t *testing.T) {
	tf := NewTypeface("TestFont", FontStyle{Weight: 400, Width: 5, Slant: 0})
	f := NewFontWithTypeface(tf)
	if f.Typeface() != tf {
		t.Error("Font should have the typeface that was set")
	}
}

func TestFontWithTypefaceAndSize(t *testing.T) {
	tf := NewTypeface("TestFont", FontStyle{Weight: 400, Width: 5, Slant: 0})
	f := NewFontWithTypefaceAndSize(tf, 24.0)
	if f.Size() != 24.0 {
		t.Errorf("Expected size 24.0, got %v", f.Size())
	}
}

func TestFontSetSize(t *testing.T) {
	f := NewFont()
	f.SetSize(20.0)
	if f.Size() != 20.0 {
		t.Errorf("Expected size 20.0, got %v", f.Size())
	}

	// Negative size should be ignored
	f.SetSize(-5.0)
	if f.Size() != 20.0 {
		t.Errorf("Negative size should be ignored, size is %v", f.Size())
	}

	// Zero size should be allowed
	f.SetSize(0.0)
	if f.Size() != 0.0 {
		t.Errorf("Zero size should be allowed, size is %v", f.Size())
	}
}

func TestFontSetScaleX(t *testing.T) {
	f := NewFont()
	f.SetScaleX(2.0)
	if f.ScaleX() != 2.0 {
		t.Errorf("Expected scaleX 2.0, got %v", f.ScaleX())
	}
}

func TestFontSetSkewX(t *testing.T) {
	f := NewFont()
	f.SetSkewX(0.5)
	if f.SkewX() != 0.5 {
		t.Errorf("Expected skewX 0.5, got %v", f.SkewX())
	}
}

func TestFontEdging(t *testing.T) {
	f := NewFont()
	f.SetEdging(enums.FontEdgingSubpixelAntiAlias)
	if f.Edging() != enums.FontEdgingSubpixelAntiAlias {
		t.Errorf("Expected edging SubpixelAntiAlias, got %v", f.Edging())
	}
}

func TestFontHinting(t *testing.T) {
	f := NewFont()
	f.SetHinting(enums.FontHintingFull)
	if f.Hinting() != enums.FontHintingFull {
		t.Errorf("Expected hinting Full, got %v", f.Hinting())
	}
}

func TestFontFlags(t *testing.T) {
	f := NewFont()

	// Test ForceAutoHinting
	if f.IsForceAutoHinting() {
		t.Error("ForceAutoHinting should be false by default")
	}
	f.SetForceAutoHinting(true)
	if !f.IsForceAutoHinting() {
		t.Error("ForceAutoHinting should be true after setting")
	}

	// Test EmbeddedBitmaps
	if f.IsEmbeddedBitmaps() {
		t.Error("EmbeddedBitmaps should be false by default")
	}
	f.SetEmbeddedBitmaps(true)
	if !f.IsEmbeddedBitmaps() {
		t.Error("EmbeddedBitmaps should be true after setting")
	}

	// Test Subpixel
	if f.IsSubpixel() {
		t.Error("Subpixel should be false by default")
	}
	f.SetSubpixel(true)
	if !f.IsSubpixel() {
		t.Error("Subpixel should be true after setting")
	}

	// Test LinearMetrics
	if f.IsLinearMetrics() {
		t.Error("LinearMetrics should be false by default")
	}
	f.SetLinearMetrics(true)
	if !f.IsLinearMetrics() {
		t.Error("LinearMetrics should be true after setting")
	}

	// Test Embolden
	if f.IsEmbolden() {
		t.Error("Embolden should be false by default")
	}
	f.SetEmbolden(true)
	if !f.IsEmbolden() {
		t.Error("Embolden should be true after setting")
	}

	// Test BaselineSnap (default is true)
	if !f.IsBaselineSnap() {
		t.Error("BaselineSnap should be true by default")
	}
	f.SetBaselineSnap(false)
	if f.IsBaselineSnap() {
		t.Error("BaselineSnap should be false after setting to false")
	}
}

func TestFontMeasureText(t *testing.T) {
	f := NewFont()
	text := []byte("Hello")
	width := f.MeasureText(text, enums.TextEncodingUTF8, nil)
	if width <= 0 {
		t.Errorf("Expected positive width, got %v", width)
	}

	// Test with bounds
	var bounds Rect
	width2 := f.MeasureText(text, enums.TextEncodingUTF8, &bounds)
	if width != width2 {
		t.Errorf("Width should be same with or without bounds: %v vs %v", width, width2)
	}
	if bounds.Right <= bounds.Left {
		t.Error("Bounds should have positive width")
	}
	if bounds.Bottom <= bounds.Top {
		t.Error("Bounds should have positive height")
	}

	// Empty text
	width = f.MeasureText([]byte{}, enums.TextEncodingUTF8, nil)
	if width != 0 {
		t.Errorf("Empty text should have zero width, got %v", width)
	}
}

func TestFontEquals(t *testing.T) {
	f1 := NewFont()
	f2 := NewFont()

	// Different typefaces means not equal (each has unique ID)
	if f1.Equals(f2) {
		t.Error("Two fonts with different typefaces should not be equal")
	}

	// Same font should equal itself
	if !f1.Equals(f1) {
		t.Error("Font should equal itself")
	}

	// Nil checks
	var nilFont *Font
	if f1.Equals(nilFont) {
		t.Error("Font should not equal nil")
	}
	if !nilFont.Equals(nilFont) {
		t.Error("Two nils should be equal")
	}
}

// ============================================================================
// Tests ported from C++ skia-source/tests/FontTest.cpp
// ============================================================================

// TestFontFlatten tests all font attribute combinations for consistency.
// Ported from: FontTest.cpp Font_flatten (lines 74-116)
// This test verifies that fonts with various attribute combinations maintain
// their properties correctly (equivalent to serialize/deserialize round-trip in C++).
func TestFontFlatten(t *testing.T) {
	// Test values from C++ FontTest.cpp lines 75-90
	sizes := []Scalar{0, 0.001, 1, 10, 10.001, 100000.01}
	scalesX := []Scalar{-5, 0, 1, 5}
	skewsX := []Scalar{-5, 0, 5}
	edgings := []enums.FontEdging{
		enums.FontEdgingAlias,
		enums.FontEdgingSubpixelAntiAlias,
	}
	hintings := []enums.FontHinting{
		enums.FontHintingNone,
		enums.FontHintingFull,
	}

	// Flag combinations
	type fontFlags struct {
		forceAutoHinting bool
		embeddedBitmaps  bool
		subpixel         bool
		linearMetrics    bool
		embolden         bool
		baselineSnap     bool
	}

	flagCombinations := []fontFlags{
		{true, false, false, false, false, true},   // ForceAutoHinting only
		{false, true, false, false, false, true},   // EmbeddedBitmaps only
		{false, false, true, false, false, true},   // Subpixel only
		{false, false, false, true, false, true},   // LinearMetrics only
		{false, false, false, false, true, true},   // Embolden only
		{false, false, false, false, false, false}, // BaselineSnap disabled
		{true, true, true, true, true, true},       // All enabled
	}

	testCount := 0
	failCount := 0

	for _, size := range sizes {
		for _, scaleX := range scalesX {
			for _, skewX := range skewsX {
				for _, edging := range edgings {
					for _, hinting := range hintings {
						for _, flags := range flagCombinations {
							testCount++

							// Create font with all attributes
							font := NewFont()
							font.SetSize(size)
							font.SetScaleX(scaleX)
							font.SetSkewX(skewX)
							font.SetEdging(edging)
							font.SetHinting(hinting)
							font.SetForceAutoHinting(flags.forceAutoHinting)
							font.SetEmbeddedBitmaps(flags.embeddedBitmaps)
							font.SetSubpixel(flags.subpixel)
							font.SetLinearMetrics(flags.linearMetrics)
							font.SetEmbolden(flags.embolden)
							font.SetBaselineSnap(flags.baselineSnap)

							// Verify all attributes were set correctly
							// Size: negative values are ignored in SetSize
							if size >= 0 && font.Size() != size {
								t.Errorf("Size: got %v, want %v", font.Size(), size)
								failCount++
							}
							if font.ScaleX() != scaleX {
								t.Errorf("ScaleX: got %v, want %v", font.ScaleX(), scaleX)
								failCount++
							}
							if font.SkewX() != skewX {
								t.Errorf("SkewX: got %v, want %v", font.SkewX(), skewX)
								failCount++
							}
							if font.Edging() != edging {
								t.Errorf("Edging: got %v, want %v", font.Edging(), edging)
								failCount++
							}
							if font.Hinting() != hinting {
								t.Errorf("Hinting: got %v, want %v", font.Hinting(), hinting)
								failCount++
							}
							if font.IsForceAutoHinting() != flags.forceAutoHinting {
								t.Errorf("ForceAutoHinting: got %v, want %v",
									font.IsForceAutoHinting(), flags.forceAutoHinting)
								failCount++
							}
							if font.IsEmbeddedBitmaps() != flags.embeddedBitmaps {
								t.Errorf("EmbeddedBitmaps: got %v, want %v",
									font.IsEmbeddedBitmaps(), flags.embeddedBitmaps)
								failCount++
							}
							if font.IsSubpixel() != flags.subpixel {
								t.Errorf("Subpixel: got %v, want %v",
									font.IsSubpixel(), flags.subpixel)
								failCount++
							}
							if font.IsLinearMetrics() != flags.linearMetrics {
								t.Errorf("LinearMetrics: got %v, want %v",
									font.IsLinearMetrics(), flags.linearMetrics)
								failCount++
							}
							if font.IsEmbolden() != flags.embolden {
								t.Errorf("Embolden: got %v, want %v",
									font.IsEmbolden(), flags.embolden)
								failCount++
							}
							if font.IsBaselineSnap() != flags.baselineSnap {
								t.Errorf("BaselineSnap: got %v, want %v",
									font.IsBaselineSnap(), flags.baselineSnap)
								failCount++
							}

							// Early exit if too many failures
							if failCount > 10 {
								t.Fatalf("Too many failures, stopping. Tested %d combinations.", testCount)
							}
						}
					}
				}
			}
		}
	}

	t.Logf("Tested %d font attribute combinations", testCount)
}

func TestFontGetWidths(t *testing.T) {
	f := NewFont()
	f.SetSize(20)
	glyphs := []uint16{1, 2, 3}
	widths := f.GetWidths(glyphs)

	if len(widths) != 3 {
		t.Errorf("Expected 3 widths, got %d", len(widths))
	}

	// Verify fallback width: 20 * 0.6 = 12
	expected := Scalar(12.0)
	for i, w := range widths {
		if w != expected {
			t.Errorf("Width %d: expected %v, got %v", i, expected, w)
		}
	}

	// Test nil/empty
	if f.GetWidths(nil) != nil {
		t.Error("Expected nil for nil glyphs")
	}
}

func TestFontGetMetrics(t *testing.T) {
	f := NewFont()
	f.SetSize(20)
	metrics := f.GetMetrics()

	// Verify fallback metrics
	// Ascent: -20 * 0.8 = -16
	// Descent: 20 * 0.2 = 4
	// Leading: 20 * 0.05 = 1
	if metrics.Ascent != -16 {
		t.Errorf("Expected Ascent -16, got %v", metrics.Ascent)
	}
	if metrics.Descent != 4 {
		t.Errorf("Expected Descent 4, got %v", metrics.Descent)
	}
	if metrics.Leading != 1 {
		t.Errorf("Expected Leading 1, got %v", metrics.Leading)
	}
}

func TestFontMeasureTextEncodings(t *testing.T) {
	f := NewFont()
	f.SetSize(10.0) // Char width = 10 * 0.6 = 6.0
	f.SetScaleX(1.0)
	expectedWidthPerChar := Scalar(6.0)

	tests := []struct {
		name      string
		encoding  enums.TextEncoding
		input     []byte
		wantChars int
	}{
		{
			name:      "UTF8",
			encoding:  enums.TextEncodingUTF8,
			input:     []byte("Hello"),
			wantChars: 5,
		},
		{
			name:     "UTF16LE",
			encoding: enums.TextEncodingUTF16,
			input: func() []byte {
				buf := new(bytes.Buffer)
				runes := []rune("Hello")
				for _, r := range runes {
					binary.Write(buf, binary.LittleEndian, uint16(r))
				}
				return buf.Bytes()
			}(),
			wantChars: 5,
		},
		{
			name:      "UTF16 Invalid Length",
			encoding:  enums.TextEncodingUTF16,
			input:     []byte{0x00, 0x00, 0x00}, // 3 bytes
			wantChars: 0,
		},
		{
			name:     "UTF32LE",
			encoding: enums.TextEncodingUTF32,
			input: func() []byte {
				buf := new(bytes.Buffer)
				runes := []rune("Hello")
				for _, r := range runes {
					binary.Write(buf, binary.LittleEndian, uint32(r))
				}
				return buf.Bytes()
			}(),
			wantChars: 5,
		},
		{
			name:     "GlyphID",
			encoding: enums.TextEncodingGlyphID,
			input: func() []byte {
				buf := new(bytes.Buffer)
				glyphs := []uint16{1, 2, 3}
				for _, g := range glyphs {
					binary.Write(buf, binary.LittleEndian, g)
				}
				return buf.Bytes()
			}(),
			wantChars: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := f.MeasureText(tt.input, tt.encoding, nil)
			expected := expectedWidthPerChar * Scalar(tt.wantChars)
			// Allow small float epsilon if needed, but here it's simple multiplication
			if width != expected {
				t.Errorf("MeasureText() = %v, want %v", width, expected)
			}
		})
	}
}
