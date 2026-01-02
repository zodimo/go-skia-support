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
// modulated by the paint's alpha. This makes it easy to create a shader
// once (e.g., bitmap tiling or gradient) and then change its transparency
// without having to modify the original shader — only the paint's alpha needs
// to be modified.
//
// Ported from: include/core/SkShader.h
// https://github.com/google/skia/blob/main/include/core/SkShader.h
type Shader interface {
	// =========================================================================
	// Core SkShader Methods (from include/core/SkShader.h)
	// =========================================================================

	// IsOpaque returns true if the shader is guaranteed to produce only opaque
	// colors, subject to the Paint using the shader to apply an opaque alpha value.
	// Subclasses should override this to allow some optimizations.
	IsOpaque() bool

	// IsAImage returns true if this shader is backed by a single SkImage.
	// If localMatrix is non-nil and this returns true, localMatrix is set to
	// the shader's local matrix. If tileMode is non-nil and this returns true,
	// tileMode[0] and tileMode[1] are set to the x and y tile modes.
	IsAImage(localMatrix *SkMatrix, tileMode []enums.TileMode) bool

	// IsAImageSimple returns true if this shader is backed by a single SkImage.
	// This is the simple form that doesn't extract matrix or tile modes.
	IsAImageSimple() bool

	// MakeWithLocalMatrix returns a shader that will apply the specified localMatrix
	// to this shader. The specified matrix will be applied before any matrix
	// associated with this shader.
	MakeWithLocalMatrix(localMatrix SkMatrix) Shader

	// MakeWithColorFilter creates a new shader that produces the same colors as
	// invoking this shader and then applying the colorfilter.
	MakeWithColorFilter(filter ColorFilter) Shader

	// MakeWithWorkingColorSpace returns a shader that will compute this shader
	// in a context such that any child shaders return RGBA values converted to
	// the inputCS colorspace.
	//
	// It is assumed that the RGBA values returned by this shader have been
	// transformed into outputCS. By default, shaders are assumed to return values
	// in the destination colorspace and premultiplied.
	//
	// A nil inputCS is assumed to be the destination CS.
	// A nil outputCS is assumed to be the inputCS.
	MakeWithWorkingColorSpace(inputCS, outputCS *models.ColorSpace) Shader

	// =========================================================================
	// Extended SkShaderBase Methods (from src/shaders/SkShaderBase.h)
	// =========================================================================

	// UniqueID returns a value unique to this shader instance.
	// Used for caching and equality comparisons.
	UniqueID() uint32

	// IsConstant returns true if the shader is guaranteed to produce only a single color.
	// If color is non-nil and this returns true, color is set to that constant color.
	// Subclasses can override this to allow loop-hoisting optimization.
	IsConstant(color *models.Color4f) bool

	// Type returns the ShaderType enum identifying this shader's concrete type.
	Type() enums.ShaderType

	// AsGradient returns the GradientType if this shader can be represented as a gradient.
	// Returns GradientTypeNone if it cannot. If info is non-nil, populates it with
	// the gradient parameters. If localMatrix is non-nil, populates it with the
	// shader's local matrix.
	//
	// See models.GradientInfo for details on how the info struct is populated
	// for different gradient types.
	AsGradient(info *models.GradientInfo, localMatrix *SkMatrix) enums.GradientType

	// =========================================================================
	// Utility Methods
	// =========================================================================

	// MakeInvertAlpha returns a shader with inverted alpha.
	MakeInvertAlpha() Shader

	// MakeWithCTM returns a shader that owns its own CTM (current transform matrix).
	MakeWithCTM(ctm SkMatrix) Shader
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
