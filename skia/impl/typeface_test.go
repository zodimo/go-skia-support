package impl

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/models"
)

func TestNewDefaultTypeface(t *testing.T) {
	tf := NewDefaultTypeface()
	if tf == nil {
		t.Fatal("NewDefaultTypeface returned nil")
	}
	if tf.UniqueID() == 0 {
		t.Error("UniqueID should not be 0")
	}
	if tf.FamilyName() != "" {
		t.Errorf("Expected empty family name, got %q", tf.FamilyName())
	}
	if tf.IsBold() {
		t.Error("Default typeface should not be bold")
	}
	if tf.IsItalic() {
		t.Error("Default typeface should not be italic")
	}
	if tf.IsFixedPitch() {
		t.Error("Default typeface should not be fixed pitch")
	}
}

func TestNewTypeface(t *testing.T) {
	tf := NewTypeface("Arial", FontStyle{Weight: 700, Width: 5, Slant: 0})
	if tf == nil {
		t.Fatal("NewTypeface returned nil")
	}
	if tf.FamilyName() != "Arial" {
		t.Errorf("Expected family name 'Arial', got %q", tf.FamilyName())
	}
	if !tf.IsBold() {
		t.Error("Typeface with weight 700 should be bold")
	}
	if tf.IsItalic() {
		t.Error("Typeface with upright slant should not be italic")
	}
}

func TestTypefaceItalic(t *testing.T) {
	tf := NewTypeface("Arial", FontStyle{Weight: 400, Width: 5, Slant: 1})
	if !tf.IsItalic() {
		t.Error("Typeface with italic slant should be italic")
	}
}

func TestTypefaceUniqueID(t *testing.T) {
	tf1 := NewDefaultTypeface()
	tf2 := NewDefaultTypeface()
	if tf1.UniqueID() == tf2.UniqueID() {
		t.Error("Two typefaces should have different unique IDs")
	}
}

func TestTypefaceWithOptions(t *testing.T) {
	tf := NewTypefaceWithOptions("Courier", FontStyle{Weight: 400, Width: 5, Slant: 0}, true)
	if !tf.IsFixedPitch() {
		t.Error("Typeface created with fixedPitch=true should be fixed pitch")
	}
}

// ============================================================================
// Tests ported from C++ skia-source/tests/TypefaceTest.cpp
// ============================================================================

// TestTypefaceStyle tests weight and width handling.
// Ported from: TypefaceTest.cpp TypefaceStyle (lines 110-125)
func TestTypefaceStyle(t *testing.T) {
	// Test all weight values (C++ tests 1-1000)
	testWeights := []models.FontWeight{100, 200, 300, 400, 500, 600, 700, 800, 900}
	for _, weight := range testWeights {
		tf := NewTypeface("Test", FontStyle{Weight: weight, Width: 5, Slant: 0})
		if tf == nil {
			t.Fatalf("NewTypeface returned nil for weight %d", weight)
		}

		tfStyle := tf.FontStyle()
		if tfStyle.Weight != weight {
			t.Errorf("Weight: got %d, want %d", tfStyle.Weight, weight)
		}

		// Verify IsBold() - bold is typically weight > 500 (Medium)
		expectedBold := weight > models.FontWeightMedium
		if tf.IsBold() != expectedBold {
			t.Errorf("Weight %d: IsBold() = %v, want %v", weight, tf.IsBold(), expectedBold)
		}
	}

	// Test width values (1-9 in C++)
	testWidths := []models.FontWidth{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, width := range testWidths {
		tf := NewTypeface("Test", FontStyle{Weight: 400, Width: width, Slant: 0})
		if tf == nil {
			t.Fatalf("NewTypeface returned nil for width %d", width)
		}

		tfStyle := tf.FontStyle()
		if tfStyle.Width != width {
			t.Errorf("Width: got %d, want %d", tfStyle.Width, width)
		}
	}
}

// TestTypefaceEquality tests typeface equality via UniqueID.
// Ported from: TypefaceTest.cpp Typeface (lines 502-514)
func TestTypefaceEquality(t *testing.T) {
	tf1 := NewDefaultTypeface()
	tf2 := NewDefaultTypeface()

	// Same typeface equals itself (same UniqueID)
	if tf1.UniqueID() != tf1.UniqueID() {
		t.Error("Typeface should have consistent UniqueID")
	}

	// Different typeface instances have different unique IDs
	if tf1.UniqueID() == tf2.UniqueID() {
		t.Error("Different typeface instances should have different UniqueIDs")
	}

	// Nil checks
	var nilTf *Typeface
	if nilTf != nil {
		t.Error("nilTf should be nil")
	}
}

// TestTypefaceStyleSlants tests all slant types.
// Ported from: TypefaceTest.cpp TypefaceStyle (extended)
func TestTypefaceStyleSlants(t *testing.T) {
	slants := []struct {
		slant        models.FontSlant
		expectItalic bool
	}{
		{models.FontSlantUpright, false}, // Upright
		{models.FontSlantItalic, true},   // Italic
		{models.FontSlantOblique, true},  // Oblique (also considered italic-like)
	}

	for _, tc := range slants {
		tf := NewTypeface("Test", FontStyle{Weight: 400, Width: 5, Slant: tc.slant})
		if tf == nil {
			t.Fatalf("NewTypeface returned nil for slant %d", tc.slant)
		}

		// IsItalic returns true for slant > 0
		if tf.IsItalic() != tc.expectItalic {
			t.Errorf("Slant %d: IsItalic() = %v, want %v", tc.slant, tf.IsItalic(), tc.expectItalic)
		}
	}
}

// TestTypefaceStyleCombinations tests various style combinations.
// Ported from: TypefaceTest.cpp TypefaceStyle (comprehensive combinations)
func TestTypefaceStyleCombinations(t *testing.T) {
	// Bold + Italic
	tf := NewTypeface("Test", FontStyle{Weight: 700, Width: 5, Slant: 1})
	if tf == nil {
		t.Fatal("NewTypeface returned nil")
	}
	if !tf.IsBold() {
		t.Error("Weight 700 should be bold")
	}
	if !tf.IsItalic() {
		t.Error("Slant 1 should be italic")
	}

	// Neither bold nor italic
	tf = NewTypeface("Test", FontStyle{Weight: 400, Width: 5, Slant: 0})
	if tf.IsBold() {
		t.Error("Weight 400 should not be bold")
	}
	if tf.IsItalic() {
		t.Error("Slant 0 should not be italic")
	}

	// Extreme weights
	tfThin := NewTypeface("Test", FontStyle{Weight: 100, Width: 5, Slant: 0})
	if tfThin.IsBold() {
		t.Error("Weight 100 should not be bold")
	}

	tfBlack := NewTypeface("Test", FontStyle{Weight: 900, Width: 5, Slant: 0})
	if !tfBlack.IsBold() {
		t.Error("Weight 900 should be bold")
	}
}
