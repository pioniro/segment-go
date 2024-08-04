package segment_int

import (
	. "github.com/pioniro/segment-go"
	"github.com/pioniro/segment-go/ordered"
	"math"
	"reflect"
	"testing"
)

func TestNewIntSegment(t *testing.T) {
	type args[T intLike] struct {
		from Border[T]
		till Border[T]
	}
	type testCase[T intLike] struct {
		name string
		args args[T]
		want *IntSegment[T]
	}
	tests := []testCase[int64]{
		{
			name: "[1;1]",
			args: args[int64]{
				from: NewBorder(Included, Int[int64](1)),
				till: NewBorder(Included, Int[int64](1)),
			},
			want: &IntSegment[int64]{ordered.NewOrderedSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](1)))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIntSegment(tt.args.from, tt.args.till); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIntSegment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSegment_IsEmpty(t *testing.T) {
	type testCase[T intLike] struct {
		name string
		s    *IntSegment[T]
		want bool
	}
	tests := []testCase[int64]{
		{
			name: "[1;2]",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](2))),
			want: false,
		},
		{
			name: "[1;2)",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Excluded, Int[int64](2))),
			want: false,
		},
		{
			name: "[1;inf)",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Unbound, Int[int64](2))),
			want: false,
		},
		{
			name: "(inf;2)",
			s:    NewIntSegment(NewBorder(Unbound, Int[int64](1)), NewBorder(Excluded, Int[int64](2))),
			want: false,
		},
		{
			name: "(inf;inf)",
			s:    NewIntSegment(NewBorder(Unbound, Int[int64](1)), NewBorder(Unbound, Int[int64](2))),
			want: false,
		},
		{
			name: "[1;1]",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](2))),
			want: false,
		},
		{
			name: "[2;1]",
			s:    NewIntSegment(NewBorder(Included, Int[int64](2)), NewBorder(Included, Int[int64](1))),
			want: true,
		},
		{
			name: "(2;1]",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](2)), NewBorder(Included, Int[int64](1))),
			want: true,
		},
		{
			name: "[2;1)",
			s:    NewIntSegment(NewBorder(Included, Int[int64](2)), NewBorder(Excluded, Int[int64](1))),
			want: true,
		},
		{
			name: "(2;1)",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](1)), NewBorder(Excluded, Int[int64](2))),
			want: true,
		},
		{
			name: "(1;1)",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](2)), NewBorder(Excluded, Int[int64](1))),
			want: true,
		},
		{
			name: "(min;min)",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](math.MinInt64)), NewBorder(Excluded, Int[int64](math.MinInt64))),
			want: true,
		},
		{
			name: "(max;max)",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](math.MaxInt64)), NewBorder(Excluded, Int[int64](math.MaxInt64))),
			want: true,
		},
		{
			name: "(min;max)",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](math.MinInt64)), NewBorder(Excluded, Int[int64](math.MaxInt64))),
			want: false,
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

func TestIntSegment_Size(t *testing.T) {
	type testCase[T intLike] struct {
		name    string
		s       *IntSegment[T]
		want    T
		wantErr bool
	}
	tests := []testCase[int64]{
		{
			name: "[1;2]",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](2))),
			want: 2,
		},
		{
			name: "[1;2)",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Excluded, Int[int64](2))),
			want: 1,
		},
		{
			name: "[1;inf)",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewUnbound[int64]()),
			want: math.MaxInt64,
		},
		{
			name:    "(inf;2)",
			s:       NewIntSegment(NewUnbound[int64](), NewBorder(Excluded, Int[int64](2))),
			want:    0,
			wantErr: true,
		},
		{
			name:    "(inf;inf)",
			s:       NewIntSegment(NewUnbound[int64](), NewUnbound[int64]()),
			want:    0,
			wantErr: true,
		},
		{
			name: "[1;1]",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](1))),
			want: 1,
		},
		{
			name: "[2;1]",
			s:    NewIntSegment(NewBorder(Included, Int[int64](2)), NewBorder(Included, Int[int64](1))),
			want: 0,
		},
		{
			name: "(2;1]",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](2)), NewBorder(Included, Int[int64](1))),
			want: 0,
		},
		{
			name: "[2;1)",
			s:    NewIntSegment(NewBorder(Included, Int[int64](2)), NewBorder(Excluded, Int[int64](1))),
			want: 0,
		},
		{
			name: "(2;1)",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](1)), NewBorder(Excluded, Int[int64](2))),
			want: 0,
		},
		{
			name: "(1;1)",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](2)), NewBorder(Excluded, Int[int64](1))),
			want: 0,
		},
		{
			name: "(min;min)",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](math.MinInt64)), NewBorder(Excluded, Int[int64](math.MinInt64))),
			want: 0,
		},
		{
			name: "(max;max)",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](math.MaxInt64)), NewBorder(Excluded, Int[int64](math.MaxInt64))),
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Size()
			if (err != nil) != tt.wantErr {
				t.Errorf("Size() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Size() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSegment_Size_Uint8(t *testing.T) {
	type testCase[T intLike] struct {
		name    string
		s       *IntSegment[T]
		want    T
		wantErr bool
	}
	tests := []testCase[uint8]{
		{
			name: "[1;2]",
			s:    NewIntSegment(NewBorder(Included, Int[uint8](1)), NewBorder(Included, Int[uint8](2))),
			want: 2,
		},
		{
			name: "[1;2)",
			s:    NewIntSegment(NewBorder(Included, Int[uint8](1)), NewBorder(Excluded, Int[uint8](2))),
			want: 1,
		},
		{
			name: "[1;inf)",
			s:    NewIntSegment(NewBorder(Included, Int[uint8](1)), NewUnbound[uint8]()),
			want: math.MaxUint8,
		},
		{
			name: "(inf;2)",
			s:    NewIntSegment(NewUnbound[uint8](), NewBorder(Excluded, Int[uint8](2))),
			want: 2,
		},
		{
			name:    "(inf;inf)",
			s:       NewIntSegment(NewUnbound[uint8](), NewUnbound[uint8]()),
			want:    0,
			wantErr: true,
		},
		{
			name: "[1;1]",
			s:    NewIntSegment(NewBorder(Included, Int[uint8](1)), NewBorder(Included, Int[uint8](1))),
			want: 1,
		},
		{
			name: "[2;1]",
			s:    NewIntSegment(NewBorder(Included, Int[uint8](2)), NewBorder(Included, Int[uint8](1))),
			want: 0,
		},
		{
			name: "(2;1]",
			s:    NewIntSegment(NewBorder(Excluded, Int[uint8](2)), NewBorder(Included, Int[uint8](1))),
			want: 0,
		},
		{
			name: "[2;1)",
			s:    NewIntSegment(NewBorder(Included, Int[uint8](2)), NewBorder(Excluded, Int[uint8](1))),
			want: 0,
		},
		{
			name: "(2;1)",
			s:    NewIntSegment(NewBorder(Excluded, Int[uint8](1)), NewBorder(Excluded, Int[uint8](2))),
			want: 0,
		},
		{
			name: "(1;1)",
			s:    NewIntSegment(NewBorder(Excluded, Int[uint8](2)), NewBorder(Excluded, Int[uint8](1))),
			want: 0,
		},
		{
			name: "(min;min)",
			s:    NewIntSegment(NewBorder(Excluded, Int[uint8](0)), NewBorder(Excluded, Int[uint8](0))),
			want: 0,
		},
		{
			name: "(max;max)",
			s:    NewIntSegment(NewBorder(Excluded, Int[uint8](math.MaxUint8)), NewBorder(Excluded, Int[uint8](math.MaxUint8))),
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Size()
			if (err != nil) != tt.wantErr {
				t.Errorf("Size() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Size() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSegment_Split(t *testing.T) {
	type args[T intLike] struct {
		size T
	}
	type testCase[T intLike] struct {
		name    string
		s       *IntSegment[T]
		args    args[T]
		want    []SplitSegment[T]
		wantErr bool
	}
	tests := []testCase[int64]{
		{
			name: "[1;2] by 1",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](2))),
			args: args[int64]{size: 1},
			want: []SplitSegment[int64]{
				NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Excluded, Int[int64](2))),
				NewIntSegment(NewBorder(Included, Int[int64](2)), NewBorder(Included, Int[int64](2))),
			},
		},
		{
			name: "[1;2] by 2",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](2))),
			args: args[int64]{size: 2},
			want: []SplitSegment[int64]{
				NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](2))),
			},
		},
		{
			name: "[1;1] by 1",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](1))),
			args: args[int64]{size: 1},
			want: []SplitSegment[int64]{
				NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](1))),
			},
		},
		{
			name: "[1;2) by 1",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](1))),
			args: args[int64]{size: 1},
			want: []SplitSegment[int64]{
				NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](1))),
			},
		},
		{
			name: "(1;2] by 1",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](1)), NewBorder(Included, Int[int64](2))),
			args: args[int64]{size: 1},
			want: []SplitSegment[int64]{
				NewIntSegment(NewBorder(Included, Int[int64](2)), NewBorder(Included, Int[int64](2))),
			},
		},
		{
			name: "(0;2) by 1",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](0)), NewBorder(Included, Int[int64](1))),
			args: args[int64]{size: 1},
			want: []SplitSegment[int64]{
				NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](1))),
			},
		},
		{
			name: "(0;1) by 1",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](0)), NewBorder(Excluded, Int[int64](1))),
			args: args[int64]{size: 1},
			want: nil,
		},
		{
			name: "[1;2] by -1",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](2))),
			args: args[int64]{size: -1},
			want: nil,
		},
		{
			name: "[1;2] by 0",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](2))),
			args: args[int64]{size: 0},
			want: nil,
		},
		{
			name: "[min;max] by max",
			s:    NewIntSegment(NewBorder(Included, Int[int64](math.MinInt64)), NewBorder(Included, Int[int64](math.MaxInt64))),
			args: args[int64]{size: math.MaxInt64},
			want: []SplitSegment[int64]{
				NewIntSegment(NewBorder(Included, Int[int64](math.MinInt64)), NewBorder(Excluded, Int[int64](-1))),
				NewIntSegment(NewBorder(Included, Int[int64](-1)), NewBorder(Excluded, Int[int64](math.MaxInt64-1))),
				NewIntSegment(NewBorder(Included, Int[int64](math.MaxInt64-1)), NewBorder(Included, Int[int64](math.MaxInt64))),
			},
		},
		{
			name: "[inf;inf] by max",
			s:    NewIntSegment(NewUnbound[int64](), NewUnbound[int64]()),
			args: args[int64]{size: math.MaxInt64},
			want: []SplitSegment[int64]{
				NewIntSegment(NewBorder(Included, Int[int64](math.MinInt64)), NewBorder(Excluded, Int[int64](-1))),
				NewIntSegment(NewBorder(Included, Int[int64](-1)), NewBorder(Excluded, Int[int64](math.MaxInt64-1))),
				NewIntSegment(NewBorder(Included, Int[int64](math.MaxInt64-1)), NewBorder(Included, Int[int64](math.MaxInt64))),
			},
		},
		{
			name: "[2;1] by 1",
			s:    NewIntSegment(NewBorder(Included, Int[int64](2)), NewBorder(Included, Int[int64](1))),
			args: args[int64]{size: 1},
			want: nil,
		},
		{
			name: "[min;min] by 1",
			s:    NewIntSegment(NewBorder(Included, Int[int64](math.MinInt64)), NewBorder(Included, Int[int64](math.MinInt64))),
			args: args[int64]{size: 1},
			want: []SplitSegment[int64]{
				NewIntSegment(NewBorder(Included, Int[int64](math.MinInt64)), NewBorder(Included, Int[int64](math.MinInt64))),
			},
		},
		{
			name: "(min;min) by 1",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](math.MinInt64)), NewBorder(Excluded, Int[int64](math.MinInt64))),
			args: args[int64]{size: 1},
			want: nil,
		},
		{
			name: "[max;max] by 1",
			s:    NewIntSegment(NewBorder(Included, Int[int64](math.MaxInt64)), NewBorder(Included, Int[int64](math.MaxInt64))),
			args: args[int64]{size: 1},
			want: []SplitSegment[int64]{
				NewIntSegment(NewBorder(Included, Int[int64](math.MaxInt64)), NewBorder(Included, Int[int64](math.MaxInt64))),
			},
		},
		{
			name: "(max;max) by 1",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](math.MaxInt64)), NewBorder(Excluded, Int[int64](math.MaxInt64))),
			args: args[int64]{size: 1},
			want: nil,
		},
		{
			name: "[min;min+1] by 1",
			s:    NewIntSegment(NewBorder(Included, Int[int64](math.MinInt64)), NewBorder(Included, Int[int64](math.MinInt64+1))),
			args: args[int64]{size: 2},
			want: []SplitSegment[int64]{
				NewIntSegment(NewBorder(Included, Int[int64](math.MinInt64)), NewBorder(Included, Int[int64](math.MinInt64+1))),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Split(tt.args.size).Collect()

			if !reflect.DeepEqual(got, tt.want) && (got != nil || tt.want != nil) {
				t.Errorf("%v .Split() by %v, got = %v, want %v", tt.s, tt.args.size, got, tt.want)
			}
		})
	}
}

func TestIntSegment_Split_Uint8(t *testing.T) {
	type args[T intLike] struct {
		size T
	}
	type testCase[T intLike] struct {
		name    string
		s       *IntSegment[T]
		args    args[T]
		want    []SplitSegment[T]
		wantErr bool
	}
	tests := []testCase[uint8]{
		{
			name: "[1;2] by 1",
			s:    NewIntSegment(NewBorder(Included, Int[uint8](1)), NewBorder(Included, Int[uint8](2))),
			args: args[uint8]{size: 1},
			want: []SplitSegment[uint8]{
				NewIntSegment(NewBorder(Included, Int[uint8](1)), NewBorder(Excluded, Int[uint8](2))),
				NewIntSegment(NewBorder(Included, Int[uint8](2)), NewBorder(Included, Int[uint8](2))),
			},
		},
		{
			name: "[1;2] by 1",
			s:    NewIntSegment(NewBorder(Included, Int[uint8](1)), NewBorder(Included, Int[uint8](2))),
			args: args[uint8]{size: 1},
			want: []SplitSegment[uint8]{
				NewIntSegment(NewBorder(Included, Int[uint8](1)), NewBorder(Excluded, Int[uint8](2))),
				NewIntSegment(NewBorder(Included, Int[uint8](2)), NewBorder(Included, Int[uint8](2))),
			},
		},
		{
			name: "[1;2] by 2",
			s:    NewIntSegment(NewBorder(Included, Int[uint8](1)), NewBorder(Included, Int[uint8](2))),
			args: args[uint8]{size: 2},
			want: []SplitSegment[uint8]{
				NewIntSegment(NewBorder(Included, Int[uint8](1)), NewBorder(Included, Int[uint8](2))),
			},
		},
		{
			name: "[1;1] by 1",
			s:    NewIntSegment(NewBorder(Included, Int[uint8](1)), NewBorder(Included, Int[uint8](1))),
			args: args[uint8]{size: 1},
			want: []SplitSegment[uint8]{
				NewIntSegment(NewBorder(Included, Int[uint8](1)), NewBorder(Included, Int[uint8](1))),
			},
		},
		{
			name: "[1;2) by 1",
			s:    NewIntSegment(NewBorder(Included, Int[uint8](1)), NewBorder(Included, Int[uint8](1))),
			args: args[uint8]{size: 1},
			want: []SplitSegment[uint8]{
				NewIntSegment(NewBorder(Included, Int[uint8](1)), NewBorder(Included, Int[uint8](1))),
			},
		},
		{
			name: "(1;2] by 1",
			s:    NewIntSegment(NewBorder(Excluded, Int[uint8](1)), NewBorder(Included, Int[uint8](2))),
			args: args[uint8]{size: 1},
			want: []SplitSegment[uint8]{
				NewIntSegment(NewBorder(Included, Int[uint8](2)), NewBorder(Included, Int[uint8](2))),
			},
		},
		{
			name: "(0;2) by 1",
			s:    NewIntSegment(NewBorder(Excluded, Int[uint8](0)), NewBorder(Included, Int[uint8](1))),
			args: args[uint8]{size: 1},
			want: []SplitSegment[uint8]{
				NewIntSegment(NewBorder(Included, Int[uint8](1)), NewBorder(Included, Int[uint8](1))),
			},
		},
		{
			name: "(0;1) by 1",
			s:    NewIntSegment(NewBorder(Excluded, Int[uint8](0)), NewBorder(Excluded, Int[uint8](1))),
			args: args[uint8]{size: 1},
			want: nil,
		},
		{
			name: "[1;2] by 0",
			s:    NewIntSegment(NewBorder(Included, Int[uint8](1)), NewBorder(Included, Int[uint8](2))),
			args: args[uint8]{size: 0},
			want: nil,
		},
		{
			name: "[min;max] by max",
			s:    NewIntSegment(NewBorder(Included, Int[uint8](0)), NewBorder(Included, Int[uint8](math.MaxUint8))),
			args: args[uint8]{size: math.MaxUint8},
			want: []SplitSegment[uint8]{
				NewIntSegment(NewBorder(Included, Int[uint8](0)), NewBorder(Excluded, Int[uint8](math.MaxUint8))),
				NewIntSegment(NewBorder(Included, Int[uint8](math.MaxUint8)), NewBorder(Included, Int[uint8](math.MaxUint8))),
			},
		},
		{
			name: "(min;min) by 1",
			s:    NewIntSegment(NewBorder(Excluded, Int[uint8](0)), NewBorder(Excluded, Int[uint8](0))),
			args: args[uint8]{size: 1},
			want: nil,
		},
		{
			name: "(max;max) by 1",
			s:    NewIntSegment(NewBorder(Excluded, Int[uint8](math.MaxUint8)), NewBorder(Excluded, Int[uint8](math.MaxUint8))),
			args: args[uint8]{size: 1},
			want: nil,
		},
		{
			name: "[inf;inf] by max",
			s:    NewIntSegment(NewUnbound[uint8](), NewUnbound[uint8]()),
			args: args[uint8]{size: math.MaxUint8},
			want: []SplitSegment[uint8]{
				NewIntSegment(NewBorder(Included, Int[uint8](0)), NewBorder(Excluded, Int[uint8](math.MaxUint8))),
				NewIntSegment(NewBorder(Included, Int[uint8](math.MaxUint8)), NewBorder(Included, Int[uint8](math.MaxUint8))),
			},
		},
		{
			name: "[2;1] by 1",
			s:    NewIntSegment(NewBorder(Included, Int[uint8](2)), NewBorder(Included, Int[uint8](1))),
			args: args[uint8]{size: 1},
			want: nil,
		},
		{
			name: "[min;min] by 1",
			s:    NewIntSegment(NewBorder(Included, Int[uint8](0)), NewBorder(Included, Int[uint8](0))),
			args: args[uint8]{size: 1},
			want: []SplitSegment[uint8]{
				NewIntSegment(NewBorder(Included, Int[uint8](0)), NewBorder(Included, Int[uint8](0))),
			},
		},
		{
			name: "[max;max] by 1",
			s:    NewIntSegment(NewBorder(Included, Int[uint8](math.MaxUint8)), NewBorder(Included, Int[uint8](math.MaxUint8))),
			args: args[uint8]{size: 1},
			want: []SplitSegment[uint8]{
				NewIntSegment(NewBorder(Included, Int[uint8](math.MaxUint8)), NewBorder(Included, Int[uint8](math.MaxUint8))),
			},
		},
		{
			name: "[min;min+1] by 1",
			s:    NewIntSegment(NewBorder(Included, Int[uint8](0)), NewBorder(Included, Int[uint8](0+1))),
			args: args[uint8]{size: 2},
			want: []SplitSegment[uint8]{
				NewIntSegment(NewBorder(Included, Int[uint8](0)), NewBorder(Included, Int[uint8](0+1))),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Split(tt.args.size).Collect()

			if !reflect.DeepEqual(got, tt.want) && (got != nil || tt.want != nil) {
				t.Errorf("%v .Split() by %v, got = %v, want %v", tt.s, tt.args.size, got, tt.want)
			}
		})
	}
}

func TestIntSegment_Split_WithErrors(t *testing.T) {
	type args[T intLike] struct {
		size T
	}
	type testCase[T intLike] struct {
		name      string
		s         *IntSegment[T]
		args      args[T]
		want      []SplitSegment[T]
		iteration int
		wantErr   bool
	}
	tests := []testCase[int64]{
		{
			name:      "[1;3] by 1, err 1",
			s:         NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](3))),
			args:      args[int64]{size: 1},
			iteration: 1,
			want: []SplitSegment[int64]{
				NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Excluded, Int[int64](2))),
			},
			wantErr: true,
		},
		{
			name:      "[1;3] by 1, err 2",
			s:         NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](3))),
			args:      args[int64]{size: 1},
			iteration: 2,
			want: []SplitSegment[int64]{
				NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Excluded, Int[int64](2))),
				NewIntSegment(NewBorder(Included, Int[int64](2)), NewBorder(Excluded, Int[int64](3))),
			},
			wantErr: true,
		},
		{
			name:      "[1;3] by 1, err 3",
			s:         NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](3))),
			args:      args[int64]{size: 1},
			iteration: 3,
			want: []SplitSegment[int64]{
				NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Excluded, Int[int64](2))),
				NewIntSegment(NewBorder(Included, Int[int64](2)), NewBorder(Excluded, Int[int64](3))),
				NewIntSegment(NewBorder(Included, Int[int64](3)), NewBorder(Included, Int[int64](3))),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []SplitSegment[int64]
			iter := 0
			tt.s.Split(tt.args.size)(func(s SplitSegment[int64], err error) bool {
				got = append(got, s)
				iter++
				if iter == tt.iteration {
					return false
				}
				return true
			})

			if !reflect.DeepEqual(got, tt.want) && (got != nil || tt.want != nil) {
				t.Errorf("%v .Split() by %v, got = %v, want %v", tt.s, tt.args.size, got, tt.want)
			}
		})
	}
}

func TestIntSegment_TryTo(t *testing.T) {
	type args struct {
		from Bound
		till Bound
	}
	type testCase[T intLike] struct {
		name    string
		s       *IntSegment[T]
		args    args
		want    TryToSegment[T]
		wantErr bool
	}
	tests := []testCase[int64]{
		// all to all
		// []
		{
			name: "[1;2] -> []",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](2))),
			args: args{
				from: Included,
				till: Included,
			},
			want: NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](2))),
		},
		{
			name: "[1;2] -> [)",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](2))),
			args: args{
				from: Included,
				till: Excluded,
			},
			want: NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Excluded, Int[int64](3))),
		},
		{
			name: "[1;2] -> (]",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](2))),
			args: args{
				from: Excluded,
				till: Included,
			},
			want: NewIntSegment(NewBorder(Excluded, Int[int64](0)), NewBorder(Included, Int[int64](2))),
		},
		{
			name: "[1;2] -> ()",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](2))),
			args: args{
				from: Excluded,
				till: Excluded,
			},
			want: NewIntSegment(NewBorder(Excluded, Int[int64](0)), NewBorder(Excluded, Int[int64](3))),
		},
		// (]
		{
			name: "(0;2] -> []",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](0)), NewBorder(Included, Int[int64](2))),
			args: args{
				from: Included,
				till: Included,
			},
			want: NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](2))),
		},
		{
			name: "(0;2] -> [)",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](0)), NewBorder(Included, Int[int64](2))),
			args: args{
				from: Included,
				till: Excluded,
			},
			want: NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Excluded, Int[int64](3))),
		},
		{
			name: "(0;2] -> (]",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](0)), NewBorder(Included, Int[int64](2))),
			args: args{
				from: Excluded,
				till: Included,
			},
			want: NewIntSegment(NewBorder(Excluded, Int[int64](0)), NewBorder(Included, Int[int64](2))),
		},
		{
			name: "(0;2] -> ()",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](0)), NewBorder(Included, Int[int64](2))),
			args: args{
				from: Excluded,
				till: Excluded,
			},
			want: NewIntSegment(NewBorder(Excluded, Int[int64](0)), NewBorder(Excluded, Int[int64](3))),
		},
		// ()
		{
			name: "(0;3) -> []",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](0)), NewBorder(Excluded, Int[int64](3))),
			args: args{
				from: Included,
				till: Included,
			},
			want: NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](2))),
		},
		{
			name: "(0;3) -> [)",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](0)), NewBorder(Excluded, Int[int64](3))),
			args: args{
				from: Included,
				till: Excluded,
			},
			want: NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Excluded, Int[int64](3))),
		},
		{
			name: "(0;3) -> (]",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](0)), NewBorder(Excluded, Int[int64](3))),
			args: args{
				from: Excluded,
				till: Included,
			},
			want: NewIntSegment(NewBorder(Excluded, Int[int64](0)), NewBorder(Included, Int[int64](2))),
		},
		{
			name: "(0;3) -> ()",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](0)), NewBorder(Excluded, Int[int64](3))),
			args: args{
				from: Excluded,
				till: Excluded,
			},
			want: NewIntSegment(NewBorder(Excluded, Int[int64](0)), NewBorder(Excluded, Int[int64](3))),
		},
		// [)
		{
			name: "[1;3) -> []",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Excluded, Int[int64](3))),
			args: args{
				from: Included,
				till: Included,
			},
			want: NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](2))),
		},
		{
			name: "[1;3) -> [)",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Excluded, Int[int64](3))),
			args: args{
				from: Included,
				till: Excluded,
			},
			want: NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Excluded, Int[int64](3))),
		},
		{
			name: "[1;3) -> (]",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Excluded, Int[int64](3))),
			args: args{
				from: Excluded,
				till: Included,
			},
			want: NewIntSegment(NewBorder(Excluded, Int[int64](0)), NewBorder(Included, Int[int64](2))),
		},
		{
			name: "[1;3) -> ()",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Excluded, Int[int64](3))),
			args: args{
				from: Excluded,
				till: Excluded,
			},
			want: NewIntSegment(NewBorder(Excluded, Int[int64](0)), NewBorder(Excluded, Int[int64](3))),
		},
		// other cases
		{
			name: "(inf;inf) -> []",
			s:    NewIntSegment(NewUnbound[int64](), NewUnbound[int64]()),
			args: args{
				from: Included,
				till: Included,
			},
			want: NewIntSegment(NewUnbound[int64](), NewUnbound[int64]()),
		},
		{
			name: "(0;inf) -> []",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](0)), NewUnbound[int64]()),
			args: args{
				from: Included,
				till: Included,
			},
			want: NewIntSegment(NewBorder(Included, Int[int64](1)), NewUnbound[int64]()),
		},
		{
			name: "(min;min) -> []",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](math.MinInt64)), NewBorder(Excluded, Int[int64](math.MinInt64))),
			args: args{
				from: Included,
				till: Included,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "(max;max) -> []",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](math.MaxInt64)), NewBorder(Excluded, Int[int64](math.MaxInt64))),
			args: args{
				from: Included,
				till: Included,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "(min;max) -> []",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](math.MinInt64)), NewBorder(Excluded, Int[int64](math.MaxInt64))),
			args: args{
				from: Included,
				till: Included,
			},
			want: NewIntSegment(NewBorder(Included, Int[int64](math.MinInt64+1)), NewBorder(Included, Int[int64](math.MaxInt64-1))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.TryTo(tt.args.from, tt.args.till)
			if (err != nil) != tt.wantErr {
				t.Errorf("TryTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TryTo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSegment_Iterate(t *testing.T) {
	type testCase[T intLike] struct {
		name string
		s    *IntSegment[T]
		want []T
	}
	tests := []testCase[int64]{
		{
			name: "[1;2]",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](2))),
			want: []int64{1, 2},
		},
		{
			name: "(MAX;MAX]",
			s:    NewIntSegment(NewBorder(Excluded, Int[int64](math.MaxInt64)), NewBorder(Included, Int[int64](math.MaxInt64))),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Iterate().Collect(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Iterate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSegment_Iterate_Interrupt(t *testing.T) {
	type testCase[T intLike] struct {
		name string
		s    *IntSegment[T]
		want []T
	}
	tests := []testCase[int64]{
		{
			name: "[1;2]",
			s:    NewIntSegment(NewBorder(Included, Int[int64](1)), NewBorder(Included, Int[int64](2))),
			want: []int64{1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []int64
			tt.s.Iterate()(func(i int64, err error) bool {
				got = append(got, i)
				return false
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Iterate() = %v, want %v", got, tt.want)
			}
		})
	}
}
