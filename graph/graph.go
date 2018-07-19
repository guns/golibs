// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package graph provides a general-purpose directed graph implementation.
package graph

import "math/bits"

// An Edge is a single weighted path to a Vertex.
type Edge struct {
	Vertex int
	Weight float64
}

// A Vertex contains a list of Edges.
type Vertex struct {
	Edges []Edge
}

// A Graph is a slice of Vertex structs.
type Graph []Vertex

// AddEdge adds a single Edge from Vertex u to v.
func (g Graph) AddEdge(u, v int, weight float64) {
	g[u].Edges = append(g[u].Edges, Edge{Vertex: v, Weight: weight})
	g.TouchVertex(v)
}

// TouchVertex makes a Vertex non-nil.
func (g Graph) TouchVertex(u int) {
	if g[u].Edges == nil {
		g[u].Edges = []Edge{}
	}
}

// LeastEdgesPath returns a path from Vertex u to v with a minimum number of
// Edges, irrespective of Edge weights. The length of the path in edges is the
// size of the returned path minus one.
//
// The path is written to the path slice, which is grown if necessary.
//
// If no path exists, an empty slice is returned.
func (g Graph) LeastEdgesPath(path []int, u, v int, w *Workspace) []int {
	w.Prepare(len(g), WA|WB, 0, -1)

	dist := w.a      // Edge distances from u
	prev := w.b      // Mapping of vertex -> previous vertex (-1 if unvisited)
	queue := w.queue // BFS queue

	// BFS
	queue.Enqueue(u)

	// If u == v, u is the endpoint, so leave it unvisited.
	if u != v {
		prev[u] = u
	}

loop:
	for queue.Len() > 0 {
		i := queue.Dequeue()

		for _, e := range g[i].Edges {
			w := e.Vertex

			if prev[w] != -1 {
				continue
			}

			prev[w] = i
			dist[w] = dist[i] + 1

			if w == v {
				break loop
			}

			queue.Enqueue(w)
		}
	}

	if prev[v] == -1 {
		// No path from u -> v was discovered
		return path[:0]
	}

	return writePath(path, v, dist[v], prev)
}

// TopologicalSort returns a slice of vertex indices in topologically
// sorted order. The offsets are written to the vs slice, which is grown if
// necessary. If a topological sort is impossible because there is a cycle in
// the graph, an empty slice is returned.
func (g Graph) TopologicalSort(tsort []int, w *Workspace) []int {
	w.Prepare(len(g), WA|WBitslice, 0, 0)

	// We require the following to use an iterative DFS for topologically
	// sorting a directed graph:
	//
	//	- A LIFO queue (stack) of nodes to visit
	//	- A table of active nodes (visited, but not fully explored)
	//	- A table of explored nodes (vertex and descendents fully explored)
	//	- A way to flag a vertex whose children have been fully explored
	//	  (this enables post-order traversal without recursion)

	active := w.a          // Map of vertex -> active?
	explored := w.bitslice // Map of vertex -> fully explored?
	stack := w.stack       // DFS stack

	tsort = resizeIntSlice(tsort, len(g)) // Prepare write buffer
	i := len(g)                           // tsort write index + 1

	for u := range g {
		if explored.Get(u) {
			continue
		}

		// DFS
		stack.Push(u)

		for stack.Len() > 0 {
			v := stack.Pop()

			// Post-order nodes are encoded as their ones' complement
			if v < 0 {
				// Vertex v's children have been explored
				v = ^v
				explored.Set(v)
				i--
				tsort[i] = v
				continue
			} else if explored.Get(v) {
				// Ignore fully explored nodes
				continue
			} else if active[v] == 1 {
				// This neighboring vertex is active but not yet
				// fully explored, so we have discovered a cycle!
				return tsort[:0]
			}

			// When all children have been explored, this parent
			// vertex will appear on top of the stack.
			stack.Push(^v)

			// Mark this vertex as visited, but not fully explored.
			active[v] = 1

			for _, e := range g[v].Edges {
				stack.Push(e.Vertex)
			}
		}
	}

	return tsort[i:]
}

// Transpose writes to h a copy of the current graph with all edges reversed.
func (g Graph) Transpose(h Graph) Graph {
	if cap(h) >= len(g) {
		h = h[:len(g)]
	} else {
		h = make(Graph, len(g))
	}

	for u := range g {
		for _, e := range g[u].Edges {
			h.AddEdge(e.Vertex, u, e.Weight)
		}
	}

	return h
}

func writePath(path []int, v, dist int, prev []int) []int {
	path = resizeIntSlice(path, dist+1)
	path[dist] = v

	for i := dist - 1; i >= 0; i-- {
		v = prev[v]
		path[i] = v
	}

	return path
}

func resizeIntSlice(s []int, size int) []int {
	if cap(s) >= size {
		return s[:size]
	}
	return make([]int, size, 1<<uint(bits.Len(uint(size-1))))
}
