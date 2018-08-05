// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"errors"

	"github.com/guns/golibs/calculate"
	"github.com/guns/golibs/memset"
)

// undefined is a sentinel value for the set of Vertex indices.
const undefined = -1

var (
	errUnexpectedCycle = errors.New("graph contains unexpected cycle")
	errNoPath          = errors.New("graph does not contain the specified path")
)

func fillUndefined(s []int) []int {
	memset.Int(s, undefined)
	return s
}

func resizeIntSlice(s []int, size int) []int {
	if cap(s) >= size {
		return s[:size]
	}
	return make([]int, size, calculate.NextCap(size))
}
