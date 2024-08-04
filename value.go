package segment

import (
	"errors"
	"fmt"
)

var (
	ErrHasNoNextValue = errors.New("has no next value")
	ErrHasNoPrevValue = errors.New("has no prev value")
)

type Value[T any] interface {
	fmt.Stringer

	Next() (Value[T], error)
	Prev() (Value[T], error)
	Value() T
}

type infValue[T any] struct {
}

// Inf returns a value, that is bigger than any other value of a given type (modulo; except inf itself)
// Distance between any value (even Inf) and Inf is Inf.
// Inf is not equal to Inf (our Inf has not a sign, so we can't compare them).
// As Inf does not have a sign, so we cant make a segment (-Inf; -Inf) or (+Inf; +Inf), only (-Inf; +Inf) which the same as (Inf; Inf).
func Inf[T any]() Value[T] {
	return &infValue[T]{}
}

func (i infValue[T]) String() string {
	return "inf"
}

// Next returns an error, because it is not possible to calculate a next value for Inf.
func (i infValue[T]) Next() (Value[T], error) {
	return Inf[T](), ErrHasNoNextValue
}

// Prev returns an error, because it is not possible to calculate a prev value for Inf.
func (i infValue[T]) Prev() (Value[T], error) {
	return Inf[T](), ErrHasNoPrevValue
}

func (i infValue[T]) Value() T {
	return *new(T)
}
