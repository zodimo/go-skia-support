package paragraph

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// MockTextWrapperOwner mocks the TextWrapperOwner interface.
type MockTextWrapperOwner struct {
	clusters       []*Cluster
	blocks         []Block
	runs           map[int]*Run
	paragraphStyle ParagraphStyle
	emptyMetrics   InternalLineMetrics
	strutMetrics   InternalLineMetrics
	strutEnabled   bool
	strutForce     bool
	roundingHack   bool
	lines          []*TextLine
	text           string
}

func (m *MockTextWrapperOwner) Styles() []Block {
	return m.blocks
}

func (m *MockTextWrapperOwner) Run(index int) *Run {
	return m.runs[index]
}

func (m *MockTextWrapperOwner) Cluster(index int) *Cluster {
	if index < 0 || index >= len(m.clusters) {
		return nil
	}
	return m.clusters[index]
}

func (m *MockTextWrapperOwner) Block(index int) Block {
	if index < 0 || index >= len(m.blocks) {
		return Block{}
	}
	return m.blocks[index]
}

func (m *MockTextWrapperOwner) GetUnicode() interfaces.SkUnicode {
	return nil
}

func (m *MockTextWrapperOwner) ParagraphStyle() ParagraphStyle {
	return m.paragraphStyle
}

func (m *MockTextWrapperOwner) FontCollection() *FontCollection {
	return nil
}

func (m *MockTextWrapperOwner) Clusters() []*Cluster {
	return m.clusters
}

func (m *MockTextWrapperOwner) StructForceHeight() bool {
	return m.strutForce
}

func (m *MockTextWrapperOwner) GetApplyRoundingHack() bool {
	return m.roundingHack
}

func (m *MockTextWrapperOwner) GetEmptyMetrics() InternalLineMetrics {
	return m.emptyMetrics
}

func (m *MockTextWrapperOwner) StrutEnabled() bool {
	return m.strutEnabled
}

func (m *MockTextWrapperOwner) StrutMetrics() InternalLineMetrics {
	return m.strutMetrics
}

func (m *MockTextWrapperOwner) Lines() []*TextLine {
	return m.lines
}

func (m *MockTextWrapperOwner) Text() string {
	return m.text
}

func (m *MockTextWrapperOwner) GetText() string {
	return m.text
}

func TestNewTextWrapper(t *testing.T) {
	tw := NewTextWrapper()
	if tw == nil {
		t.Fatal("NewTextWrapper returned nil")
	}
	if tw.lineNumber != 1 {
		t.Errorf("Expected lineNumber 1, got %d", tw.lineNumber)
	}
}

func TestTextWrapperEmptyClusters(t *testing.T) {
	tw := NewTextWrapper()
	owner := &MockTextWrapperOwner{
		clusters: []*Cluster{},
		text:     "",
	}

	lineCount := 0
	tw.BreakTextIntoLines(owner, 100, func(
		textExcludingSpaces TextRange,
		text TextRange,
		textIncludingNewlines TextRange,
		clusters ClusterRange,
		clustersWithGhosts ClusterRange,
		widthWithSpaces float32,
		startClip, endClip int,
		offset, advance models.Point,
		metrics InternalLineMetrics,
		addEllipsis bool,
	) {
		lineCount++
	})

	if lineCount != 0 {
		t.Errorf("Expected 0 lines for empty clusters, got %d", lineCount)
	}
}

func TestLineBreakerWithLittleRounding(t *testing.T) {
	breaker := NewLineBreakerWithLittleRounding(100.0, false)

	// Width clearly below should not break
	if breaker.BreakLine(90.0) {
		t.Error("Should not break at 90.0 for maxWidth 100.0")
	}

	// Width clearly above should break
	if !breaker.BreakLine(110.0) {
		t.Error("Should break at 110.0 for maxWidth 100.0")
	}

	// Width near max with rounding hack off
	if breaker.BreakLine(99.8) {
		t.Error("Should not break at 99.8 for maxWidth 100.0 (rounding hack off)")
	}
}

func TestTextStretchMethods(t *testing.T) {
	ts := NewTextStretch()
	if !ts.Empty() {
		t.Error("New TextStretch should be empty")
	}

	ts.Clean()
	if ts.Width() != 0 {
		t.Error("Cleaned TextStretch width should be 0")
	}
}
