// Package trap handles POSIX signals.
package trap

import (
	"os"
	"os/signal"

	"github.com/guns/golibs/trigger"
)

// An Action specifies the behavior of a Handler
type Action int

const (
	// None ignores signals, but executes the handler fn.
	None Action = iota
	// Restart specifies that a function should be terminated before the
	// handler is called, and restarted afterwards
	Restart
	// Exit specifies that a function should be terminated permanently
	// before the handler is called
	Exit
)

// A Handler specifies a signal handler. The exit parameter of Fn should be
// respected as an "abort ASAP" message.
type Handler struct {
	Action Action
	Fn     func(sig os.Signal, exit *trigger.Trigger)
}

// HandlerMap associates Handler objects to specific signals
type HandlerMap map[os.Signal]Handler

const sigChanLen = 8

// ExecuteWithHandlers executes f after installing the signal handlers
// specified by hmap. The Action of each Handler decides how f and the handler
// function are called:
//
//	None    call handler
//	Restart terminate f; call handler; re-call f
//	Exit    terminate f; call handler
//
// The handler function is ignored if nil, and f must exit ASAP once its
// Trigger is activated.
//
// The return value of f is returned
func ExecuteWithHandlers(hmap HandlerMap, exit *trigger.Trigger, f func(*trigger.Trigger) error) (err error) {
	if exit.Activated() {
		return nil
	}

	sigch := make(chan os.Signal, sigChanLen)
	defer close(sigch)

	// We don't want to call signal.Notify without specific signals
	if len(hmap) > 0 {
		signals := []os.Signal{}
		for k := range hmap {
			signals = append(signals, k)
		}
		signal.Notify(sigch, signals...)
		defer signal.Stop(sigch)
	}

	floop := true

	for floop && !exit.Activated() {
		fexit := trigger.New()
		fexitReply := trigger.New()
		go func() {
			err = f(fexit)
			fexitReply.Trigger()
		}()

		hloop := true

		for hloop {
			select {
			// Function exited normally
			case <-fexitReply.Channel():
				hloop = false
				floop = false
			// Emergency exit: terminate function and exit
			case <-exit.Channel():
				fexit.Trigger()
				fexitReply.Wait()
				hloop = false
				floop = false
			// Handle signal
			case sig := <-sigch:
				h := hmap[sig]
				switch h.Action {
				case None:
					// Do not disturb f, just call h.Fn
					if h.Fn != nil && !exit.Activated() {
						h.Fn(sig, exit)
					}
				case Restart:
					// Terminate f, call h.Fn, restart floop unless f returned an error
					fexit.Trigger()
					fexitReply.Wait()
					if h.Fn != nil && !exit.Activated() {
						h.Fn(sig, exit)
					}
					hloop = false
					if err != nil {
						floop = false
					}
				case Exit:
					// Terminate f, call h.Fn, exit
					fexit.Trigger()
					fexitReply.Wait()
					if h.Fn != nil && !exit.Activated() {
						h.Fn(sig, exit)
					}
					hloop = false
					floop = false
				}
			}
		}
	}

	return err
}
