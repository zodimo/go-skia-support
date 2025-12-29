package paragraph

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/models"
)

func TestTypefaceFontProvider_RegisterTypeface(t *testing.T) {
	provider := NewTypefaceFontProvider()
	tf := NewMockTypeface("Roboto", models.FontStyle{})

	count := provider.RegisterTypeface(tf)
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}

	families := provider.CountFamilies()
	if families != 1 {
		t.Errorf("expected 1 family, got %d", families)
	}

	name := provider.GetFamilyName(0)
	if name != "Roboto" {
		t.Errorf("expected 'Roboto', got '%s'", name)
	}
}

func TestTypefaceFontProvider_MatchFamily(t *testing.T) {
	provider := NewTypefaceFontProvider()
	tf := NewMockTypeface("Roboto", models.FontStyle{Weight: models.FontWeightNormal, Slant: models.FontSlantUpright})
	provider.RegisterTypeface(tf)

	set := provider.MatchFamily("Roboto")
	if set == nil {
		t.Errorf("expected style set, got nil")
	}
	if set.Count() != 1 {
		t.Errorf("expected 1 style, got %d", set.Count())
	}
}

func TestTypefaceFontProvider_MatchFamilyStyle(t *testing.T) {
	provider := NewTypefaceFontProvider()
	styleNormal := models.FontStyle{Weight: models.FontWeightNormal, Slant: models.FontSlantUpright}
	styleBold := models.FontStyle{Weight: models.FontWeightBold, Slant: models.FontSlantUpright}

	tfNormal := NewMockTypeface("Roboto", styleNormal)
	tfBold := NewMockTypeface("Roboto", styleBold)

	provider.RegisterTypeface(tfNormal)
	provider.RegisterTypeface(tfBold)

	match := provider.MatchFamilyStyle("Roboto", styleBold)
	if match != tfBold {
		t.Errorf("expected bold typeface, got %v", match)
	}
}

func TestTypefaceFontProvider_Alias(t *testing.T) {
	provider := NewTypefaceFontProvider()
	tf := NewMockTypeface("Roboto", models.FontStyle{})

	provider.RegisterTypefaceWithAlias(tf, "MyAlias")

	set := provider.MatchFamily("MyAlias")
	if set == nil {
		t.Errorf("expected match for alias 'MyAlias'")
	}
}
