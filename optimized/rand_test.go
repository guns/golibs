// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package optimized

import (
	"math/rand"
	"testing"
)

//	:: go version go1.10.3 linux/amd64
//	goos: linux
//	goarch: amd64
//	pkg: github.com/guns/golibs/optimized
//	Benchmark02OptimizedRandIntn-4      100000000        22.8 ns/op        0 B/op        0 allocs/op
//	Benchmark02RandIntn-4               50000000         28.5 ns/op        0 B/op        0 allocs/op
//	Benchmark02OptimizedRandInt63n-4    100000000        22.9 ns/op        0 B/op        0 allocs/op
//	Benchmark02RandInt63n-4             50000000         36.5 ns/op        0 B/op        0 allocs/op
//	Benchmark02OptimizedRandInt31n-4    50000000         24.5 ns/op        0 B/op        0 allocs/op
//	Benchmark02RandInt31n-4             50000000         26.7 ns/op        0 B/op        0 allocs/op
//	PASS
//	ok   github.com/guns/golibs/optimized 10.640s

func Benchmark02OptimizedRandIntn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandIntn(i + 1)
	}
}
func Benchmark02RandIntn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = rand.Intn(i + 1)
	}
}
func Benchmark02OptimizedRandInt63n(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandInt63n(int64(i) + 1)
	}
}
func Benchmark02RandInt63n(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = rand.Int63n(int64(i) + 1)
	}
}
func Benchmark02OptimizedRandInt31n(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandInt31n(int32(i) + 1)
	}
}
func Benchmark02RandInt31n(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = rand.Int31n(int32(i) + 1)
	}
}
