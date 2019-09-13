package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rafaelkperes/tcc/internal/svc/cons"
)

const (
	helpArg = "help"
)

var execName = filepath.Base(os.Args[0])

var help = flag.Bool("h", false, "display this help")
var port = flag.Int("p", 9000, "set consumer port")

func main() {
	flag.Parse()

	if *help {
		displayHelp()
		os.Exit(0)
	}

	cfg := &cons.Config{}

	log.Printf("listen and serve at port %d", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), cons.NewConsumerServer(cfg))
	log.Fatal(err)
}

func displayHelp() {
	flag.PrintDefaults()
}
