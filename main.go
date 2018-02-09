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
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aiotrc/GoRunner/codec"
	"github.com/aiotrc/GoRunner/linter"
	"github.com/aiotrc/GoRunner/scenario"
	"github.com/gin-gonic/gin"
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

	api := r.Group("/api")
	{
		api.POST("/decode/:id", decodeHandler)
		api.POST("/encode/:id", encodeHandler)

		api.GET("/about", aboutHandler)

		api.POST("/codec/:id", codecHandler)
		api.POST("/scenario/:id", scenarioHandler)

		api.POST("/lint", lintHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 Not Found"})
	})

	return r
}

func main() {
	fmt.Println("GoRunner AIoTRC @ 2017")

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
	id := c.Param("id")
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	encoder, ok := codecs[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("%q does not exit on GoRunner", id)})
		return
	}

	parsed, err := encoder.Encode(string(data))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.Data(http.StatusOK, "application/octet-stream", parsed)
	}
}

func decodeHandler(c *gin.Context) {
	id := c.Param("id")
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	decoder, ok := codecs[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("%q does not exit on GoRunner", id)})
		return
	}

	parsed, err := decoder.Decode(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		scr.Data(parsed)
		c.String(http.StatusOK, parsed)
	}
}

func codecHandler(c *gin.Context) {
	id := c.Param("id")
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	codec, err := codec.New(data, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if codecs[id] != nil {
		codecs[id].Stop()
	}
	codecs[id] = codec

	c.String(http.StatusOK, id)
}

func lintHandler(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsn, err := linter.Lint(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", []byte(jsn))
}

func scenarioHandler(c *gin.Context) {
	id := c.Param("id")
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := scr.Code(data, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, id)
}
