// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"reflect"
	"testing"
)

func TestAutoPromotingStack(t *testing.T) {
	const peek = -1
	const pop = -2
	const peekNPS = -3
	const popNPS = -4
	const transfer = -5

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
			cmds: []int{0, 1, pop},
			out:  []int{1, 0},
		},
		{
			len:  2,
			cmds: []int{0, 1, peek},
			out:  []int{1, 0, 1},
		},
		{
			len:  4,
			cmds: []int{0, 1, 2, 1, 3, 3},
			out:  []int{0, 2, 1, 3},
		},
		{
			len:  4,
			cmds: []int{3, 2, 1, transfer, peekNPS, 0, transfer, popNPS},
			out:  []int{1, 0, 3, 2, 1},
		},
		{
			len:  4,
			cmds: []int{0, 1, 2, 3, 0, 1, 2, 3, 2, 1, 0, 0, 1, 2, 3},
			out:  []int{0, 1, 2, 3},
		},
	}

	for _, row := range data {
		buf := make([]int, row.len*2)
		aps := newAutoPromotingStack(buf)
		nps := newNonPromotingStack(buf)

		for i := range buf {
			buf[i] = -1
		}

		if len(aps.s) != len(buf)/2 {
			t.Errorf("%v != %v", len(aps.s), len(buf)/2)
		}

		out := make([]int, 0, len(row.cmds))

		for _, n := range row.cmds {
			switch n {
			case peek:
				out = append(out, aps.peek())
			case pop:
				out = append(out, aps.pop())
			case peekNPS:
				out = append(out, nps.peek())
			case popNPS:
				out = append(out, nps.pop())
			case transfer:
				nps.push(aps.pop())
			default:
				aps.pushOrPromote(n)
			}
		}

		for aps.len > 0 {
			nps.push(aps.pop())
		}

		for nps.len > 0 {
			out = append(out, nps.pop())
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
