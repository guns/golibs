// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

import (
	"reflect"
	"sort"
	"testing"
)

func TestQuicksortTypeSlice(t *testing.T) {
	for i := 0; i < 100; i++ {
		r := RandPersonSlice(10 * i)

		r2 := make([]Type, len(r))
		s1 := make(PersonSlice, len(r))
		s2 := make(PersonSlice, len(r))

		for i := range r {
			s1[i] = r[i]
			r2[i] = r[i]
		}

		sort.Sort(s1)
		QuicksortTypeSlice(r2, func(a, b *Type) bool {
			if (*a).(Person).name != (*b).(Person).name {
				return (*a).(Person).name < (*b).(Person).name
			}
			return (*a).(Person).age < (*b).(Person).age
		})

		for i := range r2 {
			s2[i] = r2[i].(Person)
		}

		if !reflect.DeepEqual(s2, s1) {
			t.Logf("QuicksortNumberSlice:")
			t.Logf("%v !=", s2)
			t.Logf("%v", s1)
			t.Fail()
		}
	}
}
