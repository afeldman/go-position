package position

import (
	"fmt"
	"strings"

	"github.com/codingsince1985/geo-golang"
	"github.com/codingsince1985/geo-golang/osm"
)

const (
	osmurl string = "https://nominatim.openstreetmap.org/"
)

type (
	baseURL         string
	geocodeResponse struct {
		DisplayName string `json:"display_name"`
		Lat         string
		Lon         string
		Error       string
		Addr        osm.Address `json:"address"`
	}
)

func Geocoder() geo.Geocoder {
	return geo.HTTPGeocoder{
		EndpointBuilder:       baseURL(osmurl),
		ResponseParserFactory: func() geo.ResponseParser { return &geocodeResponse{} },
	}
}

func (b baseURL) GeocodeURL(address string) string {
	return string(b) + "search?format=json&limit=1&q=" + address
}

func (b baseURL) ReverseGeocodeURL(location geo.Location) string {
	return string(b) + "reverse?" + fmt.Sprintf("format=json&lat=%f&lon=%f", location.Lat, location.Lng)
}

func (r *geocodeResponse) Location() (*geo.Location, error) {
	if r.Error != "" {
		return nil, fmt.Errorf("geocoding error: %s", r.Error)
	}
	if r.Lat == "" && r.Lon == "" {
		return nil, nil
	}

	return &geo.Location{
		Lat: geo.ParseFloat(r.Lat),
		Lng: geo.ParseFloat(r.Lon),
	}, nil
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
