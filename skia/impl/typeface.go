package impl

import (
	"sync/atomic"
)

// Global unique ID counter for typefaces
var typefaceIDCounter uint32

// nextTypefaceID generates a unique ID for a typeface.
func nextTypefaceID() uint32 {
	return atomic.AddUint32(&typefaceIDCounter, 1)
}

// Typeface represents the typeface and intrinsic style of a font.
// This is a basic implementation for MVP; can be extended for system font integration.
//
// Ported from: skia-source/include/core/SkTypeface.h
type Typeface struct {
	style      FontStyle
	familyName string
	uniqueID   uint32
	fixedPitch bool
}

// NewDefaultTypeface creates a new typeface with default style.
func NewDefaultTypeface() *Typeface {
	return &Typeface{
		style:      FontStyle{Weight: 400, Width: 5, Slant: 0}, // Normal
		familyName: "",
		uniqueID:   nextTypefaceID(),
		fixedPitch: false,
	}
}

// NewTypeface creates a new typeface with the given family name and style.
func NewTypeface(familyName string, style FontStyle) *Typeface {
	return &Typeface{
		style:      style,
		familyName: familyName,
		uniqueID:   nextTypefaceID(),
		fixedPitch: false,
	}
}

// NewTypefaceWithOptions creates a new typeface with all options.
func NewTypefaceWithOptions(familyName string, style FontStyle, fixedPitch bool) *Typeface {
	return &Typeface{
		style:      style,
		familyName: familyName,
		uniqueID:   nextTypefaceID(),
		fixedPitch: fixedPitch,
	}
}

// FontStyle returns the typeface's intrinsic style attributes.
func (t *Typeface) FontStyle() FontStyle {
	return t.style
}

// IsBold returns true if style has the bold bit set.
func (t *Typeface) IsBold() bool {
	return t.style.IsBold()
}

// IsItalic returns true if style has the italic bit set.
func (t *Typeface) IsItalic() bool {
	return t.style.IsItalic()
}

// IsFixedPitch returns true if the typeface claims to be fixed-pitch.
func (t *Typeface) IsFixedPitch() bool {
	return t.fixedPitch
}

// UniqueID returns a 32bit value unique for this typeface.
func (t *Typeface) UniqueID() uint32 {
	return t.uniqueID
}

// FamilyName returns the family name for this typeface.
func (t *Typeface) FamilyName() string {
	return t.familyName
}

// Compile-time interface check
var _ SkTypeface = (*Typeface)(nil)
