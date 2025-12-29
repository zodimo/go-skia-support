package paragraph

import (
	"sort"

	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// TypefaceFontStyleSet implements interfaces.SkFontStyleSet.
type TypefaceFontStyleSet struct {
	typefaces []interfaces.SkTypeface
}

// NewTypefaceFontStyleSet creates a new TypefaceFontStyleSet.
func NewTypefaceFontStyleSet(typefaces []interfaces.SkTypeface) *TypefaceFontStyleSet {
	return &TypefaceFontStyleSet{
		typefaces: typefaces,
	}
}

func (s *TypefaceFontStyleSet) Count() int {
	return len(s.typefaces)
}

func (s *TypefaceFontStyleSet) GetStyle(index int, style *models.FontStyle, name *string) {
	if index < 0 || index >= len(s.typefaces) {
		return
	}
	tf := s.typefaces[index]
	if style != nil {
		*style = tf.FontStyle()
	}
	if name != nil {
		*name = tf.FamilyName()
	}
}

func (s *TypefaceFontStyleSet) CreateTypeface(index int) interfaces.SkTypeface {
	if index < 0 || index >= len(s.typefaces) {
		return nil
	}
	return s.typefaces[index]
}

func (s *TypefaceFontStyleSet) MatchStyle(pattern models.FontStyle) interfaces.SkTypeface {
	// Simple matching logic: find exact match or just return first for now.
	// TODO: Implement proper CSS3 font matching algorithm or closest match.
	for _, tf := range s.typefaces {
		if tf.FontStyle() == pattern { // Assuming model.FontStyle is comparable
			return tf
		}
	}
	if len(s.typefaces) > 0 {
		return s.typefaces[0]
	}
	return nil
}

func (s *TypefaceFontStyleSet) AppendTypeface(typeface interfaces.SkTypeface) {
	s.typefaces = append(s.typefaces, typeface)
}

// TypefaceFontProvider implements interfaces.SkFontMgr.
// It allows registering typefaces manually.
type TypefaceFontProvider struct {
	families    map[string]*TypefaceFontStyleSet
	familyNames []string
}

// NewTypefaceFontProvider creates a new TypefaceFontProvider.
func NewTypefaceFontProvider() *TypefaceFontProvider {
	return &TypefaceFontProvider{
		families:    make(map[string]*TypefaceFontStyleSet),
		familyNames: make([]string, 0),
	}
}

func (p *TypefaceFontProvider) RegisterTypeface(typeface interfaces.SkTypeface) int {
	return p.RegisterTypefaceWithAlias(typeface, "")
}

func (p *TypefaceFontProvider) RegisterTypefaceWithAlias(typeface interfaces.SkTypeface, alias string) int {
	name := alias
	if name == "" {
		name = typeface.FamilyName()
	}

	set, exists := p.families[name]
	if !exists {
		set = NewTypefaceFontStyleSet(nil)
		p.families[name] = set
		p.familyNames = append(p.familyNames, name)
		sort.Strings(p.familyNames) // Keep names sorted for deterministic index
	}
	set.AppendTypeface(typeface)
	return 1 // Return count? C++ returns index or similar? API check: RegisterTypeface returns size_t (count)
}

func (p *TypefaceFontProvider) CountFamilies() int {
	return len(p.familyNames)
}

func (p *TypefaceFontProvider) GetFamilyName(index int) string {
	if index < 0 || index >= len(p.familyNames) {
		return ""
	}
	return p.familyNames[index]
}

func (p *TypefaceFontProvider) CreateStyleSet(index int) interfaces.SkFontStyleSet {
	if index < 0 || index >= len(p.familyNames) {
		return nil
	}
	name := p.familyNames[index]
	return p.families[name]
}

func (p *TypefaceFontProvider) MatchFamily(familyName string) interfaces.SkFontStyleSet {
	if set, ok := p.families[familyName]; ok {
		return set
	}
	// Also try case-insensitive match or aliases if needed, but strict for now
	return nil
}

func (p *TypefaceFontProvider) MatchFamilyStyle(familyName string, style models.FontStyle) interfaces.SkTypeface {
	set := p.MatchFamily(familyName)
	if set == nil {
		return nil
	}
	return set.MatchStyle(style)
}

func (p *TypefaceFontProvider) MatchFamilyStyleCharacter(familyName string, style models.FontStyle, bcp47 []string, character rune) interfaces.SkTypeface {
	// Fallback logic not implemented for explicit character matching in provider yet
	return p.MatchFamilyStyle(familyName, style)
}

func (p *TypefaceFontProvider) MakeFromData(data interfaces.SkData, ttcIndex int) interfaces.SkTypeface {
	// Not implemented for provider - provider manages existing typefaces
	return nil
}

func (p *TypefaceFontProvider) MakeFromFile(path string, ttcIndex int) interfaces.SkTypeface {
	// Not implemented
	return nil
}

func (p *TypefaceFontProvider) LegacyMakeTypeface(familyName string, style models.FontStyle) interfaces.SkTypeface {
	return p.MatchFamilyStyle(familyName, style)
}
