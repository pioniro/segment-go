Segment-Go
====

[![Go Reference](https://pkg.go.dev/badge/github.com/pioniro/segment-go.svg)](https://pkg.go.dev/github.com/pioniro/segment-go)
[![Build status](https://img.shields.io/circleci/build/github/pioniro/segment-go?style=plastic)](https://app.circleci.com/pipelines/github/pioniro/segment-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/pioniro/segment-go)](https://goreportcard.com/report/github.com/pioniro/segment-go)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/pioniro/segment-go/blob/main/LICENSE)
[![Coverage](https://coveralls.io/repos/github/pioniro/segment-go/badge.svg?branch=main)](https://coveralls.io/github/pioniro/segment-go)


Segment is a Go library for working with ranges.

## Installation

To install Segment-Go, you need to install Go and set your Go workspace first.

1. Download Segment-Go: `go get github.com/pioniro/segment-go`
2. Import it in your code: `import "github.com/pioniro/segment-go"`

## Quick Start

Here's a quick example of how to use Segment-Go:

```go
package main

import (
  "fmt"
  "github.com/pioniro/segment-go"
  seg "github.com/pioniro/segment-go/integers"
)

func main() {
  // make a segment [1; 10)
  r := seg.NewIntSegment(segment.NewIncluded(seg.Int(1)), segment.NewExcluded(seg.Int(10)))
  // iterate over the segment
  for _, val := range r.Iterate().Collect() {
    // Print 123456789
    fmt.Print(val)
  }
}
```


## Features

- **Included, Excluded, Unbound**: Segment boundaries can be included in the segment, excluded, or not limited at all.
- **Split**: Segments can be split into multiple segments no larger than a specified length.
- **Includes**: Check for the inclusion of a value in a segment.
- **Iterable**: The ability to go through all the values of the segment.