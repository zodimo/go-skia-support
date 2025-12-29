package interfaces

// SkTextBlob combines multiple text runs into an immutable container.
// Each text run consists of glyphs, paint attributes, and position.
// Only parts of paint related to fonts and text rendering are used by run.
//
// Ported from: skia-source/include/core/SkTextBlob.h
type SkTextBlob interface {
	// Bounds returns the conservative bounding box.
	// Uses font associated with each glyph to determine glyph bounds,
	// and unions all bounds. Returned bounds may be larger than the
	// bounds of all glyphs in runs.
	Bounds() Rect

	// UniqueID returns a non-zero value unique among all text blobs.
	UniqueID() uint32
}
