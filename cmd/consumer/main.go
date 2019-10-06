package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/rafaelkperes/tcc/internal/svc/cons"
	log "github.com/sirupsen/logrus"
)

var (
	help = flag.Bool("h", false, "display this help")
	port = flag.Int("p", 9000, "set consumer port")
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

	log.WithFields(log.Fields{"port": *port}).Infof("listen and serve at port %d", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), cons.NewConsumerServer())
	log.Fatal(err)
}

func displayHelp() {
	flag.PrintDefaults()
}
