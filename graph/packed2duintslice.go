// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

// Packed2DUintSlice is a [][]uint backed by a single []uint.
// Out of order insertion and resizing is unsupported.
type Packed2DUintSlice [][]uint

// MakePacked2DUintSlice returns a new Packed2DUintSlice of the given size.
func MakePacked2DUintSlice(size int) Packed2DUintSlice {
	return Packed2DUintSlice{make([]uint, 0, size)}
}

// Append n to the last []uint in p. Panics if there is no more space available.
func (p Packed2DUintSlice) Append(n uint) Packed2DUintSlice {
	i := len(p) - 1
	s := p[i]
	s = s[:len(s)+1]
	s[len(s)-1] = n
	p[i] = s
	return p
}

// StartNewSlice appends a new []uint to p.
func (p Packed2DUintSlice) StartNewSlice() Packed2DUintSlice {
	s := p[len(p)-1]                       // Last slice
	return append(p, s[len(s):cap(s)][:0]) // Start new slice after last slice
}

// Cap returns the number of ints p can hold.
func (p Packed2DUintSlice) Cap() int {
	return cap(p[0])
}
