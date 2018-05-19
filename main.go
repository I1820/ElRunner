/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 17-01-2018
 * |
 * | File Name:     main.go
 * +===============================================
 */

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aiotrc/GoRunner/codec"
	"github.com/aiotrc/GoRunner/linter"
	"github.com/aiotrc/GoRunner/scenario"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/weekface/mgorus"
)

var codecs map[string]*codec.Codec
var scr *scenario.Scenario

// init initiates global variables
func init() {
	codecs = make(map[string]*codec.Codec)
	scr = scenario.New()
}

// handle registers apis and create http handler
func handle() http.Handler {
	r := gin.Default()

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 Not Found"})
	})

	r.Use(gin.ErrorLogger())

	api := r.Group("/api")
	{
		api.POST("/decode/:id", decodeHandler)
		api.POST("/encode/:id", encodeHandler)

		api.GET("/about", aboutHandler)

		api.POST("/codec", codecHandler)
		api.POST("/scenario", scenarioHandler)
		api.GET("/scenario/:id/deactivate", scenarioDeactivateHandler)
		api.GET("/scenario/:id/activate", scenarioActivateHandler)

		api.POST("/lint", lintHandler)
	}

	return r
}

func main() {
	fmt.Println("GoRunner AIoTRC @ 2017")

	// Initiate logger
	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		mongoURL = "localhost:27017"
	}
	hooker, err := mgorus.NewHooker(mongoURL, "isrc", "errors")
	if err == nil {
		log.AddHook(hooker)
		log.Infof("Logrus MongoDB Hook is %s", mongoURL)
	} else {
		log.Errorf("Logrus MongoDB Hook %q error: %s", mongoURL, err)
	}

	go func() {
		log.Fatal(scr.Start())
	}()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handle(),
	}

	go func() {
		fmt.Printf("GoRunner Listen: %s\n", srv.Addr)
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("Listen Error:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("GoRunner Shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown Error:", err)
	}
}

func aboutHandler(c *gin.Context) {
	c.String(http.StatusOK, "18.20 is leaving us")
}

func encodeHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	id := c.Param("id")

	data, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	encoder, ok := codecs[id]
	if !ok {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("%s does not exit on GoRunner", id))
		return
	}

	parsed, err := encoder.Encode(string(data))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, parsed)
	}
}

func decodeHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	id := c.Param("id")

	var json []byte
	if err := c.BindJSON(&json); err != nil {
		return
	}

	decoder, ok := codecs[id]
	if !ok {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("%s does not exit on GoRunner", id))
		return
	}

	parsed, err := decoder.Decode(json)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		if scr.Enable {
			scr.Data(parsed, id)
		}
		c.Data(http.StatusOK, "application/json", []byte(parsed))
	}
}

func codecHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var json codeReq
	if err := c.BindJSON(&json); err != nil {
		return
	}
	id := json.ID

	codec, err := codec.New([]byte(json.Code), id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if codecs[id] != nil {
		codecs[id].Stop()
	}
	codecs[id] = codec

	c.JSON(http.StatusOK, id)
}

func lintHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var code string
	if err := c.BindJSON(&code); err != nil {
		return
	}

	jsn, err := linter.Lint(code)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Data(http.StatusOK, "application/json", []byte(jsn))
}

func scenarioHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var json codeReq
	if err := c.BindJSON(&json); err != nil {
		return
	}
	id := json.ID

	if err := scr.Code([]byte(json.Code), id); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, id)
}

func scenarioActivateHandler(c *gin.Context) {
	id := c.Param("id")

	scr.Enable = true

	c.JSON(http.StatusOK, id)
}

func scenarioDeactivateHandler(c *gin.Context) {
	id := c.Param("id")

	scr.Enable = false

	c.JSON(http.StatusOK, id)
}
