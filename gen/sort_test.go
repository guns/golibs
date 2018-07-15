package gen

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func randslice(n, max int) []int {
	v := make([]int, n)
	if max > 0 {
		for i := range v {
			v[i] = rand.Intn(max)
		}
	} else {
		for i := range v {
			v[i] = i
		}
		// Fisher-Yates shuffle
		for i := len(v) - 1; i >= 0; i-- {
			j := rand.Intn(i + 1)
			v[i], v[j] = v[j], v[i]
		}
	}
	return v
}

func TestQuickSort(t *testing.T) {
	var r []int

	for i := 0; i < 100; i++ {
		if rand.Intn(2) == 0 {
			r = randslice(10*i, 0)
		} else {
			r = randslice(10*i, 5*i)
		}

		r1 := make([]Number, len(r))
		r2 := make([]Type, len(r))
		s := make([]int, len(r))
		s1 := make([]int, len(r))
		s2 := make([]int, len(r))

		copy(s, r)

		for i, n := range r {
			r1[i] = Number(n)
			r2[i] = Type(n)
		}

		sort.Ints(s)
		QuicksortNumberSlice(r1)
		QuicksortTypeSliceF(r2, func(a, b Type) bool { return a.(int) < b.(int) })

		for i := 0; i < len(s); i++ {
			s1[i] = int(r1[i])
			s2[i] = r2[i].(int)
		}

		if !reflect.DeepEqual(s1, s) {
			t.Logf("QuicksortNumberSlice:")
			t.Logf("%v !=", s1)
			t.Logf("%v", s)
			t.Fail()
		}

		if !reflect.DeepEqual(s2, s) {
			t.Logf("QuicksortTypeSliceF:")
			t.Logf("%v !=", s2)
			t.Logf("%v", s)
			t.Fail()
		}
	}
}
