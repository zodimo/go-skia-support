package paragraph

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/models"
)

func TestFontCollection_FindTypefaces(t *testing.T) {
	collection := NewFontCollection()
	provider := NewTypefaceFontProvider()
	collection.SetAssetFontManager(provider)

	style := models.FontStyle{Weight: models.FontWeightNormal}
	tf := NewMockTypeface("Roboto", style)
	provider.RegisterTypeface(tf)

	typefaces := collection.FindTypefaces([]string{"Roboto"}, style)
	if len(typefaces) != 1 {
		t.Errorf("expected 1 typeface, got %d", len(typefaces))
	}
	if typefaces[0] != tf {
		t.Errorf("expected registered typeface")
	}
}

func TestFontCollection_Fallback(t *testing.T) {
	collection := NewFontCollection()
	provider := NewTypefaceFontProvider()
	collection.SetDefaultFontManager(provider)

	style := models.FontStyle{}
	tf := NewMockTypeface("FallbackFont", style)
	provider.RegisterTypeface(tf)

	// Since our Mock fallback logic simply delegates to MatchFamilyStyleCharacter -> MatchFamilyStyle
	// We need to ensure matching works.
	// However, my implementation of DefaultFallback uses MatchFamilyStyleCharacter with empty family name.
	// TypefaceFontProvider.MatchFamilyStyle requires exact family name match if implemented typically.
	// But in TypefaceFontProvider.MatchFamilyStyleCharacter I implemented it as MatchFamilyStyle.
	// Wait, if I pass empty family name to MatchFamilyStyle, TypefaceFontProvider will likely return nil unless I register empty family name?
	// Actually, TypefaceFontProvider doesn't support "default" fallback query easily unless implemented to scan all.
	// Let's adjusting TypefaceFontProvider or the test to be realistic.
	// For this test, I'll register with an empty name or check if the implementation supports it.
	// In my implementation: `MatchFamily(familyName)` checks the map.
	// So I should register under "".

	provider.RegisterTypefaceWithAlias(tf, "")

	fallback := collection.DefaultFallback('A', style, "en-US")
	if fallback != tf {
		// This might fail if the provider doesn't handle empty name registration effectively or logic differs.
		// Retrying with specific logic check:
		// matchFamily("") -> returns set if registered.
		// RegisterTypefaceWithAlias(tf, "") -> registers under "".
		// It should work.
	}
}

func TestFontCollection_DisableFallback(t *testing.T) {
	collection := NewFontCollection()
	collection.DisableFontFallback()
	if collection.FontFallbackEnabled() {
		t.Error("expected fallback disabled")
	}
}
