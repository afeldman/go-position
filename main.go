package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	position "github.com/afeldman/go-position/position"

	"github.com/gin-gonic/gin"
)

const VERSION string = "v0.2"

var (
	release bool
	port    int
	version bool
)

func init() {
	flag.BoolVar(&release, "release", false, "set to release mode")
	flag.BoolVar(&version, "version", false, "version")
	flag.IntVar(&port, "port", 8888, "server port")
}

func main() {

	flag.Parse()

	router := gin.New()

	if version {
		fmt.Println("Version is: ", VERSION)
		return
	}

	if release {
		log.Println("start release mode")
		gin.SetMode(gin.ReleaseMode)
	} else {
		log.Println("start development mode")
		gin.SetMode(gin.DebugMode)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "404 page not found",
			"error":   "",
			"status":  http.StatusNotFound})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"message": "405 method not allowed",
			"error":   "",
			"status":  http.StatusMethodNotAllowed})
	})

	v1 := router.Group("/v1")
	{
		v1.GET("/address/:address", position.FromAddress)

		v1.GET("/geo/:lat/:lon", position.FromGeo)
	}
	router.Run(fmt.Sprintf("0.0.0.0:%d", port))
}
