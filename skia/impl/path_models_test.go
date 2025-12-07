package impl

import (
	"math"
	"testing"

	"github.com/zodimo/go-skia-support/skia/base"
	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/models"
)

// Helper function to create RRect with specific bounds and radii
func makeRRect(bounds models.Rect, radii [4]models.Point) models.RRect {
	rrect := models.RRect{}
	rrect.SetRectRadii(bounds, radii)
	return rrect
}

// TestRect_IsSorted tests the IsSorted method
func TestRect_IsSorted(t *testing.T) {
	tests := []struct {
		name     string
		rect     models.Rect
		expected bool
	}{
		{
			name: "sorted rect (normal)",
			rect: models.Rect{
				Left:   10,
				Top:    20,
				Right:  50,
				Bottom: 60,
			},
			expected: true,
		},
		{
			name: "sorted rect (equal edges)",
			rect: models.Rect{
				Left:   10,
				Top:    20,
				Right:  10, // Left == Right
				Bottom: 20, // Top == Bottom
			},
			expected: true,
		},
		{
			name: "unsorted rect (Left > Right)",
			rect: models.Rect{
				Left:   50,
				Top:    20,
				Right:  10,
				Bottom: 60,
			},
			expected: false,
		},
		{
			name: "unsorted rect (Top > Bottom)",
			rect: models.Rect{
				Left:   10,
				Top:    60,
				Right:  50,
				Bottom: 20,
			},
			expected: false,
		},
		{
			name: "unsorted rect (both reversed)",
			rect: models.Rect{
				Left:   50,
				Top:    60,
				Right:  10,
				Bottom: 20,
			},
			expected: false,
		},
		{
			name: "empty rect (zero width)",
			rect: models.Rect{
				Left:   10,
				Top:    20,
				Right:  10,
				Bottom: 60,
			},
			expected: true, // Still sorted (Left == Right)
		},
		{
			name: "empty rect (zero height)",
			rect: models.Rect{
				Left:   10,
				Top:    20,
				Right:  50,
				Bottom: 20,
			},
			expected: true, // Still sorted (Top == Bottom)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rect.IsSorted()
			if result != tt.expected {
				t.Errorf("IsSorted() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestRect_MakeOutset tests the MakeOutset method
func TestRect_MakeOutset(t *testing.T) {
	tests := []struct {
		name     string
		rect     models.Rect
		dx       base.Scalar
		dy       base.Scalar
		expected models.Rect
	}{
		{
			name: "positive outset (expand)",
			rect: models.Rect{
				Left:   10,
				Top:    20,
				Right:  50,
				Bottom: 60,
			},
			dx: 5,
			dy: 10,
			expected: models.Rect{
				Left:   5,  // 10 - 5
				Top:    10, // 20 - 10
				Right:  55, // 50 + 5
				Bottom: 70, // 60 + 10
			},
		},
		{
			name: "negative outset (contract)",
			rect: models.Rect{
				Left:   10,
				Top:    20,
				Right:  50,
				Bottom: 60,
			},
			dx: -5,
			dy: -10,
			expected: models.Rect{
				Left:   15, // 10 - (-5)
				Top:    30, // 20 - (-10)
				Right:  45, // 50 + (-5)
				Bottom: 50, // 60 + (-10)
			},
		},
		{
			name: "zero outset (no change)",
			rect: models.Rect{
				Left:   10,
				Top:    20,
				Right:  50,
				Bottom: 60,
			},
			dx: 0,
			dy: 0,
			expected: models.Rect{
				Left:   10,
				Top:    20,
				Right:  50,
				Bottom: 60,
			},
		},
		{
			name: "mixed positive/negative",
			rect: models.Rect{
				Left:   10,
				Top:    20,
				Right:  50,
				Bottom: 60,
			},
			dx: 5,
			dy: -10,
			expected: models.Rect{
				Left:   5,  // 10 - 5
				Top:    30, // 20 - (-10)
				Right:  55, // 50 + 5
				Bottom: 50, // 60 + (-10)
			},
		},
		{
			name: "very large outset",
			rect: models.Rect{
				Left:   10,
				Top:    20,
				Right:  50,
				Bottom: 60,
			},
			dx: 1000,
			dy: 2000,
			expected: models.Rect{
				Left:   -990,  // 10 - 1000
				Top:    -1980, // 20 - 2000
				Right:  1050,  // 50 + 1000
				Bottom: 2060,  // 60 + 2000
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rect.MakeOutset(tt.dx, tt.dy)

			if !NearlyEqualScalar(result.Left, tt.expected.Left) {
				t.Errorf("Left mismatch: got %f, expected %f", result.Left, tt.expected.Left)
			}
			if !NearlyEqualScalar(result.Top, tt.expected.Top) {
				t.Errorf("Top mismatch: got %f, expected %f", result.Top, tt.expected.Top)
			}
			if !NearlyEqualScalar(result.Right, tt.expected.Right) {
				t.Errorf("Right mismatch: got %f, expected %f", result.Right, tt.expected.Right)
			}
			if !NearlyEqualScalar(result.Bottom, tt.expected.Bottom) {
				t.Errorf("Bottom mismatch: got %f, expected %f", result.Bottom, tt.expected.Bottom)
			}
		})
	}
}

// TestRect_WithPath tests Rect operations when used with Path
func TestRect_WithPath(t *testing.T) {
	// Test AddRect with sorted rect
	sortedRect := models.Rect{
		Left:   10,
		Top:    20,
		Right:  50,
		Bottom: 60,
	}
	if !sortedRect.IsSorted() {
		t.Fatal("Test rect should be sorted")
	}

	path := NewSkPath(enums.PathFillTypeDefault)
	path.AddRect(sortedRect, enums.PathDirectionCW, 0)
	bounds := path.Bounds()

	// Verify bounds match rect
	if !NearlyEqualScalarDefault(bounds.Left, sortedRect.Left) {
		t.Errorf("Bounds.Left = %f, expected %f", bounds.Left, sortedRect.Left)
	}
	if !NearlyEqualScalarDefault(bounds.Top, sortedRect.Top) {
		t.Errorf("Bounds.Top = %f, expected %f", bounds.Top, sortedRect.Top)
	}
	if !NearlyEqualScalarDefault(bounds.Right, sortedRect.Right) {
		t.Errorf("Bounds.Right = %f, expected %f", bounds.Right, sortedRect.Right)
	}
	if !NearlyEqualScalarDefault(bounds.Bottom, sortedRect.Bottom) {
		t.Errorf("Bounds.Bottom = %f, expected %f", bounds.Bottom, sortedRect.Bottom)
	}

	// Test MakeOutset with Path bounds
	outsetRect := sortedRect.MakeOutset(5, 10)
	path2 := NewSkPath(enums.PathFillTypeDefault)
	path2.AddRect(outsetRect, enums.PathDirectionCW, 0)
	bounds2 := path2.Bounds()

	// Verify bounds match outset rect
	if !NearlyEqualScalarDefault(bounds2.Left, outsetRect.Left) {
		t.Errorf("Outset bounds.Left = %f, expected %f", bounds2.Left, outsetRect.Left)
	}
	if !NearlyEqualScalarDefault(bounds2.Top, outsetRect.Top) {
		t.Errorf("Outset bounds.Top = %f, expected %f", bounds2.Top, outsetRect.Top)
	}
}

// TestRRect_IsRect tests the IsRect method
func TestRRect_IsRect(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() models.RRect
		expected bool
	}{
		{
			name: "rect (all radii zero)",
			setup: func() models.RRect {
				rrect := models.RRect{}
				rrect.SetRect(models.Rect{Left: 10, Top: 20, Right: 50, Bottom: 60})
				return rrect
			},
			expected: true,
		},
		{
			name: "rect (X radius non-zero, Y zero - still rect)",
			setup: func() models.RRect {
				rrect := models.RRect{}
				rrect.SetRect(models.Rect{Left: 10, Top: 20, Right: 50, Bottom: 60})
				rrect.SetRectRadii(rrect.Rect(), [4]models.Point{
					{X: 5, Y: 0}, // X non-zero but Y zero - still rect
					{X: 0, Y: 0},
					{X: 0, Y: 0},
					{X: 0, Y: 0},
				})
				return rrect
			},
			expected: true, // IsRect returns true if at least one radius is zero per corner
		},
		{
			name: "rect (Y radius non-zero, X zero - still rect)",
			setup: func() models.RRect {
				rrect := models.RRect{}
				rrect.SetRect(models.Rect{Left: 10, Top: 20, Right: 50, Bottom: 60})
				rrect.SetRectRadii(rrect.Rect(), [4]models.Point{
					{X: 0, Y: 5}, // Y non-zero but X zero - still rect
					{X: 0, Y: 0},
					{X: 0, Y: 0},
					{X: 0, Y: 0},
				})
				return rrect
			},
			expected: true, // IsRect returns true if at least one radius is zero per corner
		},
		{
			name: "not rect (both radii non-zero)",
			setup: func() models.RRect {
				rrect := models.RRect{}
				rrect.SetRect(models.Rect{Left: 10, Top: 20, Right: 50, Bottom: 60})
				rrect.SetRectRadii(rrect.Rect(), [4]models.Point{
					{X: 5, Y: 5},
					{X: 0, Y: 0},
					{X: 0, Y: 0},
					{X: 0, Y: 0},
				})
				return rrect
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rrect := tt.setup()
			result := rrect.IsRect()
			if result != tt.expected {
				t.Errorf("IsRect() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestRRect_IsEmpty tests the IsEmpty method
func TestRRect_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() models.RRect
		expected bool
	}{
		{
			name: "non-empty rect",
			setup: func() models.RRect {
				return makeRRect(models.Rect{Left: 10, Top: 20, Right: 50, Bottom: 60},
					[4]models.Point{{X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}})
			},
			expected: false,
		},
		{
			name: "empty rect (zero width)",
			setup: func() models.RRect {
				return makeRRect(models.Rect{Left: 10, Top: 20, Right: 10, Bottom: 60},
					[4]models.Point{{X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}})
			},
			expected: true,
		},
		{
			name: "empty rect (zero height)",
			setup: func() models.RRect {
				return makeRRect(models.Rect{Left: 10, Top: 20, Right: 50, Bottom: 20},
					[4]models.Point{{X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}})
			},
			expected: true,
		},
		{
			name: "empty rect (reversed)",
			setup: func() models.RRect {
				return makeRRect(models.Rect{Left: 50, Top: 60, Right: 10, Bottom: 20},
					[4]models.Point{{X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}})
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rrect := tt.setup()
			result := rrect.IsEmpty()
			if result != tt.expected {
				t.Errorf("IsEmpty() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestRRect_IsOval tests the IsOval method
func TestRRect_IsOval(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() models.RRect
		expected bool
	}{
		{
			name: "oval (all radii equal, >= half dimensions)",
			setup: func() models.RRect {
				rrect := models.RRect{}
				rrect.SetOval(models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80})
				return rrect
			},
			expected: true,
		},
		{
			name: "not oval (radii too small)",
			setup: func() models.RRect {
				return makeRRect(models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80},
					[4]models.Point{{X: 40, Y: 30}, {X: 40, Y: 30}, {X: 40, Y: 30}, {X: 40, Y: 30}})
			},
			expected: false,
		},
		{
			name: "not oval (radii not equal)",
			setup: func() models.RRect {
				return makeRRect(models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80},
					[4]models.Point{{X: 50, Y: 40}, {X: 50, Y: 40}, {X: 50, Y: 40}, {X: 45, Y: 40}})
			},
			expected: false,
		},
		{
			name: "not oval (empty)",
			setup: func() models.RRect {
				return makeRRect(models.Rect{Left: 0, Top: 0, Right: 0, Bottom: 80},
					[4]models.Point{{X: 50, Y: 40}, {X: 50, Y: 40}, {X: 50, Y: 40}, {X: 50, Y: 40}})
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rrect := tt.setup()
			result := rrect.IsOval()
			if result != tt.expected {
				t.Errorf("IsOval() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestRRect_IsSimple tests the IsSimple method
func TestRRect_IsSimple(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() models.RRect
		expected bool
	}{
		{
			name: "simple (all radii equal, not oval)",
			setup: func() models.RRect {
				return makeRRect(models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80},
					[4]models.Point{{X: 10, Y: 10}, {X: 10, Y: 10}, {X: 10, Y: 10}, {X: 10, Y: 10}})
			},
			expected: true,
		},
		{
			name: "not simple (is rect)",
			setup: func() models.RRect {
				rrect := models.RRect{}
				rrect.SetRect(models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80})
				return rrect
			},
			expected: false,
		},
		{
			name: "not simple (is oval)",
			setup: func() models.RRect {
				rrect := models.RRect{}
				rrect.SetOval(models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80})
				return rrect
			},
			expected: false,
		},
		{
			name: "not simple (radii not equal)",
			setup: func() models.RRect {
				return makeRRect(models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80},
					[4]models.Point{{X: 10, Y: 10}, {X: 10, Y: 10}, {X: 10, Y: 10}, {X: 15, Y: 10}})
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rrect := tt.setup()
			result := rrect.IsSimple()
			if result != tt.expected {
				t.Errorf("IsSimple() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestRRect_IsNinePatch tests the IsNinePatch method
func TestRRect_IsNinePatch(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() models.RRect
		expected bool
	}{
		{
			name: "nine-patch (axis-aligned radii)",
			setup: func() models.RRect {
				rrect := models.RRect{}
				rrect.SetNinePatch(models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80}, 10, 5, 20, 5, 20, 15, 10, 15)
				return rrect
			},
			expected: true,
		},
		{
			name: "not nine-patch (is rect)",
			setup: func() models.RRect {
				rrect := models.RRect{}
				rrect.SetRect(models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80})
				return rrect
			},
			expected: false,
		},
		{
			name: "not nine-patch (is oval)",
			setup: func() models.RRect {
				rrect := models.RRect{}
				rrect.SetOval(models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80})
				return rrect
			},
			expected: false,
		},
		{
			name: "not nine-patch (is simple)",
			setup: func() models.RRect {
				rrect := models.RRect{}
				rrect.SetRectXY(models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80}, 10, 10)
				return rrect
			},
			expected: false,
		},
		{
			name: "not nine-patch (not axis-aligned)",
			setup: func() models.RRect {
				return makeRRect(models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80},
					[4]models.Point{{X: 10, Y: 5}, {X: 20, Y: 5}, {X: 20, Y: 15}, {X: 15, Y: 15}})
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rrect := tt.setup()
			result := rrect.IsNinePatch()
			if result != tt.expected {
				t.Errorf("IsNinePatch() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestRRect_Type tests the Type method
func TestRRect_Type(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() models.RRect
		expected enums.RRectType
	}{
		{
			name: "empty",
			setup: func() models.RRect {
				return makeRRect(models.Rect{Left: 0, Top: 0, Right: 0, Bottom: 80},
					[4]models.Point{{X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}})
			},
			expected: enums.RRectTypeEmpty,
		},
		{
			name: "rect",
			setup: func() models.RRect {
				rrect := models.RRect{}
				rrect.SetRect(models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80})
				return rrect
			},
			expected: enums.RRectTypeRect,
		},
		{
			name: "oval",
			setup: func() models.RRect {
				rrect := models.RRect{}
				rrect.SetOval(models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80})
				return rrect
			},
			expected: enums.RRectTypeOval,
		},
		{
			name: "simple",
			setup: func() models.RRect {
				rrect := models.RRect{}
				rrect.SetRectXY(models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80}, 10, 10)
				return rrect
			},
			expected: enums.RRectTypeSimple,
		},
		{
			name: "nine-patch",
			setup: func() models.RRect {
				rrect := models.RRect{}
				rrect.SetNinePatch(models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80}, 10, 5, 20, 5, 20, 15, 10, 15)
				return rrect
			},
			expected: enums.RRectTypeNinePatch,
		},
		{
			name: "complex",
			setup: func() models.RRect {
				return makeRRect(models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80},
					[4]models.Point{{X: 10, Y: 5}, {X: 20, Y: 5}, {X: 20, Y: 15}, {X: 15, Y: 15}})
			},
			expected: enums.RRectTypeComplex,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rrect := tt.setup()
			result := rrect.Type()
			if result != tt.expected {
				t.Errorf("Type() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestRRect_Accessors tests accessor methods
func TestRRect_Accessors(t *testing.T) {
	bounds := models.Rect{Left: 10, Top: 20, Right: 50, Bottom: 60}
	rrect := makeRRect(bounds, [4]models.Point{
		{X: 5, Y: 5},
		{X: 10, Y: 5},
		{X: 10, Y: 10},
		{X: 5, Y: 10},
	})

	// Test Rect()
	if rrect.Rect() != bounds {
		t.Errorf("Rect() = %v, expected %v", rrect.Rect(), bounds)
	}

	// Test Width()
	expectedWidth := base.Scalar(40) // 50 - 10
	if !NearlyEqualScalar(rrect.Width(), expectedWidth) {
		t.Errorf("Width() = %f, expected %f", rrect.Width(), expectedWidth)
	}

	// Test Height()
	expectedHeight := base.Scalar(40) // 60 - 20
	if !NearlyEqualScalar(rrect.Height(), expectedHeight) {
		t.Errorf("Height() = %f, expected %f", rrect.Height(), expectedHeight)
	}

	// Test RadiiAt()
	ulRadii := rrect.RadiiAt(enums.CornerUpperLeft)
	if !NearlyEqualScalar(ulRadii.X, 5) || !NearlyEqualScalar(ulRadii.Y, 5) {
		t.Errorf("RadiiAt(UpperLeft) = %v, expected {5, 5}", ulRadii)
	}

	urRadii := rrect.RadiiAt(enums.CornerUpperRight)
	if !NearlyEqualScalar(urRadii.X, 10) || !NearlyEqualScalar(urRadii.Y, 5) {
		t.Errorf("RadiiAt(UpperRight) = %v, expected {10, 5}", urRadii)
	}

	// Test GetAllRadii()
	allRadii := rrect.GetAllRadii()
	if len(allRadii) != 4 {
		t.Errorf("GetAllRadii() length = %d, expected 4", len(allRadii))
	}
	if !NearlyEqualScalar(allRadii[0].X, 5) || !NearlyEqualScalar(allRadii[0].Y, 5) {
		t.Errorf("GetAllRadii()[0] = %v, expected {5, 5}", allRadii[0])
	}
}

// TestRRect_Setters tests setter methods
func TestRRect_Setters(t *testing.T) {
	// Test SetRect
	rrect := models.RRect{}
	rect := models.Rect{Left: 10, Top: 20, Right: 50, Bottom: 60}
	rrect.SetRect(rect)

	if rrect.Rect() != rect {
		t.Errorf("SetRect failed: got %v, expected %v", rrect.Rect(), rect)
	}
	// Verify all radii are zero
	for i := 0; i < 4; i++ {
		if rrect.Radii[i].X != 0 || rrect.Radii[i].Y != 0 {
			t.Errorf("SetRect: Radii[%d] = %v, expected {0, 0}", i, rrect.Radii[i])
		}
	}

	// Test SetOval
	rrect2 := models.RRect{}
	ovalRect := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80}
	rrect2.SetOval(ovalRect)

	if rrect2.Rect() != ovalRect {
		t.Errorf("SetOval failed: got %v, expected %v", rrect2.Rect(), ovalRect)
	}
	// Verify all radii are half width/height
	expectedRx := base.Scalar(50) // 100 / 2
	expectedRy := base.Scalar(40) // 80 / 2
	for i := 0; i < 4; i++ {
		if !NearlyEqualScalar(rrect2.Radii[i].X, expectedRx) {
			t.Errorf("SetOval: Radii[%d].X = %f, expected %f", i, rrect2.Radii[i].X, expectedRx)
		}
		if !NearlyEqualScalar(rrect2.Radii[i].Y, expectedRy) {
			t.Errorf("SetOval: Radii[%d].Y = %f, expected %f", i, rrect2.Radii[i].Y, expectedRy)
		}
	}

	// Test SetRectXY
	rrect3 := models.RRect{}
	rect3 := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80}
	rrect3.SetRectXY(rect3, 10, 15)

	if rrect3.Rect() != rect3 {
		t.Errorf("SetRectXY failed: got %v, expected %v", rrect3.Rect(), rect3)
	}
	// Verify all radii are set to (10, 15)
	for i := 0; i < 4; i++ {
		if !NearlyEqualScalar(rrect3.Radii[i].X, 10) {
			t.Errorf("SetRectXY: Radii[%d].X = %f, expected 10", i, rrect3.Radii[i].X)
		}
		if !NearlyEqualScalar(rrect3.Radii[i].Y, 15) {
			t.Errorf("SetRectXY: Radii[%d].Y = %f, expected 15", i, rrect3.Radii[i].Y)
		}
	}

	// Test SetRectXY with clamping (radii too large)
	rrect4 := models.RRect{}
	rect4 := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80}
	rrect4.SetRectXY(rect4, 60, 50) // rx=60 > width/2=50, ry=50 > height/2=40

	// Should be clamped to width/2 and height/2
	expectedRxClamped := base.Scalar(50) // width / 2
	expectedRyClamped := base.Scalar(40) // height / 2
	for i := 0; i < 4; i++ {
		if !NearlyEqualScalar(rrect4.Radii[i].X, expectedRxClamped) {
			t.Errorf("SetRectXY (clamped): Radii[%d].X = %f, expected %f", i, rrect4.Radii[i].X, expectedRxClamped)
		}
		if !NearlyEqualScalar(rrect4.Radii[i].Y, expectedRyClamped) {
			t.Errorf("SetRectXY (clamped): Radii[%d].Y = %f, expected %f", i, rrect4.Radii[i].Y, expectedRyClamped)
		}
	}

	// Test SetNinePatch
	rrect5 := models.RRect{}
	rect5 := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80}
	rrect5.SetNinePatch(rect5, 10, 5, 20, 5, 20, 15, 10, 15)

	if rrect5.Rect() != rect5 {
		t.Errorf("SetNinePatch failed: got %v, expected %v", rrect5.Rect(), rect5)
	}
	// Verify radii
	if !NearlyEqualScalar(rrect5.Radii[0].X, 10) || !NearlyEqualScalar(rrect5.Radii[0].Y, 5) {
		t.Errorf("SetNinePatch: UL = %v, expected {10, 5}", rrect5.Radii[0])
	}
	if !NearlyEqualScalar(rrect5.Radii[1].X, 20) || !NearlyEqualScalar(rrect5.Radii[1].Y, 5) {
		t.Errorf("SetNinePatch: UR = %v, expected {20, 5}", rrect5.Radii[1])
	}
	if !NearlyEqualScalar(rrect5.Radii[2].X, 20) || !NearlyEqualScalar(rrect5.Radii[2].Y, 15) {
		t.Errorf("SetNinePatch: LR = %v, expected {20, 15}", rrect5.Radii[2])
	}
	if !NearlyEqualScalar(rrect5.Radii[3].X, 10) || !NearlyEqualScalar(rrect5.Radii[3].Y, 15) {
		t.Errorf("SetNinePatch: LL = %v, expected {10, 15}", rrect5.Radii[3])
	}
}

// TestRRect_WithPath tests RRect operations when used with Path
func TestRRect_WithPath(t *testing.T) {
	// Test AddRRect with rect (degenerate case)
	rect := models.Rect{Left: 10, Top: 20, Right: 50, Bottom: 60}
	rrect := models.RRect{}
	rrect.SetRect(rect)

	if !rrect.IsRect() {
		t.Fatal("Test RRect should be a rect")
	}

	path := NewSkPath(enums.PathFillTypeDefault)
	path.AddRRect(rrect, enums.PathDirectionCW)
	bounds := path.Bounds()

	// Verify bounds match rect
	if !NearlyEqualScalarDefault(bounds.Left, rect.Left) {
		t.Errorf("RRect bounds.Left = %f, expected %f", bounds.Left, rect.Left)
	}
	if !NearlyEqualScalarDefault(bounds.Top, rect.Top) {
		t.Errorf("RRect bounds.Top = %f, expected %f", bounds.Top, rect.Top)
	}

	// Test AddRRect with oval (degenerate case)
	ovalRect := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80}
	rrect2 := models.RRect{}
	rrect2.SetOval(ovalRect)

	if !rrect2.IsOval() {
		t.Fatal("Test RRect should be an oval")
	}

	path2 := NewSkPath(enums.PathFillTypeDefault)
	path2.AddRRect(rrect2, enums.PathDirectionCW)
	bounds2 := path2.Bounds()

	// Verify bounds match oval rect
	if !NearlyEqualScalarDefault(bounds2.Left, ovalRect.Left) {
		t.Errorf("Oval bounds.Left = %f, expected %f", bounds2.Left, ovalRect.Left)
	}
	if !NearlyEqualScalarDefault(bounds2.Top, ovalRect.Top) {
		t.Errorf("Oval bounds.Top = %f, expected %f", bounds2.Top, ovalRect.Top)
	}

	// Test AddRRect with simple RRect
	simpleRect := models.Rect{Left: 0, Top: 0, Right: 100, Bottom: 80}
	rrect3 := models.RRect{}
	rrect3.SetRectXY(simpleRect, 10, 10)

	if rrect3.Type() != enums.RRectTypeSimple {
		t.Fatalf("Test RRect type = %v, expected RRectTypeSimple", rrect3.Type())
	}

	path3 := NewSkPath(enums.PathFillTypeDefault)
	path3.AddRRect(rrect3, enums.PathDirectionCW)
	bounds3 := path3.Bounds()

	// Verify bounds match
	if !NearlyEqualScalarDefault(bounds3.Left, simpleRect.Left) {
		t.Errorf("Simple RRect bounds.Left = %f, expected %f", bounds3.Left, simpleRect.Left)
	}
}

// TestPoint_WithPath tests Point operations when used with Path
func TestPoint_WithPath(t *testing.T) {
	// Test Point with MoveTo
	point := models.Point{X: 10, Y: 20}
	path := NewSkPath(enums.PathFillTypeDefault)
	path.MoveTo(point.X, point.Y)

	// Verify point was added
	if path.CountPoints() != 1 {
		t.Errorf("CountPoints() = %d, expected 1", path.CountPoints())
	}

	// Test Point with LineTo
	point2 := models.Point{X: 50, Y: 60}
	path.LineTo(point2.X, point2.Y)

	if path.CountPoints() != 2 {
		t.Errorf("CountPoints() = %d, expected 2", path.CountPoints())
	}

	// Test Point transformation with matrix
	matrix := NewMatrixScale(2, 3)
	transformedPath := copyPath(path)
	transformedPath.Transform(matrix)
	bounds := transformedPath.Bounds()

	// Original bounds: (10, 20) to (50, 60)
	// Transformed: (20, 60) to (100, 180)
	expectedLeft := base.Scalar(20)    // 10 * 2
	expectedTop := base.Scalar(60)     // 20 * 3
	expectedRight := base.Scalar(100)  // 50 * 2
	expectedBottom := base.Scalar(180) // 60 * 3

	if !NearlyEqualScalarDefault(bounds.Left, expectedLeft) {
		t.Errorf("Transformed bounds.Left = %f, expected %f", bounds.Left, expectedLeft)
	}
	if !NearlyEqualScalarDefault(bounds.Top, expectedTop) {
		t.Errorf("Transformed bounds.Top = %f, expected %f", bounds.Top, expectedTop)
	}
	if !NearlyEqualScalarDefault(bounds.Right, expectedRight) {
		t.Errorf("Transformed bounds.Right = %f, expected %f", bounds.Right, expectedRight)
	}
	if !NearlyEqualScalarDefault(bounds.Bottom, expectedBottom) {
		t.Errorf("Transformed bounds.Bottom = %f, expected %f", bounds.Bottom, expectedBottom)
	}
}

// TestPoint_EdgeCases tests Point edge cases with Path
func TestPoint_EdgeCases(t *testing.T) {
	path := NewSkPath(enums.PathFillTypeDefault)

	// Test zero point
	zeroPoint := models.Point{X: 0, Y: 0}
	path.MoveTo(zeroPoint.X, zeroPoint.Y)
	if path.CountPoints() != 1 {
		t.Errorf("Zero point not added: CountPoints() = %d", path.CountPoints())
	}

	// Test negative point
	negPoint := models.Point{X: -10, Y: -20}
	path.LineTo(negPoint.X, negPoint.Y)
	if path.CountPoints() != 2 {
		t.Errorf("Negative point not added: CountPoints() = %d", path.CountPoints())
	}

	// Test very large point
	largePoint := models.Point{X: 1e10, Y: 1e10}
	path.LineTo(largePoint.X, largePoint.Y)
	if path.CountPoints() != 3 {
		t.Errorf("Large point not added: CountPoints() = %d", path.CountPoints())
	}

	// Test NaN point (should be handled gracefully)
	nanPoint := models.Point{X: base.Scalar(math.NaN()), Y: base.Scalar(math.NaN())}
	path.LineTo(nanPoint.X, nanPoint.Y)
	// Path should still function, though bounds may be invalid
	if path.CountPoints() != 4 {
		t.Errorf("NaN point not added: CountPoints() = %d", path.CountPoints())
	}
}
