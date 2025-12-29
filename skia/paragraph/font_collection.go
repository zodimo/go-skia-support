package paragraph

import (
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// FontCollection manages a collection of font managers and resolves typefaces.
//
// Ported from: skia-source/modules/skparagraph/include/FontCollection.h
type FontCollection struct {
	fontManagers       []interfaces.SkFontMgr
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
		fontManagers:       make([]interfaces.SkFontMgr, 0),
		enableFontFallback: true,
		paragraphCache:     NewParagraphCache(),
	}
}

// GetFontManagersCount returns the number of registered font managers.
func (fc *FontCollection) GetFontManagersCount() int {
	return len(fc.fontManagers)
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
	var typefaces []interfaces.SkTypeface

	// Collect all managers in order of priority (or just all of them).
	// Skia usually checks check them in specific order.
	managers := []interfaces.SkFontMgr{}
	if fc.assetFontManager != nil {
		managers = append(managers, fc.assetFontManager)
	}
	if fc.dynamicFontManager != nil {
		managers = append(managers, fc.dynamicFontManager)
	}
	if fc.testFontManager != nil {
		managers = append(managers, fc.testFontManager)
	}
	managers = append(managers, fc.fontManagers...)
	if fc.defaultFontManager != nil {
		managers = append(managers, fc.defaultFontManager)
	}

	for _, manager := range managers {
		for _, family := range familyNames {
			tf := manager.MatchFamilyStyle(family, fontStyle)
			if tf != nil {
				typefaces = append(typefaces, tf)
			}
		}
	}

	return typefaces
}

// DefaultFallback finds a fallback typeface for the given unicode character.
func (fc *FontCollection) DefaultFallback(unicode rune, fontStyle models.FontStyle, locale string) interfaces.SkTypeface {
	if !fc.enableFontFallback || fc.defaultFontManager == nil {
		return nil
	}
	// Note: bcp47 handling is simplified here.
	return fc.defaultFontManager.MatchFamilyStyleCharacter("", fontStyle, []string{locale}, unicode)
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

// GetParagraphCache returns the paragraph cache.
func (fc *FontCollection) GetParagraphCache() *ParagraphCache {
	return fc.paragraphCache
}
