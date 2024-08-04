package segment_int

import (
	rng "github.com/pioniro/segment-go"
	"math"
	"reflect"
	"testing"
)

func TestInt(t *testing.T) {
	type args[T intLike] struct {
		v T
	}
	type testCase[T intLike] struct {
		name string
		args args[T]
		want rng.Value[T]
	}
	tests := []testCase[int64]{
		{
			name: "default",
			args: args[int64]{math.MaxInt64},
			want: &intValue[int64]{math.MaxInt64},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_intValue_Next(t *testing.T) {
	type testCase[T intLike] struct {
		name  string
		v     rng.Value[T]
		want  rng.Value[T]
		want1 error
	}
	tests := []testCase[int64]{
		{
			name:  "100",
			v:     Int[int64](100),
			want:  Int[int64](101),
			want1: nil,
		},
		{
			name:  "min",
			v:     Int[int64](math.MinInt64),
			want:  Int[int64](math.MinInt64 + 1),
			want1: nil,
		},
		{
			name:  "max",
			v:     Int[int64](math.MaxInt64),
			want:  Int[int64](math.MaxInt64),
			want1: rng.ErrHasNoNextValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.v.Next()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Next() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Next() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_intValue_Prev(t *testing.T) {
	type testCase[T intLike] struct {
		name  string
		v     rng.Value[T]
		want  rng.Value[T]
		want1 error
	}
	tests := []testCase[int64]{
		{
			name:  "100",
			v:     Int[int64](100),
			want:  Int[int64](99),
			want1: nil,
		},
		{
			name:  "max",
			v:     Int[int64](math.MaxInt64),
			want:  Int[int64](math.MaxInt64 - 1),
			want1: nil,
		},
		{
			name:  "min",
			v:     Int[int64](math.MinInt64),
			want:  Int[int64](math.MinInt64),
			want1: rng.ErrHasNoPrevValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.v.Prev()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Prev() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Prev() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_intValue_String(t *testing.T) {
	type testCase[T intLike] struct {
		name string
		v    rng.Value[T]
		want string
	}
	tests := []testCase[int64]{
		{
			name: "default",
			v:    Int[int64](100),
			want: "100",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_intValue_String_Uint64(t *testing.T) {
	type testCase[T intLike] struct {
		name string
		v    rng.Value[T]
		want string
	}
	tests := []testCase[uint64]{
		{
			name: "default",
			v:    Int[uint64](math.MaxUint64),
			want: "18446744073709551615",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_intValue_Value(t *testing.T) {
	type testCase[T intLike] struct {
		name string
		v    rng.Value[T]
		want T
	}
	tests := []testCase[int64]{
		{
			name: "default",
			v:    Int[int64](100),
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Value(); got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
