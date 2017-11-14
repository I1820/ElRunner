/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 10-11-2017
 * |
 * | File Name:     decoder/decoder.go
 * +===============================================
 */

package decoder

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/aiotrc/GoRunner/runner"
)

// Decoder binded functions
type Decoder interface {
	Decode([]byte) string
	ID() string
	Stop()
}

type decoder struct {
	r runner.Runner
	i string
}

func (d *decoder) Decode(r []byte) string {
	d.r.Trigger(r)
	return string(d.r.Output())
}

func (d *decoder) Stop() {
	d.r.Stop()
}

func (d *decoder) ID() string {
	return d.i
}

// New creates decoder based on given code
func New(code []byte, id string) (Decoder, error) {
	f, err := os.Create(fmt.Sprintf("/tmp/%s.py", id))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if _, err = f.Write(code); err != nil {
		return nil, err
	}

	runner := runner.New(&runner.Task{
		Run: func(e runner.Event) runner.Output {
			cmd := exec.Command("python3", fmt.Sprintf("/tmp/%s.py", id))

			stdin, err := cmd.StdinPipe()
			if err != nil {
				return runner.Output(err.Error())
			}
			io.WriteString(stdin, base64.StdEncoding.EncodeToString(e.Data()))
			stdin.Close()

			out, err := cmd.Output()
			if err != nil {
				return runner.Output(err.Error())
			}

			return runner.Output(out)
		},
		Interval: 0,
	}, 1)
	go runner.Start()

	return &decoder{
		r: runner,
		i: id,
	}, nil
}
