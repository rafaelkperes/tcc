package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
	defaultConsumerEndpoint = "http://localhost:9000"
	bucketName              = "rkperes-storage"
	gcsCreds                = `{
		"type": "service_account",
		"project_id": "rkperes-storage",
		"private_key_id": "c826a7442fc6f5070c42138576d7cb1f60513246",
		"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC2UoQan8okMRgy\n3SzTjlhdUDnRMC8+Ua4wc70NSUwkWaf/JNSSk7MH8e8feS+r4alrFFrvF1J1FB1u\nuizELCGwV6heTDz8n2SY4rF1qhYfuSXfxmN0qy+Of/T6zOk/Hp5GidtAjwTfnQGY\nE2fOB09H+m0H+/qWgB2EP2rOKk9KdZsanF2liwLagaQ/GVt8yKA6slN90FlPAnoV\n43h7DLXPnwa7X/GMsWPsQeEFOEdIFmr+4lMXFPiFyJbHDFH/ug07V78usy8SetRL\nOJnErcK/lznSw3Zax1SqvkkxehdPPrCXku9pbBBWUt4P1jRa4PWPwWA43qAu/SGx\nyZJzGWplAgMBAAECggEAH9kaKmxvKxNIXtoz0mCzHtm8v8Xi+sfZ3azaAVAkdNUU\ne4U7fL9ALsscMitBII0ywvmzSMCSLtFsssLivwHWgK9PQemfXaGaOPqgdSVY6AG3\nk/dbuC2PCR1g9c6Fj/kRPNEn84cIGueaN65sG5k7SB9+nD5v74pBnbBWP900LJVC\nLYkLds/YcPYO25ceA4qenExugUDpiE37b1GdjQViQZX9VrOl8/siMumGysKBFC4r\nZslBPPy1H7tF16DBsKM0ylNfYjcmPSeYrx49yq5DjksX8eBD0xba+TMX5wzQ2V/M\nsq/Drw9oV5+b/XaSrN6Lsw2cyo6AxEx4ScGO5IG6AQKBgQD9H+rvcYoKAUnY1OQK\ntZqAmETmLVtzbbC2vmWZveyipeMZt9r+NlojUnl53+vgc5bG57uPe7vdG+QvXbY3\nsmPDRzLENHDtB2v97BwsmR9qtVEXP5Th2MCV6vZn2gLsY/yjAwnBxtNxmVzPUocj\nlw8ys+b0PQMwBhS074SlzdwZcQKBgQC4ZLTOoZfMSQ+52KHelUEJ1VhBLcOMw0aW\nf88opzv/r7i60WgLUU60v3JQ3ePXNMLqSSNdUVTEMSxf1T3sHaQqhfDiKjxYIUoZ\nqkV2c71BCiAR05RsXDgp6BNZgwK1bMOd5MKj6V6r0bbTlyWO/Qnr+BbkYIB8FUqJ\nPjPqLyeGNQKBgQCQ4hwPQeXJJEOooPKGTrxIrt+BAKU/xKFJeGGfRl6UGm+K4Pmw\nWFvvq91sLQdOSdsbrrhkwGYfgT9y/Si3aJxBwhcExx98DKt7hBH8VQjugyoPLI2D\nWBWjugGgH+FcfT61758+ExgkBaxh3tMLRAOm+eJQGjwg2NoxVoeOf+5jgQKBgAtT\nHAuovwLr5cxbMq3R6tmowa/XGLB3ecladiWgB75PU4Adxk8TokrVizbOOeUIt4Pe\nFA7yJMub3YbROOlcdK2r5jxtraEYAk4LOBLrTs9EyO1vWilBjK1+NFoGAs+Tq3vy\nBcY9WfQhgCIEoWjjv40/gmBqUNnOEPLW4Cdc2AeVAoGAdVEmr5Vn+0NQcjoN09sp\nOWHU+1LWJ6uvoPXA3V60x3K5MKT/z15vCOqokecjnhT+NKTl2l9BcvofBd03zbeu\noV/yjubwGVpjxFJY/3m+Tt4PNx4i8jSy3G2hl7JGJZU+hrmGjXtFW27cG6eLpbcs\nGCT80IFE91E2J6GzYm2D5X0=\n-----END PRIVATE KEY-----\n",
		"client_email": "rkperes-storage-sa@rkperes-storage.iam.gserviceaccount.com",
		"client_id": "110636495892929099742",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/rkperes-storage-sa%40rkperes-storage.iam.gserviceaccount.com"
	  }`

	noOfReqs int           = 1e2
	interval time.Duration = 0
	total    int64         = 1e5
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

	// parse port
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Warning("PORT not set")
		port = defaultPort
	}
	log.Printf("serving on HTTP port: %s", port)

	// setup logging
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)

	// start producer
	worker()

	// Set up file server
	http.HandleFunc("/_ah/start", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func worker() {
	ep, ok := os.LookupEnv("CONSUMER_ENDPOINT")
	if !ok {
		log.Warning("CONSUMER_ENDPOINT not set")
		ep = defaultConsumerEndpoint
	}

	go runAll(ep)
}

func runAll(consumerEndpoint string) {
	// prepare GCS metadata
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(gcsCreds)))
	if err != nil {
		log.Error(err)
	}

	formats := []data.Format{data.FormatJSON, data.FormatProtobuf, data.FormatMsgpack, data.FormatAvro}
	types := []data.Type{data.TypeInt, data.TypeFloat, data.TypeString, data.TypeObject}

	runID, ok := os.LookupEnv("RUN_ID")
	if !ok {
		runID = "default"
	}

	counterName := fmt.Sprintf("%s-counter", runID)
	formatCounterName := fmt.Sprintf("%s-format-counter", runID)
	for reqs := getCounter(ctx, client, counterName); reqs < noOfReqs; reqs++ {
		log.Debugf("main iteration %d/%d", reqs+1, noOfReqs)
		for idx := getCounter(ctx, client, formatCounterName); idx < len(formats); idx++ {
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

				format := strings.Join(strings.Split(string(f), "/"), "")
				basename := fmt.Sprintf("%s-%08d-%08d/%s-%s-producer.log", runID, noOfReqs, total, format, time.Now().Format(time.RFC3339))
				ow := client.Bucket(bucketName).Object(basename).NewWriter(ctx)
				defer ow.Close()

				_, err = io.Copy(ow, ff)
				if err != nil {
					log.Errorf("failed to copy to GCS: %v", err)
				}
			}()

			if err := setCounter(ctx, client, formatCounterName, idx+1); err != nil {
				log.Error(errors.Wrap(err, "could not save format counter value"))
			}
		}

		// zeroes format counter
		if err := setCounter(ctx, client, formatCounterName, 0); err != nil {
			log.Error(errors.Wrap(err, "could not save format counter value"))
		}

		if err := setCounter(ctx, client, counterName, reqs+1); err != nil {
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

func getCounter(ctx context.Context, client *storage.Client, filename string) int {
	o := client.Bucket(bucketName).Object(filename)
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

func setCounter(ctx context.Context, client *storage.Client, filename string, counter int) error {
	o := client.Bucket(bucketName).Object(filename)
	w := o.NewWriter(ctx)
	defer w.Close()

	_, err := w.Write([]byte(strconv.Itoa(counter)))
	return err
}

func run(consumerEndpoint string, format data.Format, typ data.Type) {
	d, err := data.Create(typ, total)
	if err != nil {
		log.Fatal(err.Error())
	}

	p := prod.NewProducer(consumerEndpoint)
	p.Produce(d, format, 1 /* number of requests */, interval)
}

func displayHelp() {
	flag.PrintDefaults()
}
