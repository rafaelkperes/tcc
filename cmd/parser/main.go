package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/rafaelkperes/tcc/pkg/measure"

	log "github.com/sirupsen/logrus"
)

// parser parses the producer/consumer log measures.
// The log is taken from stdin.

var (
	help = flag.Bool("h", false, "display this help")
	f    = flag.String("f", "", "input file (default to stdin)")
)

func main() {
	in := os.Stdin

	// flag-related init
	flag.Parse()
	if *help {
		flag.PrintDefaults()
	}
	if len(*f) > 0 {
		var err error
		in, err = os.Open(*f)
		if err != nil {
			log.Fatal(err)
		}
	}

	measures := make([]measure.Measure, 0)

	s := bufio.NewScanner(in)
	for s.Scan() {
		var o map[string]interface{}

		err := json.Unmarshal(s.Bytes(), &o)
		if err != nil {
			log.Errorf("on line '%s': %v", s.Text(), err)
			continue
		}

		var ok bool
		var mo map[string]interface{}
		if mo, ok = o["measures"].(map[string]interface{}); !ok {
			log.Warningf("not a measure: %s", s.Text())
			continue
		}

		m, err := measure.FromObject(mo)
		if err != nil {
			log.Errorf("on line '%s': %v", s.Text(), err)
			continue
		}

		log.Infof("parsed measure: %s", m)
		measures = append(measures, m)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(measures)
}
