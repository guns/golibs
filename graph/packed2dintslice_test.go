// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"reflect"
	"testing"
)

func TestPacked2DIntSlice(t *testing.T) {
	type P2D = Packed2DIntSlice

	const new = -1
	n := 0xff

	data := []struct {
		cmds      []int
		out, fill P2D
	}{
		{
			cmds: []int{new, new, new},
			out:  P2D{{}, {}, {}, {}},
			fill: P2D{{}, {}, {}, {}},
		},
		{
			cmds: []int{2, new, 3, new, 4},
			out:  P2D{{2}, {3}, {4}},
			fill: P2D{{n}, {n}, {n}},
		},
		{
			cmds: []int{2, 3, 4, new, 5, 6, 7},
			out:  P2D{{2, 3, 4}, {5, 6, 7}},
			fill: P2D{{n, n, n}, {n, n, n}},
		},
	}

	for _, row := range data {
		p := MakePacked2DIntSlice(len(row.cmds))

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

		// Fill the backing slice
		s := p[0][:cap(p[0])]
		for i := range s {
			s[i] = n
		}

		if !reflect.DeepEqual(p, row.fill) {
			t.Errorf("%v != %v", p, row.fill)
		}

		if p.Cap() != len(row.cmds) {
			t.Errorf("%v != %v", p.Cap(), len(row.cmds))
		}
	}
}
