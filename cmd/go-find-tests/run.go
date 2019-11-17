package main

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"

	"github.com/ShawnROGrady/go-find-tests/finder"
	"github.com/ShawnROGrady/go-find-tests/tester"
)

type runConfig struct {
	testerConf     tester.Config
	jsonFmt        bool
	lineFmt        string
	printPositions bool
}

func run(conf runConfig, path string, line, col int, dst io.Writer) error {
	t, err := tester.New(path, line, col, conf.testerConf)
	if err != nil {
		return fmt.Errorf("Error constructing tester: %s", err)
	}

	coveredBy, err := t.CoveredBy()
	if err != nil {
		return fmt.Errorf("Error determining covering tests: %s", err)
	}
	sort.Slice(coveredBy, func(i, j int) bool { return coveredBy[i] < coveredBy[j] })

	if !conf.printPositions {
		if err := printTests(dst, coveredBy, conf.jsonFmt); err != nil {
			return fmt.Errorf("Error writing output: %s", err)
		}
		return nil
	}

	dir, _ := filepath.Split(path)
	allPositions, err := finder.PackageTests(dir)
	if err != nil {
		return fmt.Errorf("Error finding tests in %s: %s", dir, err)
	}

	coveringPositions := make(map[string]finder.TestPosition)
	positionTests := []string{}
	for i := range coveredBy {
		if pos, ok := allPositions[coveredBy[i]]; ok {
			coveringPositions[coveredBy[i]] = pos
			// keep track of tests with positions to ensure consistent order when printing
			positionTests = append(positionTests, coveredBy[i])
		}
	}
	if err := printCoveringPostions(dst, coveringPositions, positionTests, conf.jsonFmt, conf.lineFmt); err != nil {
		return fmt.Errorf("Error writing output: %s", err)
	}

	return nil
}
