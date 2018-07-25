// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

import "math/bits"

// Packed2DGenericTypeBuilder is a optionally auto-growing [][]GenericType
// builder that uses a single backing slice to reduce allocations. This is
// useful for building read-only non-rectangular 2D data.
type Packed2DGenericTypeBuilder struct {
	head, tail int
	buf        []GenericType
	Rows       [][]GenericType // Contains all finished rows; shares memory with buf
	autoGrow   bool
}

// NewPacked2DGenericTypeBuilder returns a new auto-growing [][]GenericType
// builder that can accommodate at least size items.
func NewPacked2DGenericTypeBuilder(size int) *Packed2DGenericTypeBuilder {
	return NewPacked2DGenericTypeBuilderWithBuffer(
		make([]GenericType, size),
	)
}

// NewPacked2DGenericTypeBuilderFromRows returns a new auto-growing
// [][]GenericType builder from the provided 2D GenericType slice.
//
// This constructor is equivalent to
//
//	NewPacked2DGenericTypeBuilderWithBuffer(rows[0][:cap(rows[0])])
//
// and is intended to help provide a function interface for reusing
// a [][]GenericType without requiring the caller to pass in a
// Packed2DGenericTypeBuilder.
//
func NewPacked2DGenericTypeBuilderFromRows(rows [][]GenericType) *Packed2DGenericTypeBuilder {
	if cap(rows) == 0 {
		return NewPacked2DGenericTypeBuilderWithBuffer(nil)
	}
	buf := rows[:1][0]
	return NewPacked2DGenericTypeBuilderWithBuffer(buf[:cap(buf)])
}

// NewPacked2DGenericTypeBuilderWithBuffer returns a new auto-growing
// [][]GenericType builder that wraps the provided buffer, which is never
// resliced beyond its current length.
func NewPacked2DGenericTypeBuilderWithBuffer(buf []GenericType) *Packed2DGenericTypeBuilder {
	return &Packed2DGenericTypeBuilder{
		buf:      buf,
		autoGrow: true,
	}
}

// SetAutoGrow enables or disables auto-growing.
func (p *Packed2DGenericTypeBuilder) SetAutoGrow(t bool) {
	p.autoGrow = t
}

// Len returns the total number of elements added to finished rows and the
// active partition.
func (p *Packed2DGenericTypeBuilder) Len() int {
	return p.tail
}

// Cap returns the logical capacity of this builder. Note that this may be
// smaller than the capacity of the internal slice.
func (p *Packed2DGenericTypeBuilder) Cap() int {
	return len(p.buf)
}

// Append a single GenericType to the current active partition. If adding this
// element would overflow the internal buffer and auto-growing is enabled, the
// current buffer is moved to a larger buffer and p.Rows is recreated before
// adding the element.
func (p *Packed2DGenericTypeBuilder) Append(x GenericType) {
	if p.autoGrow && p.tail >= len(p.buf) {
		p.Grow(1)
	}
	p.buf[p.tail] = x
	p.tail++
}

// FinishRow appends the current active partition to p.Rows as a []GenericType.
func (p *Packed2DGenericTypeBuilder) FinishRow() {
	p.Rows = append(p.Rows, p.buf[p.head:p.tail])
	p.head = p.tail
}

// Grow internal buffer to accommodate at least n more items.
func (p *Packed2DGenericTypeBuilder) Grow(n int) {
	// We do not check to see if n <= cap(q.a) - len(q.a) because we promised
	// never to reslice the current buffer beyond its current length.
	if n <= 0 {
		return
	}

	buf := make([]GenericType, 1<<uint(bits.Len(uint(len(p.buf)+n-1))))
	copy(buf, p.buf[:p.tail])
	p.buf = buf

	// Recreate rows
	head, tail := 0, 0
	for i := range p.Rows {
		tail += len(p.Rows[i])
		p.Rows[i] = buf[head:tail]
		head = tail
	}
}

// Reset this Packed2DGenericTypeBuilder.
// Note that the internal buffer is NOT cleared.
func (p *Packed2DGenericTypeBuilder) Reset() {
	p.head = 0
	p.tail = 0
	p.Rows = p.Rows[:0]
}
