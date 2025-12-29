package paragraph

import (
	"math"
)

// ParagraphStyle controls the appearance of a paragraph.
//
// Ported from: skia-source/modules/skparagraph/include/ParagraphStyle.h
type ParagraphStyle struct {
	StrutStyle            StrutStyle
	DefaultTextStyle      TextStyle
	TextAlign             TextAlign
	TextDirection         TextDirection
	MaxLines              int
	Ellipsis              string
	EllipsisUtf16         string // Using Go string (UTF-8) for simplicity, can be adapted if needed
	Height                float32
	TextHeightBehavior    TextHeightBehavior
	HintingIsOn           bool
	ReplaceTabCharacters  bool
	FakeMissingFontStyles bool
	ApplyRoundingHack     bool
}

// NewParagraphStyle creates a new ParagraphStyle with default values.
func NewParagraphStyle() ParagraphStyle {
	return ParagraphStyle{
		StrutStyle:            NewStrutStyle(),
		DefaultTextStyle:      NewTextStyle(),
		TextAlign:             TextAlignLeft,
		TextDirection:         TextDirectionLTR,
		MaxLines:              math.MaxInt, // Unlimited
		Ellipsis:              "",
		EllipsisUtf16:         "",
		Height:                1.0,
		TextHeightBehavior:    TextHeightBehaviorAll,
		HintingIsOn:           true,
		ReplaceTabCharacters:  false,
		FakeMissingFontStyles: true,
		ApplyRoundingHack:     true,
	}
}

// GetStrutStyle returns the strut style.
func (p *ParagraphStyle) GetStrutStyle() *StrutStyle {
	return &p.StrutStyle
}

// SetStrutStyle sets the strut style.
func (p *ParagraphStyle) SetStrutStyle(strutStyle StrutStyle) {
	p.StrutStyle = strutStyle
}

// GetTextStyle returns the default text style.
func (p *ParagraphStyle) GetTextStyle() *TextStyle {
	return &p.DefaultTextStyle
}

// SetTextStyle sets the default text style.
func (p *ParagraphStyle) SetTextStyle(textStyle TextStyle) {
	p.DefaultTextStyle = textStyle
}

// GetTextDirection returns the text direction.
func (p *ParagraphStyle) GetTextDirection() TextDirection {
	return p.TextDirection
}

// SetTextDirection sets the text direction.
func (p *ParagraphStyle) SetTextDirection(direction TextDirection) {
	p.TextDirection = direction
}

// GetTextAlign returns the text alignment.
func (p *ParagraphStyle) GetTextAlign() TextAlign {
	return p.TextAlign
}

// SetTextAlign sets the text alignment.
func (p *ParagraphStyle) SetTextAlign(align TextAlign) {
	p.TextAlign = align
}

// GetMaxLines returns the maximum number of lines.
func (p *ParagraphStyle) GetMaxLines() int {
	return p.MaxLines
}

// SetMaxLines sets the maximum number of lines.
func (p *ParagraphStyle) SetMaxLines(maxLines int) {
	p.MaxLines = maxLines
}

// GetEllipsis returns the ellipsis string.
func (p *ParagraphStyle) GetEllipsis() string {
	return p.Ellipsis
}

// SetEllipsis sets the ellipsis string.
func (p *ParagraphStyle) SetEllipsis(ellipsis string) {
	p.Ellipsis = ellipsis
}

// GetHeight returns the height multiplier.
func (p *ParagraphStyle) GetHeight() float32 {
	return p.Height
}

// SetHeight sets the height multiplier.
func (p *ParagraphStyle) SetHeight(height float32) {
	p.Height = height
}

// GetTextHeightBehavior returns the text height behavior.
func (p *ParagraphStyle) GetTextHeightBehavior() TextHeightBehavior {
	return p.TextHeightBehavior
}

// SetTextHeightBehavior sets the text height behavior.
func (p *ParagraphStyle) SetTextHeightBehavior(behavior TextHeightBehavior) {
	p.TextHeightBehavior = behavior
}

// UnlimitedLines returns true if there is no line limit.
func (p *ParagraphStyle) UnlimitedLines() bool {
	return p.MaxLines == math.MaxInt
}

// Ellipsized returns true if an ellipsis is set.
func (p *ParagraphStyle) Ellipsized() bool {
	return p.Ellipsis != "" || p.EllipsisUtf16 != ""
}

// EffectiveAlign returns the effective alignment (interpreting Start/End).
func (p *ParagraphStyle) EffectiveAlign() TextAlign {
	if p.TextAlign == TextAlignStart {
		if p.TextDirection == TextDirectionLTR {
			return TextAlignLeft
		}
		return TextAlignRight
	} else if p.TextAlign == TextAlignEnd {
		if p.TextDirection == TextDirectionLTR {
			return TextAlignRight
		}
		return TextAlignLeft
	}
	return p.TextAlign
}

// IsHintingOn returns whether hinting is enabled.
func (p *ParagraphStyle) IsHintingOn() bool {
	return p.HintingIsOn
}

// TurnHintingOff disables hinting.
func (p *ParagraphStyle) TurnHintingOff() {
	p.HintingIsOn = false
}

// GetReplaceTabCharacters returns whether tab characters should be replaced.
func (p *ParagraphStyle) GetReplaceTabCharacters() bool {
	return p.ReplaceTabCharacters
}

// SetReplaceTabCharacters sets whether to replace tab characters.
func (p *ParagraphStyle) SetReplaceTabCharacters(value bool) {
	p.ReplaceTabCharacters = value
}

// GetFakeMissingFontStyles returns whether to fake missing font styles.
func (p *ParagraphStyle) GetFakeMissingFontStyles() bool {
	return p.FakeMissingFontStyles
}

// SetFakeMissingFontStyles sets whether to fake missing font styles.
func (p *ParagraphStyle) SetFakeMissingFontStyles(value bool) {
	p.FakeMissingFontStyles = value
}

// GetApplyRoundingHack returns whether to apply rounding hack.
func (p *ParagraphStyle) GetApplyRoundingHack() bool {
	return p.ApplyRoundingHack
}

// SetApplyRoundingHack sets whether to apply rounding hack.
func (p *ParagraphStyle) SetApplyRoundingHack(value bool) {
	p.ApplyRoundingHack = value
}

// Equals checks for equality between two ParagraphStyles.
func (p *ParagraphStyle) Equals(other *ParagraphStyle) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}

	if !nearlyEqual(p.Height, other.Height) {
		return false
	}

	return p.Ellipsis == other.Ellipsis &&
		p.EllipsisUtf16 == other.EllipsisUtf16 &&
		p.TextDirection == other.TextDirection &&
		p.TextAlign == other.TextAlign &&
		p.DefaultTextStyle.Equals(&other.DefaultTextStyle) &&
		p.ReplaceTabCharacters == other.ReplaceTabCharacters &&
		p.FakeMissingFontStyles == other.FakeMissingFontStyles &&
		p.StrutStyle.Equals(&other.StrutStyle)
}
