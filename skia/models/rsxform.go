package models

import "math"

// RSXform is a compressed form of a rotation+scale matrix.
//
// [ SCos     -SSin    Tx ]
// [ SSin      SCos    Ty ]
// [    0         0     1 ]
//
// Ported from SkRSXform.h
// https://github.com/google/skia/blob/main/include/core/SkRSXform.h
type RSXform struct {
	SCos Scalar
	SSin Scalar
	Tx   Scalar
	Ty   Scalar
}

// MakeRSXform creates a new RSXform.
func MakeRSXform(scos, ssin, tx, ty Scalar) RSXform {
	return RSXform{SCos: scos, SSin: ssin, Tx: tx, Ty: ty}
}

// MakeRSXformFromRadians initializes a new xform based on the scale, rotation (in radians),
// final tx,ty location and anchor-point ax,ay within the src quad.
//
// Note: the anchor point is not normalized (e.g. 0...1) but is in pixels of the src image.
func MakeRSXformFromRadians(scale, radians, tx, ty, ax, ay Scalar) RSXform {
	s := Scalar(math.Sin(float64(radians))) * scale
	c := Scalar(math.Cos(float64(radians))) * scale
	return MakeRSXform(c, s, tx+-c*ax+s*ay, ty+-s*ax-c*ay)
}

// RectStaysRect returns true if the rotation is 0, 90, 180, or 270 degrees.
func (r RSXform) RectStaysRect() bool {
	return r.SCos == 0 || r.SSin == 0
}

// SetIdentity sets the xform to identity.
func (r *RSXform) SetIdentity() {
	r.SCos = 1
	r.SSin = 0
	r.Tx = 0
	r.Ty = 0
}

// Set sets the xform values.
func (r *RSXform) Set(scos, ssin, tx, ty Scalar) {
	r.SCos = scos
	r.SSin = ssin
	r.Tx = tx
	r.Ty = ty
}

// ToQuad computes the quad points for the given width and height.
// Ported from SkRSXform.cpp
// https://github.com/google/skia/blob/main/src/core/SkRSXform.cpp
func (r RSXform) ToQuad(width, height Scalar) [4]Point {
	m00 := r.SCos
	m01 := -r.SSin
	m02 := r.Tx
	m10 := -m01
	m11 := m00
	m12 := r.Ty

	return [4]Point{
		{X: m02, Y: m12},
		{X: m00*width + m02, Y: m10*width + m12},
		{X: m00*width + m01*height + m02, Y: m10*width + m11*height + m12},
		{X: m01*height + m02, Y: m11*height + m12},
	}
}

// ToTriStrip computes the triangle strip points for the given width and height.
// Ported from SkRSXform.cpp
// https://github.com/google/skia/blob/main/src/core/SkRSXform.cpp
func (r RSXform) ToTriStrip(width, height Scalar) [4]Point {
	m00 := r.SCos
	m01 := -r.SSin
	m02 := r.Tx
	m10 := -m01
	m11 := m00
	m12 := r.Ty

	return [4]Point{
		{X: m02, Y: m12},
		{X: m01*height + m02, Y: m11*height + m12},
		{X: m00*width + m02, Y: m10*width + m12},
		{X: m00*width + m01*height + m02, Y: m10*width + m11*height + m12},
	}
}
