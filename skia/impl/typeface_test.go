package impl

import (
	"testing"
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
