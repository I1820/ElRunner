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
        f = open('/tmp/hello', 'w')
        f.write(str(data))
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

	if _, err := s.r.OutputBoundedWait(10 * time.Millisecond); err != nil {
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
		t.Fatal("%q != {'Hello': 10}", string(data))
	}
}
