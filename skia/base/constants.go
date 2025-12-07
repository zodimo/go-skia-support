package base

// Matrix indices matching C++ SkMatrix constants
const (
	KMScaleX = 0 // horizontal scale factor
	KMSkewX  = 1 // horizontal skew factor
	KMTransX = 2 // horizontal translation
	KMSkewY  = 3 // vertical skew factor
	KMScaleY = 4 // vertical scale factor
	KMTransY = 5 // vertical translation
	KMPersp0 = 6 // input x perspective factor
	KMPersp1 = 7 // input y perspective factor
	KMPersp2 = 8 // perspective bias
)

const (
	SkScalarNearlyZero = 1.0 / (1 << 12)
)

// SegmentMask constants
const (
	SegmentMaskLine  uint32 = 1 << 0
	SegmentMaskQuad  uint32 = 1 << 1
	SegmentMaskConic uint32 = 1 << 2
	SegmentMaskCubic uint32 = 1 << 3
)

// SK_ScalarRoot2Over2 is sqrt(2)/2, the weight used for quarter-circle conics
const ScalarRoot2Over2 Scalar = 0.707106781
