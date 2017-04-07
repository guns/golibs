// Copyright (c) 2016-2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package trap

import (
	"errors"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/guns/golibs/trigger"
)

var errExpected = errors.New("FEXIT")

func TestHandlerReturnsFunctionError(t *testing.T) {
	err := ExecuteWithHandlers(HandlerMap{}, nil, func(_ *trigger.T) error {
		return errExpected
	})
	if err != errExpected {
		t.Errorf("%v != %v", err, errExpected)
	}

	err = ExecuteWithHandlers(HandlerMap{}, nil, func(_ *trigger.T) error {
		return nil
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHandlerActions(t *testing.T) {
	var n int32
	data := []struct {
		hmap     HandlerMap
		err      error
		expected int32
	}{
		{HandlerMap{syscall.SIGUSR1: {None, func(_ os.Signal, _ *trigger.T) {
			atomic.AddInt32(&n, 10)
		}}}, errExpected, sigChanLen*10 + 1},
		{HandlerMap{syscall.SIGUSR1: {Restart, func(_ os.Signal, _ *trigger.T) {
			atomic.AddInt32(&n, 10)
		}}}, nil, sigChanLen*10 + sigChanLen + 1},
		{HandlerMap{syscall.SIGUSR1: {Restart, func(_ os.Signal, _ *trigger.T) {
			atomic.AddInt32(&n, 10)
		}}}, errExpected, 11},
		{HandlerMap{syscall.SIGUSR1: {Exit, func(_ os.Signal, _ *trigger.T) {
			atomic.AddInt32(&n, 10)
		}}}, nil, 11},
		{HandlerMap{syscall.SIGUSR1: {Exit, func(_ os.Signal, _ *trigger.T) {
			atomic.AddInt32(&n, 10)
		}}}, errExpected, 11},
	}

	for _, row := range data {
		n = int32(0)
		exit := trigger.New()
		start := trigger.New()
		reply := trigger.New()
		var err error

		go func() {
			err = ExecuteWithHandlers(row.hmap, exit, func(fexit *trigger.T) error {
				atomic.AddInt32(&n, 1)
				start.Trigger()
				fexit.Wait()
				return row.err
			})
			reply.Trigger()
		}()

		start.Wait()

		p, _ := os.FindProcess(os.Getpid())
		for i := 0; i < sigChanLen; i++ {
			e := p.Signal(syscall.SIGUSR1)
			if e != nil {
				t.Errorf("unexpected error: %v", e)
			}
			time.Sleep(2 * time.Millisecond) // HACK: May not work on a slow CPU
		}

		exit.Trigger()
		reply.Wait()

		if n != row.expected {
			t.Errorf("%v != %v", n, row.expected)
		}
		if err != row.err {
			t.Errorf("%v != %v", err, row.err)
		}
	}
}

func TestHandlerExitsIfPassedActivatedTrigger(t *testing.T) {
	exit := trigger.New()
	n := int32(0)
	exit.Trigger()
	err := ExecuteWithHandlers(nil, exit, func(fexit *trigger.T) error {
		atomic.AddInt32(&n, 1)
		return errExpected
	})
	if n != 0 {
		t.Errorf("%v != %v", n, 0)
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
