package enums

// AlphaType describes how pixel bits encode color
// Matches C++ SkAlphaType enum from include/core/SkAlphaType.h
type AlphaType int

const (
	AlphaTypeUnknown  AlphaType = 0
	AlphaTypeOpaque   AlphaType = 1 // pixel is opaque
	AlphaTypePremul   AlphaType = 2 // pixel components are premultiplied by alpha
	AlphaTypeUnpremul AlphaType = 3 // pixel components are independent of alpha
	AlphaTypeLast     AlphaType = AlphaTypeUnpremul
)

// ColorType describes how pixel bits encode color
// Matches C++ SkColorType enum from include/core/SkImageInfo.h
type ColorType int

const (
	ColorTypeUnknown     ColorType = 0
	ColorTypeAlpha8      ColorType = 1  // 8-bit alpha
	ColorTypeRGB565      ColorType = 2  // 5-6-5 RGB
	ColorTypeARGB4444    ColorType = 3  // 4-4-4-4 ARGB
	ColorTypeRGBA8888    ColorType = 4  // 8-8-8-8 RGBA
	ColorTypeRGB888x     ColorType = 5  // 8-8-8-x RGB
	ColorTypeBGRA8888    ColorType = 6  // 8-8-8-8 BGRA
	ColorTypeRGBA1010102 ColorType = 7  // 10-10-10-2 RGBA
	ColorTypeRGB101010x  ColorType = 8  // 10-10-10-x RGB
	ColorTypeGray8       ColorType = 9  // 8-bit gray
	ColorTypeRGBAF16Norm ColorType = 10 // 16-bit float normalized RGBA
	ColorTypeRGBAF16     ColorType = 11 // 16-bit float RGBA
	ColorTypeRGBAF32     ColorType = 12 // 32-bit float RGBA (unclamped)

	// Legacy aliases
	ColorTypeN32 = ColorTypeBGRA8888 // assuming little-endian/standard for now, can be platform dependent
)

// TileMode rules for drawing outside of the image bounds
// Matches C++ SkTileMode enum from include/core/SkTileMode.h
type TileMode int

const (
	TileModeClamp  TileMode = 0 // replicate the edge color
	TileModeRepeat TileMode = 1 // repeat the image
	TileModeMirror TileMode = 2 // mirror the image
	TileModeDecal  TileMode = 3 // draw transparent black
	TileModeLast   TileMode = TileModeDecal
)

// TextureCompressionType describes the format of compressed texture data
// Matches C++ SkTextureCompressionType enum
type TextureCompressionType int

const (
	TextureCompressionTypeNone            TextureCompressionType = 0
	TextureCompressionTypeETC2_RGB8_UNORM TextureCompressionType = 1
	TextureCompressionTypeBC1_RGB8_UNORM  TextureCompressionType = 2
)

// FilterMode describes the sampling quality
// Matches C++ SkFilterMode (subset of SkSamplingOptions)
type FilterMode int

const (
	FilterModeNearest FilterMode = 0 // single sample point (nearest neighbor)
	FilterModeLinear  FilterMode = 1 // interporate between 2x2 samples (bilinear)
	FilterModeLast    FilterMode = FilterModeLinear
)

// MipmapMode describes the sampling quality via mipmaps
// Matches C++ SkMipmapMode (subset of SkSamplingOptions)
type MipmapMode int

const (
	MipmapModeNone    MipmapMode = 0 // ignore mipmaps
	MipmapModeNearest MipmapMode = 1 // nearest mipmap level
	MipmapModeLinear  MipmapMode = 2 // interpolate between mipmap levels
	MipmapModeLast    MipmapMode = MipmapModeLinear
)
