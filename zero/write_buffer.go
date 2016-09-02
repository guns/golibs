// This file contains a reimplementation of bytes.Buffer, which is
// Copyright 2009 The Go Authors. All rights reserved.

package zero

import (
	"bytes"
	"io"
	"unicode/utf8"
)

/*

A WriteBuffer implements the same write methods as a bytes.Buffer, but zeroes
old buffers on reallocation.

*/
type WriteBuffer struct {
	buf []byte
}

// NewWriteBuffer wraps a byte slice, returning a *WriteBuffer.
func NewWriteBuffer(buf []byte) *WriteBuffer {
	return &WriteBuffer{buf: buf}
}

// Bytes returns the underlying byte slice of a WriteBuffer.
func (w *WriteBuffer) Bytes() []byte {
	return w.buf
}

// Len returns the length of (*WriteBuffer) Bytes()
func (w *WriteBuffer) Len() int { return len(w.buf) }

// Truncate zeroes and discards all but the first n bytes from the buffer. It
// panics if n is negative or greater than the length of the buffer.
func (w *WriteBuffer) Truncate(n int) {
	if n < 0 || n > len(w.buf) {
		panic("zero.WriteBuffer: truncation out of range")
	}
	memset(w.buf[n:], 0)
	w.buf = w.buf[:n]
}

// Reset zeros and resets the buffer so it has no content.
// w.Reset() is the same as w.Truncate(0).
func (w *WriteBuffer) Reset() { w.Truncate(0) }

// Grow grows the buffer's capacity, if necessary, to guarantee space for
// another n bytes. If a reallocation is needed, old memory is zeroed to
// reduce leakage of sensitive data. After Grow(n), at least n bytes can be
// written to the buffer without another allocation.
//
// If n is negative, Grow will panic.
func (w *WriteBuffer) Grow(n int) {
	if n < 0 {
		panic("(*zero.WriteBuffer) Grow: negative count")
	}
	w.buf, _ = Grow(w.buf, n)
}

// Write appends the contents of p to the buffer, growing the buffer and
// zeroing old buffers as needed. The int return value is the length of p; no
// errors are returned.
func (w *WriteBuffer) Write(p []byte) (int, error) {
	var i int
	w.buf, i = Grow(w.buf, len(p))
	return copy(w.buf[i:], p), nil
}

// WriteString appends the contents of s to the buffer, growing the buffer and
// zeroing old buffers as needed. The int return value is the length of s; no
// errors are returned.
func (w *WriteBuffer) WriteString(s string) (int, error) {
	var i int
	w.buf, i = Grow(w.buf, len(s))
	return copy(w.buf[i:], s), nil
}

// ReadFrom reads data from r until EOF and appends it to the buffer, growing
// the buffer and zeroing old buffers as needed. The int64 return value is the
// number of bytes read. Any error except io.EOF encountered during the read
// is also returned.
func (w *WriteBuffer) ReadFrom(r io.Reader) (int64, error) {
	var total int64
	for {
		// Reslice buffer while we adjust its length
		buf := w.buf
		i := len(buf)

		if cap(buf) == i {
			newcap := cap(buf) * 2
			if newcap == 0 {
				newcap = bytes.MinRead
			}
			buf, i = Grow(buf, newcap)
		}

		n, err := r.Read(buf[i:cap(buf)]) // Read to capacity
		total += int64(n)
		w.buf = buf[:i+n] // Cap to actual length

		switch err {
		case nil:
			// No errors, keep reading
		case io.EOF:
			return total, nil
		default:
			return total, err
		}
	}
}

// WriteByte appends a byte to the buffer, growing the buffer and zeroing old
// buffers as needed. The returned error is always nil, but is included to
// match bufio.Writer's WriteByte.
func (w *WriteBuffer) WriteByte(b byte) error {
	var i int
	w.buf, i = Grow(w.buf, 1)
	w.buf[i] = b
	return nil
}

// WriteRune appends the UTF-8 encoding of a Unicode code point to the buffer,
// returning its length and an error, which is always nil but is included to
// match bufio.Writer's WriteRune. The buffer is grown and old buffers are
// zeroed as needed.
func (w *WriteBuffer) WriteRune(r rune) (int, error) {
	if r < utf8.RuneSelf {
		_ = w.WriteByte(byte(r)) // never fails
		return 1, nil
	}
	bs, i := Grow(w.buf, utf8.UTFMax)
	n := utf8.EncodeRune(bs[i:], r)
	w.buf = bs[:i+n]
	return n, nil
}
