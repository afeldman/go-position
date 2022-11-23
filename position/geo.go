package position

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// error or data
type geo_response struct {
	Error string        `json:"error"`
	Geo   *GeoJSONPoint `json:"geo"`
}

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

// create a new geo object
func NewGeo() *geo_response {
	return &geo_response{
		Error: "",
		Geo:   nil,
	}
}

func parseGeoArray(c *gin.Context, jsonData []byte) {
	var request []geo_request

	if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
		var errgeopoints = make([]geo_response, 1)
		errgeopoints[0].Error = err.Error()
		c.JSON(http.StatusNotFound, errgeopoints)
	}

	var geopoints = make([]geo_response, len(request))

	//create the geolocation for all requests
	for index, address := range request {
		location, err := Geocoder().Geocode(address.Address)
		if err != nil {
			geopoints[index].Error = err.Error()
		} else {
			geopoints[index].Geo = NewGeoPoint(location.Lng, location.Lat)
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
		var errgeopoints = make([]geo_response, 1)
		errgeopoints[0].Error = err.Error()
		c.JSON(http.StatusNotFound, errgeopoints)
	}

	//create the geolocation
	location, err := Geocoder().Geocode(request.Address)
	var geopoints = make([]geo_response, 1)
	if err != nil {
		geopoints[0].Error = err.Error()
		c.JSON(http.StatusNotFound, geopoints)
	} else {
		geopoints[0].Geo = NewGeoPoint(location.Lng, location.Lat)
		c.JSON(http.StatusOK, geopoints)
	}
}

func FromAddress(c *gin.Context) {

	// build default responce
	resp := NewGeo()

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
		parseGeoArray(c, jsonData)
	case 1:
		parseGeoObject(c, jsonData)
	default:
		resp.Error = "must be object or array in the request"
		c.JSON(http.StatusInternalServerError, resp)
	}

}
