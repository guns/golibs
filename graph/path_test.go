// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"reflect"
	"testing"
)

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
