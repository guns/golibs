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
	a        []int         // Mapping of vertex -> int
	b        []int         // Mapping of vertex -> int
	bitslice bitslice.T    // Mapping of vertex -> bool
	queue    impl.IntQueue // BFS queue
	stack    impl.IntStack // DFS stack
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
		a:        buf[:size],
		b:        buf[size:alen],
		bitslice: bs,
		queue:    queue,
		stack:    stack,
	}
}

// Resize this workspace.
func (w *Workspace) Resize(size int) {
	if size <= len(w.a) {
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

	w.a = buf[:size]
	w.b = buf[size:alen]
	bsbuf := buf[alen : alen+blen]
	w.bitslice = *(*bitslice.T)(unsafe.Pointer(&bsbuf))
	(*qp) = buf[alen+blen:]
	(*sp) = buf[alen+blen:]
}

const (
	WA        = 1 << iota // Reset w.a
	WB                    // Reset w.b
	WBitslice             // Reset w.bitslice
	WAll      = WA | WB | WBitslice
)

// Reset a Workspace. The fields parameter is a bitfield of the Ws* options.
// The aVal and bVal parameters are the values to which w.a and w.b are set if
// the corresponding Ws* options are set.
//
// Note that the internal queue and stack are always reset.
func (w *Workspace) Reset(fields uint, aVal, bVal int) {
	if fields&WA > 0 {
		for i := range w.a {
			w.a[i] = aVal
		}
	}
	if fields&WB > 0 {
		for i := range w.b {
			w.b[i] = bVal
		}
	}
	if fields&WBitslice > 0 {
		w.bitslice.Reset()
	}
	// Resetting a Queue and Stack is very fast, so just do it
	w.queue.Reset()
	w.stack.Reset()
}
