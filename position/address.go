package position

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/codingsince1985/geo-golang"
	"github.com/gin-gonic/gin"
)

// get geojson format for point
func FromGeo(c *gin.Context) {

	// read the body to json string
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	}

	var request GeoJSONPoint
	var geo_pos *geo.Address

	if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	}

	address, err := Geocoder().ReverseGeocode(request.Coordinates[1], request.Coordinates[0])
	if err != nil || address == nil {
		geo_pos = nil
	} else {
		geo_pos = address
	}

	c.JSON(http.StatusOK, geo_pos)
}
