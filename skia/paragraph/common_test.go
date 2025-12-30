package paragraph

import (
	"errors"

	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// MockTypeface is a mock implementation of interfaces.SkTypeface for testing.
type MockTypeface struct {
	style      models.FontStyle
	familyName string
	uniqueID   uint32
}

func NewMockTypeface(name string, style models.FontStyle) *MockTypeface {
	return &MockTypeface{
		familyName: name,
		style:      style,
		uniqueID:   1, // Simplification
	}
}

func (m *MockTypeface) FontStyle() models.FontStyle {
	return m.style
}

func (m *MockTypeface) IsBold() bool {
	return m.style.Weight >= models.FontWeightBold
}

func (m *MockTypeface) IsItalic() bool {
	return m.style.Slant != models.FontSlantUpright
}

func (m *MockTypeface) IsFixedPitch() bool {
	return false
}

func (m *MockTypeface) UniqueID() uint32 {
	return m.uniqueID
}

func (m *MockTypeface) FamilyName() string {
	return m.familyName
}

func (m *MockTypeface) UnicharToGlyph(unichar rune) uint16 {
	return 1 // Mock: assume all characters supported
}

func (m *MockTypeface) MakeClone(args models.FontArguments) interfaces.SkTypeface {
	// Mock implementation
	return &MockTypeface{
		style:      m.style,
		familyName: m.familyName,
		uniqueID:   m.uniqueID,
	}
}

// --- Glyph Data Access Methods (stubs for interface compliance) ---

func (m *MockTypeface) UnitsPerEm() int {
	return 1000 // Common default for mock
}

func (m *MockTypeface) GetGlyphAdvance(glyphID uint16) int16 {
	return 600 // Reasonable mock advance
}

func (m *MockTypeface) GetGlyphBounds(glyphID uint16) interfaces.Rect {
	return interfaces.Rect{Left: 0, Top: -800, Right: 600, Bottom: 200}
}

func (m *MockTypeface) GetGlyphPath(glyphID uint16) (interfaces.SkPath, error) {
	return nil, errors.New("mock typeface has no glyph paths")
}
