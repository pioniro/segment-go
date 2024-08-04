package segment_int

import (
	gen "github.com/pioniro/generator-go"
	rng "github.com/pioniro/segment-go"
	"github.com/pioniro/segment-go/ordered"
)

type intLike interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}
type IntSegment[T intLike] struct {
	*ordered.OrderedSegment[T]
}

func NewIntSegment[T intLike](from, till rng.Border[T]) *IntSegment[T] {
	return &IntSegment[T]{
		OrderedSegment: ordered.NewOrderedSegment(from, till),
	}
}

// TryTo tries to create a new segment from a given segment, but with different borders if it is possible.
//
// If from value is Excluded(Max) and we want to cast it to Included, then we must to use Max.Next().
// But if we use Max.Next() we get an error, because Max.Next() is overflow. So we just return an error ErrHasNoNextValue.
//
// If till value is Included(Min) and we want to cast it to Excluded, then we must to use Min.Prev().
// But if we use Min.Prev() we get an error, because Min.Prev() is overflow. So we just return an error ErrHasNoPrevValue.
func (s *IntSegment[T]) TryTo(from rng.Bound, till rng.Bound) (rng.TryToSegment[T], error) {
	f, err := ordered.LeftBoundTo(*s.From(), from)
	if err != nil {
		return nil, err
	}

	t, err := ordered.RightBoundTo(*s.Till(), till)
	if err != nil {
		return nil, err
	}

	return NewIntSegment(f, t), nil
}

// Size returns a size of a segment.
// It can be overflow if size is bigger than max(T)
// for example,
//
//	r := NewIntSegment[Int8](NewIncluded(Int(127)), NewIncluded(Int(-1)))
//	r.Size() // -128
//
// In this case(intLike types) we can always use uint64, but go's generics doesn't support it yet, so we work with what we have
func (s *IntSegment[T]) Size() (T, error) {
	var size T
	var zero T
	// we can always cast a segment to a segment with included borders
	// corner cases:
	//	(A; MinInt) -> [A+1; MinInt-1], but size is less than 1, so we return empty gen
	//	(MaxInt; A) -> [MaxInt+1; A-1], but size is less than 1, so we return empty gen
	seg, err := mustToIncluded(s)
	if err != nil {
		return zero, nil
	}
	// (2; 4] == [3; 4] == [3; 5) == (2; 5)
	// Size( (2; 4] ) == 4 - 2 		= 2
	// Size( [3; 4] ) == 4 - 3 + 1	= 2
	// Size( [3; 5) ) == 5 - 3		= 2
	// Size( (2; 5) ) == 5 - 2 - 1	= 2
	// So we have to add 1 to difference, because seg is always [A; B].
	f := seg.From().Value().Value()
	t := seg.Till().Value().Value()
	one := zero + 1
	if t < f {
		return zero, nil
	}
	if t == f {
		return one, nil
	}
	// invariant: t >= f
	// overflow is impossible
	size = t - f
	// overflow protection
	// size is positive (t >= f as it checked below)
	// if size < 1, then we have overflow:
	// 255 + 1 = 0, 0 < 1
	// 127 + 1 = -128, -128 < 1
	// ...
	size += one
	if size < one {
		return zero, rng.ErrSegmentTooBig
	}
	return size, nil
}

// Split splits a segment into chunks of a given size.
// Bounds may change, but size will be the same: [A; B] -> [A; A+size), [A+size; A+size*2), ... [A+size*N; B] or [A + size*N; B) if B - A = size*N.
// If size less than 1, then empty gen will be returned
// If from is Unbound, then we must to start from [min(T); min(T) + size), then [min(T) + size; min(T) + size*2), etc.
// If till is Unbound, then we must to end at [max(T) - size; max(T)] or [max(T); max(T)] (this is the only case where the right border will be Included)
func (s *IntSegment[T]) Split(size T) gen.Generator[rng.SplitSegment[T]] {
	var zero T
	if size <= zero {
		return func(yield gen.Yield[rng.SplitSegment[T]]) {}
	}
	// Since we can always cast a segment to a segment with included borders, and this greatly simplifies the algorithm (compares of values),
	// we do this
	//
	// corner cases:
	//	(A; MinInt) -> [A+1; MinInt-1], but size is less than 1, so we return empty gen
	//	(MaxInt; A) -> [MaxInt+1; A-1], but size is less than 1, so we return empty gen
	inc, err := mustToIncluded(s)
	if err != nil {
		return func(yield gen.Yield[rng.SplitSegment[T]]) {}
	}
	start := inc.From().Value().Value()
	finish := inc.Till().Value().Value()
	return func(yield gen.Yield[rng.SplitSegment[T]]) {
		l := start
		// overflow protection
		// size is positive (size >= 1)
		// so, r = l + size >= l. If not - overflow
		r := l + size
		if r <= l {
			r = maxInt[T]()
		}
		for l <= finish {
			// d is a current size of a chunk
			// invariant:
			// 1. d <= size
			// 2. d > 0
			// 3. if d less than size, then we throw [..;..] range and get out of the loop
			var d T
			// SIGNED: int8
			// [-128; 5]
			//  left = 5 - (-128) = -123 // overflow, left < 0, so d = size
			// [1; 1]
			//  left = 1 - 1 = 0 // no overflow, left >= 0, so d = min(left, size) = 0
			// [-127; 5]
			//   left = 5 - (-127) = -124 // overflow, left < 0, so d = size
			// [-128; -127]
			//   left = -127 - (-128) = 1 // no overflow, left > 0, so d = min(left, size)
			// UNSIGNED: uint8
			// [0; 5]
			//   left = 5 - 0 = 5 // no overflow, left > 0, so d = min(left, size)
			// [1; 5]
			//   left = 5 - 1 = 4 // no overflow, left > 0, so d = min(left, size)
			// [0; 255]
			//   left = 255 - 0 = 255 // no overflow, left > 0, so d = min(left, size)
			// so, checking for overflow is not necessary, we can rely to left < zero
			left := finish - l
			if left < zero {
				d = size
			} else {
				d = min(size, left)
			}
			// invariant 1: d <= size
			// invariant 2: d > 0
			// overflow is impossible: d < max - l, d > 0
			r = l + d
			// or d < size (invariant 3)
			if d != size {
				yield(NewIntSegment(rng.NewIncluded(Int(l)), rng.NewIncluded(Int(r))), nil)
				return
			}
			if !yield(NewIntSegment(rng.NewIncluded(Int(l)), rng.NewExcluded(Int(r))), nil) {
				return
			}
			l = r
		}
	}
}

func mustToIncluded[T intLike](s *IntSegment[T]) (*IntSegment[T], error) {
	inc, err := s.TryTo(rng.Included, rng.Included)
	if err != nil {
		return nil, err
	}
	// We know, that Border type is always Included or Unbound, so we can use Value() method
	var start T
	var finish T
	// calculate start and finish values with taking into account bounds
	if inc.From().IsUnbound() {
		start = minInt[T]()
	} else {
		start = inc.From().Value().Value()
	}
	if inc.Till().IsUnbound() {
		finish = maxInt[T]()
	} else {
		finish = inc.Till().Value().Value()
	}
	return NewIntSegment(rng.NewIncluded(Int(start)), rng.NewIncluded(Int(finish))), nil
}

func (s *IntSegment[T]) IsEmpty() bool {
	if s.From().IsUnbound() || s.Till().IsUnbound() {
		return false
	}
	size, err := s.Size()
	if err != nil {
		return false
	}
	var zero T
	return size == zero
}

func (s *IntSegment[T]) Iterate() gen.Generator[T] {
	return func(yield gen.Yield[T]) {
		// Since we can always cast an int segment to a segment with included borders, and this greatly simplifies the algorithm (compares of values),
		// we do this
		//
		// corner cases:
		//	(A; MinInt) -> [A+1; MinInt-1], but size is less than 1, so we return empty gen
		//	(MaxInt; A) -> [MaxInt+1; A-1], but size is less than 1, so we return empty gen
		inc, err := mustToIncluded(s)
		if err != nil {
			return
		}
		start := inc.From().Value().Value()
		finish := inc.Till().Value().Value()
		for i := start; i <= finish; i++ {
			if !yield(i, nil) {
				return
			}
		}
	}
}
