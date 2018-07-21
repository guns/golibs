package graph

import (
	"reflect"
	"testing"

	"github.com/guns/golibs/generic/impl"
)

func TestGraphLeastEdgesPath(t *testing.T) {
	type Adj = map[uint][]uint

	data := []struct {
		size int
		adj  Adj
		u, v uint
		path []uint
	}{
		{
			size: 4,
			adj: Adj{
				1: {2},
				2: {1, 3},
				3: {2, 4},
				4: {3},
			},
			u:    1,
			v:    4,
			path: []uint{1, 2, 3, 4},
		},
		{
			size: 4,
			adj: Adj{
				1: {2},
				2: {1, 3, 4},
				3: {2, 4},
				4: {3},
			},
			u:    1,
			v:    4,
			path: []uint{1, 2, 4},
		},
		// No path
		{
			size: 4,
			adj: Adj{
				1: {2},
				2: {3},
				3: {},
				4: {1},
			},
			u:    1,
			v:    4,
			path: []uint{},
		},
		// Cycle
		{
			size: 4,
			adj: Adj{
				1: {2},
				2: {1, 3},
				3: {2, 4},
				4: {3},
			},
			u:    1,
			v:    1,
			path: []uint{1, 2, 1},
		},
		// Self-loop
		{
			size: 4,
			adj: Adj{
				1: {1, 2},
				2: {1, 3},
				3: {2, 4},
				4: {3},
			},
			u:    1,
			v:    1,
			path: []uint{1, 1},
		},
		{
			size: 10,
			adj: Adj{
				1:  {9},
				2:  {4, 8, 10, 3},
				3:  {9, 2, 5},
				4:  {5, 6, 2},
				5:  {3, 4},
				6:  {4, 7},
				7:  {8, 6},
				8:  {2, 7},
				9:  {3, 1, 10},
				10: {2, 9},
			},
			u:    1,
			v:    6,
			path: []uint{1, 9, 3, 2, 4, 6},
		},
		{
			size: 10,
			adj: Adj{
				1:  {9},
				2:  {4, 8, 10, 3},
				3:  {9, 2, 5},
				4:  {5, 6, 2},
				5:  {3, 4},
				6:  {4, 7},
				7:  {8, 6},
				8:  {2, 7},
				9:  {3, 1, 10},
				10: {2, 9},
			},
			u:    2,
			v:    1,
			path: []uint{2, 10, 9, 1}, // also {2, 3, 9, 1}
		},
	}

	w := NewWorkspace(0)

	for _, row := range data {
		g := MakeGraph(row.size)

		for u, edges := range row.adj {
			for _, v := range edges {
				g.AddEdge(u, v)
			}
		}

		path := g.LeastEdgesPath([]uint{}, row.u, row.v, w)

		if !reflect.DeepEqual(path, row.path) {
			t.Errorf("%v != %v", path, row.path)
		}
	}
}

func TestGraphTopologicalSort(t *testing.T) {
	type Adj = map[uint][]uint

	data := []struct {
		size   int
		adj    Adj
		cyclic bool
	}{
		{
			size: 8,
			adj: Adj{
				1: {2, 3},
				2: {4, 5},
				3: {4},
				4: {6},
				5: {7},
				6: {7, 8},
				7: {},
				8: {7},
			},
		},
		// Figure 22.7 CLRS (topologically sorted clothes)
		{
			size: 9,
			adj: Adj{
				1: {4, 5},
				2: {5},
				3: {},
				4: {5, 6},
				5: {},
				6: {9},
				7: {6, 8},
				8: {9},
				9: {},
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
				1: {},
				2: {4},
				3: {},
				4: {},
			},
		},
		// Cyclic graphs
		{
			size: 4,
			adj: Adj{
				1: {2},
				2: {3},
				3: {4},
				4: {1},
			},
			cyclic: true,
		},
		{
			size: 4,
			adj: Adj{
				1: {2},
				2: {3},
				3: {3, 4},
				4: {},
			},
			cyclic: true,
		},
	}

	w := NewWorkspace(0)

	for _, row := range data {
		g := MakeGraph(row.size)

		for u, vs := range row.adj {
			for _, v := range vs {
				g.AddEdge(u, v)
			}
		}

		tsort := g.TopologicalSort(make([]uint, 0, row.size), w)

		if row.cyclic {
			if len(tsort) != 0 {
				t.Errorf("%v != %v", len(tsort), 0)
			}
			continue
		}

		rsort := make([]uint, len(g))

		// Create a reverse mapping for easy lookup
		for i, j := range tsort {
			rsort[j] = uint(i)
		}

		// A topological sort of a dag G = (V,E) is a linear ordering of all its
		// vertices such that if G contains an edge (u,v), then u appears before v
		// in the ordering.
		for i, u := range tsort {
			for _, v := range g[u] {
				j := rsort[v]
				if j <= uint(i) {
					t.Errorf("edge (%v,%v) out of order in %v\n", u, v, tsort)
				}
			}
		}

		if len(tsort) != g.Len() {
			t.Errorf("len(tsort) %v != g.Len() %v", len(tsort), g.Len())
		}

		impl.QuicksortUintSlice(tsort)

		equal := true
		for u := 1; u < len(g); u++ {
			if tsort[u-1] != uint(u) {
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
		edges [][]uint
	}{
		{
			size:  0,
			edges: nil,
		},
		{
			size: 2,
			edges: [][]uint{
				{1, 2},
			},
		},
		{
			size: 5,
			edges: [][]uint{
				{1, 2},
				{2, 3},
				{3, 1},
				{3, 4},
				{4, 3},
				{4, 2},
				{5, 5},
			},
		},
	}

	for _, row := range data {
		g := MakeGraph(row.size)
		h := MakeGraph(row.size)
		gT := MakeGraph(row.size)

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
