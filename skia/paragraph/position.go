package paragraph

import (
	"github.com/zodimo/go-skia-support/skia/interfaces"
)

// PositionWithAffinity represents a text position along with information about
// which direction the position is associated with (for cases where a position
// falls between two runs with different directions).
//
// Ported from: skia-source/modules/skparagraph/include/DartTypes.h
type PositionWithAffinity struct {
	// Position is the code unit index in the text.
	Position int32

	// Affinity indicates whether this position is associated with the
	// preceding (Upstream) or following (Downstream) run.
	Affinity Affinity
}

// NewPositionWithAffinity creates a new PositionWithAffinity with the given values.
func NewPositionWithAffinity(position int32, affinity Affinity) PositionWithAffinity {
	return PositionWithAffinity{
		Position: position,
		Affinity: affinity,
	}
}

// NewPositionWithAffinityDefault creates a new PositionWithAffinity at position 0
// with downstream affinity (the default).
func NewPositionWithAffinityDefault() PositionWithAffinity {
	return PositionWithAffinity{
		Position: 0,
		Affinity: AffinityDownstream,
	}
}

// TextBox represents a rectangle with an associated text direction.
// It's used when querying text rects for ranges - the direction indicates
// whether the text in that box flows left-to-right or right-to-left.
//
// Ported from: skia-source/modules/skparagraph/include/DartTypes.h
type TextBox struct {
	// Rect is the bounding rectangle for this text box.
	Rect interfaces.Rect

	// Direction indicates the text direction for this box.
	Direction TextDirection
}

// NewTextBox creates a new TextBox with the given rect and direction.
func NewTextBox(rect interfaces.Rect, direction TextDirection) TextBox {
	return TextBox{
		Rect:      rect,
		Direction: direction,
	}
}
