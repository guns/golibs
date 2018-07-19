// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"reflect"
	"testing"
)

func TestPacked2DIntSlice(t *testing.T) {
	data := []struct {
		cmds       []int
		out, clear Packed2DIntSlice
	}{
		{
			cmds:  []int{0, 0, 0},
			out:   Packed2DIntSlice{{}, {}, {}, {}},
			clear: Packed2DIntSlice{{}, {}, {}, {}},
		},
		{
			cmds:  []int{1, 0, 2, 0, 3},
			out:   Packed2DIntSlice{{1}, {2}, {3}},
			clear: Packed2DIntSlice{{-1}, {-1}, {-1}},
		},
		{
			cmds:  []int{1, 2, 3, 0, 4, 5, 6},
			out:   Packed2DIntSlice{{1, 2, 3}, {4, 5, 6}},
			clear: Packed2DIntSlice{{-1, -1, -1}, {-1, -1, -1}},
		},
	}

	for _, row := range data {
		p := MakePacked2DIntSlice(len(row.cmds))

		for _, cmd := range row.cmds {
			switch cmd {
			case 0:
				p = p.StartNewSlice()
			default:
				p = p.Append(cmd)
			}
		}

		if !reflect.DeepEqual(p, row.out) {
			t.Errorf("%v != %v", p, row.out)
		}

		// Fill the backing slice with -1
		s := p[0][:cap(p[0])]
		for i := range s {
			s[i] = -1
		}

		if !reflect.DeepEqual(p, row.clear) {
			t.Errorf("%v != %v", p, row.clear)
		}

		if p.Cap() != len(row.cmds) {
			t.Errorf("%v != %v", p.Cap(), len(row.cmds))
		}
	}
}
