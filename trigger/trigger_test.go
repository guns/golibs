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

func TestMake(t *testing.T) {
	trg := struct {
		t T
		u T
	}{Make(), Make()}

	trg.u.Trigger()

	if trg.t.Activated() {
		t.Errorf("expected: !trg.t.Activated()")
	}

	if !trg.u.Activated() {
		t.Errorf("expected: trg.u.Activated()")
	}
}
