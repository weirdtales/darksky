package darksky

import (
	"fmt"
	"os"
	"time"

	"encoding/json"
	"net/http"

	"github.com/weirdtales/darksky/pkg/google"
)

// Darksky maps to an API response
type Darksky struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Currently struct {
		Time                int64   `json:"time"`
		Summary             string  `json:"summary"`
		Icon                string  `json:"icon"`
		Temperature         float64 `json:"temperature"`
		ApparentTemperature float64 `json:"apparentTemperature"`
		DewPoint            float64 `json:"dewPoint"`
		Humidity            float64 `json:"humidity"`
		WindSpeed           float64 `json:"windSpeed"`
		Visibility          float64 `json:"visibility"`
		CloudCover          float64 `json:"cloudCover"`
		Pressure            float64 `json:"pressure"`
		Ozone               float64 `json:"ozone"`
	} `json:"currently"`
	Minutely struct {
		Summary string `json:"summary"`
		Icon    string `json:"icon"`
	} `json:"minutely"`
	Daily struct {
		Summary string `json:"summary"`
		Icon    string `json:"icon"`
	} `json:"daily"`
	Hourly struct {
		Summary string `json:"summary"`
		Icon    string `json:"icon"`
	} `json:"hourly"`

	GoogleName string // use the google maps FormattedAddress value - it's good
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

// Get takes a location struct and returns a darksky one
func Get(loc geocode.Location, imperial *bool) (d Darksky, err error) {
	token := os.Getenv("DARKSKY_TOKEN")
	if token == "" {
		err = fmt.Errorf("DARKSKY_TOKEN not set in environment")
		return
	}

	d.GoogleName = loc.Results[0].FormattedAddress
	units := "si"
	if *imperial {
		units = "us"
	}
	url := "https://api.darksky.net/forecast/%s/%f,%f?units=%s"
	url = fmt.Sprintf(url, token, loc.Results[0].Geometry.Location.Lat, loc.Results[0].Geometry.Location.Lng, units)
	//fmt.Println(url)

	err = getJSON(url, &d)
	if err != nil {
		return
	}

	if d.Timezone == "" {
		fmt.Printf("%+v\n", d)
		err = fmt.Errorf("No data for that location")
	}
	return
}

func getJSON(url string, i interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return err
	}

	return json.NewDecoder(r.Body).Decode(i)
}
