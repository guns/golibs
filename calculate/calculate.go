// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package calculate provides functions for common computations like
// determining the capacity of new slice.
package calculate

import "math/bits"

// NextCap returns a recommended capacity for allocating a slice of a given
// size that is expected to grow in the future.
func NextCap(size int) int {
	switch {
	case size <= 8:
		return 8
	case size <= 4096: // Round up to next power of two
		return 1 << uint(bits.Len(uint(size-1)))
	case size <= 1<<20: // Round up to next half of current power of two
		shift := uint(bits.Len(uint(size-1))) - 2
		return (((size - 1) >> shift) + 1) << shift
	default: // Round up to next quarter of current power of two
		shift := uint(bits.Len(uint(size-1))) - 3
		return (((size - 1) >> shift) + 1) << shift
	}
}
