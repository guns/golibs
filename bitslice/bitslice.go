// Package bitslice provides a bitset implementation.
package bitslice

import "math/bits"

const shift = 5 + ((^uint(0) >> 32) & 1)
const mask = 1<<shift - 1

// SizeOf returns the minimum number of uints required to represent nbits.
func SizeOf(nbits int) int {
	return (nbits + mask) >> shift // (nbits + 63)/64
}

// T is a slice of uint.
type T []uint

// Make creates a new bitslice that accommodates at least nitems.
func Make(nitems int) T {
	return make(T, SizeOf(nitems)) // (nitems+63)/64
}

// Get returns true if the bit at offset i is set, and false otherwise.
func (bs T) Get(i int) bool {
	return bs[uint(i)>>shift]&(1<<(uint(i)&mask)) != 0 // bs[i/64] AND (1 << (i mod 64))
}

// Set a bit.
func (bs T) Set(i int) {
	bs[i>>shift] |= 1 << (uint(i) & mask) // bs[i/64] OR= (1 << (i mod 64))
}

// Clear a bit.
func (bs T) Clear(i int) {
	bs[i>>shift] &= ^(1 << (uint(i) & mask)) // bs[i/64] AND= NOT(1 << (i mod 64))
}

// Toggle a bit.
func (bs T) Toggle(i int) {
	bs[i>>shift] ^= 1 << (uint(i) & mask) // bs[i/64] XOR= (1 << (i mod 64))
}

// CompareAndSet sets a bit and returns true if the bit is clear, and returns
// false otherwise.
func (bs T) CompareAndSet(i int) bool {
	b := i >> shift
	bit := uint(1 << (uint(i) & mask))

	if bs[b]&bit == 0 {
		bs[b] |= bit
		return true
	}

	return false
}

// CompareAndClear clears a bit and returns true if the bit is set, and
// returns false otherwise.
func (bs T) CompareAndClear(i int) bool {
	b := i >> shift
	bit := uint(1 << (uint(i) & mask))

	if bs[b]&bit != 0 {
		bs[b] &= ^bit
		return true
	}

	return false
}

// CompareAndToggle toggles a bit and returns true if the bit state is equal
// to state, and returns false otherwise.
func (bs T) CompareAndToggle(i int, state bool) bool {
	if state {
		return bs.CompareAndClear(i)
	}

	return bs.CompareAndSet(i)
}

// AppendOffsets appends and returns a slice of indices of bits that are set.
func (bs T) AppendOffsets(v []int) []int {
	for i, n := range bs {
		for n != 0 {
			o := bits.TrailingZeros(n)
			v = append(v, (i<<shift)+o)
			n ^= 1 << uint(o) // Toggle bit
		}
	}

	return v
}

// Popcnt returns the number of bits set in this bitslice.
func (bs T) Popcnt() int {
	pop := 0

	for _, n := range bs {
		pop += bits.OnesCount(n)
	}

	return pop
}

// Reset this bitslice by clearing all bits.
func (bs T) Reset() {
	for i := range bs {
		bs[i] = 0
	}
}
