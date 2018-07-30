// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

// MinEdgesPath returns a path from vertex u to v with a minimum number of
// edges. The length of the path is the length of the returned path minus one.
//
// The path is written to the path slice, which is grown if necessary.
//
// If no path exists, an empty slice (path[:0) is returned.
//
// Note that trivial paths are not considered; i.e. there is no path from a
// vertex u to itself except through a cycle or self-edge.
func (g Graph) MinEdgesPath(path []int, u, v int, w *Workspace) []int {
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

func writePath(path, pred []int, v int, pathLen int) []int {
	path = resizeIntSlice(path, pathLen+1)
	path[pathLen] = v

	for i := pathLen - 1; i >= 0; i-- {
		v = pred[v]
		path[i] = v
	}

	return path
}
