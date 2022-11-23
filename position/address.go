package position

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/codingsince1985/geo-golang"
	"github.com/gin-gonic/gin"
)

// address or error
type address_response struct {
	Error   string       `json:"error"`
	Message *geo.Address `json:"message"`
}

func parseAddressArray(c *gin.Context, jsonData []byte) {
	var request []GeoJSONPoint

	if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
		var erraddr = make([]address_response, 1)
		erraddr[0].Error = err.Error()
		c.JSON(http.StatusNotFound, erraddr)
	}

	var address_responce = make([]address_response, len(request))

	//create the geolocation for all requests
	for index, geos := range request {
		address, err := Geocoder().ReverseGeocode(geos.Coordinates[1], geos.Coordinates[0])
		if err != nil {
			address_responce[index].Error = err.Error()
		} else {
			address_responce[index].Message = address
		}
	}
	c.JSON(http.StatusOK, address_responce)
}

func parseAddressObject(c *gin.Context, jsonData []byte) {
	// build request type
	request := GeoJSONPoint{}

	// json string to request object
	if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
		var erraddr = make([]address_response, 1)
		erraddr[0].Error = err.Error()
		c.JSON(http.StatusNotFound, erraddr)
	}
	// make array responce
	var addresses = make([]address_response, 1)

	//create the geolocation
	address, err := Geocoder().ReverseGeocode(request.Coordinates[1], request.Coordinates[0])
	if err != nil {
		addresses[0].Error = err.Error()
		c.JSON(http.StatusNotFound, addresses)
	} else {
		addresses[0].Message = address
		c.JSON(http.StatusOK, addresses)
	}
}

// get geojson format for point
func FromGeo(c *gin.Context) {

	resp := address_response{
		Error:   "",
		Message: nil,
	}

	// read the body to json string
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		resp.Error = err.Error()
		c.JSON(http.StatusNotFound, resp)
	}

	t, err := jsonType(bytes.NewReader(jsonData))
	if err != nil {
		resp.Error = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
	}

	switch t {
	case 0:
		parseAddressArray(c, jsonData)
	case 1:
		parseAddressObject(c, jsonData)
	default:
		resp.Error = "must be object or array in the request"
		c.JSON(http.StatusInternalServerError, resp)
	}
}
