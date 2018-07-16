// Copyright (c) 2015-2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package trigger

import "testing"

func TestTrigger(t *testing.T) {
	exit := New()
	start := New()

	go func() {
		start.Trigger()
		<-exit.Channel()
	}()

	start.Wait()

	if exit.Activated() {
		t.Errorf("expected: !exit.Activated()")
	}
	if !start.Activated() {
		t.Errorf("expected: start.Activated()")
	}

	exit.Trigger()
	exit.Trigger() // assert: should not panic

	if !exit.Activated() {
		t.Errorf("expected: exit.Activated()")
	}
}
