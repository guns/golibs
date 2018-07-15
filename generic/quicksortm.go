// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

// QuicksortComparableTypeSlice sorts a slice of ComparableType in place.
// ComparableType must define the following method:
//
//	Less(*ComparableType) bool
//
func QuicksortComparableTypeSlice(v []ComparableType) {
	switch len(v) {
	case 0, 1:
		return
	// Manually sort small slices
	case 2:
		if v[1].Less(&v[0]) {
			v[0], v[1] = v[1], v[0]
		}
		return
	case 3:
		if v[1].Less(&v[0]) {
			v[0], v[1] = v[1], v[0]
		}
		if v[2].Less(&v[1]) {
			v[1], v[2] = v[2], v[1]
		}
		if v[1].Less(&v[0]) {
			v[0], v[1] = v[1], v[0]
		}
		return
	}

	i := PartitionComparableTypeSlice(v)
	QuicksortComparableTypeSlice(v[:i+1])
	QuicksortComparableTypeSlice(v[i+1:])
}

// PartitionComparableTypeSlice partitions a slice of ComparableType in place
// such that every element 0..index is less than or equal to every element
// index+1..len(v-1). ComparableType must define the following method:
//
//	Less(*ComparableType) bool
//
func PartitionComparableTypeSlice(v []ComparableType) (index int) {
	// Hoare's partitioning with median of first, middle, and last as pivot
	var pivot ComparableType

	if len(v) > 16 {
		pivot = MedianOfThreeComparableTypeSamples(v)
	} else {
		pivot = v[(len(v)-1)/2]
	}

	i, j := -1, len(v)

	for {
		for {
			i++
			if !v[i].Less(&pivot) {
				break
			}
		}

		for {
			j--
			if !pivot.Less(&v[j]) {
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

// MedianOfThreeComparableTypeSamples returns the median of the first, middle,
// and last element. ComparableType must define the following method:
//
//	Less(*ComparableType) bool
//
func MedianOfThreeComparableTypeSamples(v []ComparableType) ComparableType {
	a := v[0]
	b := v[(len(v)-1)/2]
	c := v[len(v)-1]

	if b.Less(&a) {
		a, b = b, a
	}
	if c.Less(&b) {
		b = c
	}
	if b.Less(&a) {
		b = a
	}

	return b
}
