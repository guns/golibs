// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

// QuicksortNumberSlice sorts a slice of Number in place. Elements of type
// Number must be comparable by value.
func QuicksortNumberSlice(v []Number) {
	switch len(v) {
	case 0, 1:
		return
	// Manually sort small slices
	case 2:
		if v[1] < v[0] {
			v[0], v[1] = v[1], v[0]
		}
		return
	case 3:
		if v[1] < v[0] {
			v[0], v[1] = v[1], v[0]
		}
		if v[2] < v[1] {
			v[1], v[2] = v[2], v[1]
		}
		if v[1] < v[0] {
			v[0], v[1] = v[1], v[0]
		}
		return
	}

	i := PartitionNumberSlice(v)
	QuicksortNumberSlice(v[:i+1])
	QuicksortNumberSlice(v[i+1:])
}

// PartitionNumberSlice partitions a slice of Number in place such that every
// element 0..index is less than or equal to every element index+1..len(v-1).
// Elements of type Number must be comparable by value.
func PartitionNumberSlice(v []Number) (index int) {
	// Hoare's partitioning with median of first, middle, and last as pivot
	var pivot Number

	if len(v) > 16 {
		pivot = MedianOfThreeNumberSamples(v)
	} else {
		pivot = v[(len(v)-1)/2]
	}

	i, j := -1, len(v)

	for {
		for {
			i++
			if v[i] >= pivot {
				break
			}
		}

		for {
			j--
			if v[j] <= pivot {
				break
			}
		}

		if i < j {
			v[i], v[j] = v[j], v[i]
		} else {
			return j
		}
	}
}

// MedianOfThreeNumberSamples returns the median of the first, middle, and
// last element. Elements of type Number must be comparable by value.
func MedianOfThreeNumberSamples(v []Number) Number {
	a := v[0]
	b := v[(len(v)-1)/2]
	c := v[len(v)-1]

	if b < a {
		a, b = b, a
	}
	if c < b {
		b, c = c, b
	}
	if b < a {
		a, b = b, a
	}

	return b
}
