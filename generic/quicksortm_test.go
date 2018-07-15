// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

var names = []string{
	"Adam",
	"Adan",
	"Althea",
	"Altman",
	"Appomattox",
	"Aquinas",
	"BP",
	"Bacchanalia",
	"Belushi",
	"Berenice",
	"Bernadine",
	"Bernoulli",
	"Betty",
	"Bootes",
	"Bridget",
	"CGI",
	"CZ",
	"Chicagoan",
	"Christmastide",
	"Corinth",
	"Cyclades",
	"Czechs",
	"Daedalus",
	"Datamation",
	"Delawares",
	"Doolittle",
	"Dyson",
	"Elroy",
	"Ernest",
	"Fibonacci",
	"Fitzroy",
	"Foreman",
	"GOP",
	"Gish",
	"Golding",
	"Gonzalo",
	"Guamanian",
	"Haskell",
	"Hooke",
	"Hooters",
	"Iowans",
	"Israel",
	"Josiah",
	"Kafka",
	"Kari",
	"Kip",
	"Korans",
	"Lethe",
	"Lew",
	"Lilian",
	"Liliana",
	"Lome",
	"Londoners",
	"Luisa",
	"Lynette",
	"Maccabees",
	"Macumba",
	"Maiman",
	"Malagasy",
	"Mazarin",
	"Melpomene",
	"Millie",
	"Mirabeau",
	"Monera",
	"Moody",
	"Navajos",
	"Noyes",
	"Occident",
	"Oct",
	"Odis",
	"Olsen",
	"Percheron",
	"Perth",
	"Powhatan",
	"Randal",
	"Rene",
	"Richthofen",
	"Rorschach",
	"Russian",
	"Seleucus",
	"Seton",
	"Sinclair",
	"Slav",
	"Southwests",
	"Sumatrans",
	"Sundanese",
	"Swiss",
	"Tamika",
	"Tamra",
	"Togo",
	"Ur",
	"Ut",
	"Valarie",
	"Velveeta",
	"Walker",
	"Wilberforce",
	"Winnie",
	"Wuhan",
	"Xmases",
	"Yvonne",
}

type Person struct {
	name string
	age  int
}

func (p Person) Less(x *ComparableType) bool {
	if p.name != (*x).(Person).name {
		return p.name < (*x).(Person).name
	}
	return p.age < (*x).(Person).age
}

type PersonSlice []Person

func (p PersonSlice) Len() int      { return len(p) }
func (p PersonSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p PersonSlice) Less(i, j int) bool {
	if p[i].name != p[j].name {
		return p[i].name < p[j].name
	}
	return p[i].age < p[j].age
}

func RandPersonSlice(n int) []Person {
	v := make([]Person, n)
	for i := range v {
		v[i] = Person{name: names[rand.Intn(len(names))], age: rand.Intn(100)}
	}
	return v
}

func TestQuicksortComparableTypeSlice(t *testing.T) {
	for i := 0; i < 100; i++ {
		r := RandPersonSlice(10 * i)

		r2 := make([]ComparableType, len(r))
		s1 := make(PersonSlice, len(r))
		s2 := make(PersonSlice, len(r))

		for i := range r {
			s1[i] = r[i]
			r2[i] = r[i]
		}

		sort.Sort(s1)
		QuicksortComparableTypeSlice(r2)

		for i := range r2 {
			s2[i] = r2[i].(Person)
		}

		if !reflect.DeepEqual(s2, s1) {
			t.Logf("QuicksortNumberSlice:")
			t.Logf("%v !=", s2)
			t.Logf("%v", s1)
			t.Fail()
		}
	}
}

func TestMedianOfThreeComparableTypeSamples(t *testing.T) {
	data := [][]Person{
		{{age: 0}, {age: 1}, {age: 2}},
		{{age: 0}, {age: 2}, {age: 1}},
		{{age: 1}, {age: 0}, {age: 2}},
		{{age: 1}, {age: 2}, {age: 0}},
		{{age: 2}, {age: 0}, {age: 1}},
		{{age: 2}, {age: 1}, {age: 0}},
	}

	for _, v := range data {
		s := make([]ComparableType, len(v))
		for i := range v {
			s[i] = v[i]
		}
		if MedianOfThreeComparableTypeSamples(s).(Person).age != 1 {
			t.Errorf("MedianOfThreeNumberSamples(%v): %v != %v", v, MedianOfThreeComparableTypeSamples(s).(Person).age, 1)
		}
	}
}
