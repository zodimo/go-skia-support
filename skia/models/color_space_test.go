package models

import (
	"testing"
)

func TestColorSpace_Singletons(t *testing.T) {
	srgb := MakeSRGB()
	if srgb == nil {
		t.Fatal("MakeSRGB returned nil")
	}
	if !srgb.IsSRGB() {
		t.Error("MakeSRGB().IsSRGB() should be true")
	}

	linear := MakeSRGBLinear()
	if linear == nil {
		t.Fatal("MakeSRGBLinear returned nil")
	}
	if linear.IsSRGB() {
		t.Error("MakeSRGBLinear().IsSRGB() should be false")
	}

	// Verify singleton behavior
	srgb2 := MakeSRGB()
	if srgb != srgb2 {
		t.Error("MakeSRGB should return singleton")
	}
}

func TestColorSpace_MakeRGB(t *testing.T) {
	// Test creating sRGB manually returns singleton
	srgbFromRGB := MakeRGB(NamedTransferFnSRGB, NamedGamutSRGB)
	if srgbFromRGB != MakeSRGB() {
		t.Error("MakeRGB with sRGB params should return sRGB singleton")
	}

	// Test creating Linear sRGB manually returns singleton
	linearFromRGB := MakeRGB(NamedTransferFnLinear, NamedGamutSRGB)
	if linearFromRGB != MakeSRGBLinear() {
		t.Error("MakeRGB with Linear sRGB params should return Linear singleton")
	}

	// Test custom color space creation
	customMatrix := NamedGamutSRGB
	customMatrix.Vals[0][0] = 0.5 // modify slightly

	custom := MakeRGB(NamedTransferFnSRGB, customMatrix)
	if custom == nil {
		t.Fatal("MakeRGB returned nil for custom")
	}
	if custom == MakeSRGB() {
		t.Error("Custom color space should not match sRGB singleton")
	}

	// Verify fields
	if custom.TransferFn() != NamedTransferFnSRGB {
		t.Error("Transfer function mismatch")
	}
	if custom.ToXYZD50() != customMatrix {
		t.Error("Matrix mismatch")
	}
}

func TestColorSpace_Properties(t *testing.T) {
	srgb := MakeSRGB()
	if !srgb.GammaCloseToSRGB() {
		t.Error("sRGB should have gamma close to sRGB")
	}
	if srgb.GammaIsLinear() {
		t.Error("sRGB should not be linear")
	}

	linear := MakeSRGBLinear()
	if linear.GammaCloseToSRGB() {
		t.Error("Linear sRGB should not have gamma close to sRGB (it is linear)")
	}
	if !linear.GammaIsLinear() {
		t.Error("Linear sRGB should be linear")
	}
}

func TestColorSpace_LazyFields(t *testing.T) {
	customMatrix := NamedGamutSRGB
	// Make sure it is invertible. sRGB matrix is invertible.

	cs := MakeRGB(NamedTransferFnSRGB, customMatrix)
	if cs.fromXYZD50.Vals[0][0] != 0 {
		// This assumes 0 initialization means not calculated.
		// However, 0 is a valid value.
		// Better check logic: lazyOnce should not have run.
	}

	// Trigger lazy comp via internal method or public if exposed?
	// currently computeLazyDstFields is private/unexported
	// But it is called if we needed inverse transfer function or inverse matrix.
	// Current API doesn't expose methods that *require* lazy fields publically except maybe if we added ToXYZD50Hash or similar.

	// Wait, I didn't expose methods that use lazy fields in my Go port yet (like invTransferFn).
	// But I implemented computeLazyDstFields.
	// Let's verify it via reflection or just trust it compiles.
	// Actually, I can add a test that ensures it doesn't panic.

	// We can't easily test private fields from test unless it's in same package. This is models_test in models package.
	cs.computeLazyDstFields()

	// Check if fromXYZD50 is populated (not zero matrix)
	// sRGB inverse is roughly:
	// 3.24   -1.53  -0.49
	// -0.96   1.87   0.04
	// 0.05   -0.20   1.05

	inv := cs.fromXYZD50
	if inv.Vals[0][0] == 0 && inv.Vals[1][1] == 0 {
		t.Error("Lazy computation of fromXYZD50 failed (result looks like zero matrix)")
	}
}
