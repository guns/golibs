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

	if "false" != <-log {
		t.Errorf("%v != %v", "false", <-log)
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

	if "true" != <-log {
		t.Errorf("%v != %v", "true", <-log)
	}
	select {
	case <-trg.Channel():
	default:
		t.Fail()
	}
}
