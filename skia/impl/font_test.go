package impl

import (
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
