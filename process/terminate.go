package process

import (
	"os"
	"syscall"
	"time"
)

// alive implements `kill -0`. Note that calling alive() on a zombie process
// will return true.
func alive(p *os.Process) bool {
	if p == nil {
		return false
	}
	return p.Signal(syscall.Signal(0)) == nil
}

// Terminate sends SIGTERM to a command process, then sends SIGKILL after
// timeout if it is still alive. The process is not reaped, and zombies are
// considered to be alive.
func Terminate(p *os.Process, timeout time.Duration) {
	if !alive(p) {
		return
	}

	// Notify the process politely
	if err := p.Signal(syscall.SIGTERM); err != nil {
		return
	}

	time.Sleep(timeout)

	if !alive(p) {
		return
	}

	// If SIGKILL fails, what's really to be done?
	_ = p.Kill()
}
