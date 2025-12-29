package paragraph

import (
	"math"
	"testing"
)

func TestParagraphStyleDefaults(t *testing.T) {
	ps := NewParagraphStyle()

	if ps.GetMaxLines() != math.MaxInt {
		t.Errorf("Expected default MaxLines %d, got %d", math.MaxInt, ps.GetMaxLines())
	}
	if !ps.UnlimitedLines() {
		t.Error("Expected default to have UnlimitedLines")
	}
	if ps.GetEllipsis() != "" {
		t.Errorf("Expected empty ellipsis, got %s", ps.GetEllipsis())
	}
	if ps.Ellipsized() {
		t.Error("Expected default to not be Ellipsized")
	}
	if ps.GetTextAlign() != TextAlignLeft {
		t.Errorf("Expected default TextAlign Left, got %v", ps.GetTextAlign())
	}
	if ps.GetTextDirection() != TextDirectionLTR {
		t.Errorf("Expected default TextDirection LTR, got %v", ps.GetTextDirection())
	}
	if ps.GetHeight() != 1.0 {
		t.Errorf("Expected default Height 1.0, got %f", ps.GetHeight())
	}
	if !ps.IsHintingOn() {
		t.Error("Expected default Hinting to be On")
	}
	if ps.GetReplaceTabCharacters() {
		t.Error("Expected default ReplaceTabCharacters to be false")
	}
	if !ps.GetFakeMissingFontStyles() {
		t.Error("Expected default FakeMissingFontStyles to be true")
	}
	if !ps.GetApplyRoundingHack() {
		t.Error("Expected default ApplyRoundingHack to be true")
	}
}

func TestParagraphStyleEffectiveAlign(t *testing.T) {
	ps := NewParagraphStyle()

	// Case 1: Specific alignment
	ps.SetTextAlign(TextAlignCenter)
	if ps.EffectiveAlign() != TextAlignCenter {
		t.Errorf("Expected EffectiveAlign Center, got %v", ps.EffectiveAlign())
	}

	// Case 2: Start alignment (LTR)
	ps.SetTextAlign(TextAlignStart)
	ps.SetTextDirection(TextDirectionLTR)
	if ps.EffectiveAlign() != TextAlignLeft {
		t.Errorf("Expected EffectiveAlign Left (Start+LTR), got %v", ps.EffectiveAlign())
	}

	// Case 3: Start alignment (RTL)
	ps.SetTextDirection(TextDirectionRTL)
	if ps.EffectiveAlign() != TextAlignRight {
		t.Errorf("Expected EffectiveAlign Right (Start+RTL), got %v", ps.EffectiveAlign())
	}

	// Case 4: End alignment (LTR)
	ps.SetTextAlign(TextAlignEnd)
	ps.SetTextDirection(TextDirectionLTR)
	if ps.EffectiveAlign() != TextAlignRight {
		t.Errorf("Expected EffectiveAlign Right (End+LTR), got %v", ps.EffectiveAlign())
	}

	// Case 5: End alignment (RTL)
	ps.SetTextDirection(TextDirectionRTL)
	if ps.EffectiveAlign() != TextAlignLeft {
		t.Errorf("Expected EffectiveAlign Left (End+RTL), got %v", ps.EffectiveAlign())
	}
}

func TestParagraphStyleEquals(t *testing.T) {
	ps1 := NewParagraphStyle()
	ps2 := NewParagraphStyle()

	if !ps1.Equals(&ps2) {
		t.Error("Expected default ParagraphStyles to be equal")
	}

	ps1.SetMaxLines(2)
	// MaxLines is NOT part of Equals check in C++ or our Go implementation?
	// Let's check implementation. Go implementation does NOT include MaxLines in Equals.
	// C++ reference:
	/*
	   bool operator==(const ParagraphStyle& rhs) const {
	       return this->fHeight == rhs.fHeight &&
	              this->fEllipsis == rhs.fEllipsis &&
	              this->fEllipsisUtf16 == rhs.fEllipsisUtf16 &&
	              this->fTextDirection == rhs.fTextDirection && this->fTextAlign == rhs.fTextAlign &&
	              this->fDefaultTextStyle == rhs.fDefaultTextStyle &&
	              this->fReplaceTabCharacters == rhs.fReplaceTabCharacters &&
	              this->fFakeMissingFontStyles == rhs.fFakeMissingFontStyles;
	   }
	*/
	// C++ implementation implies MaxLines (fLinesLimit) is NOT checked.
	// This seems odd but I will follow the port.
	if !ps1.Equals(&ps2) {
		t.Error("Expected ParagraphStyles to stay equal if MaxLines changes (per C++ impl)")
	}

	ps1.SetHeight(2.0)
	if ps1.Equals(&ps2) {
		t.Error("Expected ParagraphStyles with different Height to be unequal")
	}
	ps2.SetHeight(2.0)
	if !ps1.Equals(&ps2) {
		t.Error("Expected paragraph styles to be equal again")
	}

	ps1.TurnHintingOff()
	// HintingIsOn is NOT in Equals?
	// C++: fHintingIsOn is missing from operator==. OK.
	if !ps1.Equals(&ps2) {
		t.Error("Expected ParagraphStyles to stay equal if Hinting changes (per C++ impl)")
	}

	ps1.SetTextAlign(TextAlignRight)
	if ps1.Equals(&ps2) {
		t.Error("Expected ParagraphStyles with different TextAlign to be unequal")
	}
}

func TestParagraphStyleStrutIntegration(t *testing.T) {
	ps := NewParagraphStyle()
	strut := NewStrutStyle()
	strut.SetStrutEnabled(true)

	ps.SetStrutStyle(strut)

	if !ps.GetStrutStyle().GetStrutEnabled() {
		t.Error("Expected StrutStyle to be enabled in ParagraphStyle")
	}
}
