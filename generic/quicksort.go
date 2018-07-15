// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

func QuicksortNumberSlice(v []Number) {
	switch len(v) {
	case 0, 1:
		return
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

// Hoare's partitioning with median of first, middle, and last as pivot
func PartitionNumberSlice(v []Number) int {
	var pivot Number

	if len(v) > 16 {
		pivot = MedianOfThreeNumber(v)
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

func MedianOfThreeNumber(v []Number) Number {
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
