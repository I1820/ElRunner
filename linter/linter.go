/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 08-02-2018
 * |
 * | File Name:     linter/linter.go
 * +===============================================
 */

package linter

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// Lint lints given python code with pylint
func Lint(code []byte) (string, error) {
	f, err := os.Create(fmt.Sprintf("/tmp/linter-%d.py", time.Now().Unix()))
	if err != nil {
		return "", err
	}
	defer func() {
		if err := f.Close(); err != nil {
			return
		}
	}()
	if _, err = f.Write(code); err != nil {
		return "", err
	}

	cmd := exec.Command("pylint", "--output-format", "json", "-j", "2", f.Name())

	out, err := cmd.Output()
	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			if err.Sys().(syscall.WaitStatus).ExitStatus() == 32 {
				return "", fmt.Errorf("%s: %s", err.Error(), err.Stderr)
			}
		}
	}

	return string(out), nil
}
