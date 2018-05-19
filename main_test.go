/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 17-01-2018
 * |
 * | File Name:     main_test.go
 * +===============================================
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	id = "isrc-sensor"

	herTextName = "Shamin"
	herCborName = []byte{0x66, 0x53, 0x68, 0x61, 0x6D, 0x69, 0x6E}

	locationCbor = []byte{0xA3, 0x63, 0x6C, 0x61, 0x74, 0x0A, 0x63, 0x6C, 0x6E, 0x67, 0x0A, 0x6B, 0x74, 0x65, 0x6D, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x1E}

	codecOne = `
from codec import Codec
import cbor

class ISRC(Codec):
    def decode(self, data):
        return cbor.loads(data)
    def encode(self, data):
        return cbor.dumps(data)
`
	codecTwo = `
from codec import Codec
import cbor

class ISRC(Codec):
    thing_location = 'loc'

    def decode(self, data):
        print("Hello")
        d = cbor.loads(data)
        print(d)

        d['loc'] = self.create_location(d['lat'], d['lng'])

        return d
    def encode(self, data):
        return cbor.dumps(data)

`
)

func TestAbout(t *testing.T) {
	h := handle()
	s := httptest.NewServer(h)
	defer s.Close()

	resp, err := http.Get(fmt.Sprintf("%s/api/about", s.URL))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "18.20 is leaving us" {
		t.Fatalf("who leaving us?! %q", body)
	}
}

func TestCodec1(t *testing.T) {
	h := handle()
	s := httptest.NewServer(h)
	defer s.Close()

	// Upload codec
	jsc, err := json.Marshal(codeReq{
		ID:   id,
		Code: codecOne,
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/api/codec", s.URL), "application/json", bytes.NewBuffer(jsc))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "\""+id+"\"" {
		t.Fatalf("%q != %q", string(body), id)
	}

	// Encode her name
	resp, err = http.Post(fmt.Sprintf("%s/api/encode/%s", s.URL, id), "application/json", bytes.NewBufferString("\""+herTextName+"\""))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var jse []byte
	if err := json.Unmarshal(body, &jse); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(jse, herCborName) {
		t.Fatalf("%q != %q", body, herCborName)
	}

	// Decode her name
	jsd, err := json.Marshal(herCborName)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = http.Post(fmt.Sprintf("%s/api/decode/%s", s.URL, id), "application/json", bytes.NewBuffer(jsd))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "\""+herTextName+"\"\n" {
		t.Fatalf("%q != %q", string(body), herTextName)
	}
}

func TestCodec2(t *testing.T) {
	h := handle()
	s := httptest.NewServer(h)
	defer s.Close()

	// Upload codec
	jsc, err := json.Marshal(codeReq{
		ID:   id,
		Code: codecTwo,
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/api/codec", s.URL), "application/json", bytes.NewBuffer(jsc))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "\""+id+"\"" {
		t.Fatalf("%q != %q", string(body), id)
	}

	// Decode location
	jsd, err := json.Marshal(locationCbor)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = http.Post(fmt.Sprintf("%s/api/decode/%s", s.URL, id), "application/json", bytes.NewBuffer(jsd))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
}
