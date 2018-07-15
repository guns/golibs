// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package genericbenchmarks

import "testing"

// goos: linux
// goarch: amd64
// pkg: github.com/guns/golibs/generic/genericbenchmarks
// BenchmarkChannelQueue-4            30000             52545 ns/op               0 B/op          0 allocs/op
// BenchmarkIntQueue-4               200000              7456 ns/op               0 B/op          0 allocs/op
// PASS
// ok      github.com/guns/golibs/generic/genericbenchmarks        3.683s

const queuedepth = 1000

func BenchmarkChannelQueue(b *testing.B) {
	ch := make(chan int, queuedepth)

	for i := 0; i < b.N; i++ {
		for j := 0; j < queuedepth; j++ {
			ch <- j
		}
		for len(ch) > 0 {
			_ = <-ch
		}
	}
}
func BenchmarkIntQueue(b *testing.B) {
	q := NewIntQueue(queuedepth)

	for i := 0; i < b.N; i++ {
		for j := 0; j < queuedepth; j++ {
			q.Enqueue(j)
		}
		for q.Len() > 0 {
			_ = q.Dequeue()
		}
	}
}
