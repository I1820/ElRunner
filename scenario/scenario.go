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
	"net/http"
	"net/rpc"

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
func New() *Scenario {
	s := new(Scenario)

	s.rpc = rpc.NewServer()
	s.rpc.Register(new(Endpoint))

	return s
}

// Start runs scenario server
func (s *Scenario) Start() error {
	h := jsonrpc2.HTTPHandler(s.rpc)

	return http.ListenAndServe("127.0.0.1:1373", h)
}