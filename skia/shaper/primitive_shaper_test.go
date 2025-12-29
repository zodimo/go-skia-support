package shaper

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/interfaces"
)

// MockRunHandler captures calls for verification
type MockRunHandler struct {
	BeginLineCalled    bool
	CommitLineCalled   bool
	RunInfos           []RunInfo
	CommitRunInfoCount int
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
	return Buffer{}
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

func (m *MockFontIterator) CurrentFont() interfaces.SkFont { return nil }

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
	shaper.Shape(text, nil, true, 100, handler)

	if !handler.BeginLineCalled {
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

	shaper.ShapeWithIterators(text, fontIter, bidiIter, scriptIter, langIter, 100, handler)

	if !handler.BeginLineCalled {
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
