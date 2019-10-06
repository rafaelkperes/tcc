package main

import (
	"flag"
	"math"
	"os"
	"time"

	"github.com/rafaelkperes/tcc/internal/svc/prod"

	"github.com/rafaelkperes/tcc/pkg/data"
	log "github.com/sirupsen/logrus"
)

var (
	help          = flag.Bool("h", false, "display this help")
	format        = flag.String("f", string(data.FormatJSON), "format")
	typ           = flag.String("t", string(data.TypeString), "data type")
	endpoint      = flag.String("c", "http://localhost:9000", "set consumer endpoint")
	noOfReqs      = flag.Int("r", 12, "number of total requests")
	interval      = flag.Int("i", 0, "interval in milliseconds between concurrent requests; if 0, requests are done sequentially")
	payloadLength = flag.Int64("l", 1e6, "size of the array for the payload")
	strLength     = flag.Int64("strlen", 100, "length of random strings")
	intMin        = flag.Int64("intmin", 0, "minimun value for random integers")
	intMax        = flag.Int64("intmax", math.MaxInt64, "maximum value for random integers")
)

func main() {
	// flag-related init
	flag.Parse()
	if *help {
		displayHelp()
		os.Exit(0)
	}

	// setup logging
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)

	d, err := data.Create(data.Type(*typ), *payloadLength, *intMin, *intMax, *strLength)
	if err != nil {
		log.Fatal(err.Error())
	}

	p := prod.NewProducer(*endpoint)
	p.Produce(d, data.Format(*format), *noOfReqs, time.Duration(*interval)*time.Millisecond)

	log.Debug("done")
}

func displayHelp() {
	flag.PrintDefaults()
}
