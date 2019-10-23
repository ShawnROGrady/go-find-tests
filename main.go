package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/ShawnROGrady/go-find-tests/tester"
)

func main() {
	var (
		path            string
		line, col       int
		err             error
		includeSubtests = flag.Bool("include-subs", false, "Find specific sub-tests which cover the specified block")
		helpShort       = flag.Bool("h", false, "Print a help message and exit")
		help            = flag.Bool("help", false, "Print a help message and exit")
	)
	flag.Parse()
	if *help || *helpShort {
		fmt.Fprintf(os.Stdout, "Usage: %s [-include-subs] filepath line col\n", os.Args[0])
		fmt.Fprintf(os.Stdout, "Description: %s prints the tests (and optionally sub tests) which cover a specified block of code\n", os.Args[0])
		fmt.Fprint(os.Stdout, "Required arguments:\n")
		fmt.Fprint(os.Stdout, "\tfilepath: path to the file to check\n")
		fmt.Fprint(os.Stdout, "\tline: the line of the block to check\n")
		fmt.Fprint(os.Stdout, "Optional arguments:\n")
		fmt.Fprint(os.Stdout, "\trow: the row of the block to check\n")
		fmt.Fprint(os.Stdout, "Optional flags:\n")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(os.Stdout, "\t-%s: %s [default = %v]\n", f.Name, f.Usage, f.DefValue)
		})
		os.Exit(0)
	}
	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("Path and line are required")
	}

	path = args[0]
	if line, err = strconv.Atoi(args[1]); err != nil {
		log.Fatalf("Invalid line argument: %s", err)
	}
	if len(args) > 2 {
		if col, err = strconv.Atoi(args[2]); err != nil {
			log.Fatalf("Invalid column argument: %s", err)
		}
	}

	t, err := tester.New(path, line, col, *includeSubtests)
	if err != nil {
		log.Fatalf("Error constructing tester: %s", err)
	}

	coveredBy, err := t.CoveredBy()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			log.Fatalf("Error determining covering tests: %s", exitErr.Stderr)
		} else {
			log.Fatalf("Error determining covering tests: %#v", err)
		}
	}
	fmt.Fprintf(os.Stdout, "%s\n", coveredBy)
}
