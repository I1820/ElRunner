package grifts

import (
	"github.com/I1820/ElRunner/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
