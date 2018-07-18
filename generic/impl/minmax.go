// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package impl

// MinInt2 returns the minimum of a and b.
func MinInt2(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MinInt3 returns the minimum of a, b, and c.
func MinInt3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// MinIntV returns the minimum of all int parameters.
func MinIntV(nums ...int) int {
	return MinIntSlice(nums)
}

// MinIntSlice returns the minimum of a slice of int.
func MinIntSlice(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	min := nums[0]

	for _, n := range nums[1:] {
		if n < min {
			min = n
		}
	}

	return min
}

// MaxInt2 returns the maximum of a and b.
func MaxInt2(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MaxInt3 returns the maximum of a, b, and c.
func MaxInt3(a, b, c int) int {
	if a > b {
		if a > c {
			return a
		}
		return c
	}
	if b > c {
		return b
	}
	return c
}

// MaxIntV returns the maximum of all int parameters.
func MaxIntV(nums ...int) int {
	return MaxIntSlice(nums)
}

// MaxIntSlice returns the maximum of a slice of int.
func MaxIntSlice(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	max := nums[0]

	for _, n := range nums[1:] {
		if n > max {
			max = n
		}
	}

	return max
}
