// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package calculate

import "testing"

func TestNextCap(t *testing.T) {
	data := []struct {
		size, cap int
	}{
		{0, 8},
		{1, 8},
		{8, 8},
		{9, 16},
		{15, 16},
		{16, 16},
		{20, 32},
		{200, 256},
		{256, 256},
		{1025, 2048},
		{3000, 4096},
		{4096, 4096},
		{4097, 4096 + 2048},
		{4096 + 1024, 4096 + 2048},
		{4096 + 2048, 4096 + 2048},
		{4096 + 3072, 4096 + 4096},
		{4096 + 4096, 4096 + 4096},
		{4096 + 4097, 4096*2 + 4096},
		{10000, 3 << 12},
		{25000, 1 << 15},
		{1<<19 + 1, 3 << 18},
		{1 << 20, 1 << 20},
		{123456789, 1 << 27},
		{1<<27 + 1, 5 << 25},
		{1<<27 + 1<<25 + 1, 6 << 25},
		{1<<27 + 2<<25 + 1, 7 << 25},
		{1<<27 + 3<<25 + 1, 8 << 25},
	}

	for _, row := range data {
		if NextCap(row.size) != row.cap {
			t.Errorf("NextCap(%d) -> %v != %v", row.size, NextCap(row.size), row.cap)
		}
	}
}
