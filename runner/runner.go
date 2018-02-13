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

import "time"

// Runner runs your task on sepecific events and stores
// outputs. ErrHandler calls on each run that have error
// in another thread so write it in async mode.
type Runner struct {
	task *Task
	evs  chan Event
	out  chan *output
	stp  chan int

	ErrHandler func(error)
}

// New creates new runner based on given task
func New(t *Task, backlog int) *Runner {
	return &Runner{
		task: t,
		evs:  make(chan Event, backlog),
		out:  make(chan *output, backlog),
		stp:  make(chan int),

		ErrHandler: func(err error) {},
	}
}

// Trigger runner and gets its last event
// it blocks until one event come
func (r *Runner) Trigger() Event {
	return <-r.evs
}

// DataEvent push data event (string) into runner events
func (r *Runner) DataEvent(data string) {
	r.evs <- &DataEvent{
		data,
	}
}

// Event push event into runner events
func (r *Runner) Event(e Event) {
	r.evs <- e
}

// Start starts runner and it must be call in new goroutine
func (r *Runner) Start() {
	var t <-chan time.Time
	if r.task.Interval > 0 {
		t = time.Tick(r.task.Interval)
	}
	for {
		select {
		case ev := <-r.evs:
			// channel is close let's go back
			if ev == nil {
				return
			}

			s, err := r.task.Run(ev)
			if s != "" || err != nil {
				r.out <- &output{
					s: s,
					e: err,
				}
				if err != nil {
					go r.ErrHandler(err)
				}
			}
		case <-r.stp:
			break
		case t := <-t:
			r.evs <- &IntervalEvent{
				time: t,
			}
		}
	}
}

// Stop stops runner and you cann't run it again
func (r *Runner) Stop() {
	t := time.NewTimer(time.Second * 10)
	select {
	case r.stp <- 1:
		break
	case <-t.C:
		break
	}
	close(r.evs)
	close(r.out)
	close(r.stp)
}

// Output returns last output from output queue
func (r *Runner) Output() (string, error) {
	o := <-r.out
	return o.s, o.e
}

// OutputBoundedWait returns last output from output queue in given duration
// or returns nil
func (r *Runner) OutputBoundedWait(d time.Duration) (string, error) {
	t := time.Tick(d)
	select {
	case <-t:
		return "", nil
	case o := <-r.out:
		return o.s, o.e
	}
}
