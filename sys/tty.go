// Copyright (c) 2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// +build !windows

package sys

import (
	"syscall"
	"unsafe"
)

// GetTTYState writes the TTY state of fd to termios.
func GetTTYState(fd uintptr, termios *syscall.Termios) error {
	_, _, err := Ioctl(fd, syscall.TCGETS, uintptr(unsafe.Pointer(termios)))
	return err
}

// SetTTYState sets at least one of the changes requested in termios were
// carried out. It may therefore be necessary to call SetTTYState multiple
// times to make multiple changes to the TTY indicated by fd. See tcsetattr(3)
// for more details.
func SetTTYState(fd uintptr, termios *syscall.Termios) error {
	_, _, err := Ioctl(fd, syscall.TCSETS, uintptr(unsafe.Pointer(termios)))
	return err
}

// DisableTTYCanonicalMode disables the ICANON flag in the TTY indicated by fd
// and returns a function that will return the TTY to its original state. If
// fd is not a TTY or the change was not successful, an error is returned and
// a nil restoreTTY function is returned.
func DisableTTYCanonicalMode(fd uintptr) (restoreTTY func() error, err error) {
	oldstate := syscall.Termios{}

	if err := GetTTYState(fd, &oldstate); err != nil {
		return nil, err
	}

	newstate := oldstate
	newstate.Lflag &^= syscall.ICANON

	if err := SetTTYState(fd, &newstate); err != nil {
		return nil, err
	}

	return func() error { return SetTTYState(fd, &oldstate) }, nil
}
