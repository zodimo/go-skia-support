package shaper

import (
	"github.com/zodimo/go-skia-support/skia/impl"
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// TextBlobBuilderRunHandler is a helper for shaping text directly into a SkTextBlob.
// It implements the RunHandler interface to receive shaped glyphs from a shaper
// and uses TextBlobBuilder to construct the resulting text blob.
//
// Ported from: skia-source/modules/skshaper/include/SkShaper.h (SkTextBlobBuilderRunHandler)
// Implementation: skia-source/modules/skshaper/src/SkShaper.cpp
type TextBlobBuilderRunHandler struct {
	builder         *impl.TextBlobBuilder
	utf8Text        string
	offset          models.Point
	clusters        []uint32
	clusterOffset   int
	glyphCount      int
	maxRunAscent    models.Scalar
	maxRunDescent   models.Scalar
	maxRunLeading   models.Scalar
	currentPosition models.Point
}

// NewTextBlobBuilderRunHandler creates a new TextBlobBuilderRunHandler.
// utf8Text is the original text being shaped (used for cluster mapping).
// offset is the initial position for the text.
func NewTextBlobBuilderRunHandler(utf8Text string, offset models.Point) *TextBlobBuilderRunHandler {
	return &TextBlobBuilderRunHandler{
		builder:  impl.NewTextBlobBuilder(),
		utf8Text: utf8Text,
		offset:   offset,
	}
}

// BeginLine starts processing a line.
// Called when beginning a line.
func (h *TextBlobBuilderRunHandler) BeginLine() {
	h.currentPosition = h.offset
	h.maxRunAscent = 0
	h.maxRunDescent = 0
	h.maxRunLeading = 0
}

// RunInfo provides information about the current run.
// Called once for each run in a line. Can compute baselines and offsets.
func (h *TextBlobBuilderRunHandler) RunInfo(info RunInfo) {
	// Track the maximum ascent, descent, and leading across all runs in the line.
	// This is used to compute the baseline position.
	//
	// Since SkFont interface doesn't expose GetMetrics, we estimate from font size.
	// Typical font metrics:
	// - Ascent is about -0.8 * size (negative because it's above baseline)
	// - Descent is about 0.2 * size (positive because it's below baseline)
	// - Leading is typically 0 or small positive
	if info.Font != nil {
		size := info.Font.Size()
		// Estimate metrics (ascent is negative, descent is positive in Skia convention)
		estimatedAscent := -size * 0.8
		estimatedDescent := size * 0.2
		estimatedLeading := models.Scalar(0)

		// Track minimums for ascent (negative) and maximums for descent/leading
		if estimatedAscent < h.maxRunAscent {
			h.maxRunAscent = estimatedAscent
		}
		if estimatedDescent > h.maxRunDescent {
			h.maxRunDescent = estimatedDescent
		}
		if estimatedLeading > h.maxRunLeading {
			h.maxRunLeading = estimatedLeading
		}
	}
}

// CommitRunInfo commits the run information.
// Called after all runInfo calls for a line.
func (h *TextBlobBuilderRunHandler) CommitRunInfo() {
	// Adjust Y position by the maximum ascent (ascent is negative for above baseline).
	// This positions the baseline correctly.
	h.currentPosition.Y -= h.maxRunAscent
}

// RunBuffer returns the buffer for the current run.
// Called for each run in a line after commitRunInfo. The buffer will be filled out.
func (h *TextBlobBuilderRunHandler) RunBuffer(info RunInfo) Buffer {
	glyphCount := int(info.GlyphCount)
	if glyphCount <= 0 {
		return Buffer{}
	}

	// Allocate run with full positioning (x, y for each glyph)
	runBuffer := h.builder.AllocRunPos(info.Font, glyphCount)
	if runBuffer == nil {
		return Buffer{}
	}

	h.glyphCount = glyphCount
	h.clusterOffset = info.Utf8Range.Begin

	// Create cluster array for mapping glyphs back to text positions
	h.clusters = make([]uint32, glyphCount)

	// Convert positions from Scalar slice to Point slice
	positions := make([]models.Point, glyphCount)
	for i := 0; i < glyphCount; i++ {
		// Initialize positions to zero; the shaper will fill them in
		positions[i] = models.Point{X: 0, Y: 0}
	}

	// Convert []impl.GlyphID to []uint16 for the Buffer
	glyphs := make([]uint16, glyphCount)

	return Buffer{
		Glyphs:    glyphs,
		Positions: positions,
		Offsets:   nil, // Not used in this implementation
		Clusters:  h.clusters,
		Point:     h.currentPosition,
	}
}

// CommitRunBuffer commits the run buffer.
// Called after each runBuffer is filled out.
func (h *TextBlobBuilderRunHandler) CommitRunBuffer(info RunInfo) {
	// Adjust cluster indices by subtracting the cluster offset.
	// This normalizes the cluster indices to be relative to the run's text range.
	for i := 0; i < h.glyphCount; i++ {
		if int(h.clusters[i]) >= h.clusterOffset {
			h.clusters[i] -= uint32(h.clusterOffset)
		}
	}

	// Advance the current position by the run's advance.
	h.currentPosition.X += info.Advance.X
	h.currentPosition.Y += info.Advance.Y

	// Commit the run to the builder
	h.builder.AddRun()
}

// CommitLine commits the line.
// Called when ending a line.
func (h *TextBlobBuilderRunHandler) CommitLine() {
	// Update offset for the next line.
	// Move down by the line height (descent + leading - ascent).
	h.offset.Y += h.maxRunDescent + h.maxRunLeading - h.maxRunAscent
}

// MakeBlob builds and returns the text blob from the accumulated runs.
// Returns nil if no runs were added.
func (h *TextBlobBuilderRunHandler) MakeBlob() interfaces.SkTextBlob {
	blob := h.builder.Make()
	if blob == nil {
		return nil // Return true nil interface, not nil-valued interface
	}
	return blob
}

// EndPoint returns the current offset position.
// This is useful for determining where to continue placing text.
func (h *TextBlobBuilderRunHandler) EndPoint() models.Point {
	return h.offset
}

// Compile-time interface check
var _ RunHandler = (*TextBlobBuilderRunHandler)(nil)
