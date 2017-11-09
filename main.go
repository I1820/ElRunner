package main

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"

	"github.com/aiotrc/GoRunner/runner"
	"github.com/gin-gonic/gin"
)

var runners []runner.Runner

func main() {
	fmt.Println("GoRunner by Parham Alvani")
	runners = make([]runner.Runner, 1, 100)

	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/decode/:data", decode)
		api.GET("/about", about)
	}

	runners[0] = runner.New(&runner.Task{
		Run: func(e runner.Event) runner.Output {
			cmd := exec.Command("python3", "./runner/hello.py")

			stdin, err := cmd.StdinPipe()
			if err != nil {
			}
			io.WriteString(stdin, e.Data())
			stdin.Close()

			out, err := cmd.Output()
			if err != nil {
			}

			return runner.Output(out)
		},
		Interval: 0,
	}, 1)
	go runners[0].Start()

	r.Run()
}

func about(c *gin.Context) {
	c.String(http.StatusOK, "18.20 is leaving us")
}

func decode(c *gin.Context) {
	data := c.Param("data")
	runners[0].Trigger(data)
	c.String(http.StatusOK, string(runners[0].Output()))
}
