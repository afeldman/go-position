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

	var request []geo_request

	if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	}

	type c_return struct {
		Index    int
		Position *GeoJSONPoint
	}
	c_geopoints := make(chan c_return, len(request))

	//create the geolocation for all requests
	for index, address := range request {
		go func(index int, raddress string) {
			location, err := Geocoder().Geocode(raddress)
			if err != nil || location == nil {
				c_geopoints <- c_return{Index: index, Position: nil}
			} else {
				c_geopoints <- c_return{Index: index, Position: NewGeoPoint(location.Lng, location.Lat)}
			}
		}(index, address.Address)
	}

	geoposition := make([]*GeoJSONPoint, len(request))

	for range geoposition {
		pos := <-c_geopoints
		geoposition[pos.Index] = pos.Position
	}

	c.JSON(http.StatusOK, geoposition)
}
