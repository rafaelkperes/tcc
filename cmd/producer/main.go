package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"

	"github.com/pkg/errors"
	"github.com/rafaelkperes/tcc/internal/svc/prod"

	"github.com/rafaelkperes/tcc/pkg/data"
	log "github.com/sirupsen/logrus"
)

const (
	defaultConsumerEndpoint = "http://localhost:9000"
	bucketName              = "evident-beacon-256523.appspot.com"

	noOfReqs int           = 1e2
	interval time.Duration = 0
	total    int64         = 1e6
)

var (
	defaultLogOutput = os.Stderr
	defaultPort      = "9001"
)

var (
	help    = flag.Bool("h", false, "display this help")
	started = false
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
	log.SetLevel(log.DebugLevel)

	ep, ok := os.LookupEnv("CONSUMER_ENDPOINT")
	if !ok {
		log.Warning("CONSUMER_ENDPOINT not set")
		ep = defaultConsumerEndpoint
	}

	ff, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to open log file"))
	}

	// prepare GCS metadata
	ctx := context.Background()
	// ctx := appengine.NewContext(r)
	log.SetOutput(io.MultiWriter(defaultLogOutput, ff))

	// bucketName, err := aefile.DefaultBucketName(context.Background())
	// if err != nil {
	// 	err = errors.Wrap(err, "failed to get default GCS bucket name")
	// 	log.Error(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	_, _ = w.Write([]byte(err.Error())) // ignore error
	// 	return
	// }
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create GCS client"))
	}
	runAll(ep)

	log.SetOutput(defaultLogOutput)
	defer ff.Close()
	_, err = ff.Seek(0, 0)
	if err != nil {
		log.Errorf("failed to seek: %v", err)
	}

	basename := fmt.Sprintf("%s-producer.log", time.Now().Format(time.RFC3339))
	ow := client.Bucket(bucketName).Object(basename).NewWriter(ctx)
	defer ow.Close()

	_, err = io.Copy(ow, ff)
	if err != nil {
		log.Errorf("failed to copy to GCS: %v", err)
	}

	// Set up file server
	// http.HandleFunc("/", startHandler)

	// port, ok := os.LookupEnv("PORT")
	// if !ok {
	// 	log.Warning("PORT not set")
	// 	port = defaultPort
	// }
	// log.Printf("serving on HTTP port: %s", port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	if started {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ep, ok := os.LookupEnv("CONSUMER_ENDPOINT")
	if !ok {
		log.Warning("CONSUMER_ENDPOINT not set")
		ep = defaultConsumerEndpoint
	}

	ff, err := ioutil.TempFile("", "")
	if err != nil {
		err = errors.Wrap(err, "failed to open log file")
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error())) // ignore error
		return
	}

	// prepare GCS metadata
	ctx := context.Background()
	// ctx := appengine.NewContext(r)
	log.SetOutput(io.MultiWriter(defaultLogOutput, ff))

	// start producer
	go func() {
		// bucketName, err := aefile.DefaultBucketName(context.Background())
		// if err != nil {
		// 	err = errors.Wrap(err, "failed to get default GCS bucket name")
		// 	log.Error(err)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	_, _ = w.Write([]byte(err.Error())) // ignore error
		// 	return
		// }
		bucketName := "evident-beacon-256523.appspot.com"
		client, err := storage.NewClient(ctx)
		if err != nil {
			err = errors.Wrap(err, "failed to create GCS client")
			log.Error(err)
		}
		runAll(ep)

		log.SetOutput(defaultLogOutput)
		defer ff.Close()
		_, err = ff.Seek(0, 0)
		if err != nil {
			log.Errorf("failed to seek: %v", err)
		}

		basename := fmt.Sprintf("%s-producer.log", time.Now().Format(time.RFC3339))
		ow := client.Bucket(bucketName).Object(basename).NewWriter(ctx)
		defer ow.Close()

		_, err = io.Copy(ow, ff)
		if err != nil {
			log.Errorf("failed to copy to GCS: %v", err)
		}
	}()

	started = true
	w.WriteHeader(http.StatusOK)
}

func runAll(consumerEndpoint string) {
	formats := []data.Format{data.FormatJSON, data.FormatProtobuf, data.FormatMsgpack, data.FormatAvro}
	types := []data.Type{data.TypeInt, data.TypeFloat, data.TypeString, data.TypeObject}

	for idx, f := range formats {
		for jdx, t := range types {
			log := log.WithFields(map[string]interface{}{"format": f, "typ": t})

			log.WithField("event", "progress").
				Debugf("running with %d/%d settings", (idx)*len(types)+jdx+1, len(formats)*len(types))
			run(consumerEndpoint, f, t)
			log.WithField("event", "progress").
				Debugf("done with %d/%d settings", (idx)*len(types)+jdx+1, len(formats)*len(types))
		}
	}
	log.Debug("finished")
}

func run(consumerEndpoint string, format data.Format, typ data.Type) {
	d, err := data.Create(typ, total)
	if err != nil {
		log.Fatal(err.Error())
	}

	p := prod.NewProducer(consumerEndpoint)
	p.Produce(d, format, noOfReqs, interval)
}

func displayHelp() {
	flag.PrintDefaults()
}
