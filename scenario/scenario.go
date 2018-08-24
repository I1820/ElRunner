/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 13-12-2017
 * |
 * | File Name:     scenario.go
 * +===============================================
 */

package scenario

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"time"

	"github.com/I1820/ElRunner/runner"
	"github.com/powerman/rpc-codec/jsonrpc2"
	log "github.com/sirupsen/logrus"
)

// Endpoint for scenario communication
type Endpoint struct {
	s *Scenario
}

// WaitForData waits for things incomming data
func (e *Endpoint) WaitForData(args int, reply *string) error {
	// manually read data from channel beacuse runner is in blocking section
	d := e.s.runner.Trigger()
	*reply = d.Data()
	return nil
}

// About let you about who left us alone
func (e Endpoint) About(args int, reply *string) error {
	*reply = "18.20 is leaving us"
	return nil
}

// Scenario represents rule engine scenario
type Scenario struct {
	runner *runner.Runner
	rpc    *rpc.Server

	enable bool
}

// New creates instance of Scenario
// instance contains rpc server that is not running, so Start must call.
func New() *Scenario {
	s := new(Scenario)

	s.rpc = rpc.NewServer()
	if err := s.rpc.Register(&Endpoint{s: s}); err != nil {
		panic(err)
	}

	return s
}

// Start runs scenario server
func (s *Scenario) Start() error {
	h := jsonrpc2.HTTPHandler(s.rpc)

	return http.ListenAndServe("127.0.0.1:1373", h)
}

// Stop scenario server and runner
func (s *Scenario) Stop() {
	// TODO stop json rpc
}

// Deactivate scenario runner not its server
func (s *Scenario) Deactivate() {
	if s.enable {
		s.runner.Stop()
	}
	s.enable = false
}

// Data new data is comming
func (s *Scenario) Data(d string, t string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if s.enable {
		if err := s.runner.DataEvent(ctx, d, map[string]string{
			"thing": t,
		}); err != nil {
			return
		}
	}
}

// Activate creates or replaces user scenario beacuase
// there is only one scenario is here with id
// it starts runner too
func (s *Scenario) Activate(code []byte, id string) error {
	f, err := os.Create(fmt.Sprintf("/tmp/scenario-%s.py", id))
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			return
		}
	}()
	if _, err = f.Write(code); err != nil {
		return err
	}

	s.startRunner(id)

	return nil
}

// ActivateWithoutCode just check user scenario existence
// and starts runner
func (s *Scenario) ActivateWithoutCode(id string) error {
	_, err := os.Stat(fmt.Sprintf("/tmp/scenario-%s.py", id))
	if err != nil {
		return err
	}

	s.startRunner(id)

	return nil
}

func (s *Scenario) startRunner(id string) {
	if s.enable {
		s.enable = false
		s.runner.Stop()
	}
	s.runner = runner.New(&runner.Task{
		Run: func(e runner.Event) (string, error) {
			cmd := exec.Command("runtime.py", "--job", "rule", "--id", e.Env("thing"), fmt.Sprintf("/tmp/scenario-%s.py", id))

			// stdin
			stdin, err := cmd.StdinPipe()
			if err != nil {
				return "", err
			}
			if _, err := io.WriteString(stdin, e.Data()); err != nil {
				return "", err
			}
			if err := stdin.Close(); err != nil {
				return "", err
			}

			// run and wait (blocking section)
			if _, err := cmd.Output(); err != nil {
				if err, ok := err.(*exec.ExitError); ok {
					return "", fmt.Errorf("%s: %s", err.Error(), err.Stderr)
				}
				return "", err
			}

			return "", nil
		},
		Interval: 0,
	}, 1024)
	s.enable = true

	s.runner.ErrHandler = func(err error) {
		log.WithFields(log.Fields{
			"code":      id,
			"job":       "rule",
			"project":   os.Getenv("NAME"),
			"component": "gorunner",
		}).Error(err)
	}

	go s.runner.Start()
}
