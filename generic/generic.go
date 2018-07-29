// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package generic provides generic templates for common data structures and
// functions via https://github.com/cheekybits/genny/generic
package generic

import "github.com/cheekybits/genny/generic"

type GenericType generic.Type
type GenericNumber generic.Number
type ComparableType interface{ Less(x *ComparableType) bool } // generic.Type

//go:generate genny -pkg=impl -in=math.go            -out=impl/math.go            gen GenericNumber=int
//go:generate genny -pkg=impl -in=packed2dbuilder.go -out=impl/packed2dbuilder.go gen GenericType=int
//go:generate genny -pkg=impl -in=queue.go           -out=impl/queue.go           gen GenericType=int
//go:generate genny -pkg=impl -in=quicksort.go       -out=impl/quicksort.go       gen GenericNumber=int
//go:generate genny -pkg=impl -in=stack.go           -out=impl/stack.go           gen GenericType=int
