package impl

import (
	"errors"
	"sync/atomic"

	"github.com/go-text/typesetting/font"
	ot "github.com/go-text/typesetting/font/opentype"
	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
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
// Delegates to the typeface, matching C++ SkFont::unicharToGlyph.
// Ported from: SkTypeface::unicharToGlyph
func (t *Typeface) UnicharToGlyph(unichar rune) uint16 {
	if t.goTextFace != nil {
		gid, ok := t.goTextFace.NominalGlyph(unichar)
		if ok {
			return uint16(gid)
		}
	}
	// Fallback or not found
	return 0
}

// MakeClone returns a new typeface with the specified arguments.
func (t *Typeface) MakeClone(args models.FontArguments) interfaces.SkTypeface {
	// Clone basic fields
	newTf := &Typeface{
		style:      t.style,
		familyName: t.familyName,
		uniqueID:   nextTypefaceID(),
		fixedPitch: t.fixedPitch,
	}

	if t.goTextFace != nil {
		// Create new face from the underlying Font to support thread safety and isolation
		// Accessing embedded Font field
		newFace := font.NewFace(t.goTextFace.Font)

		// Apply variations
		var vars []font.Variation
		for _, coord := range args.VariationDesignPosition.Coordinates {
			vars = append(vars, font.Variation{
				Tag:   font.Tag(coord.Axis),
				Value: coord.Value,
			})
		}
		newFace.SetVariations(vars)
		newTf.goTextFace = newFace
	}

	return newTf
}

// --- Glyph Data Access Methods ---
// Required for Font Utilities Port (font-utilities-port.md)

// UnitsPerEm returns the units-per-em value for this typeface.
// Returns 0 if there is an error or no font face is available.
// Ported from: SkTypeface::getUnitsPerEm
func (t *Typeface) UnitsPerEm() int {
	if t.goTextFace != nil {
		return int(t.goTextFace.Upem())
	}
	return 0
}

// GetGlyphAdvance returns the horizontal advance for a glyph in font units.
// This is the raw value from the font tables, not scaled by font size.
func (t *Typeface) GetGlyphAdvance(glyphID uint16) int16 {
	if t.goTextFace != nil {
		return int16(t.goTextFace.HorizontalAdvance(font.GID(glyphID)))
	}
	return 0
}

// GetGlyphBounds returns the bounding box for a glyph in font units.
// This is the raw value from the font tables, not scaled by font size.
func (t *Typeface) GetGlyphBounds(glyphID uint16) interfaces.Rect {
	if t.goTextFace != nil {
		extents, ok := t.goTextFace.GlyphExtents(font.GID(glyphID))
		if ok {
			// go-text/typesetting GlyphExtents:
			//   XBearing: left side bearing
			//   YBearing: top side bearing (positive up in font coords)
			//   Width, Height: extent dimensions
			// Skia Rect:
			//   Left, Top, Right, Bottom with Y increasing downward
			// Note: go-text/typesetting Height is negative for Y-up systems.
			// Skia uses Y-down, so we negate YBearing for Top.
			// For Bottom, we want Top + |Height|. Since Height is negative:
			// Bottom = -YBearing - Height
			return interfaces.Rect{
				Left:   Scalar(extents.XBearing),
				Top:    Scalar(-extents.YBearing),
				Right:  Scalar(extents.XBearing) + Scalar(extents.Width),
				Bottom: Scalar(-extents.YBearing) - Scalar(extents.Height),
			}
		}
	}
	return interfaces.Rect{}
}

// GetGlyphPath returns the outline path for a glyph.
// Returns an error if the glyph has no outline (e.g., space character, bitmap glyph).
func (t *Typeface) GetGlyphPath(glyphID uint16) (interfaces.SkPath, error) {
	if t.goTextFace == nil {
		return nil, errors.New("typeface has no font face")
	}

	glyphData := t.goTextFace.GlyphData(font.GID(glyphID))

	// GlyphData is an interface; we need to type-assert to GlyphOutline for vector glyphs
	outline, ok := glyphData.(font.GlyphOutline)
	if !ok || len(outline.Segments) == 0 {
		return nil, errors.New("glyph has no outline data")
	}

	path := NewSkPath(enums.PathFillTypeDefault)
	for _, seg := range outline.Segments {
		switch seg.Op {
		case ot.SegmentOpMoveTo:
			path.MoveTo(Scalar(seg.Args[0].X), Scalar(seg.Args[0].Y))
		case ot.SegmentOpLineTo:
			path.LineTo(Scalar(seg.Args[0].X), Scalar(seg.Args[0].Y))
		case ot.SegmentOpQuadTo:
			path.QuadTo(
				Scalar(seg.Args[0].X), Scalar(seg.Args[0].Y),
				Scalar(seg.Args[1].X), Scalar(seg.Args[1].Y),
			)
		case ot.SegmentOpCubeTo:
			path.CubicTo(
				Scalar(seg.Args[0].X), Scalar(seg.Args[0].Y),
				Scalar(seg.Args[1].X), Scalar(seg.Args[1].Y),
				Scalar(seg.Args[2].X), Scalar(seg.Args[2].Y),
			)
		}
	}
	return path, nil
}

// Compile-time interface check
var _ SkTypeface = (*Typeface)(nil)
