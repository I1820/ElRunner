package runner

import (
	"encoding/base64"
	"io"
	"os/exec"
	"testing"
	"time"
)

func TestIntervalEvent(t *testing.T) {
	var rz = make(chan int)
	r := New(&Task{
		Run: func(e Event) (string, error) {
			if e.Type() != IntervalEventType {
				t.Fatal("Invalid Event Type")
			}
			ts, _ := time.Parse(time.RFC3339, string(e.Data()))
			te := time.Now()
			t.Log(te.Sub(ts))
			rz <- 1
			return "", nil
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
		Run: func(e Event) (string, error) {
			if e.Type() != DataEventType {
				t.Fatal("Invalid Event Type")
			}
			t.Log(string(e.Data()))
			rz <- 1
			return "", nil
		},
		Interval: 0,
	}, 100)
	go r.Start()
	r.DataEvent([]byte("Hello"))
	r.DataEvent([]byte("Bye"))
	<-rz
	r.Stop()
}

func TestPython(t *testing.T) {
	r := New(&Task{
		Run: func(e Event) (string, error) {
			if e.Type() != IntervalEventType {
				t.Fatal("Invalid Event Type")
			}
			ts, _ := time.Parse(time.RFC3339, string(e.Data()))
			te := time.Now()
			t.Log(te.Sub(ts))

			cmd := exec.Command("python3", "hello.py")

			stdin, err := cmd.StdinPipe()
			if err != nil {
				t.Fatal(err)
			}
			io.WriteString(stdin, base64.StdEncoding.EncodeToString(e.Data()))
			stdin.Close()

			out, err := cmd.Output()
			if err != nil {
				t.Fatal(err)
			}
			t.Log(string(out))

			return string(out), nil
		},
		Interval: 1 * time.Second,
	}, 100)
	go r.Start()

	t.Log(r.Output())

	r.Stop()
}
