package gen

import (
	"reflect"
	"testing"
)

func TestStack(t *testing.T) {
	data := []struct {
		size  int
		cmds  []int
		out   []Type
		state TypeStack
	}{
		{
			size:  1,
			cmds:  []int{1, 0, 2, 0},
			out:   []Type{1, 2},
			state: TypeStack{a: []Type{2}, i: -1},
		},
		{
			size:  1,
			cmds:  []int{1, 2, 0, 0},
			out:   []Type{2, 1},
			state: TypeStack{a: []Type{1, 2}, i: -1},
		},
		{
			size:  4,
			cmds:  []int{1, 2, 3, 4, 0, 0, 5, 6, 7, 8, 0, 0},
			out:   []Type{4, 3, 8, 7},
			state: TypeStack{a: []Type{1, 2, 5, 6, 7, 8, Type(nil), Type(nil)}, i: 3},
		},
		{
			size:  5,
			cmds:  []int{1, 2, 3, 4, 0, 0, 5, 6, 7, 8, 0, 0},
			out:   []Type{4, 3, 8, 7},
			state: TypeStack{a: []Type{1, 2, 5, 6, 7, 8, Type(nil), Type(nil)}, i: 3},
		},
		{
			size:  4,
			cmds:  []int{1, 2, -1, -1, 3, 4, 0, 0},
			out:   []Type{2, 2, 4, 3},
			state: TypeStack{a: []Type{1, 2, 3, 4}, i: 1},
		},
	}

	for i, row := range data {
		s := NewTypeStack(row.size)
		out := make([]Type, 0, len(row.out))

		for _, n := range row.cmds {
			switch n {
			case -1:
				out = append(out, s.Peek())
			case 0:
				out = append(out, s.Pop())
			default:
				s.Push(n)
			}
		}

		if !reflect.DeepEqual(out, row.out) {
			t.Errorf("[%d] %v != %v", i, out, row.out)
		}

		if !reflect.DeepEqual(*s, row.state) {
			t.Errorf("[%d] %v != %v", i, s, row.state)
		}

		s.Reset()
		if s.Len() != 0 {
			t.Errorf("%v != %v", s.Len(), 0)
		}
	}
}
