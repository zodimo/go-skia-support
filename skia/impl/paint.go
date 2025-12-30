package impl

import (
	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// PaintStyleCount is the number of valid style values
const PaintStyleCount = enums.PaintStyleStrokeAndFill + 1

// PaintJoinCount is the number of valid join values
const PaintJoinCount = enums.PaintJoinLast + 1

// BlendModeCount is the number of valid blend mode values
const BlendModeCount = int(enums.BlendModeLast) + 1

// PaintDefaultsMiterLimit is the default miter limit value
const PaintDefaultsMiterLimit Scalar = 4.0

// Paint represents a paint object that specifies how geometry is drawn
type Paint struct {
	// Effect and filter objects
	PathEffect  PathEffect  // sk_sp<SkPathEffect>
	Shader      Shader      // sk_sp<SkShader>
	MaskFilter  MaskFilter  // sk_sp<SkMaskFilter>
	ColorFilter ColorFilter // sk_sp<SkColorFilter>
	ImageFilter ImageFilter // sk_sp<SkImageFilter>
	Blender     Blender     // sk_sp<SkBlender>

	// Blend mode (stored when SetBlendMode is called)
	// If Blender is nil, this represents the blend mode (defaults to BlendModeSrcOver)
	// If Blender is not nil, the blender takes precedence and this may be invalid
	blendMode *enums.BlendMode // nil means use default (SrcOver)

	// Color and stroke properties
	Color4f    Color4f // RGBA color (unpremultiplied)
	Width      Scalar  // stroke width
	MiterLimit Scalar  // miter limit

	// Bitfields for flags and enums
	Bitfields PaintBitfields
}

// NewPaint creates a new Paint with default values
func NewPaint() *Paint {
	return &Paint{
		PathEffect:  nil,
		Shader:      nil,
		MaskFilter:  nil,
		ColorFilter: nil,
		ImageFilter: nil,
		Blender:     nil,
		Color4f: Color4f{
			R: 0,
			G: 0,
			B: 0,
			A: 1, // opaque black
		},
		Width:      0, // hairline
		MiterLimit: PaintDefaultsMiterLimit,
		Bitfields: PaintBitfields{
			AntiAlias: false,
			Dither:    false,
			CapType:   enums.PaintCapDefault,
			JoinType:  enums.PaintJoinDefault,
			Style:     enums.PaintStyleFill,
		},
	}
}

// NewPaintWithColor creates a new Paint with the specified color
func NewPaintWithColor(color Color4f) *Paint {
	p := NewPaint()
	p.SetColor(color)
	return p
}

// Reset resets the paint to default values
func (p *Paint) Reset() {
	*p = *NewPaint()
	// blendMode is already nil from NewPaint()
}

// Equals compares two paints for equality
func (p *Paint) Equals(other interfaces.SkPaint) bool {
	if p == nil && other == nil {
		return true
	}
	if p == nil || other == nil {
		return false
	}

	// Cast to *Paint to access fields
	otherPaint, ok := other.(*Paint)
	if !ok {
		return false
	}

	// Compare blendMode (handle nil cases)
	var blendModeEqual bool
	if p.blendMode == nil && otherPaint.blendMode == nil {
		blendModeEqual = true
	} else if p.blendMode != nil && otherPaint.blendMode != nil {
		blendModeEqual = *p.blendMode == *otherPaint.blendMode
	} else {
		blendModeEqual = false
	}
	return p.PathEffect == otherPaint.PathEffect &&
		p.Shader == otherPaint.Shader &&
		p.MaskFilter == otherPaint.MaskFilter &&
		p.ColorFilter == otherPaint.ColorFilter &&
		p.Blender == otherPaint.Blender &&
		p.ImageFilter == otherPaint.ImageFilter &&
		blendModeEqual &&
		p.Color4f == otherPaint.Color4f &&
		p.Width == otherPaint.Width &&
		p.MiterLimit == otherPaint.MiterLimit &&
		p.Bitfields == otherPaint.Bitfields
}

// SetColor sets the paint color
// For now, color space transformation is deferred
func (p *Paint) SetColor(color Color4f) {
	p.Color4f = color.PinAlpha()
	// TODO: Apply color space transformation if colorSpace provided
}

// SetColorInt sets the paint color from a uint32 SkColor
func (p *Paint) SetColorInt(color uint32) {
	p.SetColor(Color4fFromColor(color))
}

// SetARGB sets the color from ARGB components
func (p *Paint) SetARGB(a, r, g, b uint8) {
	p.SetColor(Color4f{
		R: Scalar(r) / 255.0,
		G: Scalar(g) / 255.0,
		B: Scalar(b) / 255.0,
		A: Scalar(a) / 255.0,
	})
}

// GetColor returns the current color as Color4f
func (p *Paint) GetColor() models.Color4f {
	return p.Color4f
}

// GetColorInt returns the current color as a uint32 SkColor (ARGB format)
// This is equivalent to SkPaint::getColor() in C++
func (p *Paint) GetColorInt() uint32 {
	return p.Color4f.ToSkColor()
}

// GetAlphaf returns the alpha component as a float (0.0 to 1.0)
func (p *Paint) GetAlphaf() Scalar {
	return p.Color4f.A
}

// GetAlpha returns the alpha component as a uint8 (0 to 255)
// This is equivalent to SkPaint::getAlpha() in C++
func (p *Paint) GetAlpha() uint8 {
	// Round to nearest integer: round(alpha * 255)
	return uint8(scalarPin(p.Color4f.A, 0.0, 1.0)*255.0 + 0.5)
}

// SetAlphaf sets the alpha component, clamping to [0, 1]
func (p *Paint) SetAlphaf(a Scalar) {
	p.Color4f.A = scalarPin(a, 0.0, 1.0)
}

// SetAlpha sets the alpha component from a uint8 (0 to 255)
// This is equivalent to SkPaint::setAlpha() in C++
func (p *Paint) SetAlpha(a uint8) {
	p.SetAlphaf(Scalar(a) / 255.0)
}

// GetStyle returns the current style
func (p *Paint) GetStyle() enums.PaintStyle {
	return p.Bitfields.Style
}

// SetStyle sets the paint style
func (p *Paint) SetStyle(style enums.PaintStyle) {
	if style < PaintStyleCount {
		p.Bitfields.Style = style
	}
}

// SetStroke sets style to stroke or fill based on boolean
func (p *Paint) SetStroke(isStroke bool) {
	if isStroke {
		p.SetStyle(enums.PaintStyleStroke)
	} else {
		p.SetStyle(enums.PaintStyleFill)
	}
}

// GetStrokeWidth returns the stroke width
func (p *Paint) GetStrokeWidth() Scalar {
	return p.Width
}

// SetStrokeWidth sets the stroke width (must be >= 0)
func (p *Paint) SetStrokeWidth(width Scalar) {
	if width < 0 {
		width = 0
	}
	p.Width = width
}

// GetStrokeMiter returns the miter limit
func (p *Paint) GetStrokeMiter() Scalar {
	return p.MiterLimit
}

// SetStrokeMiter sets the miter limit (must be >= 0)
func (p *Paint) SetStrokeMiter(limit Scalar) {
	if limit < 0 {
		limit = 0
	}
	p.MiterLimit = limit
}

// GetStrokeCap returns the stroke cap style
func (p *Paint) GetStrokeCap() enums.PaintCap {
	return p.Bitfields.CapType
}

// SetStrokeCap sets the stroke cap style
func (p *Paint) SetStrokeCap(cap enums.PaintCap) {
	if cap < enums.PaintCapCount {
		p.Bitfields.CapType = cap
	}
}

// GetStrokeJoin returns the stroke join style
func (p *Paint) GetStrokeJoin() enums.PaintJoin {
	return p.Bitfields.JoinType
}

// SetStrokeJoin sets the stroke join style
func (p *Paint) SetStrokeJoin(join enums.PaintJoin) {
	if join < PaintJoinCount {
		p.Bitfields.JoinType = join
	}
}

// IsAntiAlias returns the anti-aliasing state
func (p *Paint) IsAntiAlias() bool {
	return p.Bitfields.AntiAlias
}

// SetAntiAlias sets the anti-aliasing state
func (p *Paint) SetAntiAlias(aa bool) {
	p.Bitfields.AntiAlias = aa
}

// IsDither returns the dithering state
func (p *Paint) IsDither() bool {
	return p.Bitfields.Dither
}

// SetDither sets the dithering state
func (p *Paint) SetDither(dither bool) {
	p.Bitfields.Dither = dither
}

// SetShader sets the shader
func (p *Paint) SetShader(shader Shader) {
	p.Shader = shader
}

// GetShader returns the current shader
func (p *Paint) GetShader() Shader {
	return p.Shader
}

// SetColorFilter sets the color filter
func (p *Paint) SetColorFilter(filter ColorFilter) {
	p.ColorFilter = filter
}

// GetColorFilter returns the current color filter
func (p *Paint) GetColorFilter() ColorFilter {
	return p.ColorFilter
}

// SetPathEffect sets the path effect
func (p *Paint) SetPathEffect(effect PathEffect) {
	p.PathEffect = effect
}

// GetPathEffect returns the current path effect
func (p *Paint) GetPathEffect() PathEffect {
	return p.PathEffect
}

// SetMaskFilter sets the mask filter
func (p *Paint) SetMaskFilter(filter MaskFilter) {
	p.MaskFilter = filter
}

// GetMaskFilter returns the current mask filter
func (p *Paint) GetMaskFilter() MaskFilter {
	return p.MaskFilter
}

// SetImageFilter sets the image filter
func (p *Paint) SetImageFilter(filter ImageFilter) {
	p.ImageFilter = filter
}

// GetImageFilter returns the current image filter
func (p *Paint) GetImageFilter() ImageFilter {
	return p.ImageFilter
}

// SetBlender sets the blender
func (p *Paint) SetBlender(blender Blender) {
	p.Blender = blender
	// When a custom blender is set, clear the stored blend mode
	// since the blender takes precedence
	if blender != nil {
		p.blendMode = nil
	}
}

// GetBlender returns the current blender
func (p *Paint) GetBlender() Blender {
	return p.Blender
}

// SetBlendMode sets the blend mode (convenience method)
// If mode is BlendModeSrcOver, sets blender to nil
// Otherwise, stores the mode (blender creation is deferred until Blender type is implemented)
func (p *Paint) SetBlendMode(mode enums.BlendMode) {
	if mode == enums.BlendModeSrcOver {
		p.Blender = nil
		p.blendMode = nil // nil means default (SrcOver)
	} else {
		// Store the mode for now
		// TODO: Create a blender for the mode when Blender type is implemented
		p.blendMode = &mode
		// For now, we can't create a blender, so Blender remains nil
		// This is a temporary limitation until Blender type is implemented
	}
}

// AsBlendMode returns the blend mode if it can be represented as one
// Returns (mode, true) if the blend mode can be determined, (0, false) otherwise
func (p *Paint) AsBlendMode() (enums.BlendMode, bool) {
	if p.Blender == nil {
		// If blender is nil, use stored blend mode or default to SrcOver
		if p.blendMode != nil {
			return *p.blendMode, true
		}
		return enums.BlendModeSrcOver, true
	}
	// If blender is not nil, get the blend mode from the blender
	return p.Blender.AsBlendMode()
}

// GetBlendModeOr returns the blend mode or default if not representable
func (p *Paint) GetBlendModeOr(defaultMode enums.BlendMode) enums.BlendMode {
	mode, ok := p.AsBlendMode()
	if ok {
		return mode
	}
	return defaultMode
}

// IsSrcOver returns true if blend mode is SrcOver or blender is nil
func (p *Paint) IsSrcOver() bool {
	if p.Blender == nil {
		// If blender is nil, check if stored mode is SrcOver or nil (defaults to SrcOver)
		if p.blendMode == nil {
			return true
		}
		return *p.blendMode == enums.BlendModeSrcOver
	}
	// If blender is not nil, check if it represents SrcOver
	mode, ok := p.AsBlendMode()
	if ok {
		return mode == enums.BlendModeSrcOver
	}
	// If we can't determine the blend mode, assume it's not SrcOver
	return false
}

// CanComputeFastBounds returns true if fast bounds computation is possible
// Fast bounds computation requires that ImageFilter and PathEffect (if present)
// support fast bounds computation.
func (p *Paint) CanComputeFastBounds() bool {
	// Check if ImageFilter can compute fast bounds
	if p.ImageFilter != nil {
		if !p.ImageFilter.CanComputeFastBounds() {
			return false
		}
	}

	// Check if PathEffect can compute fast bounds
	// Pass nil to determine if it can compute fast bounds
	if p.PathEffect != nil {
		if !p.PathEffect.ComputeFastBounds(nil) {
			return false
		}
	}

	return true
}

// GetInflationRadius computes the inflation radius for stroke effects.
// This is equivalent to creating a SkStrokeRec from the paint and style,
// then calling getInflationRadius() on it.
// If matrixScale is provided and > 0, it will be used for hairline strokes (width == 0).
// Otherwise, hairlines default to 1.0.
func (p *Paint) GetInflationRadius(style enums.PaintStyle, matrixScale ...Scalar) Scalar {
	var strokeWidth Scalar
	if style == enums.PaintStyleFill {
		strokeWidth = -1.0 // negative indicates fill
	} else {
		strokeWidth = p.Width
	}
	var scale Scalar = 0 // 0 means not provided
	if len(matrixScale) > 0 && matrixScale[0] > 0 {
		scale = matrixScale[0]
	}
	return GetInflationRadiusForStroke(p.Bitfields.JoinType, p.MiterLimit, p.Bitfields.CapType, strokeWidth, scale)
}

// ComputeFastBounds computes fast bounds for geometry with paint effects.
// The original bounds are adjusted to account for stroke, path effects, mask filters, and image filters.
// The storage parameter must not be nil and will be used to store the result.
func (p *Paint) ComputeFastBounds(orig Rect, storage *Rect) Rect {
	if storage == nil {
		panic("storage must not be nil")
	}
	// Things like stroking, etc... will do math on the bounds rect, assuming that it's sorted.
	// In debug builds, we could assert orig.IsSorted(), but for now we'll just proceed.

	style := p.GetStyle()
	// Ultra fast-case: filling with no effects that affect geometry
	if style == enums.PaintStyleFill {
		if p.MaskFilter == nil && p.PathEffect == nil && p.ImageFilter == nil {
			*storage = orig
			return orig
		}
	}

	return p.DoComputeFastBounds(orig, storage, style)
}

// ComputeFastStrokeBounds computes fast bounds for geometry with paint effects,
// using stroke style regardless of the current paint style.
// This is equivalent to SkPaint::computeFastStrokeBounds() in C++
func (p *Paint) ComputeFastStrokeBounds(orig Rect, storage *Rect) Rect {
	return p.DoComputeFastBounds(orig, storage, enums.PaintStyleStroke)
}

// DoComputeFastBounds is an internal method to compute fast bounds with style override.
// It applies PathEffect bounds if present, computes inflation radius for stroke,
// applies MaskFilter bounds if present, and applies ImageFilter bounds if present.
func (p *Paint) DoComputeFastBounds(origSrc Rect, storage *Rect, style enums.PaintStyle) Rect {
	if storage == nil {
		panic("storage must not be nil")
	}

	src := origSrc

	// Apply PathEffect bounds if present
	if p.PathEffect != nil {
		tmpSrc := origSrc
		if p.PathEffect.ComputeFastBounds(&tmpSrc) {
			src = tmpSrc
		}
	}

	// Compute inflation radius for stroke
	radius := p.GetInflationRadius(style)
	*storage = src.MakeOutset(radius, radius)

	// Apply MaskFilter bounds if present
	if p.MaskFilter != nil {
		p.MaskFilter.ComputeFastBounds(*storage, storage)
	}

	// Apply ImageFilter bounds if present
	if p.ImageFilter != nil {
		*storage = p.ImageFilter.ComputeFastBounds(*storage)
	}

	return *storage
}

// NothingToDraw returns true if paint prevents all drawing
// This happens when:
// - Blend mode is kDst (always draws nothing)
// - Blend mode is kSrcOver, kSrcATop, kDstOut, kDstOver, or kPlus AND alpha is 0 AND filters don't affect alpha
func (p *Paint) NothingToDraw() bool {
	mode, ok := p.AsBlendMode()
	if !ok {
		// If we can't determine the blend mode, assume something will be drawn
		return false
	}

	switch mode {
	case enums.BlendModeSrcOver, enums.BlendModeSrcATop, enums.BlendModeDstOut, enums.BlendModeDstOver, enums.BlendModePlus:
		// For these modes, if alpha is 0 and filters don't affect alpha, nothing to draw
		if p.GetAlphaf() == 0 {
			return !affectsAlphaColorFilter(p.ColorFilter) && !affectsAlphaImageFilter(p.ImageFilter)
		}
	case enums.BlendModeDst:
		// kDst always draws nothing (just returns destination)
		return true
	}

	return false
}
