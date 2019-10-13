package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/rafaelkperes/tcc/internal/svc/cons"
	log "github.com/sirupsen/logrus"
)

var (
	_help = flag.Bool("h", false, "display this help")
	_port = flag.Int("p", 9000, "set consumer port")
)

func main() {
	// flag-related init
	flag.Parse()
	if *_help {
		displayHelp()
		os.Exit(0)
	}

	port := *_port
	if sp := os.Getenv("PORT"); len(sp) > 0 {
		p, err := strconv.Atoi(sp)
		if err != nil {
			log.WithError(err).Error("invalid format in PORT environment variable")
		} else {
			port = p
		}
	}

	// setup logging
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)

	log.WithFields(log.Fields{"port": port}).Infof("listen and serve at port %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), cons.NewConsumerServer())
	log.Fatal(err)
}

func displayHelp() {
	flag.PrintDefaults()
}
