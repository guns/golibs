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
	queue := w.MakeQueue(WC)
	stack := w.MakeStack(WC)
	qp := queue.GetSlicePointer()
	sp := stack.GetSlicePointer()
	bs := w.MakeBitsliceN(2, WC)

	if w.len != 8 {
		t.Errorf("%v != %v", w.len, 8)
	}
	if w.cap != 8 {
		t.Errorf("%v != %v", w.cap, 8)
	}

	// Fill backing slice with a random value and check fields

	n := uint(rand.Int63() + 1)
	s := []uint{n, n, n, n, n, n, n, n}

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
	if !reflect.DeepEqual(bs, []bitslice.T{{n}, {n}}) {
		t.Errorf("%v != %v", bs, []bitslice.T{{n}, {n}})
	}

	// Check to see that w.c, bs, queue, and stack refer to the same memory

	for i := range w.c {
		w.c[i] = ^n
	}

	s = []uint{^n, ^n, ^n, ^n, ^n, ^n, ^n, ^n}

	if !reflect.DeepEqual(w.c, s) {
		t.Errorf("%v != %v", w.c, s)
	}
	if !reflect.DeepEqual(*qp, s) {
		t.Errorf("%v != %v", *qp, s)
	}
	if !reflect.DeepEqual(*sp, s) {
		t.Errorf("%v != %v", *sp, s)
	}
	if !reflect.DeepEqual(bs, []bitslice.T{{^n}, {^n}}) {
		t.Errorf("%v != %v", bs, []bitslice.T{{^n}, {^n}})
	}

	// Prepare the workspace for a smaller graph

	z := []uint{0, 0, 0, 0}
	s = s[:4]

	w.Prepare(4, WA|WB)

	if w.len != 4 {
		t.Errorf("%v != %v", w.len, 4)
	}
	if w.cap != 8 {
		t.Errorf("%v != %v", w.cap, 8)
	}

	if !reflect.DeepEqual(w.a, z) {
		t.Errorf("%v != %v", w.a, z)
	}
	if !reflect.DeepEqual(w.b, z) {
		t.Errorf("%v != %v", w.b, z)
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
	if !reflect.DeepEqual(bs, []bitslice.T{{^n}, {^n}}) {
		t.Errorf("%v != %v", bs, []bitslice.T{{^n}, {^n}})
	}

	w.Prepare(4, WC)
	s = []uint{0, 0, 0, 0, ^n, ^n, ^n, ^n}

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

	w.Prepare(16, 0)
	z = make([]uint, 16)

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
