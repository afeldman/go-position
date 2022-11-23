package position

import (
	"fmt"
	"strings"

	"github.com/codingsince1985/geo-golang"
	"github.com/codingsince1985/geo-golang/openstreetmap"
	"github.com/codingsince1985/geo-golang/osm"
)

type (
	geocodeResponse struct {
		DisplayName string      `json:"display_name"`
		Lat         string      `json:"lat"`
		Lon         string      `json:"lon"`
		Error       string      `json:"err"`
		Addr        osm.Address `json:"address"`
	}
)

// build geopoints for easy calculation
// first longitude and second latetude
type GeoJSONPoint struct {
	Type        string     `json:"type"`
	Coordinates [2]float64 `json:"coordinates"`
}

func Geocoder() geo.Geocoder {
	return openstreetmap.Geocoder()
}

func (r *geocodeResponse) Address() (*geo.Address, error) {
	if r.Error != "" {
		return nil, fmt.Errorf("reverse geocoding error: %s", r.Error)
	}

	return &geo.Address{
		FormattedAddress: r.DisplayName,
		HouseNumber:      r.Addr.HouseNumber,
		Street:           r.Addr.Street(),
		Postcode:         r.Addr.Postcode,
		City:             r.Addr.Locality(),
		Suburb:           r.Addr.Suburb,
		State:            r.Addr.State,
		Country:          r.Addr.Country,
		CountryCode:      strings.ToUpper(r.Addr.CountryCode),
	}, nil
}
