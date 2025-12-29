package paragraph

import (
	"bytes"
	"testing"

	"github.com/go-text/typesetting/font"
	"github.com/zodimo/go-skia-support/skia/impl"
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
	"golang.org/x/image/font/gofont/goregular"
)

func TestOneLineShaper_Shape_Basic(t *testing.T) {
	// 1. Setup Real Font
	parsed, err := font.ParseTTF(bytes.NewReader(goregular.TTF))
	if err != nil {
		t.Fatalf("Failed to parse gofont: %v", err)
	}

	skTypeface := impl.NewTypefaceWithTypefaceFace("GoRegular", models.FontStyle{}, parsed)

	// 2. Setup FontCollection
	fc := NewFontCollection()
	mgr := &FakeFontMgr{typeface: skTypeface}
	fc.SetDefaultFontManager(mgr)

	// 3. Setup Shaper
	text := "Hello OneLineShaper"
	style := NewTextStyle()
	style.FontFamilies = []string{"GoRegular"}
	style.FontSize = 16
	block := NewBlock(0, len(text), style)

	shaper := NewOneLineShaper(text, []Block{block}, nil, fc, impl.NewSkUnicode())

	// 4. Run Shape
	if !shaper.Shape() {
		t.Fatal("Shape returned false")
	}

	// 5. Verify
	if len(shaper.Runs) == 0 {
		t.Fatal("Expected runs, got 0")
	}

	run := shaper.Runs[0]
	if run.Size() != len(text) {
		t.Errorf("Expected %d glyphs, got %d", len(text), run.Size())
	}

	t.Logf("Run output: %d glyphs, Advance: %v", run.Size(), run.Advance())
}

// FakeFontMgr implements a minimal SkFontMgr for testing
type FakeFontMgr struct {
	interfaces.SkFontMgr
	typeface interfaces.SkTypeface
}

func (m *FakeFontMgr) MatchFamilyStyle(familyName string, style models.FontStyle) interfaces.SkTypeface {
	return m.typeface
}

func (m *FakeFontMgr) MatchFamilyStyleCharacter(familyName string, style models.FontStyle, bcp47 []string, character rune) interfaces.SkTypeface {
	return m.typeface
}

func (m *FakeFontMgr) LegacyMakeTypeface(familyName string, style models.FontStyle) interfaces.SkTypeface {
	return m.typeface
}
