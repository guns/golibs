package trigger

import (
	"sync"
	"testing"
)

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
		t Trigger
		u Trigger
	}{Make(), Make()}

	trg.u.Trigger()

	if trg.t.Activated() {
		t.Errorf("expected: !trg.t.Activated()")
	}

	if !trg.u.Activated() {
		t.Errorf("expected: trg.u.Activated()")
	}
}

//
// BenchmarkTrigger-4    3000000  578 ns/op  224 B/op  4 allocs/op
// BenchmarkChan-4       3000000  502 ns/op  192 B/op  2 allocs/op
// BenchmarkWaitGroup-4  3000000  500 ns/op  32  B/op  2 allocs/op
//

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
