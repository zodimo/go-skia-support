package interfaces

// SkData holds an immutable data buffer.
//
// Ported from: skia-source/include/core/SkData.h
type SkData interface {
	// Size returns the number of bytes stored.
	Size() int

	// Bytes returns the data as a byte slice.
	Bytes() []byte

	// Equals returns true if the two data objects have the same length,
	// and the contents are equal.
	Equals(other SkData) bool
}
