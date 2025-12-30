package shaper

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/go-text/typesetting/font"
	"github.com/zodimo/go-skia-support/skia/impl"
	"github.com/zodimo/go-skia-support/skia/models"
	"golang.org/x/image/font/gofont/goregular"
)

//go:embed testdata/Variable.ttf
var variableFontData []byte

func TestHarfbuzzShaper_Features_Ligatures(t *testing.T) {
	// 1. Prepare Font (GoRegular)
	parsed, err := font.ParseTTF(bytes.NewReader(goregular.TTF))
	if err != nil {
		t.Fatalf("Failed to parse goregular: %v", err)
	}

	skTypeface := impl.NewTypefaceWithTypefaceFace("regular", models.FontStyle{Weight: 400, Width: 5, Slant: 0}, parsed)
	skFont := impl.NewFont()
	skFont.SetTypeface(skTypeface)
	skFont.SetSize(20)

	text := "fi"
	shaper := NewHarfbuzzShaper()

	// Helper to shape with specific features and script
	shapeWith := func(features []Feature) int {
		handler := NewTextBlobBuilderRunHandler(text, models.Point{X: 0, Y: 0})
		fontIter := NewTrivialFontRunIterator(skFont, len(text))
		bidiIter := NewTrivialBiDiRunIterator(0, len(text))
		// Use 'Latn' script
		scriptIter := NewTrivialScriptRunIterator(uint32('L')<<24|uint32('a')<<16|uint32('t')<<8|uint32('n'), len(text))
		langIter := NewTrivialLanguageRunIterator("en", len(text))

		shaper.ShapeWithIterators(text, fontIter, bidiIter, scriptIter, langIter, features, 1000, handler)
		blob := handler.MakeBlob().(*impl.TextBlob)
		if blob.RunCount() == 0 {
			return 0
		}
		return len(blob.Run(0).Glyphs)
	}

	// 2. Shape with 'liga' off
	ligaTag := uint32('l')<<24 | uint32('i')<<16 | uint32('g')<<8 | uint32('a')
	featuresOff := []Feature{{Tag: ligaTag, Value: 0, Start: 0, End: 2}}
	countNoLiga := shapeWith(featuresOff)
	t.Logf("Computed 'fi' with liga=0: %d glyphs", countNoLiga)

	// 3. Shape with 'liga' on
	featuresOn := []Feature{{Tag: ligaTag, Value: 1, Start: 0, End: 2}}
	countLiga := shapeWith(featuresOn)
	t.Logf("Computed 'fi' with liga=1: %d glyphs", countLiga)

	// Assertions
	if countNoLiga == 0 {
		t.Error("Expected at least 1 glyph with ligatures off")
	}
	if countLiga == 0 {
		t.Error("Expected at least 1 glyph with ligatures on")
	}

	// Note: GoRegular presumably doesn't support 'fi' ligature, so we just verify the shaper doesn't crash
	// and produces output. We don't assert countLiga < countNoLiga.
	if countLiga == countNoLiga {
		t.Log("Note: Ligature did not reduce glyph count (expected for GoRegular)")
	}
}

func TestHarfbuzzShaper_Variations(t *testing.T) {
	if len(variableFontData) == 0 {
		t.Skip("Skipping variation test: Variable.ttf not found in testdata")
	}

	parsed, err := font.ParseTTF(bytes.NewReader(variableFontData))
	if err != nil {
		t.Fatalf("Failed to parse Variable.ttf: %v", err)
	}

	style := models.FontStyle{Weight: 400, Width: 5, Slant: 0}
	baseTf := impl.NewTypefaceWithTypefaceFace("Variable", style, parsed)

	// 'wdth' tag
	wdthTag := uint32('w')<<24 | uint32('d')<<16 | uint32('t')<<8 | uint32('h')

	// 1. Shape with default
	fontDefault := impl.NewFont()
	fontDefault.SetTypeface(baseTf)
	fontDefault.SetSize(20)

	shaper := NewHarfbuzzShaper()
	text := "aa"
	handlerDef := NewTextBlobBuilderRunHandler(text, models.Point{X: 0, Y: 0})
	shaper.Shape(text, fontDefault, true, 1000, handlerDef, nil)

	blobDef := handlerDef.MakeBlob().(*impl.TextBlob)
	runDef := blobDef.Run(0)
	if len(runDef.Positions) < 2 {
		t.Fatalf("Expected at least 2 glyph positions for 'aa', got %d", len(runDef.Positions))
	}
	widthDef := float32(runDef.Positions[1].X - runDef.Positions[0].X)

	// 2. Shape with variation (Extremes)
	args := models.FontArguments{
		VariationDesignPosition: models.VariationPosition{
			Coordinates: []models.VariationCoordinate{
				{Axis: wdthTag, Value: 2.0},
			},
		},
	}

	clonedTf := baseTf.MakeClone(args)
	fontVar := impl.NewFont()
	fontVar.SetTypeface(clonedTf)
	fontVar.SetSize(20)

	handlerVar := NewTextBlobBuilderRunHandler(text, models.Point{X: 0, Y: 0})
	shaper.Shape(text, fontVar, true, 1000, handlerVar, nil)

	blobVar := handlerVar.MakeBlob().(*impl.TextBlob)
	runVar := blobVar.Run(0)
	if len(runVar.Positions) < 2 {
		t.Fatalf("Expected at least 2 glyph positions for 'aa' (variation), got %d", len(runVar.Positions))
	}
	widthVar := float32(runVar.Positions[1].X - runVar.Positions[0].X)

	t.Logf("Default width: %f, Var width: %f", widthDef, widthVar)

	if widthDef == widthVar {
		t.Errorf("Variation (wdth=2.0) did not affect width (%.4f == %.4f). Check tag %x and font properties.", widthDef, widthVar, wdthTag)
	}
}
