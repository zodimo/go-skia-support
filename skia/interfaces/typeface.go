package interfaces

import (
	"github.com/zodimo/go-skia-support/skia/models"
)

// SkTypeface represents the typeface and intrinsic style of a font.
// This is used in the font, along with optionally algorithmic settings like
// textSize, textSkewX, textScaleX, to specify how text appears when drawn.
//
// Typeface objects are immutable, and so they can be shared between threads.
//
// Ported from: skia-source/include/core/SkTypeface.h
type SkTypeface interface {
	// FontStyle returns the typeface's intrinsic style attributes.
	FontStyle() models.FontStyle

	// IsBold returns true if style has the bold bit set.
	IsBold() bool

	// IsItalic returns true if style has the italic bit set.
	IsItalic() bool

	// IsFixedPitch returns true if the typeface claims to be fixed-pitch.
	// This is a style bit, advance widths may vary even if this returns true.
	IsFixedPitch() bool

	// UniqueID returns a 32bit value unique for this typeface.
	// Will never return 0.
	UniqueID() uint32

	// FamilyName returns the family name for this typeface.
	// It will always be returned encoded as UTF8.
	FamilyName() string
}
