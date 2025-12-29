package shaper

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/impl"
	"github.com/zodimo/go-skia-support/skia/models"
)

// IntegrationMockFont implements identity mapping for UnicharToGlyph
type IntegrationMockFont struct {
	*impl.Font
}

func NewIntegrationMockFont() *IntegrationMockFont {
	return &IntegrationMockFont{
		Font: impl.NewFont(),
	}
}

func (m *IntegrationMockFont) UnicharToGlyph(unichar rune) uint16 {
	return uint16(unichar)
}

func (m *IntegrationMockFont) GetWidths(glyphs []uint16) []models.Scalar {
	widths := make([]models.Scalar, len(glyphs))
	for i := range glyphs {
		widths[i] = 10.0 // Constant width
	}
	return widths
}

func TestPrimitiveShaper_Integration(t *testing.T) {
	shaper := NewPrimitiveShaper()
	text := "Hello World"
	font := NewIntegrationMockFont()
	startPoint := models.Point{X: 10, Y: 20}

	// Create handler with explicit type to access MakeBlob
	handler := NewTextBlobBuilderRunHandler(text, startPoint)

	// Shape
	shaper.Shape(text, font, true, 200, handler)

	// Get Blob
	blob := handler.MakeBlob()
	if blob == nil {
		t.Fatal("MakeBlob returned nil - verify if any runs were committed")
	}

	// Verify Blob contents
	implBlob, ok := blob.(*impl.TextBlob)
	if !ok {
		t.Fatal("Blob is not *impl.TextBlob")
	}

	if implBlob.RunCount() != 1 {
		t.Fatalf("Expected 1 run, got %d", implBlob.RunCount())
	}

	run := implBlob.Run(0)
	if run == nil {
		t.Fatal("Run(0) is nil")
	}

	if len(run.Glyphs) != len(text) {
		t.Errorf("Expected %d glyphs, got %d", len(text), len(run.Glyphs))
	}

	// Verify glyphs match identity mapping
	// If the bug exists, these will likely be all 0s
	zeros := 0
	for i, r := range text {
		expectedGlyph := uint16(r)
		gotGlyph := uint16(run.Glyphs[i])
		if gotGlyph == 0 {
			zeros++
		}
		if gotGlyph != expectedGlyph {
			t.Errorf("Glyph mismatch at %d: expected %d (%c), got %d", i, expectedGlyph, r, gotGlyph)
		}
	}

	if zeros == len(text) {
		t.Logf("All glyphs are zero - confirms bug in TextBlobBuilderRunHandler not copying data")
	}
}
