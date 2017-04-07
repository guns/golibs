// Copyright (c) 2015-2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package zero

import (
	"reflect"
	"testing"
)

func TestGrow(t *testing.T) {
	data := []struct {
		len, cap               int
		data                   string
		growth, newlen, newcap int
	}{
		{8, 8, "01234567", +8, 16, 512},
		{8, 16, "01234567", +8, 16, 16},
		{0, 0, "", +600, 600, 1024},
		{446, 512, "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			446 * 3, 446 * 4, 2048},
	}

	for _, row := range data {
		bs := make([]byte, row.len, row.cap)
		copy(bs, row.data)

		newslice, n := Grow(bs, row.growth)

		if n != row.len {
			t.Errorf("%v != %v", n, row.len)
		}
		if cap(newslice) != row.newcap {
			t.Errorf("%v != %v", cap(newslice), row.newcap)
		}
		if len(newslice) != row.newlen {
			t.Errorf("%v != %v", len(newslice), row.newlen)
		}
		if !reflect.DeepEqual(newslice[:row.len], []byte(row.data)) {
			t.Errorf("%v != %v", newslice[:row.len], []byte(row.data))
		}
		if row.cap == row.newcap {
			if !reflect.DeepEqual(bs, []byte(row.data)) {
				t.Errorf("%v != %v", bs, []byte(row.data))
			}
		} else {
			if !reflect.DeepEqual(bs, make([]byte, row.len, row.cap)) {
				t.Errorf("%v != %v", bs, make([]byte, row.len, row.cap))
			}
		}
	}
}

func TestAppend(t *testing.T) {
	a := []byte("01234567")
	b := Append(a, '8', '9', 'a')
	b = Append(b, []byte("bcdef")...)
	if !reflect.DeepEqual(a, make([]byte, 8)) {
		t.Errorf("%v != %v", a, make([]byte, 8))
	}
	if !reflect.DeepEqual(b, []byte("0123456789abcdef")) {
		t.Errorf("%v != %v", b, []byte("0123456789abcdef"))
	}

	// Test append of nil slice
	if !reflect.DeepEqual(Append(nil, []byte("01234567")...), []byte("01234567")) {
		t.Errorf("%v != %v", Append(nil, []byte("01234567")...), []byte("01234567"))
	}
}
