// Copyright (c) 2015-2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package unsafezero provides utilities for clearing data in private types.
package unsafezero

import (
	"reflect"
	"unsafe"
)

// ClearString zeroes a string's backing array. This is truly s̶t̶u̶p̶i̶d̶ dangerous.
// Here are some considerations:
//	1. The string must be not be in the read-only data segment of the
//	   program (i.e. it must be dynamically allocated).
//	2. No one expects an immutable value to change, so expect subtle bugs
//	   if the string is shared.
func ClearString(s string) {
	hdr := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	for i := 0; i < hdr.Len; i++ {
		*(*byte)(unsafe.Pointer(hdr.Data + uintptr(i))) = 0
	}
}
