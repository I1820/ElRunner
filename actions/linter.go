/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 22-08-2018
 * |
 * | File Name:     linter.go
 * +===============================================
 */

package actions

import (
	"io"
	"net/http"

	"github.com/I1820/ElRunner/linter"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

// LintHandler lints your given code
func LintHandler(c buffalo.Context) error {
	var code string
	if err := c.Bind(&code); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	jsn, err := linter.Lint([]byte(code))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.Func("application/json", func(w io.Writer, d render.Data) error {
		_, err := io.WriteString(w, jsn)
		return err
	}))
}
