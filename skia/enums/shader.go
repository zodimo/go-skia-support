// Package enums provides shader-related enumerations.
// Ported from SkShaderBase.h
// https://github.com/google/skia/blob/main/src/shaders/SkShaderBase.h

package enums

// ShaderType identifies the kind of shader.
// Matches C++ SkShaderBase::ShaderType enum.
type ShaderType int

const (
	ShaderTypeBlend ShaderType = iota
	ShaderTypeCTM
	ShaderTypeColor
	ShaderTypeColorFilter
	ShaderTypeCoordClamp
	ShaderTypeEmpty
	ShaderTypeGradientBase
	ShaderTypeImage
	ShaderTypeLocalMatrix
	ShaderTypePerlinNoise
	ShaderTypePicture
	ShaderTypeRuntime
	ShaderTypeTransform
	ShaderTypeTriColor
	ShaderTypeWorkingColorSpace
)

// GradientType identifies the kind of gradient shader.
// Matches C++ SkShaderBase::GradientType enum.
type GradientType int

const (
	// GradientTypeNone indicates the shader is not a gradient.
	GradientTypeNone GradientType = iota
	// GradientTypeLinear is a linear gradient with two end-points.
	GradientTypeLinear
	// GradientTypeRadial is a radial gradient with a center and radius.
	GradientTypeRadial
	// GradientTypeConical is a two-point conical gradient.
	GradientTypeConical
	// GradientTypeSweep is a sweep (angular) gradient around a center point.
	GradientTypeSweep
)

// GradientFlags contains optional flags for gradient shaders.
// Matches C++ SkGradientShader::Flags.
type GradientFlags uint32

const (
	// GradientFlagInterpolateColorsInPremul indicates colors are interpolated in premul space.
	GradientFlagInterpolateColorsInPremul GradientFlags = 1 << 0
)
