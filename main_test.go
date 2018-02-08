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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	id = "isrc-sensor"

	herTextName = "Kiana"
	herCborName = []byte{0x65, 0x4B, 0x69, 0x61, 0x6E, 0x61}

	code = `
from codec import Codec
import cbor

class ISRC(Codec):
    def decode(self, data):
        return cbor.loads(data)
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

func TestCodec(t *testing.T) {
	h := handle()
	s := httptest.NewServer(h)
	defer s.Close()

	// Upload codec
	resp, err := http.Post(fmt.Sprintf("%s/api/codec/%s", s.URL, id), "text/plain", bytes.NewBufferString(code))
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

	if string(body) != id {
		t.Fatalf("%q != %q", string(body), id)
	}

	// Encode her name
	resp, err = http.Post(fmt.Sprintf("%s/api/encode/%s", s.URL, id), "text/plain", bytes.NewBufferString(herTextName))
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

	if !bytes.Equal(body, herCborName) {
		t.Fatalf("%q != %q", body, herCborName)
	}

	// Decode her name
	resp, err = http.Post(fmt.Sprintf("%s/api/decode/%s", s.URL, id), "text/plain", bytes.NewBuffer(herCborName))
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

	if string(body) != herTextName+"\n" {
		t.Fatalf("%q != %q", string(body), herTextName)
	}
}
