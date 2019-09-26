package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/rafaelkperes/tcc/pkg/data"
	"github.com/rafaelkperes/tcc/pkg/measure"
	log "github.com/sirupsen/logrus"
)

var (
	help          = flag.Bool("h", false, "display this help")
	lf            = flag.String("lf", "/var/log/std.producer.log", "standard logs file")
	mf            = flag.String("mf", "/var/log/msr.producer.log", "measure log file")
	endpoint      = flag.String("c", "http://localhost:9000", "set consumer endpoint")
	noOfReqs      = flag.Int("r", 12, "number of total requests")
	interval      = flag.Int("i", 5000, "interval in milliseconds between requests")
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

	f, err := os.OpenFile(*lf, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.WithFields(log.Fields{"event": "setupLogger", "error": err}).
			Error("failed to open log file")
	} else {
		log.SetOutput(io.MultiWriter(log.StandardLogger().Out, f))
	}

	// measure logger
	msrLogger := log.New()
	msrLogger.SetFormatter(&log.JSONFormatter{})
	msrLogger.SetOutput(os.Stderr)
	f, err = os.OpenFile(*mf, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.WithFields(log.Fields{"event": "setupMeasureLogger", "error": err}).
			Error("failed to open log file")
	} else {
		log.SetOutput(io.MultiWriter(msrLogger.Out, f))
	}

	// start requests
	client := &http.Client{}

	log.Printf("requesting a total of %d requests every %d milliseconds with a length of %d to %s", *noOfReqs, *interval, *payloadLength, *endpoint)
	for idx := 0; idx < *noOfReqs; idx++ {
		time.Sleep(time.Duration(*interval) * time.Millisecond)
		log.Printf("sending request %d/%d", idx+1, *noOfReqs)

		s := data.CreateStrings(*payloadLength, 100)

		b, err := json.Marshal(s)
		if err != nil {
			panic(err)
		}

		r, err := client.Post(*endpoint, "application/json", bytes.NewReader(b))
		if err != nil {
			log.Printf("got error on request: %v", err)
			continue
		}
		defer r.Body.Close()

		rb, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("got error while reading response body: %v", err)
			continue
		}

		m := make(map[string]interface{})
		err = json.Unmarshal(rb, &m)
		if err != nil {
			log.Printf("failed to unmarshal consumer measures: %v", err)
			continue
		}

		measures, err := measure.FromObject(m)
		if err != nil {
			log.Printf("failed to parse consumer measures: %v", err)
			continue
		}
		log.Printf("got consumer measures: %v", measures.AsObject())
	}

	log.Printf("done")
}

func displayHelp() {
	flag.PrintDefaults()
}
