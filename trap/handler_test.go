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

func TestHandlerReturnsFunctionError(t *testing.T) {
	expected := errors.New("ERROR")
	err := ExecuteWithHandlers(HandlerMap{}, trigger.New(), func(_ *trigger.Trigger) error {
		return expected
	})
	if err != expected {
		t.Errorf("%v != %v", err, expected)
	}

	err = ExecuteWithHandlers(HandlerMap{}, trigger.New(), func(_ *trigger.Trigger) error {
		return nil
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHandlerExit(t *testing.T) {
	exit := trigger.New()
	start := trigger.New()
	reply := trigger.New()
	var err error
	go func() {
		err = ExecuteWithHandlers(HandlerMap{}, exit, func(fexit *trigger.Trigger) (ferr error) {
			start.Trigger()
			<-fexit.Channel()
			return errors.New("FEXIT")
		})
		reply.Trigger()
	}()
	start.Wait()
	exit.Trigger()
	reply.Wait()
	if !(err != nil && err.Error() == "FEXIT") {
		t.Errorf("expected: err.Error() == `FEXIT`, actual: %v", err)
	}
}

func TestHandlerActionNone(t *testing.T) {
	n := int32(0)
	exit := trigger.New()
	start := trigger.New()
	reply := trigger.New()
	var err error

	hmap := HandlerMap{syscall.SIGUSR1: {None, func(_ os.Signal, _ *trigger.Trigger) {
		atomic.AddInt32(&n, 1)
	}}}

	go func() {
		err = ExecuteWithHandlers(hmap, exit, func(fexit *trigger.Trigger) error {
			atomic.AddInt32(&n, sigChanLen)
			start.Trigger()
			<-fexit.Channel()
			return errors.New("fexit")
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
		time.Sleep(time.Millisecond)
	}
	exit.Trigger()
	reply.Wait()
	if n != 2*sigChanLen {
		t.Errorf("%v != %v", n, 6)
	}
	if err.Error() != "fexit" {
		t.Errorf("%v != %v", err.Error(), "fexit")
	}
}

func TestHandlerActionRestartExit(t *testing.T) {
	n := int32(0)
	reply := trigger.New()
	var err error

	hmap := HandlerMap{
		syscall.SIGUSR1: {Restart, func(_ os.Signal, _ *trigger.Trigger) {
			atomic.AddInt32(&n, 1)
		}},
		syscall.SIGTERM: {Exit, nil},
	}

	go func() {
		err = ExecuteWithHandlers(hmap, trigger.New(), func(fexit *trigger.Trigger) error {
			atomic.AddInt32(&n, 1)
			fexit.Wait()
			return nil
		})
		reply.Trigger()
	}()

	p, _ := os.FindProcess(os.Getpid())
	for i := 0; i < sigChanLen; i++ {
		e := p.Signal(syscall.SIGUSR1)
		if e != nil {
			t.Errorf("unexpected error: %v", e)
		}
		time.Sleep(time.Millisecond)
	}
	e := p.Signal(syscall.SIGTERM)
	if e != nil {
		t.Errorf("unexpected error: %v", e)
	}
	reply.Wait()

	if n != 2*sigChanLen-1 {
		t.Errorf("%v != %v", n, 2*sigChanLen-1)
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
