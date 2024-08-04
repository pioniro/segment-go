package segment

type Bound int

const (
	Unbound Bound = iota
	Included
	Excluded
)

// Border is a border of a segment. It can be Included, Excluded or Unbound.
type Border[T any] struct {
	bound Bound
	value Value[T]
}

// NewBorder creates a new border.
func NewBorder[T any](bound Bound, value Value[T]) Border[T] {
	if bound == Unbound {
		return NewUnbound[T]()
	}
	return Border[T]{
		bound: bound,
		value: value,
	}
}

func (b *Border[T]) IsUnbound() bool {
	return b.bound == Unbound
}

func (b *Border[T]) IsIncluded() bool {
	return b.bound == Included
}

func (b *Border[T]) IsExcluded() bool {
	return b.bound == Excluded
}

func (b *Border[T]) IsBound(bound Bound) bool {
	return b.bound == bound
}

func (b *Border[T]) Value() Value[T] {
	return b.value
}

func (b *Border[T]) String() string {
	if b.IsUnbound() {
		return "inf"
	}
	return b.value.String()
}

func NewIncluded[T any](value Value[T]) Border[T] {
	return NewBorder(Included, value)
}
func NewExcluded[T any](value Value[T]) Border[T] {
	return NewBorder(Excluded, value)
}

func NewUnbound[T any]() Border[T] {
	return Border[T]{
		bound: Unbound,
		value: Inf[T](),
	}
}
