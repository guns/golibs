// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"sort"
	"testing"
)

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
