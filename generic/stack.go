// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

import "math/bits"

// GenericTypeStack is an optionally auto-growing stack.
type GenericTypeStack struct {
	a        []GenericType
	next     int
	autoGrow bool
}

// NewGenericTypeStack returns a new auto-growing stack that can accommodate
// at least size items.
func NewGenericTypeStack(size int) *GenericTypeStack {
	return NewGenericTypeStackWithBuffer(
		make([]GenericType, 1<<uint(bits.Len(uint(size-1)))),
	)
}

// NewGenericTypeStackWithBuffer returns a new auto-growing stack that wraps
// the provided buffer, which is never resliced beyond its current length.
func NewGenericTypeStackWithBuffer(buf []GenericType) *GenericTypeStack {
	return &GenericTypeStack{
		a:        buf,
		next:     0,
		autoGrow: true,
	}
}

// SetAutoGrow enables or disables auto-growing.
func (s *GenericTypeStack) SetAutoGrow(t bool) {
	s.autoGrow = t
}

// Len returns the current number of elements in the stack.
func (s *GenericTypeStack) Len() int {
	return s.next
}

// Cap returns the logical capacity of the stack. Note that this may be
// smaller than the capacity of the internal slice.
func (s *GenericTypeStack) Cap() int {
	return len(s.a)
}

// Push a new element onto the stack. If adding this element would overflow
// the stack and auto-growing is enabled, the current stack is moved to a
// larger GenericTypeStack before adding the element.
func (s *GenericTypeStack) Push(x GenericType) {
	if s.autoGrow && s.Len() == len(s.a) {
		s.Grow(1)
	}
	s.a[s.next] = x
	s.next++
}

// Pop removes and returns the top element from the stack. Calling Pop on an
// empty stack results in a panic.
func (s *GenericTypeStack) Pop() GenericType {
	s.next--
	return s.a[s.next]
}

// PushSlice adds a slice of GenericType onto the stack. If adding these
// elements would overflow the stack and auto-growing is enabled, the current
// stack is moved to a larger GenericTypeStack before adding the elements.
func (s *GenericTypeStack) PushSlice(src []GenericType) {
	if len(src) == 0 {
		return
	}

	if s.autoGrow {
		newlen := s.Len() + len(src)
		if newlen > len(s.a) {
			s.Grow(newlen - len(s.a))
		}
	}

	s.next += copy(s.a[s.next:], src)
}

// PopSlice removes and writes up to len(dst) elements from the stack into
// dst. The number of popped elements is returned.
func (s *GenericTypeStack) PopSlice(dst []GenericType) (n int) {
	n = len(dst)
	if s.Len() < n {
		n = s.Len()
	}

	for i := 0; i < n; i++ {
		s.next--
		dst[i] = s.a[s.next]
	}

	return n
}

// Peek returns the top element from the stack without removing it. Peeking an
// empty stack results in a panic.
func (s *GenericTypeStack) Peek() GenericType {
	return s.a[s.next-1]
}

// Grow internal slice to accommodate at least n more items.
func (s *GenericTypeStack) Grow(n int) {
	// We do not check to see if n <= cap(q.a) - len(q.a) because we promised
	// never to reslice the current buffer beyond its current length.
	if n <= 0 {
		return
	}

	a := make([]GenericType, 1<<uint(bits.Len(uint(len(s.a)+n-1))))
	copy(a, s.a[:s.next])

	s.a = a
}

// Reset the stack so that its length is zero.
// Note that the internal slice is truncated, NOT cleared.
func (s *GenericTypeStack) Reset() {
	s.next = 0
}
