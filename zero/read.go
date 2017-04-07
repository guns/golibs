// Copyright (c) 2015-2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package zero

import (
	"bytes"
	"io"
)

// ReadAll reads data from r until EOF and appends it to buf, growing the
// buffer and zeroing old buffers as needed.
func ReadAll(buf []byte, r io.Reader) ([]byte, error) {
	for {
		i := len(buf)

		if cap(buf) == i {
			buf, i = Grow(buf, bytes.MinRead)
		}

		n, err := r.Read(buf[i:cap(buf)]) // Read to capacity
		buf = buf[:i+n]                   // Reslice to actual length

		switch err {
		case nil:
			// No errors, keep reading
		case io.EOF:
			return buf, nil
		default:
			ClearBytes(buf)
			return nil, err
		}
	}
}
