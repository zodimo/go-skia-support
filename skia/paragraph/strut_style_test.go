package paragraph

import (
	"testing"
)

func TestStrutStyleDefaults(t *testing.T) {
	s := NewStrutStyle()

	if s.GetFontSize() != DefaultFontSize {
		t.Errorf("Expected default FontSize %f, got %f", DefaultFontSize, s.GetFontSize())
	}
	if s.GetHeight() != 1.0 {
		t.Errorf("Expected default Height 1.0, got %f", s.GetHeight())
	}
	if s.GetLeading() != 0.0 {
		t.Errorf("Expected default Leading 0.0, got %f", s.GetLeading())
	}
	if s.GetStrutEnabled() {
		t.Error("Expected default StrutEnabled to be false")
	}
	if s.GetForceStrutHeight() {
		t.Error("Expected default ForceStrutHeight to be false")
	}
	if s.GetHeightOverride() {
		t.Error("Expected default HeightOverride to be false")
	}
	if s.GetHalfLeading() {
		t.Error("Expected default HalfLeading to be false")
	}
	if len(s.GetFontFamilies()) != 0 {
		t.Error("Expected empty default FontFamilies")
	}
}

func TestStrutStyleProperties(t *testing.T) {
	s := NewStrutStyle()

	s.SetFontSize(20.0)
	if s.GetFontSize() != 20.0 {
		t.Errorf("Expected FontSize 20.0, got %f", s.GetFontSize())
	}

	s.SetHeight(1.5)
	if s.GetHeight() != 1.5 {
		t.Errorf("Expected Height 1.5, got %f", s.GetHeight())
	}

	s.SetLeading(2.0)
	if s.GetLeading() != 2.0 {
		t.Errorf("Expected Leading 2.0, got %f", s.GetLeading())
	}

	s.SetStrutEnabled(true)
	if !s.GetStrutEnabled() {
		t.Error("Expected StrutEnabled to be true")
	}

	s.SetForceStrutHeight(true)
	if !s.GetForceStrutHeight() {
		t.Error("Expected ForceStrutHeight to be true")
	}

	s.SetHeightOverride(true)
	if !s.GetHeightOverride() {
		t.Error("Expected HeightOverride to be true")
	}

	s.SetHalfLeading(true)
	if !s.GetHalfLeading() {
		t.Error("Expected HalfLeading to be true")
	}

	families := []string{"Arial", "Helvetica"}
	s.SetFontFamilies(families)
	if len(s.GetFontFamilies()) != 2 {
		t.Errorf("Expected 2 FontFamilies, got %d", len(s.GetFontFamilies()))
	}
	if s.GetFontFamilies()[0] != "Arial" {
		t.Errorf("Expected first family Arial, got %s", s.GetFontFamilies()[0])
	}
}

func TestStrutStyleEquals(t *testing.T) {
	s1 := NewStrutStyle()
	s2 := NewStrutStyle()

	if !s1.Equals(&s2) {
		t.Error("Expected default StrutStyles to be equal")
	}

	s1.SetFontSize(24)
	if s1.Equals(&s2) {
		t.Error("Expected StrutStyles with different FontSize to be unequal")
	}
	s2.SetFontSize(24)
	if !s1.Equals(&s2) {
		t.Error("Expected StrutStyles with same FontSize to be equal")
	}

	s1.SetFontFamilies([]string{"Roboto"})
	if s1.Equals(&s2) {
		t.Error("Expected StrutStyles with different FontFamilies to be unequal")
	}
	s2.SetFontFamilies([]string{"Roboto"})
	if !s1.Equals(&s2) {
		t.Error("Expected StrutStyles with same FontFamilies to be equal")
	}
}
