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
	"testing"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

func TestAbout(t *testing.T) {
	New()

	var about string
	c := jsonrpc2.NewHTTPClient("http://127.0.0.1:1373")
	c.Client.Call("Endpoint.About", nil, &about)
	if about != "18.20 is leaving us" {
		t.Fatal("who leaving us?!")
	}
}
