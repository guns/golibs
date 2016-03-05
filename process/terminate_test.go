package process

import (
	"os/exec"
	"reflect"
	"testing"
	"time"
)

func TestAlive(t *testing.T) {
	cmd := exec.Command("sleep", "1")
	err := cmd.Start()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !(alive(cmd.Process)) {
		t.Errorf("expected: alive(cmd.Process)")
	}

	err = cmd.Process.Kill()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	time.Sleep(1 * time.Millisecond)

	if !(alive(cmd.Process)) {
		t.Errorf("expected: alive(cmd.Process), actual: zombie")
	}

	err = cmd.Wait()
	if err == nil {
		t.Errorf("expected err to be an error, but got nil")
	}

	if alive(cmd.Process) {
		t.Errorf("expected: !(alive(cmd.Process))")
	}
}

func TestTerminate(t *testing.T) {
	// SIGTERM
	cmd := exec.Command("sleep", "1")
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
	go Terminate(cmd.Process, 100*time.Millisecond)
	err = <-reaper
	elapsed := time.Since(start)
	if !(elapsed < 10*time.Millisecond) {
		t.Errorf("expected: elapsed < 10*time.Millisecond, actual: %v (the process did not exit immediately after SIGTERM)", elapsed)
	}
	if reflect.TypeOf(err) != reflect.TypeOf(&exec.ExitError{}) {
		t.Errorf("expected type %v, but got type %v", reflect.TypeOf(&exec.ExitError{}), reflect.TypeOf(err))
	}

	// SIGKILL
	cmd = exec.Command("sh", "-c", "trap '' TERM; sleep 1")
	start = time.Now()
	err = cmd.Start()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	reaper2 := make(chan error, 1)
	go func() {
		reaper2 <- cmd.Wait()
		close(reaper2)
	}()
	time.Sleep(2 * time.Millisecond) // Give sh a chance to set the trap
	go Terminate(cmd.Process, 30*time.Millisecond)
	err = <-reaper2
	elapsed = time.Since(start)
	if !(30*time.Millisecond < elapsed && elapsed < 35*time.Millisecond) {
		t.Errorf("expected: 30ms < elapsed < 35ms, actual: %v (the process ignored SIGTERM, but died on SIGKILL)", elapsed)
	}
	if reflect.TypeOf(err) != reflect.TypeOf(&exec.ExitError{}) {
		t.Errorf("expected type %v, but got type %v", reflect.TypeOf(&exec.ExitError{}), reflect.TypeOf(err))
	}

	// NOP
	cmd = exec.Command("true")
	start = time.Now()
	err = cmd.Start()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	reaper3 := make(chan error, 1)
	go func() {
		reaper3 <- cmd.Wait()
		close(reaper3)
	}()
	time.Sleep(2 * time.Millisecond) // Give the process a chance to exit
	go Terminate(cmd.Process, 10*time.Millisecond)
	err = <-reaper3
	elapsed = time.Since(start)
	if !(elapsed < 3*time.Millisecond) {
		t.Errorf("expected: elapsed < 3ms, actual: %v (the process was already dead before termination)", elapsed)
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
