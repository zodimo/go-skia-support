package core

import "github.com/zodimo/go-skia-support/skia/base"

// Color4f represents an RGBA color with four float components (unpremultiplied).
type Color4f struct {
	R, G, B, A base.Scalar
}

// PMColor4f represents an RGBA color with four float components, premultiplied by alpha.
type PMColor4f struct {
	R, G, B, A base.Scalar
}

var (
	Color4fTransparent = Color4f{0, 0, 0, 0}
	Color4fBlack       = Color4f{0, 0, 0, 1}
	Color4fDkGray      = Color4f{0.25, 0.25, 0.25, 1}
	Color4fGray        = Color4f{0.50, 0.50, 0.50, 1}
	Color4fLtGray      = Color4f{0.75, 0.75, 0.75, 1}
	Color4fWhite       = Color4f{1, 1, 1, 1}
	Color4fRed         = Color4f{1, 0, 0, 1}
	Color4fGreen       = Color4f{0, 1, 0, 1}
	Color4fBlue        = Color4f{0, 0, 1, 1}
	Color4fYellow      = Color4f{1, 1, 0, 1}
	Color4fCyan        = Color4f{0, 1, 1, 1}
	Color4fMagenta     = Color4f{1, 0, 1, 1}
)

func scalarPin(x, min, max base.Scalar) base.Scalar {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

// FromColor converts an SkColor (unpremultiplied ARGB) to Color4f.
func Color4fFromColor(c Color) Color4f {
	const inv255 = 1.0 / 255.0
	return Color4f{
		R: base.Scalar(ColorGetR(c)) * inv255,
		G: base.Scalar(ColorGetG(c)) * inv255,
		B: base.Scalar(ColorGetB(c)) * inv255,
		A: base.Scalar(ColorGetA(c)) * inv255,
	}
}

// ToColor converts the Color4f to an SkColor (unpremultiplied ARGB).
// Components are clamped to [0, 1].
func (c Color4f) ToColor() Color {
	r := uint8(scalarPin(c.R, 0, 1)*255.0 + 0.5)
	g := uint8(scalarPin(c.G, 0, 1)*255.0 + 0.5)
	b := uint8(scalarPin(c.B, 0, 1)*255.0 + 0.5)
	a := uint8(scalarPin(c.A, 0, 1)*255.0 + 0.5)
	return ColorARGB(a, r, g, b)
}

// Opacity methods

func (c Color4f) IsOpaque() bool {
	return c.A == 1.0
}

func (c Color4f) FitsInBytes() bool {
	return c.A >= 0 && c.A <= 1 &&
		c.R >= 0 && c.R <= 1 &&
		c.G >= 0 && c.G <= 1 &&
		c.B >= 0 && c.B <= 1
}

func (c Color4f) MakeOpaque() Color4f {
	return Color4f{c.R, c.G, c.B, 1.0}
}

func (c Color4f) PinAlpha() Color4f {
	return Color4f{c.R, c.G, c.B, scalarPin(c.A, 0, 1)}
}

// Premul converts to PMColor4f (premultiplied).
func (c Color4f) Premul() PMColor4f {
	return PMColor4f{c.R * c.A, c.G * c.A, c.B * c.A, c.A}
}

// Unpremul converts to Color4f (unpremultiplied).
func (pm PMColor4f) Unpremul() Color4f {
	if pm.A == 0 {
		return Color4f{0, 0, 0, 0}
	}
	invA := 1.0 / pm.A
	return Color4f{pm.R * invA, pm.G * invA, pm.B * invA, pm.A}
}

// Helpers for array access, useful for interacting with C++ APIs or SIMD later
func (c Color4f) Vec() [4]base.Scalar {
	return [4]base.Scalar{c.R, c.G, c.B, c.A}
}
