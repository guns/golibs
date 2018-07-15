// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

import "math/bits"

// A GenericTypeStack is an auto-growing stack.
type GenericTypeStack struct {
	a []GenericType
	i int
}

// DefaultStackLen is the default size of a GenericTypeStack that is created
// with a non-positive size.
const DefaultStackLen = 8

// NewGenericTypeStack returns a new stack that can accommodate at least size items,
// or DefaultStackLen if size <= 0.
func NewGenericTypeStack(size int) *GenericTypeStack {
	if size <= 0 {
		size = DefaultStackLen
	}
	return &GenericTypeStack{
		a: make([]GenericType, 1<<uint(bits.Len(uint(size-1)))),
		i: -1,
	}
}

// Len returns the current number of pushed elements.
func (s *GenericTypeStack) Len() int {
	return s.i + 1
}

// Push a new element onto the stack. If adding this element would overflow
// the stack, the current stack is moved to a new GenericTypeStack twice the
// size of the original before adding the element.
func (s *GenericTypeStack) Push(x GenericType) {
	if s.i == len(s.a)-1 {
		s.Grow(1)
	}
	s.i++
	s.a[s.i] = x
}

// Pop removes and returns the top element from the stack. Calling Pop on an
// empty stack results in a panic.
func (s *GenericTypeStack) Pop() GenericType {
	if s.Len() == 0 {
		panic("stack underflow")
	}
	s.i--
	return s.a[s.i+1]
}

// Peek returns the top element from the stack without removing it. Peeking an
// empty stack results in a panic.
func (s *GenericTypeStack) Peek() GenericType {
	if s.Len() == 0 {
		panic("cannot peek empty GenericTypeStack")
	}

	return s.a[s.i]
}

// Reset the stack so that its length is zero. Note that the internal slice is
// NOT cleared.
func (s *GenericTypeStack) Reset() {
	s.i = -1
}

// Grow internal slice to accommodate at least n more items.
func (s *GenericTypeStack) Grow(n int) {
	n -= cap(s.a) - len(s.a)
	if n <= 0 {
		return
	}

	a := make([]GenericType, 1<<uint(bits.Len(uint(cap(s.a)+n-1))))
	copy(a, s.a)

	s.a = a
}
