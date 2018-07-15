// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func randNumberSlice(n, max int) []Number {
	v := make([]Number, n)
	if max > 0 {
		for i := range v {
			v[i] = Number(rand.Intn(max))
		}
	} else {
		for i := range v {
			v[i] = Number(i)
		}
		// Fisher-Yates shuffle
		for i := len(v) - 1; i >= 0; i-- {
			j := rand.Intn(i + 1)
			v[i], v[j] = v[j], v[i]
		}
	}
	return v
}

func TestQuicksortNumberSlice(t *testing.T) {
	for i := 0; i < 100; i++ {
		var r []Number

		if rand.Intn(2) == 0 {
			r = randNumberSlice(10*i, 0)
		} else {
			r = randNumberSlice(10*i, 5*i)
		}

		s1 := make([]int, len(r))
		s2 := make([]int, len(r))

		for i := range r {
			s1[i] = int(r[i])
		}

		sort.Ints(s1)
		QuicksortNumberSlice(r)

		for i := range r {
			s2[i] = int(r[i])
		}

		if !reflect.DeepEqual(s2, s1) {
			t.Logf("QuicksortNumberSlice:")
			t.Logf("%v !=", s2)
			t.Logf("%v", s1)
			t.Fail()
		}
	}
}

func TestMedianOfThreeNumberSamples(t *testing.T) {
	data := [][]Number{
		{0, 1, 2},
		{0, 2, 1},
		{1, 0, 2},
		{1, 2, 0},
		{2, 0, 1},
		{2, 1, 0},
	}

	for _, v := range data {
		if MedianOfThreeNumberSamples(v) != 1 {
			t.Errorf("MedianOfThreeNumberSamples(%v): %v != %v", v, MedianOfThreeNumberSamples(v), 1)
		}
	}
}
