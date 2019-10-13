package main

import (
	"bytes"
	"flag"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rafaelkperes/tcc/internal/svc/prod"

	"github.com/rafaelkperes/tcc/pkg/data"
	"github.com/rafaelkperes/tcc/pkg/file"
	log "github.com/sirupsen/logrus"
)

const (
	defaultConsumerEndpoint = "http://localhost:9000"
)

var (
	defaultLogOutput = os.Stderr
	defaultPort      = "9001"
)

var (
	help = flag.Bool("h", false, "display this help")
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
	log.SetOutput(defaultLogOutput)
	log.SetLevel(log.DebugLevel)

	dir, ok := os.LookupEnv("RESULTS_DIR")
	if !ok {
		dir = "/tmp/tcc"
		log.Warning("RESULTS_DIR not set")
	}
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("failed to create directory %s: %v", dir, err)
	}

	ep, ok := os.LookupEnv("CONSUMER_ENDPOINT")
	if !ok {
		log.Warning("CONSUMER_ENDPOINT not set")
		ep = defaultConsumerEndpoint
	}
	go runAll(ep, dir)

	// Set up file server
	http.Handle("/", http.FileServer(http.Dir(dir)))
	http.HandleFunc("/results.tar.gz", func(w http.ResponseWriter, r *http.Request) {
		var b []byte
		buff := bytes.NewBuffer(b)

		if err := file.AsTarball(dir, buff); err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(w, buff); err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Warning("CONSUMER_ENDPOINT not set")
		port = defaultPort
	}
	log.Printf("Serving %s on HTTP port: %s\n", dir, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func runAll(consumerEndpoint, dir string) {
	formats := []data.Format{data.FormatJSON, data.FormatProtobuf, data.FormatMsgpack, data.FormatAvro}
	types := []data.Type{data.TypeInt, data.TypeFloat, data.TypeString, data.TypeObject}

	for idx, f := range formats {
		for jdx, t := range types {
			log := log.WithFields(map[string]interface{}{"format": f, "typ": t})

			if err := setLoggingFile(dir, f, t); err != nil {
				log.Error(err)
			}

			log.WithField("event", "progress").
				Debugf("running with %d/%d settings", (idx)*len(types)+jdx+1, len(formats)*len(types))
			run(consumerEndpoint, f, t)
		}
	}
	log.Debug("finished")
}

func setLoggingFile(dir string, format data.Format, typ data.Type) error {
	fns := strings.Split(string(format), "/")
	fn := fns[len(fns)-1]

	dir = filepath.Join(dir, fn)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "failed to create directory %s", dir)
	}

	filename := filepath.Join(dir, string(typ)+".log")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "failed to open file %s", filename)
	}

	log.SetOutput(io.MultiWriter(defaultLogOutput, f))
	return nil
}

func run(consumerEndpoint string, format data.Format, typ data.Type) {
	const (
		noOfReqs int           = 1
		interval time.Duration = 0
		total    int64         = 1
	)

	d, err := data.Create(typ, total)
	if err != nil {
		log.Fatal(err.Error())
	}

	p := prod.NewProducer(consumerEndpoint)
	p.Produce(d, format, noOfReqs, interval)
	log.Debug("done")
}

func displayHelp() {
	flag.PrintDefaults()
}
