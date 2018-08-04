// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package optimized provides functions that are optimized on select
// architectures.
package optimized

import "math/rand"

// RandIntn is an optimized version of math/rand.Intn that implements
// Daniel Lemire's multiplicative alternative to modulo reduction:
//
//	https://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/
//	https://lemire.me/blog/2016/06/30/fast-random-shuffling/
//
func RandIntn(n int) int
func randIntn(n int) int {
	if uintSize == 32 {
		return int(RandInt31n(int32(n)))
	}
	return int(RandInt63n(int64(n)))
}

// RandInt63n is an optimized version of math/rand.Int63n.
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

// RandInt31n is an optimized version of math/rand.Int31n.
func RandInt31n(n int32) int32 {
	if n <= 0 {
		return n / (n * 0) // panic() defeats inlining [go1.11]
	}

	lo, hi := Mul32(rand.Uint32(), uint32(n))

	if lo < uint32(n) {
		threshold := uint32(-n) % uint32(n)
		for lo < threshold {
			lo, hi = Mul32(rand.Uint32(), uint32(n))
		}
	}

	const lower31 = 0x7fffffff

	return int32(hi & lower31)
}
