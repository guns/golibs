// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

// TopologicalSort returns a slice of vertex indices in topologically sorted
// order. The indices are written to the tsort parameter, which is grown if
// necessary. If a topological sort is impossible because there is a cycle in
// the graph, an empty slice (tsort[:0]) is returned.
func (g Graph) TopologicalSort(tsort []int, w *Workspace) []int {
	w.reset(len(g), 0)

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
