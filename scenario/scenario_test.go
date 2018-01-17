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
	"net/http/httptest"
	"testing"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

func TestAbout(t *testing.T) {
	h := jsonrpc2.HTTPHandler(New().rpc)
	s := httptest.NewServer(h)
	defer s.Close()

	var about string
	c := jsonrpc2.NewHTTPClient(s.URL)
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
