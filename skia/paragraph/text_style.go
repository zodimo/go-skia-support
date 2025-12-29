package paragraph

import (
	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// DefaultFontFamily is the default font family used when none is specified.
const DefaultFontFamily = "sans-serif"

// DefaultFontSize is the default font size in points.
const DefaultFontSize = 14.0

// StyleType specifies which attributes to compare when matching text styles.
//
// Ported from: skia-source/modules/skparagraph/include/TextStyle.h
type StyleType int

const (
	// StyleTypeNone matches no attributes.
	StyleTypeNone StyleType = iota

	// StyleTypeAllAttributes matches all attributes.
	StyleTypeAllAttributes

	// StyleTypeFont matches font-related attributes.
	StyleTypeFont

	// StyleTypeForeground matches foreground color/paint.
	StyleTypeForeground

	// StyleTypeBackground matches background color/paint.
	StyleTypeBackground

	// StyleTypeShadow matches shadow attributes.
	StyleTypeShadow

	// StyleTypeDecorations matches decoration attributes.
	StyleTypeDecorations

	// StyleTypeLetterSpacing matches letter spacing.
	StyleTypeLetterSpacing

	// StyleTypeWordSpacing matches word spacing.
	StyleTypeWordSpacing
)

// TextStyle defines all styling options for a text run.
// This is the primary way to configure text appearance in paragraphs.
//
// Ported from: skia-source/modules/skparagraph/include/TextStyle.h
type TextStyle struct {
	// Decoration contains text decoration settings (underline, etc.).
	Decoration Decoration

	// FontStyle defines weight, width, and slant.
	FontStyle models.FontStyle

	// FontFamilies is the list of font family names to use, in priority order.
	FontFamilies []string

	// FontSize is the text size in logical pixels.
	FontSize float32

	// Edging controls how glyph edges are rendered.
	Edging enums.FontEdging

	// Subpixel enables subpixel text positioning.
	Subpixel bool

	// Hinting controls glyph outline adjustment level.
	Hinting enums.FontHinting

	// Height is the line height multiplier.
	Height float32

	// HeightOverride indicates whether Height should be used.
	HeightOverride bool

	// BaselineShift moves the text baseline up (negative) or down (positive).
	BaselineShift float32

	// HalfLeading enables half-leading distribution.
	HalfLeading bool

	// Locale is the locale string for locale-dependent text shaping.
	Locale string

	// LetterSpacing is additional spacing between letters.
	LetterSpacing float32

	// WordSpacing is additional spacing between words.
	WordSpacing float32

	// TextBaseline is the baseline type used for vertical alignment.
	TextBaseline TextBaseline

	// Color is the text color (used when no foreground paint is set).
	Color uint32

	// HasBackground indicates whether a background paint is set.
	HasBackground bool

	// BackgroundPaint is the paint used for text background.
	BackgroundPaint interfaces.SkPaint

	// HasForeground indicates whether a foreground paint is set.
	HasForeground bool

	// ForegroundPaint is the paint used for text foreground.
	ForegroundPaint interfaces.SkPaint

	// TextShadows is the list of shadows applied to the text.
	TextShadows []TextShadow

	// Typeface is an explicit typeface to use (overrides font family lookup).
	Typeface interfaces.SkTypeface

	// IsPlaceholder indicates this style is for a placeholder.
	IsPlaceholder bool

	// FontFeatures is the list of OpenType font features to apply.
	FontFeatures []FontFeature
}

// NewTextStyle creates a new TextStyle with default values.
func NewTextStyle() TextStyle {
	return TextStyle{
		Decoration:     NewDecoration(),
		FontStyle:      models.FontStyle{},
		FontFamilies:   []string{DefaultFontFamily},
		FontSize:       DefaultFontSize,
		Edging:         enums.FontEdgingAntiAlias,
		Subpixel:       true,
		Hinting:        enums.FontHintingSlight,
		Height:         1.0,
		HeightOverride: false,
		BaselineShift:  0.0,
		HalfLeading:    false,
		Locale:         "",
		LetterSpacing:  0.0,
		WordSpacing:    0.0,
		TextBaseline:   TextBaselineAlphabetic,
		Color:          0xFFFFFFFF, // SK_ColorWHITE
		HasBackground:  false,
		HasForeground:  false,
		TextShadows:    nil,
		IsPlaceholder:  false,
		FontFeatures:   nil,
	}
}

// --- Color methods ---

// GetColor returns the text color.
func (s *TextStyle) GetColor() uint32 {
	return s.Color
}

// SetColor sets the text color.
func (s *TextStyle) SetColor(color uint32) {
	s.Color = color
}

// --- Foreground/Background methods ---

// GetForeground returns the foreground paint. Returns nil if not set.
func (s *TextStyle) GetForeground() interfaces.SkPaint {
	if !s.HasForeground {
		return nil
	}
	return s.ForegroundPaint
}

// SetForegroundPaint sets the foreground paint.
func (s *TextStyle) SetForegroundPaint(paint interfaces.SkPaint) {
	s.HasForeground = true
	s.ForegroundPaint = paint
}

// ClearForeground removes the foreground paint.
func (s *TextStyle) ClearForeground() {
	s.HasForeground = false
	s.ForegroundPaint = nil
}

// GetBackground returns the background paint. Returns nil if not set.
func (s *TextStyle) GetBackground() interfaces.SkPaint {
	if !s.HasBackground {
		return nil
	}
	return s.BackgroundPaint
}

// SetBackgroundPaint sets the background paint.
func (s *TextStyle) SetBackgroundPaint(paint interfaces.SkPaint) {
	s.HasBackground = true
	s.BackgroundPaint = paint
}

// ClearBackground removes the background paint.
func (s *TextStyle) ClearBackground() {
	s.HasBackground = false
	s.BackgroundPaint = nil
}

// --- Decoration methods ---

// GetDecoration returns the decoration settings.
func (s *TextStyle) GetDecoration() Decoration {
	return s.Decoration
}

// GetDecorationType returns the decoration type.
func (s *TextStyle) GetDecorationType() TextDecoration {
	return s.Decoration.Type
}

// SetDecoration sets the decoration type.
func (s *TextStyle) SetDecoration(decoration TextDecoration) {
	s.Decoration.Type = decoration
}

// SetDecorationMode sets how decorations interact with descenders.
func (s *TextStyle) SetDecorationMode(mode TextDecorationMode) {
	s.Decoration.Mode = mode
}

// SetDecorationStyle sets the decoration line style.
func (s *TextStyle) SetDecorationStyle(style TextDecorationStyle) {
	s.Decoration.Style = style
}

// SetDecorationColor sets the decoration color.
func (s *TextStyle) SetDecorationColor(color uint32) {
	s.Decoration.Color = color
}

// SetDecorationThicknessMultiplier sets the thickness multiplier.
func (s *TextStyle) SetDecorationThicknessMultiplier(multiplier float32) {
	s.Decoration.ThicknessMultiplier = multiplier
}

// --- Font style methods ---

// GetFontStyle returns the font style (weight, width, slant).
func (s *TextStyle) GetFontStyle() models.FontStyle {
	return s.FontStyle
}

// SetFontStyle sets the font style.
func (s *TextStyle) SetFontStyle(style models.FontStyle) {
	s.FontStyle = style
}

// --- Shadow methods ---

// GetShadowCount returns the number of shadows.
func (s *TextStyle) GetShadowCount() int {
	return len(s.TextShadows)
}

// GetShadows returns a copy of the shadows list.
func (s *TextStyle) GetShadows() []TextShadow {
	if s.TextShadows == nil {
		return nil
	}
	result := make([]TextShadow, len(s.TextShadows))
	copy(result, s.TextShadows)
	return result
}

// AddShadow adds a shadow to the text.
func (s *TextStyle) AddShadow(shadow TextShadow) {
	s.TextShadows = append(s.TextShadows, shadow)
}

// ResetShadows removes all shadows.
func (s *TextStyle) ResetShadows() {
	s.TextShadows = nil
}

// --- Font feature methods ---

// GetFontFeatureCount returns the number of font features.
func (s *TextStyle) GetFontFeatureCount() int {
	return len(s.FontFeatures)
}

// GetFontFeatures returns a copy of the font features list.
func (s *TextStyle) GetFontFeatures() []FontFeature {
	if s.FontFeatures == nil {
		return nil
	}
	result := make([]FontFeature, len(s.FontFeatures))
	copy(result, s.FontFeatures)
	return result
}

// AddFontFeature adds a font feature.
func (s *TextStyle) AddFontFeature(name string, value int) {
	s.FontFeatures = append(s.FontFeatures, NewFontFeature(name, value))
}

// ResetFontFeatures removes all font features.
func (s *TextStyle) ResetFontFeatures() {
	s.FontFeatures = nil
}

// --- Font size methods ---

// GetFontSize returns the font size.
func (s *TextStyle) GetFontSize() float32 {
	return s.FontSize
}

// SetFontSize sets the font size.
func (s *TextStyle) SetFontSize(size float32) {
	s.FontSize = size
}

// --- Font family methods ---

// GetFontFamilies returns the font family list.
func (s *TextStyle) GetFontFamilies() []string {
	return s.FontFamilies
}

// SetFontFamilies sets the font family list.
func (s *TextStyle) SetFontFamilies(families []string) {
	s.FontFamilies = families
}

// --- Spacing methods ---

// GetLetterSpacing returns the letter spacing.
func (s *TextStyle) GetLetterSpacing() float32 {
	return s.LetterSpacing
}

// SetLetterSpacing sets the letter spacing.
func (s *TextStyle) SetLetterSpacing(spacing float32) {
	s.LetterSpacing = spacing
}

// GetWordSpacing returns the word spacing.
func (s *TextStyle) GetWordSpacing() float32 {
	return s.WordSpacing
}

// SetWordSpacing sets the word spacing.
func (s *TextStyle) SetWordSpacing(spacing float32) {
	s.WordSpacing = spacing
}

// --- Height methods ---

// GetHeight returns the height multiplier (0 if HeightOverride is false).
func (s *TextStyle) GetHeight() float32 {
	if s.HeightOverride {
		return s.Height
	}
	return 0
}

// SetHeight sets the height multiplier.
func (s *TextStyle) SetHeight(height float32) {
	s.Height = height
}

// GetHeightOverride returns whether height override is enabled.
func (s *TextStyle) GetHeightOverride() bool {
	return s.HeightOverride
}

// SetHeightOverride enables or disables height override.
func (s *TextStyle) SetHeightOverride(override bool) {
	s.HeightOverride = override
}

// GetHalfLeading returns whether half-leading is enabled.
func (s *TextStyle) GetHalfLeading() bool {
	return s.HalfLeading
}

// SetHalfLeading enables or disables half-leading.
func (s *TextStyle) SetHalfLeading(halfLeading bool) {
	s.HalfLeading = halfLeading
}

// --- Baseline methods ---

// GetBaselineShift returns the baseline shift.
func (s *TextStyle) GetBaselineShift() float32 {
	return s.BaselineShift
}

// SetBaselineShift sets the baseline shift.
func (s *TextStyle) SetBaselineShift(shift float32) {
	s.BaselineShift = shift
}

// GetTextBaseline returns the text baseline type.
func (s *TextStyle) GetTextBaseline() TextBaseline {
	return s.TextBaseline
}

// SetTextBaseline sets the text baseline type.
func (s *TextStyle) SetTextBaseline(baseline TextBaseline) {
	s.TextBaseline = baseline
}

// --- Typeface methods ---

// GetTypeface returns the typeface (may be nil).
func (s *TextStyle) GetTypeface() interfaces.SkTypeface {
	return s.Typeface
}

// SetTypeface sets the typeface.
func (s *TextStyle) SetTypeface(typeface interfaces.SkTypeface) {
	s.Typeface = typeface
}

// --- Locale methods ---

// GetLocale returns the locale string.
func (s *TextStyle) GetLocale() string {
	return s.Locale
}

// SetLocale sets the locale string.
func (s *TextStyle) SetLocale(locale string) {
	s.Locale = locale
}

// --- Font rendering methods ---

// GetFontEdging returns the font edging mode.
func (s *TextStyle) GetFontEdging() enums.FontEdging {
	return s.Edging
}

// SetFontEdging sets the font edging mode.
func (s *TextStyle) SetFontEdging(edging enums.FontEdging) {
	s.Edging = edging
}

// GetSubpixel returns whether subpixel rendering is enabled.
func (s *TextStyle) GetSubpixel() bool {
	return s.Subpixel
}

// SetSubpixel enables or disables subpixel rendering.
func (s *TextStyle) SetSubpixel(subpixel bool) {
	s.Subpixel = subpixel
}

// GetFontHinting returns the font hinting level.
func (s *TextStyle) GetFontHinting() enums.FontHinting {
	return s.Hinting
}

// SetFontHinting sets the font hinting level.
func (s *TextStyle) SetFontHinting(hinting enums.FontHinting) {
	s.Hinting = hinting
}

// --- Placeholder methods ---

// SetPlaceholder marks this style as a placeholder style.
func (s *TextStyle) SetPlaceholder() {
	s.IsPlaceholder = true
}

// --- Comparison methods ---

// CloneForPlaceholder creates a copy of this style suitable for placeholders.
func (s *TextStyle) CloneForPlaceholder() TextStyle {
	clone := *s
	clone.IsPlaceholder = true
	return clone
}

// Equals returns true if this style equals another (all attributes match).
func (s *TextStyle) Equals(other *TextStyle) bool {
	if s == other {
		return true
	}
	if other == nil {
		return false
	}

	// Compare basic fields
	if s.Color != other.Color ||
		s.HasForeground != other.HasForeground ||
		s.HasBackground != other.HasBackground ||
		s.FontSize != other.FontSize ||
		s.Height != other.Height ||
		s.HeightOverride != other.HeightOverride ||
		s.HalfLeading != other.HalfLeading ||
		s.LetterSpacing != other.LetterSpacing ||
		s.WordSpacing != other.WordSpacing ||
		s.BaselineShift != other.BaselineShift ||
		s.TextBaseline != other.TextBaseline ||
		s.Locale != other.Locale ||
		s.Edging != other.Edging ||
		s.Subpixel != other.Subpixel ||
		s.Hinting != other.Hinting ||
		s.IsPlaceholder != other.IsPlaceholder {
		return false
	}

	// Compare decoration
	if !s.Decoration.Equals(other.Decoration) {
		return false
	}

	// Compare font style
	if s.FontStyle != other.FontStyle {
		return false
	}

	// Compare font families
	if len(s.FontFamilies) != len(other.FontFamilies) {
		return false
	}
	for i, f := range s.FontFamilies {
		if f != other.FontFamilies[i] {
			return false
		}
	}

	// Compare shadows
	if len(s.TextShadows) != len(other.TextShadows) {
		return false
	}
	for i, sh := range s.TextShadows {
		if !sh.Equals(other.TextShadows[i]) {
			return false
		}
	}

	// Compare font features
	if len(s.FontFeatures) != len(other.FontFeatures) {
		return false
	}
	for i, ff := range s.FontFeatures {
		if !ff.Equals(other.FontFeatures[i]) {
			return false
		}
	}

	return true
}

// EqualsByFonts returns true if font-related attributes match.
func (s *TextStyle) EqualsByFonts(other *TextStyle) bool {
	if other == nil {
		return false
	}

	if s.FontSize != other.FontSize ||
		s.FontStyle != other.FontStyle ||
		s.Edging != other.Edging ||
		s.Subpixel != other.Subpixel ||
		s.Hinting != other.Hinting ||
		s.Locale != other.Locale {
		return false
	}

	if len(s.FontFamilies) != len(other.FontFamilies) {
		return false
	}
	for i, f := range s.FontFamilies {
		if f != other.FontFamilies[i] {
			return false
		}
	}

	return true
}

// MatchOneAttribute returns true if the specified attribute type matches.
func (s *TextStyle) MatchOneAttribute(styleType StyleType, other *TextStyle) bool {
	if other == nil {
		return false
	}

	switch styleType {
	case StyleTypeNone:
		return true
	case StyleTypeAllAttributes:
		return s.Equals(other)
	case StyleTypeFont:
		return s.EqualsByFonts(other)
	case StyleTypeForeground:
		return s.Color == other.Color &&
			s.HasForeground == other.HasForeground
	case StyleTypeBackground:
		return s.HasBackground == other.HasBackground
	case StyleTypeShadow:
		if len(s.TextShadows) != len(other.TextShadows) {
			return false
		}
		for i, sh := range s.TextShadows {
			if !sh.Equals(other.TextShadows[i]) {
				return false
			}
		}
		return true
	case StyleTypeDecorations:
		return s.Decoration.Equals(other.Decoration)
	case StyleTypeLetterSpacing:
		return s.LetterSpacing == other.LetterSpacing
	case StyleTypeWordSpacing:
		return s.WordSpacing == other.WordSpacing
	default:
		return false
	}
}
