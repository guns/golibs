// Package genericbenchmarks benchmarks the generic implementations against
// the standard library.
package genericbenchmarks

//go:generate genny -pkg=genericbenchmarks -in=../math.go       -out=math.go       gen GenericNumber=int
//go:generate genny -pkg=genericbenchmarks -in=../queue.go      -out=queue.go      gen GenericType=int
//go:generate genny -pkg=genericbenchmarks -in=../quicksort.go  -out=quicksort.go  gen GenericNumber=int
//go:generate genny -pkg=genericbenchmarks -in=../quicksortm.go -out=quicksortm.go gen ComparableType=Person

// Person is a typical struct used for benchmarks
type Person struct {
	name string
	age  int
}

// Less implements the method required by QuicksortPersonSlice
func (p *Person) Less(x *Person) bool {
	if p.name != x.name {
		return p.name < x.name
	}
	return p.age < x.age
}

// PersonSlice is a type that implements sort.Interface
type PersonSlice []Person

func (p PersonSlice) Len() int      { return len(p) }
func (p PersonSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p PersonSlice) Less(i, j int) bool {
	if p[i].name != p[j].name {
		return p[i].name < p[j].name
	}
	return p[i].age < p[j].age
}
