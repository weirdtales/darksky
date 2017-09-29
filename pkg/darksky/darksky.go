// Package darksky lightly abstracts the Darksky API.
package darksky

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/weirdtales/darksky/pkg/gmap"
)

type dataPoint struct {
	Time       int64   `json:"time"`
	Summary    string  `json:"summary"`
	Icon       string  `json:"icon"`
	Temp       float64 `json:"temperature"`
	ATemp      float64 `json:"apparentTemperature"`
	DewPoint   float64 `json:"dewPoint"`
	Humidity   float64 `json:"humidity"`
	WindSpeed  float64 `json:"windSpeed"`
	Visibility float64 `json:"visibility"`
	CloudCover float64 `json:"cloudCover"`
	Pressure   float64 `json:"pressure"`
	Ozone      float64 `json:"ozone"`
}

// Result is a single result from the API.
type Result struct {
	TZ       string    `json:"timezone"`
	Current  dataPoint `json:"currently"`
	Minutely struct {
		Summary string      `json:"summary"`
		Icon    string      `json:"icon"`
		Data    []dataPoint `json:"data"`
	} `json:"minutely"`
	Daily struct {
		Summary string      `json:"summary"`
		Icon    string      `json:"icon"`
		Data    []dataPoint `json:"data"`
	} `json:"daily"`
	Hourly struct {
		Summary string      `json:"summary"`
		Icon    string      `json:"icon"`
		Data    []dataPoint `json:"data"`
	} `json:"hourly"`
	Loc gmap.Location
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

var apiURL = "https://api.darksky.net/forecast/%s/%f,%f?units=%s"

var blocks = []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇"}

// Forecast ...
func Forecast(token string, loc gmap.Location, units string) (*Result, error) {
	u := fmt.Sprintf(apiURL, token, loc.Lat, loc.Lng, units)
	r := &Result{Loc: loc}
	err := get(u, r)
	if err != nil {
		return nil, fmt.Errorf("unable to read forecast data: %v", err)
	}
	if r.TZ == "" { // TODO is this the correct way to check?
		return nil, fmt.Errorf("no data for location: %v", loc)
	}
	return r, nil
}

// StdPrint ...
func (r Result) String() string {
	var b bytes.Buffer
	bw := bufio.NewWriter(&b)
	w := tabwriter.NewWriter(bw, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Timezone\t%s\n", r.TZ)
	fmt.Fprintf(w, "Currently\t%s\n", r.Current.Summary)
	fmt.Fprintf(w, "Temp (AT)\t%.1f (%.1f)\n", r.Current.Temp, r.Current.ATemp)

	fmt.Fprintf(w, "Hourly\t%s\n", r.Hourly.Summary)
	fmt.Fprintf(w, "\t%s\n", getBar(r.Hourly.Data))
	fmt.Fprintf(w, "Daily\t%s\n", r.Daily.Summary)
	fmt.Fprintf(w, "\t%s\n", getBar(r.Daily.Data))
	w.Flush()
	bw.Flush()
	return b.String()
}

func getBar(d []dataPoint) string {
	o := ""
	vals := []float64{}
	for _, p := range d {
		vals = append(vals, p.Temp)
	}
	sort.Float64s(vals)
	for _, p := range d {
		o += fmt.Sprint(getBlock(p.Temp, vals))
	}
	return o
}

func getMinMax(d []float64) (float64, float64) {
	min := 1000.0
	max := 0.0
	for _, v := range d {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	return min, max
}

func getBlock(v float64, d []float64) string {
	i := 0
	n := len(d)
	nb := len(blocks)
	switch {
	case n < 2:
		i = 0
	case n == nb:
		i = 1
	case n < nb:
		i = 7
	case n > nb:
		sd := sort.Float64Slice(d)
		i = sd.Search(v) / (n / nb)
	}
	return blocks[i]
}

// get ...
func get(u string, t interface{}) error {
	r, err := httpClient.Get(u)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return fmt.Errorf("non-200 status code: %d", r.StatusCode)
	}
	return json.NewDecoder(r.Body).Decode(t)
}
