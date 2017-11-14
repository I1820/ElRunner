package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aiotrc/GoRunner/decoder"
	"github.com/gin-gonic/gin"
)

var decoders map[string]decoder.Decoder

func main() {
	fmt.Println("GoRunner by Parham Alvani")

	decoders = make(map[string]decoder.Decoder)

	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/decode/:id", decodeHandler)
		api.GET("/about", aboutHandler)
		api.POST("/decoder/:id", decoderHandler)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		fmt.Printf("GoRunner Listen: %s\n", srv.Addr)
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("Listen Error:", err)
		}
	}()

	decoders["isrc-sensor"], _ = decoder.New([]byte(`
import cbor
import base64

s = input()
d = cbor.loads(s.base64.b64decode(s))
print(d)
	`), "isrc-sensor")

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

func decodeHandler(c *gin.Context) {
	id := c.Param("id")
	data, err := c.GetRawData()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	decoder, ok := decoders[id]
	if !ok {
		c.String(http.StatusNotFound, fmt.Sprintf("\"%s\" does not exit on GoRunner", id))
		return
	}

	c.String(http.StatusOK, decoder.Decode(data))
}

func decoderHandler(c *gin.Context) {
	id := c.Param("id")
	data, err := c.GetRawData()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	decoder, err := decoder.New(data, id)

	if decoders[id] != nil {
		decoders[id].Stop()
	}
	decoders[id] = decoder

	c.String(http.StatusOK, id)
}
