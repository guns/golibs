package graph

import (
	"reflect"
	"testing"

	"github.com/guns/golibs/generic/impl"
)

func TestGraphLeastEdgesPath(t *testing.T) {
	type Adj = map[int][]int

	data := []struct {
		size int
		adj  Adj
		u, v int
		path []int
	}{
		{
			size: 4,
			adj: Adj{
				0: {1},
				1: {0, 2},
				2: {1, 3},
				3: {2},
			},
			u:    0,
			v:    3,
			path: []int{0, 1, 2, 3},
		},
		{
			size: 4,
			adj: Adj{
				0: {1},
				1: {0, 2, 3},
				2: {1, 3},
				3: {2},
			},
			u:    0,
			v:    3,
			path: []int{0, 1, 3},
		},
		// No path
		{
			size: 4,
			adj: Adj{
				0: {1},
				1: {2},
				2: {},
				3: {0},
			},
			u:    0,
			v:    3,
			path: []int{},
		},
		// Cycle
		{
			size: 4,
			adj: Adj{
				0: {1},
				1: {0, 2},
				2: {1, 3},
				3: {2},
			},
			u:    0,
			v:    0,
			path: []int{0, 1, 0},
		},
		// Self-loop
		{
			size: 4,
			adj: Adj{
				0: {0, 1},
				1: {0, 2},
				2: {1, 3},
				3: {2},
			},
			u:    0,
			v:    0,
			path: []int{0, 0},
		},
		{
			size: 10,
			adj: Adj{
				0: {8},
				1: {3, 7, 9, 2},
				2: {8, 1, 4},
				3: {4, 5, 1},
				4: {2, 3},
				5: {3, 6},
				6: {7, 5},
				7: {1, 6},
				8: {2, 0, 9},
				9: {1, 8},
			},
			u:    0,
			v:    5,
			path: []int{0, 8, 2, 1, 3, 5},
		},
		{
			size: 10,
			adj: Adj{
				0: {8},
				1: {3, 7, 9, 2},
				2: {8, 1, 4},
				3: {4, 5, 1},
				4: {2, 3},
				5: {3, 6},
				6: {7, 5},
				7: {1, 6},
				8: {2, 0, 9},
				9: {1, 8},
			},
			u:    1,
			v:    0,
			path: []int{1, 9, 8, 0}, // also {1, 2, 8, 0}
		},
	}

	w := NewWorkspace(0)

	for _, row := range data {
		g := make(Graph, row.size)

		for u, edges := range row.adj {
			for _, v := range edges {
				g.AddEdge(u, v)
			}
		}

		path := g.LeastEdgesPath([]int{}, row.u, row.v, w)

		if !reflect.DeepEqual(path, row.path) {
			t.Errorf("%v != %v", path, row.path)
		}
	}
}

func TestGraphTopologicalSort(t *testing.T) {
	type Adj = map[int][]int

	data := []struct {
		size   int
		adj    Adj
		cyclic bool
	}{
		{
			size: 8,
			adj: Adj{
				0: {1, 2},
				1: {3, 4},
				2: {3},
				3: {5},
				4: {6},
				5: {6, 7},
				6: {},
				7: {6},
			},
		},
		// Figure 22.7 CLRS (topologically sorted clothes)
		{
			size: 9,
			adj: Adj{
				0: {3, 4},
				1: {4},
				2: {},
				3: {4, 5},
				4: {},
				5: {8},
				6: {5, 7},
				7: {8},
				8: {},
			},
		},
		// Unconnected vertices
		{
			size: 10,
			adj:  Adj{},
		},
		// Sparse graph
		{
			size: 4,
			adj: Adj{
				0: {},
				1: {3},
				2: {},
				3: {},
			},
		},
		// Cyclic graphs
		{
			size: 4,
			adj: Adj{
				0: {1},
				1: {2},
				2: {3},
				3: {0},
			},
			cyclic: true,
		},
		{
			size: 4,
			adj: Adj{
				0: {1},
				1: {2},
				2: {2, 3},
				3: {},
			},
			cyclic: true,
		},
	}

	w := NewWorkspace(0)

	for _, row := range data {
		g := make(Graph, row.size)

		for u, vs := range row.adj {
			for _, v := range vs {
				g.AddEdge(u, v)
			}
		}

		tsort := g.TopologicalSort(make([]int, 0, row.size), w)

		if row.cyclic {
			if len(tsort) != 0 {
				t.Errorf("%v != %v", len(tsort), 0)
			}
			continue
		}

		rsort := make([]int, len(g))

		// Create a reverse mapping for easy lookup
		for i, j := range tsort {
			rsort[j] = i
		}

		// CLRS: A topological sort of a dag G = (V,E) is a linear
		// ordering of all its vertices such that if G contains an
		// edge (u,v), then u appears before v in the ordering.
		for i, u := range tsort {
			for _, v := range g[u] {
				j := rsort[v]
				if j <= i {
					t.Errorf("edge (%v,%v) out of order in %v\n", u, v, tsort)
				}
			}
		}

		if len(tsort) != len(g) {
			t.Errorf("len(tsort) %v != len(g) %v", len(tsort), len(g))
		}

		impl.QuicksortIntSlice(tsort)

		equal := true
		for u := range g {
			if tsort[u] != u {
				equal = false
			}
		}

		if !equal {
			t.Errorf("tsort: %v does not contain all graph vertices", tsort)
		}

	}
}

func TestGraphTranspose(t *testing.T) {
	data := []struct {
		size  int
		edges [][]int
	}{
		{
			size:  0,
			edges: nil,
		},
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
	}

	for _, row := range data {
		g := make(Graph, row.size)
		h := make(Graph, row.size)
		gT := make(Graph, row.size)

		for _, e := range row.edges {
			g.AddEdge(e[0], e[1])
			h.AddEdge(e[1], e[0])
		}

		gT = g.Transpose(gT)
		if !reflect.DeepEqual(gT, h) {
			t.Errorf("%v != %v", gT, h)
		}
	}
}
