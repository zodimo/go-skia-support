package models

// Color4f represents an RGBA color with four float components (unpremultiplied)
type Color4f struct {
	R, G, B, A Scalar
}

// pinAlpha clamps alpha to [0, 1] range
func (c Color4f) PinAlpha() Color4f {
	return Color4f{
		R: c.R,
		G: c.G,
		B: c.B,
		A: scalarPin(c.A, 0.0, 1.0),
	}
}

// vec returns a pointer to the components array [R, G, B, A]
func (c *Color4f) Vec() *[4]Scalar {
	return &[4]Scalar{c.R, c.G, c.B, c.A}
}

// ToSkColor converts Color4f to uint32 SkColor (ARGB format)
// This is equivalent to SkColor4f::toSkColor() in C++
func (c Color4f) ToSkColor() uint32 {
	// Clamp components to [0, 1] and convert to uint8
	r := uint8(scalarPin(c.R, 0.0, 1.0) * 255.0)
	g := uint8(scalarPin(c.G, 0.0, 1.0) * 255.0)
	b := uint8(scalarPin(c.B, 0.0, 1.0) * 255.0)
	a := uint8(scalarPin(c.A, 0.0, 1.0) * 255.0)
	// Pack as ARGB: (a << 24) | (r << 16) | (g << 8) | b
	return uint32(a)<<24 | uint32(r)<<16 | uint32(g)<<8 | uint32(b)
}

// Local helper, cannot use helpers.ScalarPin because it would create a circular dependency
// scalarPin clamps x between lo and hi, inclusively
// Similar to SkTPin in C++
func scalarPin(x, lo, hi Scalar) Scalar {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}
