package core

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/base"
)

func TestColorSetARGB(t *testing.T) {
	c := ColorARGB(0xFF, 0x11, 0x22, 0x33)
	if c != 0xFF112233 {
		t.Errorf("Expected 0xFF112233, got %X", c)
	}
	if ColorGetA(c) != 0xFF {
		t.Error("Bad A")
	}
	if ColorGetR(c) != 0x11 {
		t.Error("Bad R")
	}
	if ColorGetG(c) != 0x22 {
		t.Error("Bad G")
	}
	if ColorGetB(c) != 0x33 {
		t.Error("Bad B")
	}
}

func TestColorSetRGB(t *testing.T) {
	c := ColorRGB(0xAA, 0xBB, 0xCC)
	if c != 0xFFAABBCC {
		t.Errorf("Expected 0xFFAABBCC, got %X", c)
	}
}

func TestPremultiply(t *testing.T) {
	// 50% alpha, red
	// a=128, r=255, g=0, b=0
	// r_premul = (255 * 128 + 128) / 255 = 128 (approx logic check)
	// SkMulDiv255Round(255, 128) -> (255*128 + 128) -> ...
	// 255*128 + 128 = 32768.
	// (32768 + (32768>>8)) >> 8 = (32768 + 128) >> 8 = 32896 >> 8 = 128.
	// Correct.

	pm := PreMultiplyARGB(128, 255, 0, 0)
	a := uint8(pm >> 24)
	r := uint8(pm >> 16)

	if a != 128 {
		t.Errorf("Expected A=128 got %d", a)
	}
	if r != 128 {
		t.Errorf("Expected R=128 got %d", r)
	}

	// Test opaque
	pm2 := PreMultiplyARGB(255, 100, 100, 100)
	if pm2 != 0xFF646464 {
		t.Errorf("Expected 0xFF646464 got %X", pm2)
	}
}

func diff(a, b base.Scalar) base.Scalar {
	d := a - b
	if d < 0 {
		return -d
	}
	return d
}

func intDiff(a, b uint8) int {
	if a > b {
		return int(a - b)
	}
	return int(b - a)
}

func TestHSV(t *testing.T) {
	tests := []struct {
		name    string
		r, g, b uint8
		h, s, v base.Scalar // approximate
	}{
		{"Red", 255, 0, 0, 0, 1, 1},
		{"Green", 0, 255, 0, 120, 1, 1},
		{"Blue", 0, 0, 255, 240, 1, 1},
		{"Cyan", 0, 255, 255, 180, 1, 1},
		{"White", 255, 255, 255, 0, 0, 1},
		{"Black", 0, 0, 0, 0, 0, 0},
	}

	for _, tt := range tests {
		var hsv [3]base.Scalar
		RGBToHSV(tt.r, tt.g, tt.b, &hsv)

		if diff(hsv[0], tt.h) > 0.1 || diff(hsv[1], tt.s) > 0.1 || diff(hsv[2], tt.v) > 0.1 {
			t.Errorf("%s: RGBToHSV got %v, expected %v", tt.name, hsv, []base.Scalar{tt.h, tt.s, tt.v})
		}

		// Round trip
		c := HSVToColor(255, hsv)
		r := ColorGetR(c)
		g := ColorGetG(c)
		b := ColorGetB(c)

		// Allow off-by-one or two due to rounding
		if intDiff(r, tt.r) > 2 || intDiff(g, tt.g) > 2 || intDiff(b, tt.b) > 2 {
			t.Errorf("%s: Round trip HSVToColor got %d,%d,%d expected %d,%d,%d", tt.name, r, g, b, tt.r, tt.g, tt.b)
		}
	}
}

func TestColor4f(t *testing.T) {
	c := ColorARGB(255, 255, 255, 255)
	c4 := Color4fFromColor(c)
	if c4.R != 1 {
		t.Error("R!=1")
	}

	cBack := c4.ToColor()
	if cBack != c {
		t.Errorf("Round trip failed: %X != %X", cBack, c)
	}

	cTransparent := Color4fTransparent
	if cTransparent.A != 0 {
		t.Error("Transparent A!=0")
	}

	cBlue := Color4fBlue
	if cBlue.B != 1 {
		t.Error("Blue B!=1")
	}

	// Premul
	// R=1, A=0.5 -> R=0.5
	cp := Color4f{1, 0, 0, 0.5}
	pm := cp.Premul()
	if pm.R != 0.5 || pm.A != 0.5 {
		t.Errorf("Premul failed: %v", pm)
	}

	// Unpremul
	up := pm.Unpremul()
	if up.R != 1 || up.A != 0.5 {
		t.Errorf("Unpremul failed: %v", up)
	}
}
