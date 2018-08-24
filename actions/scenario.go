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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	scenarios := make([]string, 0)

	files, err := ioutil.ReadDir("/tmp")
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	for _, f := range files {
		name := f.Name()
		if s := scenarioRegexp.FindStringSubmatch(name); len(s) > 0 && s[0] == name {
			scenarios = append(scenarios, s[1])
		}
	}

	return c.Render(http.StatusOK, r.JSON(scenarios))
}

// Create creates new scenario and stores it code. This function is mapped
// to the path POST /scenarios
func (ScenariosResource) Create(c buffalo.Context) error {
	var rq codeReq
	if err := c.Bind(&rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	id := rq.ID
	code := []byte(rq.Code)

	f, err := os.Create(fmt.Sprintf("/tmp/scenario-%s.py", id))
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			return
		}
	}()
	if _, err = f.Write(code); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.JSON(id))
}

// Show shows uploaded scenario code. This function is mapped
// to the path GET /scenarios/{scenario_id}
func (ScenariosResource) Show(c buffalo.Context) error {
	b, err := ioutil.ReadFile(fmt.Sprintf("/tmp/scenario-%s.py", c.Param("scenario_id")))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(string(b)))
}

// Activate activates scenario. This function is mapped
// to the path GET /scenarios/{scenario_id}/activate
func (ScenariosResource) Activate(c buffalo.Context) error {
	id := c.Param("scenario_id")

	if err := linkApp.Scenario(id); err != nil {
		c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(id))
}
