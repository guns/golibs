// Copyright (c) 2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// +build !windows

// Package sys provides convenience functions around syscalls.
package sys

import "syscall"

// Ioctl is a simple wrapper around syscall.Syscall(syscall.SYS_IOCTL, …).
// Returns nil if the syscall.Errno equals 0.
func Ioctl(a1, a2, a3 uintptr) (r1, r2 uintptr, err error) {
	r1, r2, errno := syscall.Syscall(syscall.SYS_IOCTL, a1, a2, a3)
	if errno != 0 {
		err = errno
	}
	return
}
