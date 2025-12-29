package paragraph

import (
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
