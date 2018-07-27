// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

import (
	"reflect"
	"testing"
)

func TestQueueAndStack(t *testing.T) {
	type T = GenericType
	type Queue = GenericTypeQueue
	type Stack = GenericTypeStack

	const POP = -1
	const POPS = -2
	const PEEK = -3
	const LEN = -4
	const CAP = -5
	const RESET = -6
	const GROW = -7

	data := []struct {
		size int
		cmds []T
		qout []T
		sout []T
	}{
		// Buffer size 0
		{
			size: 0,
			cmds: []T{LEN, 2, LEN, PEEK, POP, LEN},
			qout: []T{0, 1, 2, 2, 0},
			sout: []T{0, 1, 2, 2, 0},
		},
		// FIFO vs LIFO
		{
			size: 4,
			cmds: []T{1, 2, 3, 4},
			qout: []T{1, 2, 3, 4},
			sout: []T{4, 3, 2, 1},
		},
		// Fill, remove some, fill
		{
			size: 4,
			cmds: []T{1, 2, 3, LEN, 4, POP, POP, POP, 5, 6, 7},
			qout: []T{3, 1, 2, 3, 4, 5, 6, 7},
			sout: []T{3, 4, 3, 2, 7, 6, 5, 1},
		},
		{
			size: 4,
			cmds: []T{1, 2, 3, 4, POP, POP, POP, POP, 5, 6, 7, 8, POP, POP, 9, 10},
			qout: []T{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			sout: []T{4, 3, 2, 1, 8, 7, 10, 9, 6, 5},
		},
		// Grow
		{
			size: 4,
			cmds: []T{1, 2, 3, 4, LEN, CAP, 5, LEN, CAP},
			qout: []T{4, 4, 5, 8, 1, 2, 3, 4, 5},
			sout: []T{4, 4, 5, 8, 5, 4, 3, 2, 1},
		},
		{
			size: 4,
			cmds: []T{1, 2, CAP, GROW - 4, CAP, 3, 4, 5, 6, LEN},
			qout: []T{4, 8, 6, 1, 2, 3, 4, 5, 6},
			sout: []T{4, 8, 6, 6, 5, 4, 3, 2, 1},
		},
		{
			size: 4,
			cmds: []T{CAP, GROW - 5, CAP},
			qout: []T{4, 16},
			sout: []T{4, 16},
		},
		{
			size: 4,
			cmds: []T{CAP, GROW - 0, CAP},
			qout: []T{4, 4},
			sout: []T{4, 4},
		},
		// Add slices
		{
			size: 4,
			cmds: []T{[]T{}, LEN},
			qout: []T{0},
			sout: []T{0},
		},
		{
			size: 4,
			cmds: []T{[]T{1, 2, 3, 4}},
			qout: []T{1, 2, 3, 4},
			sout: []T{4, 3, 2, 1},
		},
		{
			size: 2,
			cmds: []T{[]T{1, 2, 3, 4, 5}},
			qout: []T{1, 2, 3, 4, 5},
			sout: []T{5, 4, 3, 2, 1},
		},
		{
			size: 4,
			cmds: []T{[]T{1, 2}, POP, []T{3, 4, 5}},
			qout: []T{1, 2, 3, 4, 5},
			sout: []T{2, 5, 4, 3, 1},
		},
		{
			size: 8,
			cmds: []T{[]T{1, 2, 3, 4}, POP, []T{5, 6, 7, 8}, POP, []T{9, 10}},
			qout: []T{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			sout: []T{4, 8, 10, 9, 7, 6, 5, 3, 2, 1},
		},
		// Reset
		{
			size: 4,
			cmds: []T{1, 2, 3, 4, RESET, LEN, CAP, 5, 6, 7, 8},
			qout: []T{0, 4, 5, 6, 7, 8},
			sout: []T{0, 4, 8, 7, 6, 5},
		},
	}

	for i, row := range data {
		q := NewGenericTypeQueue(row.size)
		s := NewGenericTypeStack(row.size)

		qout := make([]T, 0, len(row.qout)*2)
		sout := make([]T, 0, len(row.qout)*2)

		// Always flush
		row.cmds = append(row.cmds, POPS)

		for _, cmd := range row.cmds {
			switch cmd {
			case POP:
				qout = append(qout, q.Dequeue())
				sout = append(sout, s.Pop())
			case POPS:
				qout = qout[:len(qout)+q.DequeueSlice(qout[len(qout):cap(qout)])]
				sout = sout[:len(sout)+s.PopSlice(sout[len(sout):cap(sout)])]
			case PEEK:
				qout = append(qout, q.Peek())
				sout = append(sout, s.Peek())
			case LEN:
				qout = append(qout, q.Len())
				sout = append(sout, s.Len())
			case CAP:
				qout = append(qout, q.Cap())
				sout = append(sout, s.Cap())
			case RESET:
				q.Reset()
				s.Reset()
			default:
				x, isNum := cmd.(int)
				xs, isSlice := cmd.([]T)

				if isNum && x <= GROW {
					q.Grow(GROW - x)
					s.Grow(GROW - x)
				} else if isSlice {
					q.EnqueueSlice(xs)
					s.PushSlice(xs)
				} else {
					q.Enqueue(x)
					s.Push(x)
				}
			}
		}

		if !reflect.DeepEqual(qout, row.qout) {
			t.Errorf("[%d] qout %v != %v", i, qout, row.qout)
		}

		if !reflect.DeepEqual(sout, row.sout) {
			t.Errorf("[%d] sout %v != %v", i, sout, row.sout)
		}
	}
}
