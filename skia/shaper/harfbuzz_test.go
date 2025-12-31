package shaper

import (
	"bytes"
	"testing"

	"github.com/go-text/typesetting/font"
	"github.com/zodimo/go-skia-support/skia/impl"
	"github.com/zodimo/go-skia-support/skia/models"
	"golang.org/x/image/font/gofont/goregular"
)

func TestHarfbuzzShaper_Shape(t *testing.T) {
	// 1. Prepare Font
	parsed, err := font.ParseTTF(bytes.NewReader(goregular.TTF))
	if err != nil {
		t.Fatalf("Failed to parse goremular: %v", err)
	}

	skTypeface := impl.NewTypefaceWithTypefaceFace("regular", models.FontStyle{Weight: 400, Width: 5, Slant: 0}, parsed)
	skFont := impl.NewFont()
	skFont.SetTypeface(skTypeface)
	skFont.SetSize(16)

	// 2. Prepare Handler
	text := "Hello World"
	handler := NewTextBlobBuilderRunHandler(text, models.Point{X: 0, Y: 0})

	// 3. Shape
	shaper := NewHarfbuzzShaper()
	shaper.Shape(text, skFont, true, 1000, handler, nil)

	// 4. Verify
	blob := handler.MakeBlob()
	if blob == nil {
		t.Fatal("MakeBlob returned nil")
	}

	implBlob := blob.(*impl.TextBlob)
	if implBlob.RunCount() == 0 {
		t.Fatal("Expected runs, got 0")
	}

	run := implBlob.Run(0)
	if len(run.Glyphs) == 0 {
		t.Errorf("Expected glyphs, got 0")
	}

	t.Logf("Got %d glyphs in run 0", len(run.Glyphs))
	// We can expect strict glyph count equal to run length for "Hello World" in Latin?
	// Usually yes (1 to 1).
	if len(run.Glyphs) != len(text) {
		// Ligatures might happen, but usually not for Hello World in regular font.
		t.Logf("Glyph count %d != text length %d (ligatures?)", len(run.Glyphs), len(text))
	}
}

func TestHarfbuzzShaper_ScaleX(t *testing.T) {
	// 1. Prepare Font
	parsed, err := font.ParseTTF(bytes.NewReader(goregular.TTF))
	if err != nil {
		t.Fatalf("Failed to parse goremular: %v", err)
	}

	skTypeface := impl.NewTypefaceWithTypefaceFace("regular", models.FontStyle{Weight: 400, Width: 5, Slant: 0}, parsed)

	// Font 1: Normal scale
	skFont1 := impl.NewFont()
	skFont1.SetTypeface(skTypeface)
	skFont1.SetSize(16)
	skFont1.SetScaleX(1.0)

	// Font 2: Scaled X
	skFont2 := impl.NewFont()
	skFont2.SetTypeface(skTypeface)
	skFont2.SetSize(16)
	skFont2.SetScaleX(1.5)

	text := "Hello World"

	// Helper to get width
	getWidth := func(f *impl.Font) float32 {
		handler := NewTextBlobBuilderRunHandler(text, models.Point{X: 0, Y: 0})
		shaper := NewHarfbuzzShaper()
		shaper.Shape(text, f, true, 1000, handler, nil)
		blob := handler.MakeBlob()
		if blob == nil {
			t.Fatal("MakeBlob returned nil")
		}
		implBlob := blob.(*impl.TextBlob)
		if implBlob.RunCount() == 0 {
			t.Fatal("Expected runs, got 0")
		}

		run := implBlob.Run(0)
		// Calculate total advance
		var width float32
		if len(run.Positions) > 0 {
			lastPos := run.Positions[len(run.Positions)-1]
			// We need the advance of the last glyph to get total width,
			// but for this simple check, comparing the X position of the last glyph
			// plus its advance (if we had it) would be ideal.
			// However, since we are just checking for *difference*,
			// comparing the last glyph's position is a good proxy.
			// Actually, HarfBuzz positions are absolute if we don't reset them?
			// TextBlobBuilderRunHandler accumulates them?
			// Let's look at the positions.
			width = float32(lastPos.X)
		}
		return width
	}

	width1 := getWidth(skFont1)
	width2 := getWidth(skFont2)

	t.Logf("Width 1.0: %f", width1)
	t.Logf("Width 1.5: %f", width2)

	// We expect width2 to be approximately 1.5 * width1
	// Allow some margin for error
	expected := width1 * 1.5
	diff := width2 - expected
	if diff < 0 {
		diff = -diff
	}

	if diff > 5.0 { // Generous tolerance
		t.Errorf("Expected width2 (~%f) to be approx 1.5x width1 (%f), but got %f. Diff: %f", expected, width1, width2, diff)
	} else {
		t.Logf("Success: Widths are scaled as expected.")
	}
}
