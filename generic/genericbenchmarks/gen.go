// Package genericbenchmarks benchmarks the generic implementations against
// the standard library.
package genericbenchmarks

//go:generate genny -pkg=genericbenchmarks -in=../quicksort.go -out=quicksort.go gen Number=int
//go:generate genny -pkg=genericbenchmarks -in=../quicksortm.go -out=quicksortm.go gen ComparableType=Person

type Person struct {
	name string
	age  int
}

func (p *Person) Less(x *Person) bool {
	if p.name != x.name {
		return p.name < x.name
	}
	return p.age < x.age
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