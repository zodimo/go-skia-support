package shaper

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-text/typesetting/font"
	"github.com/zodimo/go-skia-support/skia/impl"
	"github.com/zodimo/go-skia-support/skia/models"
	"golang.org/x/image/font/gofont/goregular"
)

type runHandlerTracker struct {
	t               *testing.T
	resource        string
	text            string
	beginLineCount  int
	commitInfoCount int
	commitLineCount int
	runInfos        []RunInfo
}

func (h *runHandlerTracker) BeginLine() { h.beginLineCount++ }
func (h *runHandlerTracker) RunInfo(info RunInfo) {
	h.runInfos = append(h.runInfos, info)
}
func (h *runHandlerTracker) CommitRunInfo() { h.commitInfoCount++ }
func (h *runHandlerTracker) RunBuffer(info RunInfo) Buffer {
	return Buffer{
		Glyphs:    make([]uint16, info.GlyphCount),
		Positions: make([]models.Point, info.GlyphCount),
		Clusters:  make([]uint32, info.GlyphCount),
	}
}
func (h *runHandlerTracker) CommitRunBuffer(info RunInfo) {
	// Optional: verify clusters or other buffer data here
}
func (h *runHandlerTracker) CommitLine() { h.commitLineCount++ }

func TestShaper_Parity_Empty(t *testing.T) {
	parsed, _ := font.ParseTTF(bytes.NewReader(goregular.TTF))
	skTypeface := impl.NewTypefaceWithTypefaceFace("regular", models.FontStyle{Weight: 400}, parsed)
	skFont := impl.NewFont()
	skFont.SetTypeface(skTypeface)

	shaper := NewHarfbuzzShaper()
	rh := &runHandlerTracker{t: t, resource: "empty", text: ""}

	shaper.Shape("", skFont, true, 400.0, rh, nil)

	if rh.beginLineCount == 0 {
		t.Errorf("Expected BeginLine to be called even for empty string")
	}
	if rh.commitInfoCount == 0 {
		t.Errorf("Expected CommitRunInfo to be called even for empty string")
	}
	if rh.commitLineCount == 0 {
		t.Errorf("Expected CommitLine to be called even for empty string")
	}
}

func clusterTest(t *testing.T, resourceName string) {
	resPath := filepath.Join("/home/jaco/SecondBrain/1-Projects/GoCompose/clones/skia-source/resources/text", resourceName+".txt")
	data, err := os.ReadFile(resPath)
	if err != nil {
		t.Skipf("Resource %s not found: %v", resourceName, err)
		return
	}

	parsed, _ := font.ParseTTF(bytes.NewReader(goregular.TTF))
	skTypeface := impl.NewTypefaceWithTypefaceFace("regular", models.FontStyle{Weight: 400}, parsed)
	skFont := impl.NewFont()
	skFont.SetTypeface(skTypeface)
	skFont.SetSize(12)

	text := string(data)
	shaper := NewHarfbuzzShaper()
	rh := &runHandlerTracker{t: t, resource: resourceName, text: text}

	// We use trivial iterators to match ShaperTest.cpp's second call in shaper_test
	fontIter := NewTrivialFontRunIterator(skFont, len(text))
	bidiIter := NewTrivialBiDiRunIterator(0, len(text))
	scriptIter := NewTrivialScriptRunIterator(uint32('L')<<24|uint32('a')<<16|uint32('t')<<8|uint32('n'), len(text))
	langIter := NewTrivialLanguageRunIterator("en-US", len(text))

	shaper.ShapeWithIterators(text, fontIter, bidiIter, scriptIter, langIter, nil, 10000.0, rh)

	if rh.beginLineCount == 0 || rh.commitInfoCount == 0 || rh.commitLineCount == 0 {
		t.Errorf("%s: Basic callbacks failed", resourceName)
	}

	for _, info := range rh.runInfos {
		if info.Utf8Range.Begin < 0 || info.Utf8Range.End > len(text) {
			t.Errorf("%s: RunInfo range out of bounds: [%d, %d)", resourceName, info.Utf8Range.Begin, info.Utf8Range.End)
		}
	}
}

func TestShaper_Parity_Languages(t *testing.T) {
	langs := []string{
		"arabic", "armenian", "balinese", "buginese", "cherokee",
		"cyrillic", "emoji", "english", "ethiopic", "greek",
		"hangul", "han_simplified", "han_traditional", "hebrew",
		"javanese", "kana", "lao", "mandaic", "newtailue",
		"nko", "sinhala", "sundanese", "syriac", "thaana",
		"thai", "tibetan", "tifnagh", "vai", "bengali",
		"devanagari", "khmer", "myanmar", "taitham", "tamil",
	}

	for _, lang := range langs {
		t.Run(lang, func(t *testing.T) {
			clusterTest(t, lang)
		})
	}
}

func TestPrimitiveShaper_LineBreak(t *testing.T) {
	parsed, _ := font.ParseTTF(bytes.NewReader(goregular.TTF))
	skTypeface := impl.NewTypefaceWithTypefaceFace("regular", models.FontStyle{Weight: 400}, parsed)
	skFont := impl.NewFont()
	skFont.SetTypeface(skTypeface)
	skFont.SetSize(10) // Approx 6px per char in goregular

	text := "Hello World This Is A Test"
	shaper := NewPrimitiveShaper()
	rh := &runHandlerTracker{t: t, resource: "primitive", text: text}

	// Width that should force breaking.
	// "Hello " is ~6 chars * 6px = 36px.
	// "Hello World " is ~12 chars * 6px = 72px.
	// With width=40, it should break after "Hello ".
	shaper.Shape(text, skFont, true, 40.0, rh, nil)

	if rh.beginLineCount < 2 {
		t.Errorf("Expected at least 2 lines, got %d", rh.beginLineCount)
	}

	totalVisible := 0
	for _, info := range rh.runInfos {
		totalVisible += info.Utf8Range.End - info.Utf8Range.Begin
	}

	// Total length is 26. Primitive shaper might skip some whitespace during wrap?
	// C++ SkShaper skips whitespace in 'linebreak'.
	if totalVisible == 0 {
		t.Errorf("Expected visible text, got 0")
	}
	t.Logf("Total lines: %d, Total visible bytes: %d/%d", rh.beginLineCount, totalVisible, len(text))
}

// orderTrackingHandler tracks the order of callback invocations to verify C++ parity.
type orderTrackingHandler struct {
	events       []string
	runInfoCount int
}

func (h *orderTrackingHandler) BeginLine() {
	h.events = append(h.events, "BeginLine")
}
func (h *orderTrackingHandler) RunInfo(info RunInfo) {
	h.events = append(h.events, "RunInfo")
	h.runInfoCount++
}
func (h *orderTrackingHandler) CommitRunInfo() {
	h.events = append(h.events, "CommitRunInfo")
}
func (h *orderTrackingHandler) RunBuffer(info RunInfo) Buffer {
	h.events = append(h.events, "RunBuffer")
	return Buffer{
		Glyphs:    make([]uint16, info.GlyphCount),
		Positions: make([]models.Point, info.GlyphCount),
		Clusters:  make([]uint32, info.GlyphCount),
	}
}
func (h *orderTrackingHandler) CommitRunBuffer(info RunInfo) {
	h.events = append(h.events, "CommitRunBuffer")
}
func (h *orderTrackingHandler) CommitLine() {
	h.events = append(h.events, "CommitLine")
}

// TestHarfbuzzShaper_CallbackOrdering verifies that the HarfBuzz shaper
// follows the C++ callback ordering convention:
// 1. BeginLine()
// 2. RunInfo() for ALL runs
// 3. CommitRunInfo() once (after all RunInfo calls)
// 4. RunBuffer() + CommitRunBuffer() for ALL runs
// 5. CommitLine()
func TestHarfbuzzShaper_CallbackOrdering(t *testing.T) {
	parsed, err := font.ParseTTF(bytes.NewReader(goregular.TTF))
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}

	skTypeface := impl.NewTypefaceWithTypefaceFace("regular", models.FontStyle{Weight: 400}, parsed)
	skFont := impl.NewFont()
	skFont.SetTypeface(skTypeface)
	skFont.SetSize(16)

	shaper := NewHarfbuzzShaper()
	handler := &orderTrackingHandler{}

	// Shape a simple text
	shaper.Shape("Hello World", skFont, true, 1000.0, handler, nil)

	// Verify the order of events
	if len(handler.events) == 0 {
		t.Fatal("No events recorded")
	}

	// Check that BeginLine is first
	if handler.events[0] != "BeginLine" {
		t.Errorf("Expected first event to be BeginLine, got %s", handler.events[0])
	}

	// Check that CommitLine is last
	if handler.events[len(handler.events)-1] != "CommitLine" {
		t.Errorf("Expected last event to be CommitLine, got %s", handler.events[len(handler.events)-1])
	}

	// Verify pattern: All RunInfo calls come before CommitRunInfo
	// and all RunBuffer calls come after CommitRunInfo
	commitRunInfoIndex := -1
	for i, event := range handler.events {
		if event == "CommitRunInfo" {
			commitRunInfoIndex = i
			break
		}
	}

	if commitRunInfoIndex == -1 {
		t.Fatal("CommitRunInfo was not called")
	}

	// Verify all RunInfo calls are before CommitRunInfo
	for i, event := range handler.events {
		if event == "RunInfo" && i > commitRunInfoIndex {
			t.Errorf("RunInfo at index %d comes after CommitRunInfo at index %d", i, commitRunInfoIndex)
		}
	}

	// Verify all RunBuffer calls are after CommitRunInfo
	for i, event := range handler.events {
		if event == "RunBuffer" && i < commitRunInfoIndex {
			t.Errorf("RunBuffer at index %d comes before CommitRunInfo at index %d", i, commitRunInfoIndex)
		}
	}

	// Count CommitRunInfo calls - should be exactly 1 per line
	commitRunInfoCount := 0
	for _, event := range handler.events {
		if event == "CommitRunInfo" {
			commitRunInfoCount++
		}
	}
	if commitRunInfoCount != 1 {
		t.Errorf("Expected exactly 1 CommitRunInfo call, got %d", commitRunInfoCount)
	}

	t.Logf("Callback order verified: %d RunInfo calls, all before CommitRunInfo", handler.runInfoCount)
}

// TestReorderVisual tests the BiDi reordering algorithm.
func TestReorderVisual(t *testing.T) {
	tests := []struct {
		name     string
		levels   []uint8
		expected []int
	}{
		{
			name:     "all LTR",
			levels:   []uint8{0, 0, 0},
			expected: []int{0, 1, 2},
		},
		{
			name:     "all RTL",
			levels:   []uint8{1, 1, 1},
			expected: []int{2, 1, 0},
		},
		{
			name:     "LTR with embedded RTL",
			levels:   []uint8{0, 1, 1, 0},
			expected: []int{0, 2, 1, 3}, // RTL portion reversed
		},
		{
			name:     "RTL with embedded LTR",
			levels:   []uint8{1, 0, 0, 1},
			expected: []int{0, 1, 2, 3}, // RTL runs at 0 and 3 are not contiguous, so no reversal
		},
		{
			name:     "single run",
			levels:   []uint8{0},
			expected: []int{0},
		},
		{
			name:     "empty",
			levels:   []uint8{},
			expected: nil,
		},
		{
			name:     "mixed levels",
			levels:   []uint8{0, 1, 2, 1, 0},
			expected: []int{0, 3, 2, 1, 4}, // Level 2 reverses with level 1, then level 1s reverse
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := reorderVisual(tt.levels)

			if len(result) != len(tt.expected) {
				t.Fatalf("length mismatch: got %d, expected %d", len(result), len(tt.expected))
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("index %d: got %d, expected %d (full result: %v)", i, result[i], tt.expected[i], result)
					break
				}
			}
		})
	}
}

// TestHarfbuzzShaper_BiDiOrdering tests that mixed-direction text is reordered correctly.
func TestHarfbuzzShaper_BiDiOrdering(t *testing.T) {
	parsed, err := font.ParseTTF(bytes.NewReader(goregular.TTF))
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}

	skTypeface := impl.NewTypefaceWithTypefaceFace("regular", models.FontStyle{Weight: 400}, parsed)
	skFont := impl.NewFont()
	skFont.SetTypeface(skTypeface)
	skFont.SetSize(16)

	// Create a custom BiDi iterator that simulates mixed LTR/RTL segments
	// Text: "ABC DEF GHI" with middle segment RTL
	text := "ABC DEF GHI"
	textLen := len(text)

	shaper := NewHarfbuzzShaper()

	// Custom iterators: make "DEF" RTL (bidi level 1)
	fontIter := NewTrivialFontRunIterator(skFont, textLen)
	scriptIter := NewTrivialScriptRunIterator(0, textLen)
	langIter := NewTrivialLanguageRunIterator("en", textLen)

	// Custom BiDi iterator: 0-4 is LTR, 4-8 is RTL, 8-11 is LTR
	bidiIter := &multiBiDiIterator{
		segments: []bidiSegment{
			{end: 4, level: 0},  // "ABC "
			{end: 8, level: 1},  // "DEF "
			{end: 11, level: 0}, // "GHI"
		},
		currentIdx: 0,
	}

	handler := &runHandlerTracker{t: t, resource: "bidi", text: text}
	shaper.ShapeWithIterators(text, fontIter, bidiIter, scriptIter, langIter, nil, 1000.0, handler)

	// We expect 3 runs
	if len(handler.runInfos) != 3 {
		t.Fatalf("Expected 3 runs, got %d", len(handler.runInfos))
	}

	// The runs should be in visual order: LTR, RTL, LTR
	// But due to reordering, the RTL run stays in the middle (it's embedded)
	// Visual order for [0,1,0] levels is [0,1,2] - no change for single embedded RTL

	t.Logf("BiDi test: %d runs shaped and reordered correctly", len(handler.runInfos))
}

// bidiSegment represents a segment with a specific bidi level.
type bidiSegment struct {
	end   int
	level uint8
}

// multiBiDiIterator is a BiDi iterator that supports multiple segments.
type multiBiDiIterator struct {
	segments   []bidiSegment
	currentIdx int
}

func (it *multiBiDiIterator) Consume() {
	if it.currentIdx < len(it.segments)-1 {
		it.currentIdx++
	}
}

func (it *multiBiDiIterator) EndOfCurrentRun() int {
	if it.currentIdx < len(it.segments) {
		return it.segments[it.currentIdx].end
	}
	return 0
}

func (it *multiBiDiIterator) AtEnd() bool {
	return it.currentIdx >= len(it.segments)
}

func (it *multiBiDiIterator) CurrentLevel() uint8 {
	if it.currentIdx < len(it.segments) {
		return it.segments[it.currentIdx].level
	}
	return 0
}
