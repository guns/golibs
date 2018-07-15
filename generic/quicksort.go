// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

// QuicksortGenericNumberSlice sorts a slice of GenericNumber in place.
// Elements of type GenericNumber must be comparable by value.
func QuicksortGenericNumberSlice(v []GenericNumber) {
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

	i := PartitionGenericNumberSlice(v)
	QuicksortGenericNumberSlice(v[:i+1])
	QuicksortGenericNumberSlice(v[i+1:])
}

// PartitionGenericNumberSlice partitions a slice of GenericNumber in place
// such that every element 0..index is less than or equal to every element
// index+1..len(v-1). Elements of type GenericNumber must be comparable by value.
func PartitionGenericNumberSlice(v []GenericNumber) (index int) {
	// Hoare's partitioning with median of first, middle, and last as pivot
	var pivot GenericNumber

	if len(v) > 16 {
		pivot = MedianOfThreeGenericNumberSamples(v)
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

// MedianOfThreeGenericNumberSamples returns the median of the first, middle,
// and last element. Elements of type GenericNumber must be comparable by value.
func MedianOfThreeGenericNumberSamples(v []GenericNumber) GenericNumber {
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
