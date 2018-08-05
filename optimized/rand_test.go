// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package optimized

import (
	"math/rand"
	"testing"
)

//	go version go1.10.3 linux/amd64
//	goos: linux
//	goarch: amd64
//	pkg: github.com/guns/golibs/optimized
//	Benchmark02RandIntn             50000000                27.8 ns/op             0 B/op          0 allocs/op
//	Benchmark02OptimizedRandIntn    100000000               23.7 ns/op             0 B/op          0 allocs/op
//	Benchmark02RandInt63n           50000000                35.9 ns/op             0 B/op          0 allocs/op
//	Benchmark02OptimizedRandInt63n  100000000               22.6 ns/op             0 B/op          0 allocs/op
//	Benchmark02RandInt31n           50000000                26.1 ns/op             0 B/op          0 allocs/op
//	Benchmark02OptimizedRandInt31n  100000000               23.9 ns/op             0 B/op          0 allocs/op
//	PASS
//	ok      github.com/guns/golibs/optimized        11.682s

func Benchmark02RandIntn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = rand.Intn(0xffff)
	}
}
func Benchmark02OptimizedRandIntn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandIntn(0xffff)
	}
}
func Benchmark02RandInt63n(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = rand.Int63n(0xffff)
	}
}
func Benchmark02OptimizedRandInt63n(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandInt63n(0xffff)
	}
}
func Benchmark02RandInt31n(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = rand.Int31n(0xffff)
	}
}
func Benchmark02OptimizedRandInt31n(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandInt31n(0xffff)
	}
}
