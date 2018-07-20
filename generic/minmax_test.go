// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

import "testing"

func TestMinMax(t *testing.T) {
	type N = GenericNumber

	data := []struct {
		v        []N
		min, max [8]N
	}{
		{
			v:   []N{-2, -1, 0, 1, 2},
			min: [8]N{-2, -2, 0, -2, -2, -2, -2, -2},
			max: [8]N{-1, 0, 0, -2, -1, 0, 1, 2},
		},
		{
			v:   []N{2, 3, -1, 4, 5},
			min: [8]N{2, -1, 0, 2, 2, -1, -1, -1},
			max: [8]N{3, 3, 0, 2, 3, 3, 4, 5},
		},
		{
			v:   []N{-2, -3, -1, 0, 1},
			min: [8]N{-3, -3, 0, -2, -3, -3, -3, -3},
			max: [8]N{-2, -1, 0, -2, -2, -1, 0, 1},
		},
		{
			v:   []N{1, 0, -1, 0, 1},
			min: [8]N{0, -1, 0, 1, 0, -1, -1, -1},
			max: [8]N{1, 1, 0, 1, 1, 1, 1, 1},
		},
	}

	for _, row := range data {
		v := row.v
		min := [8]N{}
		max := [8]N{}

		min[0] = MinGenericNumber2(v[0], v[1])
		min[1] = MinGenericNumber3(v[0], v[1], v[2])
		min[2] = MinGenericNumberV()
		min[3] = MinGenericNumberV(v[0])
		min[4] = MinGenericNumberV(v[0], v[1])
		min[5] = MinGenericNumberV(v[0], v[1], v[2])
		min[6] = MinGenericNumberV(v[0], v[1], v[2], v[3])
		min[7] = MinGenericNumberV(v[0], v[1], v[2], v[3], v[4])

		max[0] = MaxGenericNumber2(v[0], v[1])
		max[1] = MaxGenericNumber3(v[0], v[1], v[2])
		max[2] = MaxGenericNumberV()
		max[3] = MaxGenericNumberV(v[0])
		max[4] = MaxGenericNumberV(v[0], v[1])
		max[5] = MaxGenericNumberV(v[0], v[1], v[2])
		max[6] = MaxGenericNumberV(v[0], v[1], v[2], v[3])
		max[7] = MaxGenericNumberV(v[0], v[1], v[2], v[3], v[4])

		if min != row.min {
			t.Errorf("%v != %v", min, row.min)
		}
		if max != row.max {
			t.Errorf("%v != %v", max, row.max)
		}
	}
}
