// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"math/bits"

	"github.com/guns/golibs/bitslice"
	"github.com/guns/golibs/generic/impl"
)

// A Workspace is used while traversing Graphs.
type Workspace struct {
	size    int           // Size of this Workspace
	buf     []int         // General scratch buffer
	prev    []int         // Mapping of vertex -> previous vertex
	queue   impl.IntQueue // BFS queue
	stack   impl.IntStack // DFS stack
	visited bitslice.T    // Boolean set of visited vertices
}

// NewWorkspace returns a new Workspace suitable for a Graph of a given size.
func NewWorkspace(size int) *Workspace {
	// Single shared int buffer
	alen := size * 2
	qlen := 1 << uint(bits.Len(uint(size/2-1)))
	buf := make([]int, alen+qlen)

	// The queue and stack share a slice and cannot be used concurrently.
	queue := impl.IntQueue{}
	stack := impl.IntStack{}
	(*queue.GetSlicePointer()) = buf[alen:]
	(*stack.GetSlicePointer()) = buf[alen:]

	return &Workspace{
		size:    size,
		buf:     buf[:size],
		prev:    buf[size:alen],
		queue:   queue,
		stack:   stack,
		visited: bitslice.Make(size),
	}
}

// Resize this workspace.
func (w *Workspace) Resize(size int) {
	if size <= w.size {
		return
	}

	var buf []int
	alen := size * 2
	qlen := 1 << uint(bits.Len(uint(size/2-1)))
	buflen := alen + qlen
	bitslicelen := (size + 63) / 64

	// The queue and stack buffers may have been replaced by larger slices
	// that we might be able to reuse.
	qp, sp := w.queue.GetSlicePointer(), w.stack.GetSlicePointer()
	if len(*qp) >= buflen {
		buf = *qp
	} else if len(*sp) >= buflen {
		buf = *sp
	} else {
		buf = make([]int, buflen)
	}

	w.size = size
	w.buf = buf[:size]
	w.prev = buf[size:alen]
	(*qp) = buf[alen:]
	(*sp) = buf[alen:]

	for i := len(w.visited); i < bitslicelen; i++ {
		w.visited = append(w.visited, 0)
	}
}

// Reset a Workspace.
func (w *Workspace) Reset() {
	for i := range w.prev {
		w.buf[i] = 0
		w.prev[i] = -1
	}
	w.queue.Reset()
	w.stack.Reset()
	w.visited.Reset()
}
