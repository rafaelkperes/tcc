package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	helpArg = "help"
)

var execName = filepath.Base(os.Args[0])

var help = flag.Bool("h", false, "display this help")

func main() {
	flag.Parse()

	if *help {
		displayHelp()
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		fmt.Fprintf(flag.CommandLine.Output(), "expected a single argument, got %d arguments\n\n", flag.NArg())
		displayHelp()
		os.Exit(2)
	}

	switch flag.Arg(0) {
	default:
		fmt.Fprintf(flag.CommandLine.Output(), "unexpected argument '%s'\n\n", flag.Arg(0))
		displayHelp()
		os.Exit(2)
	}
}

func displayHelp() {
	w := flag.CommandLine.Output()
	fmt.Fprintln(w, "additional flags:")
	flag.PrintDefaults()
}
