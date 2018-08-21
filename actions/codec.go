/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 21-08-2018
 * |
 * | File Name:     codec.go
 * +===============================================
 */

package actions

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/I1820/ElRunner/codec"
	"github.com/gobuffalo/buffalo"
)

// scenario, codec request payload
type codeReq struct {
	ID   string `json:"id" binding:"required"`
	Code string `json:"code" binding:"required"`
}

// CodecsResource manages existing codecs
type CodecsResource struct {
	buffalo.Resource
}

var codecRegexp *regexp.Regexp

func init() {
	rg, err := regexp.Compile(`codec-(\w*).py`)
	if err == nil {
		codecRegexp = rg
	}
}

// List lists available codecs. This function is mapped
// to the path GET /codecs
func (CodecsResource) List(c buffalo.Context) error {
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

// Create creates new codec and stores it code. This function is mapped
// to the path POST /codecs
func (CodecsResource) Create(c buffalo.Context) error {
	var rq codeReq
	if err := c.Bind(&rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	id := rq.ID

	_, err := codec.New([]byte(rq.Code), id)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(id))
}

// Show shows uploaded codec code. This function is mapped
// to the path GET /codecs/{codec_id}
func (CodecsResource) Show(c buffalo.Context) error {
	b, err := ioutil.ReadFile(fmt.Sprintf("/tmp/codec-%s.py", c.Param("codec_id")))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(string(b)))
}

// Destroy removes uploaded codec. This function is mapped
// to the path DELETE /codecs/{codec_id}
func (CodecsResource) Destroy(c buffalo.Context) error {
	if err := os.Remove(fmt.Sprintf("/tmp/codec-%s.py", c.Param("codec_id"))); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(c.Param("codec_id")))
}
