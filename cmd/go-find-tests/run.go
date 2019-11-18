package main

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

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

	var (
		coveringPositions map[string]*testPosition
		positionTests     []string
	)

	if !conf.testerConf.IncludeSubtests {
		coveringPositions = make(map[string]*testPosition)
		positionTests = []string{}
		for i := range coveredBy {
			if pos, ok := allPositions[coveredBy[i]]; ok {
				coveringPositions[coveredBy[i]] = &testPosition{TestPosition: pos}
				// keep track of tests with positions to ensure consistent order when printing
				positionTests = append(positionTests, coveredBy[i])
			}
		}
	} else {
		coveringPositions, positionTests = positionSubs(allPositions, coveredBy)
	}
	if err := printCoveringPostions(dst, coveringPositions, positionTests, conf.jsonFmt, conf.lineFmt); err != nil {
		return fmt.Errorf("Error writing output: %s", err)
	}

	return nil
}

func positionSubs(allPositions map[string]finder.TestPosition, coveredBy []string) (map[string]*testPosition, []string) {
	positions := make(map[string]*testPosition)
	positionTests := []string{}

	for i := range coveredBy {
		if pos, ok := allPositions[coveredBy[i]]; ok {
			// top level test
			positionTests = append(positionTests, coveredBy[i]) // keep track of tests with positions to ensure consistent order when printing
			if posInfo, ok := positions[coveredBy[i]]; ok {
				posInfo.TestPosition = pos
			} else {
				positions[coveredBy[i]] = &testPosition{
					TestPosition: pos,
				}
			}
			continue
		}

		if parts := strings.Split(coveredBy[i], "/"); len(parts) >= 2 {
			if posInfo, ok := positions[parts[0]]; ok {
				if len(posInfo.SubTests) == 0 {
					posInfo.SubTests = []string{coveredBy[i]}
				} else {
					posInfo.SubTests = append(posInfo.SubTests, coveredBy[i])
				}
			} else {
				positions[parts[0]] = &testPosition{
					SubTests: []string{coveredBy[i]},
				}
			}
		}
	}

	return positions, positionTests
}
