// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/guns/golibs/bitslice"
)

func TestWorkspaceBackingSlice(t *testing.T) {
	w := NewWorkspace(4)
	buf := w.a[:cap(w.a)] // Backing slice

	// Fill backing slice with a sentinel value and check fields
	n := int(rand.Int31() + 1)
	for i := range buf {
		buf[i] = n
	}

	a := []int{n, n, n, n}

	if !reflect.DeepEqual(w.a, a) {
		t.Errorf("%v != %v", w.a, a)
	}
	if !reflect.DeepEqual(w.b, a) {
		t.Errorf("%v != %v", w.b, a)
	}
	if !reflect.DeepEqual(w.bitslice, bitslice.T{uint(n)}) {
		t.Errorf("%v != %v", w.bitslice, bitslice.T{uint(n)})
	}

	qp, sp := w.queue.GetSlicePointer(), w.stack.GetSlicePointer()

	if !reflect.DeepEqual(*qp, []int{n, n}) {
		t.Errorf("%v != %v", *qp, []int{n, n})
	}
	if !reflect.DeepEqual(*sp, []int{n, n}) {
		t.Errorf("%v != %v", *sp, []int{n, n})
	}

	// The stack and queue should share the same storage
	for i := range *qp {
		(*qp)[i] = n - 1
	}

	if !reflect.DeepEqual(*qp, *sp) {
		t.Errorf("%v != %v", *qp, *sp)
	}
}
