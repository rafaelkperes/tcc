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
	bucketName = "rkperes-storage"
	gcsCreds   = `{
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
)

var (
	help = flag.Bool("h", false, "display this help")
	_f   = flag.String("f", "", "use file as input (instead of stdin)")
	_o   = flag.String("o", "", "output to file (instead of stdout)")
)

type entry struct {
	msr    measure.Measure
	size   int64
	typ    string
	format string
}

func main() {
	// flag-related init
	flag.Parse()
	if *help {
		flag.PrintDefaults()
	}

	// in := os.Stdin
	// if len(*_f) > 0 {
	// 	var err error
	// 	in, err = os.Open(*_f)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	defer in.Close()
	// }

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(gcsCreds)))
	if err != nil {
		log.Error(err)
	}

	in := bytes.NewBuffer([]byte{})

	it := client.Bucket(bucketName).Objects(ctx, &storage.Query{Prefix: "real4-00000100-00100000/"})
	for obj, err := it.Next(); err != iterator.Done; obj, err = it.Next() {
		if err != nil {
			log.Fatal(err)
		}

		func() {
			r, err := client.Bucket(obj.Bucket).Object(obj.Name).NewReader(ctx)
			if err != nil {
				log.Fatal(errors.Wrapf(err, "failed to open %s", obj.Name))
			}
			defer r.Close()

			_, err = io.Copy(in, r)
			if err != nil {
				log.Fatal(errors.Wrapf(err, "failed to copy from %s to buffer", obj.Name))
			}
		}()
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
	err = w.Write([]string{
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
}
