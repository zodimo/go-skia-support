package impl

import (
	"bytes"
	"testing"

	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/font/opentype"
	"golang.org/x/image/font/gofont/goregular"
)

// Helper to create a typeface with a real font
func newTypefaceWithGoRegular(t *testing.T) *Typeface {
	t.Helper()

	// Parse the embedded Go Regular font
	loader, err := opentype.NewLoader(bytes.NewReader(goregular.TTF))
	if err != nil {
		t.Fatalf("failed to load Go Regular font: %v", err)
	}

	goFont, err := font.NewFont(loader)
	if err != nil {
		t.Fatalf("failed to parse Go Regular font: %v", err)
	}

	face := font.NewFace(goFont)
	return NewTypefaceWithTypefaceFace("Go", FontStyle{Weight: 400, Width: 5, Slant: 0}, face)
}

// ============================================================================
// Tests for glyph data access methods with real Go Regular font
// ============================================================================

func TestTypeface_UnitsPerEm_RealFont(t *testing.T) {
	tf := newTypefaceWithGoRegular(t)

	upem := tf.UnitsPerEm()
	if upem <= 0 {
		t.Errorf("UnitsPerEm should be positive, got %d", upem)
	}
	// Go fonts typically have 2048 units per em
	if upem != 2048 {
		t.Logf("Note: UnitsPerEm is %d (expected 2048 for Go fonts)", upem)
	}
}

func TestTypeface_GetGlyphAdvance_RealFont(t *testing.T) {
	tf := newTypefaceWithGoRegular(t)

	// Get the glyph ID for 'A'
	glyphA := tf.UnicharToGlyph('A')
	if glyphA == 0 {
		t.Fatal("UnicharToGlyph should return non-zero for 'A'")
	}

	// Get the advance width
	advance := tf.GetGlyphAdvance(glyphA)
	if advance <= 0 {
		t.Errorf("GetGlyphAdvance should return positive value for 'A', got %d", advance)
	}

	// Space should also have an advance
	glyphSpace := tf.UnicharToGlyph(' ')
	advanceSpace := tf.GetGlyphAdvance(glyphSpace)
	if advanceSpace <= 0 {
		t.Errorf("GetGlyphAdvance should return positive value for space, got %d", advanceSpace)
	}

	// 'A' should typically be wider than 'i'
	glyphI := tf.UnicharToGlyph('i')
	advanceI := tf.GetGlyphAdvance(glyphI)
	if advanceI <= 0 {
		t.Errorf("GetGlyphAdvance should return positive value for 'i', got %d", advanceI)
	}
	// In a proportional font, 'A' is usually wider than 'i'
	if advance <= advanceI {
		t.Logf("Note: 'A' advance (%d) should typically be > 'i' advance (%d) in proportional font", advance, advanceI)
	}
}

func TestTypeface_GetGlyphBounds_RealFont(t *testing.T) {
	tf := newTypefaceWithGoRegular(t)

	// Get bounds for 'A'
	glyphA := tf.UnicharToGlyph('A')
	boundsA := tf.GetGlyphBounds(glyphA)

	// Bounds should have positive width and height
	width := boundsA.Right - boundsA.Left
	height := boundsA.Bottom - boundsA.Top

	if width <= 0 {
		t.Errorf("Glyph 'A' should have positive width, got %f", width)
	}
	if height <= 0 {
		t.Errorf("Glyph 'A' should have positive height, got %f", height)
	}

	// Space should have zero or near-zero bounds (no ink)
	glyphSpace := tf.UnicharToGlyph(' ')
	boundsSpace := tf.GetGlyphBounds(glyphSpace)
	spaceWidth := boundsSpace.Right - boundsSpace.Left
	if spaceWidth > 0 {
		t.Logf("Space glyph has non-zero bounds width: %f (may contain ink)", spaceWidth)
	}

	t.Logf("'A' bounds: Left=%f, Top=%f, Right=%f, Bottom=%f",
		boundsA.Left, boundsA.Top, boundsA.Right, boundsA.Bottom)
}

func TestTypeface_GetGlyphPath_RealFont(t *testing.T) {
	tf := newTypefaceWithGoRegular(t)

	// Get path for 'A' - should have outline data
	glyphA := tf.UnicharToGlyph('A')
	pathA, err := tf.GetGlyphPath(glyphA)

	if err != nil {
		t.Errorf("GetGlyphPath for 'A' should not error, got: %v", err)
	}
	if pathA == nil {
		t.Error("GetGlyphPath for 'A' should return a path")
	}

	// Path should have points and verbs
	if pathA != nil {
		pointCount := pathA.CountPoints()
		verbCount := pathA.CountVerbs()

		if pointCount == 0 {
			t.Error("Path for 'A' should have points")
		}
		if verbCount == 0 {
			t.Error("Path for 'A' should have verbs")
		}

		t.Logf("'A' path: %d points, %d verbs", pointCount, verbCount)
	}

	// Space glyph should have no outline (error expected)
	glyphSpace := tf.UnicharToGlyph(' ')
	pathSpace, err := tf.GetGlyphPath(glyphSpace)
	if err == nil && pathSpace != nil && pathSpace.CountPoints() > 0 {
		t.Logf("Space glyph unexpectedly has outline with %d points", pathSpace.CountPoints())
	}
}

func TestTypeface_NoFontFace_ReturnsZeroDefaults(t *testing.T) {
	// Test with no goTextFace set
	tf := NewDefaultTypeface()

	if tf.UnitsPerEm() != 0 {
		t.Errorf("UnitsPerEm without font face should be 0, got %d", tf.UnitsPerEm())
	}

	if tf.GetGlyphAdvance(1) != 0 {
		t.Errorf("GetGlyphAdvance without font face should be 0")
	}

	bounds := tf.GetGlyphBounds(1)
	if bounds.Left != 0 || bounds.Right != 0 || bounds.Top != 0 || bounds.Bottom != 0 {
		t.Errorf("GetGlyphBounds without font face should return zero rect")
	}

	_, err := tf.GetGlyphPath(1)
	if err == nil {
		t.Error("GetGlyphPath without font face should return error")
	}
}
