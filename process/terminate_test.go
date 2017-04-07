// Copyright (c) 2015-2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package process

import (
	"os/exec"
	"reflect"
	"testing"
	"time"
)

func TestIsAlive(t *testing.T) {
	cmd := exec.Command("sleep", "1")
	err := cmd.Start()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !IsAlive(cmd.Process) {
		t.Errorf("expected: IsAlive(cmd.Process)")
	}

	err = cmd.Process.Kill()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	time.Sleep(5 * time.Millisecond)

	if !IsAlive(cmd.Process) {
		t.Errorf("expected: IsAlive(cmd.Process)") // zombies are alive
	}

	err = cmd.Wait()
	if err == nil {
		t.Errorf("expected err to be an error, but got nil")
	}

	if IsAlive(cmd.Process) {
		t.Errorf("expected: !IsAlive(cmd.Process)")
	}

	if IsAlive(nil) {
		t.Errorf("expected: !IsAlive(nil)")
	}
}

func TestTerminate(t *testing.T) {
	data := []struct {
		cmd          []string
		lower, upper time.Duration
		err          error
	}{
		{
			[]string{"sleep", "1"},
			0,
			15 * time.Millisecond,
			&exec.ExitError{},
		},
		{
			[]string{"sh", "-c", "trap '' TERM; sleep 1"},
			100 * time.Millisecond,
			900 * time.Millisecond,
			&exec.ExitError{},
		},
		{
			[]string{"true"},
			0,
			10 * time.Millisecond,
			nil,
		},
	}

	for _, row := range data {
		cmd := exec.Command(row.cmd[0], row.cmd[1:]...)
		start := time.Now()
		err := cmd.Start()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		reaper := make(chan error, 1)
		go func() {
			reaper <- cmd.Wait()
			close(reaper)
		}()
		time.Sleep(5 * time.Millisecond) // Process setup time
		go Terminate(cmd.Process, 100*time.Millisecond)
		err = <-reaper
		elapsed := time.Since(start)
		if !(row.lower <= elapsed && elapsed <= row.upper) {
			t.Errorf("expected: %v ≤ elapsed ≤ %v, actual: %v", row.lower, row.upper, elapsed)
		}
		if reflect.TypeOf(err) != reflect.TypeOf(row.err) {
			t.Errorf("expected type %v, but got type %v", reflect.TypeOf(row.err), reflect.TypeOf(err))
		}
	}
}
