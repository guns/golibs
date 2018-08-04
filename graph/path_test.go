// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"reflect"
	"testing"
)

func TestPath(t *testing.T) {
	type pred map[int]int
	type edgew struct {
		src, dst int
		w        float64
	}

	data := []struct {
		size          int
		pred          pred
		weights       []edgew
		defaultWeight float64
		path          []int
		edgeCount     int
		pathWeight    float64
	}{
		{
			size:          0,
			pred:          nil,
			weights:       nil,
			defaultWeight: 0,
			path:          nil,
			edgeCount:     0,
			pathWeight:    0,
		},
		{
			size:          4,
			pred:          pred{1: 0, 2: 1, 3: 2},
			weights:       nil,
			defaultWeight: 0.5,
			path:          []int{0, 1, 2, 3},
			edgeCount:     3,
			pathWeight:    1.5,
		},
		{
			size: 8,
			pred: pred{7: 4, 4: 1, 1: 5, 5: 2},
			weights: []edgew{
				{2, 5, 1},
				{5, 1, 2},
				{4, 7, 8},
			},
			defaultWeight: 4,
			path:          []int{2, 5, 1, 4, 7},
			pathWeight:    15,
			edgeCount:     4,
		},
	}

	for _, row := range data {
		path := Path{}
		path.reset(row.size)

		for u, v := range row.pred {
			path.pred[u] = v
		}

		if row.edgeCount > 0 {
			path.finish(row.path[len(row.path)-1], row.edgeCount)
		}

		m := MakeWeightMap(row.defaultWeight, row.size)

		for _, e := range row.weights {
			m.SetWeight(e.src, e.dst, e.w)
		}

		if path.EdgeCount() != row.edgeCount {
			t.Errorf("%v != %v", path.EdgeCount(), row.edgeCount)
		}
		if len(path.path) != 0 {
			t.Errorf("%v != %v", len(path.path), 0)
		}
		if !reflect.DeepEqual(path.Path(), row.path) {
			t.Errorf("%v != %v", path.Path(), row.path)
		}
		if path.PathWeight(m) != row.pathWeight {
			t.Errorf("%v != %v", path.PathWeight(m), row.pathWeight)
		}
	}
}

func TestGraphMinEdgesPath(t *testing.T) {
	type Adj = map[int][]int

	data := []struct {
		size     int
		adj      Adj
		src, dst int
		path     []int
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
			src:  0,
			dst:  3,
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
			src:  0,
			dst:  3,
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
			src:  0,
			dst:  3,
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
			src:  0,
			dst:  3,
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
			src:  0,
			dst:  0,
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
			src:  0,
			dst:  0,
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
			src:  0,
			dst:  5,
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
			src:  1,
			dst:  0,
			path: []int{1, 9, 8, 0}, // also {1, 2, 8, 0}
		},
	}

	var g Graph
	w := &Workspace{}
	path := &Path{}

	for _, row := range data {
		g = g.Reset(row.size)

		for u, edges := range row.adj {
			for _, v := range edges {
				g.AddEdge(u, v)
			}
		}

		err := g.MinEdgesPath(path, row.src, row.dst, w)

		if !reflect.DeepEqual(path.Path(), row.path) {
			t.Errorf("%v != %v", path.Path(), row.path)
		}

		if err == nil && path.EdgeCount() == 0 {
			t.Errorf("expected err to be nil, but have %v", err)
		} else if err != nil && path.EdgeCount() > 0 {
			t.Error("expected err to be non-nil")
		}
	}
}
