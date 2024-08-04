package segment

import (
	"reflect"
	"testing"
)

func TestInf(t *testing.T) {
	type testCase[T any] struct {
		name string
		want Value[T]
	}
	tests := []testCase[int64]{
		{
			name: "default",
			want: &infValue[int64]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Inf[int64](); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Inf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_infValue_Next(t *testing.T) {
	type testCase[T any] struct {
		name  string
		i     infValue[T]
		want  Value[T]
		want1 error
	}
	tests := []testCase[int64]{
		{
			name:  "default",
			i:     infValue[int64]{},
			want:  Inf[int64](),
			want1: ErrHasNoNextValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.i.Next()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Next() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Next() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_infValue_Prev(t *testing.T) {
	type testCase[T any] struct {
		name  string
		i     infValue[T]
		want  Value[T]
		want1 error
	}
	tests := []testCase[int64]{
		{
			name:  "default",
			i:     infValue[int64]{},
			want:  Inf[int64](),
			want1: ErrHasNoPrevValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.i.Prev()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Prev() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Prev() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_infValue_Value(t *testing.T) {
	type testCase[T any] struct {
		name string
		i    infValue[T]
		want T
	}
	tests := []testCase[int64]{
		{
			name: "default",
			i:    infValue[int64]{},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_infValue_String(t *testing.T) {
	inf := Inf[int64]()
	want := "inf"
	if got := inf.String(); got != want {
		t.Errorf("String() = %v, want %v", got, want)
	}
}
