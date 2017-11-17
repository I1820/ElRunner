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

func (c *codec) Decode(r []byte) (string, error) {
	c.r.Event(DecodeEvent(r))
	return c.r.Output()
}

func (c *codec) Encode(p string) ([]byte, error) {
	c.r.Event(EncodeEvent(p))
	s, e := c.r.Output()
	return []byte(s), e
}

func (c *codec) Stop() {
	c.r.Stop()
}

func (c *codec) ID() string {
	return c.i
}

// New creates decoder based on given code
func New(code []byte, id string) (Codec, error) {
	f, err := os.Create(fmt.Sprintf("/tmp/%s.py", id))
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
			switch e.(type) {
			case DecodeEvent:
				cmd = exec.Command("runtime.py", "--job", "decode", fmt.Sprintf("/tmp/%s.py", id))
			case EncodeEvent:
				cmd = exec.Command("runtime.py", "--job", "encode", fmt.Sprintf("/tmp/%s.py", id))
			default:
				return "", nil
			}

			stdin, err := cmd.StdinPipe()
			if err != nil {
				return "", err
			}
			io.WriteString(stdin, base64.StdEncoding.EncodeToString(e.Data()))
			stdin.Close()

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
