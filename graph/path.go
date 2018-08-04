// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

// A Path is a sequence of vertices (v₁, v₂, …, vₙ) that represents a path
// through a graph's edges from v₁ to vₙ. It is lazily constructed from
// a predecessor table to avoid the cost of actually computing the path when
// only the edge count of the path is desired.
type Path struct {
	path        []int
	pred        []int
	endVertex   int
	edgeCount   int
	initialized bool
}

func (p *Path) reset(size int) {
	p.pred = fillUndefined(resizeIntSlice(p.pred, size))
	p.endVertex = undefined
	p.edgeCount = 0
	p.initialized = false
}

func (p *Path) finish(endVertex, edgeCount int) {
	p.endVertex = endVertex
	p.edgeCount = edgeCount
}

// Path returns the path v₁ -> vₙ as []int{v₁, v₂, …, vₙ}. This method
// constructs and caches the path from a predecessor table on first call.
func (p *Path) Path() []int {
	if p.initialized {
		return p.path
	}

	if p.endVertex == undefined {
		p.path = p.path[:0]
		p.initialized = true
		return p.path
	}

	p.path = resizeIntSlice(p.path, p.edgeCount+1)
	p.path[p.edgeCount] = p.endVertex

	v := p.endVertex

	for i := p.edgeCount - 1; i >= 0; i-- {
		v = p.pred[v]
		p.path[i] = v
	}

	p.initialized = true
	return p.path
}

// EdgeCount returns the number of edges in this path.
func (p *Path) EdgeCount() int {
	return p.edgeCount
}

// PathWeight returns the combined edge weights of this path.
func (p *Path) PathWeight(m WeightMapper) float64 {
	n := p.EdgeCount()
	if n == 0 {
		return 0
	}

	path := p.Path()
	var weight float64

	for i := 0; i < n; i++ {
		weight += m.Weight(path[i], path[i+1])
	}

	return weight
}

// MinEdgesPath searches for a path from src to dst with a minimum number of
// edges. If no such path exists, a non-nil error is returned. The path is
// written to the path parameter.
//
// Note that trivial paths are not considered; i.e. there is no path from a
// vertex src to itself except through a cycle or self-edge.
//
// Worst-case time: O(|V| + |E|)
func (g Graph) MinEdgesPath(path *Path, src, dst int, w *Workspace) error {
	w.reset(len(g), wB)
	path.reset(len(g))

	queue := w.queue(wA) // |V|w · BFS queue
	dist := w.b          // |V|w · Slice of vertex -> edge distance from src
	pred := path.pred    // |V|w · Slice of vertex -> predecessor vertex

	// BFS
	queue.Enqueue(src)

	// If src == dst, src is the endpoint, so leave it unvisited.
	if src != dst {
		pred[src] = src
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

			if v == dst {
				break loop
			}

			queue.Enqueue(v)
		}
	}

	if pred[dst] == undefined {
		// No path from u -> v was discovered
		return errNoPath
	}

	path.finish(dst, dist[dst])

	return nil
}
