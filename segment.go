package segment

import (
	"errors"
	"github.com/pioniro/generator-go"
)

var (
	ErrSegmentTooBig = errors.New("segment is too big")
)

type (
	// ISegment is a segment interface, that is used to define borders of a segment
	ISegment[T any] interface {
		From() *Border[T]
		Till() *Border[T]
	}
	// TryToSegment is an extension of ISegment interface, that is used to try to create a
	// new segment from a given segment, but with different borders if it is possible.
	TryToSegment[T any] interface {
		ISegment[T]
		// TryTo tries to create a new segment from a given segment, but with different borders if it is possible.
		// Unbound borders will remain unbound.
		// Left Included can be cast to Excluded with calling s.From().Next(), eg. [1;2) -> (2;2)
		// Left Excluded can be cast to Included with calling s.From().Prev(), eg. (1;2] -> [1;1]
		// Right Included can be cast to Excluded with calling s.Till().Prev(), eg. [1;2) -> [1;1]
		// Right Excluded can be cast to Included with calling s.Till().Next(), eg. (1;2] -> (2;2)
		// If it is not possible to cast a border, then an error will be returned.
		TryTo(from, till Bound) (TryToSegment[T], error)
		IsEmpty() bool
	}
	// IsEmptySegment is an extension of ISegment interface, that is used to check if a segment is empty.
	IsEmptySegment[T any] interface {
		ISegment[T]
		// IsEmpty returns true if a segment is empty.
		IsEmpty() bool
	}
	SizeSegment[T any] interface {
		// Size returns a size of a segment.
		// Some segments can have infinite size, so it is not possible to
		// Some T cant calculate size, so it is not possible to calculate size for them.
		// In this case an error will be returned.
		Size() (T, error)
	}
	SplitSegment[T any] interface {
		ISegment[T]
		// Split splits a segment into chunks of a given size.
		// In some cases we cannot know the number of chunks in advance, so we can't use a slice,
		// To do this we use a generator for a potentially infinite number of pieces
		Split(size T) generator.Generator[SplitSegment[T]]
	}
	IncludedSegment[T any] interface {
		ISegment[T]
		// IsIncludes returns true if a segment includes a point.
		IsIncludes(point T) bool
	}
	IterableSegment[T any] interface {
		ISegment[T]
		// Iterate returns a generator that iterates over values of a segment.
		Iterate() generator.Generator[T]
	}
)

// Segment is a generic implementation of ISegment interface.
type Segment[T any] struct {
	F Border[T]
	T Border[T]
}

func (s *Segment[T]) From() *Border[T] {
	return &s.F
}

func (s *Segment[T]) Till() *Border[T] {
	return &s.T
}

func NewSegment[T any](from, till Border[T]) *Segment[T] {
	return &Segment[T]{
		F: from,
		T: till,
	}
}
