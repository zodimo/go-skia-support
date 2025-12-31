package models

import (
	"math"
	"testing"
)

func floatEquals(a, b Scalar) bool {
	return math.Abs(float64(a-b)) < 0.00001
}

func TestRSXform_MakeFromRadians(t *testing.T) {
	// Identity: scale=1, rad=0, tx=0, ty=0, center=0,0
	r := MakeRSXformFromRadians(1, 0, 0, 0, 0, 0)
	if !floatEquals(r.SCos, 1) || !floatEquals(r.SSin, 0) || !floatEquals(r.Tx, 0) || !floatEquals(r.Ty, 0) {
		t.Errorf("MakeFromRadians(1, 0, ...) = %v, want Identity", r)
	}

	// Rotate 90 degrees around 0,0
	rad := Scalar(math.Pi / 2)
	r = MakeRSXformFromRadians(1, rad, 10, 20, 0, 0)
	// cos(90) = 0, sin(90) = 1
	// tx = 10, ty = 20
	// We allow some epsilon error for PI/2 calculation
	if math.Abs(float64(r.SCos)) > 0.00001 || math.Abs(float64(r.SSin)-1) > 0.00001 || !floatEquals(r.Tx, 10) || !floatEquals(r.Ty, 20) {
		t.Errorf("MakeFromRadians(90 deg) = %v, want SCos~0, SSin~1, Tx=10, Ty=20", r)
	}

	// Rotate 90 degrees around 10, 10 (anchor) at 0, 0 (tx, ty for anchor)
	// SkRSXform::MakeFromRadians logic:
	// s = sin(rad)*scale
	// c = cos(rad)*scale
	// tx' = tx + -c*ax + s*ay
	// ty' = ty + -s*ax - c*ay
	//
	// Here: scale=1, rad=90 => s=1, c=0
	// tx=0, ty=0, ax=10, ay=10
	// tx' = 0 + 0 + 1*10 = 10
	// ty' = 0 - 1*10 - 0 = -10
	r = MakeRSXformFromRadians(1, rad, 0, 0, 10, 10)
	if !floatEquals(r.Tx, 10) || !floatEquals(r.Ty, -10) {
		t.Errorf("MakeFromRadians(anchor) Tx,Ty = %f,%f, want 10,-10", r.Tx, r.Ty)
	}
}

func TestRSXform_ToQuad(t *testing.T) {
	width := Scalar(100)
	height := Scalar(100)
	// Identity translated by 10, 20
	r := MakeRSXform(1, 0, 10, 20)

	quad := r.ToQuad(width, height)

	// Expected:
	// 0: (tx, ty) = (10, 20)
	// 1: (w + tx, ty) = (110, 20)
	// 2: (w + tx, h + ty) = (110, 120)
	// 3: (tx, h + ty) = (10, 120)

	expected := [4]Point{
		{10, 20},
		{110, 20},
		{110, 120},
		{10, 120},
	}

	for i, p := range quad {
		if !floatEquals(p.X, expected[i].X) || !floatEquals(p.Y, expected[i].Y) {
			t.Errorf("ToQuad[%d] = %v, want %v", i, p, expected[i])
		}
	}
}

func TestRSXform_ToTriStrip(t *testing.T) {
	width := Scalar(100)
	height := Scalar(100)
	// Identity translated by 10, 20
	r := MakeRSXform(1, 0, 10, 20)

	strip := r.ToTriStrip(width, height)

	// Expected TriStrip logic:
	// 0: (tx, ty) = (10, 20)
	// 1: (tx, h + ty) = (10, 120)  <-- diff from Quad
	// 2: (w + tx, ty) = (110, 20)
	// 3: (w + tx, h + ty) = (110, 120)

	expected := [4]Point{
		{10, 20},
		{10, 120},
		{110, 20},
		{110, 120},
	}

	for i, p := range strip {
		if !floatEquals(p.X, expected[i].X) || !floatEquals(p.Y, expected[i].Y) {
			t.Errorf("ToTriStrip[%d] = %v, want %v", i, p, expected[i])
		}
	}
}

func TestRSXform_RectStaysRect(t *testing.T) {
	r := MakeRSXform(1, 0, 0, 0) // Identity
	if !r.RectStaysRect() {
		t.Error("Identity should be RectStaysRect")
	}

	r = MakeRSXform(0, 1, 0, 0) // 90 deg rotation
	if !r.RectStaysRect() {
		t.Error("90 deg rotation should be RectStaysRect")
	}

	r = MakeRSXform(0.707, 0.707, 0, 0) // 45 deg
	if r.RectStaysRect() {
		t.Error("45 deg rotation should NOT be RectStaysRect")
	}
}
