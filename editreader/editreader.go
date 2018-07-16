// Copyright (c) 2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package editreader provides a buffered reader that enables basic editing.
package editreader

import (
	"io"
	"unicode"
)

// An Op communicates an edit operation.
type Op byte

// When multiple Ops are combined, each action will be executed in
// least-significant-order.
//
//	Erase | Append | Flush | Close
//
// For example, the value above instructs an editreader to erase the last byte
// in the buffer, append the current byte, flush the buffer, and finally close.
const (
	Erase     Op = 1 << iota // Clear and remove one byte at end of buffer
	EraseWord                // Clear and remove non-whitespace sequence at end of buffer
	Kill                     // Truncate entire buffer
	Append                   // Append to buffer
	Flush                    // Flush the buffer to readers
	Close                    // Close this editreader and send EOF to future readers
)

// An EditFn specifies the semantics of an editreader. i is the index of the
// buffer where b would be appended.
type EditFn func(i int, b byte) Op

// T wraps a Reader and implements an editable buffer by reading one byte at a
// time from the wrapped Reader and executing the Ops returned by an EditFn.
type T struct {
	editFn    EditFn
	src       io.Reader // Source
	rbuf      []byte    // 1-byte buffer for reading from r
	buf       []byte    // Edit buffer
	ridx      int       // Next edit buffer read index
	widx      int       // Next edit buffer write index
	err       error     // Error to return readers after close
	available bool      // Set on Flush, cleared after full read
	overflow  bool      // Set when an Append would overflow the edit buffer
	done      bool      // Set on read error/EOF or Close Op
	secure    bool      // Clear edit buffer when appropriate
}

// New returns an editreader with an buffer buflen bytes long. If f is nil,
// BasicLineEdit is used as the edit function. If buflen is non-positive, a
// default buflen of 4096 is used.
//
// Unflushed input that exceeds buflen is flushed automatically.
//
// If secure is true, the buffer is cleared on Kill, Close, and after reads.
// Note that input that is flushed but unread is not cleared on Close. Since
// the data will be cleared on read, callers should make sure to drain the
// editreader if they care about zeroing buffers. Erase and EraseWord always
// zero the bytes they erase.
func New(r io.Reader, buflen int, secure bool, f EditFn) *T {
	if buflen < 1 {
		buflen = 4096
	}
	if f == nil {
		f = BasicLineEdit
	}
	buf := make([]byte, buflen+1)
	return &T{
		editFn: f,
		src:    r,
		rbuf:   buf[:1],
		buf:    buf[1:],
		secure: secure,
	}
}

// Read reads from the edit buffer if it is available for read. If there is no
// data available, data is read and processed from the source reader until it
// is flushed and available for read.
func (e *T) Read(dst []byte) (n int, err error) {
	for {
		if e.available {
			return e.readAvailable(dst)
		} else if e.done {
			return 0, e.err
		}
		e.scan()
	}
}

// WriteTo implements WriterTo, and has been explicitly provided to avoid use
// of a transfer buffer in io.Copy, which is called by exec.Cmd to pipe data
// from a non-file stdin.
func (e *T) WriteTo(w io.Writer) (n int64, err error) {
	buf := make([]byte, 4096)
	for {
		i, rerr := e.Read(buf)
		if i > 0 {
			j, werr := w.Write(buf[:i])
			n += int64(j)
			if werr != nil {
				err = werr
				break
			}
		}
		if rerr != nil {
			if rerr != io.EOF {
				err = rerr
			}
			break
		}
	}
	if e.secure {
		clearbytes(buf)
	}
	return n, err
}

// readAvailable copies unread data into dst. The available flag is cleared if
// no unread data remains.
// WARNING: This method assumes the buffer is available for read!
func (e *T) readAvailable(dst []byte) (n int, err error) {
	n = copy(dst, e.buf[e.ridx:e.widx])
	i := e.ridx + n
	if e.secure {
		clearbytes(e.buf[e.ridx:i])
	}
	e.ridx = i
	if e.ridx >= e.widx {
		e.ridx = 0
		e.widx = 0
		e.available = false
	}
	return n, nil
}

// scan reads and processes one byte from the source reader. If there was an
// overflow from the last scan, no read occurs and the previously read byte is
// processed instead.
func (e *T) scan() {
	if e.overflow {
		e.overflow = false
		e.process(e.rbuf[0])
		return
	}

	n, err := e.src.Read(e.rbuf)
	if n > 0 {
		e.process(e.rbuf[0])
	}
	if err != nil {
		e.closeWithError(err)
	}
}

// process executes the Ops specified by the EditFn for byte b.
func (e *T) process(b byte) {
	ops := e.editFn(e.widx, b)

	if ops&Erase > 0 {
		e.erase()
	}
	if ops&EraseWord > 0 {
		e.eraseWord()
	}
	if ops&Kill > 0 {
		e.kill()
	}
	if ops&Append > 0 {
		e.append(b)
	}
	if ops&Flush > 0 {
		e.available = true
	}
	if ops&Close > 0 {
		e.closeWithError(nil)
	}
}

// erase prunes and zeroes the last byte in the edit buffer.
func (e *T) erase() {
	if e.widx <= 0 {
		return
	}
	e.widx--
	e.buf[e.widx] = 0
}

// eraseWord prunes the last sequence of non-whitespace bytes in the write
// buffer. Multibyte character sequences are unsupported.
func (e *T) eraseWord() {
	if e.widx <= 0 {
		return
	}

	// number of boundary transitions
	n := 2

	e.widx--
	if isWordRune(rune(e.buf[e.widx])) {
		n--
	}
	e.buf[e.widx] = 0

	for e.widx > 0 {
		e.widx--
		isword := isWordRune(rune(e.buf[e.widx]))
		if n == 2 && isword {
			n--
		} else if n == 1 && !isword {
			e.widx++
			break
		}
		e.buf[e.widx] = 0
	}
}

// kill truncates the edit buffer
func (e *T) kill() {
	if e.secure {
		clearbytes(e.buf[:e.widx])
	}
	e.widx = 0
}

// append appends a byte to the edit buffer. If this would overflow the
// buffer, the overflow flag is set and the buffer is marked as available.
func (e *T) append(b byte) {
	if e.widx >= len(e.buf) {
		e.overflow = true
		e.available = true
		return
	}
	e.buf[e.widx] = b
	e.widx++
}

// closeWithError marks e as done, and sets err as the error to send to readers. If err
// is nil, io.EOF is sent.
func (e *T) closeWithError(err error) {
	if e.done {
		return
	}
	e.done = true
	e.err = err
	if err == nil {
		e.err = io.EOF
	}
	if e.secure {
		e.rbuf[0] = 0
		n := e.widx
		if e.available {
			n = e.ridx
		}
		clearbytes(e.buf[:n])
	}
}

// BasicLineEdit specifies a simple line editor.
func BasicLineEdit(i int, b byte) Op {
	switch b {
	case '\b', 0x7f: // ^H, ^?
		return Erase
	case 0x17: // ^W
		return EraseWord
	case 0x15: // ^U
		return Kill
	case 0x04: // ^D
		if i == 0 {
			return Close
		}
		return Flush
	case '\n':
		return Append | Flush
	default:
		return Append
	}
}

// clearbytes zeros a byte slice.
func clearbytes(bs []byte) {
	for i := range bs {
		bs[i] = 0
	}
}

func isWordRune(r rune) bool {
	return r != ' ' && unicode.IsPrint(r)
}
