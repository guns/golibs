package zero

import (
	"bytes"
	"unicode/utf8"
	"unsafe"
)

// ClearBytes zeroes a byte slice
func ClearBytes(bs []byte) {
	for i := range bs {
		bs[i] = 0
	}
}

// Copy of bytes.Buffer
type publicBuffer struct {
	buf       []byte            // contents are the bytes buf[off : len(buf)]
	off       int               // read at &buf[off], write at &buf[len(buf)]
	runeBytes [utf8.UTFMax]byte // avoid allocation of slice on each WriteByte or Rune
	bootstrap [64]byte          // memory to hold first slice; helps small buffers (Printf) avoid allocation.
	lastRead  int               // last read operation, so that Unread* can work correctly.
}

// ClearBuffer zeroes ALL data in a bytes.Buffer
func ClearBuffer(bbuf *bytes.Buffer) {
	b := (*publicBuffer)(unsafe.Pointer(bbuf))
	ClearBytes(b.buf)
	b.buf = b.buf[:0]
	b.off = 0
	for i := range b.runeBytes {
		b.runeBytes[i] = 0
	}
	for i := range b.bootstrap {
		b.bootstrap[i] = 0
	}
	b.lastRead = 0
}

// Copy of reflect.stringHeader
type stringHeader struct {
	Data unsafe.Pointer
	Len  int
}

// ClearString zeroes a string's backing array. This is truly s̶t̶u̶p̶i̶d̶ dangerous.
// Here are some considerations:
//	1. The string must be not be in the read-only data segment of the
//	   program (i.e. it must be dynamically allocated).
//	2. No one expects an immutable value to change, so expect data races
//	   if the string is shared.
func ClearString(s string) {
	hdr := *(*stringHeader)(unsafe.Pointer(&s))
	for i := 0; i < hdr.Len; i++ {
		*(*byte)(unsafe.Pointer(uintptr(hdr.Data) + uintptr(i))) = 0
	}
}
