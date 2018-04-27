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
	"fmt"
	"io"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"

	"github.com/aiotrc/GoRunner/runner"
	"github.com/powerman/rpc-codec/jsonrpc2"
	log "github.com/sirupsen/logrus"
)

// Endpoint for scenario communication
type Endpoint struct {
	s *Scenario
}

// WaitForData waits for things incomming data
func (e *Endpoint) WaitForData(args int, reply *string) error {
	d := e.s.r.Trigger()
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
	r   *runner.Runner
	e   bool
	rpc *rpc.Server

	Enable bool
}

// New creates instance of Scenario
// instance contains rpc server that is not running, so Start must call.
func New() *Scenario {
	s := new(Scenario)

	s.Enable = true

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

// Stop stops scenario
func (s *Scenario) Stop() {
	if s.e {
		s.r.Stop()
	}
}

// Data new data is comming
func (s *Scenario) Data(d string, t string) {
	if s.e {
		s.r.DataEvent(d, map[string]string{
			"thing": t,
		})
	}
}

// Code creates or replaces scenario beacuase
// there is only on scenario
func (s *Scenario) Code(code []byte, id string) error {
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

	if s.e {
		s.r.Stop()
	}
	s.e = true
	s.r = runner.New(&runner.Task{
		Run: func(e runner.Event) (string, error) {
			cmd := exec.Command("runtime.py", "--job", "rule", "--id", e.Env("thing"), f.Name())

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

			// run
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

	s.r.ErrHandler = func(err error) {
		log.WithFields(log.Fields{
			"code":      id,
			"job":       "rule",
			"project":   os.Getenv("NAME"),
			"component": "gorunner",
		}).Error(err)
	}

	go s.r.Start()

	return nil
}
