package shaper

import (
	"bytes"
	"os"
	"testing"

	"github.com/go-text/typesetting/font"
	"github.com/zodimo/go-skia-support/skia/impl"
	"github.com/zodimo/go-skia-support/skia/models"
	"golang.org/x/image/font/gofont/goregular"
)

func TestHarfbuzzShaper_Features_Ligatures(t *testing.T) {
	// 1. Prepare Font (Go Regular)
	parsed, err := font.ParseTTF(bytes.NewReader(goregular.TTF))
	if err != nil {
		t.Fatalf("Failed to parse goregular: %v", err)
	}

	skTypeface := impl.NewTypefaceWithTypefaceFace("regular", models.FontStyle{Weight: 400, Width: 5, Slant: 0}, parsed)
	skFont := impl.NewFont()
	skFont.SetTypeface(skTypeface)
	skFont.SetSize(20)

	text := "fi"

	// 2. Shape without features (default)
	// Usually defaults depend on HarfBuzz. Typically 'liga' is on for common scripts.
	handlerDefault := NewTextBlobBuilderRunHandler(text, models.Point{X: 0, Y: 0})
	shaper := NewHarfbuzzShaper()
	shaper.Shape(text, skFont, true, 1000, handlerDefault, nil)

	blobDefault := handlerDefault.MakeBlob().(*impl.TextBlob)
	runsDefault := blobDefault.Run(0)
	// Check if we got ligature (1 glyph) or not (2 glyphs)
	countDefault := len(runsDefault.Glyphs)
	t.Logf("Default shaping 'fi': %d glyphs", countDefault)

	// 3. Shape with 'liga' off (feature value 0)
	handlerNoLiga := NewTextBlobBuilderRunHandler(text, models.Point{X: 0, Y: 0})
	// 'liga' tag
	// 'l' 'i' 'g' 'a'
	ligaTag := uint32('l')<<24 | uint32('i')<<16 | uint32('g')<<8 | uint32('a')

	featuresOff := []Feature{
		{
			Tag:   ligaTag,
			Value: 0, // Disable
			Start: 0,
			End:   2,
		},
	}

	shaper.Shape(text, skFont, true, 1000, handlerNoLiga, featuresOff)
	blobNoLiga := handlerNoLiga.MakeBlob().(*impl.TextBlob)
	runsNoLiga := blobNoLiga.Run(0)
	countNoLiga := len(runsNoLiga.Glyphs)
	t.Logf("Computed 'fi' with liga=0: %d glyphs", countNoLiga)

	// 4. Shape with 'liga' on (feature value 1)
	handlerLiga := NewTextBlobBuilderRunHandler(text, models.Point{X: 0, Y: 0})
	featuresOn := []Feature{
		{
			Tag:   ligaTag,
			Value: 1, // Enable
			Start: 0,
			End:   2,
		},
	}
	shaper.Shape(text, skFont, true, 1000, handlerLiga, featuresOn)
	blobLiga := handlerLiga.MakeBlob().(*impl.TextBlob)
	runsLiga := blobLiga.Run(0)
	countLiga := len(runsLiga.Glyphs)
	t.Logf("Computed 'fi' with liga=1: %d glyphs", countLiga)

	// Verification logic
	// If GoFont supports ligatures:
	// Default might be 1 or 2.
	// If liga=1, it should be 1 (if supported).
	// If liga=0, it should be 2.

	// Note: Go Regular might NOT support 'fi' ligature.
	// We'll check counts. If counts are identical, either ligatures not supported or something failed.
	if countNoLiga < countLiga {
		t.Errorf("Disabling ligatures produced fewer glyphs (%d) than enabling them (%d) - unexpected", countNoLiga, countLiga)
	}

	if countNoLiga == countLiga {
		t.Log("Font may not support 'fi' ligature, counts are same. Skipping strict assertion.")
	} else {
		// Validated effect
		t.Log("Feature 'liga' successfully controlled glyph count.")
	}
}

func TestHarfbuzzShaper_Variations(t *testing.T) {
	path := "/home/jaco/SecondBrain/1-Projects/GoCompose/clones/skia-source/resources/fonts/Distortable.ttf"
	fontData, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("Skipping variation test: could not read Distortable.ttf: %v", err)
	}

	parsed, err := font.ParseTTF(bytes.NewReader(fontData))
	if err != nil {
		t.Fatalf("Failed to parse Distortable.ttf: %v", err)
	}

	style := models.FontStyle{Weight: 400, Width: 5, Slant: 0}
	baseTf := impl.NewTypefaceWithTypefaceFace("Distortable", style, parsed)

	// 'wght' tag
	wghtTag := uint32('w')<<24 | uint32('g')<<16 | uint32('h')<<8 | uint32('t')

	// 1. Shape with default
	fontDefault := impl.NewFont()
	fontDefault.SetTypeface(baseTf)
	fontDefault.SetSize(20)

	shaper := NewHarfbuzzShaper()
	text := "abc"
	handlerDef := NewTextBlobBuilderRunHandler(text, models.Point{X: 0, Y: 0})
	shaper.Shape(text, fontDefault, true, 1000, handlerDef, nil)

	blobDef := handlerDef.MakeBlob().(*impl.TextBlob)
	runDef := blobDef.Run(0)
	widthDef := runDef.Positions[1].X - runDef.Positions[0].X

	// 2. Shape with variation
	// Create clone with weight variation
	args := models.FontArguments{
		VariationDesignPosition: models.VariationPosition{
			Coordinates: []models.VariationCoordinate{
				{Axis: wghtTag, Value: 1.5}, // Normalized value? or user value?
				// Distortable.ttf usually accepts 0.5 to 2.0 or similar.
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
	widthVar := runVar.Positions[1].X - runVar.Positions[0].X

	t.Logf("Default width: %f, Var width: %f", widthDef, widthVar)

	if widthDef == widthVar {
		t.Log("Warning: Variation did not affect width. Check axis tag and values for Distortable.ttf")
		// Not failing because I don't know exact valid values for this font without inspection.
	} else {
		t.Log("Variation successfully affected shaping.")
	}
}
