package shaper

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/impl"
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// MockFontMgr is a mock implementation of SkFontMgr for testing.
type MockFontMgr struct {
	matchFamilyStyleCharacterFunc func(familyName string, style models.FontStyle, bcp47 []string, character rune) interfaces.SkTypeface
}

func NewMockFontMgr() *MockFontMgr {
	return &MockFontMgr{}
}

func (m *MockFontMgr) CountFamilies() int                                      { return 0 }
func (m *MockFontMgr) GetFamilyName(index int) string                          { return "" }
func (m *MockFontMgr) CreateStyleSet(index int) interfaces.SkFontStyleSet      { return nil }
func (m *MockFontMgr) MatchFamily(familyName string) interfaces.SkFontStyleSet { return nil }
func (m *MockFontMgr) MatchFamilyStyle(familyName string, style models.FontStyle) interfaces.SkTypeface {
	return nil
}
func (m *MockFontMgr) MatchFamilyStyleCharacter(familyName string, style models.FontStyle, bcp47 []string, character rune) interfaces.SkTypeface {
	if m.matchFamilyStyleCharacterFunc != nil {
		return m.matchFamilyStyleCharacterFunc(familyName, style, bcp47, character)
	}
	return nil
}
func (m *MockFontMgr) MakeFromData(data interfaces.SkData, ttcIndex int) interfaces.SkTypeface {
	return nil
}
func (m *MockFontMgr) MakeFromFile(path string, ttcIndex int) interfaces.SkTypeface { return nil }
func (m *MockFontMgr) LegacyMakeTypeface(familyName string, style models.FontStyle) interfaces.SkTypeface {
	return nil
}

func TestFontMgrRunIterator_Lifecycle(t *testing.T) {
	font := impl.NewFont()
	fontMgr := NewMockFontMgr()
	text := "Hello"

	iter := MakeFontMgrRunIterator(text, font, fontMgr)

	if iter.AtEnd() {
		t.Error("Iterator should not be at end initially for non-empty text")
	}

	iter.Consume()

	if iter.EndOfCurrentRun() == 0 {
		t.Error("EndOfCurrentRun should be > 0 after Consume")
	}

	if iter.CurrentFont() == nil {
		t.Error("CurrentFont should not be nil")
	}
}

func TestFontMgrRunIterator_EmptyText(t *testing.T) {
	font := impl.NewFont()
	fontMgr := NewMockFontMgr()
	text := ""

	iter := MakeFontMgrRunIterator(text, font, fontMgr)

	if !iter.AtEnd() {
		t.Error("Iterator should be at end for empty text")
	}
}

func TestFontMgrRunIterator_NilFontMgr_FallsBackToTrivial(t *testing.T) {
	font := impl.NewFont()
	text := "Hello"

	iter := MakeFontMgrRunIterator(text, font, nil)

	// Should use TrivialFontRunIterator
	_, isTrivial := iter.(*TrivialFontRunIterator)
	if !isTrivial {
		t.Error("Expected TrivialFontRunIterator when fontMgr is nil")
	}
}

func TestFontMgrRunIterator_ConsumesFullText(t *testing.T) {
	font := impl.NewFont()
	fontMgr := NewMockFontMgr()
	text := "Hello World"

	iter := MakeFontMgrRunIterator(text, font, fontMgr)

	// Consume all runs
	for !iter.AtEnd() {
		iter.Consume()
	}

	if iter.EndOfCurrentRun() != len(text) {
		t.Errorf("EndOfCurrentRun should be %d, got %d", len(text), iter.EndOfCurrentRun())
	}
}

func TestFontMgrRunIterator_WithOptions(t *testing.T) {
	font := impl.NewFont()
	fontMgr := NewMockFontMgr()
	text := "Test"
	requestName := "Arial"
	requestStyle := models.FontStyleBold()

	iter := MakeFontMgrRunIteratorWithOptions(text, font, fontMgr, requestName, requestStyle, nil)

	if iter.AtEnd() {
		t.Error("Iterator should not be at end initially")
	}

	iter.Consume()

	if iter.AtEnd() == false && iter.EndOfCurrentRun() == 0 {
		t.Error("Should have advanced after Consume")
	}
}

func TestFontMgrRunIterator_UsesInitialFont(t *testing.T) {
	font := impl.NewFont()
	fontMgr := NewMockFontMgr()
	text := "ABC"

	iter := MakeFontMgrRunIterator(text, font, fontMgr)
	iter.Consume()

	// Since our stub implementation always returns 1 from UnicharToGlyph,
	// the initial font should always be used
	currentFont := iter.CurrentFont()
	if currentFont == nil {
		t.Error("CurrentFont should not be nil")
	}
}

func TestFontMgrRunIterator_UnicodeText(t *testing.T) {
	font := impl.NewFont()
	fontMgr := NewMockFontMgr()
	text := "Hello 世界" // Mixed ASCII and CJK

	iter := MakeFontMgrRunIterator(text, font, fontMgr)

	// Should be able to consume all runs without error
	for !iter.AtEnd() {
		iter.Consume()
	}

	if iter.EndOfCurrentRun() != len(text) {
		t.Errorf("Should have consumed full text: expected %d, got %d", len(text), iter.EndOfCurrentRun())
	}
}

// Compile-time interface check
var _ interfaces.SkFontMgr = (*MockFontMgr)(nil)
