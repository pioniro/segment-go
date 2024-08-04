package segment_int

import (
	"github.com/pioniro/segment-go"
	"strconv"
)

type intValue[T intLike] struct {
	value T
}

func Int[T intLike](v T) segment.Value[T] {
	return &intValue[T]{
		value: v,
	}
}

func (v *intValue[T]) Value() T {
	return v.value
}

func (v *intValue[T]) String() string {
	typ, _ := intType(v.value)
	switch typ {
	case "uint64":
		fallthrough
	case "uint32":
		return strconv.FormatUint(uint64(v.value), 10)
	default:
		return strconv.FormatInt(int64(v.value), 10)
	}
}

// Next returns a next value of a given value. Or an error if it is not possible to calculate a next value (overflow for example).
func (v *intValue[T]) Next() (segment.Value[T], error) {
	newValue := v.value + 1
	if newValue < v.value {
		return Int(v.value), segment.ErrHasNoNextValue
	}
	return Int(newValue), nil
}

// Prev returns a prev value of a given value. Or an error if it is not possible to calculate a prev value (overflow for example).
func (v *intValue[T]) Prev() (segment.Value[T], error) {
	newValue := v.value - 1
	if newValue > v.value {
		return Int(v.value), segment.ErrHasNoPrevValue
	}
	return Int(newValue), nil
}
