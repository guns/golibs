// Copyright (c) 2015-2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package trigger

import (
	"sync"
	"testing"
)

//	go version go1.10.3 linux/amd64
//	goos: linux
//	goarch: amd64
//	pkg: github.com/guns/golibs/trigger
//	BenchmarkTrigger         3000000               452 ns/op             224 B/op          4 allocs/op
//	BenchmarkChan            5000000               395 ns/op             192 B/op          2 allocs/op
//	BenchmarkWaitGroup       3000000               407 ns/op              32 B/op          2 allocs/op
//	PASS
//	ok      github.com/guns/golibs/trigger  5.830s

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
