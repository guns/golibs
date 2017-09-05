// Copyright (c) 2015-2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package trigger provides a simple, but flexible way to communicate a single
// state transition.
package trigger

import "sync/atomic"

/*

A trigger is a flexible, synchronized way to communicate a single state
transition. Like sync.Mutex, triggers should never be copied.

*/
type T struct {
	done uint32
	ch   chan struct{}
}

// New creates a trigger.
func New() *T {
	return &T{ch: make(chan struct{})}
}

// Make returns a trigger value. This is only useful when embedding a trigger
// in a struct.
func Make() T {
	return T{ch: make(chan struct{})}
}

// Activated quickly checks to see if this trigger has been activated.
func (t *T) Activated() bool {
	return atomic.LoadUint32(&t.done) != 0
}

// Trigger communicates a state transition. This method is idempotent.
func (t *T) Trigger() {
	if atomic.CompareAndSwapUint32(&t.done, 0, 1) {
		close(t.ch)
	}
}

// Channel returns a read channel that can be used to receive a transition
// notification in a select operation.
func (t *T) Channel() <-chan struct{} {
	return t.ch
}

// Wait blocks the current goroutine until this trigger is activated.
func (t *T) Wait() {
	// Check fast path first
	if atomic.LoadUint32(&t.done) == 0 {
		<-t.ch
	}
}
