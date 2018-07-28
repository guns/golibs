// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package optimized provides functions that are optimized on select
// architectures.
package optimized

import (
	"math/big"
	"math/bits"
)

// Mul returns the 128-bit unsigned product of a and b as a pair of 64-bit
// unsigned integers.
func Mul64(x, y uint64) (lower, upper uint64)

func mul64(x, y uint64) (lower, upper uint64) {
	a := new(big.Int).SetUint64(x)
	b := new(big.Int).SetUint64(y)
	p := new(big.Int).Mul(a, b)

	if p.IsUint64() {
		return p.Uint64(), 0
	}

	w := p.Bits()

	switch bits.UintSize {
	case 64:
		lower = uint64(w[0])
		upper = uint64(w[1])
	case 32:
		lower = uint64(w[1])<<32 + uint64(w[0])
		if len(w) == 3 {
			upper = uint64(w[2])
		} else {
			upper = uint64(w[3])<<32 + uint64(w[2])
		}
	}

	return lower, upper
}
