/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 10-11-2017
 * |
 * | File Name:     decoder/decoder_test.go
 * +===============================================
 */

package decoder

import "testing"

func TestHelloDeocder(t *testing.T) {
	code := []byte(`
class ISRC(Codec, requirements=[]):
    def decode(self, data):
        return data.decode('ascii')
    def encode(self, data):
        pass
	`)
	d, err := New(code, "hi")
	if err != nil {
		t.Fatal(err)
	}
	r := d.Decode([]byte("Hi"))
	if r != "Hi\n" {
		t.Fatal("Invalid Decode Result \"", r, "\"")
	}
	d.Stop()
}
