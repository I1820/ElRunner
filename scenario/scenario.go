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
)

// Endpoint for scenario communication
type Endpoint struct {
}

// WaitForData waits for things incomming data
func (e *Endpoint) WaitForData(args int, reply *string) error {
	fmt.Println("Helloooooo")
	*reply = "YesYes"
	return nil
}

// About let you about who left us alone
func (e Endpoint) About(args int, reply *string) error {
	*reply = "18.20 is leaving us"
	return nil
}

// Scenario represents rule engine scenario
type Scenario struct {
	r   runner.Runner
	rpc *rpc.Server
}

// New creates instance of Scenario
// instance contains rpc server that is not running, so Start must call.
func New() *Scenario {
	s := new(Scenario)

	s.rpc = rpc.NewServer()
	if err := s.rpc.Register(new(Endpoint)); err != nil {
		panic(err)
	}

	return s
}

// Start runs scenario server
func (s *Scenario) Start() error {
	h := jsonrpc2.HTTPHandler(s.rpc)

	return http.ListenAndServe("127.0.0.1:1373", h)
}

// Data new data is comming
func (s *Scenario) Data(d string) {
	s.r.DataEvent(d)
}

// Code creates or replaces scenario beacuase
// there is only on scenario
func (s *Scenario) Code(code []byte, id string) error {
	f, err := os.Create(fmt.Sprintf("/tmp/scenario-%s.py", id))
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.Write(code); err != nil {
		return err
	}

	s.r = runner.New(&runner.Task{
		Run: func(e runner.Event) (string, error) {
			cmd := exec.Command("runtime.py", "--job", "rule", f.Name())

			// stdin
			stdin, err := cmd.StdinPipe()
			if err != nil {
				return "", err
			}
			if _, err := io.WriteString(stdin, e.Data()); err != nil {
				return "", err
			}
			stdin.Close()

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

	return nil
}
