package graph

import (
	"reflect"
	"sort"
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

func TestGraphLeastEdgesPath(t *testing.T) {
	type Adj = map[int][]int

	data := []struct {
		size int
		adj  Adj
		u, v int
		path []int
	}{
		// Attempt to overflow queue (should be first case)
		{
			size: 4,
			adj: Adj{
				0: {0, 1, 2, 3},
				1: {},
				2: {},
				3: {},
			},
			u:    0,
			v:    3,
			path: []int{0, 3},
		},
		// Typical case
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

	var g Graph
	w := &Workspace{}
	var path []int

	for _, row := range data {
		g = g.Resize(row.size)

		for u, edges := range row.adj {
			for _, v := range edges {
				g.AddEdge(u, v)
			}
		}

		path = g.LeastEdgesPath(path, row.u, row.v, w)

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
		// Attempt to overflow stack (should be first case)
		{
			size: 4,
			adj: Adj{
				0: {1},
				1: {2},
				2: {0, 1, 2, 3},
				3: {},
			},
			cyclic: true,
		},
		// Typical case
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

	var g Graph
	w := &Workspace{}
	var tsort []int

	for i, row := range data {
		g = g.Resize(row.size)

		for u, vs := range row.adj {
			for _, v := range vs {
				g.AddEdge(u, v)
			}
		}

		tsort = g.TopologicalSort(tsort, w)

		if row.cyclic {
			if len(tsort) != 0 {
				t.Errorf("[%d] %v != %v", i, len(tsort), 0)
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
					t.Errorf("[%d] edge (%v,%v) out of order in %v\n", i, u, v, tsort)
				}
			}
		}

		if len(tsort) != len(g) {
			t.Errorf("[%d] len(tsort) %v != len(g) %v", i, len(tsort), len(g))
		}

		sort.Ints(tsort)

		equal := true
		for u := range g {
			if tsort[u] != u {
				equal = false
			}
		}

		if !equal {
			t.Errorf("[%d] tsort: %v does not contain all graph vertices", i, tsort)
		}

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
		g = g.Resize(row.size)
		h = h.Resize(row.size)

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

func TestGraphStronglyConnectedComponents(t *testing.T) {
	data := []struct {
		size int
		adj  map[int][]int
		scc  [][]int
	}{
		{
			size: 8,
			adj: map[int][]int{
				0: {4},
				1: {0},
				2: {1, 3},
				3: {2},
				4: {1},
				5: {1, 4, 6},
				6: {2, 5},
				7: {3, 6, 7},
			},
			scc: [][]int{
				{0, 1, 4},
				{2, 3},
				{5, 6},
				{7},
			},
		},
		{
			size: 11,
			adj: map[int][]int{
				0:  {1, 3},
				1:  {2, 5},
				2:  {8},
				3:  {4, 6, 9},
				4:  {5, 7, 10},
				5:  {2},
				6:  {0},
				7:  {9},
				8:  {5},
				9:  {4},
				10: {7, 8},
			},
			scc: [][]int{
				{2, 5, 8},
				{4, 7, 9, 10},
				{1},
				{0, 3, 6},
			},
		},
	}

	var g Graph
	w := &Workspace{}
	var scc [][]int

	for _, row := range data {
		g = g.Resize(row.size)

		for u, edges := range row.adj {
			for _, v := range edges {
				g.AddEdge(u, v)
			}
		}

		scc = g.StronglyConnectedComponents(scc, w)

		for i := range scc {
			sort.Ints(scc[i])
		}

		if !reflect.DeepEqual(scc, row.scc) {
			t.Errorf("%v != %v", scc, row.scc)
		}
	}
}
