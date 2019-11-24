package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"github.com/rafaelkperes/tcc/pkg/measure"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	log "github.com/sirupsen/logrus"
)

// parser parses the producer/consumer log measures.

const (
	timeFormat = time.RFC3339Nano
)

var (
	help    = flag.Bool("h", false, "display this help")
	_bucket = flag.String("bucket", "", "bucket to retrieve data from")
	_prefix = flag.String("prefix", "", "prefix in the bucket where the logs are")
	_auth   = flag.String("auth", "", "path to file with service account key")
)

type entry struct {
	msr    measure.Measure
	size   int64
	typ    string
	format string
}

func main() {
	log.SetLevel(log.DebugLevel)

	// flag-related init
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// validate parameters
	auth := *_auth
	bucket := *_bucket
	prefix := *_prefix

	if len(bucket) == 0 {
		log.Fatal("unexpected empty bucket parameter")
	}

	// create client
	var opts []option.ClientOption
	if len(auth) > 0 {
		// use auth key
		log.Info("using auth parameter")
		opts = append(opts, option.WithCredentialsFile(auth))
	}
	ctx := context.Background()
	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		log.Fatalf("for GCS client: %s", err)
	}

	in := bytes.NewBuffer([]byte{})
	if err := readAll(ctx, client, bucket, prefix, in); err != nil {
		log.Fatal(errors.Wrap(err, "while reading storage entries"))
	}

	if err := parse(in, os.Stdout); err != nil {
		log.Fatal(errors.Wrap(err, "while parsing entries"))
	}
}

func readAll(ctx context.Context, client *storage.Client, bucket, prefix string, buff io.Writer) error {
	log.Infof("read from bucket %s with prefix %s", bucket, prefix)
	it := client.Bucket(bucket).Objects(ctx, &storage.Query{Prefix: prefix})
	for obj, err := it.Next(); err != iterator.Done; obj, err = it.Next() {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(obj.Name, ".log") {
			// skip non-log files
			continue
		}

		func() {
			r, err := client.Bucket(obj.Bucket).Object(obj.Name).NewReader(ctx)
			if err != nil {
				log.Fatal(errors.Wrapf(err, "failed to open %s", obj.Name))
			}
			defer r.Close()

			_, err = io.Copy(buff, r)
			if err != nil {
				log.Fatal(errors.Wrapf(err, "failed to copy from %s to buffer", obj.Name))
			}
		}()
	}
	return nil
}

func parse(buff io.Reader, out io.Writer) error {
	// parse entries
	entries := make([]entry, 0)
	s := bufio.NewScanner(buff)
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
		datasize, ok := o["datasize"].(float64)
		if !ok {
			log.Warningf("not a valid data size: %v", datasize)
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
			size:   int64(datasize),
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
		"size",
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
			strconv.FormatInt(it.size, 10),
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

	return nil
}
