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

// Grow returns a graph with n more vertices.
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

func resizeIntSlice(s []int, size int) []int {
	if cap(s) >= size {
		return s[:size]
	}
	return make([]int, size, 1<<uint(bits.Len(uint(size-1))))
}
