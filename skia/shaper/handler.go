package shaper

import (
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// RunHandler is the interface for handling the results of text shaping.
// It maps to SkShaper::RunHandler in C++.
type RunHandler interface {
	// BeginLine starts processing a line.
	BeginLine()

	// RunInfo provides information about the current run.
	RunInfo(info RunInfo)

	// CommitRunInfo commits the run information.
	CommitRunInfo()

	// RunBuffer returns the buffer for the current run.
	RunBuffer(info RunInfo) Buffer

	// CommitRunBuffer commits the run buffer.
	CommitRunBuffer(info RunInfo)

	// CommitLine commits the line.
	CommitLine()
}

// Range represents a range of indices in the text.
type Range struct {
	Begin int
	End   int
}

// RunInfo contains information about a shaped run.
// It maps to SkShaper::RunHandler::RunInfo in C++.
type RunInfo struct {
	Font       interfaces.SkFont
	BidiLevel  uint8
	Script     uint32
	Language   string
	Advance    models.Point
	GlyphCount uint64
	Utf8Range  Range
}

// Buffer contains the shaped glyphs and positions for a run.
// It maps to SkShaper::RunHandler::Buffer in C++.
type Buffer struct {
	Glyphs    []uint16       // Glyph indices
	Positions []models.Point // Glyph positions
	Offsets   []models.Point // Glyph offsets (optional, depending on usage) - wait, C++ might use Point for offsets or just implicit.
	// C++ definition:
	// struct Buffer {
	//     SkGlyphID* glyphs;
	//     SkPoint* positions;
	//     SkPoint* offsets;
	//     uint32_t* clusters;
	//     SkPoint point;
	// };

	// Re-checking requirement: "Offsets ([]uint64)" was in my plan?
	// Plan said: "Offsets ([]uint64), Clusters ([]uint32), Point (Point)."
	// Wait, Offsets being []uint64 seems wrong if they are geometric offsets.
	// In C++, offsets are usually SkPoint (x/y offset from position).
	// But `SkShaper::RunHandler::Buffer` often has `offsets` as `SkPoint*`.
	// My plan said `Offsets ([]uint64)`? That might be a typo in my plan.
	// I will check the C++ source or assume SkPoint for offsets.
	// Clusters are uint32 (indices into text).

	Clusters []uint32
	Point    models.Point // Starting point of the run?
}
