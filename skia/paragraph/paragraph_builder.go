package paragraph

// ParagraphBuilder is the entry point for building a Paragraph.
// It follows a builder pattern to construct a paragraph with multiple styles.
//
// Ported from: skia-source/modules/skparagraph/include/ParagraphBuilder.h
type ParagraphBuilder interface {
	// PushStyle pushes a new text style onto the style stack.
	PushStyle(style *TextStyle)

	// Pop pops the current text style from the style stack.
	Pop()

	// PeekStyle returns the current text style from the top of the stack.
	PeekStyle() TextStyle

	// AddText adds text to the builder.
	AddText(text string)

	// AddTextBytes adds text bytes to the builder.
	AddTextBytes(text []byte)

	// AddPlaceholder adds a placeholder to the builder.
	AddPlaceholder(style PlaceholderStyle)

	// Build constructs and returns the Paragraph.
	Build() Paragraph

	// Reset resets the builder to its initial state.
	Reset()

	// GetText returns the current text accumulated in the builder.
	GetText() string

	// GetParagraphStyle returns the paragraph style used by the builder.
	GetParagraphStyle() ParagraphStyle
}

// MakeParagraphBuilder creates a new ParagraphBuilder with the given style and font collection.
func MakeParagraphBuilder(style ParagraphStyle, fontCollection *FontCollection) ParagraphBuilder {
	return &paragraphBuilderImpl{
		paragraphStyle: style,
		fontCollection: fontCollection,
		styleStack:     []TextStyle{style.DefaultTextStyle},
	}
}

// paragraphBuilderImpl is the default implementation of ParagraphBuilder.
type paragraphBuilderImpl struct {
	paragraphStyle ParagraphStyle
	fontCollection *FontCollection
	styleStack     []TextStyle
	text           string
	placeholders   []PlaceholderStyle // Simplified storage
}

// PushStyle pushes a new text style onto the style stack.
func (pb *paragraphBuilderImpl) PushStyle(style *TextStyle) {
	if style != nil {
		pb.styleStack = append(pb.styleStack, *style)
	}
}

// Pop pops the current text style from the style stack.
func (pb *paragraphBuilderImpl) Pop() {
	if len(pb.styleStack) > 1 {
		pb.styleStack = pb.styleStack[:len(pb.styleStack)-1]
	}
}

// PeekStyle returns the current text style from the top of the stack.
func (pb *paragraphBuilderImpl) PeekStyle() TextStyle {
	if len(pb.styleStack) > 0 {
		return pb.styleStack[len(pb.styleStack)-1]
	}
	return pb.paragraphStyle.DefaultTextStyle
}

// AddText adds text to the builder.
func (pb *paragraphBuilderImpl) AddText(text string) {
	pb.text += text
}

// AddTextBytes adds text bytes to the builder.
func (pb *paragraphBuilderImpl) AddTextBytes(text []byte) {
	pb.text += string(text)
}

// AddPlaceholder adds a placeholder to the builder.
func (pb *paragraphBuilderImpl) AddPlaceholder(style PlaceholderStyle) {
	pb.placeholders = append(pb.placeholders, style)
}

// Build constructs and returns the Paragraph.
func (pb *paragraphBuilderImpl) Build() Paragraph {
	// Collect text blocks with their styles
	blocks := []Block{
		{
			Range: NewTextRange(0, len(pb.text)),
			Style: pb.paragraphStyle.DefaultTextStyle,
		},
	}

	// Create the paragraph implementation
	// Placeholders are converted to internal format during layout
	return NewParagraphImpl(
		pb.text,
		pb.paragraphStyle,
		blocks,
		nil, // Placeholders will be handled properly when full tracking is implemented
		pb.fontCollection,
		nil, // Unicode interface can be nil for basic usage
	)
}

// Reset resets the builder to its initial state.
func (pb *paragraphBuilderImpl) Reset() {
	pb.styleStack = []TextStyle{pb.paragraphStyle.DefaultTextStyle}
	pb.text = ""
	pb.placeholders = nil
}

// GetText returns the current text accumulated in the builder.
func (pb *paragraphBuilderImpl) GetText() string {
	return pb.text
}

// GetParagraphStyle returns the paragraph style used by the builder.
func (pb *paragraphBuilderImpl) GetParagraphStyle() ParagraphStyle {
	return pb.paragraphStyle
}
