// Copyright (c) 2015-2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package zero provides utilities for working with sensitive byte slices that
// should be cleared before being garbage collected.
package zero

// ClearBytes zeroes a byte slice
func ClearBytes(bs []byte) {
	for i := range bs {
		bs[i] = 0
	}
}
