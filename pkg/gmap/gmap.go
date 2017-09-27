// Package gmap is a real basic interface into the Google maps API.
package gmap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Location is a location that Google Maps knows about.
type Location struct {
	Addr string
	Lat  float64
	Lng  float64
}

// location is the raw JSON response from Google.
type location struct {
	Results []struct {
		Addr string `json:"formatted_address"`
		Geom struct {
			Loc struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

const apiURL = "http://maps.google.com/maps/api/geocode/json?address=%s&sensor=no"

// Find asks Google Maps to convert the supplied string into a Location struct.
// An error is returned if marshaling fails or no results are returned. If more
// than one result is returned, the first one is returned.
func Find(q string) (*Location, error) {
	q = url.QueryEscape(q)
	l := &location{}
	err := get(q, l)
	if err != nil {
		return nil, err
	}
	if len(l.Results) < 1 {
		return nil, fmt.Errorf("no results from maps API for: %v", q)
	}
	loc := &Location{}
	loc.Addr = l.Results[0].Addr
	loc.Lat = l.Results[0].Geom.Loc.Lat
	loc.Lng = l.Results[0].Geom.Loc.Lng
	return loc, nil
}

// Location stringer.
func (l Location) String() string {
	return l.Addr
}

func get(q string, t interface{}) error {
	u := fmt.Sprintf(apiURL, q)
	r, err := httpClient.Get(u)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return fmt.Errorf("non-200 from google maps: %d", r.StatusCode)
	}
	return json.NewDecoder(r.Body).Decode(t)
}
