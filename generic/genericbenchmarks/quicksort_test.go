// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package genericbenchmarks

import (
	"math/rand"
	"sort"
	"testing"
)

//	:: go version go1.10.3 linux/amd64
//	goos: linux
//	goarch: amd64
//	pkg: github.com/guns/golibs/generic/genericbenchmarks
//	BenchmarkSortInts-4                         20000      83534 ns/op       32 B/op        1 allocs/op
//	BenchmarkQuicksortIntSlice-4                50000      34116 ns/op        0 B/op        0 allocs/op
//	BenchmarkSortSortPersonSlice-4              10000     203769 ns/op       32 B/op        1 allocs/op
//	BenchmarkQuicksortPersonSliceMethod-4       10000     169927 ns/op        0 B/op        0 allocs/op
//	PASS
//	ok   github.com/guns/golibs/generic/genericbenchmarks 8.338s

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
func BenchmarkSortSortPersonSlice(b *testing.B) {
	r := randPersonSlice(slicelen)
	s := make([]Person, len(r))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(s, r)
		sort.Sort(PersonSlice(s))
	}
}
func BenchmarkQuicksortPersonSliceMethod(b *testing.B) {
	r := randPersonSlice(slicelen)
	s := make([]Person, len(r))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(s, r)
		QuicksortPersonSlice(s)
	}
}

func randslice(n int) []int {
	v := make([]int, n)
	for i := range v {
		v[i] = rand.Intn(i + 1)
	}
	return v
}

func randPersonSlice(n int) []Person {
	v := make([]Person, n)
	for i := range v {
		v[i] = Person{name: names[rand.Intn(len(names))], age: rand.Intn(100)}
	}
	return v
}

var names = []string{
	"Alberto",
	"Apr",
	"Assyria",
	"Auriga",
	"Bahia",
	"Barron",
	"Baylor",
	"Bert",
	"Bradstreet",
	"Brecht",
	"Bulgari",
	"C",
	"Camelots",
	"Cameroonian",
	"Capitol",
	"Carter",
	"Caterpillar",
	"Chaldea",
	"Chechen",
	"Cheviot",
	"Cheyenne",
	"Christine",
	"Constantinople",
	"Corvallis",
	"Curie",
	"Curtis",
	"Cyclopes",
	"Dale",
	"Danish",
	"Delawares",
	"Dennis",
	"Deuteronomy",
	"Dzerzhinsky",
	"Enterprise",
	"Everest",
	"FDA",
	"Francisca",
	"Frenchwomen",
	"Hank",
	"Harlem",
	"Hogan",
	"ING",
	"Indonesian",
	"Italians",
	"Jennie",
	"Kilauea",
	"Kilroy",
	"Klimt",
	"Krasnoyarsk",
	"Laurie",
	"Lillian",
	"Malaprop",
	"Malaysia",
	"Manichean",
	"McNamara",
	"Melvin",
	"Mendocino",
	"Methuselah",
	"Mich",
	"Milagros",
	"Missouri",
	"Nigerians",
	"Nimitz",
	"Noble",
	"Ojibwa",
	"Omani",
	"Oppenheimer",
	"Oprah",
	"Palembang",
	"Pavlovian",
	"PayPal",
	"Peabody",
	"Petrarch",
	"Piraeus",
	"Pompadour",
	"Ramsey",
	"Reinhold",
	"Rowe",
	"Santayana",
	"Semarang",
	"Serb",
	"Shaka",
	"Shasta",
	"Shepherd",
	"Siamese",
	"Stacy",
	"Stephan",
	"Suzuki",
	"Thompson",
	"Thu",
	"Timex",
	"UAR",
	"Uzi",
	"Vanuatu",
	"Venetian",
	"Verdi",
	"Vijayawada",
	"Vindemiatrix",
	"Watt",
	"Zanuck",
}
