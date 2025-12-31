package core

import (
	"math"

	"github.com/zodimo/go-skia-support/skia/base"
)

// RGBToHSV converts RGB components to HSV.
// hsv[0] is Hue [0 .. 360)
// hsv[1] is Saturation [0 .. 1]
// hsv[2] is Value [0 .. 1]
func RGBToHSV(r, g, b uint8, hsv *[3]base.Scalar) {
	min := r
	if g < min {
		min = g
	}
	if b < min {
		min = b
	}

	max := r
	if g > max {
		max = g
	}
	if b > max {
		max = b
	}

	delta := max - min

	v := base.Scalar(max) / 255.0

	var s base.Scalar
	if max == 0 {
		s = 0
	} else {
		s = base.Scalar(delta) / base.Scalar(max)
	}

	var h base.Scalar
	if delta == 0 {
		h = 0
	} else {
		d := base.Scalar(delta)
		if r == max {
			h = (base.Scalar(g) - base.Scalar(b)) / d
		} else if g == max {
			h = 2.0 + (base.Scalar(b)-base.Scalar(r))/d
		} else {
			h = 4.0 + (base.Scalar(r)-base.Scalar(g))/d
		}

		h *= 60.0
		if h < 0 {
			h += 360.0
		}
	}

	hsv[0] = h
	hsv[1] = s
	hsv[2] = v
}

// ColorToHSV converts a Color to HSV.
func ColorToHSV(c Color, hsv *[3]base.Scalar) {
	RGBToHSV(ColorGetR(c), ColorGetG(c), ColorGetB(c), hsv)
}

// HSVToColor converts HSV components to an ARGB Color.
// Alpha is passed through unchanged.
// hsv[0] is Hue [0 .. 360)
// hsv[1] is Saturation [0 .. 1]
// hsv[2] is Value [0 .. 1]
func HSVToColor(alpha uint8, hsv [3]base.Scalar) Color {
	s := hsv[1]
	v := hsv[2]

	// Pin
	if s < 0 {
		s = 0
	} else if s > 1 {
		s = 1
	}
	if v < 0 {
		v = 0
	} else if v > 1 {
		v = 1
	}

	vByte := uint8(v*255.0 + 0.5)

	if s <= 0 { // Shade of gray
		return ColorARGB(alpha, vByte, vByte, vByte)
	}

	hx := hsv[0]
	if hx < 0 || hx >= 360 {
		hx = 0
	} else {
		hx /= 60.0
	}

	w := float32(math.Floor(float64(hx)))
	f := hx - w

	// Formulas from SkColor.cpp
	// p = (1 - s) * v
	// q = (1 - (s * f)) * v
	// t = (1 - (s * (1 - f))) * v

	p := uint8((1.0-s)*v*255.0 + 0.5)
	q := uint8((1.0-(s*f))*v*255.0 + 0.5)
	t := uint8((1.0-(s*(1.0-f)))*v*255.0 + 0.5)

	var r, g, b uint8
	switch int(w) {
	case 0:
		r, g, b = vByte, t, p
	case 1:
		r, g, b = q, vByte, p
	case 2:
		r, g, b = p, vByte, t
	case 3:
		r, g, b = p, q, vByte
	case 4:
		r, g, b = t, p, vByte
	default: // 5
		r, g, b = vByte, p, q
	}

	return ColorARGB(alpha, r, g, b)
}
