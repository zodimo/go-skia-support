package impl

import (
	"math"
	"math/rand"
	"testing"

	"github.com/zodimo/go-skia-support/skia/base"
	"github.com/zodimo/go-skia-support/skia/models"
)

// TestMatrixInversion tests matrix inversion with various matrix types.
// Ported from: skia-source/tests/MatrixTest.cpp:DEF_TEST(Matrix, reporter)
func TestMatrixInversion(t *testing.T) {
	// Test 1: Translate matrix inversion
	mat := NewMatrixIdentity()
	mat.SetTranslate(1, 1)
	inverse, ok := mat.Invert()
	if !ok {
		t.Fatal("Expected translate matrix to be invertible")
	}
	iden1 := NewMatrixIdentity()
	iden1.SetConcat(mat, inverse)
	if !IsIdentity(iden1) {
		t.Error("M * M^-1 should equal identity for translate matrix")
	}

	// Test 2: Scale matrix inversion
	mat.SetScale(2, 4)
	inverse, ok = mat.Invert()
	if !ok {
		t.Fatal("Expected scale matrix to be invertible")
	}
	iden1 = NewMatrixIdentity()
	iden1.SetConcat(mat, inverse)
	if !IsIdentity(iden1) {
		t.Error("M * M^-1 should equal identity for scale matrix")
	}

	// Test 3: Scale matrix with fractional scale inversion
	mat.SetScale(0.5, 2) // SK_Scalar1/2 = 0.5
	inverse, ok = mat.Invert()
	if !ok {
		t.Fatal("Expected fractional scale matrix to be invertible")
	}
	iden1 = NewMatrixIdentity()
	iden1.SetConcat(mat, inverse)
	if !IsIdentity(iden1) {
		t.Error("M * M^-1 should equal identity for fractional scale matrix")
	}

	// Test 4: Scale+rotate matrix inversion
	// C++: mat.setScale(3, 5, 20, 0).postRotate(25)
	// In Go, we need to set scale, translate to pivot, rotate, then translate back
	mat.SetScale(3, 5)
	mat.PostTranslate(20, 0) // Translate to pivot
	mat.PostRotate(25, 0, 0) // Rotate about origin (which is now at pivot)
	mat.PostTranslate(-20, 0) // Translate back
	inverse, ok = mat.Invert()
	if !ok {
		t.Fatal("Expected scale+rotate matrix to be invertible")
	}
	iden1 = NewMatrixIdentity()
	iden1.SetConcat(mat, inverse)
	if !IsIdentity(iden1) {
		t.Error("M * M^-1 should equal identity for scale+rotate matrix")
	}
	iden2 := NewMatrixIdentity()
	iden2.SetConcat(inverse, mat)
	if !IsIdentity(iden2) {
		t.Error("M^-1 * M should equal identity for scale+rotate matrix")
	}

	// Test 5: Zero scale X (non-invertible)
	mat.SetScale(0, 1)
	_, ok = mat.Invert()
	if ok {
		t.Error("Expected zero X scale matrix to be non-invertible")
	}

	// Test 6: Zero scale Y (non-invertible)
	mat.SetScale(1, 0)
	_, ok = mat.Invert()
	if ok {
		t.Error("Expected zero Y scale matrix to be non-invertible")
	}

	// Test 7: Matrix that results in non-finite inverse
	// C++: mat.setAll(0.0f, 1.0f, 2.0f, 0.0f, 1.0f, -3.40277175e+38f, 1.00003040f, 1.0f, 0.0f)
	mat.SetAll(0.0, 1.0, 2.0, 0.0, 1.0, -3.40277175e+38, 1.00003040, 1.0, 0.0)
	_, ok = mat.Invert()
	if ok {
		t.Error("Expected matrix with non-finite inverse to be non-invertible")
	}

	// Test 8: Denormalized scale matrix (non-invertible due to overflow)
	// C++: mat.setAll(std::numeric_limits<float>::denorm_min(), 0.f, 0.f, 0.f, 1.f, 0.f, 0.f, 0.f, 1.f)
	denormMin := base.Scalar(math.SmallestNonzeroFloat32)
	mat.SetAll(denormMin, 0, 0, 0, 1, 0, 0, 0, 1)
	if !mat.IsScaleTranslate() {
		t.Error("Expected denorm scale matrix to be scale+translate")
	}
	_, ok = mat.Invert()
	if ok {
		t.Error("Expected denorm scale matrix to be non-invertible (1/scale overflows)")
	}

	// Test 9: Scale+translate matrix with NaN translation (non-invertible)
	// C++: mat.setAll(2.f, 0.f, std::numeric_limits<float>::quiet_NaN(), 0.f, 2.f, 0.f, 0.f, 0.f, 1.f)
	nanVal := base.Scalar(math.NaN())
	mat.SetAll(2, 0, nanVal, 0, 2, 0, 0, 0, 1)
	if !mat.IsScaleTranslate() {
		t.Error("Expected NaN translate matrix to be scale+translate")
	}
	_, ok = mat.Invert()
	if ok {
		t.Error("Expected scale+translate matrix with NaN translation to be non-invertible")
	}

	// Test 10: Translate-only matrix with NaN translation (non-invertible)
	// C++: mat.setAll(1.f, 0.f, std::numeric_limits<float>::quiet_NaN(), 0.f, 1.f, 0.f, 0.f, 0.f, 1.f)
	mat.SetAll(1, 0, nanVal, 0, 1, 0, 0, 0, 1)
	if !mat.IsTranslate() {
		t.Error("Expected NaN translate-only matrix to be translate")
	}
	_, ok = mat.Invert()
	if ok {
		t.Error("Expected translate-only matrix with NaN translation to be non-invertible")
	}

	// Test 11: Scale+translate matrix that becomes non-finite after inversion
	// C++: mat.setAll(std::numeric_limits<float>::min(), 0.f, std::numeric_limits<float>::max(), 0.f, 1.f, 0.f, 0.f, 0.f, 1.f)
	minFloat := base.Scalar(math.SmallestNonzeroFloat32)
	maxFloat := base.Scalar(math.MaxFloat32)
	mat.SetAll(minFloat, 0, maxFloat, 0, 1, 0, 0, 0, 1)
	if !mat.IsScaleTranslate() {
		t.Error("Expected extreme value matrix to be scale+translate")
	}
	if !IsFiniteMatrix(mat) {
		t.Error("Expected extreme value matrix to be finite")
	}
	_, ok = mat.Invert()
	if ok {
		t.Error("Expected extreme value matrix to be non-invertible (trans/scale becomes non-finite)")
	}
}

// TestMatrixConcatenation tests matrix concatenation operations.
// Ported from: skia-source/tests/MatrixTest.cpp:DEF_TEST(Matrix_Concat, r)
func TestMatrixConcatenation(t *testing.T) {
	// Test basic concatenation: Translate * Scale
	a := NewMatrixIdentity()
	a.SetTranslate(10, 20)

	b := NewMatrixIdentity()
	b.SetScale(3, 5)

	expected := NewMatrixIdentity()
	expected.SetConcat(a, b)

	// Verify SetConcat produces correct result
	// Expected: Scale first, then translate
	// Result should be: scale(3,5) then translate(10,20)
	// Which gives: scaleX=3, scaleY=5, transX=10, transY=20
	if !NearlyEqualScalar(expected.GetScaleX(), 3) {
		t.Errorf("Expected scaleX=3, got %f", expected.GetScaleX())
	}
	if !NearlyEqualScalar(expected.GetScaleY(), 5) {
		t.Errorf("Expected scaleY=5, got %f", expected.GetScaleY())
	}
	if !NearlyEqualScalar(expected.GetTranslateX(), 10) {
		t.Errorf("Expected transX=10, got %f", expected.GetTranslateX())
	}
	if !NearlyEqualScalar(expected.GetTranslateY(), 20) {
		t.Errorf("Expected transY=20, got %f", expected.GetTranslateY())
	}

	// Test associativity: (A * B) * C == A * (B * C)
	a = NewMatrixIdentity()
	a.SetTranslate(5, 10)
	b = NewMatrixIdentity()
	b.SetScale(2, 3)
	c := NewMatrixIdentity()
	c.SetRotate(45, 0, 0)

	ab := NewMatrixIdentity()
	ab.SetConcat(a, b)
	abc1 := NewMatrixIdentity()
	abc1.SetConcat(ab, c)

	bc := NewMatrixIdentity()
	bc.SetConcat(b, c)
	abc2 := NewMatrixIdentity()
	abc2.SetConcat(a, bc)

	if !NearlyEqual(abc1, abc2) {
		t.Error("Matrix concatenation should be associative: (A * B) * C == A * (B * C)")
	}

	// Test identity matrix concatenation: M * I == I * M == M
	m := NewMatrixIdentity()
	m.SetScale(2, 3)
	m.PostTranslate(5, 7)

	identity := NewMatrixIdentity()

	mi := NewMatrixIdentity()
	mi.SetConcat(m, identity)
	if !NearlyEqual(mi, m) {
		t.Error("M * I should equal M")
	}

	im := NewMatrixIdentity()
	im.SetConcat(identity, m)
	if !NearlyEqual(im, m) {
		t.Error("I * M should equal M")
	}
}

// TestMatrixMapRect tests matrix rect transformation.
// Ported from: skia-source/tests/MatrixTest.cpp:DEF_TEST(Matrix_maprects, r)
func TestMatrixMapRect(t *testing.T) {
	const scale = 1000.0

	mat := NewMatrixIdentity()
	mat.SetScale(2, 3)
	mat.PostTranslate(1, 4)

	// Test that mapRect matches mapPoints on rect corners
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 1000; i++ { // Reduced from 10000 for faster tests
		src := models.Rect{
			Left:   base.Scalar(rng.Float32()*scale - scale/2),
			Top:    base.Scalar(rng.Float32()*scale - scale/2),
			Right:  base.Scalar(rng.Float32()*scale - scale/2),
			Bottom: base.Scalar(rng.Float32()*scale - scale/2),
		}

		// Ensure rect is sorted
		if src.Left > src.Right {
			src.Left, src.Right = src.Right, src.Left
		}
		if src.Top > src.Bottom {
			src.Top, src.Bottom = src.Bottom, src.Top
		}

		// Map rect corners using MapPoints
		corners := []models.Point{
			{X: src.Left, Y: src.Top},
			{X: src.Right, Y: src.Top},
			{X: src.Right, Y: src.Bottom},
			{X: src.Left, Y: src.Bottom},
		}
		mappedCorners := make([]models.Point, 4)
		mat.MapPoints(mappedCorners, corners)

		// Find bounding box of mapped corners
		minX := mappedCorners[0].X
		maxX := mappedCorners[0].X
		minY := mappedCorners[0].Y
		maxY := mappedCorners[0].Y
		for _, pt := range mappedCorners[1:] {
			if pt.X < minX {
				minX = pt.X
			}
			if pt.X > maxX {
				maxX = pt.X
			}
			if pt.Y < minY {
				minY = pt.Y
			}
			if pt.Y > maxY {
				maxY = pt.Y
			}
		}
		expectedRect := models.Rect{Left: minX, Top: minY, Right: maxX, Bottom: maxY}

		// Map rect using MapRect
		mappedRect := mat.MapRect(src)

		// Compare results (within tolerance)
		if !NearlyEqualScalar(mappedRect.Left, expectedRect.Left) ||
			!NearlyEqualScalar(mappedRect.Right, expectedRect.Right) ||
			!NearlyEqualScalar(mappedRect.Top, expectedRect.Top) ||
			!NearlyEqualScalar(mappedRect.Bottom, expectedRect.Bottom) {
			t.Errorf("MapRect result doesn't match MapPoints on corners:\n  MapRect: %v\n  Expected: %v",
				mappedRect, expectedRect)
		}
	}

	// Test non-finite rect handling after mapping with large scale
	{
		// Scale matrix with very large scale factors
		m0 := NewMatrixScale(1e20, 1e20)
		rect := models.Rect{Left: 0, Top: 0, Right: 1e20, Bottom: 1e20}
		mapped := m0.MapRect(rect)
		// Result should be non-finite (infinity)
		if IsFinite(mapped.Left) && IsFinite(mapped.Right) && IsFinite(mapped.Top) && IsFinite(mapped.Bottom) {
			t.Error("Expected mapped rect with large scale to be non-finite")
		}
	}

	// Test identity matrix mapping
	identity := NewMatrixIdentity()
	testRect := models.Rect{Left: 10, Top: 20, Right: 30, Bottom: 40}
	mapped := identity.MapRect(testRect)
	if !NearlyEqualScalar(mapped.Left, testRect.Left) ||
		!NearlyEqualScalar(mapped.Right, testRect.Right) ||
		!NearlyEqualScalar(mapped.Top, testRect.Top) ||
		!NearlyEqualScalar(mapped.Bottom, testRect.Bottom) {
		t.Error("Identity matrix should not change rect")
	}

	// Test translate-only matrix
	translate := NewMatrixTranslate(5, 10)
	mapped = translate.MapRect(testRect)
	if !NearlyEqualScalar(mapped.Left, testRect.Left+5) ||
		!NearlyEqualScalar(mapped.Right, testRect.Right+5) ||
		!NearlyEqualScalar(mapped.Top, testRect.Top+10) ||
		!NearlyEqualScalar(mapped.Bottom, testRect.Bottom+10) {
		t.Error("Translate matrix should only translate rect")
	}

	// Test scale-only matrix
	scaleMat := NewMatrixScale(2, 3)
	mapped = scaleMat.MapRect(testRect)
	if !NearlyEqualScalar(mapped.Left, testRect.Left*2) ||
		!NearlyEqualScalar(mapped.Right, testRect.Right*2) ||
		!NearlyEqualScalar(mapped.Top, testRect.Top*3) ||
		!NearlyEqualScalar(mapped.Bottom, testRect.Bottom*3) {
		t.Error("Scale matrix should only scale rect")
	}
}

// TestMatrixGetterSetter tests matrix getter and setter methods.
// Ported from: skia-source/tests/MatrixTest.cpp:test_set9()
func TestMatrixGetterSetter(t *testing.T) {
	// Test 1: Identity matrix via Set9
	m := NewMatrixIdentity()
	var buffer [9]base.Scalar
	buffer[base.KMScaleX] = 1
	buffer[base.KMScaleY] = 1
	buffer[base.KMPersp2] = 1
	m.Set9(buffer)
	if !m.IsIdentity() {
		t.Error("Set9 with identity values should create identity matrix")
	}

	// Verify Get9 returns correct values
	values := m.Get9()
	expected := [9]base.Scalar{1, 0, 0, 0, 1, 0, 0, 0, 1}
	for i := 0; i < 9; i++ {
		if !NearlyEqualScalar(values[i], expected[i]) {
			t.Errorf("Get9[%d] = %f, expected %f", i, values[i], expected[i])
		}
	}

	// Verify RC() accessor matches Get9()
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			idx := row*3 + col
			rcVal := m.GetRC(row, col)
			if !NearlyEqualScalar(rcVal, values[idx]) {
				t.Errorf("RC(%d,%d) = %f, Get9[%d] = %f", row, col, rcVal, idx, values[idx])
			}
		}
	}

	// Test 2: Scale matrix via Set9
	m.SetScale(2, 3)
	values = m.Get9()
	expected = [9]base.Scalar{2, 0, 0, 0, 3, 0, 0, 0, 1}
	for i := 0; i < 9; i++ {
		if !NearlyEqualScalar(values[i], expected[i]) {
			t.Errorf("Get9[%d] = %f, expected %f for scale(2,3)", i, values[i], expected[i])
		}
	}

	// Test 3: Set9 after post-translate
	m.SetScale(2, 3)
	m.PostTranslate(4, 5)
	values = m.Get9()
	expected = [9]base.Scalar{2, 0, 4, 0, 3, 5, 0, 0, 1}
	for i := 0; i < 9; i++ {
		if !NearlyEqualScalar(values[i], expected[i]) {
			t.Errorf("Get9[%d] = %f, expected %f for scale(2,3).postTranslate(4,5)", i, values[i], expected[i])
		}
	}

	// Test 4: Verify Set9 round-trip
	testValues := [9]base.Scalar{1.5, 2.5, 3.5, 4.5, 5.5, 6.5, 7.5, 8.5, 9.5}
	m.Set9(testValues)
	roundTrip := m.Get9()
	for i := 0; i < 9; i++ {
		if !NearlyEqualScalar(roundTrip[i], testValues[i]) {
			t.Errorf("Set9/Get9 round-trip failed at index %d: got %f, expected %f", i, roundTrip[i], testValues[i])
		}
	}

	// Test 5: Verify RC() matches Get9() for all positions
	for i := 0; i < 9; i++ {
		row := i / 3
		col := i % 3
		rcVal := m.GetRC(row, col)
		if !NearlyEqualScalar(rcVal, testValues[i]) {
			t.Errorf("RC(%d,%d) = %f, Get9[%d] = %f", row, col, rcVal, i, testValues[i])
		}
	}
}

// TestMatrixHelperFunctions tests matrix helper functions independently.
// AC: 3 - Test coverage for matrix_helpers.go functions
func TestMatrixHelperFunctions(t *testing.T) {
	// Test rowcol3() helper function
	// rowcol3 computes dot product of row vector with column vector
	// row = [r0, r1, r2], col = [c0, c3, c6] (striding by 3)
	t.Run("rowcol3", func(t *testing.T) {
		row := []base.Scalar{1, 2, 3}
		col := []base.Scalar{4, 0, 0, 5, 0, 0, 6, 0, 0} // Stride by 3: [4, 5, 6]
		result := rowcol3(row, col)
		expected := base.Scalar(1*4 + 2*5 + 3*6) // = 4 + 10 + 18 = 32
		if !NearlyEqualScalar(result, expected) {
			t.Errorf("rowcol3() = %f, expected %f", result, expected)
		}

		// Test with zero values
		row = []base.Scalar{0, 0, 0}
		result = rowcol3(row, col)
		if result != 0 {
			t.Errorf("rowcol3() with zero row should return 0, got %f", result)
		}

		// Test with identity-like values
		row = []base.Scalar{1, 0, 0}
		col = []base.Scalar{1, 0, 0, 0, 1, 0, 0, 0, 1}
		result = rowcol3(row, col)
		if !NearlyEqualScalar(result, 1) {
			t.Errorf("rowcol3() with identity row/col should return 1, got %f", result)
		}
	})

	// Test muladdmul() helper function
	// muladdmul computes a*b + c*d
	t.Run("muladdmul", func(t *testing.T) {
		result := muladdmul(2, 3, 4, 5)
		expected := base.Scalar(2*3 + 4*5) // = 6 + 20 = 26
		if !NearlyEqualScalar(result, expected) {
			t.Errorf("muladdmul(2,3,4,5) = %f, expected %f", result, expected)
		}

		// Test with zero values
		result = muladdmul(0, 5, 3, 0)
		if result != 0 {
			t.Errorf("muladdmul(0,5,3,0) should return 0, got %f", result)
		}

		// Test with negative values
		result = muladdmul(-2, 3, 4, -5)
		expected = base.Scalar(-2*3 + 4*(-5)) // = -6 + -20 = -26
		if !NearlyEqualScalar(result, expected) {
			t.Errorf("muladdmul(-2,3,4,-5) = %f, expected %f", result, expected)
		}

		// Test with near-zero values
		small := base.Scalar(1e-10)
		result = muladdmul(small, 1, small, 1)
		if !NearlyEqualScalar(result, 2*small) {
			t.Errorf("muladdmul with small values failed")
		}
	})

	// Test scross_dscale() helper function
	// scross_dscale computes (a*b - c*d) * scale
	t.Run("scross_dscale", func(t *testing.T) {
		result := scross_dscale(2, 3, 4, 5, 0.5)
		expected := float64(2*3-4*5) * 0.5 // = (6 - 20) * 0.5 = -14 * 0.5 = -7
		if !NearlyEqualScalar(base.Scalar(result), base.Scalar(expected)) {
			t.Errorf("scross_dscale(2,3,4,5,0.5) = %f, expected %f", result, expected)
		}

		// Test with zero scale
		result = scross_dscale(2, 3, 4, 5, 0)
		if result != 0 {
			t.Errorf("scross_dscale with zero scale should return 0, got %f", result)
		}

		// Test with zero cross product
		result = scross_dscale(2, 3, 2, 3, 1.0)
		if result != 0 {
			t.Errorf("scross_dscale with equal vectors should return 0, got %f", result)
		}
	})
}

// TestMatrixModelOperations tests matrix-related model operations (Point, Rect transformations).
// AC: 4 - Test coverage for matrix-related model operations
func TestMatrixModelOperations(t *testing.T) {
	// Test Point transformation via Matrix
	t.Run("MapPoint", func(t *testing.T) {
		// Identity matrix should not change point
		identity := NewMatrixIdentity()
		pt := models.Point{X: 10, Y: 20}
		mapped := identity.MapPoint(pt)
		if !NearlyEqualScalar(mapped.X, pt.X) || !NearlyEqualScalar(mapped.Y, pt.Y) {
			t.Errorf("Identity matrix should not change point: got (%f,%f), expected (%f,%f)",
				mapped.X, mapped.Y, pt.X, pt.Y)
		}

		// Translate matrix
		translate := NewMatrixTranslate(5, 10)
		mapped = translate.MapPoint(pt)
		if !NearlyEqualScalar(mapped.X, pt.X+5) || !NearlyEqualScalar(mapped.Y, pt.Y+10) {
			t.Errorf("Translate matrix failed: got (%f,%f), expected (%f,%f)",
				mapped.X, mapped.Y, pt.X+5, pt.Y+10)
		}

		// Scale matrix
		scale := NewMatrixScale(2, 3)
		mapped = scale.MapPoint(pt)
		if !NearlyEqualScalar(mapped.X, pt.X*2) || !NearlyEqualScalar(mapped.Y, pt.Y*3) {
			t.Errorf("Scale matrix failed: got (%f,%f), expected (%f,%f)",
				mapped.X, mapped.Y, pt.X*2, pt.Y*3)
		}

		// Rotate matrix (90 degrees should map (1,0) to (0,1))
		rotate := NewMatrixRotate(90)
		pt = models.Point{X: 1, Y: 0}
		mapped = rotate.MapPoint(pt)
		if !NearlyEqualScalar(mapped.X, 0) || !NearlyEqualScalar(mapped.Y, 1) {
			t.Errorf("Rotate 90° failed: got (%f,%f), expected (0,1)", mapped.X, mapped.Y)
		}

		// Zero point
		pt = models.Point{X: 0, Y: 0}
		mapped = scale.MapPoint(pt)
		if mapped.X != 0 || mapped.Y != 0 {
			t.Errorf("Zero point should map to zero: got (%f,%f)", mapped.X, mapped.Y)
		}
	})

	// Test MapPoints (batch point transformation)
	t.Run("MapPoints", func(t *testing.T) {
		mat := NewMatrixScale(2, 3)
		src := []models.Point{
			{X: 1, Y: 2},
			{X: 3, Y: 4},
			{X: 5, Y: 6},
		}
		dst := make([]models.Point, len(src))
		count := mat.MapPoints(dst, src)
		if count != len(src) {
			t.Errorf("MapPoints returned count %d, expected %d", count, len(src))
		}
		for i := range src {
			expected := mat.MapPoint(src[i])
			if !NearlyEqualScalar(dst[i].X, expected.X) || !NearlyEqualScalar(dst[i].Y, expected.Y) {
				t.Errorf("MapPoints[%d] = (%f,%f), expected (%f,%f)",
					i, dst[i].X, dst[i].Y, expected.X, expected.Y)
			}
		}

		// Test with empty slices
		count = mat.MapPoints([]models.Point{}, []models.Point{})
		if count != 0 {
			t.Errorf("MapPoints with empty slices should return 0, got %d", count)
		}

		// Test with mismatched sizes
		dst = make([]models.Point, 2)
		count = mat.MapPoints(dst, src)
		if count != 2 {
			t.Errorf("MapPoints with mismatched sizes should return min(len(dst),len(src)), got %d", count)
		}
	})

	// Test Rect transformation via Matrix
	t.Run("MapRect", func(t *testing.T) {
		// Identity matrix should not change rect
		identity := NewMatrixIdentity()
		rect := models.Rect{Left: 10, Top: 20, Right: 30, Bottom: 40}
		mapped := identity.MapRect(rect)
		if !NearlyEqualScalar(mapped.Left, rect.Left) ||
			!NearlyEqualScalar(mapped.Right, rect.Right) ||
			!NearlyEqualScalar(mapped.Top, rect.Top) ||
			!NearlyEqualScalar(mapped.Bottom, rect.Bottom) {
			t.Error("Identity matrix should not change rect")
		}

		// Translate matrix
		translate := NewMatrixTranslate(5, 10)
		mapped = translate.MapRect(rect)
		if !NearlyEqualScalar(mapped.Left, rect.Left+5) ||
			!NearlyEqualScalar(mapped.Right, rect.Right+5) ||
			!NearlyEqualScalar(mapped.Top, rect.Top+10) ||
			!NearlyEqualScalar(mapped.Bottom, rect.Bottom+10) {
			t.Error("Translate matrix should only translate rect")
		}

		// Scale matrix
		scale := NewMatrixScale(2, 3)
		mapped = scale.MapRect(rect)
		if !NearlyEqualScalar(mapped.Left, rect.Left*2) ||
			!NearlyEqualScalar(mapped.Right, rect.Right*2) ||
			!NearlyEqualScalar(mapped.Top, rect.Top*3) ||
			!NearlyEqualScalar(mapped.Bottom, rect.Bottom*3) {
			t.Error("Scale matrix should only scale rect")
		}

		// Empty rect (zero area)
		emptyRect := models.Rect{Left: 10, Top: 20, Right: 10, Bottom: 20}
		mapped = scale.MapRect(emptyRect)
		if mapped.Left != mapped.Right || mapped.Top != mapped.Bottom {
			t.Error("Empty rect should map to empty rect")
		}

		// Negative rect (unsorted)
		negRect := models.Rect{Left: 30, Top: 40, Right: 10, Bottom: 20}
		mapped = scale.MapRect(negRect)
		// MapRect should handle unsorted rects correctly
		if mapped.Left > mapped.Right || mapped.Top > mapped.Bottom {
			t.Error("Mapped rect should be sorted")
		}
	})
}

