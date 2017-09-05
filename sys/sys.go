// Copyright (c) 2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// +build !windows

// Package sys provides convenience functions around syscalls.
package sys

import (
	"os"
	"syscall"
)

// Ioctl is a simple wrapper around syscall.Syscall(syscall.SYS_IOCTL, …).
// Returns nil if the syscall.Errno equals 0.
func Ioctl(a1, a2, a3 uintptr) (r1, r2 uintptr, err error) {
	r1, r2, errno := syscall.Syscall(syscall.SYS_IOCTL, a1, a2, a3)
	if errno != 0 {
		err = errno
	}
	return
}

// Pipe is a wrapper around pipe(2). Returns nil *os.File objects on error.
func Pipe() (r, w *os.File, err error) {
	fildes := make([]int, 2)
	if err := syscall.Pipe(fildes); err != nil {
		return nil, nil, err
	}

	return os.NewFile(uintptr(fildes[0]), "r"), os.NewFile(uintptr(fildes[1]), "w"), nil
}
