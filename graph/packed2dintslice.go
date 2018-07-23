// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

// Packed2DIntSlice is a [][]int backed by a single []int.
// Out of order insertion and resizing is unsupported.
type Packed2DIntSlice [][]int

// MakePacked2DIntSlice returns a new Packed2DIntSlice of the given size.
func MakePacked2DIntSlice(size int) Packed2DIntSlice {
	return Packed2DIntSlice{make([]int, 0, size)}
}

// Append n to the last []int in p. Panics if there is no more space available.
func (p Packed2DIntSlice) Append(n int) Packed2DIntSlice {
	i := len(p) - 1
	s := p[i]
	s = s[:len(s)+1]
	s[len(s)-1] = n
	p[i] = s
	return p
}

// StartNewSlice appends a new []int to p.
func (p Packed2DIntSlice) StartNewSlice() Packed2DIntSlice {
	s := p[len(p)-1]                       // Last slice
	return append(p, s[len(s):cap(s)][:0]) // Start new slice after last slice
}

// Cap returns the number of ints p can hold.
func (p Packed2DIntSlice) Cap() int {
	return cap(p[0])
}
