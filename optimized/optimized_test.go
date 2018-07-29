// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package optimized

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	os.Exit(m.Run())
}

func TestMul64(t *testing.T) {
	data := []struct {
		x, y, lower, upper uint64
	}{
		{
			x:     0,
			y:     0,
			upper: 0,
			lower: 0,
		},
		{
			x:     0,
			y:     1,
			upper: 0,
			lower: 0,
		},
		{
			x:     0x2c51755c,
			y:     0x9e00f286,
			upper: 0,
			lower: 0x1b5a706afb946628,
		},
		{
			x:     1<<64 - 1,
			y:     2,
			upper: 1,
			lower: 1<<64 - 2,
		},
		{
			x:     1<<64 - 1,
			y:     1<<64 - 1,
			upper: 0xfffffffffffffffe,
			lower: 1,
		},
		{
			x:     0x736b9f3f93cad329,
			y:     0x341afaad2b00aaf4,
			upper: 0x177e064c431c4a9b,
			lower: 0xe097ceea708a7d14,
		},
	}

	fns := []func(x, y uint64) (lo, hi uint64){
		Mul64,
		mul64,
	}

	for _, row := range data {
		for i, f := range fns {
			lo, hi := f(row.x, row.y)

			if lo != row.lower || hi != row.upper {
				t.Logf("fns[%d](0x%x, 0x%x) ->", i, row.x, row.y)
				t.Logf("\t(0x%x, 0x%x) !=", lo, hi)
				t.Logf("\t(0x%x, 0x%x)", row.lower, row.upper)
				t.Fail()
			}
		}
	}

	// Test implementations against each other
	const rounds = 1000000
	for i := 0; i < rounds; i++ {
		x, y := rand.Uint64(), rand.Uint64()
		lo0, hi0 := Mul64(x, y)
		lo1, hi1 := mul64(x, y)

		if lo0 != lo1 || hi0 != hi1 {
			t.Logf("Mul64(0x%x, 0x%x) != fallback", x, y)
			t.Logf("\t(0x%x, 0x%x) !=", lo0, hi0)
			t.Logf("\t(0x%x, 0x%x)", lo1, hi1)
			t.Fail()
			break
		}
	}
}

func TestMul32(t *testing.T) {
	data := []struct {
		x, y, lower, upper uint32
	}{
		{
			x:     0,
			y:     0,
			upper: 0,
			lower: 0,
		},
		{
			x:     0,
			y:     1,
			upper: 0,
			lower: 0,
		},
		{
			x:     0x755c,
			y:     0xf286,
			upper: 0,
			lower: 0x6f2e6628,
		},
		{
			x:     1<<32 - 1,
			y:     2,
			upper: 1,
			lower: 1<<32 - 2,
		},
		{
			x:     1<<32 - 1,
			y:     1<<32 - 1,
			upper: 0xfffffffe,
			lower: 1,
		},
		{
			x:     0x93cad329,
			y:     0x2b00aaf4,
			upper: 0x18d37429,
			lower: 0x708a7d14,
		},
	}

	fns := []func(x, y uint32) (lo, hi uint32){
		Mul32,
		mul32,
	}

	for _, row := range data {
		for i, f := range fns {
			lo, hi := f(row.x, row.y)

			if lo != row.lower || hi != row.upper {
				t.Logf("fns[%d](0x%x, 0x%x) ->", i, row.x, row.y)
				t.Logf("\t(0x%x, 0x%x) !=", lo, hi)
				t.Logf("\t(0x%x, 0x%x)", row.lower, row.upper)
				t.Fail()
			}
		}
	}

	// Test implementations against each other
	const rounds = 1000000
	for i := 0; i < rounds; i++ {
		x, y := rand.Uint32(), rand.Uint32()
		lo0, hi0 := Mul32(x, y)
		lo1, hi1 := mul32(x, y)

		if lo0 != lo1 || hi0 != hi1 {
			t.Logf("Mul32(0x%x, 0x%x) != fallback", x, y)
			t.Logf("\t(0x%x, 0x%x) !=", lo0, hi0)
			t.Logf("\t(0x%x, 0x%x)", lo1, hi1)
			t.Fail()
			break
		}
	}
}

//	:: go version go1.10.3 linux/amd64
//	goos: linux
//	goarch: amd64
//	pkg: github.com/guns/golibs/optimized
//	BenchmarkMul-4              1000000000         1.98 ns/op        0 B/op        0 allocs/op
//	BenchmarkMulFallback-4      500000000          3.82 ns/op        0 B/op        0 allocs/op
//	BenchmarkMul64-4            1000000000         1.99 ns/op        0 B/op        0 allocs/op
//	BenchmarkMul64Fallback-4    500000000          3.89 ns/op        0 B/op        0 allocs/op
//	BenchmarkMul32-4            2000000000         1.87 ns/op        0 B/op        0 allocs/op
//	BenchmarkMul32Fallback-4    500000000          3.83 ns/op        0 B/op        0 allocs/op
//	PASS
//	ok   github.com/guns/golibs/optimized 15.351s

func BenchmarkMul(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Mul(0x93cad329, 0x2b00aaf4)
	}
}
func BenchmarkMulFallback(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = mul(0x93cad329, 0x2b00aaf4)
	}
}
func BenchmarkMul64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Mul64(0x736b9f3f93cad329, 0x341afaad2b00aaf4)
	}
}
func BenchmarkMul64Fallback(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = mul64(0x736b9f3f93cad329, 0x341afaad2b00aaf4)
	}
}
func BenchmarkMul32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Mul32(0x93cad329, 0x2b00aaf4)
	}
}
func BenchmarkMul32Fallback(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = mul32(0x93cad329, 0x2b00aaf4)
	}
}

//	:: go version go1.10.3 linux/amd64
//	goos: linux
//	goarch: amd64
//	pkg: github.com/guns/golibs/optimized
//	BenchmarkOptimizedRandIntn-4      100000000        23.0 ns/op        0 B/op        0 allocs/op
//	BenchmarkRandIntn-4               50000000         28.8 ns/op        0 B/op        0 allocs/op
//	BenchmarkOptimizedRandInt63n-4    100000000        22.8 ns/op        0 B/op        0 allocs/op
//	BenchmarkRandInt63n-4             50000000         36.6 ns/op        0 B/op        0 allocs/op
//	BenchmarkOptimizedRandInt31n-4    50000000         24.5 ns/op        0 B/op        0 allocs/op
//	BenchmarkRandInt31n-4             50000000         26.9 ns/op        0 B/op        0 allocs/op
//	PASS
//	ok   github.com/guns/golibs/optimized 10.694s

func BenchmarkOptimizedRandIntn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandIntn(i + 1)
	}
}
func BenchmarkRandIntn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = rand.Intn(i + 1)
	}
}
func BenchmarkOptimizedRandInt63n(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandInt63n(int64(i) + 1)
	}
}
func BenchmarkRandInt63n(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = rand.Int63n(int64(i) + 1)
	}
}
func BenchmarkOptimizedRandInt31n(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandInt31n(int32(i) + 1)
	}
}
func BenchmarkRandInt31n(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = rand.Int31n(int32(i) + 1)
	}
}
