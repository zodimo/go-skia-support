package paragraph

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/interfaces"
)

func TestNewPositionWithAffinity(t *testing.T) {
	tests := []struct {
		name         string
		position     int32
		affinity     Affinity
		wantPosition int32
		wantAffinity Affinity
	}{
		{"upstream at 0", 0, AffinityUpstream, 0, AffinityUpstream},
		{"downstream at 10", 10, AffinityDownstream, 10, AffinityDownstream},
		{"negative position", -5, AffinityUpstream, -5, AffinityUpstream},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPositionWithAffinity(tt.position, tt.affinity)
			if p.Position != tt.wantPosition {
				t.Errorf("Position = %d, want %d", p.Position, tt.wantPosition)
			}
			if p.Affinity != tt.wantAffinity {
				t.Errorf("Affinity = %d, want %d", p.Affinity, tt.wantAffinity)
			}
		})
	}
}

func TestNewPositionWithAffinityDefault(t *testing.T) {
	p := NewPositionWithAffinityDefault()
	if p.Position != 0 {
		t.Errorf("Default Position = %d, want 0", p.Position)
	}
	if p.Affinity != AffinityDownstream {
		t.Errorf("Default Affinity = %d, want Downstream (%d)", p.Affinity, AffinityDownstream)
	}
}

func TestNewTextBox(t *testing.T) {
	rect := interfaces.Rect{Left: 10, Top: 20, Right: 100, Bottom: 50}

	// Test LTR direction
	boxLTR := NewTextBox(rect, TextDirectionLTR)
	if boxLTR.Rect != rect {
		t.Errorf("TextBox.Rect = %+v, want %+v", boxLTR.Rect, rect)
	}
	if boxLTR.Direction != TextDirectionLTR {
		t.Errorf("TextBox.Direction = %d, want LTR (%d)", boxLTR.Direction, TextDirectionLTR)
	}

	// Test RTL direction
	boxRTL := NewTextBox(rect, TextDirectionRTL)
	if boxRTL.Direction != TextDirectionRTL {
		t.Errorf("TextBox.Direction = %d, want RTL (%d)", boxRTL.Direction, TextDirectionRTL)
	}
}

func TestTextBoxWithZeroRect(t *testing.T) {
	rect := interfaces.Rect{}
	box := NewTextBox(rect, TextDirectionLTR)

	if box.Rect.Left != 0 || box.Rect.Top != 0 || box.Rect.Right != 0 || box.Rect.Bottom != 0 {
		t.Errorf("Zero rect not preserved: %+v", box.Rect)
	}
}
