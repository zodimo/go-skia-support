package paragraph

import (
	"math"
	"testing"

	"github.com/zodimo/go-skia-support/skia/models"
)

func TestMetricsInitialization(t *testing.T) {
	style := NewTextStyle()
	fontMetrics := models.FontMetrics{
		Ascent:  -10.0,
		Descent: 3.0,
		Leading: 1.0,
	}

	sm := NewStyleMetrics(&style, fontMetrics)

	if sm.TextStyle != &style {
		t.Error("StyleMetrics expected to hold reference to TextStyle")
	}
	if sm.FontMetrics.Ascent != -10.0 {
		t.Error("StyleMetrics expected to hold FontMetrics")
	}
}

func TestLineMetricsDefaults(t *testing.T) {
	lm := NewLineMetrics()

	if lm.Ascent != math.MaxFloat64 {
		t.Errorf("Expected Ascent MaxFloat64, got %f", lm.Ascent)
	}
	if lm.Descent != -math.MaxFloat64 {
		t.Errorf("Expected Descent -MaxFloat64, got %f", lm.Descent)
	}
	if lm.Height != 0.0 {
		t.Errorf("Expected Height 0.0, got %f", lm.Height)
	}
	if len(lm.LineMetrics) != 0 {
		t.Error("Expected empty LineMetrics map")
	}
}

func TestFontStyleHelpers(t *testing.T) {
	// Verify FontMetrics helper methods
	fm := models.FontMetrics{
		Flags:              models.FontMetricsUnderlineThicknessIsValidFlag,
		UnderlineThickness: 1.5,
	}

	valid, thickness := fm.HasUnderlineThickness()
	if !valid {
		t.Error("Expected valid underline thickness")
	}
	if thickness != 1.5 {
		t.Errorf("Expected thickness 1.5, got %f", thickness)
	}

	valid, _ = fm.HasUnderlinePosition()
	if valid {
		t.Error("Expected invalid underline position")
	}
}
