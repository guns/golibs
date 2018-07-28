// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package optimized provides functions that are optimized on select
// architectures.
package optimized

// Mul returns the 128-bit unsigned product of a and b as a pair of 64-bit
// unsigned integers.
func Mul64(x, y uint64) (lower, upper uint64)

func mul64(x, y uint64) (lower, upper uint64) {
	const lo = 0xffffffff

	a, b := x>>32, x&lo
	c, d := y>>32, y&lo

	ac := a * c
	ad := a * d
	bc := b * c
	bd := b * d

	mid := ad&lo + bc&lo + bd>>32 // 34 bits

	lower = (mid&lo)<<32 + bd&lo
	upper = ac + ad>>32 + bc>>32 + mid>>32

	return lower, upper
}
