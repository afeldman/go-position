package position

import (
	"encoding/json"
	"fmt"
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

func FromAddress(c *gin.Context) {

	// build default responce
	resp := NewGeo()

	// read the body to json string
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		resp.Error = err.Error()
		c.JSON(http.StatusNotFound, resp)
	}

	// build request type
	request := geo_request{
		Address: "",
	}

	// json string to request object
	if err = json.Unmarshal([]byte(jsonData), &request); err != nil {
		resp.Error = err.Error()
		c.JSON(http.StatusNotFound, resp)
	}

	fmt.Println(request.Address)

	//create the geolocation
	location, err := Geocoder().Geocode(request.Address)
	if err != nil {
		resp.Error = err.Error()
		c.JSON(http.StatusNotFound, resp)
	} else {
		resp.Geo = NewGeoPoint(location.Lng, location.Lat)
		c.JSON(http.StatusOK, resp)
	}
}
