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
import time

s = input()
print("hello from python", s)
	`)
	d, err := New(code, "hi")
	if err != nil {
		t.Fatal(err)
	}
	r := d.Decode("Hi")
	if r != "hello from python Hi\n" {
		t.Fatal("Invalid Decode Result", r)
	}
	d.Stop()
}
