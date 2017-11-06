/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 16-10-2017
 * |
 * | File Name:     task.go
 * +===============================================
 */

package runner

import "time"

// Task represents single task must run on the runner
type Task struct {
	Run      func(ev Event)
	Interval time.Duration
}