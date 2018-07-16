package bitslice

import (
	"reflect"
	"testing"
)

func TestBitSlice(t *testing.T) {
	data := []struct {
		bs      T
		n       int
		get     bool
		set     T
		clear   T
		toggle  T
		cas     bool
		cac     bool
		catarg  bool
		cat     bool
		offsets []int
	}{
		{
			bs:      T{0x7},
			n:       1,
			get:     true,
			set:     T{0x7},
			clear:   T{0x5},
			toggle:  T{0x5},
			cas:     false,
			cac:     true,
			catarg:  true,
			cat:     true,
			offsets: []int{0, 1, 2},
		},
		{
			bs:      T{0x0, 0x1},
			n:       64,
			get:     true,
			set:     T{0x0, 0x1},
			clear:   T{0x0, 0x0},
			toggle:  T{0x0, 0x0},
			cas:     false,
			cac:     true,
			catarg:  false,
			cat:     false,
			offsets: []int{64},
		},
		{
			bs:      T{0x0, 0x7},
			n:       65,
			get:     true,
			set:     T{0x0, 0x7},
			clear:   T{0x0, 0x5},
			toggle:  T{0x0, 0x5},
			cas:     false,
			cac:     true,
			catarg:  true,
			cat:     true,
			offsets: []int{64, 65, 66},
		},
		{
			bs:      T{0xffff, 0x5},
			n:       65,
			get:     false,
			set:     T{0xffff, 0x7},
			clear:   T{0xffff, 0x5},
			toggle:  T{0xffff, 0x7},
			cas:     true,
			cac:     false,
			catarg:  false,
			cat:     true,
			offsets: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 64, 66},
		},
	}

	for i, row := range data {
		n := row.n
		bs := Make(len(row.bs) * 64)
		copy(bs, row.bs)

		if bs.Get(n) != row.get {
			t.Errorf("[%v] expected: bs.Get(%v) == %v", i, n, row.get)
		}

		bs.Set(row.n)
		if !reflect.DeepEqual(bs, row.set) {
			t.Logf("[%v] bs.Set(%v):", i, n)
			t.Logf("\t%064b !=", bs)
			t.Logf("\t%064b", row.set)
			t.Fail()
		}

		copy(bs, row.bs)
		bs.Clear(n)
		if !reflect.DeepEqual(bs, row.clear) {
			t.Logf("[%v] bs.Clear(%v):", i, n)
			t.Logf("\t%064b !=", bs)
			t.Logf("\t%064b", row.set)
			t.Fail()
		}

		copy(bs, row.bs)
		bs.Toggle(n)
		if !reflect.DeepEqual(bs, row.toggle) {
			t.Logf("[%v] bs.Toggle(%v):", i, n)
			t.Logf("\t%064b !=", bs)
			t.Logf("\t%064b", row.set)
			t.Fail()
		}

		copy(bs, row.bs)
		if bs.CompareAndSet(n) != row.cas {
			t.Errorf("[%v] expected: bs.CompareAndSet(%v) == %v", i, n, row.cas)
		} else if !reflect.DeepEqual(bs, row.set) {
			t.Logf("[%v] bs.CompareAndSet(%v):", i, n)
			t.Logf("\t%064b !=", bs)
			t.Logf("\t%064b", row.set)
			t.Fail()
		}

		copy(bs, row.bs)
		if bs.CompareAndClear(n) != row.cac {
			t.Errorf("[%v] expected: bs.CompareAndClear(%v) == %v", i, n, row.cac)
		} else if !reflect.DeepEqual(bs, row.clear) {
			t.Logf("[%v] bs.CompareAndClear(%v):", i, n)
			t.Logf("\t%064b !=", bs)
			t.Logf("\t%064b", row.clear)
			t.Fail()
		}

		copy(bs, row.bs)
		if bs.CompareAndToggle(n, row.catarg) != row.cat {
			t.Errorf("[%v] expected: bs.CompareAndToggle(%v, %v) == %v", i, n, row.catarg, row.cat)
		} else if row.cat && !reflect.DeepEqual(bs, row.toggle) {
			t.Logf("[%v] bs.CompareAndToggle(%v, %v):", i, n, row.catarg)
			t.Logf("\t%064b !=", bs)
			t.Logf("\t%064b", row.toggle)
			t.Fail()
		} else if !row.cat && !reflect.DeepEqual(bs, row.bs) {
			t.Logf("[%v] bs.CompareAndToggle(%v, %v):", i, n, row.catarg)
			t.Logf("\t%064b !=", bs)
			t.Logf("\t%064b", row.bs)
			t.Fail()
		}

		copy(bs, row.bs)
		v := make([]int, 0, len(row.offsets))
		if !reflect.DeepEqual(bs.AppendOffsets(v), row.offsets) {
			t.Errorf("%v != %v", bs.AppendOffsets(v), row.offsets)
		}
		if bs.Popcnt() != len(row.offsets) {
			t.Errorf("%v != %v", bs.Popcnt(), len(row.offsets))
		}

		copy(bs, row.bs)
		bs.Reset()
		if !reflect.DeepEqual(bs, Make(len(bs)*64)) {
			t.Errorf("%v != %v", bs, Make(len(bs)*64))
		}
	}
}
