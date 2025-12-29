package paragraph

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// MockTextLineOwner mocks the TextLineOwner interface manually.
type MockTextLineOwner struct {
	Blocks     []Block
	RunMap     map[int]*Run
	ClusterMap map[int]*Cluster
	PStyle     ParagraphStyle
}

func (m *MockTextLineOwner) Styles() []Block {
	return m.Blocks
}

func (m *MockTextLineOwner) Run(index int) *Run {
	return m.RunMap[index]
}

func (m *MockTextLineOwner) Cluster(index int) *Cluster {
	return m.ClusterMap[index]
}

func (m *MockTextLineOwner) Block(index int) Block {
	if index < 0 || index >= len(m.Blocks) {
		return Block{}
	}
	return m.Blocks[index]
}

func (m *MockTextLineOwner) GetUnicode() interfaces.SkUnicode {
	return nil
}

func (m *MockTextLineOwner) ParagraphStyle() ParagraphStyle {
	return m.PStyle
}

func (m *MockTextLineOwner) FontCollection() *FontCollection {
	return nil
}

func TestNewTextLine(t *testing.T) {
	owner := &MockTextLineOwner{
		Blocks: []Block{{
			Range: NewTextRange(0, 10),
			Style: NewTextStyle(),
		}},
	}

	offset := models.Point{X: 10, Y: 20}
	advance := models.Point{X: 100, Y: 14}
	metrics := NewInternalLineMetrics()

	blocks := NewBlockRange(0, 1)

	tl := NewTextLine(
		owner,
		offset,
		advance,
		blocks,
		NewTextRange(0, 10),
		NewTextRange(0, 10),
		NewTextRange(0, 10),
		Range[int]{Start: 0, End: 0},
		Range[int]{Start: 0, End: 0},
		100.0,
		metrics,
	)

	if tl == nil {
		t.Fatal("NewTextLine returned nil")
	}
	if tl.Width() != 100.0 {
		t.Errorf("Expected width 100.0, got %f", tl.Width())
	}
	if tl.Height() != 14.0 {
		t.Errorf("Expected height 14.0, got %f", tl.Height())
	}
}

func TestTextLineFormat(t *testing.T) {
	style := NewTextStyle()
	owner := &MockTextLineOwner{
		Blocks: []Block{{
			Range: NewTextRange(0, 10),
			Style: style,
		}},
		PStyle: ParagraphStyle{TextDirection: TextDirectionLTR},
	}

	offset := models.Point{X: 0, Y: 0}
	advance := models.Point{X: 50, Y: 14}
	metrics := NewInternalLineMetrics()

	tl := NewTextLine(
		owner,
		offset,
		advance, // width 50
		NewBlockRange(0, 1),
		NewTextRange(0, 10),
		NewTextRange(0, 10),
		NewTextRange(0, 10),
		Range[int]{Start: 0, End: 0},
		Range[int]{Start: 0, End: 0},
		50.0,
		metrics,
	)

	// Test Right Alignment
	// MaxWidth = 100. Line Width = 50. Delta = 50.
	// Right align should set shift to 50.
	tl.Format(TextAlignRight, 100.0)

	if tl.shift != 50.0 {
		t.Errorf("Expected shift 50.0, got %f", tl.shift)
	}

	// Test Center Alignment
	tl.shift = 0
	tl.Format(TextAlignCenter, 100.0)
	if tl.shift != 25.0 {
		t.Errorf("Expected shift 25.0, got %f", tl.shift)
	}
}
