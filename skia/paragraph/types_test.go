package paragraph

import (
	"testing"
)

func TestAffinity(t *testing.T) {
	tests := []struct {
		name     string
		affinity Affinity
		expected int
	}{
		{"Upstream", AffinityUpstream, 0},
		{"Downstream", AffinityDownstream, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.affinity) != tt.expected {
				t.Errorf("Affinity %s = %d, want %d", tt.name, tt.affinity, tt.expected)
			}
		})
	}
}

func TestRectHeightStyle(t *testing.T) {
	tests := []struct {
		name  string
		style RectHeightStyle
		value int
	}{
		{"Tight", RectHeightStyleTight, 0},
		{"Max", RectHeightStyleMax, 1},
		{"IncludeLineSpacingMiddle", RectHeightStyleIncludeLineSpacingMiddle, 2},
		{"IncludeLineSpacingTop", RectHeightStyleIncludeLineSpacingTop, 3},
		{"IncludeLineSpacingBottom", RectHeightStyleIncludeLineSpacingBottom, 4},
		{"Strut", RectHeightStyleStrut, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.style) != tt.value {
				t.Errorf("RectHeightStyle %s = %d, want %d", tt.name, tt.style, tt.value)
			}
		})
	}
}

func TestRectWidthStyle(t *testing.T) {
	if RectWidthStyleTight != 0 {
		t.Errorf("RectWidthStyleTight = %d, want 0", RectWidthStyleTight)
	}
	if RectWidthStyleMax != 1 {
		t.Errorf("RectWidthStyleMax = %d, want 1", RectWidthStyleMax)
	}
}

func TestTextAlign(t *testing.T) {
	tests := []struct {
		name  string
		align TextAlign
		value int
	}{
		{"Left", TextAlignLeft, 0},
		{"Right", TextAlignRight, 1},
		{"Center", TextAlignCenter, 2},
		{"Justify", TextAlignJustify, 3},
		{"Start", TextAlignStart, 4},
		{"End", TextAlignEnd, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.align) != tt.value {
				t.Errorf("TextAlign %s = %d, want %d", tt.name, tt.align, tt.value)
			}
		})
	}
}

func TestTextDirection(t *testing.T) {
	if TextDirectionRTL != 0 {
		t.Errorf("TextDirectionRTL = %d, want 0", TextDirectionRTL)
	}
	if TextDirectionLTR != 1 {
		t.Errorf("TextDirectionLTR = %d, want 1", TextDirectionLTR)
	}
}

func TestTextBaseline(t *testing.T) {
	if TextBaselineAlphabetic != 0 {
		t.Errorf("TextBaselineAlphabetic = %d, want 0", TextBaselineAlphabetic)
	}
	if TextBaselineIdeographic != 1 {
		t.Errorf("TextBaselineIdeographic = %d, want 1", TextBaselineIdeographic)
	}
}

func TestTextHeightBehavior(t *testing.T) {
	tests := []struct {
		name     string
		behavior TextHeightBehavior
		value    int
	}{
		{"All", TextHeightBehaviorAll, 0x0},
		{"DisableFirstAscent", TextHeightBehaviorDisableFirstAscent, 0x1},
		{"DisableLastDescent", TextHeightBehaviorDisableLastDescent, 0x2},
		{"DisableAll", TextHeightBehaviorDisableAll, 0x3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.behavior) != tt.value {
				t.Errorf("TextHeightBehavior %s = %d, want %d", tt.name, tt.behavior, tt.value)
			}
		})
	}

	// Test that DisableAll is the combination of first and last
	combined := TextHeightBehaviorDisableFirstAscent | TextHeightBehaviorDisableLastDescent
	if combined != TextHeightBehaviorDisableAll {
		t.Errorf("DisableFirstAscent | DisableLastDescent = %d, want %d", combined, TextHeightBehaviorDisableAll)
	}
}

func TestLineMetricStyle(t *testing.T) {
	if LineMetricStyleTypographic != 0 {
		t.Errorf("LineMetricStyleTypographic = %d, want 0", LineMetricStyleTypographic)
	}
	if LineMetricStyleCSS != 1 {
		t.Errorf("LineMetricStyleCSS = %d, want 1", LineMetricStyleCSS)
	}
}
