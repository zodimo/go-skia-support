package impl

import "github.com/zodimo/go-skia-support/skia/enums"

// PaintBitfields represents packed bitfields for paint flags
// In Go, we use a simple struct instead of bitfields for readability
type PaintBitfields struct {
	AntiAlias bool             // anti-aliasing enabled
	Dither    bool             // dithering enabled
	CapType   enums.PaintCap   // stroke cap type (2 bits)
	JoinType  enums.PaintJoin  // stroke join type (2 bits)
	Style     enums.PaintStyle // paint style (2 bits)
	// Padding: 24 bits unused (not needed in Go struct)
}
