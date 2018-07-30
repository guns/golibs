package graph

import (
	"reflect"
	"testing"
)

func TestGraphGrow(t *testing.T) {
	var g Graph
	var adj [][]int

	for i := 1; i < 41; i++ {
		g = g.Grow(1)
		g.AddEdge(i-1, i)
		adj = append(adj, []int{i})
	}

	if !reflect.DeepEqual(g, Graph(adj)) {
		t.Logf("%v !=", g)
		t.Logf("%v", adj)
		t.Fail()
	}
}

func TestGraphTranspose(t *testing.T) {
	data := []struct {
		size  int
		edges [][]int
	}{
		{
			size: 2,
			edges: [][]int{
				{0, 1},
			},
		},
		{
			size: 5,
			edges: [][]int{
				{0, 1},
				{1, 2},
				{2, 0},
				{2, 3},
				{3, 2},
				{3, 1},
				{4, 4},
			},
		},
		{
			size: 3,
			edges: [][]int{
				{0, 1},
				{1, 2},
				{2, 0},
			},
		},
		{
			size:  0,
			edges: nil,
		},
	}

	var g, h, gT Graph

	for i, row := range data {
		g = g.Reset(row.size)
		h = h.Reset(row.size)

		for _, e := range row.edges {
			g.AddEdge(e[0], e[1])
			h.AddEdge(e[1], e[0])
		}

		gT = g.Transpose(gT)

		// Create empty edge lists for equality testing
		for i := range gT {
			if h[i] == nil {
				h[i] = []int{}
			}
			if gT[i] == nil {
				gT[i] = []int{}
			}
		}

		if !reflect.DeepEqual(gT, h) {
			t.Errorf("[%d] %v != %v", i, gT, h)
		}
	}
}
