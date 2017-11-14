/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 17-10-2017
 * |
 * | File Name:     event.go
 * +===============================================
 */

package runner

import "time"

// Types of Events that exist in system
const (
	IntervalEventType int = iota
	DataEventType
)

// Event represents type and data of occurred events
type Event interface {
	Type() int
	Data() []byte
}

// IntervalEvent occurs when user specific interval finishes
type IntervalEvent struct {
	time time.Time
}

// Type returns type of event
func (i *IntervalEvent) Type() int {
	return IntervalEventType
}

// Data returns data associated with event
func (i *IntervalEvent) Data() []byte {
	return []byte(i.time.Format(time.RFC3339))
}

// DataEvent occurs when new data comes from push service
type DataEvent struct {
	data []byte
}

// Type returns type of event
func (d *DataEvent) Type() int {
	return DataEventType
}

// Data returns data associated with event
func (d *DataEvent) Data() []byte {
	return d.data
}
