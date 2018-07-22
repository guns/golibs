// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package genericbenchmarks

import "testing"

var numbers = []int{45, 76, 82, 77, 60, 29, 40, 94, 32, 14, 12, 95, 92, 36, 38, 70, 43, 90, 20, 46, 8, 71, 80, 30, 67, 33, 47, 74, 35, 61, 25, 98, 91, 63, 42, 54, 5, 55, 23, 41, 11, 34, 68, 99, 15, 78, 31, 6, 26, 56, 83, 57, 58, 87, 28, 21, 73, 13, 10, 44, 86, 88, 75, 96, 52, 65, 59, 27, 93, 66, 17, 69, 3, 62, 81, 53, 7, 72, 1, 22, 16, 37, 85, 18, 50, 19, 2, 4, 0, 9, 64, 49, 24, 39, 97, 84, 48, 89, 51, 79}

// goos: linux
// goarch: amd64
// pkg: github.com/guns/golibs/generic/genericbenchmarks
// BenchmarkMinInt2-4              2000000000               0.36 ns/op            0 B/op          0 allocs/op
// BenchmarkMinIntV2-4             50000000                22.7 ns/op             0 B/op          0 allocs/op
// BenchmarkMinInt3-4              100000000               16.0 ns/op             0 B/op          0 allocs/op
// BenchmarkMinIntV3-4             50000000                31.5 ns/op             0 B/op          0 allocs/op
// BenchmarkMinIntLoop-4           20000000                79.5 ns/op             0 B/op          0 allocs/op
// BenchmarkMinIntSlice-4          20000000               107 ns/op               0 B/op          0 allocs/op
// PASS
// ok      github.com/guns/golibs/generic/genericbenchmarks        9.082s

func BenchmarkMinInt2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MinInt2(i%len(numbers), (i+1)%len(numbers))
	}
}

func BenchmarkMinIntV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MinIntV(i%len(numbers), (i+1)%len(numbers))
	}
}

func BenchmarkMinInt3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MinInt3(i%len(numbers), (i+1)%len(numbers), (i+2)%len(numbers))
	}
}

func BenchmarkMinIntV3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MinIntV(i%len(numbers), (i+1)%len(numbers), (i+2)%len(numbers))
	}
}

func BenchmarkMinIntLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		min := numbers[0]
		for _, n := range numbers[1:] {
			if n < min {
				min = n
			}
		}
		_ = min
	}
}

func BenchmarkMinIntSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = MinIntSlice(numbers)
	}
}
