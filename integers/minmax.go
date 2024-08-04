package segment_int

import (
	"fmt"
	"math"
)

var maxValues = map[string]uint64{
	"int":    uint64(math.MaxInt),
	"int8":   uint64(math.MaxInt8),
	"int16":  uint64(math.MaxInt16),
	"int32":  uint64(math.MaxInt32),
	"int64":  uint64(math.MaxInt64),
	"uint":   uint64(math.MaxUint),
	"uint8":  uint64(math.MaxUint8),
	"uint16": uint64(math.MaxUint16),
	"uint32": uint64(math.MaxUint32),
	"uint64": uint64(math.MaxUint64),
}

var minValues = map[string]int64{
	"int":    int64(math.MinInt),
	"int8":   int64(math.MinInt8),
	"int16":  int64(math.MinInt16),
	"int32":  int64(math.MinInt32),
	"int64":  int64(math.MinInt64),
	"uint":   int64(0),
	"uint8":  int64(0),
	"uint16": int64(0),
	"uint32": int64(0),
	"uint64": int64(0),
}

// intType returns a string representation of a given type.
// Require [1] All return values must be listed in minValues and maxValues
func intType(t interface{}) (string, error) {
	switch t.(type) {
	case int:
		return "int", nil
	case *int:
		return "int", nil
	case int8:
		return "int8", nil
	case *int8:
		return "int8", nil
	case int16:
		return "int16", nil
	case *int16:
		return "int16", nil
	case int32:
		return "int32", nil
	case *int32:
		return "int32", nil
	case int64:
		return "int64", nil
	case *int64:
		return "int64", nil
	case uint:
		return "uint", nil
	case *uint:
		return "uint", nil
	case uint8:
		return "uint8", nil
	case *uint8:
		return "uint8", nil
	case uint16:
		return "uint16", nil
	case *uint16:
		return "uint16", nil
	case uint32:
		return "uint32", nil
	case *uint32:
		return "uint32", nil
	case uint64:
		return "uint64", nil
	case *uint64:
		return "uint64", nil
	}
	return "", fmt.Errorf("unknown type %T", t)
}

// Ensure [1] T exists in intType because intLike Require [1]
func minInt[T intLike]() T {
	t := new(T)
	// errors are ignored, because Ensure [1]
	tt, _ := intType(t)
	// Ensure [2] T exists in minValues because intLike Require [2]
	return T(minValues[tt])
}

// Ensure [1] T exists in intType because intLike Require [1]
func maxInt[T intLike]() T {
	t := new(T)
	// errors are ignored, because we know, that t is intLike
	tt, _ := intType(t)
	// Ensure [2] T exists in minValues because intLike Require [2]
	return T(maxValues[tt])
}
