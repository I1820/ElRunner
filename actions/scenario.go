/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 21-08-2018
 * |
 * | File Name:     scenario.go
 * +===============================================
 */

package actions

import (
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gobuffalo/buffalo"
)

// ScenariosResource manages existing scenarios
type ScenariosResource struct {
	buffalo.Resource
}

var scenarioRegexp *regexp.Regexp

func init() {
	rg, err := regexp.Compile(`scenario-(\w*).py`)
	if err == nil {
		scenarioRegexp = rg
	}
}

// List lists available scenarios. This function is mapped
// to the path GET /scenarios
func (ScenariosResource) List(c buffalo.Context) error {
	codecs := make([]string, 0)

	files, err := ioutil.ReadDir("/tmp")
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	for _, f := range files {
		name := f.Name()
		if s := codecRegexp.FindStringSubmatch(name); len(s) > 0 && s[0] == name {
			codecs = append(codecs, s[1])
		}
	}

	return c.Render(http.StatusOK, r.JSON(codecs))
}
