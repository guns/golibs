// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package debug provides debugging facilities.
package debug

import (
	"os"
	"runtime"
	"runtime/pprof"
)

// StartCPUProfile starts CPU profiling and returns a stop function that
// writes the profile to the given path.
func StartCPUProfile(path string) (stopCPUProfile func()) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}

	return func() {
		pprof.StopCPUProfile()
		if err := f.Close(); err != nil {
			panic(err)
		}
	}
}

// StartMemProfile starts memory profiling with 100% sampling and returns
// a stop function that writes the profile to the given path.
func StartMemProfile(path string) (stopMemProfile func()) {
	runtime.MemProfileRate = 1

	return func() {
		runtime.MemProfileRate = 0

		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}

		runtime.GC() // get up-to-date statistics

		if err := pprof.WriteHeapProfile(f); err != nil {
			panic(err)
		}
		if err := f.Close(); err != nil {
			panic(err)
		}
	}
}
