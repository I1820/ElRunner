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
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/aiotrc/GoRunner/runner"
)

// Codec provides Encoder/Decoder binded functions
type Codec interface {
	Decode([]byte) (string, error)
	Encode(string) ([]byte, error)
	ID() string
	Stop()
}

type codec struct {
	r runner.Runner
	i string
}

// Decode sends base64 coded data into stdin of
// user codec script and reads string represntation of json data object
// from its stdout
func (c *codec) Decode(r []byte) (string, error) {
	c.r.Event(DecodeEvent(base64.StdEncoding.EncodeToString(r)))
	return c.r.Output()
}

// Encode sends string represntation of json data object into stdin of
// user codec script and reads base64 coded data
// from its stdout
func (c *codec) Encode(p string) ([]byte, error) {
	c.r.Event(EncodeEvent(p))
	s, e := c.r.Output()
	if e != nil {
		return nil, e
	}
	return base64.StdEncoding.DecodeString(s)
}

func (c *codec) Stop() {
	c.r.Stop()
}

func (c *codec) ID() string {
	return c.i
}

// New creates encoder/decoder based on given data
func New(code []byte, id string) (Codec, error) {
	f, err := os.Create(fmt.Sprintf("/tmp/codec-%s.py", id))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if _, err = f.Write(code); err != nil {
		return nil, err
	}

	runner := runner.New(&runner.Task{
		Run: func(e runner.Event) (string, error) {
			var cmd *exec.Cmd

			// command
			switch e.(type) {
			case DecodeEvent:
				cmd = exec.Command("runtime.py", "--job", "decode", f.Name())
			case EncodeEvent:
				cmd = exec.Command("runtime.py", "--job", "encode", f.Name())
			default:
				return "", nil
			}

			// stdin
			stdin, err := cmd.StdinPipe()
			if err != nil {
				return "", err
			}
			if _, err := io.WriteString(stdin, e.Data()); err != nil {
				return "", err
			}
			stdin.Close()

			// stdout
			out, err := cmd.Output()
			if err != nil {
				if err, ok := err.(*exec.ExitError); ok {
					return "", fmt.Errorf("%s: %s", err.Error(), err.Stderr)
				}
				return "", err
			}

			return string(out), nil
		},
		Interval: 0,
	}, 1)
	go runner.Start()

	return &codec{
		r: runner,
		i: id,
	}, nil
}
