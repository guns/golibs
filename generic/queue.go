// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

import "math/bits"

// GenericTypeQueue is an auto-growing queue backed by a ring buffer.
type GenericTypeQueue struct {
	a          []GenericType
	head, tail int
}

// DefaultGenericTypeQueueLen is the default size of a GenericTypeQueue that
// is created with a non-positive size.
const DefaultGenericTypeQueueLen = 8

// NewGenericTypeQueue returns a new queue that can accommodate at least size
// items, or DefaultQueueLen if size <= 0.
func NewGenericTypeQueue(size int) *GenericTypeQueue {
	if size <= 0 {
		size = DefaultGenericTypeQueueLen
	}
	return &GenericTypeQueue{
		a:    make([]GenericType, 1<<uint(bits.Len(uint(size-1)))),
		head: -1,
		tail: -1,
	}
}

// Len returns the current number of queued elements.
func (q *GenericTypeQueue) Len() int {
	switch {
	case q.head == -1:
		return 0
	case q.head < q.tail:
		return q.tail - q.head
	default:
		return len(q.a) - q.head + q.tail
	}
}

// Enqueue a new element into the queue. If adding this element would overflow
// the queue, the current queue is moved to a new GenericTypeQueue twice the
// size of the original before adding the element.
func (q *GenericTypeQueue) Enqueue(x GenericType) {
	if q.tail == -1 {
		q.head = 0
		q.tail = 0
	} else if q.head == q.tail {
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

// Peek returns the next element from the queue without removing it. Peeking
// an empty queue results in a panic.
func (q *GenericTypeQueue) Peek() GenericType {
	return q.a[q.head]
}

// Reset the queue so that its length is zero.
// Note that the internal slice is NOT cleared.
func (q *GenericTypeQueue) Reset() {
	q.head = -1
	q.tail = -1
}

// Grow internal slice to accommodate at least n more items.
func (q *GenericTypeQueue) Grow(n int) {
	// We do not check to see if n <= cap(q.a) - len(q.a) because we'll
	// never have unused capacity.
	if n <= 0 {
		return
	}

	a := make([]GenericType, 1<<uint(bits.Len(uint(len(q.a)+n-1))))

	switch {
	// Queue is empty
	case q.head == -1:
		q.a = a
	// Elements are in order
	case q.head < q.tail:
		copy(a, q.a[q.head:q.tail])
		q.a = a
		q.tail -= q.head
		q.head = 0
	// First segment of elements are at the rear of the array
	default:
		n := copy(a, q.a[q.head:])
		n += copy(a[n:], q.a[:q.tail])
		q.a = a
		q.head = 0
		q.tail = n
	}
}

// GetSlicePointer returns a pointer to the backing slice of this GenericTypeQueue.
// *WARNING* Use at your own risk.
func (q *GenericTypeQueue) GetSlicePointer() *[]GenericType {
	return &q.a
}
