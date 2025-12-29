package impl

// RunBuffer supplies storage for glyphs and positions within a run.
// A run is a sequence of glyphs sharing font metrics and positioning.
//
// Ported from: skia-source/include/core/SkTextBlob.h (SkTextBlobBuilder::RunBuffer)
type RunBuffer struct {
	Glyphs    []GlyphID // storage for glyph indexes in run
	Positions []Scalar  // storage for glyph positions in run
}

// Points returns the positions as a slice of Points (for 2D positioning).
func (rb *RunBuffer) Points() []Point {
	if len(rb.Positions) < 2 {
		return nil
	}
	count := len(rb.Positions) / 2
	points := make([]Point, count)
	for i := 0; i < count; i++ {
		points[i] = Point{
			X: rb.Positions[i*2],
			Y: rb.Positions[i*2+1],
		}
	}
	return points
}

// TextBlobBuilder is a helper class for constructing TextBlob.
//
// Ported from: skia-source/include/core/SkTextBlob.h (SkTextBlobBuilder)
type TextBlobBuilder struct {
	runs           []TextBlobRun
	bounds         Rect
	currentBuffer  RunBuffer
	currentFont    SkFont
	currentOffset  Point
	runCount       int
	deferredBounds bool
}

// NewTextBlobBuilder creates a new empty TextBlobBuilder.
func NewTextBlobBuilder() *TextBlobBuilder {
	return &TextBlobBuilder{
		runs:           nil,
		bounds:         Rect{},
		deferredBounds: true,
	}
}

// AllocRun returns run with storage for glyphs.
// Caller must write count glyphs to RunBuffer.Glyphs before next call.
// Glyphs are positioned on a baseline at (x, y), using font metrics to
// determine their relative placement.
func (b *TextBlobBuilder) AllocRun(font SkFont, count int, x, y Scalar) *RunBuffer {
	if count <= 0 || font == nil {
		return nil
	}

	b.currentFont = font
	b.currentOffset = Point{X: x, Y: y}
	b.currentBuffer = RunBuffer{
		Glyphs:    make([]GlyphID, count),
		Positions: nil, // Not used for default run
	}

	return &b.currentBuffer
}

// AllocRunPosH returns run with storage for glyphs and positions along baseline.
// Caller must write count glyphs to RunBuffer.Glyphs and count scalars to
// RunBuffer.Positions before next call.
// Glyphs are positioned on a baseline at y, using x-axis positions
// written by caller to RunBuffer.Positions.
func (b *TextBlobBuilder) AllocRunPosH(font SkFont, count int, y Scalar) *RunBuffer {
	if count <= 0 || font == nil {
		return nil
	}

	b.currentFont = font
	b.currentOffset = Point{X: 0, Y: y}
	b.currentBuffer = RunBuffer{
		Glyphs:    make([]GlyphID, count),
		Positions: make([]Scalar, count), // X positions only
	}

	return &b.currentBuffer
}

// AllocRunPos returns run with storage for glyphs and Point positions.
// Caller must write count glyphs to RunBuffer.Glyphs and count*2 scalars to
// RunBuffer.Positions (x,y pairs) before next call.
func (b *TextBlobBuilder) AllocRunPos(font SkFont, count int) *RunBuffer {
	if count <= 0 || font == nil {
		return nil
	}

	b.currentFont = font
	b.currentOffset = Point{X: 0, Y: 0}
	b.currentBuffer = RunBuffer{
		Glyphs:    make([]GlyphID, count),
		Positions: make([]Scalar, count*2), // X,Y pairs
	}

	return &b.currentBuffer
}

// AddRun commits the current run buffer to the builder.
// Call this after filling the RunBuffer returned by AllocRun/AllocRunPosH/AllocRunPos.
func (b *TextBlobBuilder) AddRun() {
	if len(b.currentBuffer.Glyphs) == 0 || b.currentFont == nil {
		return
	}

	// Convert buffer to run
	run := TextBlobRun{
		Font:   b.currentFont,
		Glyphs: make([]GlyphID, len(b.currentBuffer.Glyphs)),
	}
	copy(run.Glyphs, b.currentBuffer.Glyphs)

	// Calculate positions
	if len(b.currentBuffer.Positions) == 0 {
		// Default positioning (AllocRun case)
		run.Positions = calculateGlyphPositions(run.Glyphs, b.currentFont, b.currentOffset.X, b.currentOffset.Y)
	} else if len(b.currentBuffer.Positions) == len(b.currentBuffer.Glyphs) {
		// Horizontal positioning (AllocRunPosH case)
		run.Positions = make([]Point, len(b.currentBuffer.Glyphs))
		for i, xPos := range b.currentBuffer.Positions {
			run.Positions[i] = Point{X: xPos, Y: b.currentOffset.Y}
		}
	} else if len(b.currentBuffer.Positions) == len(b.currentBuffer.Glyphs)*2 {
		// Full positioning (AllocRunPos case)
		run.Positions = b.currentBuffer.Points()
	}

	b.runs = append(b.runs, run)
	b.runCount++
	b.deferredBounds = true

	// Clear current buffer
	b.currentBuffer = RunBuffer{}
	b.currentFont = nil
}

// Make returns TextBlob built from runs of glyphs added by builder.
// Returned TextBlob is immutable; it may be copied, but its contents may not be altered.
// Returns nil if no runs of glyphs were added by builder.
// Resets TextBlobBuilder to its initial empty state, allowing it to be
// reused to build a new set of runs.
func (b *TextBlobBuilder) Make() *TextBlob {
	// Commit any pending run
	if len(b.currentBuffer.Glyphs) > 0 {
		b.AddRun()
	}

	if len(b.runs) == 0 {
		return nil
	}

	// Calculate total bounds
	bounds := b.computeBounds()

	// Create the text blob
	blob := &TextBlob{
		runs:     b.runs,
		bounds:   bounds,
		uniqueID: nextTextBlobID(),
	}

	// Reset builder
	b.runs = nil
	b.bounds = Rect{}
	b.runCount = 0
	b.deferredBounds = true

	return blob
}

// computeBounds calculates the union of all run bounds.
func (b *TextBlobBuilder) computeBounds() Rect {
	if len(b.runs) == 0 {
		return Rect{}
	}

	var totalBounds Rect
	first := true

	for _, run := range b.runs {
		runBounds := calculateTextBounds(run.Glyphs, run.Positions, run.Font)

		if first {
			totalBounds = runBounds
			first = false
		} else {
			// Union bounds
			if runBounds.Left < totalBounds.Left {
				totalBounds.Left = runBounds.Left
			}
			if runBounds.Top < totalBounds.Top {
				totalBounds.Top = runBounds.Top
			}
			if runBounds.Right > totalBounds.Right {
				totalBounds.Right = runBounds.Right
			}
			if runBounds.Bottom > totalBounds.Bottom {
				totalBounds.Bottom = runBounds.Bottom
			}
		}
	}

	return totalBounds
}
