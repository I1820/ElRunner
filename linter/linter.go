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
	"time"
)

// Lint lints given python code with pylint
func Lint(code []byte) (string, error) {
	f, err := os.Create(fmt.Sprintf("/tmp/%d.py", time.Now().Unix()))
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err = f.Write(code); err != nil {
		return "", err
	}

	cmd := exec.Command("pylint", "--output-format", "json", "-j", "2", "-E", f.Name())

	out, err := cmd.Output()
	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("%s: %s", err.Error(), err.Stderr)
		}
		return "", err
	}

	return string(out), nil
}
