package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/rafaelkperes/tcc/pkg/gen"
)

const (
	genArg  = "gen"
	helpArg = "help"
)

var execName = filepath.Base(os.Args[0])

var help = flag.Bool("h", false, "display this help")
var dir = flag.String("d", "./data", fmt.Sprintf("%s %s flag: sets the target dir", execName, genArg))
var size = flag.Int("n", 1e6, fmt.Sprintf("%s %s flag: sets the amount of entries generated", execName, genArg))
var seed = flag.Int64("s", 42, fmt.Sprintf("%s %s flag: sets the random seed", execName, genArg))

func main() {
	flag.Parse()

	if *help {
		displayHelp()
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		fmt.Fprintf(flag.CommandLine.Output(), "expected a single argument, got %d arguments\n\n", flag.NArg())
		displayHelp()
		os.Exit(2)
	}

	switch flag.Arg(0) {
	case genArg:
		genMain()
	case helpArg:
		displayHelp()
		os.Exit(0)
	default:
		fmt.Fprintf(flag.CommandLine.Output(), "unexpected argument '%s'\n\n", flag.Arg(0))
		displayHelp()
		os.Exit(2)
	}
}

func genMain() {
	rand.Seed(*seed)

	err := gen.All(*dir, *size)
	if err != nil {
		fmt.Printf("failed to generate data: %+s\n", err)
		os.Exit(1)
	}
}

func displayHelp() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, "%s %s  - generate data files\n", execName, genArg)
	fmt.Fprintf(w, "%s %s - display this help\n\n", execName, helpArg)
	fmt.Fprintln(w, "additional flags:")
	flag.PrintDefaults()
}
