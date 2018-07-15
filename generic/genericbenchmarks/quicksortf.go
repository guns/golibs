// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package genericbenchmarks

func QuicksortIntSliceF(v []int, less func(a, b int) bool) {
	switch len(v) {
	case 0, 1:
		return
	case 2:
		if less(v[1], v[0]) {
			v[0], v[1] = v[1], v[0]
		}
		return
	case 3:
		if less(v[1], v[0]) {
			v[0], v[1] = v[1], v[0]
		}
		if less(v[2], v[1]) {
			v[1], v[2] = v[2], v[1]
		}
		if less(v[1], v[0]) {
			v[0], v[1] = v[1], v[0]
		}
		return
	}

	i := PartitionIntSliceF(v, less)
	QuicksortIntSliceF(v[:i+1], less)
	QuicksortIntSliceF(v[i+1:], less)
}

// Hoare's partitioning with median of first, middle, and last as pivot
func PartitionIntSliceF(v []int, less func(a, b int) bool) int {
	var pivot int

	if len(v) > 16 {
		pivot = MedianOfThreeIntF(v, less)
	} else {
		pivot = v[(len(v)-1)/2]
	}

	i, j := -1, len(v)

	for {
		for {
			i++
			if !less(v[i], pivot) {
				break
			}
		}

		for {
			j--
			if !less(pivot, v[j]) {
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

func MedianOfThreeIntF(v []int, less func(a, b int) bool) int {
	a := v[0]
	b := v[(len(v)-1)/2]
	c := v[len(v)-1]

	if less(b, a) {
		a, b = b, a
	}
	if less(c, b) {
		b, c = c, b
	}
	if less(b, a) {
		a, b = b, a
	}

	return b
}
