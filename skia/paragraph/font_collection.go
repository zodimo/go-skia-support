package paragraph

import (
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// FontCollection manages a collection of font managers and resolves typefaces.
//
// Ported from: skia-source/modules/skparagraph/include/FontCollection.h
type FontCollection struct {
	typefaces          map[string][]interfaces.SkTypeface
	defaultFontManager interfaces.SkFontMgr
	assetFontManager   interfaces.SkFontMgr
	dynamicFontManager interfaces.SkFontMgr
	testFontManager    interfaces.SkFontMgr
	enableFontFallback bool
	paragraphCache     *ParagraphCache
}

// NewFontCollection creates a new FontCollection.
func NewFontCollection() *FontCollection {
	return &FontCollection{
		typefaces:          make(map[string][]interfaces.SkTypeface),
		enableFontFallback: true,
		paragraphCache:     NewParagraphCache(),
	}
}

// GetFontManagersCount returns the number of registered font managers.
func (fc *FontCollection) GetFontManagersCount() int {
	return len(fc.getFontManagerOrder())
}

func (fc *FontCollection) getFontManagerOrder() []interfaces.SkFontMgr {
	order := make([]interfaces.SkFontMgr, 0, 4)
	if fc.dynamicFontManager != nil {
		order = append(order, fc.dynamicFontManager)
	}
	if fc.assetFontManager != nil {
		order = append(order, fc.assetFontManager)
	}
	if fc.testFontManager != nil {
		order = append(order, fc.testFontManager)
	}
	// Note: The C++ implementation checks enableFontFallback here for the default manager
	if fc.defaultFontManager != nil && fc.enableFontFallback {
		order = append(order, fc.defaultFontManager)
	}
	return order
}

// SetAssetFontManager sets the asset font manager.
func (fc *FontCollection) SetAssetFontManager(fontManager interfaces.SkFontMgr) {
	fc.assetFontManager = fontManager
}

// SetDynamicFontManager sets the dynamic font manager.
func (fc *FontCollection) SetDynamicFontManager(fontManager interfaces.SkFontMgr) {
	fc.dynamicFontManager = fontManager
}

// SetTestFontManager sets the test font manager.
func (fc *FontCollection) SetTestFontManager(fontManager interfaces.SkFontMgr) {
	fc.testFontManager = fontManager
}

// SetDefaultFontManager sets the default font manager.
func (fc *FontCollection) SetDefaultFontManager(fontManager interfaces.SkFontMgr) {
	fc.defaultFontManager = fontManager
}

// GetFallbackManager returns the fallback font manager (usually the default one).
func (fc *FontCollection) GetFallbackManager() interfaces.SkFontMgr {
	return fc.defaultFontManager
}

// FindTypefaces finds typefaces for the given family names and style.
func (fc *FontCollection) FindTypefaces(familyNames []string, fontStyle models.FontStyle) []interfaces.SkTypeface {
	key := fc.makeFamilyKey(familyNames, fontStyle)
	if cached, ok := fc.typefaces[key]; ok {
		return cached
	}

	var typefaces []interfaces.SkTypeface
	managers := fc.getFontManagerOrder()

	for _, family := range familyNames {
		match := fc.matchTypeface(family, fontStyle, managers)
		if match != nil {
			typefaces = append(typefaces, match)
		}
	}

	if len(typefaces) == 0 {
		match := fc.matchTypeface("", fontStyle, managers)
		if match == nil {
			for _, manager := range managers {
				match = manager.LegacyMakeTypeface("", fontStyle)
				if match != nil {
					break
				}
			}
		}
		if match != nil {
			typefaces = append(typefaces, match)
		}
	}

	fc.typefaces[key] = typefaces
	return typefaces
}

func (fc *FontCollection) matchTypeface(familyName string, fontStyle models.FontStyle, managers []interfaces.SkFontMgr) interfaces.SkTypeface {
	for _, manager := range managers {
		match := manager.MatchFamilyStyle(familyName, fontStyle)
		if match != nil {
			return match
		}
	}
	return nil
}

func (fc *FontCollection) makeFamilyKey(familyNames []string, fontStyle models.FontStyle) string {
	// Simple key generation strategy
	// Joined family names + style attributes
	key := ""
	for _, f := range familyNames {
		key += f + "|"
	}
	key += string(rune(fontStyle.Weight)) + "|"
	key += string(rune(fontStyle.Width)) + "|"
	key += string(rune(fontStyle.Slant))
	return key
}

// DefaultFallback finds a fallback typeface for the given unicode character.
func (fc *FontCollection) DefaultFallback(unicode rune, fontStyle models.FontStyle, locale string) interfaces.SkTypeface {
	for _, manager := range fc.getFontManagerOrder() {
		// Go strings are UTF-8, but locally we just pass the slice.
		// simplified bcp47 handling
		locales := []string{}
		if locale != "" {
			locales = append(locales, locale)
		}
		match := manager.MatchFamilyStyleCharacter("", fontStyle, locales, unicode)
		if match != nil {
			return match
		}
	}
	return nil
}

// DefaultFallbackTypeface returns the default fallback typeface.
func (fc *FontCollection) DefaultFallbackTypeface() interfaces.SkTypeface {
	if fc.defaultFontManager == nil {
		return nil
	}
	return fc.defaultFontManager.MatchFamilyStyle("", models.FontStyle{})
}

// DefaultEmojiFallback finds an emoji font.
func (fc *FontCollection) DefaultEmojiFallback(emojiStart rune, fontStyle models.FontStyle, locale string) interfaces.SkTypeface {
	// Simplified implementation: Look for common emoji fonts or use available managers
	emojiFonts := []string{"Apple Color Emoji", "Noto Color Emoji", "Segoe UI Emoji"}
	managers := fc.getFontManagerOrder()

	for _, rangeName := range emojiFonts {
		for _, manager := range managers {
			match := manager.MatchFamilyStyle(rangeName, fontStyle)
			if match != nil {
				return match
			}
		}
	}
	// Fallback to character matching if specific family not found
	return fc.DefaultFallback(emojiStart, fontStyle, locale)
}

// DisableFontFallback disables font fallback.
func (fc *FontCollection) DisableFontFallback() {
	fc.enableFontFallback = false
}

// EnableFontFallback enables font fallback.
func (fc *FontCollection) EnableFontFallback() {
	fc.enableFontFallback = true
}

// FontFallbackEnabled returns true if font fallback is enabled.
func (fc *FontCollection) FontFallbackEnabled() bool {
	return fc.enableFontFallback
}

// ClearCaches clears the caches.
func (fc *FontCollection) ClearCaches() {
	fc.typefaces = make(map[string][]interfaces.SkTypeface)
	fc.paragraphCache = NewParagraphCache() // Reset paragraph cache
}

// GetParagraphCache returns the paragraph cache.
func (fc *FontCollection) GetParagraphCache() *ParagraphCache {
	return fc.paragraphCache
}
