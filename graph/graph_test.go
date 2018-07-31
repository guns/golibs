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
	type Adj map[int][]int

	data := []struct {
		size int
		adj  Adj
	}{
		{
			size: 2,
			adj: Adj{
				0: {1},
			},
		},
		{
			size: 5,
			adj: Adj{
				0: {1},
				1: {2},
				2: {0, 3},
				3: {1, 2},
				4: {4},
			},
		},
		{
			size: 3,
			adj: Adj{
				0: {1},
				1: {2},
				2: {0},
			},
		},
		{
			size: 0,
			adj:  nil,
		},
	}

	var g, h, gT Graph

	for i, row := range data {
		g = g.Reset(row.size)
		h = h.Reset(row.size)

		for u, vs := range row.adj {
			for _, v := range vs {
				g.AddEdge(u, v)
				h.AddEdge(v, u)
			}
		}

		gT = g.Transpose(gT)

		// Prep edge lists for equality testing
		for u := range gT {
			if h[u] == nil {
				h[u] = []int{}
			}
			if gT[u] == nil {
				gT[u] = []int{}
			}

			sort.Ints(h[u])
			sort.Ints(gT[u])
		}

		if !reflect.DeepEqual(gT, h) {
			t.Errorf("[%d] %v != %v", i, gT, h)
		}
	}
}
