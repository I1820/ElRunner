/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 08-02-2018
 * |
 * | File Name:     linter/linter_test.go
 * +===============================================
 */

package linter

import (
	"encoding/json"
	"testing"
)

func TestBasic(t *testing.T) {
	code := []byte(`
from codec import Codec

class ISRC(Codec):
    def decode(self, data):
        return data.decode('ascii')
    def encode(self, data):
        pass
	`)

	out, err := Lint(code)
	if err != nil {
		t.Fatalf("Linter error: %s", err)
	}

	if out == "" {
		return
	}

	var jsn interface{}
	if err := json.Unmarshal([]byte(out), &jsn); err != nil {
		t.Fatalf("Linter results error: %s data: %q", err, out)
	}
	t.Log(jsn)
}
