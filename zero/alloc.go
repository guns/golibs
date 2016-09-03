package zero

// Grow returns a byte slice that can accommodate n more bytes, and the index
// where bytes should be appended. If a reallocation is needed, old memory is
// zeroed to reduce leakage of sensitive data.
func Grow(bs []byte, n int) ([]byte, int) {
	newlen := len(bs) + n
	if cap(bs) >= newlen {
		return bs[:newlen], len(bs)
	}
	var newcap int
	if newlen < 1024 {
		newcap = 1024
	} else if newlen < 2048 {
		newcap = 2048
	} else {
		newcap = newlen + 4096 - (newlen % 4096)
	}
	newslice := make([]byte, len(bs), newcap)
	copy(newslice, bs)
	ClearBytes(bs)
	return newslice[:newlen], len(bs)
}

// Append appends byte slices, but uses Grow for reallocation to reduce
// leakage of sensitive data.
func Append(dst []byte, src ...byte) []byte {
	dst, n := Grow(dst, len(src))
	copy(dst[n:], src)
	return dst
}
