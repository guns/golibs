// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

// MinGenericNumber2 returns the minimum of a and b.
func MinGenericNumber2(a, b GenericNumber) GenericNumber {
	if a < b {
		return a
	}
	return b
}

// MinGenericNumber3 returns the minimum of a, b, and c.
func MinGenericNumber3(a, b, c GenericNumber) GenericNumber {
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

// MinGenericNumberV returns the minimum of all GenericNumber parameters.
func MinGenericNumberV(nums ...GenericNumber) GenericNumber {
	return MinGenericNumberSlice(nums)
}

// MinGenericNumberSlice returns the minimum of a slice of GenericNumber.
func MinGenericNumberSlice(nums []GenericNumber) GenericNumber {
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

// MaxGenericNumber2 returns the maximum of a and b.
func MaxGenericNumber2(a, b GenericNumber) GenericNumber {
	if a > b {
		return a
	}
	return b
}

// MaxGenericNumber3 returns the maximum of a, b, and c.
func MaxGenericNumber3(a, b, c GenericNumber) GenericNumber {
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

// MaxGenericNumberV returns the maximum of all GenericNumber parameters.
func MaxGenericNumberV(nums ...GenericNumber) GenericNumber {
	return MaxGenericNumberSlice(nums)
}

// MaxGenericNumberSlice returns the maximum of a slice of GenericNumber.
func MaxGenericNumberSlice(nums []GenericNumber) GenericNumber {
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
