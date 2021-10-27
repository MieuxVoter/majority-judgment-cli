package benchmark

// from https://github.com/piersy/iterate-slice-bench-go
// with fixes (methink)

import (
	"testing"
)

var Val = uint64(42)
var slice = make([]uint64, 500*500)

func BenchmarkRangeReadSliceByIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := range slice {
			Val = slice[j] + 1
		}
	}
}

func BenchmarkRangeReadSliceByValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, v := range slice {
			Val = v + 1
		}
	}
}

func BenchmarkRangeWriteSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := range slice {
			slice[j] = Val + 1
		}
	}
}

func BenchmarkRangeReadAndWriteSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j, v := range slice {
			slice[j] = v + 1
		}
	}
}

func BenchmarkForIterReadSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := len(slice)
		for j := 0; j < l; j++ {
			Val = slice[j] + 1
		}
	}
}

func BenchmarkForIterWriteSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := len(slice)
		for j := 0; j < l; j++ {
			slice[j] = Val + 1
		}
	}
}

func BenchmarkForIterReadAndWriteSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := len(slice)
		for j := 0; j < l; j++ {
			slice[j] = slice[j] + 1
		}
	}
}

func BenchmarkDecreasingForIterReadAndWriteSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := len(slice)
		for j := l - 1; j >= 0; j-- {
			slice[j] = slice[j] + 1
		}
	}
}
