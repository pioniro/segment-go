package ordered

import (
	"fmt"
	"github.com/pioniro/segment-go"
)

// ordered is a types, that can be ordered with built-in operators <, >, <=, >=
type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// OrderedSegment is an implementation of ISegment, TryToSegment interfaces for ordered values .
type OrderedSegment[T ordered] struct {
	from segment.Border[T]
	till segment.Border[T]
}

func NewOrderedSegment[T ordered](from, till segment.Border[T]) *OrderedSegment[T] {
	return &OrderedSegment[T]{
		from: from,
		till: till,
	}
}
func (s *OrderedSegment[T]) From() *segment.Border[T] {
	return &s.from
}

func (s *OrderedSegment[T]) Till() *segment.Border[T] {
	return &s.till
}

func (s *OrderedSegment[T]) TryTo(from segment.Bound, till segment.Bound) (segment.TryToSegment[T], error) {
	f, err := LeftBoundTo(s.from, from)
	if err != nil {
		return nil, err
	}

	t, err := RightBoundTo(s.till, till)
	if err != nil {
		return nil, err
	}

	return NewOrderedSegment(f, t), nil
}

// String returns a string representation of a segment.
// example: [1;2], (1;2], [1;2), (1;2), [1;inf), (1;inf), (inf;2], (inf;2), (inf;inf)
func (s *OrderedSegment[T]) String() string {
	leftBound := "("
	rightBound := ")"
	from := s.from
	till := s.till
	if s.from.IsIncluded() {
		leftBound = "["
	}
	if s.till.IsIncluded() {
		rightBound = "]"
	}
	return fmt.Sprintf("%s%s;%s%s", leftBound, from.Value(), till.Value(), rightBound)
}

func (s *OrderedSegment[T]) IsEmpty() bool {
	if s.From().IsUnbound() || s.Till().IsUnbound() {
		return false
	}
	inc, err := s.TryTo(segment.Included, segment.Included)
	// this is possible only if rightsegment.Border is segment.Excluded minimum or leftsegment.Border is segment.Excluded maximum
	if err != nil {
		return true
	}
	return inc.From().Value().Value() > inc.Till().Value().Value()
}

func (s *OrderedSegment[T]) IsIncludes(point T) bool {
	if s.From().IsUnbound() && s.Till().IsUnbound() {
		return true
	}
	inc, err := s.TryTo(segment.Included, segment.Included)
	// this is possible only if rightsegment.Border is segment.Excluded minimum or leftsegment.Border is segment.Excluded maximum.
	// In both cases segment does not include anything.
	if err != nil {
		return false
	}
	if s.From().IsUnbound() {
		return inc.Till().Value().Value() >= point
	}
	if s.Till().IsUnbound() {
		return inc.From().Value().Value() <= point
	}
	return inc.From().Value().Value() <= point && inc.Till().Value().Value() >= point
}

func LeftBoundTo[T ordered](b segment.Border[T], to segment.Bound) (segment.Border[T], error) {
	if !b.IsUnbound() && !b.IsBound(to) {
		switch to {
		case segment.Included:
			value, err := b.Value().Next()
			if err != nil {
				return b, err
			}
			b = segment.NewBorder(segment.Included, value)
		case segment.Excluded:
			value, err := b.Value().Prev()
			if err != nil {
				return b, err
			}
			b = segment.NewBorder(segment.Excluded, value)
		case segment.Unbound:
		}
	}
	return b, nil
}

func RightBoundTo[T ordered](b segment.Border[T], to segment.Bound) (segment.Border[T], error) {
	if !b.IsUnbound() && !b.IsBound(to) {
		switch to {
		case segment.Included:
			value, err := b.Value().Prev()
			if err != nil {
				return b, err
			}
			b = segment.NewBorder(segment.Included, value)
		case segment.Excluded:
			value, err := b.Value().Next()
			if err != nil {
				return b, err
			}
			b = segment.NewBorder(segment.Excluded, value)
		case segment.Unbound:
		}
	}
	return b, nil
}
