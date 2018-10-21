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
	"context"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/stretchr/testify/assert"
)

func TestAbout(t *testing.T) {
	h := jsonrpc2.HTTPHandler(New().rpc)
	s := httptest.NewServer(h)
	defer s.Close()

	var about string
	c := jsonrpc2.NewHTTPClient(s.URL)
	assert.NoError(t, c.Client.Call("Endpoint.About", nil, &about))
	assert.Equal(t, about, "18.20 is leaving us")
}

func TestHelloScenario(t *testing.T) {
	code := []byte(`
from scenario import Scenario

class Fanco(Scenario):
    async def run(self, data):
        f = open('/tmp/hello', 'w+')
        f.write(str(data))
        f.close()
	`)

	s := New()

	h := jsonrpc2.HTTPHandler(s.rpc)
	srv := httptest.NewServer(h)
	defer srv.Close()

	if err := s.Activate(code, "Hello"); err != nil {
		t.Fatal(err)
	}
	defer s.Deactivate()

	s.Data("{\"Hello\": 10}", "Parham")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := s.runner.Output(ctx)
	assert.NoError(t, err)

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

class Fanco(Scenario):
    async def run(self, data):
        f = open('/tmp/rpc', 'w+')
        f.write(str(data))
        f.write('\n')

        response = await self.wait_for_data(timeout=1000)
        if response is not None:
            f.write(str(response))
        else:
            f.write('timeout')
        f.close()
        `)

	s := New()
	go func() {
		assert.NoError(t, s.Start())
	}()

	assert.NoError(t, s.Activate(code, "RPC"))
	defer s.Deactivate()

	s.Data("{\"Hello\": 10}", "Parham")
	s.Data("{\"Hello\": 9}", "Parham")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := s.runner.Output(ctx)
	assert.NoError(t, err)

	f, err := os.Open("/tmp/rpc")
	assert.NoError(t, err)

	data, err := ioutil.ReadAll(f)
	assert.NoError(t, err)

	assert.Equal(t, data, []byte("{\"Hello\": 10}\n{\"Hello\": 9}"))

	t.Log(string(data))
}

func TestEmailScenario(t *testing.T) {
	code := []byte(`
from scenario import Scenario


class Fanco(Scenario):
    async def run(self, data):
        sender = 'platform.avidnetco@gmail.com'
        receivers = ['parham.alvani@gmail.com']

        message = 'From: From Travis CI of Avidnet <platform.avidnetco@gmail.com>\n' \
                  'To: To Parham Alvani <parham.alvani@gmail.com>\n' \
                  'Subject: Rule Engine Notification (Avidnet-Travis)\n\n' \
                  'Data:' + str(data) + '\n' \
                                        'Sent by Rule Engine. Scenario'
        self.send_email(host='smtp.gmail.com', port=587,
                        username="platform.avidnetco@gmail.com",
                        password="fancopass(1397)",
                        sender=sender,
                        receivers=receivers, message=message)
	`)

	s := New()

	h := jsonrpc2.HTTPHandler(s.rpc)
	srv := httptest.NewServer(h)
	defer srv.Close()

	assert.NoError(t, s.Activate(code, "Email"))
	defer s.Deactivate()

	s.Data("{\"Hello\": 10}", "Parham")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := s.runner.Output(ctx)
	assert.NoError(t, err)
}
