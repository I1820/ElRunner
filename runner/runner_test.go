package runner

import (
	"context"
	"encoding/base64"
	"io"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	assert.True(t, r.Status())
	r.Stop()
	assert.False(t, r.Status())
}

func TestDataEvent(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	var mx sync.WaitGroup

	r := New(&Task{
		Run: func(e Event) (string, error) {
			if e.Type() != DataEventType {
				t.Fatal("Invalid Event Type")
			}
			t.Log(string(e.Data()))
			mx.Add(1)
			wg.Done()
			return "", nil
		},
		Interval: 0,
	}, 100)

	go r.Start()
	go r.Start()
	go r.Start()

	assert.NoError(t, r.DataEvent(context.Background(), "Hello"))
	assert.NoError(t, r.DataEvent(context.Background(), "Bye"))

	// all runs are done
	wg.Wait()

	// there are 2 runs at all
	mx.Done()
	mx.Done()
	mx.Wait()

	r.Stop()

	assert.Error(t, r.DataEvent(context.Background(), "Error"))
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
			if _, err := io.WriteString(stdin, base64.StdEncoding.EncodeToString([]byte(e.Data()))); err != nil {
				t.Fatal(err)
			}
			if err := stdin.Close(); err != nil {
				t.Fatal(err)
			}

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

	o, err := r.Output(context.Background())
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(o, "hello from python"), "Invalid message from python: %q do not have prefix \"Hello from python\"", o)

	r.Stop()
}
