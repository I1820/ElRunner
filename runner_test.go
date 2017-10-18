package runner

import (
	"testing"
	"time"
)

func TestIntervalEvent(t *testing.T) {
	var rz = make(chan int)
	r := New(&Task{
		Run: func(e Event) {
			if e.Type() != IntervalEventType {
				t.Fatal("Invalid Event Type")
			}
			ts, _ := time.Parse(time.RFC3339, e.Data())
			te := time.Now()
			t.Log(te.Sub(ts))
			rz <- 1
		},
		Interval: 1 * time.Second,
	}, 100)
	go r.Start()
	<-rz
	r.Stop()
}

func TestDataEvent(t *testing.T) {
	var rz = make(chan int)
	r := New(&Task{
		Run: func(e Event) {
			if e.Type() != DataEventType {
				t.Fatal("Invalid Event Type")
			}
			t.Log(e.Data())
			rz <- 1
		},
		Interval: 0,
	}, 100)
	go r.Start()
	r.Trigger(&DataEvent{
		data: "Hello",
	})
	r.Trigger(&DataEvent{
		data: "Bye",
	})
	<-rz
	r.Stop()
}
