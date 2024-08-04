package segment

import (
	"reflect"
	"strconv"
	"testing"
)

type testValue struct {
	value int64
}

func NewTestValue(v int64) Value[int64] {
	return &testValue{
		value: v,
	}
}

func (v *testValue) Value() int64 {
	return v.value
}

func (v *testValue) String() string {
	return strconv.FormatInt(v.value, 10)
}
func (v *testValue) Next() (Value[int64], error) {
	return NewTestValue(v.value + 1), nil
}
func (v *testValue) Prev() (Value[int64], error) {
	return NewTestValue(v.value - 1), nil
}

func TestBorder_IsBound(t *testing.T) {
	type args struct {
		bound Bound
	}
	type testCase[T any] struct {
		name string
		b    Border[T]
		args args
		want bool
	}
	tests := []testCase[int64]{
		{
			name: "[1] is included",
			b:    NewBorder(Included, NewTestValue(1)),
			args: args{
				bound: Included,
			},
			want: true,
		},
		{
			name: "(1) is included",
			b:    NewBorder(Excluded, NewTestValue(1)),
			args: args{
				bound: Included,
			},
			want: false,
		},
		{
			name: "(inf) is included",
			b:    NewBorder(Unbound, NewTestValue(1)),
			args: args{
				bound: Included,
			},
			want: false,
		},
		{
			name: "[1] is excluded",
			b:    NewBorder(Included, NewTestValue(1)),
			args: args{
				bound: Excluded,
			},
			want: false,
		},
		{
			name: "(1) is excluded",
			b:    NewBorder(Excluded, NewTestValue(1)),
			args: args{
				bound: Excluded,
			},
			want: true,
		},
		{
			name: "(inf) is included",
			b:    NewBorder(Unbound, NewTestValue(1)),
			args: args{
				bound: Excluded,
			},
			want: false,
		},
		{
			name: "(inf) is unbound",
			b:    NewBorder(Unbound, NewTestValue(1)),
			args: args{
				bound: Unbound,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IsBound(tt.args.bound); got != tt.want {
				t.Errorf("IsBound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBorder_IsExcluded(t *testing.T) {
	type testCase[T any] struct {
		name string
		b    Border[T]
		want bool
	}
	tests := []testCase[int64]{
		{
			name: "[1]",
			b:    NewBorder(Included, NewTestValue(1)),
			want: false,
		},
		{
			name: "(1)",
			b:    NewBorder(Excluded, NewTestValue(1)),
			want: true,
		},
		{
			name: "(inf)",
			b:    NewUnbound[int64](),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IsExcluded(); got != tt.want {
				t.Errorf("IsExcluded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBorder_IsIncluded(t *testing.T) {
	type testCase[T any] struct {
		name string
		b    Border[T]
		want bool
	}
	tests := []testCase[int64]{
		{
			name: "[1]",
			b:    NewBorder(Included, NewTestValue(1)),
			want: true,
		},
		{
			name: "(1)",
			b:    NewBorder(Excluded, NewTestValue(1)),
			want: false,
		},
		{
			name: "(inf)",
			b:    NewUnbound[int64](),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IsIncluded(); got != tt.want {
				t.Errorf("IsIncluded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBorder_IsUnbound(t *testing.T) {
	type testCase[T any] struct {
		name string
		b    Border[T]
		want bool
	}
	tests := []testCase[int64]{
		{
			name: "[1]",
			b:    NewBorder(Included, NewTestValue(1)),
			want: false,
		},
		{
			name: "(1)",
			b:    NewBorder(Excluded, NewTestValue(1)),
			want: false,
		},
		{
			name: "(inf)",
			b:    NewUnbound[int64](),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IsUnbound(); got != tt.want {
				t.Errorf("IsUnbound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBorder_String(t *testing.T) {
	type testCase[T any] struct {
		name string
		b    Border[T]
		want string
	}
	tests := []testCase[int64]{
		{
			name: "[1]",
			b:    NewBorder(Included, NewTestValue(1)),
			want: "1",
		},
		{
			name: "(1)",
			b:    NewBorder(Excluded, NewTestValue(1)),
			want: "1",
		},
		{
			name: "(inf)",
			b:    NewUnbound[int64](),
			want: "inf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBorder_Value(t *testing.T) {
	type testCase[T any] struct {
		name string
		b    Border[T]
		want Value[T]
	}
	tests := []testCase[int64]{
		{
			name: "[1]",
			b:    NewBorder(Included, NewTestValue(1)),
			want: NewTestValue(1),
		},
		{
			name: "(1)",
			b:    NewBorder(Excluded, NewTestValue(1)),
			want: NewTestValue(1),
		},
		{
			name: "(inf)",
			b:    NewUnbound[int64](),
			want: Inf[int64](),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBorder(t *testing.T) {
	type args[T any] struct {
		bound Bound
		value Value[T]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want Border[T]
	}
	tests := []testCase[int64]{
		{
			name: "[1]",
			args: args[int64]{
				bound: Included,
				value: NewTestValue(1),
			},
			want: Border[int64]{
				bound: Included,
				value: NewTestValue(int64(1)),
			},
		},
		{
			name: "(1)",
			args: args[int64]{
				bound: Excluded,
				value: NewTestValue(1),
			},
			want: Border[int64]{
				bound: Excluded,
				value: NewTestValue(int64(1)),
			},
		},
		{
			name: "(inf)",
			args: args[int64]{
				bound: Unbound,
				value: NewTestValue(1),
			},
			want: Border[int64]{
				bound: Unbound,
				value: Inf[int64](),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBorder(tt.args.bound, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBorder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewExcluded(t *testing.T) {
	type args[T any] struct {
		value Value[T]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want Border[T]
	}
	tests := []testCase[int64]{
		{
			name: "(1)",
			args: args[int64]{
				value: NewTestValue(1),
			},
			want: Border[int64]{
				bound: Excluded,
				value: NewTestValue(int64(1)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewExcluded(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewExcluded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIncluded(t *testing.T) {
	type args[T any] struct {
		value Value[T]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want Border[T]
	}
	tests := []testCase[int64]{
		{
			name: "[1]",
			args: args[int64]{
				value: NewTestValue(1),
			},
			want: Border[int64]{
				bound: Included,
				value: NewTestValue(int64(1)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIncluded(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIncluded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUnbound(t *testing.T) {
	type testCase[T any] struct {
		name string
		want Border[T]
	}
	tests := []testCase[int64]{
		{
			name: "(inf)",
			want: Border[int64]{
				bound: Unbound,
				value: Inf[int64](),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUnbound[int64](); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUnbound() = %v, want %v", got, tt.want)
			}
		})
	}
}
