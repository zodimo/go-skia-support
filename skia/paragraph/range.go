package paragraph

import (
	"math"
)

// EmptyIndex represents an invalid/empty index value.
// This matches C++ std::numeric_limits<size_t>::max().
const EmptyIndex = math.MaxInt

// Range represents a range of values with start and end.
// This is a generic type matching C++ SkRange<T>.
//
// Ported from: skia-source/modules/skparagraph/include/DartTypes.h
type Range[T ~int | ~int32 | ~int64 | ~uint | ~uint32 | ~uint64] struct {
	Start T
	End   T
}

// NewRange creates a new Range with the given start and end values.
func NewRange[T ~int | ~int32 | ~int64 | ~uint | ~uint32 | ~uint64](start, end T) Range[T] {
	return Range[T]{Start: start, End: end}
}

// Width returns the width of the range (end - start).
func (r Range[T]) Width() T {
	return r.End - r.Start
}

// Shift moves the range by the given delta.
func (r *Range[T]) Shift(delta T) {
	r.Start += delta
	r.End += delta
}

// Contains returns true if this range fully contains the other range.
func (r Range[T]) Contains(other Range[T]) bool {
	return r.Start <= other.Start && r.End >= other.End
}

// Intersects returns true if this range overlaps with the other range.
func (r Range[T]) Intersects(other Range[T]) bool {
	return max(r.Start, other.Start) <= min(r.End, other.End)
}

// Intersection returns the intersection of this range with another.
// If ranges don't intersect, returns an empty or inverted range.
func (r Range[T]) Intersection(other Range[T]) Range[T] {
	return Range[T]{
		Start: max(r.Start, other.Start),
		End:   min(r.End, other.End),
	}
}

// Empty returns true if this range is empty (represents no valid range).
// A range is empty when both start and end equal EmptyIndex.
func (r Range[T]) Empty() bool {
	return int64(r.Start) == int64(EmptyIndex) && int64(r.End) == int64(EmptyIndex)
}

// IsValid returns true if this is a valid non-negative-width range.
func (r Range[T]) IsValid() bool {
	return r.Start <= r.End
}

// Equals returns true if this range equals another range.
func (r Range[T]) Equals(other Range[T]) bool {
	return r.Start == other.Start && r.End == other.End
}

// TextRange is a range of text indices (code unit positions).
type TextRange = Range[int]

// BlockRange represents a range of blocks.
type BlockRange = Range[int]

// EmptyRange represents an empty/invalid range.
var EmptyRange = Range[int]{Start: EmptyIndex, End: EmptyIndex}

// NewTextRange creates a new TextRange with the given start and end.
func NewTextRange(start, end int) TextRange {
	return Range[int]{Start: start, End: end}
}

// NewBlockRange creates a new BlockRange with the given start and end.
func NewBlockRange(start, end int) BlockRange {
	return Range[int]{Start: start, End: end}
}

// NewClusterRange creates a new ClusterRange with the given start and end.
// ClusterRange type is defined in run.go.
func NewClusterRange(start, end int) Range[int] {
	return Range[int]{Start: start, End: end}
}

// helper functions for generic min/max
func min[T ~int | ~int32 | ~int64 | ~uint | ~uint32 | ~uint64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func max[T ~int | ~int32 | ~int64 | ~uint | ~uint32 | ~uint64](a, b T) T {
	if a > b {
		return a
	}
	return b
}
