// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"reflect"
	"unsafe"
)

type listNode struct {
	prev, next int
}

func (node listNode) undefined() bool {
	return node.prev == undefined && node.next == undefined
}

func makeListNodeSlice(buf []int) []listNode {
	s := (*[]listNode)(unsafe.Pointer(&buf))

	// listNode is two ints long
	header := (*reflect.SliceHeader)(unsafe.Pointer(s))
	header.Len = len(buf) / 2
	header.Cap = cap(buf) / 2

	return *s
}

// autoPromotingStack and nonPromotingStack implement the memory-sharing dual
// stack data structure described by Tim Leslie at http://www.timl.id.au/SCC.
//
// autoPromotingStack is a stack with the following interesting property:
// pushing an int that is already on the stack promotes that int to the top.
// This can be done in constant time with a slice-backed doubly linked list
// whose node values are the indices of the slice.
//
// nonPromotingStack is a simple stack that can share memory with an
// autoPromotingStack provided the following constraints hold:
//
//	1) A nonPromotingStack never contains duplicates.
//	2) The sets of elements in an autoPromotingStack and nonPromotingStack
//	   are always disjoint.

type autoPromotingStack struct {
	s   []listNode
	top int
	len int
}

func newAutoPromotingStack(buf []int) *autoPromotingStack {
	return &autoPromotingStack{
		s:   makeListNodeSlice(buf),
		top: undefined,
		len: 0,
	}
}

func (aps *autoPromotingStack) peek() int {
	return aps.top
}

func (aps *autoPromotingStack) pop() int {
	oldtop := aps.top
	aps.top = aps.s[oldtop].prev

	aps.s[oldtop].prev = undefined
	if aps.top != undefined {
		aps.s[aps.top].next = undefined
	}
	aps.len--

	return oldtop
}

// If index n is on the stack, it is promoted to be the top element.
func (aps *autoPromotingStack) pushOrPromote(n int) {
	if aps.s[n].undefined() {
		// Standard doubly linked list append
		if aps.top != undefined {
			aps.s[aps.top].next = n
		}
		aps.s[n].prev = aps.top
		aps.top = n
		aps.len++
		return
	} else if n == aps.top {
		// n is already promoted
		return
	}

	//
	// Promote n
	//

	node := aps.s[n]

	// Unlink node and splice
	if node.prev != undefined {
		aps.s[node.prev].next = node.next
	}
	aps.s[node.next].prev = node.prev

	// Promote node
	if aps.top != undefined {
		aps.s[aps.top].next = n
	}
	aps.s[n].prev = aps.top
	aps.s[n].next = undefined
	aps.top = n
}

type nonPromotingStack struct {
	s   []listNode
	top int
	len int
}

func newNonPromotingStack(buf []int) *nonPromotingStack {
	return &nonPromotingStack{
		s:   makeListNodeSlice(buf),
		top: undefined,
		len: 0,
	}
}

func (nps *nonPromotingStack) peek() int {
	return nps.top
}

func (nps *nonPromotingStack) pop() int {
	top := nps.top
	nps.top = nps.s[top].prev
	nps.len--
	return top
}

func (nps *nonPromotingStack) push(n int) {
	nps.s[n].prev = nps.top
	nps.top = n
	nps.len++
}
