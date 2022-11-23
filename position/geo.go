package position

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// geo request data is a address
type geo_request struct {
	Address string `json:"address"`
}

// build new geopoint
func NewGeoPoint(lon, lat float64) *GeoJSONPoint {
	return &GeoJSONPoint{
		Type:        "Point",
		Coordinates: [2]float64{lon, lat},
	}
}

func parseGeoArray(c *gin.Context, jsonData []byte) {
	var request []geo_request

	if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	}

	var geopoints = make([]GeoJSONPoint, len(request))

	//create the geolocation for all requests
	for index, address := range request {
		location, err := Geocoder().Geocode(address.Address)
		if err != nil {
			log.Println(err.Error())
		} else {
			geopoints[index] = *NewGeoPoint(location.Lng, location.Lat)
		}
	}

	c.JSON(http.StatusOK, geopoints)
}

func parseGeoObject(c *gin.Context, jsonData []byte) {
	// build request type
	request := geo_request{
		Address: "",
	}

	// json string to request object
	if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	}

	//create the geolocation
	location, err := Geocoder().Geocode(request.Address)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	} else {
		var geopoints = make([]GeoJSONPoint, 1)
		geopoints[0] = *NewGeoPoint(location.Lng, location.Lat)
		c.JSON(http.StatusOK, geopoints)
	}
}

func FromAddress(c *gin.Context) {

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
		parseGeoArray(c, jsonData)
	case 1:
		parseGeoObject(c, jsonData)
	default:
		c.JSON(http.StatusBadRequest, "must be object or array in the request")
	}

}
