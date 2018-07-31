// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

// A WeightMapper returns the unique weight of an edge (u, v).
type WeightMapper interface {
	Weight(u, v int) float64
}

type pair [2]int

// WeightMap is a simple WeightMapper implementation that returns a default
// weight for edges without explicit weights.
type WeightMap struct {
	m             map[pair]float64
	defaultWeight float64
}

// MakeWeightMap returns a new WeightMap. Note that the defaultWeight of
// a WeightMap is immutable.
func MakeWeightMap(defaultWeight float64, sizeHint int) WeightMap {
	return WeightMap{
		m:             make(map[pair]float64, sizeHint),
		defaultWeight: defaultWeight,
	}
}

// Weight returns the weight of an edge (u, v), or the default edge weight if
// that edge's weight is undefined.
func (m WeightMap) Weight(u, v int) float64 {
	if w, ok := m.m[pair{u, v}]; ok {
		return w
	}
	return m.defaultWeight
}

// SetWeight sets the weight of an edge (u, v) to w. If w is equal to the
// default edge weight, no assignment occurs.
func (m WeightMap) SetWeight(u, v int, w float64) {
	if w != m.defaultWeight {
		m.m[pair{u, v}] = w
	}
}

// DeleteWeight removes the weight associated with edge (u, v), if any.
func (m WeightMap) DeleteWeight(u, v int) {
	delete(m.m, pair{u, v})
}
