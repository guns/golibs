package trigger

import "sync/atomic"

/*

A Trigger is a fast, synchronized way to communicate a single state
transition.

*/
type Trigger struct {
	status uint32
	ch     chan struct{}
}

// New creates a Trigger.
func New() *Trigger {
	return &Trigger{ch: make(chan struct{})}
}

// Activated quickly checks to see if this Trigger has been activated.
func (t *Trigger) Activated() bool {
	return atomic.LoadUint32(&t.status) != 0
}

// Trigger communicates a state transition. This function is idempotent.
func (t *Trigger) Trigger() {
	if atomic.CompareAndSwapUint32(&t.status, 0, 1) {
		close(t.ch)
	}
}

// Channel returns a read channel that can be used to receive a transition
// notification in a select operation.
func (t *Trigger) Channel() <-chan struct{} {
	return t.ch
}

// Wait blocks the current goroutine until this Trigger is activated.
func (t *Trigger) Wait() {
	// Check fast path first
	if !t.Activated() {
		<-t.ch
	}
}
