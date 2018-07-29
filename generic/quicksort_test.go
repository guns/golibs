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

func randGenericNumberSlice(n, max int) []GenericNumber {
	v := make([]GenericNumber, n)
	if max > 0 {
		for i := range v {
			v[i] = GenericNumber(rand.Intn(max))
		}
	} else {
		for i := range v {
			v[i] = GenericNumber(i)
		}
		rand.Shuffle(len(v), func(i, j int) {
			v[i], v[j] = v[j], v[i]
		})
	}
	return v
}

func TestQuicksortGenericNumberSlice(t *testing.T) {
	for i := 0; i < 100; i++ {
		var r []GenericNumber

		if rand.Intn(2) == 0 {
			r = randGenericNumberSlice(10*i, 0)
		} else {
			r = randGenericNumberSlice(10*i, 5*i)
		}

		s1 := make([]int, len(r))
		s2 := make([]int, len(r))

		for i := range r {
			s1[i] = int(r[i])
		}

		sort.Ints(s1)
		QuicksortGenericNumberSlice(r)

		for i := range r {
			s2[i] = int(r[i])
		}

		if !reflect.DeepEqual(s2, s1) {
			t.Logf("QuicksortGenericNumberSlice:")
			t.Logf("%v !=", s2)
			t.Logf("%v", s1)
			t.Fail()
		}
	}
}

func TestMedianOfThreeGenericNumberSamples(t *testing.T) {
	data := [][]GenericNumber{
		{0, 1, 2},
		{0, 2, 1},
		{1, 0, 2},
		{1, 2, 0},
		{2, 0, 1},
		{2, 1, 0},
	}

	for _, v := range data {
		if MedianOfThreeGenericNumberSamples(v) != 1 {
			t.Errorf("MedianOfThreeGenericNumberSamples(%v): %v != %v", v, MedianOfThreeGenericNumberSamples(v), 1)
		}
	}
}
