package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	position "github.com/afeldman/go-position/position"

	"github.com/codingsince1985/geo-golang"
	"github.com/gin-gonic/gin"
)

var geocoder geo.Geocoder = position.Geocoder()

const VERSION string = "v0.1.0"

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
		c.JSON(http.StatusNotFound, gin.H{"message": "404 page not found", "status": http.StatusNotFound})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "405 method not allowed", "status": http.StatusMethodNotAllowed})
	})

	v1 := router.Group("/v1")
	{
		v1.GET("/address/:address", func(c *gin.Context) {
			adress := c.Param("address")

			location, err := geocoder.Geocode(adress)

			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"status":  http.StatusNotFound,
					"message": err.Error,
				})
			} else {
				var b []byte
				b, err = json.Marshal(location)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  http.StatusBadRequest,
						"message": err.Error(),
					})
				} else {
					c.JSON(http.StatusOK, gin.H{
						"status":  http.StatusOK,
						"message": string(b),
					})
				}
			}

		})

		v1.GET("/geo/:lat/:lon", func(c *gin.Context) {
			lat := c.Param("lat")
			lon := c.Param("lon")

			slat, _ := strconv.ParseFloat(lat, 64)
			slon, _ := strconv.ParseFloat(lon, 64)

			adress, err := geocoder.ReverseGeocode(slat, slon)

			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"status":  http.StatusNotFound,
					"message": err.Error,
				})
			} else {
				var b []byte
				b, err = json.Marshal(adress)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  http.StatusBadRequest,
						"message": err.Error(),
					})
				} else {
					c.JSON(http.StatusOK, gin.H{
						"status":  http.StatusOK,
						"message": string(b),
					})
				}
			}

		})
	}
	router.Run(fmt.Sprintf("0.0.0.0:%d", port))
}
