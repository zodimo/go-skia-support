package shaper

import (
	"bytes"
	"testing"

	"github.com/go-text/typesetting/font"
	"github.com/zodimo/go-skia-support/skia/impl"
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
	"golang.org/x/image/font/gofont/goregular"
)

// MockTypefaceForHarfbuzz implements SkTypeface and UseGoTextFace.
type MockTypefaceForHarfbuzz struct {
	interfaces.SkTypeface
	face *font.Face
}

func (m *MockTypefaceForHarfbuzz) GoTextFace() *font.Face {
	return m.face
}

// We rely on embedding interfaces.SkTypeface to satisfy the interface,
// but we need a concrete implementation to forward to.
// impl.Typeface provides basic impl.
// However, impl.Typeface usage in testing:
// we should probably just embed *impl.Typeface struct.

func NewMockTypefaceForHarfbuzz(face *font.Face) *MockTypefaceForHarfbuzz {
	return &MockTypefaceForHarfbuzz{
		SkTypeface: impl.NewDefaultTypeface(),
		face:       face,
	}
}

func TestHarfbuzzShaper_Shape(t *testing.T) {
	// 1. Prepare Font
	parsed, err := font.ParseTTF(bytes.NewReader(goregular.TTF))
	if err != nil {
		t.Fatalf("Failed to parse goremular: %v", err)
	}

	mockTypeface := NewMockTypefaceForHarfbuzz(parsed)
	skFont := impl.NewFont()
	skFont.SetTypeface(mockTypeface)
	skFont.SetSize(16)

	// 2. Prepare Handler
	text := "Hello World"
	handler := NewTextBlobBuilderRunHandler(text, models.Point{X: 0, Y: 0})

	// 3. Shape
	shaper := NewHarfbuzzShaper()
	shaper.Shape(text, skFont, true, 1000, handler)

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
