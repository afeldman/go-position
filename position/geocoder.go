package position

import (
	"github.com/codingsince1985/geo-golang"
	"github.com/codingsince1985/geo-golang/openstreetmap"
)

// build geopoints for easy calculation
// first longitude and second latetude
type GeoJSONPoint struct {
	Type        string     `json:"type"`
	Coordinates [2]float64 `json:"coordinates"`
}

// return a geocoder.
// default is openstreetmap but we can change in future
func Geocoder() geo.Geocoder {
	return openstreetmap.Geocoder()
}
