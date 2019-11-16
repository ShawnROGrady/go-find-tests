package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ShawnROGrady/go-find-tests/finder"
	"github.com/ShawnROGrady/go-find-tests/tester"
)

func main() {
	var (
		path            string
		line, col       int
		err             error
		includeSubtests = flag.Bool("include-subs", false, "Find specific sub-tests which cover the specified block")
		printPositions  = flag.Bool("print-positions", false, "Print the positions of the found tests (NOTE: this does not currently work with subtests)")
		jsonFmt         = flag.Bool("json", false, "Print the output in json format")
		lineFmt         = flag.String("line-fmt", "%t:%f:%l:%c", "With -print-positions: the fmt to use when writing the postions of found tests. Structure:\n\t\t'%t': test name\n\t\t'%f': file\n\t\t'%l': line\n\t\t'%c': column\n\t\t'%o': offset\n\t")
		helpShort       = flag.Bool("h", false, "Print a help message and exit")
		help            = flag.Bool("help", false, "Print a help message and exit")
	)
	flag.Parse()
	if *help || *helpShort {
		fmt.Fprintf(os.Stdout, "Usage: %s [-include-subs] [-json|-line-fmt] filepath line col\n", os.Args[0])
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
		log.Fatalf("Error determining covering tests: %s", err)
	}

	if !*printPositions {
		if err := printTests(os.Stdout, coveredBy, *jsonFmt); err != nil {
			log.Fatalf("Error writing output: %s", err)
		}
		return
	}

	dir, _ := filepath.Split(path)
	allPositions, err := finder.PackageTests(dir)
	if err != nil {
		log.Fatalf("Error finding tests in %s: %s", dir, err)
	}

	coveringPositions := make(map[string]finder.TestPosition)
	for i := range coveredBy {
		if pos, ok := allPositions[coveredBy[i]]; ok {
			coveringPositions[coveredBy[i]] = pos
		}
	}
	if err := printCoveringPostions(os.Stdout, coveringPositions, *jsonFmt, *lineFmt); err != nil {
		log.Fatalf("Error writing output: %s", err)
	}
}
