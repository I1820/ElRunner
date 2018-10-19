package actions

import (
	"fmt"

	linkapp "github.com/I1820/ElRunner/app"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	contenttype "github.com/gobuffalo/mw-contenttype"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/sirupsen/logrus"
	"github.com/weekface/mgorus"

	"github.com/gobuffalo/x/sessions"
	"github.com/rs/cors"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var linkApp *linkapp.Application

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			PreWares: []buffalo.PreWare{
				cors.Default().Handler,
			},
			SessionName: "_el_runner_session",
		})

		// If no content type is sent by the client
		// the application/json will be set, otherwise the client's
		// content type will be used.
		app.Use(contenttype.Add("application/json"))

		// Create mongodb connection
		url := envy.Get("DB_URL", "mongodb://172.18.0.1:27017")
		hooker, err := mgorus.NewHooker(url, "i1820", fmt.Sprintf("projects.logs.%s", envy.Get("PROJECT", "ElRunner")))
		if err == nil {
			logrus.AddHook(hooker)
			logrus.Infof("Logrus MongoDB Hook is %s", url)
		} else {
			logrus.Errorf("Logrus MongoDB Hook %q error: %s", url, err)
		}

		if ENV == "development" {
			app.Use(paramlogger.ParameterLogger)
		}

		// LinkApp initiation
		linkApp = linkapp.New(envy.Get("PROJECT", "ElRunner"))
		linkApp.Run()

		// Routes
		app.GET("/about", AboutHandler)
		api := app.Group("/api")
		{
			cr := CodecsResource{}
			api.Resource("/codecs", cr)
			api.POST("/codecs/{codec_id}/decode", cr.Decode)
			api.POST("/codecs/{codec_id}/encode", cr.Encode)

			sr := ScenariosResource{}
			api.GET("/scenarios/main", sr.Main)
			api.GET("/scenarios/deactivate", sr.Deactivate)
			api.Resource("/scenarios", sr)
			api.GET("/scenarios/{scenario_id}/activate", sr.Activate)

			api.POST("/lint", LintHandler)
		}
	}

	return app
}
