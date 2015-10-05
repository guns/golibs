package trigger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, "false", <-log)
	select {
	case <-trg.Channel():
		t.Fail()
	default:
	}

	trg.Trigger()
	assert.Equal(t, "true", <-log)
	select {
	case <-trg.Channel():
	default:
		t.Fail()
	}
}
