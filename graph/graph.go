// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package graph provides a general-purpose directed graph implementation.
package graph

import (
	"math/bits"

	"github.com/guns/golibs/generic/impl"
)

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

// Grow returns a graph with n more vertices. The
func (g Graph) Grow(n int) Graph {
	if n <= 0 {
		return g
	} else if len(g)+n < cap(g) {
		i := len(g)
		g = g[:i+n]
		g[i:].Reset()
		return g
	}

	h := make(Graph, 1<<uint(bits.Len(uint(len(g)+n-1))))
	copy(h, g)
	return h[:len(g)+n]
}

// Resize this graph, then reset it. Use Grow to enlarge the graph without
// resetting its state.
func (g Graph) Resize(n int) Graph {
	if n > cap(g) {
		return make(Graph, 1<<uint(bits.Len(uint(n-1))))[:n]
	}

	g = g[:n]
	g.Reset()
	return g
}

// Reset all edge slices in the graph.
// Note that the slices are truncated, NOT cleared.
func (g Graph) Reset() {
	for i := range g {
		g[i] = g[i][:0]
	}
}

// LeastEdgesPath returns a path from vertex u to v with a minimum number of
// edges. The length of the path is the length of the returned path minus one.
//
// The path is written to the path slice, which is grown if necessary.
//
// If no path exists, an empty slice (path[:0) is returned.
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
	active := bs[0]                            // |V|   · Bitslice of vertex -> active?
	explored := bs[1]                          // |V|   · Bitslice of vertex -> fully explored?
	stack := w.makeAutoPromotingStack(wB | wC) // 2|V|w · DFS stack

	tsort = resizeIntSlice(tsort, len(g)) // Prepare write buffer
	idx := len(g)                         // tsort write index + 1

	for u := range g {
		if explored.Get(u) {
			continue
		}

		// DFS
		stack.PushOrPromote(u)

		for stack.Len() > 0 {
			u := stack.Peek()

			// Finish nodes whose children have been explored.
			if active.Get(u) {
				stack.Pop()
				explored.Set(u)
				idx--
				tsort[idx] = u
				continue
			}

			// Mark this vertex as visited, but not fully explored.
			active.Set(u)

			// Visit children nodes
			for _, v := range g[u] {
				if explored.Get(v) {
					// Ignore fully explored nodes
					continue
				} else if active.Get(v) {
					// This neighboring vertex is active but not yet
					// fully explored, so we have discovered a cycle!
					return tsort[:0]
				}
				stack.PushOrPromote(v)
			}
		}
	}

	return tsort[idx:]
}

// Transpose writes to h a copy of the current graph with all edges reversed.
func (g Graph) Transpose(h Graph) Graph {
	h = h.Resize(len(g))

	for u := range g {
		for _, v := range g[u] {
			h.AddEdge(v, u)
		}
	}

	return h
}

// StronglyConnectedComponents returns a forest of strongly connected
// components in reverse topological order. This forest of vertex indices is
// written to scc, which is grown if necessary.
//
// Note that the returned [][]int is actually backed by a single []int
// (scc[0][:cap(scc[0])]) to minimize allocations, and is therefore NOT fit
// for modification.
//
// Passing the same [][]int returned by this function as the scc parameter
// reuses memory and can eliminate unnecessary allocations.
func (g Graph) StronglyConnectedComponents(scc [][]int, w *Workspace) [][]int {
	w.prepare(len(g), wANeg)

	// Tim Leslie's iterative implementation [1] of David Pearce's
	// memory-efficient strongly connected components algorithm. [2]
	//
	// [1]: http://www.timl.id.au/SCC
	// [2]: http://homepages.ecs.vuw.ac.nz/~djp/files/IPL15-preprint.pdf

	rindex := w.a                             // |V|w  · Array of v -> local root index
	dfs := w.makeAutoPromotingStack(wB | wC)  // 2|V|w · Auto-promoting DFS stack
	backtrack := *newNonPromotingStack(dfs.s) // 0     · Backtrack stack; shares memory with dfs

	builder := impl.NewPacked2DIntBuilderFromRows(scc)
	builder.Grow(len(g) - builder.Cap())
	builder.SetAutoGrow(false)

	i := 1
	component := len(g) - 1

	for u := range g {
		if rindex[u] != undefined {
			continue
		}

		// DFS
		dfs.PushOrPromote(u)

		for dfs.Len() > 0 {
			u := dfs.Peek()

			if rindex[u] == undefined {
				// Top of dfs is unvisited, so assign it an index and push or promote its
				// unvisited successors.

				rindex[u] = i
				i++

				for _, v := range g[u] {
					if rindex[v] == undefined {
						dfs.PushOrPromote(v)
					}
				}
			} else {
				// Top of dfs has been visited, so compare it against successors to find a
				// minimum local root index.

				dfs.Pop()
				root := true

				for _, v := range g[u] {
					if rindex[v] < rindex[u] {
						rindex[u] = rindex[v]
						root = false
					}
				}

				if root {
					// u is the local component root, so everything on the backtrack stack
					// that has an rindex >= u's rindex is part of this component.
					for backtrack.len > 0 && rindex[u] <= rindex[backtrack.Peek()] {
						v := backtrack.Pop()
						rindex[v] = component
						builder.Append(v)
						i--
					}

					rindex[u] = component
					builder.Append(u)
					i--

					builder.FinishRow()
					component--
				} else {
					// u is not a local root, so push it on the backtrack stack.
					backtrack.Push(u)
				}
			}
		}
	}

	return builder.Rows
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
