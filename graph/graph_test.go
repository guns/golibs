package graph

import (
	"reflect"
	"sort"
	"testing"
)

func TestGraphLeastEdgesPath(t *testing.T) {
	data := []struct {
		size int
		adj  map[int][]int // For readability
		u, v int
		path []int
	}{
		{
			size: 4,
			adj: map[int][]int{
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
			adj: map[int][]int{
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
			adj: map[int][]int{
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
			adj: map[int][]int{
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
			adj: map[int][]int{
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
			adj: map[int][]int{
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
			adj: map[int][]int{
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

	w := NewWorkspace(4)

	for _, row := range data {
		g := make(Graph, row.size)

		for u, edges := range row.adj {
			for _, v := range edges {
				g.AddEdge(u, v, 1)
			}
		}

		path := g.LeastEdgesPath([]int{}, row.u, row.v, w)

		if !reflect.DeepEqual(path, row.path) {
			t.Errorf("%v != %v", path, row.path)
		}
	}
}

func TestGraphTopologicalSort(t *testing.T) {
	data := []struct {
		size   int
		adj    map[int][]int // For readability
		cyclic bool
	}{
		{
			size: 8,
			adj: map[int][]int{
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
			adj: map[int][]int{
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
			adj:  map[int][]int{},
		},
		// Sparse graph
		{
			size: 4,
			adj: map[int][]int{
				0: {},
				1: {3},
				2: {},
				3: {},
			},
		},
		// Cyclic graphs
		{
			size: 4,
			adj: map[int][]int{
				0: {1},
				1: {2},
				2: {3},
				3: {0},
			},
			cyclic: true,
		},
		{
			size: 4,
			adj: map[int][]int{
				0: {1},
				1: {2},
				2: {2, 3},
				3: {},
			},
			cyclic: true,
		},
	}

	w := NewWorkspace(4)

	for _, row := range data {
		g := make(Graph, row.size)

		for u, adj := range row.adj {
			for _, v := range adj {
				g.AddEdge(u, v, 1)
			}
		}

		tsort := g.TopologicalSort(make([]int, 0, row.size), w)

		if row.cyclic {
			if len(tsort) != 0 {
				t.Errorf("%v != %v", len(tsort), 0)
			}
			continue
		}

		if len(tsort) != len(g) {
			t.Errorf("len(tsort) %v != len(g) %v", len(tsort), len(g))
		}

		copy(w.a, tsort)
		sort.Ints(w.a)

		equal := true
		for i := range g {
			if w.a[i] != i {
				equal = false
			}
		}

		if !equal {
			t.Errorf("tsort: %v does not contain all graph vertices", tsort)
		}

		// Create a reverse mapping for easy lookup
		for i, j := range tsort {
			w.a[j] = i
		}

		// A topological sort of a dag G = (V,E) is a linear ordering
		// of all its vertices such that if G contains an edge (u,v),
		// then u appears before v in the ordering.
		for i, u := range tsort {
			for _, e := range g[u].Edges {
				j := w.a[e.Vertex]
				if j <= i {
					t.Errorf("edge (%v,%v) out of order in %v\n", u, e.Vertex, tsort)
				}
			}
		}
	}
}

func TestGraphTranspose(t *testing.T) {
	type edge struct {
		u, v int
		w    float64
	}
	data := []struct {
		size  int
		edges []edge
	}{
		{size: 0, edges: nil},
		{
			size:  2,
			edges: []edge{{0, 1, 1}},
		},
		{
			size: 5,
			edges: []edge{
				{0, 1, 2.1},
				{1, 2, 3.2},
				{2, 0, 0.7},
				{2, 3, 4.1},
				{3, 2, 0.2},
				{3, 1, 3.6},
				{4, 4, 0.5},
			},
		},
	}

	for _, row := range data {
		g := make(Graph, row.size)
		h := make(Graph, row.size)
		gT := make(Graph, row.size)

		for _, e := range row.edges {
			g.AddEdge(e.u, e.v, e.w)
			h.AddEdge(e.v, e.u, e.w)
		}

		gT = g.Transpose(gT)
		if !reflect.DeepEqual(gT, h) {
			t.Errorf("%v != %v", gT, h)
		}
	}
}
