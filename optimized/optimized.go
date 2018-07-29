// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package optimized provides functions that are optimized on select
// architectures.
package optimized

import "math/rand"

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

// RandInt63n is an optimized version of math/rand.Int63n that implements
// Daniel Lemire's multiplicative alternative to modulo reduction:
//
//	https://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/
//	https://lemire.me/blog/2016/06/30/fast-random-shuffling/
//
// cf. math/rand.int31n
func RandInt63n(n int64) int64 {
	if n <= 0 {
		return n / (n * 0) // panic() defeats inlining [go1.11]
	}

	lo, hi := Mul64(rand.Uint64(), uint64(n))

	if lo < uint64(n) {
		threshold := uint64(-n) % uint64(n)
		for lo < threshold {
			lo, hi = Mul64(rand.Uint64(), uint64(n))
		}
	}

	const lower63 = 0x7fffffffffffffff

	return int64(hi & lower63)
}
