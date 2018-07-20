// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"unsafe"

	"github.com/guns/golibs/bitslice"
	"github.com/guns/golibs/generic/impl"
)

// A Workspace provides general-purpose storage while traversing Graphs.
type Workspace struct {
	len, cap int // Not buffer length/capacity
	a, b, c  []int
	bs       bitslice.T
}

func workspaceInternalOffsets(size int) (a, b, c, buflen int) {
	a = size
	b = a + size
	c = b + size
	buflen = c + bitslice.UintLen(size)
	return
}

// NewWorkspace returns a new Workspace for a Graph of a given size.
func NewWorkspace(size int) *Workspace {
	// Single shared int buffer
	a, b, c, buflen := workspaceInternalOffsets(size)
	buf := make([]int, buflen)
	bs := buf[c:]

	return &Workspace{
		len: size,
		cap: size,
		a:   buf[:a],
		b:   buf[a:b],
		c:   buf[b:c],
		bs:  *(*bitslice.T)(unsafe.Pointer(&bs)),
	}
}

// Resize this workspace. Returns true if a reallocation was necessary, and
// false if not. Note that all buffers are zeroed on reallocation, so a Reset
// may not be necessary after a Resize that triggers a reallocation.
func (w *Workspace) Resize(size int) bool {
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

// WorkspaceField values represent fields of a Workspace.
type WorkspaceField uint

const (
	WA    WorkspaceField = 1 << iota // Select w.a or reset w.a with 0
	WANeg                            // Reset w.a with -1
	WB                               // Select w.b or reset w.b with 0
	WBNeg                            // Reset w.b with -1
	WC                               // Select w.c or reset w.c with 0
	WCNeg                            // Reset w.c with -1
	WBS                              // Select w.bs or reset w.bs
)

func (w *Workspace) selectSlice(field WorkspaceField) []int {
	switch field {
	case WA:
		return w.a
	case WB:
		return w.b
	case WC:
		return w.c
	default:
		return nil // panic() defeats inlining
	}

}

// MakeQueue returns an IntQueue with the specified field as a backing slice.
func (w *Workspace) MakeQueue(field WorkspaceField) impl.IntQueue {
	buf := w.selectSlice(field)[:w.cap]
	q := impl.IntQueue{}

	*q.GetSlicePointer() = buf
	q.Reset()

	return q
}

// MakeStack returns an IntStack with the specified field as a backing slice.
func (w *Workspace) MakeStack(field WorkspaceField) impl.IntStack {
	buf := w.selectSlice(field)[:w.cap]
	s := impl.IntStack{}

	*s.GetSlicePointer() = buf
	s.Reset()

	return s
}

// Reset a Workspace. The fields parameter is a bitfield of WorkspaceField
// values that specify which fields to reset.
func (w *Workspace) Reset(fields WorkspaceField) {
	if fields == 0 {
		return
	}

	if fields&WA > 0 {
		for i := range w.a {
			w.a[i] = 0
		}
	} else if fields&WANeg > 0 {
		for i := range w.a {
			w.a[i] = -1
		}
	}

	if fields&WB > 0 {
		for i := range w.b {
			w.b[i] = 0
		}
	} else if fields&WBNeg > 0 {
		for i := range w.b {
			w.b[i] = -1
		}
	}

	if fields&WC > 0 {
		for i := range w.c {
			w.c[i] = 0
		}
	} else if fields&WCNeg > 0 {
		for i := range w.c {
			w.c[i] = -1
		}
	}

	if fields&WBS > 0 {
		w.bs.Reset()
	}
}

// Prepare a Workspace for a Graph of a given size. The fields parameter is a
// bitfield of WorkspaceField values that specify which fields to reset.
func (w *Workspace) Prepare(size int, fields WorkspaceField) {
	if w.Resize(size) {
		// New workspaces are zero-filled, so avoid unnecessary work.
		fields &= ^(WA | WB | WC | WBS)
	}

	w.Reset(fields)
}
