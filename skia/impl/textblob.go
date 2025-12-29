package impl

import (
	"sync/atomic"

	"github.com/zodimo/go-skia-support/skia/enums"
)

// Global unique ID counter for text blobs
var textBlobIDCounter uint32

// nextTextBlobID generates a unique ID for a text blob.
func nextTextBlobID() uint32 {
	return atomic.AddUint32(&textBlobIDCounter, 1)
}

// GlyphID represents a glyph index. This matches C++ SkGlyphID (uint16).
type GlyphID uint16

// TextBlobRun represents a single run of glyphs with the same font.
type TextBlobRun struct {
	Font      SkFont    // Font for this run
	Glyphs    []GlyphID // Glyph indices
	Positions []Point   // Position for each glyph
}

// TextBlob combines multiple text runs into an immutable container.
// Each text run consists of glyphs, font reference, and positions.
//
// Ported from: skia-source/include/core/SkTextBlob.h
type TextBlob struct {
	runs     []TextBlobRun
	bounds   Rect
	uniqueID uint32
}

// Bounds returns the conservative bounding box.
func (tb *TextBlob) Bounds() Rect {
	return tb.bounds
}

// UniqueID returns a non-zero value unique among all text blobs.
func (tb *TextBlob) UniqueID() uint32 {
	return tb.uniqueID
}

// RunCount returns the number of runs in this text blob.
func (tb *TextBlob) RunCount() int {
	return len(tb.runs)
}

// Run returns the run at the given index.
func (tb *TextBlob) Run(index int) *TextBlobRun {
	if index < 0 || index >= len(tb.runs) {
		return nil
	}
	return &tb.runs[index]
}

// MakeTextBlobFromString creates a TextBlob from a string.
// This is a convenience function that uses UTF-8 encoding.
func MakeTextBlobFromString(text string, font SkFont) *TextBlob {
	return MakeTextBlobFromText([]byte(text), enums.TextEncodingUTF8, font)
}

// MakeTextBlobFromText creates a TextBlob from text bytes with the given encoding.
func MakeTextBlobFromText(text []byte, encoding enums.TextEncoding, font SkFont) *TextBlob {
	if len(text) == 0 || font == nil {
		return nil
	}

	// Convert text to glyphs based on encoding
	glyphs := textToGlyphs(text, encoding)
	if len(glyphs) == 0 {
		return nil
	}

	// Calculate positions (simple left-to-right layout)
	positions := calculateGlyphPositions(glyphs, font, 0, 0)

	// Calculate bounds
	bounds := calculateTextBounds(glyphs, positions, font)

	run := TextBlobRun{
		Font:      font,
		Glyphs:    glyphs,
		Positions: positions,
	}

	return &TextBlob{
		runs:     []TextBlobRun{run},
		bounds:   bounds,
		uniqueID: nextTextBlobID(),
	}
}

// textToGlyphs converts text bytes to glyph IDs.
// This is a simplified implementation that maps characters directly to glyph IDs.
// A real implementation would use the typeface's character-to-glyph mapping.
func textToGlyphs(text []byte, encoding enums.TextEncoding) []GlyphID {
	var glyphs []GlyphID

	switch encoding {
	case enums.TextEncodingUTF8:
		// Simplified: treat each byte as a character
		// Real implementation would decode UTF-8 properly
		for _, b := range text {
			glyphs = append(glyphs, GlyphID(b))
		}
	case enums.TextEncodingUTF16:
		// Simplified: each pair of bytes is a character
		for i := 0; i+1 < len(text); i += 2 {
			ch := uint16(text[i]) | (uint16(text[i+1]) << 8)
			glyphs = append(glyphs, GlyphID(ch))
		}
	case enums.TextEncodingUTF32:
		// Simplified: each 4 bytes is a character, truncate to uint16
		for i := 0; i+3 < len(text); i += 4 {
			ch := uint32(text[i]) | (uint32(text[i+1]) << 8) |
				(uint32(text[i+2]) << 16) | (uint32(text[i+3]) << 24)
			glyphs = append(glyphs, GlyphID(ch))
		}
	case enums.TextEncodingGlyphID:
		// Already glyph IDs
		for i := 0; i+1 < len(text); i += 2 {
			gid := uint16(text[i]) | (uint16(text[i+1]) << 8)
			glyphs = append(glyphs, GlyphID(gid))
		}
	}

	return glyphs
}

// calculateGlyphPositions calculates positions for each glyph.
// This is a simplified implementation using estimated advances.
// A real implementation would use actual glyph metrics.
func calculateGlyphPositions(glyphs []GlyphID, font SkFont, startX, startY Scalar) []Point {
	positions := make([]Point, len(glyphs))
	x := startX

	// Estimate advance width as 0.6 * font size (reasonable average for proportional fonts)
	advance := font.Size() * 0.6 * font.ScaleX()

	for i := range glyphs {
		positions[i] = Point{X: x, Y: startY}
		x += advance
	}

	return positions
}

// calculateTextBounds calculates the bounding box for glyphs at their positions.
func calculateTextBounds(glyphs []GlyphID, positions []Point, font SkFont) Rect {
	if len(glyphs) == 0 {
		return Rect{}
	}

	// Estimate bounds based on font size and glyph count
	size := font.Size()
	ascent := size * 0.8  // typical ascent is ~80% of em size
	descent := size * 0.2 // typical descent is ~20% of em size
	advance := size * 0.6 * font.ScaleX()

	// Find min/max positions
	minX := positions[0].X
	maxX := positions[len(positions)-1].X + advance
	minY := positions[0].Y - ascent
	maxY := positions[0].Y + descent

	for _, pos := range positions {
		if pos.X < minX {
			minX = pos.X
		}
		if pos.Y-ascent < minY {
			minY = pos.Y - ascent
		}
		if pos.Y+descent > maxY {
			maxY = pos.Y + descent
		}
	}

	return Rect{
		Left:   minX,
		Top:    minY,
		Right:  maxX,
		Bottom: maxY,
	}
}

// Compile-time interface check
var _ SkTextBlob = (*TextBlob)(nil)
