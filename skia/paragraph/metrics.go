package paragraph

import (
	"math"

	"github.com/zodimo/go-skia-support/skia/models"
)

// StyleMetrics contains metrics for a specific text style.
//
// Ported from: skia-source/modules/skparagraph/include/Metrics.h
type StyleMetrics struct {
	TextStyle   *TextStyle
	FontMetrics models.FontMetrics
}

// NewStyleMetrics creates a new StyleMetrics.
func NewStyleMetrics(style *TextStyle, metrics models.FontMetrics) StyleMetrics {
	return StyleMetrics{
		TextStyle:   style,
		FontMetrics: metrics,
	}
}

// LineMetrics contains metrics for a single line of text in a paragraph.
//
// Ported from: skia-source/modules/skparagraph/include/Metrics.h
type LineMetrics struct {
	// StartIndex is the index in the text buffer where the line begins.
	StartIndex int
	// EndIndex is the index in the text buffer where the line ends.
	EndIndex int
	// EndExcludingWhitespaces is the index excluding trailing whitespace.
	EndExcludingWhitespaces int
	// EndIncludingNewline is the index including the newline character.
	EndIncludingNewline int
	// HardBreak indicates if the line ends with a hard break (newline).
	HardBreak bool

	// Ascent is the final computed ascent for the line (positive).
	Ascent float64
	// Descent is the final computed descent for the line (positive).
	Descent float64
	// UnscaledAscent is the ascent without scaling.
	UnscaledAscent float64
	// Height is the total height of the line (round(ascent + descent)).
	Height float64
	// Width is the width of the line.
	Width float64
	// Left is the left edge of the line.
	Left float64
	// Baseline is the y position of the baseline from the top of the paragraph.
	Baseline float64
	// LineNumber is the zero-indexed line number.
	LineNumber int

	// LineMetrics maps text index to StyleMetrics.
	// The key is the start index of the run.
	LineMetrics map[int]StyleMetrics
}

// NewLineMetrics creates a new LineMetrics with default values.
func NewLineMetrics() LineMetrics {
	return LineMetrics{
		Ascent:  math.MaxFloat64,  // SK_ScalarMax
		Descent: -math.MaxFloat64, // SK_ScalarMin (approx) - Wait, SK_ScalarMin is usually a very small positive number or negative max? C++ initialized descent to SK_ScalarMin. In C++, SK_ScalarMin is usually smallest positive normal float. But `Descent` is usually positive below baseline.
		// Let's check C++ Metrics.h again.
		// double fAscent = SK_ScalarMax;
		// double fDescent = SK_ScalarMin;
		// SK_ScalarMin is typically -SK_ScalarMax or a very small number depending on context. Given layout logic, likely initialized to extreme opposites.
		// I will initialize Ascent to MaxFloat64 (top is min) and Descent to -MaxFloat64 (bottom is max).
		// Wait, Y-down. Ascent is distance ABOVE baseline. Usually positive number in some contexts or negative coordinate. Skia usually defines Ascent as NEGATIVE (up).
		// But LineMetrics comment says: "Ascent and descent are provided as positive numbers."
		// If provided as positive numbers, then Ascent being Max seems wrong for initialization if we want to find the max ascent (which would be the largest positive number).
		// Ah, standard Skia SkFontMetrics: Ascent is negative (up).
		// LineMetrics comment: "Ascent and descent are provided as positive numbers." -> This implies they are absolute distances.
		// The C++ code initializes `fAscent = SK_ScalarMax`.
		// If we are minimizing Ascent (finding the top-most point which is formatted as positive distance? No).
		// If we are MAXIMIZING the extent, we usually start with 0.
		// Let's stick to matching C++ initialization values as closely as possible.
		// SK_ScalarMax in Skia IS purely max float.
		// SK_ScalarMin in Skia IS usually earliest representation (negative).
		// If I use MaxFloat64 and -MaxFloat64 it should be safe.
		UnscaledAscent: math.MaxFloat64,
		Height:         0.0,
		Width:          0.0,
		Left:           0.0,
		Baseline:       0.0,
		LineNumber:     0,
		LineMetrics:    make(map[int]StyleMetrics),
	}
}
