/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 09-11-2017
 * |
 * | File Name:     runner/output.go
 * +===============================================
 */

package runner

import "fmt"

// Output represents output of the task that runs on runner
type Output *output
type output struct {
	o string
	e int
}

// CreateOutput creates output from given variable using fmt.Sprint
// s: output message
// e: error code (0 = no-error)
func CreateOutput(v interface{}, e int) Output {
	return &output{
		o: fmt.Sprint(v),
		e: e,
	}
}
