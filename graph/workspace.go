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
	WA    WorkspaceField = 1 << iota // Select or reset (*Workspace).a
	WB                               // Select or reset (*Workspace).b
	WC                               // Select or reset (*Workspace).c
	WANeg                            // Fill (*Workspace).a with -1
	WBNeg                            // Fill (*Workspace).b with -1
	WCNeg                            // Fill (*Workspace).c with -1
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
		return nil // panic() defeats inlining [go1.11]
	}

}

// MakeQueue returns an empty IntQueue with the specified field as a backing slice.
func (w *Workspace) MakeQueue(field WorkspaceField) impl.IntQueue {
	buf := w.selectSlice(field)[:w.cap]
	q := impl.IntQueue{}

	*q.GetSlicePointer() = buf
	q.Reset()

	return q
}

// MakeStack returns an empty IntStack with the specified field as a backing slice.
func (w *Workspace) MakeStack(field WorkspaceField) impl.IntStack {
	buf := w.selectSlice(field)[:w.cap]
	s := impl.IntStack{}

	*s.GetSlicePointer() = buf
	s.Reset()

	return s
}

// MakeBitsliceN returns a slice of n empty bitslice.T with the specified
// field as a backing slice. Each bitslice has a capacity equal to the current
// size of the workspace. The maximum number of bitslices that can be returned
// is equal to:
//
//	currentWorkspaceLen / bitslice.UintLen(currentWorkspaceLen)
//
func (w *Workspace) MakeBitsliceN(n int, field WorkspaceField) []bitslice.T {
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
}

// Prepare a Workspace for a Graph of a given size. The fields parameter is a
// bitfield of WorkspaceField values that specify which fields to reset.
func (w *Workspace) Prepare(size int, fields WorkspaceField) {
	if w.Resize(size) {
		// New workspaces are zero-filled, so avoid unnecessary work.
		fields &= ^(WA | WB | WC)
	}

	w.Reset(fields)
}
