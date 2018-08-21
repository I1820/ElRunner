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

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloDeocder(t *testing.T) {
	code := []byte(`
from codec import Codec

class ISRC(Codec):
    def decode(self, data):
        return data.decode('ascii')
    def encode(self, data):
        pass
`)
	d, err := New(code, "hello")
	assert.NoError(t, err)

	r, err := d.Decode(context.Background(), []byte("Hi"))
	assert.NoError(t, err)

	assert.Equalf(t, "Hi", r, "Invalid Decode Result %q", r)
}

func TestPersistHelloDecoder(t *testing.T) {
	d, err := NewWithoutCode("hello")
	assert.NoError(t, err)

	r, err := d.Decode(context.Background(), []byte("Hi"))
	assert.NoError(t, err)

	assert.Equalf(t, "Hi", r, "Invalid Decode Result %q", r)
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
	d, err := New(code, "fault")
	assert.NoError(t, err)

	_, err = d.Decode(context.Background(), []byte("Hi"))
	assert.Error(t, err)
}

func TestHiEncoder(t *testing.T) {
	code := []byte(`
from codec import Codec

class ISRC(Codec):
    def decode(self, data):
        pass
    def encode(self, data):
        return data.encode('ascii')
`)
	d, err := New(code, "hi")
	assert.NoError(t, err)

	r, err := d.Encode(context.Background(), "Hi")
	assert.NoError(t, err)

	assert.Equalf(t, "Hi", string(r), "Invalid Encode Result %q", r)
}
