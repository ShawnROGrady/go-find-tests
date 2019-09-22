package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/ShawnROGrady/go-find-tests/tester"
)

func main() {
	// TODO: parse args containing file + row + col
	// TODO: output tests covering combination
	var (
		path      string
		line, col int
		err       error
	)
	if len(os.Args) < 3 {
		log.Fatal("Path and line are required")
	}

	path = os.Args[1]
	if line, err = strconv.Atoi(os.Args[2]); err != nil {
		log.Fatalf("Invalid line argument: %s", err)
	}
	if len(os.Args) > 3 {
		if col, err = strconv.Atoi(os.Args[3]); err != nil {
			log.Fatalf("Invalid column argument: %s", err)
		}
	}

	t, err := tester.New(path, line, col)
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
