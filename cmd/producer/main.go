package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rafaelkperes/tcc/internal/svc/cons"

	"github.com/rafaelkperes/tcc/pkg/data"
)

var help = flag.Bool("h", false, "display this help")
var endpoint = flag.String("c", "http://localhost:9000", "set consumer endpoint")
var noOfReqs = flag.Int("r", 12, "number of total requests")
var interval = flag.Int("i", 5000, "interval in milliseconds between requests")
var payloadLength = flag.Int64("l", 1e6, "size of the array for the payload")

func main() {
	flag.Parse()

	if *help {
		displayHelp()
		os.Exit(0)
	}

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

		_, err = cons.MeasuresFromJSON(rb)
		if err != nil {
			log.Printf("failed to parse consumer measures: %v", err)
			continue
		}
		log.Printf("got consumer measures: %s", string(rb))
	}

	log.Printf("done")
}

func displayHelp() {
	flag.PrintDefaults()
}
