package interfaces

import (
	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/models"
)

type SkPaint interface {
	AsBlendMode() (enums.BlendMode, bool)
	CanComputeFastBounds() bool
	ComputeFastBounds(orig Rect, storage *Rect) Rect
	ComputeFastStrokeBounds(orig Rect, storage *Rect) Rect
	DoComputeFastBounds(origSrc Rect, storage *Rect, style enums.PaintStyle) Rect
	Equals(other SkPaint) bool
	GetAlpha() uint8
	GetAlphaf() Scalar
	GetBlendModeOr(defaultMode enums.BlendMode) enums.BlendMode
	GetBlender() Blender
	GetColor() models.Color4f
	GetColorFilter() ColorFilter
	GetColorInt() uint32
	GetImageFilter() ImageFilter
	GetInflationRadius(style enums.PaintStyle, matrixScale ...Scalar) Scalar
	GetMaskFilter() MaskFilter
	GetPathEffect() PathEffect
	GetShader() Shader
	GetStrokeCap() enums.PaintCap
	GetStrokeJoin() enums.PaintJoin
	GetStrokeMiter() Scalar
	GetStrokeWidth() Scalar
	GetStyle() enums.PaintStyle
	IsAntiAlias() bool
	IsDither() bool
	IsSrcOver() bool
	NothingToDraw() bool
	Reset()
	SetARGB(a uint8, r uint8, g uint8, b uint8)
	SetAlpha(a uint8)
	SetAlphaf(a Scalar)
	SetAntiAlias(aa bool)
	SetBlendMode(mode enums.BlendMode)
	SetBlender(blender Blender)
	SetColor(color models.Color4f)
	SetColorFilter(filter ColorFilter)
	SetColorInt(color uint32)
	SetDither(dither bool)
	SetImageFilter(filter ImageFilter)
	SetMaskFilter(filter MaskFilter)
	SetPathEffect(effect PathEffect)
	SetShader(shader Shader)
	SetStroke(isStroke bool)
	SetStrokeCap(cap enums.PaintCap)
	SetStrokeJoin(join enums.PaintJoin)
	SetStrokeMiter(limit Scalar)
	SetStrokeWidth(width Scalar)
	SetStyle(style enums.PaintStyle)
}

// PathEffect is the interface for objects that affect the geometry of a drawing primitive
// before it is transformed by the canvas' matrix and drawn.
// Dashing is implemented as a subclass of PathEffect.
type PathEffect interface {
	// ComputeFastBounds computes fast bounds for the path effect.
	// If bounds is nil, returns true if fast bounds computation is possible.
	// If bounds is not nil, modifies bounds in place and returns true if successful.
	ComputeFastBounds(bounds *Rect) bool
}

// Shader specifies the premultiplied source color(s) for what is being drawn.
// If a paint has no shader, then the paint's color is used. If the paint has a
// shader, then the shader's color(s) are used instead, but they are
// modulated by the paint's alpha.
type Shader interface {
	// IsOpaque returns true if the shader is guaranteed to produce only opaque colors,
	// subject to the Paint using the shader to apply an opaque alpha value.
	// This method is optional - implementations may return false if unknown.
	IsOpaque() bool
}

// MaskFilter is the interface for objects that modify the mask before it is used.
// Blur and emboss are implemented as subclasses of MaskFilter.
type MaskFilter interface {
	// ComputeFastBounds computes fast bounds for the mask filter.
	// Modifies storage in place with the adjusted bounds.
	ComputeFastBounds(bounds Rect, storage *Rect)
}

// ColorFilter is the interface for objects that modify colors in the drawing pipeline.
// When present in a paint, they are called with the "src" colors, and return new colors,
// which are then passed onto the next stage (either ImageFilter or BlendMode).
type ColorFilter interface {
	// IsAlphaUnchanged returns true if the filter is guaranteed to never change
	// the alpha of a color it filters.
	IsAlphaUnchanged() bool
}

// ImageFilter is the interface for objects that transform an input image into an output image.
// Image filters can be chained together to create complex effects.
type ImageFilter interface {
	// CanComputeFastBounds returns true if fast bounds computation is possible.
	// Fast bounds computation requires that the filter doesn't affect transparent black.
	CanComputeFastBounds() bool
	// ComputeFastBounds computes fast bounds for the image filter.
	// Returns the adjusted bounds.
	ComputeFastBounds(bounds Rect) Rect
}

// Blender is the interface for objects that specify how source and destination colors
// are combined. It can represent a simple BlendMode or a custom blending operation.
type Blender interface {
	// AsBlendMode returns the blend mode if it can be represented as one.
	// Returns (mode, true) if the blend mode can be determined, (0, false) otherwise.
	AsBlendMode() (enums.BlendMode, bool)
}
