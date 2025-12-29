package interfaces

// CodeUnitFlags represents properties of a code unit.
type CodeUnitFlags int

const (
	CodeUnitFlagNone             CodeUnitFlags = 0
	CodeUnitFlagPartOfWhitespace CodeUnitFlags = 1 << 0
	CodeUnitFlagGraphemeStart    CodeUnitFlags = 1 << 1
	CodeUnitFlagControl          CodeUnitFlags = 1 << 2
	CodeUnitFlagHardLineBreak    CodeUnitFlags = 1 << 3
)

// SkUnicode provides Unicode properties and segmentation logic.
//
// Ported from: skia-source/modules/skunicode/include/SkUnicode.h
type SkUnicode interface {
	// FindPreviousGraphemeBoundary finds the start of the grapheme cluster containing the offset.
	FindPreviousGraphemeBoundary(text string, offset int) int

	// IsEmoji returns true if the rune is an emoji.
	IsEmoji(r rune) bool

	// IsEmojiComponent returns true if the rune is an emoji component.
	IsEmojiComponent(r rune) bool

	// IsRegionalIndicator returns true if the rune is a regional indicator.
	IsRegionalIndicator(r rune) bool

	// CodeUnitHasProperty returns true if the code unit at the given index has the specified property.
	CodeUnitHasProperty(text string, offset int, property CodeUnitFlags) bool
}
