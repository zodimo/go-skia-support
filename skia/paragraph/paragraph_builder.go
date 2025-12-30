package paragraph

// placeholderInfo tracks a placeholder during building.
type placeholderInfo struct {
	style     PlaceholderStyle
	textStyle TextStyle
	textStart int // Text position where placeholder is inserted
}

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
		paragraphStyle:   style,
		fontCollection:   fontCollection,
		styleStack:       []TextStyle{style.DefaultTextStyle},
		blocks:           make([]Block, 0),
		placeholderInfos: make([]placeholderInfo, 0),
	}
}

// paragraphBuilderImpl is the default implementation of ParagraphBuilder.
type paragraphBuilderImpl struct {
	paragraphStyle   ParagraphStyle
	fontCollection   *FontCollection
	styleStack       []TextStyle
	text             string
	blocks           []Block           // Styled text blocks
	placeholderInfos []placeholderInfo // Placeholder tracking
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
	if len(text) == 0 {
		return
	}

	startPos := len(pb.text)
	pb.text += text
	endPos := len(pb.text)

	// Create or extend a block with the current style
	currentStyle := pb.PeekStyle()

	if len(pb.blocks) > 0 {
		lastBlock := &pb.blocks[len(pb.blocks)-1]
		// If contiguous, extend the block (simplified: always extend if contiguous)
		if lastBlock.Range.End == startPos {
			lastBlock.Range.End = endPos
			return
		}
	}

	// Create new block
	pb.blocks = append(pb.blocks, Block{
		Range: NewTextRange(startPos, endPos),
		Style: currentStyle,
	})
}

// AddTextBytes adds text bytes to the builder.
func (pb *paragraphBuilderImpl) AddTextBytes(text []byte) {
	pb.AddText(string(text))
}

// AddPlaceholder adds a placeholder to the builder.
// Placeholders are represented as a special character (object replacement character U+FFFC)
// in the text stream and tracked separately for layout.
func (pb *paragraphBuilderImpl) AddPlaceholder(style PlaceholderStyle) {
	textStart := len(pb.text)

	// Insert placeholder marker character (U+FFFC = Object Replacement Character)
	pb.text += "\uFFFC"

	// Track the placeholder with its position and current style
	pb.placeholderInfos = append(pb.placeholderInfos, placeholderInfo{
		style:     style,
		textStyle: pb.PeekStyle(),
		textStart: textStart,
	})

	// Create a block for the placeholder character
	pb.blocks = append(pb.blocks, Block{
		Range: NewTextRange(textStart, len(pb.text)),
		Style: pb.PeekStyle(),
	})
}

// Build constructs and returns the Paragraph.
func (pb *paragraphBuilderImpl) Build() Paragraph {
	// Ensure we have at least one block
	blocks := pb.blocks
	if len(blocks) == 0 {
		blocks = []Block{
			{
				Range: NewTextRange(0, len(pb.text)),
				Style: pb.paragraphStyle.DefaultTextStyle,
			},
		}
	}

	// Convert placeholder infos to internal Placeholder format
	placeholders := make([]Placeholder, len(pb.placeholderInfos)+1)

	// First placeholder is always a sentinel at position 0
	placeholders[0] = Placeholder{
		Range:        NewTextRange(0, 0),
		Style:        NewPlaceholderStyle(),
		TextStyle:    pb.paragraphStyle.DefaultTextStyle,
		BlocksBefore: NewBlockRange(0, 0),
		TextBefore:   NewTextRange(0, 0),
	}

	// Convert each tracked placeholder
	for i, info := range pb.placeholderInfos {
		// Find blocks before this placeholder
		blocksBefore := 0
		for j, block := range blocks {
			if block.Range.Start < info.textStart {
				blocksBefore = j + 1
			}
		}

		placeholders[i+1] = Placeholder{
			Range:        NewTextRange(info.textStart, info.textStart+3), // U+FFFC is 3 bytes in UTF-8
			Style:        info.style,
			TextStyle:    info.textStyle,
			BlocksBefore: NewBlockRange(0, blocksBefore),
			TextBefore:   NewTextRange(0, info.textStart),
		}
	}

	// Create the paragraph implementation
	return NewParagraphImpl(
		pb.text,
		pb.paragraphStyle,
		blocks,
		placeholders,
		pb.fontCollection,
		nil, // Unicode interface can be nil for basic usage
	)
}

// Reset resets the builder to its initial state.
func (pb *paragraphBuilderImpl) Reset() {
	pb.styleStack = []TextStyle{pb.paragraphStyle.DefaultTextStyle}
	pb.text = ""
	pb.blocks = make([]Block, 0)
	pb.placeholderInfos = make([]placeholderInfo, 0)
}

// GetText returns the current text accumulated in the builder.
func (pb *paragraphBuilderImpl) GetText() string {
	return pb.text
}

// GetParagraphStyle returns the paragraph style used by the builder.
func (pb *paragraphBuilderImpl) GetParagraphStyle() ParagraphStyle {
	return pb.paragraphStyle
}
