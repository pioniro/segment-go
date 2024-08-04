package main

import (
	"fmt"
	"github.com/pioniro/segment-go"
	rngint "github.com/pioniro/segment-go/integers"
)

func simpleExample() {
	// make a segment [1; 10)
	r := rngint.NewIntSegment(segment.NewIncluded(rngint.Int(1)), segment.NewExcluded(rngint.Int(10)))
	// iterate over the segment
	for _, val := range r.Iterate().Collect() {
		// Print 123456789
		fmt.Print(val)
	}
}

func splitExample() {
	// make a segment [1; 10)
	r := rngint.NewIntSegment(segment.NewIncluded(rngint.Int(1)), segment.NewExcluded(rngint.Int(10)))
	// split segment into 5 parts: [1; 3), [3; 5), [5; 7), [7; 9), [9; 9]
	for _, rr := range r.Split(2).Collect() {
		fmt.Println(rr)
		// Print
		// [1;3)
		// [3;5)
		// [5;7)
		// [7;9)
		// [9;9]

	}
}
