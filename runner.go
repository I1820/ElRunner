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

// Runner binded functions
type Runner interface {
	Start()
}

type runner struct {
	Task *Task
}

// New creates new runner based on given task
func New(t *Task) Runner {
	return &runner{Task: t}
}

func (r *runner) listen() {
	// TODO: Listen for upcomming events
}

func (r *runner) interval(i int) {
	// TODO: Listen for interval timeouts
}

func (r *runner) Start() {
	evs := make(chan Event, 100)
	go r.listen()
	go r.interval(r.Task.Interval)
	for {
		ev := <-evs
		r.Task.Run()
	}
}
