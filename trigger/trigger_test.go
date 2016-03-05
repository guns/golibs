package trigger

import (
	"fmt"
	"testing"
)

func TestTrigger(t *testing.T) {
	log := make(chan string, 2)
	trg := New()

	go func() {
		log <- fmt.Sprint(trg.Activated())
		trg.Wait()
		log <- fmt.Sprint(trg.Activated())
		close(log)
	}()

	if <-log != "false" {
		t.Errorf("%v != %v", <-log, "false")
	}

	select {
	case <-trg.Channel():
		t.Fail()
	default:
	}

	// Test idempotency of Trigger
	for i := 0; i < 5; i++ {
		go trg.Trigger()
	}

	if <-log != "true" {
		t.Errorf("%v != %v", <-log, "true")
	}

	select {
	case <-trg.Channel():
	default:
		t.Fail()
	}
}
