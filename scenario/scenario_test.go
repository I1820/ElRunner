/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 24-12-2017
 * |
 * | File Name:     scenario_test.go
 * +===============================================
 */

package scenario

import (
	"fmt"
	"os"
	"testing"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

func TestMain(m *testing.M) {
	s := New()

	go func() {
		err := s.Start()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	os.Exit(m.Run())
}

func TestAbout(t *testing.T) {
	var about string
	c := jsonrpc2.NewHTTPClient("http://127.0.0.1:1373")
	err := c.Client.Call("Endpoint.About", nil, &about)
	if err != nil {
		t.Fatalf("Call Endpoint.About: %s", err)
	}
	if about != "18.20 is leaving us" {
		t.Fatalf("who leaving us?! %q", about)
	}
}

func TestListen(t *testing.T) {
}
