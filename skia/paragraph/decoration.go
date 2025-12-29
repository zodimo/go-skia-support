package paragraph

// TextDecoration specifies text decoration types.
// Multiple decorations can be combined using bitwise OR.
//
// Ported from: skia-source/modules/skparagraph/include/TextStyle.h
type TextDecoration int

const (
	// TextDecorationNone specifies no decoration.
	TextDecorationNone TextDecoration = 0x0

	// TextDecorationUnderline adds an underline below the text.
	TextDecorationUnderline TextDecoration = 0x1

	// TextDecorationOverline adds a line above the text.
	TextDecorationOverline TextDecoration = 0x2

	// TextDecorationLineThrough adds a strikethrough line through the text.
	TextDecorationLineThrough TextDecoration = 0x4
)

// HasUnderline returns true if this decoration includes underline.
func (d TextDecoration) HasUnderline() bool {
	return d&TextDecorationUnderline != 0
}

// HasOverline returns true if this decoration includes overline.
func (d TextDecoration) HasOverline() bool {
	return d&TextDecorationOverline != 0
}

// HasLineThrough returns true if this decoration includes line-through.
func (d TextDecoration) HasLineThrough() bool {
	return d&TextDecorationLineThrough != 0
}

// TextDecorationStyle specifies the style of text decorations.
//
// Ported from: skia-source/modules/skparagraph/include/TextStyle.h
type TextDecorationStyle int

const (
	// TextDecorationStyleSolid draws a solid line.
	TextDecorationStyleSolid TextDecorationStyle = iota

	// TextDecorationStyleDouble draws a double line.
	TextDecorationStyleDouble

	// TextDecorationStyleDotted draws a dotted line.
	TextDecorationStyleDotted

	// TextDecorationStyleDashed draws a dashed line.
	TextDecorationStyleDashed

	// TextDecorationStyleWavy draws a wavy line.
	TextDecorationStyleWavy
)

// TextDecorationMode specifies how decorations interact with descenders.
//
// Ported from: skia-source/modules/skparagraph/include/TextStyle.h
type TextDecorationMode int

const (
	// TextDecorationModeGaps skips gaps around descenders.
	TextDecorationModeGaps TextDecorationMode = iota

	// TextDecorationModeThrough draws through descenders.
	TextDecorationModeThrough
)

// Decoration holds all decoration-related properties for text styling.
//
// Ported from: skia-source/modules/skparagraph/include/TextStyle.h
type Decoration struct {
	// Type is the decoration type (underline, overline, line-through, or combinations).
	Type TextDecoration

	// Mode determines how decoration interacts with descenders.
	Mode TextDecorationMode

	// Color is the decoration color.
	Color uint32

	// Style is the decoration line style (solid, dashed, etc.).
	Style TextDecorationStyle

	// ThicknessMultiplier is applied to the default thickness.
	ThicknessMultiplier float32
}

// NewDecoration creates a Decoration with default values.
func NewDecoration() Decoration {
	return Decoration{
		Type:                TextDecorationNone,
		Mode:                TextDecorationModeThrough,
		Color:               0x00000000, // SK_ColorTRANSPARENT
		Style:               TextDecorationStyleSolid,
		ThicknessMultiplier: 1.0,
	}
}

// Equals returns true if this decoration equals another.
func (d Decoration) Equals(other Decoration) bool {
	return d.Type == other.Type &&
		d.Mode == other.Mode &&
		d.Color == other.Color &&
		d.Style == other.Style &&
		d.ThicknessMultiplier == other.ThicknessMultiplier
}
