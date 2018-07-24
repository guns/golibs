// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

import "math/bits"

// GenericTypeStack is an optionally auto-growing stack.
type GenericTypeStack struct {
	a        []GenericType
	i        int
	autoGrow bool
}

// NewGenericTypeStack returns a new auto-growing stack that can accommodate
// at least size items.
func NewGenericTypeStack(size int) *GenericTypeStack {
	if size <= 0 {
		size = 8 // Sane minimum length
	}
	return NewGenericTypeStackWithBuffer(
		make([]GenericType, 1<<uint(bits.Len(uint(size-1)))),
	)
}

// NewGenericTypeStackWithBuffer returns a new auto-growing stack that wraps
// the provided buffer, which is never resliced beyond its current length.
func NewGenericTypeStackWithBuffer(buf []GenericType) *GenericTypeStack {
	return &GenericTypeStack{
		a:        buf,
		i:        0,
		autoGrow: true,
	}
}

// SetAutoGrow enables or disables auto-growing.
func (s *GenericTypeStack) SetAutoGrow(t bool) {
	s.autoGrow = t
}

// Len returns the current number of pushed elements.
func (s *GenericTypeStack) Len() int {
	return s.i
}

// Push a new element onto the stack. If adding this element would overflow
// the stack and auto-growing is enabled, the current stack is moved to a
// larger GenericTypeStack before adding the element.
func (s *GenericTypeStack) Push(x GenericType) {
	if s.autoGrow && s.Len() == len(s.a) {
		s.Grow(1)
	}
	s.a[s.i] = x
	s.i++
}

// PushSlice adds a slice of GenericType onto the stack. If adding these
// elements would overflow the stack and auto-growing is enabled, the current
// stack is moved to a larger GenericTypeStack before adding the elements.
// Note that the slice is copied into the stack in-order instead of being
// pushed onto the stack one by one.
func (s *GenericTypeStack) PushSlice(xs []GenericType) {
	if len(xs) == 0 {
		return
	}

	newlen := s.Len() + len(xs)
	if s.autoGrow && newlen > len(s.a) {
		s.Grow(newlen - len(s.a))
	}

	copy(s.a[s.i:], xs)
	s.i += len(xs)
}

// Pop removes and returns the top element from the stack. Calling Pop on an
// empty stack results in a panic.
func (s *GenericTypeStack) Pop() GenericType {
	s.i--
	return s.a[s.i]
}

// Peek returns the top element from the stack without removing it. Peeking an
// empty stack results in a panic.
func (s *GenericTypeStack) Peek() GenericType {
	return s.a[s.i-1]
}

// Grow internal slice to accommodate at least n more items.
func (s *GenericTypeStack) Grow(n int) {
	// We do not check to see if n <= cap(q.a) - len(q.a) because we promised
	// never to reslice the current buffer beyond its current length.
	if n <= 0 {
		return
	}

	a := make([]GenericType, 1<<uint(bits.Len(uint(len(s.a)+n-1))))
	copy(a, s.a)

	s.a = a
}

// Reset the stack so that its length is zero.
// Note that the internal slice is NOT cleared.
func (s *GenericTypeStack) Reset() {
	s.i = 0
}
