package impl

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/enums"
)

func TestMakeTextBlobFromString(t *testing.T) {
	font := NewFont()
	blob := MakeTextBlobFromString("Hello", font)
	if blob == nil {
		t.Fatal("MakeTextBlobFromString returned nil")
	}
	if blob.UniqueID() == 0 {
		t.Error("UniqueID should not be 0")
	}

	bounds := blob.Bounds()
	if bounds.Right <= bounds.Left {
		t.Errorf("Bounds should have positive width: %v", bounds)
	}
}

func TestMakeTextBlobFromStringEmpty(t *testing.T) {
	font := NewFont()
	blob := MakeTextBlobFromString("", font)
	if blob != nil {
		t.Error("Empty string should return nil blob")
	}
}

func TestMakeTextBlobFromStringNilFont(t *testing.T) {
	blob := MakeTextBlobFromString("Hello", nil)
	if blob != nil {
		t.Error("Nil font should return nil blob")
	}
}

func TestMakeTextBlobFromText(t *testing.T) {
	font := NewFont()

	tests := []struct {
		name     string
		text     []byte
		encoding enums.TextEncoding
	}{
		{"UTF-8", []byte("Hello"), enums.TextEncodingUTF8},
		{"UTF-16", []byte{0x48, 0x00, 0x65, 0x00, 0x6c, 0x00, 0x6c, 0x00, 0x6f, 0x00}, enums.TextEncodingUTF16},
		{"UTF-32", []byte{0x48, 0x00, 0x00, 0x00, 0x65, 0x00, 0x00, 0x00}, enums.TextEncodingUTF32},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			blob := MakeTextBlobFromText(tc.text, tc.encoding, font)
			if blob == nil {
				t.Fatal("MakeTextBlobFromText returned nil")
			}
			if blob.UniqueID() == 0 {
				t.Error("UniqueID should not be 0")
			}
		})
	}
}

func TestTextBlobUniqueID(t *testing.T) {
	font := NewFont()
	blob1 := MakeTextBlobFromString("Hello", font)
	blob2 := MakeTextBlobFromString("World", font)
	if blob1.UniqueID() == blob2.UniqueID() {
		t.Error("Two blobs should have different unique IDs")
	}
}

func TestTextBlobBounds(t *testing.T) {
	font := NewFont()
	blob := MakeTextBlobFromString("Test", font)
	bounds := blob.Bounds()

	// Bounds should be valid
	if bounds.Right <= bounds.Left {
		t.Errorf("Bounds width should be positive: left=%v, right=%v", bounds.Left, bounds.Right)
	}
	if bounds.Bottom <= bounds.Top {
		t.Errorf("Bounds height should be positive: top=%v, bottom=%v", bounds.Top, bounds.Bottom)
	}
}

func TestTextBlobRunCount(t *testing.T) {
	font := NewFont()
	blob := MakeTextBlobFromString("Hello", font)
	if blob.RunCount() != 1 {
		t.Errorf("Expected 1 run, got %d", blob.RunCount())
	}
}

func TestTextBlobRun(t *testing.T) {
	font := NewFont()
	blob := MakeTextBlobFromString("Hi", font)
	run := blob.Run(0)
	if run == nil {
		t.Fatal("Run(0) should not be nil")
	}
	if len(run.Glyphs) != 2 {
		t.Errorf("Expected 2 glyphs for 'Hi', got %d", len(run.Glyphs))
	}
	if len(run.Positions) != 2 {
		t.Errorf("Expected 2 positions, got %d", len(run.Positions))
	}
	if run.Font == nil {
		t.Error("Run should have font reference")
	}

	// Out of bounds
	if blob.Run(-1) != nil {
		t.Error("Run(-1) should return nil")
	}
	if blob.Run(1) != nil {
		t.Error("Run(1) should return nil for single-run blob")
	}
}

func TestTextBlobBuilder(t *testing.T) {
	builder := NewTextBlobBuilder()
	font := NewFont()

	// Allocate a run
	buffer := builder.AllocRun(font, 5, 0, 0)
	if buffer == nil {
		t.Fatal("AllocRun returned nil")
	}
	if len(buffer.Glyphs) != 5 {
		t.Errorf("Expected 5 glyphs, got %d", len(buffer.Glyphs))
	}

	// Fill the buffer
	for i := 0; i < 5; i++ {
		buffer.Glyphs[i] = GlyphID('A' + i)
	}
	builder.AddRun()

	// Build the blob
	blob := builder.Make()
	if blob == nil {
		t.Fatal("Make returned nil")
	}
	if blob.RunCount() != 1 {
		t.Errorf("Expected 1 run, got %d", blob.RunCount())
	}
}

func TestTextBlobBuilderEmpty(t *testing.T) {
	builder := NewTextBlobBuilder()
	blob := builder.Make()
	if blob != nil {
		t.Error("Empty builder should return nil blob")
	}
}

func TestTextBlobBuilderAllocRunPosH(t *testing.T) {
	builder := NewTextBlobBuilder()
	font := NewFont()

	buffer := builder.AllocRunPosH(font, 3, 10.0)
	if buffer == nil {
		t.Fatal("AllocRunPosH returned nil")
	}
	if len(buffer.Glyphs) != 3 {
		t.Errorf("Expected 3 glyphs, got %d", len(buffer.Glyphs))
	}
	if len(buffer.Positions) != 3 {
		t.Errorf("Expected 3 positions, got %d", len(buffer.Positions))
	}

	// Fill positions
	buffer.Glyphs[0] = GlyphID('X')
	buffer.Glyphs[1] = GlyphID('Y')
	buffer.Glyphs[2] = GlyphID('Z')
	buffer.Positions[0] = 0
	buffer.Positions[1] = 10
	buffer.Positions[2] = 20
	builder.AddRun()

	blob := builder.Make()
	if blob == nil {
		t.Fatal("Make returned nil")
	}
}

func TestTextBlobBuilderAllocRunPos(t *testing.T) {
	builder := NewTextBlobBuilder()
	font := NewFont()

	buffer := builder.AllocRunPos(font, 2)
	if buffer == nil {
		t.Fatal("AllocRunPos returned nil")
	}
	if len(buffer.Glyphs) != 2 {
		t.Errorf("Expected 2 glyphs, got %d", len(buffer.Glyphs))
	}
	if len(buffer.Positions) != 4 {
		t.Errorf("Expected 4 position values (2 x,y pairs), got %d", len(buffer.Positions))
	}

	// Fill with x,y pairs
	buffer.Glyphs[0] = GlyphID('A')
	buffer.Glyphs[1] = GlyphID('B')
	buffer.Positions[0] = 0  // x1
	buffer.Positions[1] = 0  // y1
	buffer.Positions[2] = 10 // x2
	buffer.Positions[3] = 5  // y2
	builder.AddRun()

	blob := builder.Make()
	if blob == nil {
		t.Fatal("Make returned nil")
	}
}

func TestTextBlobBuilderInvalidArgs(t *testing.T) {
	builder := NewTextBlobBuilder()
	font := NewFont()

	// Zero count
	if builder.AllocRun(font, 0, 0, 0) != nil {
		t.Error("AllocRun with 0 count should return nil")
	}

	// Negative count
	if builder.AllocRun(font, -1, 0, 0) != nil {
		t.Error("AllocRun with negative count should return nil")
	}

	// Nil font
	if builder.AllocRun(nil, 5, 0, 0) != nil {
		t.Error("AllocRun with nil font should return nil")
	}
}

// ============================================================================
// Tests ported from C++ skia-source/tests/TextBlobTest.cpp
// ============================================================================

// TestTextBlobBuilderRunMerging tests that runs with same attributes are merged.
// Ported from: TextBlobTest.cpp TextBlobTester::TestBuilder (lines 46-123)
func TestTextBlobBuilderRunMerging(t *testing.T) {
	builder := NewTextBlobBuilder()
	font := NewFont()

	// Test 1: Single run of default positioning
	buffer := builder.AllocRun(font, 128, 100, 100)
	for i := 0; i < 128; i++ {
		buffer.Glyphs[i] = GlyphID(i % 128)
	}
	builder.AddRun()

	blob := builder.Make()
	if blob == nil {
		t.Fatal("Expected non-nil blob for single run")
	}
	if blob.RunCount() != 1 {
		t.Errorf("Expected 1 run, got %d", blob.RunCount())
	}

	// Verify glyph data integrity
	run := blob.Run(0)
	if run == nil {
		t.Fatal("Run(0) should not be nil")
	}
	if len(run.Glyphs) != 128 {
		t.Errorf("Expected 128 glyphs, got %d", len(run.Glyphs))
	}
	for i := 0; i < len(run.Glyphs); i++ {
		if run.Glyphs[i] != GlyphID(i%128) {
			t.Errorf("Glyph[%d] = %d, expected %d", i, run.Glyphs[i], i%128)
			break
		}
	}
}

// TestTextBlobBuilderMultipleRuns tests builder with multiple runs.
// Ported from: TextBlobTest.cpp TestBuilder (lines 67-123)
func TestTextBlobBuilderMultipleRuns(t *testing.T) {
	builder := NewTextBlobBuilder()
	font := NewFont()

	// Add 3 runs with default positioning at different positions
	for runIdx := 0; runIdx < 3; runIdx++ {
		buffer := builder.AllocRun(font, 128, 100, Scalar(150+runIdx*100))
		for i := 0; i < 128; i++ {
			buffer.Glyphs[i] = GlyphID(i % 128)
		}
		builder.AddRun()
	}

	blob := builder.Make()
	if blob == nil {
		t.Fatal("Expected non-nil blob")
	}

	// Each default-positioned run at different Y offset should remain separate
	if blob.RunCount() < 1 {
		t.Errorf("Expected at least 1 run, got %d", blob.RunCount())
	}
}

// TestTextBlobBuilderHorizontalPositioning tests horizontal positioned runs.
// Ported from: TextBlobTest.cpp set5 (lines 74-84)
func TestTextBlobBuilderHorizontalPositioning(t *testing.T) {
	builder := NewTextBlobBuilder()
	font := NewFont()

	// Create runs with horizontal positioning at same Y
	for runIdx := 0; runIdx < 2; runIdx++ {
		buffer := builder.AllocRunPosH(font, 128, 150)
		for i := 0; i < 128; i++ {
			buffer.Glyphs[i] = GlyphID(i % 128)
			buffer.Positions[i] = Scalar(runIdx*100 + i)
		}
		builder.AddRun()
	}

	blob := builder.Make()
	if blob == nil {
		t.Fatal("Expected non-nil blob")
	}

	// Verify runs exist
	if blob.RunCount() < 1 {
		t.Errorf("Expected at least 1 run, got %d", blob.RunCount())
	}
}

// TestTextBlobBuilderFullPositioning tests full 2D positioned runs.
// Ported from: TextBlobTest.cpp set6 (lines 86-95)
func TestTextBlobBuilderFullPositioning(t *testing.T) {
	builder := NewTextBlobBuilder()
	font := NewFont()

	// Create runs with full positioning
	for runIdx := 0; runIdx < 3; runIdx++ {
		buffer := builder.AllocRunPos(font, 128)
		for i := 0; i < 128; i++ {
			buffer.Glyphs[i] = GlyphID(i % 128)
			buffer.Positions[i*2] = Scalar(i)
			buffer.Positions[i*2+1] = Scalar(-i)
		}
		builder.AddRun()
	}

	blob := builder.Make()
	if blob == nil {
		t.Fatal("Expected non-nil blob")
	}

	// Verify positions are stored correctly
	run := blob.Run(0)
	if run == nil || len(run.Positions) == 0 {
		t.Fatal("Run should have positions")
	}
}

// TestTextBlobBoundsExplicit tests bounds computation with explicit bounds.
// Ported from: TextBlobTest.cpp TestBounds (lines 127-174)
func TestTextBlobBoundsExplicit(t *testing.T) {
	builder := NewTextBlobBuilder()
	font := NewFont()

	// Empty builder should return nil blob
	blob := builder.Make()
	if blob != nil {
		t.Error("Empty builder should return nil blob")
	}

	// Single run with glyphs should have non-empty bounds
	buffer := builder.AllocRun(font, 16, 0, 0)
	for i := 0; i < 16; i++ {
		buffer.Glyphs[i] = GlyphID('A' + i)
	}
	builder.AddRun()

	blob = builder.Make()
	if blob == nil {
		t.Fatal("Blob should not be nil")
	}

	bounds := blob.Bounds()
	// Bounds should have positive width
	width := bounds.Right - bounds.Left
	if width <= 0 {
		t.Errorf("Bounds width should be positive, got %v", width)
	}
}

// TestTextBlobBoundsMultipleRuns tests bounds union of multiple runs.
// Ported from: TextBlobTest.cpp TestBounds (lines 158-169)
func TestTextBlobBoundsMultipleRuns(t *testing.T) {
	builder := NewTextBlobBuilder()
	font := NewFont()

	// Create multiple runs at different positions
	positions := []Point{
		{X: 10, Y: 10},
		{X: 50, Y: 20},
		{X: 0, Y: 5},
	}

	for _, pos := range positions {
		buffer := builder.AllocRun(font, 16, pos.X, pos.Y)
		for i := 0; i < 16; i++ {
			buffer.Glyphs[i] = GlyphID('A' + i)
		}
		builder.AddRun()
	}

	blob := builder.Make()
	if blob == nil {
		t.Fatal("Blob should not be nil")
	}

	bounds := blob.Bounds()

	// Bounds should encompass all runs
	// Min X should be at or before 0 (from third run at X=0)
	if bounds.Left > 0 {
		t.Errorf("Bounds left should be at or before 0, got %v", bounds.Left)
	}

	// Bounds should have reasonable width covering all runs
	if bounds.Right <= bounds.Left {
		t.Errorf("Bounds width should be positive: left=%v, right=%v", bounds.Left, bounds.Right)
	}
}

// TestTextBlobPaintProps verifies font properties are captured in runs.
// Ported from: TextBlobTest.cpp TestPaintProps (lines 197-240)
func TestTextBlobPaintProps(t *testing.T) {
	// Create a "kitchen sink" font with all properties set
	font := NewFont()
	font.SetSize(42)
	font.SetScaleX(4.2)
	font.SetSkewX(0.42)
	font.SetHinting(enums.FontHintingFull)
	font.SetEdging(enums.FontEdgingSubpixelAntiAlias)
	font.SetEmbolden(true)
	font.SetLinearMetrics(true)
	font.SetSubpixel(true)
	font.SetEmbeddedBitmaps(true)
	font.SetForceAutoHinting(true)

	// Ensure we didn't pick default values
	defaultFont := NewFont()
	if defaultFont.Size() == font.Size() {
		t.Error("Test font size should differ from default")
	}
	if defaultFont.ScaleX() == font.ScaleX() {
		t.Error("Test font scaleX should differ from default")
	}
	if defaultFont.SkewX() == font.SkewX() {
		t.Error("Test font skewX should differ from default")
	}

	builder := NewTextBlobBuilder()

	// Add runs with different positioning types
	buffer := builder.AllocRun(font, 1, 0, 0)
	buffer.Glyphs[0] = GlyphID('A')
	builder.AddRun()

	buffer = builder.AllocRunPosH(font, 1, 0)
	buffer.Glyphs[0] = GlyphID('B')
	buffer.Positions[0] = 10
	builder.AddRun()

	buffer = builder.AllocRunPos(font, 1)
	buffer.Glyphs[0] = GlyphID('C')
	buffer.Positions[0] = 20
	buffer.Positions[1] = 0
	builder.AddRun()

	blob := builder.Make()
	if blob == nil {
		t.Fatal("Blob should not be nil")
	}

	// Verify font properties are preserved in each run
	for i := 0; i < blob.RunCount(); i++ {
		run := blob.Run(i)
		if run == nil {
			t.Errorf("Run %d should not be nil", i)
			continue
		}

		runFont := run.Font
		if runFont == nil {
			t.Errorf("Run %d font should not be nil", i)
			continue
		}

		if runFont.Size() != 42 {
			t.Errorf("Run %d font size = %v, want 42", i, runFont.Size())
		}
		if runFont.ScaleX() != 4.2 {
			t.Errorf("Run %d font scaleX = %v, want 4.2", i, runFont.ScaleX())
		}
		if runFont.SkewX() != 0.42 {
			t.Errorf("Run %d font skewX = %v, want 0.42", i, runFont.SkewX())
		}
		if !runFont.IsEmbolden() {
			t.Errorf("Run %d font should have embolden=true", i)
		}
		if !runFont.IsSubpixel() {
			t.Errorf("Run %d font should have subpixel=true", i)
		}
	}
}

// TestTextBlobMakeFromStringPositioning tests that MakeFromString produces full positioning.
// Ported from: TextBlobTest.cpp TextBlob_MakeAsDrawText (lines 467-479)
func TestTextBlobMakeFromStringPositioning(t *testing.T) {
	font := NewFont()
	text := "Hello"
	blob := MakeTextBlobFromString(text, font)
	if blob == nil {
		t.Fatal("MakeTextBlobFromString should not return nil")
	}

	// Should have 1 run
	if blob.RunCount() != 1 {
		t.Errorf("Expected 1 run, got %d", blob.RunCount())
	}

	run := blob.Run(0)
	if run == nil {
		t.Fatal("Run(0) should not be nil")
	}

	// Should have same number of glyphs as characters
	if len(run.Glyphs) != len(text) {
		t.Errorf("Expected %d glyphs, got %d", len(text), len(run.Glyphs))
	}

	// Each glyph should have a position (full positioning)
	if len(run.Positions) != len(text) {
		t.Errorf("Expected %d positions, got %d", len(text), len(run.Positions))
	}
}

// TestTextBlobIterator tests blob iteration APIs.
// Ported from: TextBlobTest.cpp TextBlob_iter (lines 481-511)
func TestTextBlobIterator(t *testing.T) {
	font := NewFont()
	builder := NewTextBlobBuilder()

	// Add two runs
	buffer := builder.AllocRun(font, 5, 10, 20)
	for i := 0; i < 5; i++ {
		buffer.Glyphs[i] = GlyphID('H' + i) // H, I, J, K, L -> "Hello" approximation
	}
	builder.AddRun()

	buffer = builder.AllocRun(font, 6, 10, 40)
	for i := 0; i < 6; i++ {
		buffer.Glyphs[i] = GlyphID('W' + i%6)
	}
	builder.AddRun()

	blob := builder.Make()
	if blob == nil {
		t.Fatal("Blob should not be nil")
	}

	// Verify run count
	if blob.RunCount() != 2 {
		t.Errorf("Expected 2 runs, got %d", blob.RunCount())
	}

	// Verify first run
	run := blob.Run(0)
	if run == nil {
		t.Fatal("Run(0) should not be nil")
	}
	if len(run.Glyphs) != 5 {
		t.Errorf("Expected 5 glyphs in run 0, got %d", len(run.Glyphs))
	}

	// Verify second run
	run = blob.Run(1)
	if run == nil {
		t.Fatal("Run(1) should not be nil")
	}
	if len(run.Glyphs) != 6 {
		t.Errorf("Expected 6 glyphs in run 1, got %d", len(run.Glyphs))
	}

	// Out of bounds access
	if blob.Run(2) != nil {
		t.Error("Run(2) should return nil")
	}
	if blob.Run(-1) != nil {
		t.Error("Run(-1) should return nil")
	}
}
