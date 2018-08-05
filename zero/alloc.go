// Copyright (c) 2015-2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package zero provides utilities for working with sensitive byte slices that
// should be cleared before being garbage collected.
package zero

import (
	"bytes"

	"github.com/guns/golibs/calculate"
	"github.com/guns/golibs/optimized"
)

// Grow returns a byte slice that can accommodate n more bytes, and the index
// where bytes should be appended. If a reallocation is needed, old memory is
// zeroed to reduce leakage of sensitive data.
func Grow(bs []byte, n int) ([]byte, int) {
	newlen := len(bs) + n
	if cap(bs) >= newlen {
		return bs[:newlen], len(bs)
	}

	newcap := calculate.NextCap(newlen)
	if newcap < bytes.MinRead {
		newcap = bytes.MinRead
	}

	newslice := make([]byte, len(bs), newcap)
	copy(newslice, bs)
	optimized.MemsetByteSlice(bs, 0)
	return newslice[:newlen], len(bs)
}

// Append appends byte slices, but uses Grow for reallocation to reduce
// leakage of sensitive data.
func Append(dst []byte, src ...byte) []byte {
	dst, n := Grow(dst, len(src))
	copy(dst[n:], src)
	return dst
}
