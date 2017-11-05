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
	Trigger(e *DataEvent)
}

type runner struct {
	task *Task
	evs  chan Event
	stp  chan int
}

// New creates new runner based on given task
func New(t *Task, backlog int) Runner {
	return &runner{
		task: t,
		evs:  make(chan Event, backlog),
		stp:  make(chan int),
	}
}

func (r *runner) Trigger(e *DataEvent) {
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
			r.task.Run(ev)
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
