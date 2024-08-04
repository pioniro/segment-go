package ordered

import (
	"github.com/pioniro/segment-go"
	"math"
	"reflect"
	"strconv"
	"testing"
)

type testValue struct {
	value int64
}

func NewTestValue(v int64) segment.Value[int64] {
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
func (v *testValue) Next() (segment.Value[int64], error) {
	if v.value == math.MaxInt64 {
		return nil, segment.ErrHasNoNextValue
	}
	return NewTestValue(v.value + 1), nil
}
func (v *testValue) Prev() (segment.Value[int64], error) {
	if v.value == math.MinInt64 {
		return nil, segment.ErrHasNoPrevValue
	}
	return NewTestValue(v.value - 1), nil
}

func TestNewOrderedSegment(t *testing.T) {
	type args[T int64] struct {
		from segment.Border[T]
		till segment.Border[T]
	}
	type testCase[T int64] struct {
		name string
		args args[T]
		want *OrderedSegment[T]
	}
	tests := []testCase[int64]{
		{
			name: "[1;1]",
			args: args[int64]{
				from: segment.NewBorder(segment.Included, NewTestValue(1)),
				till: segment.NewBorder(segment.Included, NewTestValue(1)),
			},
			want: &OrderedSegment[int64]{from: segment.NewBorder(segment.Included, NewTestValue(1)), till: segment.NewBorder(segment.Included, NewTestValue(1))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOrderedSegment(tt.args.from, tt.args.till); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrderedSegment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderedSegment_String(t *testing.T) {
	type testCase[T int64] struct {
		name string
		s    *OrderedSegment[T]
		want string
	}
	tests := []testCase[int64]{
		{
			name: "[1;100]",
			s:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewBorder(segment.Included, NewTestValue(100))),
			want: "[1;100]",
		},
		{
			name: "(1;100]",
			s:    NewOrderedSegment(segment.NewBorder(segment.Excluded, NewTestValue(1)), segment.NewBorder(segment.Included, NewTestValue(100))),
			want: "(1;100]",
		},
		{
			name: "(1;100)",
			s:    NewOrderedSegment(segment.NewBorder(segment.Excluded, NewTestValue(1)), segment.NewBorder(segment.Excluded, NewTestValue(100))),
			want: "(1;100)",
		},
		{
			name: "[1;100)",
			s:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewBorder(segment.Excluded, NewTestValue(100))),
			want: "[1;100)",
		},
		{
			name: "[1;inf)",
			s:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewUnbound[int64]()),
			want: "[1;inf)",
		},
		{
			name: "(inf;100)",
			s:    NewOrderedSegment(segment.NewUnbound[int64](), segment.NewBorder(segment.Excluded, NewTestValue(100))),
			want: "(inf;100)",
		},
		{
			name: "(inf;inf)",
			s:    NewOrderedSegment(segment.NewUnbound[int64](), segment.NewUnbound[int64]()),
			want: "(inf;inf)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderedSegment_TryTo(t *testing.T) {
	type args struct {
		from segment.Bound
		till segment.Bound
	}
	type testCase[T int64] struct {
		name    string
		s       *OrderedSegment[T]
		args    args
		want    *OrderedSegment[T]
		wantErr bool
	}
	tests := []testCase[int64]{
		{
			name: "[1;100]->(0;100]",
			s:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewBorder(segment.Included, NewTestValue(100))),
			args: args{
				from: segment.Excluded,
				till: segment.Included,
			},
			want:    NewOrderedSegment(segment.NewBorder(segment.Excluded, NewTestValue(0)), segment.NewBorder(segment.Included, NewTestValue(100))),
			wantErr: false,
		},
		{
			name: "(0;100)->[1;99]",
			s:    NewOrderedSegment(segment.NewBorder(segment.Excluded, NewTestValue(0)), segment.NewBorder(segment.Excluded, NewTestValue(100))),
			args: args{
				from: segment.Included,
				till: segment.Included,
			},
			want:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewBorder(segment.Included, NewTestValue(99))),
			wantErr: false,
		},
		{
			name: "(0;inf)->[1;inf)",
			s:    NewOrderedSegment(segment.NewBorder(segment.Excluded, NewTestValue(0)), segment.NewUnbound[int64]()),
			args: args{
				from: segment.Included,
				till: segment.Included,
			},
			want:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewUnbound[int64]()),
			wantErr: false,
		},
		{
			name: "(inf;inf)->(inf;inf)",
			s:    NewOrderedSegment(segment.NewUnbound[int64](), segment.NewUnbound[int64]()),
			args: args{
				from: segment.Included,
				till: segment.Included,
			},
			want:    NewOrderedSegment(segment.NewUnbound[int64](), segment.NewUnbound[int64]()),
			wantErr: false,
		},
		{
			name: "[1;100]->[1;101)",
			s:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewBorder(segment.Included, NewTestValue(100))),
			args: args{
				from: segment.Included,
				till: segment.Excluded,
			},
			want:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewBorder(segment.Excluded, NewTestValue(101))),
			wantErr: false,
		},
		{
			name: "[1;max]->[1;max+1)",
			s:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewBorder(segment.Included, NewTestValue(math.MaxInt64))),
			args: args{
				from: segment.Included,
				till: segment.Excluded,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "[min;2]->(min-1;2]",
			s:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(math.MinInt64)), segment.NewBorder(segment.Included, NewTestValue(2))),
			args: args{
				from: segment.Excluded,
				till: segment.Included,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "[min;max]->(min-1;max+1)",
			s:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(math.MinInt64)), segment.NewBorder(segment.Included, NewTestValue(math.MaxInt64))),
			args: args{
				from: segment.Excluded,
				till: segment.Excluded,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "[1;2]->[1;2]",
			s:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewBorder(segment.Included, NewTestValue(2))),
			args: args{
				from: segment.Unbound,
				till: segment.Unbound,
			},
			want:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewBorder(segment.Included, NewTestValue(2))),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.TryTo(tt.args.from, tt.args.till)
			if (err != nil) != tt.wantErr {
				t.Errorf("TryTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) && got != nil && tt.want != nil {
				t.Errorf("TryTo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderedSegment_From(t *testing.T) {
	type testCase[T ordered] struct {
		name string
		s    *OrderedSegment[T]
		want segment.Border[T]
	}
	tests := []testCase[int64]{
		{
			name: "[1;2]",
			s:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewBorder(segment.Included, NewTestValue(2))),
			want: segment.NewBorder(segment.Included, NewTestValue(1)),
		},
		{
			name: "(1;2]",
			s:    NewOrderedSegment(segment.NewBorder(segment.Excluded, NewTestValue(1)), segment.NewBorder(segment.Included, NewTestValue(2))),
			want: segment.NewBorder(segment.Excluded, NewTestValue(1)),
		},
		{
			name: "(inf;2]",
			s:    NewOrderedSegment(segment.NewBorder(segment.Unbound, NewTestValue(1)), segment.NewBorder(segment.Included, NewTestValue(2))),
			want: segment.NewBorder(segment.Unbound, segment.Inf[int64]()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.From(); !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("From() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestOrderedSegment_Till(t *testing.T) {
	type testCase[T ordered] struct {
		name string
		s    *OrderedSegment[T]
		want segment.Border[T]
	}
	tests := []testCase[int64]{
		{
			name: "[1;2]",
			s:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewBorder(segment.Included, NewTestValue(2))),
			want: segment.NewBorder(segment.Included, NewTestValue(2)),
		},
		{
			name: "[1;2)",
			s:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewBorder(segment.Excluded, NewTestValue(2))),
			want: segment.NewBorder(segment.Excluded, NewTestValue(2)),
		},
		{
			name: "[1;inf)",
			s:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewBorder(segment.Unbound, NewTestValue(2))),
			want: segment.NewBorder(segment.Unbound, segment.Inf[int64]()),
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

func TestOrderedSegment_IsEmpty(t *testing.T) {
	type testCase[T ordered] struct {
		name string
		s    *OrderedSegment[T]
		want bool
	}
	tests := []testCase[int64]{
		{
			name: "[1;2]",
			s:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewBorder(segment.Included, NewTestValue(2))),
			want: false,
		},
		{
			name: "[1;2)",
			s:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewBorder(segment.Excluded, NewTestValue(2))),
			want: false,
		},
		{
			name: "[1;inf)",
			s:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewBorder(segment.Unbound, NewTestValue(2))),
			want: false,
		},
		{
			name: "(inf;2)",
			s:    NewOrderedSegment(segment.NewBorder(segment.Unbound, NewTestValue(1)), segment.NewBorder(segment.Excluded, NewTestValue(2))),
			want: false,
		},
		{
			name: "(inf;inf)",
			s:    NewOrderedSegment(segment.NewBorder(segment.Unbound, NewTestValue(1)), segment.NewBorder(segment.Unbound, NewTestValue(2))),
			want: false,
		},
		{
			name: "[1;1]",
			s:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewBorder(segment.Included, NewTestValue(2))),
			want: false,
		},
		{
			name: "[2;1]",
			s:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(2)), segment.NewBorder(segment.Included, NewTestValue(1))),
			want: true,
		},
		{
			name: "(2;1]",
			s:    NewOrderedSegment(segment.NewBorder(segment.Excluded, NewTestValue(2)), segment.NewBorder(segment.Included, NewTestValue(1))),
			want: true,
		},
		{
			name: "[2;1)",
			s:    NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(2)), segment.NewBorder(segment.Excluded, NewTestValue(1))),
			want: true,
		},
		{
			name: "(2;1)",
			s:    NewOrderedSegment(segment.NewBorder(segment.Excluded, NewTestValue(1)), segment.NewBorder(segment.Excluded, NewTestValue(2))),
			want: true,
		},
		{
			name: "(1;1)",
			s:    NewOrderedSegment(segment.NewBorder(segment.Excluded, NewTestValue(2)), segment.NewBorder(segment.Excluded, NewTestValue(1))),
			want: true,
		},
		{
			name: "(min;min)",
			s:    NewOrderedSegment(segment.NewBorder(segment.Excluded, NewTestValue(math.MinInt64)), segment.NewBorder(segment.Excluded, NewTestValue(math.MinInt64))),
			want: true,
		},
		{
			name: "(max;max)",
			s:    NewOrderedSegment(segment.NewBorder(segment.Excluded, NewTestValue(math.MaxInt64)), segment.NewBorder(segment.Excluded, NewTestValue(math.MaxInt64))),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderedSegment_IsIncludes(t *testing.T) {
	type testCase[T ordered] struct {
		name  string
		s     OrderedSegment[T]
		cases map[T]bool
	}
	tests := []testCase[int64]{
		{
			name: "[1;2]",
			s:    *NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewBorder(segment.Included, NewTestValue(2))),
			cases: map[int64]bool{
				0: false,
				1: true,
				2: true,
				3: false,
			},
		},
		{
			name: "(1;2]",
			s:    *NewOrderedSegment(segment.NewBorder(segment.Excluded, NewTestValue(1)), segment.NewBorder(segment.Included, NewTestValue(2))),
			cases: map[int64]bool{
				0: false,
				1: false,
				2: true,
				3: false,
			},
		},
		{
			name: "[1;2)",
			s:    *NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewBorder(segment.Excluded, NewTestValue(2))),
			cases: map[int64]bool{
				0: false,
				1: true,
				2: false,
				3: false,
			},
		},
		{
			name: "(1;2)",
			s:    *NewOrderedSegment(segment.NewBorder(segment.Excluded, NewTestValue(1)), segment.NewBorder(segment.Excluded, NewTestValue(2))),
			cases: map[int64]bool{
				0: false,
				1: false,
				2: false,
				3: false,
			},
		},
		{
			name: "(1;3)",
			s:    *NewOrderedSegment(segment.NewBorder(segment.Excluded, NewTestValue(1)), segment.NewBorder(segment.Excluded, NewTestValue(3))),
			cases: map[int64]bool{
				0: false,
				1: false,
				2: true,
				3: false,
				4: false,
			},
		},
		{
			name: "[1;inf)",
			s:    *NewOrderedSegment(segment.NewBorder(segment.Included, NewTestValue(1)), segment.NewUnbound[int64]()),
			cases: map[int64]bool{
				0:             false,
				1:             true,
				2:             true,
				3:             true,
				math.MaxInt64: true,
			},
		},
		{
			name: "(inf;1]",
			s:    *NewOrderedSegment(segment.NewUnbound[int64](), segment.NewBorder(segment.Included, NewTestValue(1))),
			cases: map[int64]bool{
				math.MinInt64: true,
				0:             true,
				1:             true,
				2:             false,
			},
		},
		{
			name: "(inf;inf)",
			s:    *NewOrderedSegment(segment.NewUnbound[int64](), segment.NewUnbound[int64]()),
			cases: map[int64]bool{
				math.MinInt64: true,
				0:             true,
				1:             true,
				math.MaxInt64: true,
			},
		},
		{
			name: "(+max;1]",
			s:    *NewOrderedSegment(segment.NewBorder(segment.Excluded, NewTestValue(math.MaxInt64)), segment.NewBorder(segment.Included, NewTestValue(1))),
			cases: map[int64]bool{
				math.MinInt64: false,
				0:             false,
				1:             false,
				math.MaxInt64: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for point, want := range tt.cases {
				t.Run(strconv.Itoa(int(NewTestValue(point).Value())), func(t *testing.T) {
					if got := tt.s.IsIncludes(point); got != want {
						t.Errorf("IsIncludes() = %v, want %v", got, want)
					}
				})
			}
		})
	}
}
