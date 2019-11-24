package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"github.com/rafaelkperes/tcc/internal/svc/prod"
	"google.golang.org/api/option"

	"github.com/rafaelkperes/tcc/pkg/data"
	log "github.com/sirupsen/logrus"
)

const (
	envPort             = "PORT"
	envBucket           = "BUCKET"
	envPrefix           = "PREFIX"
	envNoOfReqs         = "NUMBER_OF_REQUESTS"
	envPayloadSize      = "PAYLOAD_SIZE"
	envConsumerEndpoint = "CONSUMER_ENDPOINT"

	defaultConsumerEndpoint               = "http://localhost:9000"
	interval                time.Duration = 0
)

var (
	defaultLogOutput       = os.Stderr
	defaultPort            = "9001"
	defaultPrefix          = "default"
	defaultNoOfReqs    int = 1
	defaultPayloadSize int = 1e3

	noOfReqs    int
	payloadSize int
)

var (
	help  = flag.Bool("h", false, "display this help")
	_auth = flag.String("auth", "", "path to file with service account key")
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

	// parse variables
	port, ok := os.LookupEnv(envPort)
	if !ok {
		log.Warningf("%s not set", envPort)
		port = defaultPort
	}
	log.Infof("serving on HTTP port %s", port)

	bucket, ok := os.LookupEnv(envBucket)
	if !ok {
		log.Fatalf("%s not set", envBucket)
	}
	prefix, ok := os.LookupEnv(envPrefix)
	if !ok {
		log.Warningf("%s not set", envPrefix)
		prefix = defaultPrefix
	}
	log.Infof("writing to bucket %s with prefix %s", bucket, prefix)

	if n, err := strconv.ParseFloat(os.Getenv(envNoOfReqs), 64); err != nil {
		log.Warningf("%s not valid", envNoOfReqs)
		noOfReqs = defaultNoOfReqs
	} else {
		noOfReqs = int(n)
	}
	if pSize, err := strconv.ParseFloat(os.Getenv(envPayloadSize), 64); err != nil {
		log.Warningf("%s not valid", envPayloadSize)
		noOfReqs = defaultPayloadSize
	} else {
		payloadSize = int(pSize)
	}
	log.Infof("requesting a total of %d requests with size %d", noOfReqs, payloadSize)
	consumerEP, ok := os.LookupEnv(envConsumerEndpoint)
	if !ok {
		log.Warningf("%s not set", envConsumerEndpoint)
		consumerEP = defaultConsumerEndpoint
	}
	log.Infof("requesting to consumer %s", consumerEP)

	// create storage client
	auth := *_auth
	var opts []option.ClientOption
	if len(auth) > 0 {
		// use auth key
		opts = append(opts, option.WithCredentialsFile(auth))
	}
	ctx := context.Background()
	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		log.Fatal(errors.Wrap(err, "while creating storage client"))
	}

	// start producer
	worker(ctx, consumerEP, client, bucket, prefix)
	// func() {
	// 	w := client.Bucket(bucket).Object(path.Join(prefix, "test")).NewWriter(ctx)
	// 	defer w.Close()
	// 	if _, err := io.Copy(w, strings.NewReader("Hello, Producer!")); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Info("finished writing")
	// }()

	// setup server
	http.HandleFunc("/_ah/start", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/_ah/stop", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("up and running"))
		w.WriteHeader(http.StatusOK)
	})
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func worker(ctx context.Context, ep string, client *storage.Client, bucket, prefix string) {
	go runAll(ctx, ep, client, bucket, prefix)
}

func runAll(ctx context.Context, consumerEndpoint string, client *storage.Client, bucket, prefix string) {
	formats := []data.Format{data.FormatJSON, data.FormatProtobuf, data.FormatMsgpack, data.FormatAvro}
	types := []data.Type{data.TypeInt, data.TypeFloat, data.TypeString, data.TypeObject}

	counterName := path.Join(prefix, "counter")
	formatCounterName := path.Join(prefix, "counter-format")
	for reqs := getCounter(ctx, client, bucket, counterName); reqs < noOfReqs; reqs++ {
		log.Debugf("main iteration %d/%d", reqs+1, noOfReqs)
		for idx := getCounter(ctx, client, bucket, formatCounterName); idx < len(formats); idx++ {
			f := formats[idx]
			log.Debugf("format iteration %d/%d (%s)", idx+1, len(formats), string(f))

			func() {
				ff, err := ioutil.TempFile("", "")
				if err != nil {
					log.Fatal(errors.Wrap(err, "failed to open log file"))
				}
				defer ff.Close()

				log.SetOutput(io.MultiWriter(defaultLogOutput, ff))
				runFormat(f, types, consumerEndpoint)
				log.SetOutput(defaultLogOutput)

				_, err = ff.Seek(0, 0)
				if err != nil {
					log.Errorf("failed to seek: %v", err)
				}

				formatSplit := strings.Split(string(f), "/")
				basename := path.Join(prefix, fmt.Sprintf("%s-%s.log", formatSplit[len(formatSplit)-1], time.Now().Format(time.RFC3339)))
				ow := client.Bucket(bucket).Object(basename).NewWriter(ctx)
				defer ow.Close()

				_, err = io.Copy(ow, ff)
				if err != nil {
					log.Errorf("failed to copy to GCS: %v", err)
				}
			}()

			if err := setCounter(ctx, client, bucket, formatCounterName, idx+1); err != nil {
				log.Error(errors.Wrap(err, "could not save format counter value"))
			}
		}

		// zeroes format counter
		if err := setCounter(ctx, client, bucket, formatCounterName, 0); err != nil {
			log.Error(errors.Wrap(err, "could not save format counter value"))
		}

		if err := setCounter(ctx, client, bucket, counterName, reqs+1); err != nil {
			log.Error(errors.Wrap(err, "could not save counter value"))
		}
	}
	log.Debug("finished")
}

func runFormat(f data.Format, types []data.Type, consumerEndpoint string) {
	for _, t := range types {
		run(consumerEndpoint, f, t)
	}
}

func getCounter(ctx context.Context, client *storage.Client, bucket, filename string) int {
	o := client.Bucket(bucket).Object(filename)
	r, err := o.NewReader(ctx)
	if err != nil {
		return 0
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return 0
	}

	i, err := strconv.Atoi(string(b))
	if err != nil {
		return 0
	}

	return i
}

func setCounter(ctx context.Context, client *storage.Client, bucket, filename string, counter int) error {
	o := client.Bucket(bucket).Object(filename)
	w := o.NewWriter(ctx)
	defer w.Close()

	_, err := w.Write([]byte(strconv.Itoa(counter)))
	return err
}

func run(consumerEndpoint string, format data.Format, typ data.Type) {
	d, err := data.Create(typ, int64(payloadSize))
	if err != nil {
		log.Fatal(err.Error())
	}

	p := prod.NewProducer(consumerEndpoint)
	p.Produce(d, format, 1 /* number of requests */, interval)
}

func displayHelp() {
	flag.PrintDefaults()
}
