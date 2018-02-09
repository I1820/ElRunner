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
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"
	"time"

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

func TestHelloScenario(t *testing.T) {
	code := []byte(`
from scenario import Scenario

class ISRC(Scenario):
    def run(self, data):
        f = open('/tmp/hello', 'w+')
        f.write(str(data))
        f.close()
	`)

	s := New()

	h := jsonrpc2.HTTPHandler(s.rpc)
	srv := httptest.NewServer(h)
	defer srv.Close()

	if err := s.Code(code, "Hello"); err != nil {
		t.Fatal(err)
	}
	defer s.Stop()

	s.Data("{\"Hello\": 10}")

	if _, err := s.r.OutputBoundedWait(1 * time.Second); err != nil {
		t.Fatalf("Runs error: %s", err)
	}

	f, err := os.Open("/tmp/hello")
	if err != nil {
		t.Fatalf("Could not open /tmp/hello: %s", err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("Could not read /tmp/hello: %s", err)
	}

	if string(data) != "{'Hello': 10}" {
		t.Fatalf("%q != {'Hello': 10}", string(data))
	}
}

func TestRPCScenario(t *testing.T) {
	code := []byte(`
import asyncio
from scenario import Scenario

class ISRC(Scenario):
    def run(self, data):
        f = open('/tmp/rpc', 'w+')
        f.write(str(data))
        f.write('\n')

        loop = asyncio.get_event_loop()
        t = asyncio.ensure_future(self.wait_for_data(timeout=1000))
        loop.run_until_complete(t)
        loop.close()
        response = t.result()
        if response is not None:
            f.write(str(response))
        else:
            f.write(str(t.done()))
        f.close()
	`)

	s := New()
	go func() {
		if err := s.Start(); err != nil {
			t.Fatal(err)
		}
	}()

	if err := s.Code(code, "RPC"); err != nil {
		t.Fatal(err)
	}
	defer s.Stop()

	s.Data("{\"Hello\": 10}")
	s.Data("{\"Hello\": 9}")

	if _, err := s.r.OutputBoundedWait(1 * time.Second); err != nil {
		t.Fatalf("Runs error: %s", err)
	}

	f, err := os.Open("/tmp/rpc")
	if err != nil {
		t.Fatalf("Could not open /tmp/rpc: %s", err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("Could not read /tmp/rpc: %s", err)
	}

	fmt.Println(string(data))
}
