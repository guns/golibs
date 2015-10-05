package process

import (
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAlive(t *testing.T) {
	cmd := exec.Command("sleep", "1")
	err := cmd.Start()
	assert.Nil(t, err)
	assert.True(t, alive(cmd.Process))
	err = cmd.Process.Kill()
	assert.Nil(t, err)
	time.Sleep(1 * time.Millisecond)
	assert.True(t, alive(cmd.Process), "Zombies are still alive()")
	err = cmd.Wait()
	assert.Error(t, err)
	assert.False(t, alive(cmd.Process))
}

func TestTerminate(t *testing.T) {
	// SIGTERM
	cmd := exec.Command("sleep", "1")
	start := time.Now()
	err := cmd.Start()
	assert.Nil(t, err)
	reaper := make(chan error, 1)
	go func() {
		reaper <- cmd.Wait()
		close(reaper)
	}()
	go Terminate(cmd.Process, 100*time.Millisecond)
	err = <-reaper
	elapsed := time.Since(start)
	assert.True(t, elapsed < 10*time.Millisecond, "The process exited immediately after SIGTERM")
	assert.IsType(t, &exec.ExitError{}, err)

	// SIGKILL
	cmd = exec.Command("sh", "-c", "trap '' TERM; sleep 1")
	start = time.Now()
	err = cmd.Start()
	assert.Nil(t, err)
	reaper2 := make(chan error, 1)
	go func() {
		reaper2 <- cmd.Wait()
		close(reaper2)
	}()
	time.Sleep(2 * time.Millisecond) // Give sh a chance to set the trap
	go Terminate(cmd.Process, 30*time.Millisecond)
	err = <-reaper2
	elapsed = time.Since(start)
	assert.True(t, 30*time.Millisecond < elapsed && elapsed < 35*time.Millisecond,
		"The process ignored SIGTERM, but died on SIGKILL")
	assert.IsType(t, &exec.ExitError{}, err)

	// NOP
	cmd = exec.Command("true")
	start = time.Now()
	err = cmd.Start()
	assert.Nil(t, err)
	reaper3 := make(chan error, 1)
	go func() {
		reaper3 <- cmd.Wait()
		close(reaper3)
	}()
	time.Sleep(2 * time.Millisecond) // Give the process a chance to exit
	go Terminate(cmd.Process, 10*time.Millisecond)
	err = <-reaper3
	elapsed = time.Since(start)
	assert.True(t, elapsed < 3*time.Millisecond, "The process was already dead before termination")
	assert.Nil(t, err)
}
