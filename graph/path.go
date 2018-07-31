// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

// A Path is a sequence of vertices (v₁, v₂, …, vₙ) that represents a
// path through a graph from v₁ to vₙ.
type Path []int

// EdgeCount returns the number of edges in a path.
func (p Path) EdgeCount() int {
	if len(p) == 0 {
		return 0
	}
	return len(p) - 1
}

// Weight returns the combined edge weights of a path.
func (p Path) Weight(m WeightMapper) float64 {
	n := p.EdgeCount()
	if n == 0 {
		return 0
	}

	var weight float64
	for i := 0; i < n; i++ {
		weight += m.Weight(p[i], p[i+1])
	}
	return weight
}

// MinEdgesPath returns a path from vertex u to v with a minimum number of
// edges. The path is written to Path, which is grown if necessary.
//
// If no path exists, an empty slice (path[:0]) is returned.
//
// Note that trivial paths are not considered; i.e. there is no path from a
// vertex u to itself except through a cycle or self-edge.
func (g Graph) MinEdgesPath(path Path, u, v int, w *Workspace) Path {
	w.reset(len(g), wA|wBNeg)

	dist := w.a              // |V|w · Slice of vertex -> edge distance from u
	pred := w.b              // |V|w · Slice of vertex -> predecessor vertex
	queue := w.makeQueue(wC) // |V|w · BFS queue

	// BFS
	queue.Enqueue(u)

	// If u == v, u is the endpoint, so leave it unvisited.
	target := v
	if u != target {
		pred[u] = u
	}

loop:
	for queue.Len() > 0 {
		u := queue.Dequeue()

		for _, v := range g[u] {
			if pred[v] != undefined {
				continue
			}

			pred[v] = u
			dist[v] = dist[u] + 1

			if v == target {
				break loop
			}

			queue.Enqueue(v)
		}
	}

	if pred[v] == undefined {
		// No path from u -> v was discovered
		return path[:0]
	}

	return writePath(path, pred, v, dist[v])
}

func writePath(path Path, pred []int, v int, pathLen int) Path {
	path = resizeIntSlice(path, pathLen+1)
	path[pathLen] = v

	for i := pathLen - 1; i >= 0; i-- {
		v = pred[v]
		path[i] = v
	}

	return path
}
