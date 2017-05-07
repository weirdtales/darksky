// darksky displays data from the Dark Sky API.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/weirdtales/darksky/pkg/darksky"
	"github.com/weirdtales/darksky/pkg/google"
)

var (
	imperial = flag.Bool("i", false, "use non-SI units")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	q := flag.Arg(0)
	if q == "" {
		usage()
	}
	fmt.Println(q)

	loc, err := geocode.Find(q)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", loc)

	d, err := darksky.Get(loc, imperial)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", d)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: darksky [options] query...\n")
	fmt.Fprintf(os.Stderr, "flags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}
