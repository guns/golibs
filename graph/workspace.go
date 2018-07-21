// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"github.com/guns/golibs/bitslice"
	"github.com/guns/golibs/generic/impl"
)

// A Workspace provides general-purpose storage while traversing Graphs.
type Workspace struct {
	len, cap int // Not buffer length/capacity
	a, b, c  []uint
}

// NewWorkspace returns a new Workspace for a Graph of a given size.
func NewWorkspace(size int) *Workspace {
	// Single shared uint buffer
	buf := make([]uint, size*3)

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
	WA WorkspaceField = 1 << iota // Select or reset (*Workspace).a
	WB                            // Select or reset (*Workspace).b
	WC                            // Select or reset (*Workspace).c
)

func (w *Workspace) selectSlice(field WorkspaceField) []uint {
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

// MakeQueue returns an empty UintQueue with the specified field as a backing slice.
func (w *Workspace) MakeQueue(field WorkspaceField) impl.UintQueue {
	buf := w.selectSlice(field)[:w.cap]
	q := impl.UintQueue{}

	*q.GetSlicePointer() = buf
	q.Reset()

	return q
}

// MakeStack returns an empty UintStack with the specified field as a backing slice.
func (w *Workspace) MakeStack(field WorkspaceField) impl.UintStack {
	buf := w.selectSlice(field)[:w.cap]
	s := impl.UintStack{}

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
		bs[i] = buf[offset : offset+blen]
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
	if fields&WA > 0 {
		for i := range w.a {
			w.a[i] = 0
		}
	}

	if fields&WB > 0 {
		for i := range w.b {
			w.b[i] = 0
		}
	}

	if fields&WC > 0 {
		for i := range w.c {
			w.c[i] = 0
		}
	}
}

// Prepare a Workspace for a Graph of a given size. The fields parameter is a
// bitfield of WorkspaceField values that specify which fields to reset.
func (w *Workspace) Prepare(size int, fields WorkspaceField) {
	// Reallocated workspaces are zero-filled, so avoid unnecessary work.
	if !w.Resize(size) {
		w.Reset(fields)
	}
}
