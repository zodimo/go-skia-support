package impl

import (
	"sync/atomic"

	"github.com/go-text/typesetting/font"
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
	goTextFace *font.Face
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

// NewTypefaceWithTypefaceFace creates a new typeface with a go-text/typesetting Face.
func NewTypefaceWithTypefaceFace(familyName string, style FontStyle, face *font.Face) *Typeface {
	return &Typeface{
		style:      style,
		familyName: familyName,
		uniqueID:   nextTypefaceID(),
		fixedPitch: false,
		goTextFace: face,
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

// GoTextFace returns the underlying go-text/typesetting Face, if any.
func (t *Typeface) GoTextFace() *font.Face {
	return t.goTextFace
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

// UnicharToGlyph returns the glyph ID for the given Unicode character.
// This is a stub implementation - real implementation requires cmap table parsing.
// Returns 1 for all characters to indicate "supported" for MVP.
// Ported from: SkTypeface::unicharToGlyph
func (t *Typeface) UnicharToGlyph(unichar rune) uint16 {
	// Stub: assume all characters are supported
	// Real implementation would call onCharsToGlyphs with platform-specific cmap parsing
	return 1
}

// Compile-time interface check
var _ SkTypeface = (*Typeface)(nil)
