// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestPacked2DGenericTypeBuilder(t *testing.T) {
	type T = GenericType

	const FINISH = -1
	const RESET = -2
	const LEN = -3
	const CAP = -4
	const GROW = -5

	data := []struct {
		size int
		cmds []T
		rows [][]T
		out  []T
		grow bool
	}{
		{
			size: 4,
			cmds: []T{FINISH},
			rows: [][]T{{}},
		},
		{
			size: 4,
			cmds: []T{0, 1, LEN, FINISH, LEN, 2, 3, LEN, FINISH, LEN},
			rows: [][]T{{0, 1}, {2, 3}},
			out:  []T{2, 2, 4, 4},
		},
		{
			size: 4,
			cmds: []T{FINISH, FINISH, 0, 1, 2, FINISH, FINISH, 4},
			rows: [][]T{{}, {}, {0, 1, 2}, {}},
		},
		// Grow
		{
			size: 4,
			cmds: []T{1, 2, 3, 4, FINISH, LEN, CAP, 5, 6, 7, 8, FINISH, LEN, CAP},
			rows: [][]T{{1, 2, 3, 4}, {5, 6, 7, 8}},
			out:  []T{4, 4, 8, 8},
			grow: true,
		},
		{
			size: 8,
			cmds: []T{1, 2, FINISH, 3, 4, FINISH, 5, 6, FINISH, 7, 8, FINISH, 9, 10, FINISH},
			rows: [][]T{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}},
			grow: true,
		},
		{
			size: 2,
			cmds: []T{1, 2, GROW - 0, FINISH, GROW - 5, 3, 4, FINISH},
			rows: [][]T{{1, 2}, {3, 4}},
			grow: true,
		},
		// Reset
		{
			size: 4,
			cmds: []T{1, 2, FINISH, 3, 4, FINISH, RESET, 5, 6, 7, 8, FINISH},
			rows: [][]T{{5, 6, 7, 8}},
		},
		{
			size: 4,
			cmds: []T{1, 2, 3, 4, RESET, 5, 6, 7, 8},
			rows: [][]T(nil),
		},
	}

	for i, row := range data {
		p := NewPacked2DGenericTypeBuilder(row.size)
		var out []T

		for _, n := range row.cmds {
			switch n {
			case FINISH:
				p.FinishRow()
			case RESET:
				p.Reset()
			case LEN:
				out = append(out, p.Len())
			case CAP:
				out = append(out, p.Cap())
			default:
				if n.(int) <= GROW {
					p.Grow(GROW - n.(int))
				} else {
					p.Append(n)
				}
			}
		}

		if !reflect.DeepEqual(p.Rows, row.rows) {
			t.Errorf("[%d] %v != %v", i, p.Rows, row.rows)
		}

		if !reflect.DeepEqual(out, row.out) {
			t.Errorf("[%d] %v != %v", i, out, row.out)
		}

		for i := range p.buf {
			p.buf[i] = -1
		}

		if !row.grow {
		loop:
			for i := range p.Rows {
				for j := range p.Rows[i] {
					if p.Rows[i][j] != -1 {
						t.Logf("[%d] p.buf and p.Rows should share memory", i)
						t.Logf("p.buf:  %v", p.buf)
						t.Logf("p.Rows: %v", p.Rows)
						t.Fail()
						break loop
					}
				}
			}
		}
	}
}

func TestPacked2DGenericTypeBuilderFromRows(t *testing.T) {
	p := NewPacked2DGenericTypeBuilder(8)
	p.FinishRow()
	q := NewPacked2DGenericTypeBuilderFromRows(p.Rows)
	pdata := (*reflect.SliceHeader)(unsafe.Pointer(&p.buf)).Data
	qdata := (*reflect.SliceHeader)(unsafe.Pointer(&q.buf)).Data

	if pdata != qdata {
		t.Errorf("%v != %v", pdata, qdata)
	}
}
