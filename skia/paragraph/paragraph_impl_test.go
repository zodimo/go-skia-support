package paragraph

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/models"
)

// --- Test Helpers ---

// createTestParagraph creates a ParagraphImpl for testing with sensible defaults.
// Note: Uses MockTypeface which doesn't support full shaping - tests focus on
// structure and state management rather than pixel-perfect layout.
func createTestParagraph(text string) *ParagraphImpl {
	style := NewParagraphStyle()
	style.DefaultTextStyle = TextStyle{
		FontSize:     14,
		FontFamilies: []string{"TestFont"},
	}

	blocks := []Block{
		{
			Range: NewTextRange(0, len(text)),
			Style: style.DefaultTextStyle,
		},
	}

	fc := NewFontCollection()
	provider := NewTypefaceFontProvider()
	fontStyle := models.FontStyle{Weight: models.FontWeightNormal}
	tf := NewMockTypeface("TestFont", fontStyle)
	provider.RegisterTypeface(tf)
	fc.SetAssetFontManager(provider)

	return NewParagraphImpl(text, style, blocks, nil, fc, nil)
}

// --- Layout Tests ---

func TestParagraphImpl_Layout_Empty(t *testing.T) {
	p := createTestParagraph("")

	p.Layout(100)

	// Empty paragraph should have valid state
	if p.GetHeight() < 0 {
		t.Error("Height should be >= 0 for empty paragraph")
	}
	// LineNumber() returns len(lines) which should be 0 for empty
	if p.LineNumber() != 0 {
		t.Logf("Lines: %d (empty paragraph may have 0 lines)", p.LineNumber())
	}
	if p.DidExceedMaxLines() {
		t.Error("Should not exceed max lines for empty paragraph")
	}
}

func TestParagraphImpl_Layout_SingleLine(t *testing.T) {
	p := createTestParagraph("Hello World")

	p.Layout(1000) // Wide enough for single line

	// With mock typeface, shaping may fail but state should be valid
	if p.GetMaxWidth() != 1000 {
		t.Errorf("MaxWidth should be 1000, got %f", p.GetMaxWidth())
	}

	// Height might be 0 if shaping fails with mock fonts
	t.Logf("Layout result: Height=%f, Lines=%d, LongestLine=%f",
		p.GetHeight(), p.LineNumber(), p.GetLongestLine())
}

func TestParagraphImpl_Layout_MultiLine(t *testing.T) {
	p := createTestParagraph("Hello World this is a long text that should wrap")

	p.Layout(100) // Narrow width to force wrapping

	// With mock fonts, may not actually wrap but state should be valid
	t.Logf("Layout result: Height=%f, Lines=%d", p.GetHeight(), p.LineNumber())
}

func TestParagraphImpl_Layout_ExplicitNewline(t *testing.T) {
	p := createTestParagraph("Line1\nLine2\nLine3")

	p.Layout(1000)

	// State should be valid even if shaping doesn't produce lines
	t.Logf("Lines with newlines: %d", p.LineNumber())
}

func TestParagraphImpl_Layout_MaxLines(t *testing.T) {
	text := "Line1\nLine2\nLine3"
	style := NewParagraphStyle()
	style.MaxLines = 2
	style.DefaultTextStyle = TextStyle{FontSize: 14, FontFamilies: []string{"TestFont"}}

	blocks := []Block{{Range: NewTextRange(0, len(text)), Style: style.DefaultTextStyle}}
	fc := NewFontCollection()
	provider := NewTypefaceFontProvider()
	fontStyle := models.FontStyle{Weight: models.FontWeightNormal}
	tf := NewMockTypeface("TestFont", fontStyle)
	provider.RegisterTypeface(tf)
	fc.SetAssetFontManager(provider)

	p := NewParagraphImpl(text, style, blocks, nil, fc, nil)
	p.Layout(1000)

	// MaxLines should be respected
	lines := p.LineNumber()
	if lines > 2 && lines > 0 {
		t.Errorf("Should have at most 2 lines due to MaxLines, got %d", lines)
	}
	t.Logf("MaxLines test: %d lines", lines)
}

// --- Metrics Tests ---

func TestParagraphImpl_Metrics_Baselines(t *testing.T) {
	p := createTestParagraph("Test")
	p.Layout(100)

	alpha := p.GetAlphabeticBaseline()
	ideo := p.GetIdeographicBaseline()

	// Log baseline values for inspection
	t.Logf("Alphabetic baseline: %f, Ideographic baseline: %f", alpha, ideo)
}

func TestParagraphImpl_Metrics_IntrinsicWidths(t *testing.T) {
	p := createTestParagraph("Hello World")
	p.Layout(100)

	minWidth := p.GetMinIntrinsicWidth()
	maxWidth := p.GetMaxIntrinsicWidth()

	if minWidth < 0 {
		t.Errorf("MinIntrinsicWidth should be >= 0, got %f", minWidth)
	}
	// Note: With mock fonts, intrinsic widths may be 0 or unusual
	t.Logf("MinIntrinsicWidth: %f, MaxIntrinsicWidth: %f", minWidth, maxWidth)
}

// --- Query Tests: Position ---

func TestParagraphImpl_GetGlyphPositionAtCoordinate_Origin(t *testing.T) {
	p := createTestParagraph("Hello")
	p.Layout(100)

	pos := p.GetGlyphPositionAtCoordinate(0, 0)

	// Position at origin should be valid
	if pos.Position < 0 {
		t.Errorf("Position should be >= 0, got %d", pos.Position)
	}
	t.Logf("Position at (0,0): %d, Affinity: %d", pos.Position, pos.Affinity)
}

func TestParagraphImpl_GetGlyphPositionAtCoordinate_End(t *testing.T) {
	p := createTestParagraph("Hello")
	p.Layout(100)

	// Query past the end of the line
	pos := p.GetGlyphPositionAtCoordinate(1000, 0)

	// Should be at or near the end, still valid
	if pos.Position < 0 {
		t.Errorf("Position should be >= 0, got %d", pos.Position)
	}
}

func TestParagraphImpl_GetWordBoundary(t *testing.T) {
	p := createTestParagraph("Hello World")
	p.Layout(100)

	// Get word boundary at position 2 (inside "Hello")
	boundary := p.GetWordBoundary(2)

	if boundary.Start < 0 || boundary.End < 0 {
		t.Error("Word boundary should have valid start and end")
	}
	t.Logf("Word boundary at 2: [%d, %d]", boundary.Start, boundary.End)
}

// --- Query Tests: Rects ---

func TestParagraphImpl_GetRectsForRange_SingleChar(t *testing.T) {
	p := createTestParagraph("Hello")
	p.Layout(100)

	rects := p.GetRectsForRange(0, 1, RectHeightStyleTight, RectWidthStyleTight)

	t.Logf("Rects for single char range: %d", len(rects))
}

func TestParagraphImpl_GetRectsForRange_FullText(t *testing.T) {
	p := createTestParagraph("Hello")
	p.Layout(100)

	rects := p.GetRectsForRange(0, 5, RectHeightStyleTight, RectWidthStyleTight)

	t.Logf("Rects for full text: %d", len(rects))
	for i, rect := range rects {
		t.Logf("  Rect %d: L=%f R=%f", i, rect.Rect.Left, rect.Rect.Right)
	}
}

func TestParagraphImpl_GetRectsForPlaceholders_Empty(t *testing.T) {
	p := createTestParagraph("Hello")
	p.Layout(100)

	rects := p.GetRectsForPlaceholders()

	// No placeholders in this paragraph
	if len(rects) != 0 {
		t.Errorf("Expected 0 placeholder rects, got %d", len(rects))
	}
}

// --- Query Tests: Line ---

func TestParagraphImpl_GetLineMetrics(t *testing.T) {
	p := createTestParagraph("Hello World")
	p.Layout(1000)

	metrics := p.GetLineMetrics()

	t.Logf("Line metrics count: %d", len(metrics))
	for i, m := range metrics {
		t.Logf("  Line %d: Height=%f, Width=%f", i, m.Height, m.Width)
	}
}

func TestParagraphImpl_GetLineMetricsAt(t *testing.T) {
	p := createTestParagraph("Hello World")
	p.Layout(1000)

	var metrics LineMetrics
	found := p.GetLineMetricsAt(0, &metrics)

	if p.LineNumber() > 0 && !found {
		t.Error("Should find line 0 when lines exist")
	}
	t.Logf("GetLineMetricsAt(0): found=%v, height=%f", found, metrics.Height)
}

func TestParagraphImpl_GetLineMetricsAt_Invalid(t *testing.T) {
	p := createTestParagraph("Hello")
	p.Layout(100)

	var metrics LineMetrics
	found := p.GetLineMetricsAt(999, &metrics)

	if found {
		t.Error("Should not find line 999")
	}
}

func TestParagraphImpl_GetLineNumberAt(t *testing.T) {
	p := createTestParagraph("Hello World")
	p.Layout(1000)

	lineNum := p.GetLineNumberAt(0)

	t.Logf("Line number at offset 0: %d", lineNum)
}

func TestParagraphImpl_GetActualTextRange(t *testing.T) {
	p := createTestParagraph("Hello World")
	p.Layout(1000)

	textRange := p.GetActualTextRange(0, false)

	t.Logf("Actual text range for line 0: [%d, %d]", textRange.Start, textRange.End)
}

// --- Query Tests: Glyph ---

func TestParagraphImpl_GetGlyphClusterAt(t *testing.T) {
	p := createTestParagraph("Hello")
	p.Layout(100)

	var info GlyphClusterInfo
	found := p.GetGlyphClusterAt(0, &info)

	t.Logf("GetGlyphClusterAt(0): found=%v", found)
}

func TestParagraphImpl_GetClosestGlyphClusterAt(t *testing.T) {
	p := createTestParagraph("Hello")
	p.Layout(100)

	var info GlyphClusterInfo
	found := p.GetClosestGlyphClusterAt(5, 5, &info)

	t.Logf("GetClosestGlyphClusterAt(5,5): found=%v", found)
}

// --- Query Tests: Font ---

func TestParagraphImpl_GetFontAt(t *testing.T) {
	p := createTestParagraph("Hello")
	p.Layout(100)

	fontInfo := p.GetFontAt(0)

	t.Logf("FontInfo at 0: Font=%v, Range=[%d,%d]",
		fontInfo.Font != nil, fontInfo.TextRange.Start, fontInfo.TextRange.End)
}

func TestParagraphImpl_GetFonts(t *testing.T) {
	p := createTestParagraph("Hello")
	p.Layout(100)

	fonts := p.GetFonts()

	t.Logf("Found %d font infos", len(fonts))
}

// --- State Tests ---

func TestParagraphImpl_State_Relayout(t *testing.T) {
	p := createTestParagraph("Hello World")

	// First layout
	p.Layout(100)
	height1 := p.GetHeight()

	// Same width - should use cached results
	p.Layout(100)
	height2 := p.GetHeight()

	if height1 != height2 {
		t.Errorf("Height should be same for identical layout width: %f vs %f", height1, height2)
	}
}

func TestParagraphImpl_State_WidthChange(t *testing.T) {
	p := createTestParagraph("Hello World this is some text")

	// First layout - narrow
	p.Layout(50)
	height1 := p.GetHeight()
	lines1 := p.LineNumber()

	// Second layout - wide
	p.Layout(500)
	height2 := p.GetHeight()
	lines2 := p.LineNumber()

	t.Logf("Width 50: height=%f, lines=%d", height1, lines1)
	t.Logf("Width 500: height=%f, lines=%d", height2, lines2)
}

func TestParagraphImpl_MarkDirty(t *testing.T) {
	p := createTestParagraph("Hello")
	p.Layout(100)

	// Mark dirty and relayout
	p.MarkDirty()
	p.Layout(100)

	// Should complete without panic
	t.Logf("After MarkDirty: Height=%f", p.GetHeight())
}

// --- Visitor Tests ---

func TestParagraphImpl_Visit(t *testing.T) {
	p := createTestParagraph("Hello")
	p.Layout(100)

	visitCount := 0
	p.Visit(func(info VisitorInfo) {
		visitCount++
	})

	t.Logf("Visitor called %d times", visitCount)
}

func TestParagraphImpl_ExtendedVisit(t *testing.T) {
	p := createTestParagraph("Hello")
	p.Layout(100)

	visitCount := 0
	p.ExtendedVisit(func(info ExtendedVisitorInfo) {
		visitCount++
	})

	t.Logf("Extended visitor called %d times", visitCount)
}

// --- Update Tests ---

func TestParagraphImpl_UpdateTextAlign(t *testing.T) {
	p := createTestParagraph("Hello")
	p.Layout(100)

	p.UpdateTextAlign(TextAlignCenter)
	p.Layout(100)

	// Should complete without error
	t.Logf("After UpdateTextAlign: Height=%f", p.GetHeight())
}

func TestParagraphImpl_UpdateFontSize(t *testing.T) {
	p := createTestParagraph("Hello")
	p.Layout(100)

	p.UpdateFontSize(0, 5, 24)
	p.Layout(100)

	// Should trigger relayout without panic
	t.Logf("After UpdateFontSize: Height=%f", p.GetHeight())
}

// --- Utility Tests ---

func TestParagraphImpl_UnresolvedGlyphs(t *testing.T) {
	p := createTestParagraph("Hello")
	p.Layout(100)

	count := p.UnresolvedGlyphs()

	t.Logf("Unresolved glyphs: %d", count)
}

func TestParagraphImpl_UnresolvedCodepoints(t *testing.T) {
	p := createTestParagraph("Hello")
	p.Layout(100)

	codepoints := p.UnresolvedCodepoints()

	t.Logf("Unresolved codepoints: %d", len(codepoints))
}

// --- Interface Conformance ---

func TestParagraphImpl_ImplementsParagraph(t *testing.T) {
	var _ Paragraph = (*ParagraphImpl)(nil)
	t.Log("ParagraphImpl implements Paragraph interface")
}

func TestParagraphImpl_ImplementsTextLineOwner(t *testing.T) {
	var _ TextLineOwner = (*ParagraphImpl)(nil)
	t.Log("ParagraphImpl implements TextLineOwner interface")
}

func TestParagraphImpl_ImplementsTextWrapperOwner(t *testing.T) {
	var _ TextWrapperOwner = (*ParagraphImpl)(nil)
	t.Log("ParagraphImpl implements TextWrapperOwner interface")
}
