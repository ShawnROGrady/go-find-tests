package main

import (
	"flag"
	"fmt"
	"log"
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
	)
	flag.Parse()
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
	fmt.Printf("%s\n", coveredBy)
}
