package paragraph

// Block represents a range of text with a specific style.
// Blocks are used internally to track styled text runs.
//
// Ported from: skia-source/modules/skparagraph/include/TextStyle.h
type Block struct {
	// Range is the text range this block covers.
	Range TextRange

	// Style is the text style applied to this range.
	Style TextStyle
}

// NewBlock creates a new Block with the given range and style.
func NewBlock(start, end int, style TextStyle) Block {
	return Block{
		Range: NewTextRange(start, end),
		Style: style,
	}
}

// NewBlockFromRange creates a new Block with the given TextRange and style.
func NewBlockFromRange(textRange TextRange, style TextStyle) Block {
	return Block{
		Range: textRange,
		Style: style,
	}
}

// Add extends this block's range to include a tail range.
// The tail must start where this block ends.
func (b *Block) Add(tail TextRange) {
	// Assertion: b.Range.End == tail.Start
	b.Range = NewTextRange(b.Range.Start, b.Range.Start+b.Range.Width()+tail.Width())
}

// Placeholder represents a placeholder element in the text with its styling.
// Placeholders are non-text elements like images that participate in layout.
//
// Ported from: skia-source/modules/skparagraph/include/TextStyle.h
type Placeholder struct {
	// Range is the text range this placeholder occupies.
	Range TextRange

	// Style is the placeholder's dimensional and alignment settings.
	Style PlaceholderStyle

	// TextStyle is the text style context for this placeholder.
	TextStyle TextStyle

	// BlocksBefore is the range of blocks before this placeholder.
	BlocksBefore BlockRange

	// TextBefore is the text range before this placeholder.
	TextBefore TextRange
}

// NewPlaceholder creates a new Placeholder with the given parameters.
func NewPlaceholder(start, end int, style PlaceholderStyle, textStyle TextStyle,
	blocksBefore BlockRange, textBefore TextRange) Placeholder {
	return Placeholder{
		Range:        NewTextRange(start, end),
		Style:        style,
		TextStyle:    textStyle,
		BlocksBefore: blocksBefore,
		TextBefore:   textBefore,
	}
}

// NewPlaceholderDefault creates a Placeholder with default/empty values.
func NewPlaceholderDefault() Placeholder {
	return Placeholder{
		Range:        EmptyRange,
		Style:        NewPlaceholderStyle(),
		TextStyle:    NewTextStyle(),
		BlocksBefore: EmptyRange,
		TextBefore:   EmptyRange,
	}
}
