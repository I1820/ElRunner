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
}

type runner struct {
	task *Task
	evs  chan Event
	stp  chan int
}

// New creates new runner based on given task
func New(t *Task) Runner {
	return &runner{
		task: t,
		evs:  make(chan Event, 100),
		stp:  make(chan int, 1),
	}
}

func (r *runner) listen() {
	// TODO: Listen for upcomming events
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
	go r.listen()
	go r.interval(r.task.Interval)
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
	r.stp <- 1
}
