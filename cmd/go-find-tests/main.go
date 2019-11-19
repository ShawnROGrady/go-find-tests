package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/ShawnROGrady/go-find-tests/tester"
)

const (
	defaultLineFmt = "%t:%f:%l:%c:%s"
)

func main() {
	var (
		includeSubtests = flag.Bool("include-subs", false, "Find specific sub-tests which cover the specified block")
		short           = flag.Bool("short", false, "Sets '-short' flag when testing for coverage")
		runExpr         = flag.String("run", ".", "Check only tests matching the regular expression")
		printPositions  = flag.Bool("print-positions", false, "Print the positions of the found tests")
		jsonFmt         = flag.Bool("json", false, "Print the output in json format")
		lineFmt         = flag.String("line-fmt", defaultLineFmt, "With -print-positions: the fmt to use when writing the postions of found tests. Structure:\n\t\t'%t': test name\n\t\t'%f': file\n\t\t'%l': line\n\t\t'%c': column\n\t\t'%o': offset\n\t'%s': subtests (printed as comma separated list)")
		runSeq          = flag.Bool("seq", false, "Run all tests sequentially. Greatly reduces performance but may be neccessary for integration tests")
		helpShort       = flag.Bool("h", false, "Print a help message and exit")
		help            = flag.Bool("help", false, "Print a help message and exit")
	)
	flag.Parse()
	if *help || *helpShort {
		fmt.Fprintf(os.Stdout, "Usage: %s [-include-subs] [-short] [-run regexp] [-json|-line-fmt regexp] filepath:line[.col]\n", os.Args[0])
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
	if len(args) == 0 {
		log.Fatal("Position argument (fmt = 'file:line[.col]') required")
	}

	pos, err := parsePosition(args[0])
	if err != nil {
		log.Fatalf("Error parsing position arg: %s", err)
	}

	conf := runConfig{
		testerConf: tester.Config{
			IncludeSubtests: *includeSubtests,
			Short:           *short,
			Run:             *runExpr,
			Seq:             *runSeq,
		},
		jsonFmt:        *jsonFmt,
		lineFmt:        *lineFmt,
		printPositions: *printPositions,
	}

	if err := run(conf, pos.file, pos.line, pos.col, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

type pos struct {
	file      string
	line, col int
}

func parsePosition(arg string) (*pos, error) {
	argFmt := `^([a-zA-Z0-9\/\.\-_]+.go):([0-9]+)(?:\.([0-9]+))?`
	argReg := regexp.MustCompile(argFmt)

	subexps := argReg.FindStringSubmatch(arg)
	if len(subexps) == 0 {
		return nil, errors.New("provided position doesn't match format 'file:line.column'")
	}

	line, err := strconv.Atoi(subexps[2])
	if err != nil {
		return nil, err
	}

	var col int
	if len(subexps) == 4 && subexps[3] != "" {
		col, err = strconv.Atoi(subexps[3])
		if err != nil {
			return nil, err
		}
	}

	return &pos{
		file: subexps[1],
		line: line,
		col:  col,
	}, nil
}
