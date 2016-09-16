package zero

import (
	"bytes"
	"io"
)

// ReadAllInto reads data from r until EOF and appends it to buf, growing the
// buffer and zeroing old buffers as needed.
func ReadAllInto(buf []byte, r io.Reader) ([]byte, error) {
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
