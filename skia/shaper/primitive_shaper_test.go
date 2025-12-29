package shaper

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/impl"
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// MockRunHandler captures calls for verification
type MockRunHandler struct {
	BeginLineCalled    bool
	CommitLineCalled   bool
	RunInfos           []RunInfo
	CommitRunInfoCount int
	Buffers            []Buffer
}

func (m *MockRunHandler) BeginLine() {
	m.BeginLineCalled = true
}

func (m *MockRunHandler) RunInfo(info RunInfo) {
	m.RunInfos = append(m.RunInfos, info)
}

func (m *MockRunHandler) CommitRunInfo() {
	m.CommitRunInfoCount++
}

func (m *MockRunHandler) RunBuffer(info RunInfo) Buffer {
	// Return a buffer with allocated slices to avoid nil pointer panic if shaper writes to it
	count := int(info.GlyphCount)
	buffer := Buffer{
		Glyphs:    make([]uint16, count),
		Positions: make([]models.Point, count),
		Clusters:  make([]uint32, count),
	}
	m.Buffers = append(m.Buffers, buffer)
	// Return pointer to the last appended buffer?
	// Slice append copies struct. We need to return the one that will be inspected.
	// But RunBuffer returns Buffer struct, not pointer.
	// So PrimitiveShaper writes to the returned Buffer struct which is a copy?
	// In C++, Buffer contains pointers to storage.
	// In Go, Buffer definition in handler.go has slices:
	// type Buffer struct { Glyphs []uint16, Positions []models.Point ... }
	// So copying Buffer struct copies slice headers. Writing to slices works.
	return buffer
}

func (m *MockRunHandler) CommitRunBuffer(info RunInfo) {
}

func (m *MockRunHandler) CommitLine() {
	m.CommitLineCalled = true
}

// MockIterators to simulate breaks
type MockIterator struct {
	breaks []int // Sorted breakpoints
	idx    int
	length int
}

func NewMockIterator(breaks []int, length int) *MockIterator {
	return &MockIterator{breaks: breaks, length: length}
}

func (m *MockIterator) Consume() {
	if m.idx < len(m.breaks) {
		m.idx++
	}
}

func (m *MockIterator) EndOfCurrentRun() int {
	if m.idx < len(m.breaks) {
		return m.breaks[m.idx]
	}
	return m.length
}

func (m *MockIterator) AtEnd() bool {
	return m.idx >= len(m.breaks) && (len(m.breaks) == 0 || m.breaks[len(m.breaks)-1] == m.length)
}

// Implement specific iterator types embedding MockIterator
type MockFontIterator struct{ *MockIterator }

func (m *MockFontIterator) CurrentFont() interfaces.SkFont { return impl.NewFont() }

type MockBiDiIterator struct{ *MockIterator }

func (m *MockBiDiIterator) CurrentLevel() uint8 { return 0 }

type MockScriptIterator struct{ *MockIterator }

func (m *MockScriptIterator) CurrentScript() uint32 { return 0 }

type MockLangIterator struct{ *MockIterator }

func (m *MockLangIterator) CurrentLanguage() string { return "" }

func TestPrimitiveShaper_Shape_LoopLogic(t *testing.T) {
	shaper := NewPrimitiveShaper()
	text := "0123456789"
	// length := len(text) // Unused variable removed

	// Case 1: All trivial (no breaks)
	// Expect 1 run covering 0-10
	handler := &MockRunHandler{}
	font := impl.NewFont()
	shaper.Shape(text, font, true, 100, handler, nil)

	if !handler.BeginLineCalled { // Changed from handler.beginLineCalled to handler.BeginLineCalled
		t.Error("BeginLine not called")
	}
	if !handler.CommitLineCalled {
		t.Error("CommitLine not called")
	}
}

func TestPrimitiveShaper_ShapeWithIterators_Breaks(t *testing.T) {
	shaper := NewPrimitiveShaper()
	text := "0123456789"
	length := len(text)

	// Custom iterators with breaks
	// Font breaks at 5
	fontIter := &MockFontIterator{NewMockIterator([]int{5, 10}, length)}
	// BiDi breaks at 3 and 8
	bidiIter := &MockBiDiIterator{NewMockIterator([]int{3, 8, 10}, length)}
	// Script no breaks
	scriptIter := &MockScriptIterator{NewMockIterator([]int{}, length)}
	// Lang no breaks
	langIter := &MockLangIterator{NewMockIterator([]int{}, length)}

	handler := &MockRunHandler{}

	shaper.ShapeWithIterators(text, fontIter, bidiIter, scriptIter, langIter, nil, 100, handler)

	if !handler.BeginLineCalled { // Changed from handler.beginLineCalled to handler.BeginLineCalled
		t.Error("BeginLine not called")
	}
	if !handler.CommitLineCalled {
		t.Error("CommitLine not called")
	}

	// Verify iterators are consumed
	if !fontIter.AtEnd() {
		if fontIter.EndOfCurrentRun() != length {
			t.Errorf("Font iterator not at end, current end: %d", fontIter.EndOfCurrentRun())
		}
	}

	// Verify runs
	// Expected runs:
	// Run 1: 0-3 (limited by BiDi)
	// Run 2: 3-5 (limited by Font)
	// Run 3: 5-8 (limited by BiDi)
	// Run 4: 8-10 (limited by End)

	expectedRanges := []struct{ begin, end int }{
		{0, 3},
		{3, 5},
		{5, 8},
		{8, 10},
	}

	if len(handler.RunInfos) != len(expectedRanges) {
		t.Fatalf("Expected %d runs, got %d", len(expectedRanges), len(handler.RunInfos))
	}

	for i, exp := range expectedRanges {
		got := handler.RunInfos[i].Utf8Range
		if got.Begin != exp.begin || got.End != exp.end {
			t.Errorf("Run %d mismatch: expected [%d, %d), got [%d, %d)", i, exp.begin, exp.end, got.Begin, got.End)
		}
	}
}

// MockFont for testing
type MockFont struct {
	interfaces.SkFont
}

func (m *MockFont) UnicharToGlyph(unichar rune) uint16 {
	return uint16(unichar) // Identity mapping for test
}

func (m *MockFont) GetWidths(glyphs []uint16) []models.Scalar {
	widths := make([]models.Scalar, len(glyphs))
	for i := range glyphs {
		widths[i] = 10.0 // Constant width
	}
	return widths
}

func TestPrimitiveShaper_Shape_SimpleRun(t *testing.T) {
	shaper := NewPrimitiveShaper()
	text := "ABC"
	font := &MockFont{}
	handler := &MockRunHandler{}

	shaper.Shape(text, font, true, 100, handler, nil)

	if len(handler.RunInfos) != 1 {
		t.Fatalf("Expected 1 run, got %d", len(handler.RunInfos))
	}

	info := handler.RunInfos[0]
	if info.GlyphCount != 3 {
		t.Errorf("Expected 3 glyphs, got %d", info.GlyphCount)
	}
	if info.Advance.X != 30.0 {
		t.Errorf("Expected advance 30.0, got %f", info.Advance.X)
	}
}
