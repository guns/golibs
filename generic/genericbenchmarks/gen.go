package genericbenchmarks

//go:generate genny -pkg=genericbenchmarks -in=../quicksort.go -out=quicksort.go gen Number=int
//go:generate genny -pkg=genericbenchmarks -in=../quicksortf.go -out=quicksortf.go gen Type=int
