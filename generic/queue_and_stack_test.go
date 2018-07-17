// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

import (
	"reflect"
	"testing"
)

func TestQueueAndStack(t *testing.T) {
	data := []struct {
		size     int
		cmds     []int
		queueOut []GenericType
		queue    GenericTypeQueue
		stackOut []GenericType
		stack    GenericTypeStack
	}{
		{
			size:     1,
			cmds:     []int{1, 0, 2, 0},
			queueOut: []GenericType{1, 2},
			queue:    GenericTypeQueue{a: []GenericType{2}, head: -1, tail: -1},
			stackOut: []GenericType{1, 2},
			stack:    GenericTypeStack{a: []GenericType{2}, i: -1},
		},
		{
			size:     1,
			cmds:     []int{1, 2, 0, 0},
			queueOut: []GenericType{1, 2},
			queue:    GenericTypeQueue{a: []GenericType{1, 2}, head: -1, tail: -1},
			stackOut: []GenericType{2, 1},
			stack:    GenericTypeStack{a: []GenericType{1, 2}, i: -1},
		},
		{
			size:     4,
			cmds:     []int{1, 2, 3, 4, 0, 0, 5, 6, 7, 8, 0, 0},
			queueOut: []GenericType{1, 2, 3, 4},
			queue:    GenericTypeQueue{a: []GenericType{3, 4, 5, 6, 7, 8, GenericType(nil), GenericType(nil)}, head: 2, tail: 6},
			stackOut: []GenericType{4, 3, 8, 7},
			stack:    GenericTypeStack{a: []GenericType{1, 2, 5, 6, 7, 8, GenericType(nil), GenericType(nil)}, i: 3},
		},
		{
			size:     5,
			cmds:     []int{1, 2, 3, 4, 0, 0, 5, 6, 7, 8, 0, 0},
			queueOut: []GenericType{1, 2, 3, 4},
			queue:    GenericTypeQueue{a: []GenericType{1, 2, 3, 4, 5, 6, 7, 8}, head: 4, tail: 0},
			stackOut: []GenericType{4, 3, 8, 7},
			stack:    GenericTypeStack{a: []GenericType{1, 2, 5, 6, 7, 8, GenericType(nil), GenericType(nil)}, i: 3},
		},
		{
			size:     4,
			cmds:     []int{1, 2, -1, -1, 3, 4, 0, 0},
			queueOut: []GenericType{1, 1, 1, 2},
			queue:    GenericTypeQueue{a: []GenericType{1, 2, 3, 4}, head: 2, tail: 0},
			stackOut: []GenericType{2, 2, 4, 3},
			stack:    GenericTypeStack{a: []GenericType{1, 2, 3, 4}, i: 1},
		},
		{
			size:     0,
			cmds:     []int{1, 2, 3, 4, 0, 0, 0, 0},
			queueOut: []GenericType{1, 2, 3, 4},
			queue:    GenericTypeQueue{a: []GenericType{1, 2, 3, 4, GenericType(nil), GenericType(nil), GenericType(nil), GenericType(nil)}, head: -1, tail: -1},
			stackOut: []GenericType{4, 3, 2, 1},
			stack:    GenericTypeStack{a: []GenericType{1, 2, 3, 4, GenericType(nil), GenericType(nil), GenericType(nil), GenericType(nil)}, i: -1},
		},
	}

	for i, row := range data {
		q := NewGenericTypeQueue(row.size)
		s := NewGenericTypeStack(row.size)
		qout := make([]GenericType, 0, len(row.queueOut))
		sout := make([]GenericType, 0, len(row.stackOut))

		for _, n := range row.cmds {
			switch n {
			case -1:
				qout = append(qout, q.Peek())
				sout = append(sout, s.Peek())
			case 0:
				qout = append(qout, q.Dequeue())
				sout = append(sout, s.Pop())
			default:
				q.Enqueue(n)
				s.Push(n)
			}
		}

		if !reflect.DeepEqual(qout, row.queueOut) {
			t.Errorf("[%d] %v != %v", i, qout, row.queueOut)
		}
		if !reflect.DeepEqual(sout, row.stackOut) {
			t.Errorf("[%d] %v != %v", i, sout, row.stackOut)
		}

		if !reflect.DeepEqual(*q, row.queue) {
			t.Errorf("[%d] %v != %v", i, *q, row.queue)
		}
		if !reflect.DeepEqual(*s, row.stack) {
			t.Errorf("[%d] %v != %v", i, *s, row.stack)
		}

		qp := q.GetSlicePointer()
		if qp != &q.a {
			t.Errorf("%p != %p", qp, &q.a)
		}
		sp := s.GetSlicePointer()
		if sp != &s.a {
			t.Errorf("%p != %p", sp, &s.a)
		}

		q.Reset()
		if q.Len() != 0 {
			t.Errorf("%v != %v", q.Len(), 0)
		}
		s.Reset()
		if s.Len() != 0 {
			t.Errorf("%v != %v", s.Len(), 0)
		}
	}
}
