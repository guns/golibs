// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

import "github.com/guns/golibs/calculate"

// GenericTypeQueue is an optionally auto-growing queue backed by a ring buffer.
type GenericTypeQueue struct {
	a          []GenericType
	head, tail int
	autoGrow   bool
}

// NewGenericTypeQueue returns a new auto-growing queue that can accommodate
// at least size items.
func NewGenericTypeQueue(size int) *GenericTypeQueue {
	return NewGenericTypeQueueWithBuffer(
		make([]GenericType, size),
	)
}

// NewGenericTypeQueueWithBuffer returns a new auto-growing queue that wraps
// the provided buffer, which is never resliced beyond its current length.
func NewGenericTypeQueueWithBuffer(buf []GenericType) *GenericTypeQueue {
	return &GenericTypeQueue{
		a:        buf,
		head:     -1,
		tail:     -1,
		autoGrow: true,
	}
}

// SetAutoGrow enables or disables auto-growing.
func (q *GenericTypeQueue) SetAutoGrow(t bool) {
	q.autoGrow = t
}

// Len returns the current number of elements in the queue.
func (q *GenericTypeQueue) Len() int {
	switch {
	case q.head == -1:
		// Queue is empty:
		//
		//	h
		//	 [_ _ _ _ _ _]
		//	t
		//
		return 0
	case q.head < q.tail:
		// Elements are in order:
		//
		//	    h
		//	 [_ 1 2 3 _ _]
		//	          t
		//
		return q.tail - q.head
	default:
		// Elements begin at rear and continue at front:
		//
		//	        h
		//	 [0 1 _ 3 4 5]
		//	      t
		//
		return len(q.a) - q.head + q.tail
	}
}

// Cap returns the logical capacity of the queue. Note that this may be
// smaller than the capacity of the internal slice.
func (q *GenericTypeQueue) Cap() int {
	return len(q.a)
}

// Enqueue a new element into the queue. If adding this element would overflow
// the queue and auto-growing is enabled, the current queue is moved to a
// larger GenericTypeQueue before adding the element.
func (q *GenericTypeQueue) Enqueue(x GenericType) {
	if q.tail == -1 {
		q.head = 0
		q.tail = 0
		if len(q.a) == 0 {
			q.Grow(1)
		}
	} else if q.autoGrow && q.head == q.tail {
		q.Grow(1)
	}

	q.a[q.tail] = x

	q.tail++
	if q.tail >= len(q.a) {
		q.tail -= len(q.a)
	}
}

// Dequeue removes and returns the next element from the queue. Calling
// Dequeue on an empty queue results in a panic.
func (q *GenericTypeQueue) Dequeue() GenericType {
	x := q.a[q.head]

	q.head++
	if q.head >= len(q.a) {
		q.head -= len(q.a)
	}

	if q.head == q.tail {
		q.Reset()
	}

	return x
}

// EnqueueSlice adds a slice of GenericType into the queue. If adding these
// elements would overflow the queue and auto-growing is enabled, the current
// queue is moved to a larger GenericTypeQueue before adding the elements.
func (q *GenericTypeQueue) EnqueueSlice(src []GenericType) {
	if len(src) == 0 {
		return
	}

	newlen := q.Len() + len(src)
	if newlen > len(q.a) {
		if q.autoGrow {
			q.Grow(newlen - len(q.a))
		} else {
			panic("insufficient space for EnqueueSlice")
		}
	}

	switch {
	case q.head == -1:
		// Queue is empty:
		//
		//	h
		//	 [_ _ _ _ _ _]
		//	t
		//
		q.head = 0
		q.tail = 0

		q.tail += copy(q.a, src)
		if q.tail >= len(q.a) {
			q.tail -= len(q.a)
		}
	case q.tail < q.head:
		// Free segment is contiguous:
		//
		//	          h
		//	 [0 _ _ _ 4 5]
		//	    t
		//
		q.tail += copy(q.a[q.tail:], src)
	default:
		// Free segment begins at rear and continues at front:
		//
		//	      h
		//	 [_ _ 2 3 _ _]
		//	          t
		//
		n := copy(q.a[q.tail:], src)
		if n < len(src) {
			n += copy(q.a, src[n:])
		}

		q.tail += n
		if q.tail >= len(q.a) {
			q.tail -= len(q.a)
		}
	}
}

// DequeueSlice removes and writes up to len(dst) elements from the queue into
// dst. The number of dequeued elements is returned.
func (q *GenericTypeQueue) DequeueSlice(dst []GenericType) (n int) {
	switch {
	case q.head == -1:
		// Queue is empty:
		//
		//	h
		//	 [_ _ _ _ _ _]
		//	t
		return 0
	case q.head < q.tail:
		// Elements are in order:
		//
		//	    h
		//	 [_ 1 2 3 _ _]
		//	          t
		//
		n := copy(dst, q.a[q.head:q.tail])
		q.head += n

		if q.head == q.tail {
			q.Reset()
		}

		return n
	default:
		// Elements begin at rear and continue at front:
		//
		//	        h
		//	 [0 1 _ 3 4 5]
		//	      t
		//
		n := copy(dst, q.a[q.head:])
		if n < len(dst) {
			n += copy(dst[n:], q.a[:q.tail])
		}

		q.head += n
		if q.head >= len(q.a) {
			q.head -= len(q.a)
		}

		if q.head == q.tail {
			q.Reset()
		}

		return n
	}
}

// Peek returns the next element from the queue without removing it. Peeking
// an empty queue results in a panic.
func (q *GenericTypeQueue) Peek() GenericType {
	return q.a[q.head]
}

// Grow internal slice to accommodate at least n more items.
func (q *GenericTypeQueue) Grow(n int) {
	// We do not check to see if n <= cap(q.a) - len(q.a) because we promised
	// never to reslice the current buffer beyond its current length.
	if n <= 0 {
		return
	}

	a := make([]GenericType, calculate.NextCap(len(q.a)+n))

	switch {
	case q.head == -1:
		// Queue is empty:
		//
		//	h
		//	 [_ _ _ _ _ _]
		//	t
		q.a = a
	case q.head < q.tail:
		// Elements are in order:
		//
		//	    h
		//	 [_ 1 2 3 _ _]
		//	          t
		//
		q.tail = copy(a, q.a[q.head:q.tail])
		q.head = 0
		q.a = a
	default:
		// Elements begin at rear and continue at front:
		//
		//	        h
		//	 [0 1 _ 3 4 5]
		//	      t
		//
		n := copy(a, q.a[q.head:])
		n += copy(a[n:], q.a[:q.tail])
		q.tail = n
		q.head = 0
		q.a = a
	}
}

// Reset the queue so that its length is zero.
// Note that the internal slice is truncated, NOT cleared.
func (q *GenericTypeQueue) Reset() {
	q.head = -1
	q.tail = -1
}
