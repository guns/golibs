// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"math/rand"
	"reflect"
	"testing"
	"unsafe"

	"github.com/guns/golibs/bitslice"
)

func TestWorkspace(t *testing.T) {
	w := NewWorkspace(8)

	buf := w.a[:cap(w.a)] // Backing slice
	queue := w.queue(wC)
	stack := w.stack(wC)
	bs := w.bitslices(2, wC)

	type IntQueue struct {
		a          []int
		head, tail int
		autoGrow   bool
	}

	type IntStack struct {
		a        []int
		next     int
		autoGrow bool
	}

	qbuf := (*(*IntQueue)(unsafe.Pointer(&queue))).a
	sbuf := (*(*IntStack)(unsafe.Pointer(&stack))).a

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
	if !reflect.DeepEqual(qbuf, s) {
		t.Errorf("%v != %v", qbuf, s)
	}
	if !reflect.DeepEqual(sbuf, s) {
		t.Errorf("%v != %v", sbuf, s)
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
	if !reflect.DeepEqual(qbuf, s) {
		t.Errorf("%v != %v", qbuf, s)
	}
	if !reflect.DeepEqual(sbuf, s) {
		t.Errorf("%v != %v", sbuf, s)
	}
	if !reflect.DeepEqual(bs, []bitslice.T{{uint(^n)}, {uint(^n)}}) {
		t.Errorf("%v != %v", bs, []bitslice.T{{uint(^n)}, {uint(^n)}})
	}

	// Prepare the workspace for a smaller graph

	z := []int{0, 0, 0, 0}
	zneg := []int{-1, -1, -1, -1}
	s = s[:4]

	w.reset(4, wA)
	fillUndefined(w.b)

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
	if !reflect.DeepEqual(qbuf, s[:8]) {
		t.Errorf("%v != %v", qbuf, s)
	}
	if !reflect.DeepEqual(qbuf, s[:8]) {
		t.Errorf("%v != %v", qbuf, s)
	}
	if !reflect.DeepEqual(bs, []bitslice.T{{uint(^n)}, {uint(^n)}}) {
		t.Errorf("%v != %v", bs, []bitslice.T{{uint(^n)}, {uint(^n)}})
	}

	w.reset(4, wC)
	s = []int{0, 0, 0, 0, ^n, ^n, ^n, ^n}

	if !reflect.DeepEqual(w.c, z) {
		t.Errorf("%v != %v", w.c, z)
	}
	if !reflect.DeepEqual(qbuf, s) {
		t.Errorf("%v != %v", qbuf, s)
	}
	if !reflect.DeepEqual(sbuf, s) {
		t.Errorf("%v != %v", sbuf, s)
	}
	if !reflect.DeepEqual(bs, []bitslice.T{{0}, {0}}) {
		t.Errorf("%v != %v", bs, []bitslice.T{{0}, {0}})
	}

	// Grow the workspace

	w.reset(16, 0)
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

	if !reflect.DeepEqual(qbuf, s) {
		t.Errorf("%v != %v", qbuf, s)
	}
	if !reflect.DeepEqual(sbuf, s) {
		t.Errorf("%v != %v", sbuf, s)
	}
	if !reflect.DeepEqual(bs, []bitslice.T{{0}, {0}}) {
		t.Errorf("%v != %v", bs, []bitslice.T{{0}, {0}})
	}

	// Test shared stacks

	fillUndefined(w.a)
	fillUndefined(w.b)
	aps := w.autoPromotingStack(wA | wB)
	nps := newNonPromotingStack(aps.s)

	for i := 0; i < w.len; i++ {
		aps.PushOrPromote(i)
	}

	if !reflect.DeepEqual(aps.s, nps.s) {
		t.Errorf("%v != %v", aps.s, nps.s)
	}

	apsSlice := *(*[]int)(unsafe.Pointer(&aps.s))
	header := (*reflect.SliceHeader)(unsafe.Pointer(&apsSlice))
	header.Len *= 2
	header.Cap *= 2

	if !reflect.DeepEqual(apsSlice, w.a[:w.len*2]) {
		t.Logf("%v !=", apsSlice)
		t.Logf("%v", w.a[:w.len*2])
		t.Fail()
	}
}
