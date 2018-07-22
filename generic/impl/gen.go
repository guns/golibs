// Package impl provides common concrete implementations of generic templates
// from github.com/guns/golibs/generic
package impl

//go:generate genny -pkg=impl -in=../math.go -out=math.go gen GenericNumber=int,uint
//go:generate genny -pkg=impl -in=../queue.go -out=queue.go gen GenericType=int,uint
//go:generate genny -pkg=impl -in=../quicksort.go -out=quicksort.go gen GenericNumber=int,uint
//go:generate genny -pkg=impl -in=../stack.go -out=stack.go gen GenericType=int,uint
