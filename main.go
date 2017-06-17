// darksky displays data from the Dark Sky API.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/weirdtales/darksky/pkg/darksky"
	"github.com/weirdtales/darksky/pkg/google"
)

var (
	imperial = flag.Bool("i", false, "use non-SI units")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	q := strings.Join(flag.Args(), " ")
	if q == "" {
		usage()
	}

	loc, err := geocode.Find(q)
	if err != nil {
		panic(err)
	}

	d, err := darksky.Get(loc, imperial)
	if err != nil {
		panic(err)
	}

	printDarksky(&d)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: darksky [options] query...\n")
	fmt.Fprintf(os.Stderr, "flags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func printDarksky(d *darksky.Darksky) {
	tz, _ := time.LoadLocation(d.Timezone)
	t := time.Unix(d.Currently.Time, 0)
	fmt.Printf("%s @ %s\n", d.GoogleName, t.In(tz).Format("2006-01-02 15:04:05"))
	fmt.Printf("Current temperature:\t%f\n", d.Currently.Temperature)
	fmt.Printf("Current apparent temp:\t%f\n", d.Currently.ApparentTemperature)
	fmt.Printf("Current summary:\t%s\n", d.Currently.Summary)
	fmt.Printf("Hourly summary:\t\t%s\n", d.Hourly.Summary)
	fmt.Printf("Daily summary:\t\t%s\n", d.Daily.Summary)
}
