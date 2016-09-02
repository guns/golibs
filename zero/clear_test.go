package zero

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"
)

func TestClearBuffer(t *testing.T) {
	buf := &bytes.Buffer{}

	// Fill the internal bootstrap array
	for i := 0; i < 128; i++ {
		err := buf.WriteByte(0xff)
		if err != nil {
			t.Fail()
		}
	}

	// As well as the runeBytes array
	if _, err := buf.WriteRune('❤'); err != nil {
		t.Fail()
	}

	if _, err := buf.Read(make([]byte, 8)); err != nil {
		t.Fail()
	}

	ClearBuffer(buf)
	if !reflect.DeepEqual(buf, bytes.NewBuffer([]byte{})) {
		t.Errorf("%v != %v", buf, bytes.NewBuffer([]byte{}))
	}
}

func TestClearString(t *testing.T) {
	// Only dynamic strings can be mutated
	bs := make([]byte, 8)
	for i := range bs {
		bs[i] = byte(rand.Uint32() & 0x7f)
	}
	a := string(bs)
	b := a
	z := "\000\000\000\000\000\000\000\000"
	if !(a != z) {
		t.Errorf("expected: a != z, actual: %v == %v", a, z)
	}
	ClearString(a)
	if a != z {
		t.Errorf("%v != %v", a, z)
	}
	if b != z {
		t.Errorf("%v != %v", b, z)
	}
}
