package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/rafaelkperes/tcc/internal/svc/cons"
	log "github.com/sirupsen/logrus"
)

const (
	envPort = "PORT"

	defaultPort = "9000"
)

var (
	_help = flag.Bool("h", false, "display this help")
)

func main() {
	// flag-related init
	flag.Parse()
	if *_help {
		displayHelp()
		os.Exit(0)
	}

	port, ok := os.LookupEnv(envPort)
	if !ok {
		log.Warningf("%s not set", envPort)
		port = defaultPort
	}
	log.Infof("serving on HTTP port %s", port)

	// setup logging
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)

	log.WithFields(log.Fields{"port": port}).Infof("listen and serve at port %s", port)
	// setup server
	http.HandleFunc("/_ah/start", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/_ah/stop", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.Handle("/", cons.NewConsumerServer())
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func displayHelp() {
	flag.PrintDefaults()
}
