// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import "math/bits"

// undefined is a sentinel value for the set of Vertex indices.
const undefined = -1

func fillUndefined(s []int) []int {
	for i := range s {
		s[i] = undefined
	}
	return s
}

func resizeIntSlice(s []int, size int) []int {
	if cap(s) >= size {
		return s[:size]
	}
	return make([]int, size, 1<<uint(bits.Len(uint(size-1))))
}
