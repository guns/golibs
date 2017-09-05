// Copyright (c) 2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package editpipe provides a buffered pipe between a reader and a writer
// that allows basic editing operations.
package editpipe

import (
	"errors"
	"io"
	"unicode"

	"github.com/guns/golibs/errjoin"
)

// Op communicates an edit operation.
type Op byte

// Multiple Ops can be combined and each action will be executed in order.
//
//	Erase|Append|Flush|Close
//
// For example, the value above instructs an editpipe to erase the last byte,
// append the current byte, flush the buffer, and finally close the embedded
// writer.
const (
	Erase     Op = 1 << iota // Prune one byte at end of buffer
	EraseWord                // Prune one word (non-whitespace sequence) at end of buffer
	Kill                     // Truncate entire buffer
	Append                   // Append to buffer
	Flush                    // Write the buffer to the underlying file
	Close                    // Close the underlying file
)

var errClosed = errors.New("writer already closed")

// An EditFn specifies the semantics of an editpipe. The index parameter is
// the current index of the write buffer.
type EditFn func(index int, b byte) Op

// T implements buffering from a reader to a writer by reading one byte at a
// time from the reader and executing the Ops returned by an EditFn
type T struct {
	f      EditFn
	r      io.Reader      // Source
	w      io.WriteCloser // Destination
	rbuf   []byte         // Read buffer (length 1)
	wbuf   []byte         // Write buffer
	i      int            // Next write index
	secure bool           // Clear write buffer whenever appropriate
	done   bool           // Writer closed
}

// New returns a T with an internal write buffer buflen bytes long. If f is
// nil, BasicLineEdit is used as the edit function. If buflen is non-positive,
// a default buflen of 4096 is used.
//
// Unflushed input that exceeds buflen is flushed automatically to w.
//
// If secure is true, the write buffer is cleared when it is flushed, killed,
// or when w is closed.
//
// w is closed when a read from r returns an error or EOF.
func New(r io.Reader, w io.WriteCloser, secure bool, buflen int, f EditFn) *T {
	if buflen < 1 {
		buflen = 4096
	}
	if f == nil {
		f = BasicLineEdit
	}
	return &T{
		f:      f,
		r:      r,
		w:      w,
		rbuf:   make([]byte, 1),
		wbuf:   make([]byte, buflen),
		i:      0,
		secure: secure,
		done:   false,
	}
}

// append appends the byte to the write buffer. If this would overflow the
// buffer, the buffer is flushed first.
func (p *T) append(b byte) error {
	if p.i >= len(p.wbuf) {
		if err := p.flush(); err != nil {
			return err
		}
	}
	p.wbuf[p.i] = b
	p.i++
	return nil
}

// clear zeroes the write buffer.
func (p *T) clear() {
	for i := 0; i < p.i; i++ {
		p.wbuf[i] = 0
	}
}

// flush writes the write buffer to the embedded writer.
func (p *T) flush() error {
	if p.done {
		return errClosed
	}

	_, err := p.w.Write(p.wbuf[:p.i])

	if p.secure {
		p.clear()
	}

	p.i = 0
	return err
}

// erase prunes and zeroes the last byte in the write buffer.
func (p *T) erase() {
	if p.i <= 0 {
		return
	}
	p.i--
	p.wbuf[p.i] = 0 // Clearing this byte is faster than a branch
}

func isWordRune(r rune) bool {
	return r != ' ' && unicode.IsPrint(r)
}

// eraseWord prunes the last sequence of non-whitespace bytes in the write
// buffer. Multibyte character sequences are unsupported.
func (p *T) eraseWord() {
	if p.i <= 0 {
		return
	}

	// number of boundary transitions needed
	n := 2

	p.i--
	if isWordRune(rune(p.wbuf[p.i])) {
		n--
	}
	p.wbuf[p.i] = 0

	for p.i > 0 {
		p.i--
		isword := isWordRune(rune(p.wbuf[p.i]))
		if n == 2 && isword {
			n--
		} else if n == 1 && !isword {
			p.i++
			break
		}
		p.wbuf[p.i] = 0
	}
}

// kill truncates the write buffer.
func (p *T) kill() {
	if p.secure {
		p.clear()
	}
	p.i = 0
}

// close closes the embedded writer.
func (p *T) close() error {
	if p.done {
		return errClosed
	}
	p.done = true
	if p.secure {
		p.clear()
	}
	return p.w.Close()
}

// process reads and processes one byte from the embedded reader. On EOF and
// error the embedded writer is closed, and non-EOF errors are returned.
func (p *T) process() (ok bool, err error) {
	if p.done {
		return false, errClosed
	}

	n, err := p.r.Read(p.rbuf)
	if err == io.EOF {
		return false, p.close()
	} else if err != nil {
		return false, errjoin.Join("; ", err, p.close())
	} else if n <= 0 {
		return true, nil
	}

	b := p.rbuf[0]
	p.rbuf[0] = 0 // Clearing this byte is faster than a branch
	ops := p.f(p.i, b)

	if ops&Erase > 0 {
		p.erase()
	}
	if ops&EraseWord > 0 {
		p.eraseWord()
	}
	if ops&Kill > 0 {
		p.kill()
	}
	if ops&Append > 0 {
		if err := p.append(b); err != nil {
			return false, errjoin.Join("; ", err, p.close())
		}
	}
	if ops&Flush > 0 {
		if err := p.flush(); err != nil {
			return false, errjoin.Join("; ", err, p.close())
		}
	}
	if ops&Close > 0 {
		return false, p.close()
	}

	return true, nil
}

// ProcessAll reads and processes bytes from the embedded reader until EOF or
// error. The embedded writer is closed, and non-EOF errors are returned.
func (p *T) ProcessAll() error {
	for {
		if ok, err := p.process(); err != nil {
			return err
		} else if !ok {
			return nil
		}
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
