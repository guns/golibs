// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

import (
	"reflect"
	"testing"
)

const pop = -1
const peek = -2
const grow = -3

func TestQueueAndStack(t *testing.T) {
	type T = GenericType
	type Queue = GenericTypeQueue
	type Stack = GenericTypeStack

	data := []struct {
		size     int
		cmds     []interface{}
		queueOut []T
		queue    Queue
		stackOut []T
		stack    Stack
		len      int
	}{
		{
			size:     1,
			cmds:     []interface{}{1, pop, 2, pop},
			queueOut: []T{1, 2},
			queue:    Queue{a: []T{2}, head: -1, tail: -1, autoGrow: true},
			stackOut: []T{1, 2},
			stack:    Stack{a: []T{2}, i: 0, autoGrow: true},
			len:      0,
		},
		{
			size:     1,
			cmds:     []interface{}{1, 2, pop, pop},
			queueOut: []T{1, 2},
			queue:    Queue{a: []T{1, 2}, head: -1, tail: -1, autoGrow: true},
			stackOut: []T{2, 1},
			stack:    Stack{a: []T{1, 2}, i: 0, autoGrow: true},
			len:      0,
		},
		{
			size:     4,
			cmds:     []interface{}{1, 2, 3, 4, pop, pop, 5, 6, 7, 8, pop, pop},
			queueOut: []T{1, 2, 3, 4},
			queue:    Queue{a: []T{3, 4, 5, 6, 7, 8, T(nil), T(nil)}, head: 2, tail: 6, autoGrow: true},
			stackOut: []T{4, 3, 8, 7},
			stack:    Stack{a: []T{1, 2, 5, 6, 7, 8, T(nil), T(nil)}, i: 4, autoGrow: true},
			len:      4,
		},
		{
			size:     5,
			cmds:     []interface{}{1, 2, 3, 4, pop, pop, 5, 6, 7, 8, pop, pop},
			queueOut: []T{1, 2, 3, 4},
			queue:    Queue{a: []T{1, 2, 3, 4, 5, 6, 7, 8}, head: 4, tail: 0, autoGrow: true},
			stackOut: []T{4, 3, 8, 7},
			stack:    Stack{a: []T{1, 2, 5, 6, 7, 8, T(nil), T(nil)}, i: 4, autoGrow: true},
			len:      4,
		},
		{
			size:     4,
			cmds:     []interface{}{1, 2, peek, peek, 3, 4, pop, pop},
			queueOut: []T{1, 1, 1, 2},
			queue:    Queue{a: []T{1, 2, 3, 4}, head: 2, tail: 0, autoGrow: true},
			stackOut: []T{2, 2, 4, 3},
			stack:    Stack{a: []T{1, 2, 3, 4}, i: 2, autoGrow: true},
			len:      2,
		},
		// Default size
		{
			size:     0,
			cmds:     []interface{}{1, 2, 3, 4, pop, pop, pop, pop},
			queueOut: []T{1, 2, 3, 4},
			queue:    Queue{a: []T{1, 2, 3, 4, T(nil), T(nil), T(nil), T(nil)}, head: -1, tail: -1, autoGrow: true},
			stackOut: []T{4, 3, 2, 1},
			stack:    Stack{a: []T{1, 2, 3, 4, T(nil), T(nil), T(nil), T(nil)}, i: 0, autoGrow: true},
			len:      0,
		},
		// Grow
		{
			size:     1,
			cmds:     []interface{}{grow - 0, grow - 0},
			queueOut: []T{},
			queue:    *NewGenericTypeQueue(1, true),
			stackOut: []T{},
			stack:    *NewGenericTypeStack(1, true),
			len:      0,
		},
		{
			size:     1,
			cmds:     []interface{}{grow - 1, 1, grow - 2, 2, 3, 4, pop, pop, 5, 6},
			queueOut: []T{1, 2},
			queue:    Queue{a: []T{5, 6, 3, 4}, head: 2, tail: 2, autoGrow: true},
			stackOut: []T{4, 3},
			stack:    Stack{a: []T{1, 2, 5, 6}, i: 4, autoGrow: true},
			len:      4,
		},
		// Add slices
		{
			size:     4,
			cmds:     []interface{}{[]T{1, 2, 3}, pop, []T{4, 5}},
			queueOut: []T{1},
			queue:    Queue{a: []T{5, 2, 3, 4}, head: 1, tail: 1, autoGrow: true},
			stackOut: []T{3},
			stack:    Stack{a: []T{1, 2, 4, 5}, i: 4, autoGrow: true},
			len:      4,
		},
		{
			size:     4,
			cmds:     []interface{}{[]T{1, 2, 3, 4}, pop, pop, []T{}, pop, 5, []T{6, 7}},
			queueOut: []T{1, 2, 3},
			queue:    Queue{a: []T{5, 6, 7, 4}, head: 3, tail: 3, autoGrow: true},
			stackOut: []T{4, 3, 2},
			stack:    Stack{a: []T{1, 5, 6, 7}, i: 4, autoGrow: true},
			len:      4,
		},
	}

	for i, row := range data {
		q := NewGenericTypeQueue(row.size, true)
		s := NewGenericTypeStack(row.size, true)
		qout := make([]T, 0, len(row.queueOut))
		sout := make([]T, 0, len(row.stackOut))

		for _, x := range row.cmds {
			n, isNum := x.(int)
			xs, isSlice := x.([]T)

			switch {
			case isNum && n == pop:
				qout = append(qout, q.Dequeue())
				sout = append(sout, s.Pop())
			case isNum && n == peek:
				qout = append(qout, q.Peek())
				sout = append(sout, s.Peek())
			case isNum && n <= grow:
				q.Grow(-x.(int) - -grow)
				s.Grow(-x.(int) - -grow)
			default:
				if isSlice {
					q.EnqueueSlice(xs)
					s.PushSlice(xs)
				} else {
					q.Enqueue(x)
					s.Push(x)
				}
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

		if q.Len() != row.len {
			t.Errorf("%v != %v", q.Len(), row.len)
		}
		if s.Len() != row.len {
			t.Errorf("%v != %v", s.Len(), row.len)
		}
		q.Reset()
		s.Reset()
		if q.Len() != 0 {
			t.Errorf("%v != %v", q.Len(), 0)
		}
		if s.Len() != 0 {
			t.Errorf("%v != %v", s.Len(), 0)
		}
	}
}
