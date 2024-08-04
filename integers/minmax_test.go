package segment_int

import (
	"fmt"
	"math"
	"testing"
)

func Test_intType(t *testing.T) {
	testsTypes := map[string]interface{}{
		"int":    int(1),
		"int8":   int8(1),
		"int16":  int16(1),
		"int32":  int32(1),
		"int64":  int64(1),
		"uint":   uint(1),
		"uint8":  uint8(1),
		"uint16": uint16(1),
		"uint32": uint32(1),
		"uint64": uint64(1),
	}
	testsPointers := map[string]interface{}{
		"int":    new(int),
		"int8":   new(int8),
		"int16":  new(int16),
		"int32":  new(int32),
		"int64":  new(int64),
		"uint":   new(uint),
		"uint8":  new(uint8),
		"uint16": new(uint16),
		"uint32": new(uint32),
		"uint64": new(uint64),
	}
	for _, tests := range map[string]map[string]interface{}{"value": testsTypes, "pointer": testsPointers} {
		for want, val := range tests {
			name := fmt.Sprintf("%s %T", want, val)
			t.Run(name, func(t *testing.T) {
				got, err := intType(val)
				if err != nil {
					t.Errorf("intType() error = %v", err)
					return
				}
				if got != want {
					t.Errorf("intType() got = %v, want %v", got, want)
				}
			})
		}
	}

	t.Run("float32", func(t *testing.T) {
		_, err := intType(float32(1))
		if (err != nil) && err.Error() != "unknown type float32" {
			t.Errorf("intType() error = %v, wantErr %v", err, "unknown type float32")
			return
		}
	})
}

func maxIntRun[T intLike](t *testing.T, name string, want T) {
	t.Run(name, func(t *testing.T) {
		got := maxInt[T]()
		if got != want {
			t.Errorf("maxInt() got = %v, want %v", got, want)
		}
	})
}

func Test_maxInt(t *testing.T) {
	maxIntRun[int](t, "int", int(math.MaxInt))
	maxIntRun[int8](t, "int8", int8(math.MaxInt8))
	maxIntRun[int16](t, "int16", int16(math.MaxInt16))
	maxIntRun[int32](t, "int32", int32(math.MaxInt32))
	maxIntRun[int64](t, "int64", int64(math.MaxInt64))
	maxIntRun[uint](t, "uint", uint(math.MaxUint))
	maxIntRun[uint8](t, "uint8", uint8(math.MaxUint8))
	maxIntRun[uint16](t, "uint16", uint16(math.MaxUint16))
	maxIntRun[uint32](t, "uint32", uint32(math.MaxUint32))
	maxIntRun[uint64](t, "uint64", uint64(math.MaxUint64))
}

func minIntRun[T intLike](t *testing.T, name string, want T) {
	t.Run(name, func(t *testing.T) {
		got := minInt[T]()
		if got != want {
			t.Errorf("minInt() got = %v, want %v", got, want)
		}
	})
}

func Test_minInt(t *testing.T) {
	minIntRun[int](t, "int", int(math.MinInt))
	minIntRun[int8](t, "int8", int8(math.MinInt8))
	minIntRun[int16](t, "int16", int16(math.MinInt16))
	minIntRun[int32](t, "int32", int32(math.MinInt32))
	minIntRun[int64](t, "int64", int64(math.MinInt64))
	minIntRun[uint](t, "uint", uint(0))
	minIntRun[uint8](t, "uint8", uint8(0))
	minIntRun[uint16](t, "uint16", uint16(0))
	minIntRun[uint32](t, "uint32", uint32(0))
	minIntRun[uint64](t, "uint64", uint64(0))
}
