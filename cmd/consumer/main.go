package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rafaelkperes/tcc/internal/svc/cons"
	log "github.com/sirupsen/logrus"
)

var (
	help = flag.Bool("h", false, "display this help")
	port = flag.Int("p", 9000, "set consumer port")
	lf   = flag.String("lf", "/var/log/consumer.std.log", "standard logs file")
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

	log.WithFields(log.Fields{"port": *port}).Infof("listen and serve at port %d", *port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), cons.NewConsumerServer())
	log.Fatal(err)
}

func displayHelp() {
	flag.PrintDefaults()
}
