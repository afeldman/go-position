package position

import (
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

	// build request type
	request := GeoJSONPoint{}

	// json string to request object
	if err = json.Unmarshal([]byte(jsonData), &request); err != nil {
		resp.Error = err.Error()
		c.JSON(http.StatusNotFound, resp)
	}

	address, err := Geocoder().ReverseGeocode(request.Coordinates[1], request.Coordinates[0])
	if err != nil {
		resp.Error = err.Error()
		c.JSON(http.StatusNotFound, resp)
	} else {
		resp.Message = address
		c.JSON(http.StatusOK, resp)
	}

}
