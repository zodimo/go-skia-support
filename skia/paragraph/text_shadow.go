package paragraph

import (
	"github.com/zodimo/go-skia-support/skia/interfaces"
)

// TextShadow represents a shadow effect applied to text.
// Multiple shadows can be applied to the same text.
//
// Ported from: skia-source/modules/skparagraph/include/TextShadow.h
type TextShadow struct {
	// Color is the shadow color. Default is black.
	Color uint32

	// Offset is the shadow offset from the text.
	Offset interfaces.Point

	// BlurSigma is the blur radius as a sigma value. 0 means no blur.
	BlurSigma float64
}

// NewTextShadow creates a new TextShadow with the given parameters.
func NewTextShadow(color uint32, offset interfaces.Point, blurSigma float64) TextShadow {
	return TextShadow{
		Color:     color,
		Offset:    offset,
		BlurSigma: blurSigma,
	}
}

// NewTextShadowDefault creates a TextShadow with default values (black, no offset, no blur).
func NewTextShadowDefault() TextShadow {
	return TextShadow{
		Color:     0xFF000000, // SK_ColorBLACK
		Offset:    interfaces.Point{},
		BlurSigma: 0.0,
	}
}

// HasShadow returns true if this shadow has any visible effect.
// A shadow has no effect if it has zero blur and zero offset.
func (s TextShadow) HasShadow() bool {
	return s.BlurSigma > 0 || s.Offset.X != 0 || s.Offset.Y != 0
}

// Equals returns true if this shadow equals another shadow.
func (s TextShadow) Equals(other TextShadow) bool {
	return s.Color == other.Color &&
		s.Offset.X == other.Offset.X &&
		s.Offset.Y == other.Offset.Y &&
		s.BlurSigma == other.BlurSigma
}
