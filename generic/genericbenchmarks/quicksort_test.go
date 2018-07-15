package genericbenchmarks

import (
	"math/rand"
	"sort"
	"testing"
)

func randslice(n int) []int {
	v := make([]int, n)
	for i := range v {
		v[i] = rand.Intn(i + 1)
		// Fisher-Yates shuffle
		for i := len(v) - 1; i >= 0; i-- {
			j := rand.Intn(i + 1)
			v[i], v[j] = v[j], v[i]
		}
	}
	return v
}

// goos: linux
// goarch: amd64
// pkg: github.com/guns/golibs/generic/genericbenchmarks
// BenchmarkSortInts-4              30000     62471 ns/op      32 B/op       1 allocs/op
// BenchmarkQuicksortIntSlice-4     50000     26011 ns/op       0 B/op       0 allocs/op
// BenchmarkQuicksortIntSliceF-4    30000     56733 ns/op       0 B/op       0 allocs/op

const slicelen = 1000

func BenchmarkSortInts(b *testing.B) {
	r := randslice(slicelen)
	s := make([]int, len(r))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(s, r)
		sort.Ints(s)
	}
}
func BenchmarkQuicksortIntSlice(b *testing.B) {
	r := randslice(slicelen)
	s := make([]int, len(r))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(s, r)
		QuicksortIntSlice(s)
	}
}
func BenchmarkQuicksortIntSliceF(b *testing.B) {
	r := randslice(slicelen)
	s := make([]int, len(r))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(s, r)
		QuicksortIntSliceF(s, func(a, b int) bool { return a < b })
	}
}
