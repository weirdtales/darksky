package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/dedelala/round"
	"github.com/dedelala/sysexits"
	"github.com/weirdtales/darksky/pkg/darksky"
	"github.com/weirdtales/darksky/pkg/gmap"
)

func main() {
	units := flag.String("u", "si", "units: auto, ca, uk2, us, si")
	flag.Parse()

	token := os.Getenv("DARKSKY_TOKEN")
	if token == "" {
		fmt.Fprintf(os.Stderr, "DARKSKY_TOKEN environ var missing or empty")
		os.Exit(sysexits.Usage)
	}

	q := strings.Join(flag.Args(), " ")
	if q == "" {
		flag.Usage()
		os.Exit(sysexits.Usage)
	}

	loc, err := gmap.Find(q)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(sysexits.Unavailable)
	}
	spout := round.NewSpinMe(os.Stdout, round.Block)
	fmt.Fprintf(&spout, "%s ", loc)

	res, err := darksky.Forecast(token, *loc, *units)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(sysexits.Unavailable)
	}
	spout.Close()
	fmt.Printf("\n%s\n", strings.Repeat("-", utf8.RuneCountInString(loc.String())))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(sysexits.Unavailable)
	}
	fmt.Println(res)
	os.Exit(sysexits.OK)
}
