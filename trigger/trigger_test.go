package trigger

import (
	"fmt"
	"sync"
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
		trg.Trigger()
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

func TestConstruct(t *testing.T) {
	trg := struct {
		t Trigger
		u Trigger
	}{Construct(), Construct()}

	trg.u.Trigger()

	if trg.t.Activated() {
		t.Errorf("expected: !trg.t.Activated()")
	}

	if !trg.u.Activated() {
		t.Errorf("expected: trg.u.Activated()")
	}
}

func BenchmarkTrigger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		trg := New()
		reply := New()
		go func() {
			trg.Wait()
			reply.Trigger()
		}()
		trg.Trigger()
		reply.Wait()
	}
}

func BenchmarkChan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		trg := make(chan struct{})
		reply := make(chan struct{})
		go func() {
			<-trg
			close(reply)
		}()
		close(trg)
		<-reply
	}
}

func BenchmarkWaitGroup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		trg := sync.WaitGroup{}
		trg.Add(1)
		reply := sync.WaitGroup{}
		reply.Add(1)
		go func() {
			trg.Wait()
			reply.Done()
		}()
		trg.Done()
		reply.Wait()
	}
}
