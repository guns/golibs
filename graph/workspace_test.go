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

func TestWorkspace(t *testing.T) {
	w := NewWorkspace(8)
	buf := w.a[:cap(w.a)] // Backing slice
	queue := w.makeQueue(wC)
	stack := w.makeStack(wC)
	qp := queue.GetSlicePointer()
	sp := stack.GetSlicePointer()
	bs := w.makeBitsliceN(2, wC)

	if w.len != 8 {
		t.Errorf("%v != %v", w.len, 8)
	}
	if w.cap != 8 {
		t.Errorf("%v != %v", w.cap, 8)
	}

	// Fill backing slice with a random value and check fields

	n := int(rand.Int31() + 1)
	s := []int{n, n, n, n, n, n, n, n}

	for i := range buf {
		buf[i] = n
	}

	if !reflect.DeepEqual(w.a, s) {
		t.Errorf("%v != %v", w.a, s)
	}
	if !reflect.DeepEqual(w.b, s) {
		t.Errorf("%v != %v", w.b, s)
	}
	if !reflect.DeepEqual(w.c, s) {
		t.Errorf("%v != %v", w.c, s)
	}
	if !reflect.DeepEqual(*qp, s) {
		t.Errorf("%v != %v", *qp, s)
	}
	if !reflect.DeepEqual(*sp, s) {
		t.Errorf("%v != %v", *sp, s)
	}
	if !reflect.DeepEqual(bs, []bitslice.T{{uint(n)}, {uint(n)}}) {
		t.Errorf("%v != %v", bs, []bitslice.T{{uint(n)}, {uint(n)}})
	}

	// Check to see that w.c, bs, queue, and stack refer to the same memory

	for i := range w.c {
		w.c[i] = ^n
	}

	s = []int{^n, ^n, ^n, ^n, ^n, ^n, ^n, ^n}

	if !reflect.DeepEqual(w.c, s) {
		t.Errorf("%v != %v", w.c, s)
	}
	if !reflect.DeepEqual(*qp, s) {
		t.Errorf("%v != %v", *qp, s)
	}
	if !reflect.DeepEqual(*sp, s) {
		t.Errorf("%v != %v", *sp, s)
	}
	if !reflect.DeepEqual(bs, []bitslice.T{{uint(^n)}, {uint(^n)}}) {
		t.Errorf("%v != %v", bs, []bitslice.T{{uint(^n)}, {uint(^n)}})
	}

	// Prepare the workspace for a smaller graph

	z := []int{0, 0, 0, 0}
	zneg := []int{-1, -1, -1, -1}
	s = s[:4]

	w.prepare(4, wA|wBNeg)

	if w.len != 4 {
		t.Errorf("%v != %v", w.len, 4)
	}
	if w.cap != 8 {
		t.Errorf("%v != %v", w.cap, 8)
	}

	if !reflect.DeepEqual(w.a, z) {
		t.Errorf("%v != %v", w.a, z)
	}
	if !reflect.DeepEqual(w.b, zneg) {
		t.Errorf("%v != %v", w.b, zneg)
	}
	if !reflect.DeepEqual(w.c, s) {
		t.Errorf("%v != %v", w.c, s)
	}
	if !reflect.DeepEqual(*qp, s[:8]) {
		t.Errorf("%v != %v", *qp, s)
	}
	if !reflect.DeepEqual(*qp, s[:8]) {
		t.Errorf("%v != %v", *qp, s)
	}
	if !reflect.DeepEqual(bs, []bitslice.T{{uint(^n)}, {uint(^n)}}) {
		t.Errorf("%v != %v", bs, []bitslice.T{{uint(^n)}, {uint(^n)}})
	}

	w.prepare(4, wC)
	s = []int{0, 0, 0, 0, ^n, ^n, ^n, ^n}

	if !reflect.DeepEqual(w.c, z) {
		t.Errorf("%v != %v", w.c, z)
	}
	if !reflect.DeepEqual(*qp, s) {
		t.Errorf("%v != %v", *qp, s)
	}
	if !reflect.DeepEqual(*sp, s) {
		t.Errorf("%v != %v", *sp, s)
	}
	if !reflect.DeepEqual(bs, []bitslice.T{{0}, {0}}) {
		t.Errorf("%v != %v", bs, []bitslice.T{{0}, {0}})
	}

	// Grow the workspace

	w.prepare(16, 0)
	z = make([]int, 16)

	if w.len != 16 {
		t.Errorf("%v != %v", w.len, 16)
	}
	if w.cap != 16 {
		t.Errorf("%v != %v", w.cap, 16)
	}

	if !reflect.DeepEqual(w.a, z) {
		t.Errorf("%v != %v", w.a, z)
	}
	if !reflect.DeepEqual(w.b, z) {
		t.Errorf("%v != %v", w.b, z)
	}
	if !reflect.DeepEqual(w.c, z) {
		t.Errorf("%v != %v", w.c, z)
	}

	if !reflect.DeepEqual(*qp, s) {
		t.Errorf("%v != %v", *qp, s)
	}
	if !reflect.DeepEqual(*sp, s) {
		t.Errorf("%v != %v", *sp, s)
	}
	if !reflect.DeepEqual(bs, []bitslice.T{{0}, {0}}) {
		t.Errorf("%v != %v", bs, []bitslice.T{{0}, {0}})
	}
}
