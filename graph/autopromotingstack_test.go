// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"reflect"
	"testing"
)

func TestAutoPromotingStack(t *testing.T) {
	const PEEK = -1
	const POP = -2
	const PEEKNPS = -3
	const POPNPS = -4
	const TRANSFER = -5

	data := []struct {
		len  int
		cmds []int
		out  []int
	}{
		{
			len:  2,
			cmds: []int{0, 1},
			out:  []int{0, 1},
		},
		{
			len:  2,
			cmds: []int{0, 1, POP},
			out:  []int{1, 0},
		},
		{
			len:  2,
			cmds: []int{0, 1, PEEK},
			out:  []int{1, 0, 1},
		},
		{
			len:  4,
			cmds: []int{0, 1, 2, 1, 3, 3},
			out:  []int{0, 2, 1, 3},
		},
		{
			len:  4,
			cmds: []int{3, 2, 1, TRANSFER, PEEKNPS, 0, TRANSFER, POPNPS},
			out:  []int{1, 0, 3, 2, 1},
		},
		{
			len:  4,
			cmds: []int{0, 1, 2, 3, 0, 1, 2, 3, 2, 1, 0, 0, 1, 2, 3},
			out:  []int{0, 1, 2, 3},
		},
		{
			len:  4,
			cmds: []int{0, 0, 1, 1, 2, 2, 3, 3},
			out:  []int{0, 1, 2, 3},
		},
	}

	for _, row := range data {
		buf := make([]int, row.len*2)
		aps := newAutoPromotingStack(makeListNodeSlice(buf))
		nps := newNonPromotingStack(aps.s)

		for i := range buf {
			buf[i] = -1
		}

		if len(aps.s) != len(buf)/2 {
			t.Errorf("%v != %v", len(aps.s), len(buf)/2)
		}

		out := make([]int, 0, len(row.cmds))

		for _, n := range row.cmds {
			switch n {
			case PEEK:
				out = append(out, aps.Peek())
			case POP:
				out = append(out, aps.Pop())
			case PEEKNPS:
				out = append(out, nps.Peek())
			case POPNPS:
				out = append(out, nps.Pop())
			case TRANSFER:
				nps.Push(aps.Pop())
			default:
				aps.PushOrPromote(n)
			}
		}

		for aps.Len() > 0 {
			nps.Push(aps.Pop())
		}

		for nps.Len() > 0 {
			out = append(out, nps.Pop())
		}

		if !reflect.DeepEqual(out, row.out) {
			t.Errorf("%v != %v", out, row.out)
		}

		for i := range buf {
			buf[i] = 0xff
		}

		for i := range aps.s {
			if aps.s[i].prev != 0xff || aps.s[i].next != 0xff ||
				nps.s[i].prev != 0xff || nps.s[i].next != 0xff {
				t.Log("All values should be 0xff")
				t.Logf("buf: %v", buf)
				t.Logf("aps: %v", aps)
				t.Logf("nps: %v", nps)
				t.Fail()
				break
			}
		}
	}
}
