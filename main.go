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
	"github.com/gin-gonic/gin"
)

var codecs map[string]codec.Codec

// init initiates global variables
func init() {
	codecs = make(map[string]codec.Codec)
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
	}

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
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	encoder, ok := codecs[id]
	if !ok {
		c.String(http.StatusNotFound, fmt.Sprintf("\"%s\" does not exit on GoRunner", id))
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
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	decoder, ok := codecs[id]
	if !ok {
		c.String(http.StatusNotFound, fmt.Sprintf("\"%s\" does not exit on GoRunner", id))
		return
	}

	parsed, err := decoder.Decode(data)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.String(http.StatusOK, parsed)
	}
}

func codecHandler(c *gin.Context) {
	id := c.Param("id")
	data, err := c.GetRawData()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	codec, err := codec.New(data, id)

	if codecs[id] != nil {
		codecs[id].Stop()
	}
	codecs[id] = codec

	c.String(http.StatusOK, id)
}
