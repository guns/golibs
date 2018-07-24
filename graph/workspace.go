// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"unsafe"

	"github.com/guns/golibs/bitslice"
	"github.com/guns/golibs/generic/impl"
)

// A Workspace provides general-purpose scratch storage for Graph methods.
type Workspace struct {
	len, cap int // Not buffer length/capacity
	a, b, c  []int
}

// NewWorkspace returns a new Workspace for a Graph of a given size.
func NewWorkspace(size int) *Workspace {
	// Single shared int buffer
	buf := make([]int, size*3)

	return &Workspace{
		len: size,
		cap: size,
		a:   buf[:size],
		b:   buf[size : size*2],
		c:   buf[size*2:],
	}
}

// resize this workspace. Returns true if a reallocation was necessary, and
// false if not. Note that all buffers are zeroed on reallocation, so a Reset
// may not be necessary after a resize that triggers a reallocation.
func (w *Workspace) resize(size int) bool {
	if size == w.len {
		return false
	} else if size <= w.cap {
		w.len = size
		w.a = w.a[:size]
		w.b = w.b[:size]
		w.c = w.c[:size]
		return false
	}

	*w = *NewWorkspace(size)
	return true
}

// workspaceField values represent fields of a Workspace.
type workspaceField uint

const (
	wA    workspaceField = 1 << iota // Select or reset (*Workspace).a
	wB                               // Select or reset (*Workspace).b
	wC                               // Select or reset (*Workspace).c
	wANeg                            // Fill (*Workspace).a with -1
	wBNeg                            // Fill (*Workspace).b with -1
	wCNeg                            // Fill (*Workspace).c with -1
)

func (w *Workspace) selectSlice(field workspaceField) []int {
	switch field {
	case wA:
		return w.a
	case wB:
		return w.b
	case wC:
		return w.c
	default:
		return nil // panic() defeats inlining [go1.11]
	}

}

// makeQueue returns an empty IntQueue with the specified field as a backing slice.
func (w *Workspace) makeQueue(field workspaceField) impl.IntQueue {
	buf := w.selectSlice(field)[:w.cap]
	q := impl.IntQueue{}

	*q.GetSlicePointer() = buf
	q.Reset()

	return q
}

// makeStack returns an empty IntStack with the specified field as a backing slice.
func (w *Workspace) makeStack(field workspaceField) impl.IntStack {
	buf := w.selectSlice(field)[:w.cap]
	s := impl.IntStack{}

	*s.GetSlicePointer() = buf
	s.Reset()

	return s
}

// makeBitsliceN returns a slice of n empty bitslice.T with the specified
// field as a backing slice. Each bitslice has a capacity equal to the current
// size of the workspace. The maximum number of bitslices that can be returned
// is equal to:
//
//	currentWorkspaceLen / bitslice.UintLen(currentWorkspaceLen)
//
func (w *Workspace) makeBitsliceN(n int, field workspaceField) []bitslice.T {
	buf := w.selectSlice(field)[:w.cap]
	bs := make([]bitslice.T, n)
	blen := bitslice.UintLen(w.len)
	offset := 0

	for i := 0; i < n; i++ {
		b := buf[offset : offset+blen]
		bs[i] = *(*bitslice.T)(unsafe.Pointer(&b))
		offset += blen
	}

	s := buf[:offset]
	for i := range s {
		s[i] = 0
	}

	return bs
}

// makeSharedStacks returns an autoPromotingStack and nonPromotingStack that
// share memory and are backed by the given fields. The fields parameter must
// specify two contiguous internal fields.
func (w *Workspace) makeSharedStacks(fields workspaceField) (autoPromotingStack, nonPromotingStack) {
	var buf []int

	switch fields {
	case wA | wB:
		buf = w.a[:w.len*2]
	case wB | wC:
		buf = w.b[:w.len*2]
	}

	return *newAutoPromotingStack(buf), *newNonPromotingStack(buf)
}

// reset a Workspace. The fields parameter is a bitfield of workspaceField
// values that specify which fields to reset.
func (w *Workspace) reset(fields workspaceField) {
	if fields == 0 {
		return
	}

	if fields&wA > 0 {
		for i := range w.a {
			w.a[i] = 0
		}
	} else if fields&wANeg > 0 {
		for i := range w.a {
			w.a[i] = -1
		}
	}

	if fields&wB > 0 {
		for i := range w.b {
			w.b[i] = 0
		}
	} else if fields&wBNeg > 0 {
		for i := range w.b {
			w.b[i] = -1
		}
	}

	if fields&wC > 0 {
		for i := range w.c {
			w.c[i] = 0
		}
	} else if fields&wCNeg > 0 {
		for i := range w.c {
			w.c[i] = -1
		}
	}
}

// prepare a Workspace for a Graph of a given size. The fields parameter is a
// bitfield of workspaceField values that specify which fields to reset.
func (w *Workspace) prepare(size int, fields workspaceField) {
	if w.resize(size) {
		// New workspaces are zero-filled, so avoid unnecessary work.
		fields &= ^(wA | wB | wC)
	}

	w.reset(fields)
}
