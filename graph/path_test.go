// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"reflect"
	"testing"
)

func TestPath(t *testing.T) {
	type edgew struct {
		u, v int
		w    float64
	}
	data := []struct {
		path          Path
		weights       []edgew
		defaultWeight float64
		weight        float64
		edgeCount     int
	}{
		{
			path:          Path{},
			weights:       nil,
			defaultWeight: 0,
			weight:        0,
			edgeCount:     0,
		},
		{
			path:          Path{0, 1, 2, 3},
			weights:       nil,
			defaultWeight: 0.5,
			weight:        1.5,
			edgeCount:     3,
		},
		{
			path: Path{2, 5, 1, 4, 7},
			weights: []edgew{
				{2, 5, 1},
				{5, 1, 2},
				{4, 7, 8},
			},
			defaultWeight: 4,
			weight:        15,
			edgeCount:     4,
		},
	}

	for _, row := range data {
		m := MakeWeightMap(row.defaultWeight, row.path.EdgeCount())

		for _, e := range row.weights {
			m.SetWeight(e.u, e.v, e.w)
		}

		if row.path.Weight(m) != row.weight {
			t.Errorf("%v != %v", row.path.Weight(m), row.weight)
		}
		if row.path.EdgeCount() != row.edgeCount {
			t.Errorf("%v != %v", row.path.EdgeCount(), row.edgeCount)
		}
	}
}

func TestGraphMinEdgesPath(t *testing.T) {
	type Adj = map[int][]int

	data := []struct {
		size int
		adj  Adj
		u, v int
		path Path
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
			path: Path{0, 3},
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
			path: Path{0, 1, 2, 3},
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
			path: Path{0, 1, 3},
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
			path: Path{},
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
			path: Path{0, 1, 0},
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
			path: Path{0, 0},
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
			path: Path{0, 8, 2, 1, 3, 5},
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
			path: Path{1, 9, 8, 0}, // also {1, 2, 8, 0}
		},
	}

	var g Graph
	w := &Workspace{}
	var path Path

	for _, row := range data {
		g = g.Reset(row.size)

		for u, edges := range row.adj {
			for _, v := range edges {
				g.AddEdge(u, v)
			}
		}

		path = g.MinEdgesPath(path, row.u, row.v, w)

		if !reflect.DeepEqual(path, row.path) {
			t.Errorf("%v != %v", path, row.path)
		}
	}
}
