/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 10-11-2017
 * |
 * | File Name:     codec/codec.go
 * +===============================================
 */

package codec

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type command int

const (
	encodeCommand command = iota
	decodeCommand command = iota
)

// Codec provides Encoder/Decode functions
// based on user codec script
type Codec struct {
	id string
}

// Decode sends base64 coded data into stdin of
// user codec script and reads string represntation of json data object
// from its stdout
func (c *Codec) Decode(ctx context.Context, r []byte) (string, error) {
	return c.exec(ctx, decodeCommand, base64.StdEncoding.EncodeToString(r))
}

// Encode sends string represntation of json data object into stdin of
// user codec script and reads base64 coded data
// from its stdout
func (c *Codec) Encode(ctx context.Context, p string) ([]byte, error) {
	s, err := c.exec(ctx, encodeCommand, p)
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(s)
}

// exec runs user script for encode/decode
func (c *Codec) exec(ctx context.Context, t command, input string) (string, error) {
	var cmd *exec.Cmd

	switch t {
	case decodeCommand:
		cmd = exec.CommandContext(ctx, "runtime.py", "--job", "decode", fmt.Sprintf("/tmp/codec-%s.py", c.id))
	case encodeCommand:
		cmd = exec.CommandContext(ctx, "runtime.py", "--job", "encode", fmt.Sprintf("/tmp/codec-%s.py", c.id))
	}

	// stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	if _, err := io.WriteString(stdin, input); err != nil {
		return "", err
	}
	if err := stdin.Close(); err != nil {
		return "", err
	}

	// stdout
	out, err := cmd.Output()
	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("%s: %s", err.Error(), err.Stderr)
		}
		return "", err
	}

	return string(out), nil

}

// New creates encoder/decoder based on given code
// it store user script under /tmp/codec-id.py
func New(code []byte, id string) (*Codec, error) {
	f, err := os.Create(fmt.Sprintf("/tmp/codec-%s.py", id))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			return
		}
	}()
	if _, err = f.Write(code); err != nil {
		return nil, err
	}

	return &Codec{
		id: id,
	}, nil
}

// NewWithoutCode creates encoder/decoder without any file creation
// but it check user script existence under /tmp/codec-id.py
func NewWithoutCode(id string) (*Codec, error) {
	if _, err := os.Stat(fmt.Sprintf("/tmp/codec-%s.py", id)); err != nil {
		return nil, err
	}

	return &Codec{
		id: id,
	}, nil
}
