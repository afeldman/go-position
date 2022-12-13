package position

import (
	"encoding/json"
	"io/ioutil"
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

func FromAddress(c *gin.Context) {
	// read the body to json string
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	}

	var request geo_request
	var geoposition *GeoJSONPoint

	if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	}

	location, err := Geocoder().Geocode(request.Address)
	if err != nil || location == nil {
		geoposition = nil
	} else {
		geoposition = NewGeoPoint(location.Lng, location.Lat)
	}

	c.JSON(http.StatusOK, geoposition)
}
