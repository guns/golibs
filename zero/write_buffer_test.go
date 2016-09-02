package zero

import (
	"reflect"
	"strings"
	"testing"
)

func TestWriteBufferWrite(t *testing.T) {
	w := NewWriteBuffer(make([]byte, 0, 8))
	n, err := w.Write([]byte("01234567"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 8 {
		t.Errorf("%v != %v", n, 8)
	}
	if cap(w.buf) != 8 {
		t.Errorf("%v != %v", cap(w.buf), 8)
	}
	if !reflect.DeepEqual(w.buf, []byte("01234567")) {
		t.Errorf("%v != %v", w.buf, []byte("01234567"))
	}

	bs := w.buf
	n, err = w.Write([]byte("89abcdef"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 8 {
		t.Errorf("%v != %v", n, 8)
	}
	if cap(w.Bytes()) != 1024 {
		t.Errorf("%v != %v", cap(w.Bytes()), 1024)
	}
	if !reflect.DeepEqual(w.Bytes(), []byte("0123456789abcdef")) {
		t.Errorf("%v != %v", w.Bytes(), []byte("0123456789abcdef"))
	}
	if w.Len() != 16 {
		t.Errorf("%v != %v", w.Len(), 16)
	}
	if !reflect.DeepEqual(bs, make([]byte, 8)) {
		t.Errorf("%v != %v", bs, make([]byte, 8))
	}

	bs = make([]byte, 8)
	copy(bs, "01234567")
	w = NewWriteBuffer(bs)
	n, err = w.WriteString("89abcdef")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 8 {
		t.Errorf("%v != %v", n, 8)
	}
	if cap(w.Bytes()) != 1024 {
		t.Errorf("%v != %v", cap(w.Bytes()), 1024)
	}
	if !reflect.DeepEqual(w.Bytes(), []byte("0123456789abcdef")) {
		t.Errorf("%v != %v", w.Bytes(), []byte("0123456789abcdef"))
	}
	if w.Len() != 16 {
		t.Errorf("%v != %v", w.Len(), 16)
	}
	if !reflect.DeepEqual(bs, make([]byte, 8)) {
		t.Errorf("%v != %v", bs, make([]byte, 8))
	}

	bs = make([]byte, 8)
	copy(bs, "01234567")
	w = NewWriteBuffer(bs)
	err = w.WriteByte('8')
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if cap(w.Bytes()) != 1024 {
		t.Errorf("%v != %v", cap(w.Bytes()), 1024)
	}
	if !reflect.DeepEqual(w.Bytes(), []byte("012345678")) {
		t.Errorf("%v != %v", w.Bytes(), []byte("012345678"))
	}
	if w.Len() != 9 {
		t.Errorf("%v != %v", w.Len(), 9)
	}
	if !reflect.DeepEqual(bs, make([]byte, 8)) {
		t.Errorf("%v != %v", bs, make([]byte, 8))
	}

	bs = make([]byte, 8)
	copy(bs, "01234567")
	w = NewWriteBuffer(bs)
	n, err = w.WriteRune('❤')
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("%v != %v", n, 3)
	}
	if cap(w.Bytes()) != 1024 {
		t.Errorf("%v != %v", cap(w.Bytes()), 1024)
	}
	if !reflect.DeepEqual(w.Bytes(), []byte("01234567❤")) {
		t.Errorf("%v != %v", w.Bytes(), []byte("01234567❤"))
	}
	if w.Len() != 11 {
		t.Errorf("%v != %v", w.Len(), 11)
	}
	if !reflect.DeepEqual(bs, make([]byte, 8)) {
		t.Errorf("%v != %v", bs, make([]byte, 8))
	}

	bs = make([]byte, 8)
	copy(bs, "01234567")
	w = NewWriteBuffer(bs)
	n, err = w.WriteRune('8')
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("%v != %v", n, 1)
	}
	if cap(w.Bytes()) != 1024 {
		t.Errorf("%v != %v", cap(w.Bytes()), 1024)
	}
	if !reflect.DeepEqual(w.Bytes(), []byte("012345678")) {
		t.Errorf("%v != %v", w.Bytes(), []byte("012345678"))
	}
	if w.Len() != 9 {
		t.Errorf("%v != %v", w.Len(), 9)
	}
	if !reflect.DeepEqual(bs, make([]byte, 8)) {
		t.Errorf("%v != %v", bs, make([]byte, 8))
	}
}

func TestTruncateWriteBuffer(t *testing.T) {
	bs := make([]byte, 16)
	copy(bs, "0123456789abcdef")
	w := NewWriteBuffer(bs)
	w.Truncate(8)
	if cap(w.Bytes()) != 16 {
		t.Errorf("%v != %v", cap(w.Bytes()), 16)
	}
	if !reflect.DeepEqual(w.Bytes(), []byte("01234567")) {
		t.Errorf("%v != %v", w.Bytes(), []byte("01234567"))
	}
	if !reflect.DeepEqual(bs[:16], []byte("01234567\x00\x00\x00\x00\x00\x00\x00\x00")) {
		t.Errorf("%v != %v", bs[:16], []byte("01234567\x00\x00\x00\x00\x00\x00\x00\x00"))
	}
	if !reflect.DeepEqual(w.Bytes()[:16], []byte("01234567\x00\x00\x00\x00\x00\x00\x00\x00")) {
		t.Errorf("%v != %v", w.Bytes()[:16], []byte("01234567\x00\x00\x00\x00\x00\x00\x00\x00"))
	}

	w.Reset()
	if cap(w.Bytes()) != 16 {
		t.Errorf("%v != %v", cap(w.Bytes()), 16)
	}
	if !reflect.DeepEqual(w.Bytes(), []byte{}) {
		t.Errorf("%v != %v", w.Bytes(), []byte{})
	}
	if !reflect.DeepEqual(bs[:16], []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")) {
		t.Errorf("%v != %v", bs[:16], []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"))
	}
	if !reflect.DeepEqual(w.Bytes()[:16], []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")) {
		t.Errorf("%v != %v", w.Bytes()[:16], []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"))
	}
}

func TestGrowWriteBuffer(t *testing.T) {
	bs := make([]byte, 8)
	copy(bs, "01234567")
	w := NewWriteBuffer(bs)
	w.Grow(8)
	if !reflect.DeepEqual(bs, []byte("\x00\x00\x00\x00\x00\x00\x00\x00")) {
		t.Errorf("%v != %v", bs, []byte("\x00\x00\x00\x00\x00\x00\x00\x00"))
	}
	if w.Len() != 16 {
		t.Errorf("%v != %v", w.Len(), 16)
	}
	if !reflect.DeepEqual(w.Bytes(), []byte("01234567\x00\x00\x00\x00\x00\x00\x00\x00")) {
		t.Errorf("%v != %v", w.Bytes(), []byte("01234567\x00\x00\x00\x00\x00\x00\x00\x00"))
	}
}

func TestWriteBufferReadFrom(t *testing.T) {
	w := new(WriteBuffer)
	lorem := `Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.`
	n, err := w.ReadFrom(strings.NewReader(lorem))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != int64(len(lorem)) {
		t.Errorf("%v != %v", n, int64(len(lorem)))
	}
	if !reflect.DeepEqual(w.Bytes(), []byte(lorem)) {
		t.Errorf("%v != %v", w.Bytes(), []byte(lorem))
	}
}
