// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package graph provides a general-purpose directed graph implementation.
package graph

import "github.com/guns/golibs/calculate"

// A Graph is a set { (V, E) : V ⊆ ℕ⁰ and E ⊆ { (u, v) ∈ V } }.
//
// Concretely, a Graph is implemented as an adjacency list of (ℕ⁰, []ℕ⁰).
// Note that vertices are numbered from zero.
//
// All edges are directed and unweighted. Undirected graphs can be constructed
// by simply adding the reverse of each edge, and edge weights can be stored
// in a WeightMapper.
//
// Because vertices are represented as signed integers, the maximum size of a
// graph is machineUintLen/2.
//
type Graph [][]int

// AddEdge adds a single directed edge from vertex u to v.
func (g Graph) AddEdge(u, v int) {
	g[u] = append(g[u], v)
}

// RemoveEdge removes a single directed edge from vertex u to v. Note that
// this requires a linear search of u's edge list.
func (g Graph) RemoveEdge(u, v int) {
	edges := g[u]
	for i, succ := range edges {
		if succ == v {
			edges[i] = edges[len(edges)-1] // Overwrite with last vertex
			g[u] = edges[:len(edges)-1]    // Exclude last vertex
		}
	}
}

// Grow returns a graph with n more vertices.
func (g Graph) Grow(n int) Graph {
	switch {
	case n <= 0:
		return g
	case len(g)+n < cap(g):
		i := len(g)
		g = g[:i+n]
		g[i:].ResetEdges()
		return g
	default:
		h := make(Graph, calculate.NextCap(len(g)+n))
		copy(h, g)
		return h[:len(g)+n]
	}
}

// Reset this graph by resizing it and resetting its edges.
// Use Grow to enlarge the graph without truncating its state.
func (g Graph) Reset(n int) Graph {
	if n > cap(g) {
		return make(Graph, calculate.NextCap(n))[:n]
	}

	g = g[:n]
	g.ResetEdges()
	return g
}

// ResetEdges reset all edge slices in the graph.
// Note that edge slices are NOT cleared.
func (g Graph) ResetEdges() {
	for u := range g {
		g[u] = g[u][:0]
	}
}

// Transpose writes to h a copy of the current graph with all edges reversed.
func (g Graph) Transpose(h Graph) Graph {
	h = h.Reset(len(g))

	for u, vs := range g {
		for _, v := range vs {
			h.AddEdge(v, u)
		}
	}

	return h
}
