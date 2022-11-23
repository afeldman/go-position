package position

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/codingsince1985/geo-golang"
	"github.com/gin-gonic/gin"
)

func parseAddressArray(c *gin.Context, jsonData []byte) {
	var request []GeoJSONPoint

	if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	}

	var address_responce = make([]geo.Address, len(request))

	//create the geolocation for all requests
	for index, geos := range request {
		address, err := Geocoder().ReverseGeocode(geos.Coordinates[1], geos.Coordinates[0])
		if err != nil {
			log.Println(err.Error())
		} else {
			address_responce[index] = *address
		}
	}
	c.JSON(http.StatusOK, address_responce)
}

func parseAddressObject(c *gin.Context, jsonData []byte) {
	// build request type
	request := GeoJSONPoint{}

	// json string to request object
	if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	}
	// make array responce
	var addresses = make([]geo.Address, 1)

	//create the geolocation
	address, err := Geocoder().ReverseGeocode(request.Coordinates[1], request.Coordinates[0])
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	} else {
		addresses[0] = *address
		c.JSON(http.StatusOK, addresses)
	}
}

// get geojson format for point
func FromGeo(c *gin.Context) {

	// read the body to json string
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	}

	t, err := jsonType(bytes.NewReader(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	switch t {
	case 0:
		parseAddressArray(c, jsonData)
	case 1:
		parseAddressObject(c, jsonData)
	default:
		c.JSON(http.StatusInternalServerError, "must be object or array in the request")
	}
}
