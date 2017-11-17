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

// Runner binded functions
type Runner interface {
	Start()
	Stop()
	Output() Output
	DataEvent([]byte)
	Event(Event)
}

type runner struct {
	task *Task
	evs  chan Event
	out  chan Output
	stp  chan int
}

// New creates new runner based on given task
func New(t *Task, backlog int) Runner {
	return &runner{
		task: t,
		evs:  make(chan Event, backlog),
		out:  make(chan Output, backlog),
		stp:  make(chan int),
	}
}

func (r *runner) DataEvent(data []byte) {
	r.evs <- &DataEvent{
		data,
	}
}

func (r *runner) Event(e Event) {
	r.evs <- e
}

func (r *runner) interval(i time.Duration) {
	for {
		time.Sleep(i)
		r.evs <- &IntervalEvent{
			time: time.Now(),
		}
	}

}

func (r *runner) Start() {
	if r.task.Interval > 0 {
		go r.interval(r.task.Interval)
	}
	for {
		select {
		case ev := <-r.evs:
			o := r.task.Run(ev)
			if o != "" {
				r.out <- o
			}
		case <-r.stp:
			break
		}
	}
}

func (r *runner) Stop() {
	t := time.NewTimer(time.Second * 10)
	select {
	case r.stp <- 1:
		break
	case <-t.C:
		break
	}
}

func (r *runner) Output() Output {
	o := <-r.out
	return o
}
