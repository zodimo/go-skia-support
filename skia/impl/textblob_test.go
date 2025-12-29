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
