package paragraph

import (
	"bytes"
	"testing"

	"github.com/go-text/typesetting/font"
	"github.com/zodimo/go-skia-support/skia/impl"
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
	"github.com/zodimo/go-skia-support/skia/shaper"
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

	// 4. Run Shape
	bidiRegions := []BidiRegion{{Start: 0, End: len(text), Level: 0}}
	shaper := NewOneLineShaper(text, []Block{block}, nil, fc, impl.NewSkUnicode(), bidiRegions)

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

func TestOneLineShaper_ScriptDetection(t *testing.T) {
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

	// 3. Setup Mixed Text: Latin + Greek
	// "Hello " (Latin + Common space) -> Should be Latin
	// "Γεια" (Greek) -> Should be Greek
	text := "Hello Γεια"
	style := NewTextStyle()
	style.FontFamilies = []string{"GoRegular"}
	style.FontSize = 16
	block := NewBlock(0, len(text), style)

	// 4. Run Shape
	bidiRegions := []BidiRegion{{Start: 0, End: len(text), Level: 0}}
	ols := NewOneLineShaper(text, []Block{block}, nil, fc, impl.NewSkUnicode(), bidiRegions)

	// 4. Run Shape
	if !ols.Shape() {
		t.Fatal("Shape returned false")
	}

	// 5. Verify
	// Print runs
	for i, r := range ols.Runs {
		t.Logf("Run %d: Script %x, Glyphs %d, Range %v", i, r.Script(), r.Size(), r.TextRange())
	}

	if len(ols.Runs) < 2 {
		t.Errorf("Expected at least 2 runs for mixed script, got %d", len(ols.Runs))
	}

	// Check scripts
	// We expect Latin then Greek
	// Run 0: Latin (includes space)
	// Run 1: Greek

	if len(ols.Runs) >= 2 {
		run0 := ols.Runs[0]
		run1 := ols.Runs[1]

		if run0.Script() != shaper.ScriptLatin {
			t.Errorf("Run 0: Expected script Latin (%x), got %x", shaper.ScriptLatin, run0.Script())
		}

		greekTag := makeTag("Grek")
		if run1.Script() != greekTag {
			t.Errorf("Run 1: Expected script Greek (%x), got %x", greekTag, run1.Script())
		}
	}
}

// Helper to match script_iterator.go
func makeTag(s string) uint32 {
	if len(s) != 4 {
		return 0
	}
	return uint32(s[0])<<24 | uint32(s[1])<<16 | uint32(s[2])<<8 | uint32(s[3])
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
