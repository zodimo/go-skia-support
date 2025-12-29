package paragraph

import (
	"math"

	"github.com/zodimo/go-skia-support/skia/interfaces"
)

// InternalLineMetrics tracks line metrics during layout.
//
// Ported from: skia-source/modules/skparagraph/src/Run.h (InternalLineMetrics class)
type InternalLineMetrics struct {
	Ascent  float32
	Descent float32
	Leading float32

	RawAscent  float32
	RawDescent float32
	RawLeading float32

	ForceStrut bool
}

// NewInternalLineMetrics creates a new InternalLineMetrics with default values.
func NewInternalLineMetrics() InternalLineMetrics {
	return InternalLineMetrics{
		Ascent:     math.MaxFloat32,
		Descent:    -math.MaxFloat32,
		Leading:    0,
		RawAscent:  math.MaxFloat32,
		RawDescent: -math.MaxFloat32,
		RawLeading: 0,
		ForceStrut: false,
	}
}

// NewInternalLineMetricsFromValues creates a new InternalLineMetrics with specific values.
func NewInternalLineMetricsFromValues(a, d, l float32) InternalLineMetrics {
	return InternalLineMetrics{
		Ascent:     a,
		Descent:    d,
		Leading:    l,
		RawAscent:  a,
		RawDescent: d,
		RawLeading: l,
		ForceStrut: false,
	}
}

// NewInternalLineMetricsFromFont creates a new InternalLineMetrics from a font.
func NewInternalLineMetricsFromFont(font interfaces.SkFont, forceStrut bool) InternalLineMetrics {
	// TODO: Get actual metrics from font. For now using estimation or stub.
	// In strict port, we need font.getMetrics().
	// See run.go helper `getFontMetrics`
	metrics := getFontMetrics(font) // Defined in run.go

	return InternalLineMetrics{
		Ascent:     float32(metrics.Ascent),
		Descent:    float32(metrics.Descent),
		Leading:    float32(metrics.Leading),
		RawAscent:  float32(metrics.Ascent),
		RawDescent: float32(metrics.Descent),
		RawLeading: float32(metrics.Leading),
		ForceStrut: forceStrut,
	}
}

// AddRun updates metrics with a run's metrics.
func (ilm *InternalLineMetrics) AddRun(run *Run) {
	if ilm.ForceStrut {
		return
	}
	ilm.Ascent = minScalar(ilm.Ascent, run.CorrectAscent())
	ilm.Descent = maxScalar(ilm.Descent, run.CorrectDescent())
	ilm.Leading = maxScalar(ilm.Leading, run.CorrectLeading())

	ilm.RawAscent = minScalar(ilm.RawAscent, run.Ascent())
	ilm.RawDescent = maxScalar(ilm.RawDescent, run.Descent())
	ilm.RawLeading = maxScalar(ilm.RawLeading, run.Leading())
}

// Add updates metrics with another metrics object.
func (ilm *InternalLineMetrics) Add(other InternalLineMetrics) {
	ilm.Ascent = minScalar(ilm.Ascent, other.Ascent)
	ilm.Descent = maxScalar(ilm.Descent, other.Descent)
	ilm.Leading = maxScalar(ilm.Leading, other.Leading)
	ilm.RawAscent = minScalar(ilm.RawAscent, other.RawAscent)
	ilm.RawDescent = maxScalar(ilm.RawDescent, other.RawDescent)
	ilm.RawLeading = maxScalar(ilm.RawLeading, other.RawLeading)
}

// Clean resets the metrics to initial state.
func (ilm *InternalLineMetrics) Clean() {
	ilm.Ascent = math.MaxFloat32
	ilm.Descent = -math.MaxFloat32
	ilm.Leading = 0
	ilm.RawAscent = math.MaxFloat32
	ilm.RawDescent = -math.MaxFloat32
	ilm.RawLeading = 0
}

// IsClean checks if the metrics are in initial state.
func (ilm *InternalLineMetrics) IsClean() bool {
	return ilm.Ascent == math.MaxFloat32 &&
		ilm.Descent == -math.MaxFloat32 &&
		ilm.Leading == 0 &&
		ilm.RawAscent == math.MaxFloat32 &&
		ilm.RawDescent == -math.MaxFloat32 &&
		ilm.RawLeading == 0
}

// Delta returns the delta between height and ideographic baseline.
func (ilm *InternalLineMetrics) Delta() float32 {
	return ilm.Height() - ilm.IdeographicBaseline()
}

// RunTop calculates the top position for a run.
func (ilm *InternalLineMetrics) RunTop(run *Run, ascentStyle LineMetricStyle) float32 {
	ascent := run.CorrectAscent()
	if ascentStyle == LineMetricStyleTypographic {
		ascent = run.Ascent()
	}
	// Formula: fLeading / 2 - fAscent + (styleAscent) + delta
	return ilm.Leading/2 - ilm.Ascent + ascent + ilm.Delta()
}

// Height returns the total line height.
func (ilm *InternalLineMetrics) Height() float32 {
	return float32(math.Round(float64(ilm.Descent - ilm.Ascent + ilm.Leading)))
}

// UpdateLineMetrics updates the target metrics based on this metrics (and force strut).
func (ilm *InternalLineMetrics) UpdateLineMetrics(metrics *InternalLineMetrics) {
	if metrics.ForceStrut {
		metrics.Ascent = ilm.Ascent
		metrics.Descent = ilm.Descent
		metrics.Leading = ilm.Leading
		metrics.RawAscent = ilm.RawAscent
		metrics.RawDescent = ilm.RawDescent
		metrics.RawLeading = ilm.RawLeading
	} else {
		metrics.Ascent = minScalar(metrics.Ascent, ilm.Ascent-ilm.Leading/2.0)
		metrics.Descent = maxScalar(metrics.Descent, ilm.Descent+ilm.Leading/2.0)
		metrics.RawAscent = minScalar(metrics.RawAscent, ilm.RawAscent-ilm.RawLeading/2.0)
		metrics.RawDescent = maxScalar(metrics.RawDescent, ilm.RawDescent+ilm.RawLeading/2.0)
	}
}

// Baselines
func (ilm *InternalLineMetrics) AlphabeticBaseline() float32 {
	return ilm.Leading/2 - ilm.Ascent
}

func (ilm *InternalLineMetrics) IdeographicBaseline() float32 {
	return ilm.Descent - ilm.Ascent + ilm.Leading
}

func (ilm *InternalLineMetrics) Baseline() float32 {
	return ilm.AlphabeticBaseline()
}

// Helpers
func minScalar(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func maxScalar(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
