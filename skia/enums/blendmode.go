package enums

// BlendMode represents the blend mode for compositing
type BlendMode uint8

const (
	BlendModeClear    BlendMode = iota // r = 0
	BlendModeSrc                       // r = s
	BlendModeDst                       // r = d
	BlendModeSrcOver                   // r = s + (1-sa)*d
	BlendModeDstOver                   // r = d + (1-da)*s
	BlendModeSrcIn                     // r = s * da
	BlendModeDstIn                     // r = d * sa
	BlendModeSrcOut                    // r = s * (1-da)
	BlendModeDstOut                    // r = d * (1-sa)
	BlendModeSrcATop                   // r = s*da + d*(1-sa)
	BlendModeDstATop                   // r = d*sa + s*(1-da)
	BlendModeXor                       // r = s*(1-da) + d*(1-sa)
	BlendModePlus                      // r = min(s + d, 1)
	BlendModeModulate                  // r = s*d
	BlendModeScreen                    // r = s + d - s*d

	BlendModeOverlay    // multiply or screen, depending on destination
	BlendModeDarken     // rc = s + d - max(s*da, d*sa), ra = kSrcOver
	BlendModeLighten    // rc = s + d - min(s*da, d*sa), ra = kSrcOver
	BlendModeColorDodge // brighten destination to reflect source
	BlendModeColorBurn  // darken destination to reflect source
	BlendModeHardLight  // multiply or screen, depending on source
	BlendModeSoftLight  // lighten or darken, depending on source
	BlendModeDifference // rc = s + d - 2*(min(s*da, d*sa)), ra = kSrcOver
	BlendModeExclusion  // rc = s + d - two(s*d), ra = kSrcOver
	BlendModeMultiply   // r = s*(1-da) + d*(1-sa) + s*d

	BlendModeHue        // hue of source with saturation and luminosity of destination
	BlendModeSaturation // saturation of source with hue and luminosity of destination
	BlendModeColor      // hue and saturation of source with luminosity of destination
	BlendModeLuminosity // luminosity of source with hue and saturation of destination

	BlendModeLastCoeffMode     = BlendModeScreen     // last porter duff blend mode
	BlendModeLastSeparableMode = BlendModeMultiply   // last blend mode operating separately on components
	BlendModeLast              = BlendModeLuminosity // last valid value
)
