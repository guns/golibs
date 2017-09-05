// Copyright (c) 2015-2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package unsafezero

import (
	"bytes"
	"errors"
	"runtime"
	"unicode/utf8"
	"unsafe"

	"github.com/guns/golibs/zero"
)

const (
	goVersionUnknown = iota
	go1_0
	go1_8
	go1_9
)

var goBufferVersion = goVersionUnknown

func init() {
	v := runtime.Version()
	if len(v) < 5 {
		if v == "go1" {
			goBufferVersion = go1_0
		}
		return
	}
	switch v[:5] {
	case "go1.8":
		goBufferVersion = go1_8
	case "go1.9":
		goBufferVersion = go1_9
	}
}

type go1_0BytesBuffer struct {
	buf       []byte            // contents are the bytes buf[off : len(buf)]
	off       int               // read at &buf[off], write at &buf[len(buf)]
	runeBytes [utf8.UTFMax]byte // avoid allocation of slice on each WriteByte or Rune
	bootstrap [64]byte          // memory to hold first slice; helps small buffers (Printf) avoid allocation.
	lastRead  int               // last read operation, so that Unread* can work correctly.
}

type go1_8BytesBuffer struct {
	buf       []byte   // contents are the bytes buf[off : len(buf)]
	off       int      // read at &buf[off], write at &buf[len(buf)]
	bootstrap [64]byte // memory to hold first slice; helps small buffers (Printf) avoid allocation.
	lastRead  int      // last read operation, so that Unread* can work correctly.
}

type go1_9BytesBuffer struct {
	buf       []byte   // contents are the bytes buf[off : len(buf)]
	off       int      // read at &buf[off], write at &buf[len(buf)]
	lastRead  int      // last read operation, so that Unread* can work correctly.
	bootstrap [64]byte // memory to hold first slice; helps small buffers avoid allocation.
}

// ClearBuffer zeroes ALL data in a bytes.Buffer
func ClearBuffer(bbuf *bytes.Buffer) {
	switch goBufferVersion {
	case go1_9:
		b := (*go1_9BytesBuffer)(unsafe.Pointer(bbuf))
		zero.ClearBytes(b.buf)
		b.buf = b.buf[:0]
		b.off = 0
		for i := range b.bootstrap {
			b.bootstrap[i] = 0
		}
		b.lastRead = 0
	case go1_8:
		b := (*go1_8BytesBuffer)(unsafe.Pointer(bbuf))
		zero.ClearBytes(b.buf)
		b.buf = b.buf[:0]
		b.off = 0
		for i := range b.bootstrap {
			b.bootstrap[i] = 0
		}
		b.lastRead = 0
	case go1_0:
		b := (*go1_0BytesBuffer)(unsafe.Pointer(bbuf))
		zero.ClearBytes(b.buf)
		b.buf = b.buf[:0]
		b.off = 0
		for i := range b.runeBytes {
			b.runeBytes[i] = 0
		}
		for i := range b.bootstrap {
			b.bootstrap[i] = 0
		}
		b.lastRead = 0
	default:
		panic(errors.New("unable to determine go version in ClearBuffer"))
	}
}
