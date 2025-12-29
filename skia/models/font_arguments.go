package models

// FontArguments represents a set of arguments to filter a font from a stream or file.
// Ported from: SkFontArguments.h
type FontArguments struct {
	CollectionIndex         int
	VariationDesignPosition VariationPosition
	// Palette ignored for now
}

type VariationPosition struct {
	Coordinates []VariationCoordinate
}

type VariationCoordinate struct {
	Axis  uint32
	Value float32
}
