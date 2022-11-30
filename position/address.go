package position

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/codingsince1985/geo-golang"
	"github.com/gin-gonic/gin"
)

func parseAddressArray(c *gin.Context, jsonData []byte) {
	var request []GeoJSONPoint

	if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	}

	type c_return struct {
		Index    int
		Position *geo.Address
	}
	c_address := make(chan c_return, len(request))

	//create the geolocation for all requests
	for index, geos := range request {
		go func(index int, lng, lat float64) {
			address, err := Geocoder().ReverseGeocode(lng, lat)
			if err != nil || address == nil {
				c_address <- c_return{Index: index, Position: nil}
			} else {
				c_address <- c_return{Index: index, Position: address}
			}
		}(index, geos.Coordinates[1], geos.Coordinates[0])
	}

	address_responce := make([]*geo.Address, len(request))

	for range address_responce {
		pos := <-c_address
		address_responce[pos.Index] = pos.Position
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
