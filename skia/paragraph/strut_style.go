package paragraph

import (
	"math"

	"github.com/zodimo/go-skia-support/skia/models"
)

// StrutStyle defines the strut configuration for a paragraph.
// Strut is a mechanism to force consistent line heights.
//
// Ported from: skia-source/modules/skparagraph/include/ParagraphStyle.h
type StrutStyle struct {
	FontFamilies     []string
	FontStyle        models.FontStyle
	FontSize         float32
	Height           float32
	Leading          float32
	ForceStrutHeight bool
	StrutEnabled     bool
	HeightOverride   bool
	HalfLeading      bool
}

// NewStrutStyle creates a new StrutStyle with default values.
func NewStrutStyle() StrutStyle {
	return StrutStyle{
		FontFamilies:     nil,
		FontStyle:        models.FontStyle{}, // Default regular
		FontSize:         DefaultFontSize,
		Height:           1.0,
		Leading:          0.0,
		ForceStrutHeight: false,
		StrutEnabled:     false,
		HeightOverride:   false,
		HalfLeading:      false,
	}
}

// GetFontFamilies returns the font families.
func (s *StrutStyle) GetFontFamilies() []string {
	return s.FontFamilies
}

// SetFontFamilies sets the font families.
func (s *StrutStyle) SetFontFamilies(families []string) {
	s.FontFamilies = families
}

// GetFontStyle returns the font style.
func (s *StrutStyle) GetFontStyle() models.FontStyle {
	return s.FontStyle
}

// SetFontStyle sets the font style.
func (s *StrutStyle) SetFontStyle(fontStyle models.FontStyle) {
	s.FontStyle = fontStyle
}

// GetFontSize returns the font size.
func (s *StrutStyle) GetFontSize() float32 {
	return s.FontSize
}

// SetFontSize sets the font size.
func (s *StrutStyle) SetFontSize(size float32) {
	s.FontSize = size
}

// GetHeight returns the height multiplier.
func (s *StrutStyle) GetHeight() float32 {
	return s.Height
}

// SetHeight sets the height multiplier.
func (s *StrutStyle) SetHeight(height float32) {
	s.Height = height
}

// GetLeading returns the leading.
func (s *StrutStyle) GetLeading() float32 {
	return s.Leading
}

// SetLeading sets the leading.
func (s *StrutStyle) SetLeading(leading float32) {
	s.Leading = leading
}

// GetStrutEnabled returns whether the strut is enabled.
func (s *StrutStyle) GetStrutEnabled() bool {
	return s.StrutEnabled
}

// SetStrutEnabled enables or disables the strut.
func (s *StrutStyle) SetStrutEnabled(enabled bool) {
	s.StrutEnabled = enabled
}

// GetForceStrutHeight returns whether the strut height is forced.
func (s *StrutStyle) GetForceStrutHeight() bool {
	return s.ForceStrutHeight
}

// SetForceStrutHeight sets whether to force the strut height.
func (s *StrutStyle) SetForceStrutHeight(force bool) {
	s.ForceStrutHeight = force
}

// GetHeightOverride returns whether height override is active.
func (s *StrutStyle) GetHeightOverride() bool {
	return s.HeightOverride
}

// SetHeightOverride sets whether height override is active.
func (s *StrutStyle) SetHeightOverride(override bool) {
	s.HeightOverride = override
}

// GetHalfLeading returns whether half leading is enabled.
func (s *StrutStyle) GetHalfLeading() bool {
	return s.HalfLeading
}

// SetHalfLeading sets whether half leading is enabled.
func (s *StrutStyle) SetHalfLeading(halfLeading bool) {
	s.HalfLeading = halfLeading
}

// Equals checks for equality between two StrutStyles.
func (s *StrutStyle) Equals(other *StrutStyle) bool {
	if s == other {
		return true
	}
	if other == nil {
		return false
	}
	if s.StrutEnabled != other.StrutEnabled ||
		s.ForceStrutHeight != other.ForceStrutHeight ||
		s.HeightOverride != other.HeightOverride ||
		s.HalfLeading != other.HalfLeading ||
		s.FontStyle != other.FontStyle {
		return false
	}

	if !nearlyEqual(s.FontSize, other.FontSize) ||
		!nearlyEqual(s.Height, other.Height) ||
		!nearlyEqual(s.Leading, other.Leading) {
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

// nearlyEqual compares two float32 values with a small epsilon.
func nearlyEqual(a, b float32) bool {
	const epsilon = 1e-5 // Small epsilon for float comparison
	return math.Abs(float64(a-b)) < epsilon
}
