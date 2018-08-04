// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package optimized

const uintSize = 32 << ((^uint(0) >> 32) & 1)

// Mul returns the unsigned product of a and b as a pair of unsigned integers.
func Mul(x, y uint) (lower, upper uint)
func mul(x, y uint) (lower, upper uint) {
	if uintSize == 32 {
		lo, hi := Mul32(uint32(x), uint32(y))
		return uint(lo), uint(hi)
	}
	lo, hi := Mul64(uint64(x), uint64(y))
	return uint(lo), uint(hi)
}

// Mul64 returns the 128-bit unsigned product of a and b as a pair of 64-bit
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

// Mul32 returns the 64-bit unsigned product of a and b as a pair of 32-bit
// unsigned integers.
func Mul32(x, y uint32) (lower, upper uint32)
func mul32(x, y uint32) (lower, upper uint32) {
	const lo = 0xffff

	a, b := x>>16, x&lo
	c, d := y>>16, y&lo

	ac := a * c
	ad := a * d
	bc := b * c
	bd := b * d

	mid := ad&lo + bc&lo + bd>>16 // 18 bits

	lower = (mid&lo)<<16 + bd&lo
	upper = ac + ad>>16 + bc>>16 + mid>>16

	return lower, upper
}
