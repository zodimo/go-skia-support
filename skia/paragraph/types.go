// Package paragraph provides the SkParagraph text layout module types.
// This module enables rich, multi-styled paragraph rendering with proper line breaking,
// text shaping, and layout metrics.
//
// Ported from: skia-source/modules/skparagraph/include/
package paragraph

// Affinity indicates the visual position when a text position is between two runs.
// This is used when a position could be at the end of one run or the start of another.
//
// Ported from: skia-source/modules/skparagraph/include/DartTypes.h
type Affinity int

const (
	// AffinityUpstream indicates the position associates with the preceding run.
	AffinityUpstream Affinity = iota

	// AffinityDownstream indicates the position associates with the following run.
	AffinityDownstream
)

// RectHeightStyle specifies the strategy for computing rectangle heights
// when querying text boxes for a range.
//
// Ported from: skia-source/modules/skparagraph/include/DartTypes.h
type RectHeightStyle int

const (
	// RectHeightStyleTight provides tight bounding boxes that fit heights per run.
	RectHeightStyleTight RectHeightStyle = iota

	// RectHeightStyleMax makes all boxes in a line the maximum height of all runs.
	RectHeightStyleMax

	// RectHeightStyleIncludeLineSpacingMiddle extends boxes to cover half the
	// space above and half below the line.
	RectHeightStyleIncludeLineSpacingMiddle

	// RectHeightStyleIncludeLineSpacingTop adds line spacing to the top of the rect.
	RectHeightStyleIncludeLineSpacingTop

	// RectHeightStyleIncludeLineSpacingBottom adds line spacing to the bottom of the rect.
	RectHeightStyleIncludeLineSpacingBottom

	// RectHeightStyleStrut uses strut-based height calculation.
	RectHeightStyleStrut
)

// RectWidthStyle specifies the strategy for computing rectangle widths
// when querying text boxes for a range.
//
// Ported from: skia-source/modules/skparagraph/include/DartTypes.h
type RectWidthStyle int

const (
	// RectWidthStyleTight provides tight bounding boxes that fit widths per run.
	RectWidthStyleTight RectWidthStyle = iota

	// RectWidthStyleMax extends the last rect of each line to match the widest rect.
	RectWidthStyleMax
)

// TextAlign specifies text alignment within a paragraph.
//
// Ported from: skia-source/modules/skparagraph/include/DartTypes.h
type TextAlign int

const (
	// TextAlignLeft aligns text to the left edge.
	TextAlignLeft TextAlign = iota

	// TextAlignRight aligns text to the right edge.
	TextAlignRight

	// TextAlignCenter centers text horizontally.
	TextAlignCenter

	// TextAlignJustify stretches text to fill the line width.
	TextAlignJustify

	// TextAlignStart aligns to the start edge based on text direction (LTR: left, RTL: right).
	TextAlignStart

	// TextAlignEnd aligns to the end edge based on text direction (LTR: right, RTL: left).
	TextAlignEnd
)

// TextDirection specifies the base direction for text layout.
//
// Ported from: skia-source/modules/skparagraph/include/DartTypes.h
type TextDirection int

const (
	// TextDirectionRTL specifies right-to-left text direction.
	TextDirectionRTL TextDirection = iota

	// TextDirectionLTR specifies left-to-right text direction.
	TextDirectionLTR
)

// TextBaseline specifies the baseline type used for text vertical alignment.
//
// Ported from: skia-source/modules/skparagraph/include/DartTypes.h
type TextBaseline int

const (
	// TextBaselineAlphabetic uses the alphabetic baseline (Latin, Cyrillic, Greek).
	TextBaselineAlphabetic TextBaseline = iota

	// TextBaselineIdeographic uses the ideographic baseline (CJK characters).
	TextBaselineIdeographic
)

// TextHeightBehavior controls how text height is calculated, particularly
// affecting the first and last lines of a paragraph.
//
// Ported from: skia-source/modules/skparagraph/include/DartTypes.h
type TextHeightBehavior int

const (
	// TextHeightBehaviorAll applies height to all lines normally.
	TextHeightBehaviorAll TextHeightBehavior = 0x0

	// TextHeightBehaviorDisableFirstAscent removes extra space above the first line.
	TextHeightBehaviorDisableFirstAscent TextHeightBehavior = 0x1

	// TextHeightBehaviorDisableLastDescent removes extra space below the last line.
	TextHeightBehaviorDisableLastDescent TextHeightBehavior = 0x2

	// TextHeightBehaviorDisableAll removes extra space from both first and last lines.
	TextHeightBehaviorDisableAll TextHeightBehavior = 0x1 | 0x2
)

// LineMetricStyle specifies how line metrics are calculated.
//
// Ported from: skia-source/modules/skparagraph/include/DartTypes.h
type LineMetricStyle int

const (
	// LineMetricStyleTypographic uses ascent, descent, etc. from a fixed baseline.
	LineMetricStyleTypographic LineMetricStyle = iota

	// LineMetricStyleCSS uses CSS-style metrics with leading split and height adjustments.
	LineMetricStyleCSS
)
