// Copyright (c) 2015-2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// +build !windows

// Package process provides tools for working with OS processes.
package process

import (
	"os"
	"syscall"
	"time"
)

// IsAlive implements `kill -0`. Note that calling IsAlive() on a zombie process
// will return true.
func IsAlive(p *os.Process) bool {
	if p == nil {
		return false
	}
	return p.Signal(syscall.Signal(0)) == nil
}

// Terminate sends SIGTERM to a command process, then sends SIGKILL after
// timeout if it is still alive. The process is not reaped, and zombies are
// considered to be alive.
func Terminate(p *os.Process, timeout time.Duration) {
	if !IsAlive(p) {
		return
	}

	// Notify the process politely
	if err := p.Signal(syscall.SIGTERM); err != nil {
		return
	}

	time.Sleep(timeout)

	if !IsAlive(p) {
		return
	}

	_ = p.Kill() // errcheck: If SIGKILL fails, what's really to be done?
}
