package paragraph

// FontFeature represents an OpenType font feature with a name and value.
// Font features control typographic features like ligatures, small caps, etc.
//
// Ported from: skia-source/modules/skparagraph/include/TextStyle.h
type FontFeature struct {
	// Name is the OpenType feature tag (4-character string, e.g., "liga", "smcp").
	Name string

	// Value is the feature setting value. Typically 0 = off, 1 = on,
	// but some features support other values.
	Value int
}

// NewFontFeature creates a new FontFeature with the given name and value.
func NewFontFeature(name string, value int) FontFeature {
	return FontFeature{
		Name:  name,
		Value: value,
	}
}

// Equals returns true if this font feature equals another.
func (f FontFeature) Equals(other FontFeature) bool {
	return f.Name == other.Name && f.Value == other.Value
}
