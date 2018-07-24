// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

import "math/bits"

// GenericTypeQueue is an optionally auto-growing queue backed by a ring buffer.
type GenericTypeQueue struct {
	a          []GenericType
	head, tail int
	autoGrow   bool
}

// NewGenericTypeQueue returns a new auto-growing queue that can accommodate
// at least size items.
func NewGenericTypeQueue(size int) *GenericTypeQueue {
	if size <= 0 {
		size = 8 // Sane minimum length
	}
	return NewGenericTypeQueueWithBuffer(
		make([]GenericType, 1<<uint(bits.Len(uint(size-1)))),
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

// Len returns the current number of queued elements.
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

// Enqueue a new element into the queue. If adding this element would overflow
// the queue and auto-growing is enabled, the current queue is moved to a
// larger GenericTypeQueue before adding the element.
func (q *GenericTypeQueue) Enqueue(x GenericType) {
	if q.tail == -1 {
		q.head = 0
		q.tail = 0
	} else if q.autoGrow && q.head == q.tail {
		q.Grow(1)
	}

	q.a[q.tail] = x

	q.tail++
	if q.tail >= len(q.a) {
		q.tail -= len(q.a)
	}
}

// EnqueueSlice adds a slice of GenericType into the queue. If adding these
// elements would overflow the queue and auto-growing is enabled, the current
// queue is moved to a larger GenericTypeQueue before adding the elements.
func (q *GenericTypeQueue) EnqueueSlice(xs []GenericType) {
	if len(xs) == 0 {
		return
	}

	newlen := q.Len() + len(xs)
	if q.autoGrow && newlen > len(q.a) {
		q.Grow(newlen - len(q.a))
	}

	switch {
	case q.head == -1:
		// Queue is empty:
		//
		//	h
		//	 [_ _ _ _ _ _]
		//	t
		//
		copy(q.a, xs)
		q.head = 0
		q.tail = len(xs)

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
		copy(q.a[q.tail:], xs)
		q.tail += len(xs)
	default:
		// Free segment begins at rear and continues at front:
		//
		//	      h
		//	 [_ _ 2 3 _ _]
		//	          t
		//
		n := copy(q.a[q.tail:], xs)

		if n < len(xs) {
			n += copy(q.a, xs[n:])
		}

		q.tail += n
		if q.tail >= len(q.a) {
			q.tail -= len(q.a)
		}
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

	a := make([]GenericType, 1<<uint(bits.Len(uint(len(q.a)+n-1))))

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
		copy(a, q.a[q.head:q.tail])
		q.a = a
		q.tail -= q.head
		q.head = 0
	default:
		// Elements begin at rear and continue at front:
		//
		//	        h
		//	 [0 1 _ 3 4 5]
		//	      t
		//
		n := copy(a, q.a[q.head:])
		n += copy(a[n:], q.a[:q.tail])
		q.a = a
		q.head = 0
		q.tail = n
	}
}

// Reset the queue so that its length is zero.
// Note that the internal slice is NOT cleared.
func (q *GenericTypeQueue) Reset() {
	q.head = -1
	q.tail = -1
}
