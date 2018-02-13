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
	log "github.com/sirupsen/logrus"
)

// Codec provides Encoder/Decoder binded functions
type Codec struct {
	r *runner.Runner
	i string
}

// Decode sends base64 coded data into stdin of
// user codec script and reads string represntation of json data object
// from its stdout
func (c *Codec) Decode(r []byte) (string, error) {
	c.r.Event(DecodeEvent(base64.StdEncoding.EncodeToString(r)))
	return c.r.Output()
}

// Encode sends string represntation of json data object into stdin of
// user codec script and reads base64 coded data
// from its stdout
func (c *Codec) Encode(p string) ([]byte, error) {
	c.r.Event(EncodeEvent(p))
	s, e := c.r.Output()
	if e != nil {
		return nil, e
	}
	return base64.StdEncoding.DecodeString(s)
}

// Stop stops codec runner
func (c *Codec) Stop() {
	c.r.Stop()
}

// ID returns codec identification
func (c *Codec) ID() string {
	return c.i
}

// New creates encoder/decoder based on given data
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
		},
		Interval: 0,
	}, 1)
	runner.ErrHandler = func(err error) {
		log.WithFields(log.Fields{
			"Project": id,
		}).Error(err)
	}

	go runner.Start()

	return &Codec{
		r: runner,
		i: id,
	}, nil
}
