package paragraph

import (
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// VisitorFlags specifies options for the Visit method.
//
// Ported from: skia-source/modules/skparagraph/include/Paragraph.h
type VisitorFlags int

const (
	// VisitorFlagsNone indicates no special flags.
	VisitorFlagsNone VisitorFlags = 0

	// VisitorFlagsWhiteSpace indicates that the visitor should be called for whitespace.
	VisitorFlagsWhiteSpace VisitorFlags = 1 << 0
)

// VisitorInfo contains information passed to the Visitor callback.
//
// Ported from: skia-source/modules/skparagraph/include/Paragraph.h
type VisitorInfo struct {
	Font       interfaces.SkFont
	Origin     models.Point
	Advance    float32
	Glyphs     []uint16 // SkGlyphID
	Positions  []models.Point
	Utf8Starts []uint32
	Flags      VisitorFlags
}

// ExtendedVisitorInfo extends VisitorInfo with bounds information.
//
// Ported from: skia-source/modules/skparagraph/include/Paragraph.h
type ExtendedVisitorInfo struct {
	VisitorInfo
	Bounds models.Rect
}

// Visitor is a callback function for visiting text runs.
type Visitor func(info VisitorInfo)

// ExtendedVisitor is a callback function for visiting text runs with extended info.
type ExtendedVisitor func(info ExtendedVisitorInfo)

// GlyphClusterInfo contains information about a glyph cluster.
//
// Ported from: skia-source/modules/skparagraph/include/Paragraph.h
type GlyphClusterInfo struct {
	Bounds    models.Rect
	TextRange TextRange
	Direction TextDirection
}

// GlyphInfo contains detailed information about a single glyph or grapheme.
//
// Ported from: skia-source/modules/skparagraph/include/Paragraph.h
type GlyphInfo struct {
	GraphemeBounds models.Rect
	TextRange      TextRange
	Direction      TextDirection
	IsEllipsis     bool
}

// FontInfo contains information about the font used for a text range.
//
// Ported from: skia-source/modules/skparagraph/include/Paragraph.h
type FontInfo struct {
	Font      interfaces.SkFont
	TextRange TextRange
}

// Paragraph is the primary interface for layout and rendering of rich text.
//
// Ported from: skia-source/modules/skparagraph/include/Paragraph.h
type Paragraph interface {
	// --- Layout metrics ---
	GetMaxWidth() float32
	GetHeight() float32
	GetMinIntrinsicWidth() float32
	GetMaxIntrinsicWidth() float32
	GetAlphabeticBaseline() float32
	GetIdeographicBaseline() float32
	GetLongestLine() float32
	DidExceedMaxLines() bool

	// --- Layout ---
	Layout(width float32)

	// --- Paint ---
	Paint(canvas interfaces.SkCanvas, x, y float32)
	PaintWithPainter(painter ParagraphPainter, x, y float32)

	// --- Query (rects) ---
	GetRectsForRange(start, end int, heightStyle RectHeightStyle, widthStyle RectWidthStyle) []TextBox
	GetRectsForPlaceholders() []TextBox

	// --- Query (position) ---
	GetGlyphPositionAtCoordinate(dx, dy float32) PositionWithAffinity
	GetWordBoundary(offset int) Range[int]

	// --- Line metrics ---
	GetLineMetrics() []LineMetrics
	LineNumber() int
	GetLineMetricsAt(lineNumber int, lineMetrics *LineMetrics) bool
	GetActualTextRange(lineNumber int, includeSpaces bool) TextRange

	// --- Glyph info ---
	GetGlyphClusterAt(codeUnitIndex int, glyphInfo *GlyphClusterInfo) bool
	GetClosestGlyphClusterAt(dx, dy float32, glyphInfo *GlyphClusterInfo) bool
	GetGlyphInfoAtUTF16Offset(codeUnitIndex int, glyphInfo *GlyphInfo) bool
	GetClosestUTF16GlyphInfoAt(dx, dy float32, glyphInfo *GlyphInfo) bool

	// --- Font info ---
	GetFontAt(codeUnitIndex int) FontInfo
	GetFontAtUTF16Offset(codeUnitIndex int) FontInfo
	GetFonts() []FontInfo

	// --- Line number queries ---
	GetLineNumberAt(codeUnitIndex int) int
	GetLineNumberAtUTF16Offset(codeUnitIndex int) int

	// --- Visitor pattern ---
	Visit(visitor Visitor)
	ExtendedVisit(visitor ExtendedVisitor)
	GetPath(lineNumber int) interfaces.SkPath

	// --- Emoji/color checks ---
	ContainsEmoji(blob interfaces.SkTextBlob) bool
	ContainsColorFontOrBitmap(blob interfaces.SkTextBlob) bool

	// --- Updates ---
	UpdateTextAlign(textAlign TextAlign)
	UpdateFontSize(from, to int, fontSize float32)
	UpdateForegroundPaint(from, to int, paint interfaces.SkPaint)
	UpdateBackgroundPaint(from, to int, paint interfaces.SkPaint)

	// --- State ---
	MarkDirty()
	UnresolvedGlyphs() int
	UnresolvedCodepoints() []rune
}
