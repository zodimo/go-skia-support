package paragraph

import (
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// EmptyText is an empty text range.
var EmptyText = NewTextRange(0, 0)

// --- Paint methods ---

// Paint renders the paragraph to a canvas.
func (p *ParagraphImpl) Paint(canvas interfaces.SkCanvas, x, y float32) {
	// Create a simple canvas painter wrapper
	painter := &canvasParagraphPainter{canvas: canvas}
	p.PaintWithPainter(painter, x, y)
}

// PaintWithPainter renders the paragraph using a custom painter.
func (p *ParagraphImpl) PaintWithPainter(painter ParagraphPainter, x, y float32) {
	for _, line := range p.lines {
		line.Paint(painter, x, y)
	}
}

// canvasParagraphPainter wraps a canvas for basic painting.
type canvasParagraphPainter struct {
	canvas interfaces.SkCanvas
}

func (c *canvasParagraphPainter) DrawTextBlob(blob interfaces.SkTextBlob, x, y float32, paint interfaces.SkPaint) {
}
func (c *canvasParagraphPainter) DrawTextShadow(blob interfaces.SkTextBlob, x, y float32, color models.Color4f, blurSigma float64) {
}
func (c *canvasParagraphPainter) DrawRect(rect models.Rect, paint interfaces.SkPaint) {
}
func (c *canvasParagraphPainter) DrawFilledRect(rect models.Rect, style DecorationStyle) {
}
func (c *canvasParagraphPainter) DrawPath(path interfaces.SkPath, style DecorationStyle) {
}
func (c *canvasParagraphPainter) DrawLine(x0, y0, x1, y1 float32, style DecorationStyle) {
}
func (c *canvasParagraphPainter) ClipRect(rect models.Rect) {
}
func (c *canvasParagraphPainter) Translate(dx, dy float32) {
}
func (c *canvasParagraphPainter) Save() {
}
func (c *canvasParagraphPainter) Restore() {
}

// --- Query methods: Rects ---

// GetRectsForRange returns bounding boxes for the given text range.
func (p *ParagraphImpl) GetRectsForRange(start, end int, rectHeightStyle RectHeightStyle, rectWidthStyle RectWidthStyle) []TextBox {
	results := make([]TextBox, 0)

	if len(p.text) == 0 {
		if start == 0 && end > 0 {
			rect := models.Rect{Left: 0, Top: 0, Right: 0, Bottom: models.Scalar(p.height)}
			results = append(results, NewTextBox(rect, p.paragraphStyle.TextDirection))
		}
		return results
	}

	p.ensureUTF16Mapping()

	if start >= end || start >= len(p.utf8IndexForUTF16Index) || end == 0 {
		return results
	}

	// Convert UTF-16 indices to UTF-8
	textRange := NewTextRange(len(p.text), len(p.text))
	if start < len(p.utf8IndexForUTF16Index) {
		utf8 := p.utf8IndexForUTF16Index[start]
		textRange.Start = p.findNextGraphemeBoundary(utf8)
	}
	if end < len(p.utf8IndexForUTF16Index) {
		utf8 := p.findPreviousGraphemeBoundary(p.utf8IndexForUTF16Index[end])
		textRange.End = utf8
	}

	// Get rects from each line (simplified)
	for _, line := range p.lines {
		lineText := line.textIncludingNewlines
		intersect := lineText.Intersection(textRange)
		if intersect.Width() == 0 && lineText.Start != textRange.Start {
			continue
		}
		// Calculate line bounds as approximation
		rect := models.Rect{
			Left:   line.offset.X + models.Scalar(line.shift),
			Top:    line.offset.Y,
			Right:  line.offset.X + line.advance.X,
			Bottom: line.offset.Y + line.advance.Y,
		}
		results = append(results, NewTextBox(rect, p.paragraphStyle.TextDirection))
	}

	return results
}

// GetRectsForPlaceholders returns bounding boxes for all placeholders.
func (p *ParagraphImpl) GetRectsForPlaceholders() []TextBox {
	boxes := make([]TextBox, 0)

	if len(p.text) == 0 || len(p.placeholders) == 1 {
		return boxes
	}

	// Iterate through clusters to find placeholders
	for _, cluster := range p.clusters {
		if run := cluster.Run(); run != nil && run.IsPlaceholder() {
			rect := models.Rect{
				Left:   models.Scalar(cluster.StartPos()),
				Top:    0,
				Right:  models.Scalar(cluster.EndPos()),
				Bottom: models.Scalar(p.height),
			}
			boxes = append(boxes, NewTextBox(rect, p.paragraphStyle.TextDirection))
		}
	}

	return boxes
}

// --- Query methods: Position ---

// GetGlyphPositionAtCoordinate returns the position nearest to the given coordinates.
func (p *ParagraphImpl) GetGlyphPositionAtCoordinate(dx, dy float32) PositionWithAffinity {
	if len(p.text) == 0 {
		return NewPositionWithAffinityDefault()
	}

	p.ensureUTF16Mapping()

	for i, line := range p.lines {
		offsetY := float32(line.offset.Y)
		lineHeight := line.Height()

		// Check if this is our line
		if dy < offsetY+lineHeight || i == len(p.lines)-1 {
			// Calculate approximate position based on x
			lineWidth := line.Width()
			if lineWidth <= 0 {
				return NewPositionWithAffinityDefault()
			}

			// Simple linear approximation
			lineStart := line.textExcludingSpaces.Start
			lineEnd := line.textExcludingSpaces.End
			lineChars := lineEnd - lineStart

			relativeX := dx - float32(line.offset.X) - line.shift
			fraction := relativeX / lineWidth
			if fraction < 0 {
				fraction = 0
			}
			if fraction > 1 {
				fraction = 1
			}

			charPos := lineStart + int(fraction*float32(lineChars))
			if charPos < len(p.utf16IndexForUTF8Index) {
				return NewPositionWithAffinity(int32(p.utf16IndexForUTF8Index[charPos]), AffinityDownstream)
			}
			return NewPositionWithAffinityDefault()
		}
	}

	return NewPositionWithAffinityDefault()
}

// GetWordBoundary returns the word boundaries at the given offset.
func (p *ParagraphImpl) GetWordBoundary(offset int) Range[int] {
	if len(p.words) == 0 {
		p.words = p.computeWords()
	}

	start := 0
	end := 0
	for i := 0; i < len(p.words); i++ {
		word := p.words[i]
		if word <= offset {
			start = word
			end = word
		} else {
			end = word
			break
		}
	}

	return NewRange(start, end)
}

func (p *ParagraphImpl) computeWords() []int {
	words := make([]int, 0)
	inWord := false

	for i := 0; i < len(p.text); i++ {
		ch := p.text[i]
		isWordChar := ch != ' ' && ch != '\t' && ch != '\n' && ch != '\r'

		if isWordChar && !inWord {
			words = append(words, i)
			inWord = true
		} else if !isWordChar && inWord {
			words = append(words, i)
			inWord = false
		}
	}

	if inWord {
		words = append(words, len(p.text))
	}

	return words
}

// --- Query methods: Line ---

// GetLineMetrics returns metrics for all lines.
func (p *ParagraphImpl) GetLineMetrics() []LineMetrics {
	metrics := make([]LineMetrics, len(p.lines))
	for i, line := range p.lines {
		metrics[i] = LineMetrics{
			StartIndex: line.textExcludingSpaces.Start,
			EndIndex:   line.textExcludingSpaces.End,
			Ascent:     float64(-line.sizes.Ascent),
			Descent:    float64(line.sizes.Descent),
			Height:     float64(line.Height()),
			Width:      float64(line.Width()),
			Left:       float64(line.offset.X),
			Baseline:   float64(line.Baseline()),
			LineNumber: i,
			HardBreak:  line.isHardBreak(),
		}
	}
	return metrics
}

// GetLineMetricsAt returns metrics for a specific line.
func (p *ParagraphImpl) GetLineMetricsAt(lineNumber int, lineMetrics *LineMetrics) bool {
	if lineNumber < 0 || lineNumber >= len(p.lines) {
		return false
	}
	if lineMetrics != nil {
		line := p.lines[lineNumber]
		*lineMetrics = LineMetrics{
			StartIndex: line.textExcludingSpaces.Start,
			EndIndex:   line.textExcludingSpaces.End,
			Ascent:     float64(-line.sizes.Ascent),
			Descent:    float64(line.sizes.Descent),
			Height:     float64(line.Height()),
			Width:      float64(line.Width()),
			Left:       float64(line.offset.X),
			Baseline:   float64(line.Baseline()),
			LineNumber: lineNumber,
			HardBreak:  line.isHardBreak(),
		}
	}
	return true
}

// GetLineNumberAt returns the line number at the given code unit index.
func (p *ParagraphImpl) GetLineNumberAt(codeUnitIndex int) int {
	if codeUnitIndex >= len(p.text) || len(p.lines) == 0 {
		return -1
	}

	for i, line := range p.lines {
		if codeUnitIndex >= line.textIncludingNewlines.Start && codeUnitIndex < line.textIncludingNewlines.End {
			return i
		}
	}

	return -1
}

// GetLineNumberAtUTF16Offset returns the line number at the given UTF-16 offset.
func (p *ParagraphImpl) GetLineNumberAtUTF16Offset(codeUnitIndex int) int {
	p.ensureUTF16Mapping()
	if codeUnitIndex >= len(p.utf8IndexForUTF16Index) {
		return -1
	}
	utf8 := p.utf8IndexForUTF16Index[codeUnitIndex]
	return p.GetLineNumberAt(utf8)
}

// GetActualTextRange returns the actual text range for a line.
func (p *ParagraphImpl) GetActualTextRange(lineNumber int, includeSpaces bool) TextRange {
	if lineNumber < 0 || lineNumber >= len(p.lines) {
		return EmptyText
	}
	line := p.lines[lineNumber]
	if includeSpaces {
		return line.text
	}
	return line.textExcludingSpaces
}

// --- Query methods: Glyph/Cluster ---

// GetGlyphClusterAt returns the glyph cluster at the given code unit index.
func (p *ParagraphImpl) GetGlyphClusterAt(codeUnitIndex int, glyphInfo *GlyphClusterInfo) bool {
	lineNumber := p.GetLineNumberAt(codeUnitIndex)
	if lineNumber == -1 {
		return false
	}

	line := p.lines[lineNumber]
	for c := line.clusterRange.Start; c < line.clusterRange.End; c++ {
		cluster := p.clusters[c]
		if cluster.Contains(codeUnitIndex) {
			if glyphInfo != nil {
				rect := models.Rect{
					Left:  models.Scalar(cluster.StartPos()),
					Right: models.Scalar(cluster.EndPos()),
				}
				*glyphInfo = GlyphClusterInfo{
					Bounds:    rect,
					TextRange: cluster.TextRange(),
					Direction: p.paragraphStyle.TextDirection,
				}
			}
			return true
		}
	}

	return false
}

// GetClosestGlyphClusterAt returns the closest glyph cluster to the coordinates.
func (p *ParagraphImpl) GetClosestGlyphClusterAt(dx, dy float32, glyphInfo *GlyphClusterInfo) bool {
	res := p.GetGlyphPositionAtCoordinate(dx, dy)

	utf16Offset := int(res.Position)
	if res.Affinity == AffinityUpstream && utf16Offset > 0 {
		utf16Offset = utf16Offset - 1
	}

	p.ensureUTF16Mapping()
	if utf16Offset >= len(p.utf8IndexForUTF16Index) {
		return false
	}

	return p.GetGlyphClusterAt(p.utf8IndexForUTF16Index[utf16Offset], glyphInfo)
}

// GetGlyphInfoAtUTF16Offset returns glyph info at the given UTF-16 offset.
func (p *ParagraphImpl) GetGlyphInfoAtUTF16Offset(codeUnitIndex int, glyphInfo *GlyphInfo) bool {
	p.ensureUTF16Mapping()
	if codeUnitIndex >= len(p.utf8IndexForUTF16Index) {
		return false
	}

	utf8 := p.utf8IndexForUTF16Index[codeUnitIndex]
	lineNumber := p.GetLineNumberAt(utf8)
	if lineNumber == -1 {
		return false
	}

	if glyphInfo != nil {
		glyphClusterIndex := p.clusterIndex(utf8)
		run := (*Run)(nil)
		if glyphClusterIndex >= 0 && glyphClusterIndex < len(p.clusters) {
			run = p.clusters[glyphClusterIndex].Run()
		}

		*glyphInfo = GlyphInfo{
			GraphemeBounds: models.Rect{},
			TextRange:      NewTextRange(utf8, utf8+1),
			Direction:      p.paragraphStyle.TextDirection,
			IsEllipsis:     run != nil && run.IsEllipsis(),
		}
	}
	return true
}

// GetClosestUTF16GlyphInfoAt returns the closest glyph info to the coordinates.
func (p *ParagraphImpl) GetClosestUTF16GlyphInfoAt(dx, dy float32, glyphInfo *GlyphInfo) bool {
	res := p.GetGlyphPositionAtCoordinate(dx, dy)

	utf16Offset := int(res.Position)
	if res.Affinity == AffinityUpstream && utf16Offset > 0 {
		utf16Offset = utf16Offset - 1
	}

	return p.GetGlyphInfoAtUTF16Offset(utf16Offset, glyphInfo)
}

// --- Query methods: Font ---

// GetFontAt returns the font at the given code unit index.
func (p *ParagraphImpl) GetFontAt(codeUnitIndex int) FontInfo {
	for _, run := range p.runs {
		textRange := run.TextRange()
		if textRange.Start <= codeUnitIndex && codeUnitIndex < textRange.End {
			return FontInfo{Font: run.Font(), TextRange: textRange}
		}
	}
	return FontInfo{}
}

// GetFontAtUTF16Offset returns the font at the given UTF-16 offset.
func (p *ParagraphImpl) GetFontAtUTF16Offset(codeUnitIndex int) FontInfo {
	p.ensureUTF16Mapping()
	if codeUnitIndex >= len(p.utf8IndexForUTF16Index) {
		return FontInfo{}
	}
	utf8 := p.utf8IndexForUTF16Index[codeUnitIndex]
	return p.GetFontAt(utf8)
}

// GetFonts returns all fonts used in the paragraph with their text ranges.
func (p *ParagraphImpl) GetFonts() []FontInfo {
	results := make([]FontInfo, 0, len(p.runs))
	for _, run := range p.runs {
		results = append(results, FontInfo{
			Font:      run.Font(),
			TextRange: run.TextRange(),
		})
	}
	return results
}

// --- Visitor methods ---

// Visit calls the visitor function for each run in each line.
func (p *ParagraphImpl) Visit(visitor Visitor) {
	if visitor == nil {
		return
	}

	for _, line := range p.lines {
		for _, runIdx := range line.runsInVisualOrder {
			if runIdx < 0 || runIdx >= len(p.runs) {
				continue
			}
			run := p.runs[runIdx]
			if run == nil {
				continue
			}

			info := VisitorInfo{
				Font:    run.Font(),
				Origin:  models.Point{X: line.offset.X + models.Scalar(line.shift), Y: line.offset.Y},
				Advance: float32(run.Advance().X),
				Glyphs:  run.Glyphs(),
			}
			visitor(info)
		}
	}
}

// ExtendedVisit calls the extended visitor function.
func (p *ParagraphImpl) ExtendedVisit(visitor ExtendedVisitor) {
	if visitor == nil {
		return
	}

	for _, line := range p.lines {
		for _, runIdx := range line.runsInVisualOrder {
			if runIdx < 0 || runIdx >= len(p.runs) {
				continue
			}
			run := p.runs[runIdx]
			if run == nil {
				continue
			}

			info := ExtendedVisitorInfo{
				VisitorInfo: VisitorInfo{
					Font:    run.Font(),
					Origin:  models.Point{X: line.offset.X + models.Scalar(line.shift), Y: line.offset.Y},
					Advance: float32(run.Advance().X),
					Glyphs:  run.Glyphs(),
				},
				Bounds: run.Clip(),
			}
			visitor(info)
		}
	}
}

// GetPath returns the glyph outlines for a line.
func (p *ParagraphImpl) GetPath(lineNumber int) interfaces.SkPath {
	// Path generation requires deeper integration
	return nil
}

// --- Utility methods ---

// ContainsEmoji checks if a text blob contains emoji.
func (p *ParagraphImpl) ContainsEmoji(blob interfaces.SkTextBlob) bool {
	if p.unicode == nil {
		return false
	}
	for _, r := range p.text {
		if p.unicode.IsEmoji(r) {
			return true
		}
	}
	return false
}

// ContainsColorFontOrBitmap checks if a text blob contains color fonts.
func (p *ParagraphImpl) ContainsColorFontOrBitmap(blob interfaces.SkTextBlob) bool {
	return false
}

// --- Update methods ---

// UpdateTextAlign updates the text alignment.
func (p *ParagraphImpl) UpdateTextAlign(align TextAlign) {
	if p.paragraphStyle.TextAlign != align {
		p.paragraphStyle.TextAlign = align
		if p.state >= StateLineBroken {
			p.state = StateLineBroken
		}
	}
}

// UpdateFontSize updates the font size for a text range.
func (p *ParagraphImpl) UpdateFontSize(from, to int, size float32) {
	for i := range p.textStyles {
		block := &p.textStyles[i]
		if block.Range.Start >= from && block.Range.End <= to {
			block.Style.FontSize = size
		}
	}
	p.state = StateUnknown
}

// UpdateForegroundPaint updates the foreground paint for a text range.
func (p *ParagraphImpl) UpdateForegroundPaint(from, to int, paint interfaces.SkPaint) {
	// Would update PaintOrID in TextStyle
}

// UpdateBackgroundPaint updates the background paint for a text range.
func (p *ParagraphImpl) UpdateBackgroundPaint(from, to int, paint interfaces.SkPaint) {
	// Would update background in TextStyle
}

// --- UTF-16 Mapping ---

func (p *ParagraphImpl) ensureUTF16Mapping() {
	p.utf16MappingOnce.Do(func() {
		p.utf8IndexForUTF16Index = make([]int, 0, len(p.text))
		p.utf16IndexForUTF8Index = make([]int, len(p.text)+1)

		utf16Idx := 0
		for utf8Idx := 0; utf8Idx < len(p.text); {
			r, size := decodeRuneInString(p.text[utf8Idx:])

			p.utf16IndexForUTF8Index[utf8Idx] = utf16Idx
			p.utf8IndexForUTF16Index = append(p.utf8IndexForUTF16Index, utf8Idx)

			if r > 0xFFFF {
				utf16Idx++
				p.utf8IndexForUTF16Index = append(p.utf8IndexForUTF16Index, utf8Idx)
			}

			utf8Idx += size
			utf16Idx++
		}
		p.utf16IndexForUTF8Index[len(p.text)] = utf16Idx
	})
}

func decodeRuneInString(s string) (rune, int) {
	if len(s) == 0 {
		return 0, 0
	}
	b := s[0]
	if b < 0x80 {
		return rune(b), 1
	}
	if b < 0xC0 {
		return 0xFFFD, 1
	}
	if b < 0xE0 && len(s) >= 2 {
		return rune(b&0x1F)<<6 | rune(s[1]&0x3F), 2
	}
	if b < 0xF0 && len(s) >= 3 {
		return rune(b&0x0F)<<12 | rune(s[1]&0x3F)<<6 | rune(s[2]&0x3F), 3
	}
	if b < 0xF8 && len(s) >= 4 {
		return rune(b&0x07)<<18 | rune(s[1]&0x3F)<<12 | rune(s[2]&0x3F)<<6 | rune(s[3]&0x3F), 4
	}
	return 0xFFFD, 1
}

func (p *ParagraphImpl) findPreviousGraphemeBoundary(offset int) int {
	if p.unicode != nil {
		return p.unicode.FindPreviousGraphemeBoundary(p.text, offset)
	}
	if offset > 0 {
		return offset
	}
	return 0
}

func (p *ParagraphImpl) findNextGraphemeBoundary(offset int) int {
	if offset < len(p.text) {
		return offset
	}
	return len(p.text)
}

// Enable the Paragraph interface check
var _ Paragraph = (*ParagraphImpl)(nil)
