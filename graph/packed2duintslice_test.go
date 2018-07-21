// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"reflect"
	"testing"
)

func TestPacked2DUintSlice(t *testing.T) {
	type p2d = Packed2DUintSlice

	const new = 0
	const max = ^uint(0)

	data := []struct {
		cmds       []uint
		out, clear p2d
	}{
		{
			cmds:  []uint{new, new, new},
			out:   p2d{{}, {}, {}, {}},
			clear: p2d{{}, {}, {}, {}},
		},
		{
			cmds:  []uint{2, new, 3, new, 4},
			out:   p2d{{2}, {3}, {4}},
			clear: p2d{{max}, {max}, {max}},
		},
		{
			cmds:  []uint{2, 3, 4, new, 5, 6, 7},
			out:   p2d{{2, 3, 4}, {5, 6, 7}},
			clear: p2d{{max, max, max}, {max, max, max}},
		},
	}

	for _, row := range data {
		p := MakePacked2DUintSlice(len(row.cmds))

		for _, cmd := range row.cmds {
			switch cmd {
			case new:
				p = p.StartNewSlice()
			default:
				p = p.Append(cmd)
			}
		}

		if !reflect.DeepEqual(p, row.out) {
			t.Errorf("%v != %v", p, row.out)
		}

		// Fill the backing slice with ^0
		s := p[0][:cap(p[0])]
		for i := range s {
			s[i] = max
		}

		if !reflect.DeepEqual(p, row.clear) {
			t.Errorf("%v != %v", p, row.clear)
		}

		if p.Cap() != len(row.cmds) {
			t.Errorf("%v != %v", p.Cap(), len(row.cmds))
		}
	}
}
