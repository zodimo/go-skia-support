package paragraph

// PlaceholderAlignment specifies how a placeholder is vertically aligned
// relative to the surrounding text.
//
// Ported from: skia-source/modules/skparagraph/include/TextStyle.h
type PlaceholderAlignment int

const (
	// PlaceholderAlignmentBaseline matches the placeholder baseline with the text baseline.
	PlaceholderAlignmentBaseline PlaceholderAlignment = iota

	// PlaceholderAlignmentAboveBaseline aligns the bottom edge of the placeholder
	// with the baseline, so the placeholder sits on top of the baseline.
	PlaceholderAlignmentAboveBaseline

	// PlaceholderAlignmentBelowBaseline aligns the top edge of the placeholder
	// with the baseline, so the placeholder hangs below the baseline.
	PlaceholderAlignmentBelowBaseline

	// PlaceholderAlignmentTop aligns the top edge of the placeholder with the
	// top edge of the font. When the placeholder is very tall, extra space
	// hangs from the top and extends through the bottom of the line.
	PlaceholderAlignmentTop

	// PlaceholderAlignmentBottom aligns the bottom edge of the placeholder with
	// the bottom edge of the font. When the placeholder is very tall, extra space
	// rises from the bottom and extends through the top of the line.
	PlaceholderAlignmentBottom

	// PlaceholderAlignmentMiddle aligns the middle of the placeholder with the
	// middle of the text. When the placeholder is very tall, extra space grows
	// equally from the top and bottom of the line.
	PlaceholderAlignmentMiddle
)

// PlaceholderStyle defines the dimensions and alignment of a placeholder.
// Placeholders are non-text elements (like images or inline widgets) that
// participate in text layout.
//
// Ported from: skia-source/modules/skparagraph/include/TextStyle.h
type PlaceholderStyle struct {
	// Width of the placeholder in logical pixels.
	Width float32

	// Height of the placeholder in logical pixels.
	Height float32

	// Alignment specifies how the placeholder is vertically aligned.
	Alignment PlaceholderAlignment

	// Baseline is the text baseline type used for alignment calculations.
	Baseline TextBaseline

	// BaselineOffset is the distance from the top edge of the rect to the
	// baseline position. This baseline will be aligned against the alphabetic
	// baseline of the surrounding text.
	// Positive values drop the baseline lower (positions the rect higher).
	// Small or negative values cause the rect to be positioned underneath the line.
	// When BaselineOffset == Height, the bottom edge rests on the alphabetic baseline.
	BaselineOffset float32
}

// NewPlaceholderStyle creates a new PlaceholderStyle with default values.
func NewPlaceholderStyle() PlaceholderStyle {
	return PlaceholderStyle{
		Width:          0,
		Height:         0,
		Alignment:      PlaceholderAlignmentBaseline,
		Baseline:       TextBaselineAlphabetic,
		BaselineOffset: 0,
	}
}

// NewPlaceholderStyleWithParams creates a new PlaceholderStyle with the given parameters.
func NewPlaceholderStyleWithParams(width, height float32, alignment PlaceholderAlignment,
	baseline TextBaseline, baselineOffset float32) PlaceholderStyle {
	return PlaceholderStyle{
		Width:          width,
		Height:         height,
		Alignment:      alignment,
		Baseline:       baseline,
		BaselineOffset: baselineOffset,
	}
}

// Equals returns true if this placeholder style equals another.
func (p PlaceholderStyle) Equals(other PlaceholderStyle) bool {
	return p.Width == other.Width &&
		p.Height == other.Height &&
		p.Alignment == other.Alignment &&
		p.Baseline == other.Baseline &&
		p.BaselineOffset == other.BaselineOffset
}
