// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

import (
	"reflect"
	"testing"
)

func TestQueue(t *testing.T) {
	data := []struct {
		size  int
		cmds  []int
		out   []Type
		state TypeQueue
	}{
		{
			size:  1,
			cmds:  []int{1, 0, 2, 0},
			out:   []Type{1, 2},
			state: TypeQueue{a: []Type{2}, head: -1, tail: -1},
		},
		{
			size:  1,
			cmds:  []int{1, 2, 0, 0},
			out:   []Type{1, 2},
			state: TypeQueue{a: []Type{1, 2}, head: -1, tail: -1},
		},
		{
			size:  4,
			cmds:  []int{1, 2, 3, 4, 0, 0, 5, 6, 7, 8, 0, 0},
			out:   []Type{1, 2, 3, 4},
			state: TypeQueue{a: []Type{3, 4, 5, 6, 7, 8, Type(nil), Type(nil)}, head: 2, tail: 6},
		},
		{
			size:  5,
			cmds:  []int{1, 2, 3, 4, 0, 0, 5, 6, 7, 8, 0, 0},
			out:   []Type{1, 2, 3, 4},
			state: TypeQueue{a: []Type{1, 2, 3, 4, 5, 6, 7, 8}, head: 4, tail: 0},
		},
		{
			size:  0,
			cmds:  []int{1, 2, 3, 4, 0, 0, 0, 0},
			out:   []Type{1, 2, 3, 4},
			state: TypeQueue{a: []Type{1, 2, 3, 4, Type(nil), Type(nil), Type(nil), Type(nil)}, head: -1, tail: -1},
		},
		{
			size:  4,
			cmds:  []int{1, 2, -1, -1, 3, 4, 0, 0},
			out:   []Type{1, 1, 1, 2},
			state: TypeQueue{a: []Type{1, 2, 3, 4}, head: 2, tail: 0},
		},
	}

	for i, row := range data {
		q := NewTypeQueue(row.size)
		out := make([]Type, 0, len(row.out))

		for _, n := range row.cmds {
			switch n {
			case -1:
				out = append(out, q.Peek())
			case 0:
				out = append(out, q.Dequeue())
			default:
				q.Enqueue(n)
			}
		}

		if !reflect.DeepEqual(out, row.out) {
			t.Errorf("[%d] %v != %v", i, out, row.out)
		}

		if !reflect.DeepEqual(*q, row.state) {
			t.Errorf("[%d] %v != %v", i, *q, row.state)
		}

		q.Reset()
		if q.Len() != 0 {
			t.Errorf("%v != %v", q.Len(), 0)
		}
	}
}
