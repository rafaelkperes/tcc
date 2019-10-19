package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rafaelkperes/tcc/pkg/measure"

	log "github.com/sirupsen/logrus"
)

// parser parses the producer/consumer log measures.

const (
	timeFormat = time.RFC3339Nano
)

var (
	help = flag.Bool("h", false, "display this help")
	_f   = flag.String("f", "", "use file as input (instead of stdin)")
	_o   = flag.String("o", "", "output to file (instead of stdout)")
)

type entry struct {
	msr    measure.Measure
	typ    string
	format string
}

func main() {
	// flag-related init
	flag.Parse()
	if *help {
		flag.PrintDefaults()
	}

	in := os.Stdin
	if len(*_f) > 0 {
		var err error
		in, err = os.Open(*_f)
		if err != nil {
			log.Fatal(err)
		}
		defer in.Close()
	}

	out := os.Stdout
	if len(*_o) > 0 {
		var err error
		out, err = os.OpenFile(*_o, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
	}

	// parse entries
	entries := make([]entry, 0)
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
			log.Debugf("skipping entry: not a measure: %s", s.Text())
			continue
		}

		format, ok := o["format"].(string)
		if !ok {
			log.Warningf("not a valid format: %v", format)
		}
		typ, ok := o["typ"].(string)
		if !ok {
			log.Warningf("not a valid type: %v", typ)
		}

		m, err := measure.FromObject(mo)
		if err != nil {
			log.Errorf("on line '%s': %v", s.Text(), err)
			continue
		}

		log.Debugf("parsed measure: %s", m)
		entries = append(entries, entry{
			msr:    m,
			typ:    typ,
			format: format,
		})
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	// produce output
	w := csv.NewWriter(out)
	defer w.Flush()

	// headers
	err := w.Write([]string{
		"format",
		"type",
		"dTotal",
		"dSrlz",
		"dSent",
		"dConsRed",
		"dBodyRead",
		"dDsrlz",
		"tStrt",
		"tSrld",
		"tSent",
		"tTend",
		"tRecv",
		"tRbod",
		"tDsrl",
	})
	if err != nil {
		log.Fatal(err)
	}

	// data
	for _, it := range entries {
		rec := []string{
			string(it.format),
			string(it.typ),
		}

		strt, _ := it.msr.Get("strt")
		srld, _ := it.msr.Get("srld")
		sent, _ := it.msr.Get("sent")
		tend, _ := it.msr.Get("tend")

		recv, _ := it.msr.Get("recv")
		rbod, _ := it.msr.Get("rbod")
		dsrl, _ := it.msr.Get("dsrl")

		dTotal := tend.Sub(strt)

		dSrlz := srld.Sub(strt)
		dSent := sent.Sub(srld)
		dConsReq := tend.Sub(sent)

		dBodyRead := rbod.Sub(recv)
		dDsrlz := dsrl.Sub(rbod)

		rec = append(rec,
			fmt.Sprint(dTotal.Nanoseconds()),
			fmt.Sprint(dSrlz.Nanoseconds()),
			fmt.Sprint(dSent.Nanoseconds()),
			fmt.Sprint(dConsReq.Nanoseconds()),
			fmt.Sprint(dBodyRead.Nanoseconds()),
			fmt.Sprint(dDsrlz.Nanoseconds()),
			strt.Format(timeFormat),
			srld.Format(timeFormat),
			sent.Format(timeFormat),
			tend.Format(timeFormat),
			recv.Format(timeFormat),
			rbod.Format(timeFormat),
			dsrl.Format(timeFormat),
		)

		err := w.Write(rec)
		if err != nil {
			log.Fatal(err)
		}
	}
}
