package core

// Ported from SkColor.h
// https://github.com/google/skia/blob/main/include/core/SkColor.h

type Alpha uint8
type Color uint32
type PMColor uint32

const (
	AlphaTransparent Alpha = 0x00
	AlphaOpaque      Alpha = 0xFF

	ColorTransparent Color = 0x00000000
	ColorBlack       Color = 0xFF000000
	ColorDkGray      Color = 0xFF444444
	ColorGray        Color = 0xFF888888
	ColorLtGray      Color = 0xFFCCCCCC
	ColorWhite       Color = 0xFFFFFFFF
	ColorRed         Color = 0xFFFF0000
	ColorGreen       Color = 0xFF00FF00
	ColorBlue        Color = 0xFF0000FF
	ColorYellow      Color = 0xFFFFFF00
	ColorCyan        Color = 0xFF00FFFF
	ColorMagenta     Color = 0xFFFF00FF
)

func ColorARGB(a, r, g, b uint8) Color {
	return (Color(a) << 24) | (Color(r) << 16) | (Color(g) << 8) | Color(b)
}

func ColorRGBA(r, g, b, a uint8) Color {
	return (Color(a) << 24) | (Color(r) << 16) | (Color(g) << 8) | Color(b)
}

func ColorRGB(r, g, b uint8) Color {
	return ColorARGB(0xFF, r, g, b)
}

func ColorGetA(c Color) uint8 { return uint8((c >> 24) & 0xFF) }
func ColorGetR(c Color) uint8 { return uint8((c >> 16) & 0xFF) }
func ColorGetG(c Color) uint8 { return uint8((c >> 8) & 0xFF) }
func ColorGetB(c Color) uint8 { return uint8((c >> 0) & 0xFF) }

func ColorSetA(c Color, a uint8) Color {
	return (c & 0x00FFFFFF) | (Color(a) << 24)
}

// Ported from SkColorPriv.h
// https://github.com/google/skia/blob/main/include/core/SkColorPriv.h

func mulDiv255Round(a, b uint8) uint8 {
	prod := uint32(a)*uint32(b) + 128
	return uint8((prod + (prod >> 8)) >> 8)
}

func PreMultiplyARGB(a, r, g, b uint8) PMColor {
	if a != 255 {
		r = mulDiv255Round(r, a)
		g = mulDiv255Round(g, a)
		b = mulDiv255Round(b, a)
	}
	// Using same packing as Color (ARGB)
	return (PMColor(a) << 24) | (PMColor(r) << 16) | (PMColor(g) << 8) | PMColor(b)
}

func PreMultiplyColor(c Color) PMColor {
	return PreMultiplyARGB(ColorGetA(c), ColorGetR(c), ColorGetG(c), ColorGetB(c))
}
