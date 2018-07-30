// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"reflect"
	"sort"
	"testing"
)

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
