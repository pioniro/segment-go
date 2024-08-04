package segment

import (
	"fmt"
	"reflect"
	"testing"
)

type customtype struct {
	a int
	b string
}

type customvalue[T customtype] struct {
	t customtype
}

func (c *customvalue[T]) String() string {
	return fmt.Sprintf("{%d, %s}", c.t.a, c.t.b)
}

func (c *customvalue[T]) Next() (Value[T], error) {
	return &customvalue[T]{t: customtype{
		a: c.t.a + 1,
		b: c.t.b,
	}}, nil
}

func (c *customvalue[T]) Prev() (Value[T], error) {
	return &customvalue[T]{t: customtype{
		a: c.t.a - 1,
		b: c.t.b,
	}}, nil
}

func (c *customvalue[T]) Value() T {
	return T(c.t)
}
func newcustomvalue(a int, b string) Value[customtype] {
	return &customvalue[customtype]{t: customtype{
		a: a,
		b: b,
	}}
}

func TestNewSegment(t *testing.T) {
	type args[T any] struct {
		from Border[T]
		till Border[T]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want *Segment[T]
	}
	tests := []testCase[customtype]{
		{
			name: "[{1,a}; {2,b}]",
			args: args[customtype]{
				from: NewIncluded(newcustomvalue(1, "a")),
				till: NewIncluded(newcustomvalue(2, "b")),
			},
			want: &Segment[customtype]{
				F: NewIncluded(newcustomvalue(1, "a")),
				T: NewIncluded(newcustomvalue(2, "b")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSegment(tt.args.from, tt.args.till); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSegment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSegment_From(t *testing.T) {
	type testCase[T any] struct {
		name string
		s    Segment[T]
		want Border[T]
	}
	tests := []testCase[customtype]{
		{
			name: "[{1,a}; {2,b}]",
			s: Segment[customtype]{
				F: NewIncluded(newcustomvalue(1, "a")),
				T: NewIncluded(newcustomvalue(2, "b")),
			},
			want: NewIncluded(newcustomvalue(1, "a")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.From(); !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("From() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSegment_Till(t *testing.T) {
	type testCase[T any] struct {
		name string
		s    Segment[T]
		want Border[T]
	}
	tests := []testCase[customtype]{
		{
			name: "[{1,a}; {2,b}]",
			s: Segment[customtype]{
				F: NewIncluded(newcustomvalue(1, "a")),
				T: NewIncluded(newcustomvalue(2, "b")),
			},
			want: NewIncluded(newcustomvalue(2, "b")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Till(); !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("Till() = %v, want %v", got, tt.want)
			}
		})
	}
}
