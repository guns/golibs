// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package optimized

import "testing"

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
