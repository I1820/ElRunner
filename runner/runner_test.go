package runner

import (
	"io"
	"os/exec"
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

func TestPython(t *testing.T) {
	var rz = make(chan int)
	r := New(&Task{
		Run: func(e Event) {
			if e.Type() != IntervalEventType {
				t.Fatal("Invalid Event Type")
			}
			ts, _ := time.Parse(time.RFC3339, e.Data())
			te := time.Now()
			t.Log(te.Sub(ts))

			cmd := exec.Command("python3", "hello.py")

			stdin, err := cmd.StdinPipe()
			if err != nil {
				t.Fatal(err)
			}
			io.WriteString(stdin, e.Data())
			stdin.Close()

			out, err := cmd.Output()
			if err != nil {
				t.Fatal(err)
			}
			t.Log(string(out))

			rz <- 1
		},
		Interval: 1 * time.Second,
	}, 100)
	go r.Start()
	<-rz
	r.Stop()

}
