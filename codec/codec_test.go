/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 10-11-2017
 * |
 * | File Name:     codec/codec_test.go
 * +===============================================
 */

package codec

import "testing"

func TestHelloDeocder(t *testing.T) {
	code := []byte(`
from codec import Codec

class ISRC(Codec):
    def decode(self, data):
        return data.decode('ascii')
    def encode(self, data):
        pass
	`)
	d, err := New(code, "hi")
	if err != nil {
		t.Fatal(err)
	}
	r, err := d.Decode([]byte("Hi"))
	if err != nil {
		t.Fatal(err)
	}
	if r != "\"Hi\"\n" {
		t.Fatalf("Invalid Decode Result %q", r)
	}
	d.Stop()
}

func TestFaultyDecoder(t *testing.T) {
	code := []byte(`
from codec import Codec

class ISRC(Codec):
    def decode(self, data):
        khar
    def encode(self, data):
        pass
	`)
	d, err := New(code, "hi")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := d.Decode([]byte("Hi")); err != nil {
		t.Log(err)
	} else {
		t.Fatal("There is no error?!")
	}
	d.Stop()

}

func TestHelloEncoder(t *testing.T) {
	code := []byte(`
from codec import Codec

class ISRC(Codec):
    def decode(self, data):
        pass
    def encode(self, data):
        return data.encode('ascii')
	`)
	d, err := New(code, "hi")
	if err != nil {
		t.Fatal(err)
	}
	r, err := d.Encode("Hi")
	if err != nil {
		t.Fatal(err)
	}
	if string(r) != "Hi" {
		t.Fatalf("Invalid Decode Result %q", r)
	}
	d.Stop()

}
