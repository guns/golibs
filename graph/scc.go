// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import "github.com/guns/golibs/generic/impl"

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
	w.reset(len(g), wANeg)

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
