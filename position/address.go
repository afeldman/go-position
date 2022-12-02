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
			if lng == 0 && lat == 0 {
				c_address <- c_return{Index: index, Position: nil}
			} else {
				address, err := Geocoder().ReverseGeocode(lng, lat)
				if err != nil || address == nil {
					c_address <- c_return{Index: index, Position: nil}
				} else {
					c_address <- c_return{Index: index, Position: address}
				}
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
