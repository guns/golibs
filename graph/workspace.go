// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"math/bits"
	"unsafe"

	"github.com/guns/golibs/bitslice"
	"github.com/guns/golibs/generic/impl"
)

// A Workspace is used while traversing Graphs.
type Workspace struct {
	buf     []int         // General scratch buffer
	prev    []int         // Mapping of vertex -> previous vertex
	visited bitslice.T    // Boolean set of visited vertices
	queue   impl.IntQueue // BFS queue
	stack   impl.IntStack // DFS stack
}

// NewWorkspace returns a new Workspace suitable for a Graph of a given size.
func NewWorkspace(size int) *Workspace {
	// Single shared int buffer
	alen := size * 2
	blen := bitslice.UintLen(size)
	qlen := 1 << uint(bits.Len(uint(size/2-1)))
	buf := make([]int, alen+blen+qlen)

	// The bitslice is between the int buffers and the queue buffers
	bsbuf := buf[alen : alen+blen]
	bs := *(*bitslice.T)(unsafe.Pointer(&bsbuf))

	// The queue and stack share a slice and cannot be used concurrently.
	queue := impl.IntQueue{}
	stack := impl.IntStack{}
	(*queue.GetSlicePointer()) = buf[alen+blen:]
	(*stack.GetSlicePointer()) = buf[alen+blen:]

	return &Workspace{
		buf:     buf[:size],
		prev:    buf[size:alen],
		visited: bs,
		queue:   queue,
		stack:   stack,
	}
}

// Resize this workspace.
func (w *Workspace) Resize(size int) {
	if size <= len(w.buf) {
		return
	}

	var buf []int
	alen := size * 2
	blen := bitslice.UintLen(size)
	qlen := 1 << uint(bits.Len(uint(size/2-1)))
	buflen := alen + blen + qlen

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

	w.buf = buf[:size]
	w.prev = buf[size:alen]
	bsbuf := buf[alen : alen+blen]
	w.visited = *(*bitslice.T)(unsafe.Pointer(&bsbuf))
	(*qp) = buf[alen+blen:]
	(*sp) = buf[alen+blen:]
}

// Reset a Workspace.
func (w *Workspace) Reset() {
	for i := range w.buf {
		w.buf[i] = 0
		w.prev[i] = -1
	}
	w.visited.Reset()
	w.queue.Reset()
	w.stack.Reset()
}
