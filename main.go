package main

import (
	"flag"
	"fmt"
	"net/http"

	position "github.com/afeldman/go-position/position"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"

	"github.com/gin-gonic/gin"
)

var (
	release bool
	port    int
	version bool
	server  string
)

func init() {
	flag.BoolVar(&release, "release", false, "set to release mode")
	flag.BoolVar(&version, "version", false, "version")
	flag.IntVar(&port, "port", 8888, "server port")
	flag.StringVar(&server, "server", "0.0.0.0", "server url")
}

// create a rest api for geoposition and address formates
func main() {

	flag.Parse()

	// start logger
	log := logrus.New()

	// set the router
	router := gin.New()
	router.Use(ginlogrus.Logger(log), gin.Recovery())

	// version informaiton
	if version {
		fmt.Println("Version is: ", position.VERSION)
		return
	}

	// set the router server in release mode if requested
	if release {
		log.Println("start release mode")
		gin.SetMode(gin.ReleaseMode)
	} else {
		log.Println("start development mode")
		gin.SetMode(gin.DebugMode)
	}

	// send 404 page not found for all wrong requests
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "404 page not found",
			"error":   ""})
	})

	// send not allowed
	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"message": "405 method not allowed",
			"error":   ""})
	})

	router.GET("/", func(c *gin.Context) {

		c.Header("Content-Type", "application/json; charset=utf-8")
		c.String(http.StatusOK, "{\"version\":\"%s\"}", position.VERSION)
	})

	//set the entry group for the api
	v1 := router.Group("/v1")
	{
		// get the geocodes of an address
		v1.POST("/address2geo", position.FromAddress)

		// get the addres of the geo positon
		v1.POST("/geo2address", position.FromGeo)
	}

	// start server
	router.Run(fmt.Sprintf("%s:%d", server, port))
}
