package paragraph

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/impl"
	"github.com/zodimo/go-skia-support/skia/models"
	"github.com/zodimo/go-skia-support/skia/shaper"
)

// MockFont is a simple wrapper to provide controlling font properties in tests.
// Although we can use impl.Font, this mocked version makes intent clearer for some tests
// if needed, but for now we'll mostly use impl.Font since it's concrete.
// Actually, using impl.Font is better.

func TestNewRun(t *testing.T) {
	font := impl.NewFontWithTypefaceAndSize(nil, 20.0)

	// Create a dummy RunInfo
	info := shaper.RunInfo{
		Font:       font,
		BidiLevel:  0,
		Script:     1234, // Random script
		Language:   "en",
		Advance:    models.Point{X: 100, Y: 0},
		GlyphCount: 5,
		Utf8Range:  shaper.Range{Begin: 0, End: 5},
	}

	firstChar := 10
	run := NewRun(info, firstChar, 1.0, false, 0, 0, 0)

	if run.Size() != 5 {
		t.Errorf("Expected size 5, got %d", run.Size())
	}
	if run.Font() != font {
		t.Error("Expected font to match")
	}
	if run.TextRange().Start != 10 || run.TextRange().End != 15 {
		t.Errorf("Expected text range [10, 15], got %v", run.TextRange())
	}
	if run.clusterStart != 10 {
		t.Errorf("Expected cluster start 10, got %d", run.clusterStart)
	}
	if len(run.Positions()) != 6 {
		t.Errorf("Expected 6 positions (including trailing), got %d", len(run.Positions()))
	}
	if run.Positions()[5].X != 100 {
		t.Errorf("Expected trailing position X 100, got %f", run.Positions()[5].X)
	}
}

func TestRun_Metrics(t *testing.T) {
	// Size 20 font
	font := impl.NewFontWithTypefaceAndSize(nil, 20.0)
	// Default impl.Font metrics:
	// Ascent: -16 (0.8 * 20)
	// Descent: 4  (0.2 * 20)
	// Leading: 1  (0.05 * 20)

	info := shaper.RunInfo{
		Font:       font,
		GlyphCount: 0,
		Utf8Range:  shaper.Range{Begin: 0, End: 0},
	}

	// Case 1: Standard metrics (no multiplier, no half-leading)
	run := NewRun(info, 0, 0, false, 0, 0, 0)

	expectedAscent := float32(-16.0 - 1.0*0.5) // -16.5
	expectedDescent := float32(4.0 + 1.0*0.5)  // 4.5
	expectedLeading := float32(0)              // 0
	_ = expectedLeading                        // Suppress unused error if strictly enforcing, but we can verify it
	if !nearlyEqual(run.CorrectLeading(), expectedLeading) {
		t.Errorf("Expected leading %f, got %f", expectedLeading, run.CorrectLeading())
	}

	if !nearlyEqual(run.CorrectAscent(), expectedAscent) {
		t.Errorf("Expected ascent %f, got %f", expectedAscent, run.CorrectAscent())
	}
	if !nearlyEqual(run.CorrectDescent(), expectedDescent) {
		t.Errorf("Expected descent %f, got %f", expectedDescent, run.CorrectDescent())
	}

	// Case 2: Height Multiplier
	run = NewRun(info, 0, 2.0, false, 0, 0, 0)

	// Intrinsic height = 4.5 - (-16.5) = 21.0
	// Target height = 2.0 * 20.0 = 40.0
	// Multiplier = 40.0 / 21.0 = 1.90476
	// Ascent = -16.5 * 1.90476 = -31.4285
	// Descent = 4.5 * 1.90476 = 8.5714

	expectedAscent = -16.5 * (40.0 / 21.0)
	expectedDescent = 4.5 * (40.0 / 21.0)

	if !nearlyEqual(run.CorrectAscent(), expectedAscent) {
		t.Errorf("Expected multiplied ascent %f, got %f", expectedAscent, run.CorrectAscent())
	}

	// Case 3: Half-Leading
	run = NewRun(info, 0, 2.0, true, 0, 0, 0)

	// Intrinsic height = 21.0
	// Target height = 40.0
	// Extra leading = (40.0 - 21.0) / 2 = 9.5
	// Ascent = -16.5 - 9.5 = -26.0
	// Descent = 4.5 + 9.5 = 14.0

	expectedAscent = -16.5 - 9.5
	expectedDescent = 4.5 + 9.5

	if !nearlyEqual(run.CorrectAscent(), expectedAscent) {
		t.Errorf("Expected half-leading ascent %f, got %f", expectedAscent, run.CorrectAscent())
	}
	if !nearlyEqual(run.CorrectDescent(), expectedDescent) {
		t.Errorf("Expected half-leading descent %f, got %f", expectedDescent, run.CorrectDescent())
	}
}

func TestRun_Direction(t *testing.T) {
	// LTR
	infoLTR := shaper.RunInfo{BidiLevel: 0, GlyphCount: 1, Utf8Range: shaper.Range{Begin: 0, End: 1}, Font: impl.NewFont()}
	runLTR := NewRun(infoLTR, 0, 1.0, false, 0, 0, 0)
	if !runLTR.LeftToRight() {
		t.Error("Expected LTR for level 0")
	}
	if runLTR.TextDirection() != TextDirectionLTR {
		t.Error("Expected TextDirection LTR")
	}

	// RTL
	infoRTL := shaper.RunInfo{BidiLevel: 1, GlyphCount: 1, Utf8Range: shaper.Range{Begin: 0, End: 1}, Font: impl.NewFont()}
	runRTL := NewRun(infoRTL, 0, 1.0, false, 0, 0, 0)
	if runRTL.LeftToRight() {
		t.Error("Expected RTL for level 1")
	}
	if runRTL.TextDirection() != TextDirectionRTL {
		t.Error("Expected TextDirection RTL")
	}
}

func TestRun_NewRunBuffer(t *testing.T) {
	info := shaper.RunInfo{
		Font:       impl.NewFont(),
		GlyphCount: 5,
		Utf8Range:  shaper.Range{Begin: 0, End: 5},
	}
	run := NewRun(info, 0, 0, false, 0, 0, 0)
	buffer := run.NewRunBuffer()

	if len(buffer.Glyphs) != 5 {
		t.Errorf("Expected buffer glyphs length 5, got %d", len(buffer.Glyphs))
	}
	if len(buffer.Positions) != 6 {
		t.Errorf("Expected buffer positions length 6, got %d", len(buffer.Positions))
	}
	if len(buffer.Clusters) != 6 {
		t.Errorf("Expected buffer clusters length 6, got %d", len(buffer.Clusters))
	}

	// Verify buffer points to run data by modifying buffer and checking run
	buffer.Glyphs[0] = 100
	if run.Glyphs()[0] != 100 {
		t.Error("Buffer should point to Run's glyph storage")
	}
}

func TestRun_CursiveScript(t *testing.T) {
	// ARAB
	info := shaper.RunInfo{
		Script:     makeFourByteTag('A', 'r', 'a', 'b'),
		GlyphCount: 1,
		Utf8Range:  shaper.Range{Begin: 0, End: 1},
		Font:       impl.NewFont(),
	}
	run := NewRun(info, 0, 0, false, 0, 0, 0)
	if !run.IsCursiveScript() {
		t.Error("Expected Arabic to be cursive")
	}

	// Latn
	info.Script = makeFourByteTag('L', 'a', 't', 'n')
	run = NewRun(info, 0, 0, false, 0, 0, 0)
	if run.IsCursiveScript() {
		t.Error("Expected Latin to NOT be cursive")
	}
}

func TestRun_CalculateWidth(t *testing.T) {
	info := shaper.RunInfo{
		Font:       impl.NewFont(),
		GlyphCount: 3,
		Utf8Range:  shaper.Range{Begin: 0, End: 3},
	}
	run := NewRun(info, 0, 0, false, 0, 0, 0)

	// Setup positions: 0, 10, 20, 30
	run.Positions()[0] = models.Point{X: 0, Y: 0}
	run.Positions()[1] = models.Point{X: 10, Y: 0}
	run.Positions()[2] = models.Point{X: 20, Y: 0}
	run.Positions()[3] = models.Point{X: 30, Y: 0}

	width := run.CalculateWidth(0, 3, false)
	if width != 30 {
		t.Errorf("Expected width 30, got %f", width)
	}

	width = run.CalculateWidth(1, 2, false)
	if width != 10 { // 20 - 10
		t.Errorf("Expected width 10, got %f", width)
	}
}

func TestRun_IsResolved(t *testing.T) {
	info := shaper.RunInfo{
		Font:       impl.NewFont(),
		GlyphCount: 2,
		Utf8Range:  shaper.Range{Begin: 0, End: 2},
	}
	run := NewRun(info, 0, 0, false, 0, 0, 0)

	// All zeros (default)
	if run.IsResolved() {
		// Actually, 0 means .notdef usually? Or missing glyph?
		// IsResolved checks if glyph != 0. 0 is traditionally .notdef in many fonts, but handled as resolved?
		// Wait, Run.cpp implementation:
		// bool Run::isResolved() const {
		//   for (auto& glyph :fGlyphs) {
		//       if (glyph == 0) {
		//           return false;
		//       }
		//   }
		//   return true;
		// }
		// So 0 is considered UNRESOLVED (or failed/missing).
		t.Error("Expected unresolved for zero glyphs")
	}

	run.Glyphs()[0] = 1
	run.Glyphs()[1] = 2
	if !run.IsResolved() {
		t.Error("Expected resolved for non-zero glyphs")
	}
}
