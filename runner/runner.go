/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 16-10-2017
 * |
 * | File Name:     runner.go
 * +===============================================
 */

package runner

import (
	"context"
	"fmt"
	"time"
)

// Runner runs your task on sepecific events and stores
// outputs. ErrHandler calls on each run that have error
// in another thread so write it in async mode.
type Runner struct {
	task         *Task
	eventStream  chan Event
	outputStream chan *output
	done         chan struct{}

	ErrHandler func(error)
}

// New creates new runner based on given task
func New(t *Task, backlog int) *Runner {
	return &Runner{
		task:         t,
		eventStream:  make(chan Event, backlog),
		outputStream: make(chan *output, backlog),
		done:         make(chan struct{}),

		ErrHandler: func(err error) {},
	}
}

// NewWithoutOutput creates new runner without any output channel
func NewWithoutOutput(t *Task, backlog int) *Runner {
	return &Runner{
		task:         t,
		eventStream:  make(chan Event, backlog),
		outputStream: nil,
		done:         make(chan struct{}),

		ErrHandler: func(err error) {},
	}
}

// Trigger runner and gets its last event
// it blocks until one event come
func (r *Runner) Trigger() Event {
	return <-r.eventStream
}

// DataEvent push data event (string + environments) into runner events
func (r *Runner) DataEvent(ctx context.Context, data string, envs ...map[string]string) error {
	e := make(map[string]string)

	for _, env := range envs {
		for k, v := range env {
			e[k] = v
		}
	}

	return r.Event(ctx, &DataEvent{
		data,
		e,
	})
}

// Event push event into runner events
func (r *Runner) Event(ctx context.Context, e Event) error {
	if !r.Status() {
		return fmt.Errorf("Cannot push event on stopped runner")
	}

	select {
	case r.eventStream <- e:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Status returns true when runner is in the running state
// otherwise false
func (r *Runner) Status() bool {
	select {
	case <-r.done:
		return false
	default:
		return true
	}
}

// Start starts runner and it must be call in new goroutine
// you can start many routine by call this function many times
func (r *Runner) Start() {
	var t <-chan time.Time
	if r.task.Interval > 0 {
		t = time.Tick(r.task.Interval)
	}
	for {
		select {
		case ev, open := <-r.eventStream:
			if !open {
				return
			}

			wait := make(chan struct{})

			go func() {
				o, err := r.task.Run(ev)

				// write into output channel
				select {
				case r.outputStream <- &output{
					output: o,
					err:    err,
				}:
				default:
				}

				// call user handler
				if err != nil {
					go r.ErrHandler(err)
				}

				close(wait)
			}()

			select {
			case <-wait: // run function response is ready
			case <-r.done: // runner is done
				return
			}
		case t := <-t:
			select {
			case <-r.done:
				return
			default:
			}

			select {
			case r.eventStream <- &IntervalEvent{
				time: t,
			}:
			default:
				continue
			}
		}
	}
}

// Stop stops runner and you cann't run it again
func (r *Runner) Stop() {
	close(r.eventStream)
	close(r.outputStream)
	close(r.done)
}

// Output returns last output from output queue
func (r *Runner) Output(ctx context.Context) (string, error) {
	select {
	case o := <-r.outputStream:
		return o.output, o.err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
