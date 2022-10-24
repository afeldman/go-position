package position

import (
	"math"
	"net/http"
	"strconv"

	"github.com/codingsince1985/geo-golang"

	"github.com/gin-gonic/gin"
)

var (
	geocoder geo.Geocoder = Geocoder()
)

type response_location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type geo_response struct {
	Status  int               `json:"status"`
	Error   string            `json:"error"`
	Message response_location `json:"message"`
}

type address_response struct {
	Status  int         `json:"status"`
	Error   string      `json:"error"`
	Message geo.Address `json:"message"`
}

func FromGeo(c *gin.Context) {
	lat := c.Param("lat")
	lon := c.Param("lon")

	slat, _ := strconv.ParseFloat(lat, 64)
	slon, _ := strconv.ParseFloat(lon, 64)

	resp := address_response{
		Status: http.StatusNotFound,
		Error:  "",
		Message: geo.Address{
			FormattedAddress: "",
			HouseNumber:      "",
			Street:           "",
			Postcode:         "",
			City:             "",
			Suburb:           "",
			State:            "",
			Country:          "",
			CountryCode:      "",
		},
	}

	address, err := geocoder.ReverseGeocode(slat, slon)
	if err != nil {
		resp.Error = err.Error()
		c.JSON(http.StatusNotFound, resp)
	} else {
		resp.Message = *address
		resp.Status = http.StatusOK
		c.JSON(http.StatusOK, resp)
	}

}

func FromAddress(c *gin.Context) {
	address := c.Param("address")

	location, err := geocoder.Geocode(address)

	resp := geo_response{
		Status: http.StatusNotFound,
		Error:  "",
		Message: response_location{
			Lat: math.Inf(1),
			Lng: math.Inf(1),
		},
	}

	if err != nil {
		resp.Error = err.Error()
		c.JSON(http.StatusNotFound, resp)
	} else {
		resp.Status = http.StatusOK
		resp.Message.Lat = location.Lat
		resp.Message.Lng = location.Lng
		c.JSON(http.StatusOK, resp)
	}
}
