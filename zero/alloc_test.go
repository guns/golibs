package zero

import (
	"reflect"
	"testing"
)

func TestGrow(t *testing.T) {
	bs := []byte("01234567")
	if cap(bs) != 8 {
		t.Errorf("%v != %v", cap(bs), 8)
	}
	newslice, n := Grow(bs, 5000)
	if n != 8 {
		t.Errorf("%v != %v", n, 8)
	}
	if cap(newslice) != 8192 {
		t.Errorf("%v != %v", cap(newslice), 8192)
	}
	if len(newslice) != 5008 {
		t.Errorf("%v != %v", len(newslice), 5008)
	}
	if !reflect.DeepEqual(newslice[:8], []byte("01234567")) {
		t.Errorf("%v != %v", newslice[:8], []byte("01234567"))
	}
	if !reflect.DeepEqual(bs, make([]byte, 8)) {
		t.Errorf("%v != %v", bs, make([]byte, 8))
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
}
