package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("GoRunner by Parham Alvani")

	r := gin.Default()

	r.GET("/api/about", about)

	r.Run()
}

func about(c *gin.Context) {
	c.String(http.StatusOK, "18.20 is leaving us")
}

func decoder(w http.ResponseWriter, r *http.Request) {
}
