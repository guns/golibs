// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package graph provides a general-purpose directed graph implementation.
package graph

import "math/bits"

// A Graph is a set { (V, E) : V ⊆ ℤ⁺ and E ⊆ { (u, v) ∈ V } }.
//
// Concretely, a Graph is implemented as an adjacency list of (v, []v), where
// v is a positive unsigned number. Notice that because Graph indices begin
// at 0 and vertices begin at 1, the first entry of a Graph is always empty.
// Accordingly, iterating through a graph's edges should be done like this:
//
//	for u := 1; u < len(g); u++ {
//		for _, v := range g[u] {
//			// Edge (u, v)
//		}
//	}
//
// All edges are directed and unweighted. Undirected graphs can be constructed
// by simply adding the reverse of each edge, and edge weights can be stored
// in a parallel data structure.
//
type Graph [][]uint

// Undefined is an invalid sentinel vertex.
const Undefined = 0

// MakeGraph returns a new Graph with the given logical size.
func MakeGraph(size int) Graph {
	return make(Graph, size+1)
}

// Len returns the logical size of the graph.
func (g Graph) Len() int {
	return len(g) - 1
}

// AddEdge adds a single directed edge from vertex u to v.
func (g Graph) AddEdge(u, v uint) {
	g[u] = append(g[u], v)
	g.Touch(v)
}

// Touch makes a vertex non-nil.
func (g Graph) Touch(u uint) {
	if g[u] == nil {
		g[u] = []uint{}
	}
}

// LeastEdgesPath returns a path from vertex u to v with a minimum number of
// edges. The length of the path is the length of the returned path minus one.
//
// The path is written to the path slice, which is grown if necessary.
//
// If no path exists, an empty slice is returned.
//
// Note that trivial paths are not considered; i.e. there is no path from a
// vertex u to itself except through a cycle or self-edge.
func (g Graph) LeastEdgesPath(path []uint, u, v uint, w *Workspace) []uint {
	w.Prepare(len(g), WA|WB)

	dist := w.a              // |V|w · Slice of vertex -> edge distance from u
	pred := w.b              // |V|w · Slice of vertex -> predecessor vertex
	queue := w.MakeQueue(WC) // |V|w · BFS queue

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
			if pred[v] != Undefined {
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

	if pred[v] == Undefined {
		// No path from u -> v was discovered
		return path[:0]
	}

	return writePath(path, pred, v, int(dist[v]))
}

// TopologicalSort returns a slice of vertex indices in topologically sorted
// order. The offsets are written to the tsort parameter, which is grown if
// necessary. If a topological sort is impossible because there is a cycle in
// the graph, an empty slice (tsort[:0]) is returned.
func (g Graph) TopologicalSort(tsort []uint, w *Workspace) []uint {
	w.Prepare(len(g), 0)

	bs := w.MakeBitsliceN(3, WA)
	active := bs[0]          // |V|  · Bitslice of vertex -> active?
	explored := bs[1]        // |V|  · Bitslice of vertex -> fully explored?
	postorder := bs[2]       // |V|  · Bitslice of stack depth -> children fully explored?
	stack := w.MakeStack(WB) // |V|w · DFS stack

	tsort = resizeUintSlice(tsort, g.Len()) // Prepare write buffer
	idx := len(tsort)                       // tsort write index + 1

	for u := 1; u < len(g); u++ {
		if explored.Get(u) {
			continue
		}

		// DFS
		stack.Push(uint(u))

		// visit(u)
		for stack.Len() > 0 {
			u := stack.Pop()

			// Post-order visit nodes whose children have been explored.
			if postorder.CompareAndClear(stack.Len()) {
				explored.Set(int(u))
				idx--
				tsort[idx] = u
				continue
			}

			if explored.Get(int(u)) {
				// Ignore fully explored nodes
				continue
			} else if active.Get(int(u)) {
				// This neighboring vertex is active but not yet
				// fully explored, so we have discovered a cycle!
				return tsort[:0]
			}

			// Mark this vertex as visited, but not fully explored.
			active.Set(int(u))

			// Postorder visit this parent vertex after all its
			// children have been fully explored.
			postorder.Set(stack.Len())
			stack.Push(u)

			for _, v := range g[u] {
				stack.Push(v)
			}
		}
	}

	return tsort[idx:]
}

// Transpose writes to h a copy of the current graph with all edges reversed.
func (g Graph) Transpose(h Graph) Graph {
	if cap(h) >= len(g) {
		h = h[:len(g)]
	} else {
		h = make(Graph, len(g))
	}

	for u := 1; u < len(g); u++ {
		for _, v := range g[u] {
			h.AddEdge(v, uint(u))
		}
	}

	return h
}

func writePath(path, pred []uint, v uint, pathLen int) []uint {
	path = resizeUintSlice(path, pathLen+1)
	path[pathLen] = v

	for i := pathLen - 1; i >= 0; i-- {
		v = pred[v]
		path[i] = v
	}

	return path
}

func resizeUintSlice(s []uint, size int) []uint {
	if cap(s) >= size {
		return s[:size]
	}
	return make([]uint, size, 1<<uint(bits.Len(uint(size-1))))
}
