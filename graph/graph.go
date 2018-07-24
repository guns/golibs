// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package graph provides a general-purpose directed graph implementation.
package graph

import "math/bits"

// A Graph is a set { (V, E) : V ⊆ ℕ⁰ and E ⊆ { (u, v) ∈ V } }.
//
// Concretely, a Graph is implemented as an adjacency list of (ℕ⁰, []ℕ⁰).
// Note that vertices are numbered from zero.
//
// All edges are directed and unweighted. Undirected graphs can be constructed
// by simply adding the reverse of each edge, and edge weights can be stored
// in a parallel data structure.
//
// Because vertices are represented as signed integers, the maximum size of a
// graph is machineUintLen/2.
//
type Graph [][]int

// undefined is a sentinel value for the set of Vertex indices.
const undefined = -1

// AddEdge adds a single directed edge from vertex u to v.
func (g Graph) AddEdge(u, v int) {
	g[u] = append(g[u], v)
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
func (g Graph) LeastEdgesPath(path []int, u, v int, w *Workspace) []int {
	w.prepare(len(g), wA|wBNeg)

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

	return writePath(path, pred, v, int(dist[v]))
}

// TopologicalSort returns a slice of vertex indices in topologically sorted
// order. The indices are written to the tsort parameter, which is grown if
// necessary. If a topological sort is impossible because there is a cycle in
// the graph, an empty slice (tsort[:0]) is returned.
func (g Graph) TopologicalSort(tsort []int, w *Workspace) []int {
	w.prepare(len(g), 0)

	bs := w.makeBitsliceN(2, wA)
	active := bs[0]          // |V|  · Bitslice of vertex -> active?
	explored := bs[1]        // |V|  · Bitslice of vertex -> fully explored?
	stack := w.makeStack(wB) // |V|w · DFS stack

	tsort = resizeIntSlice(tsort, len(g)) // Prepare write buffer
	idx := len(g)                         // tsort write index + 1

	for u := range g {
		if explored.Get(u) {
			continue
		}

		// DFS
		stack.Push(u)

		// visit(u)
		for stack.Len() > 0 {
			u := stack.Pop()

			// Post-order visit nodes whose children have been explored.
			// These nodes are encoded as their ones' complement.
			if u < 0 {
				u = ^u
				explored.Set(u)
				idx--
				tsort[idx] = u
				continue
			}

			if explored.Get(u) {
				// Ignore fully explored nodes
				continue
			} else if active.Get(u) {
				// This neighboring vertex is active but not yet
				// fully explored, so we have discovered a cycle!
				return tsort[:0]
			}

			// Mark this vertex as visited, but not fully explored.
			active.Set(u)

			// Post-order visit this parent vertex after all its
			// children have been fully explored.
			stack.Push(^u)

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

	for u := range g {
		for _, v := range g[u] {
			h.AddEdge(v, u)
		}
	}

	return h
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

func resizeIntSlice(s []int, size int) []int {
	if cap(s) >= size {
		return s[:size]
	}
	return make([]int, size, 1<<uint(bits.Len(uint(size-1))))
}
