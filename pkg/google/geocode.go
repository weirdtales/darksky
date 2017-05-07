// geocode uses google's public maps API to convert a human-writable string
// into a formatted address and a pair of coordinates.
package geocode

import (
	"fmt"
	"time"

	"encoding/json"
	"net/http"
)

type Location struct {
	Results []struct {
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

// GetLocation returns a struct of location details for the supplied query string.
func Find(q string) (loc Location, err error) {
	url := "http://maps.googleapis.com/maps/api/geocode/json?address=%s&sensor=false"
	url = fmt.Sprintf(url, q)

	err = getJson(url, &loc)
	if err != nil {
		return
	}
	if len(loc.Results) < 1 {
		err = fmt.Errorf("no results for \"%s\"", q)
	}

	return
}

func getJson(url string, i interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		err = fmt.Errorf("%+v", r.Body)
		return err
	}

	return json.NewDecoder(r.Body).Decode(i)
}
