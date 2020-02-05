package tester

import (
	"context"
	"fmt"

	"github.com/ShawnROGrady/go-find-tests/cover"
	"golang.org/x/sync/errgroup"
)

type coverFinder interface {
	coveringTests(t *Tester, testBin, outputDir string, allTests []string, includeSubtests bool) ([]string, error)
}

/*
sequentialFinder runs reach test in sequence
primarily useful for integration tests
as well as establishing a baseline for benchmarks
*/
type sequentialFinder struct{}

func (s sequentialFinder) coveringTests(t *Tester, testBin, outputDir string, allTests []string, includeSubtests bool) ([]string, error) {
	coveredBy := []string{}
	for i := range allTests {
		coverout, stdout, err := t.runCompiledTest(allTests[i], testBin, outputDir)
		if err != nil {
			return []string{}, fmt.Errorf("error running test '%s': %s", allTests[i], err)
		}

		prof, err := cover.New(coverout)
		if err != nil {
			return []string{}, fmt.Errorf("error parsing coverage output: %s", err)
		}
		if err := coverout.Close(); err != nil {
			return []string{}, err
		}

		if prof.Covers(t.testPos.file, t.testPos.line, t.testPos.col) {
			coveredBy = append(coveredBy, allTests[i])
			if includeSubtests {
				subTests, err := subtests(stdout)
				if err != nil {
					return []string{}, fmt.Errorf("error finding subtests: %s", err)
				}
				if len(subTests) == 0 {
					continue
				}
				if len(subTests) == 1 {
					coveredBy = append(coveredBy, subTests[0])
					continue
				}
				coveringSubs, err := s.coveringTests(t, testBin, outputDir, subTests, false)
				if err != nil {
					return coveredBy, err
				}
				coveredBy = append(coveredBy, coveringSubs...)
			}
		}
	}
	return coveredBy, nil
}

// errGroupFinder runs each test in a separate go routine managed by an error group
type errGroupFinder struct{}

func (e errGroupFinder) coveringTests(t *Tester, testBin, outputDir string, allTests []string, includeSubtests bool) ([]string, error) {
	var (
		ctx       = context.Background()
		coveredBy = []string{}
		tests     = make([]string, len(allTests))
		subs      = make([][]string, len(allTests))
	)
	g, ctx := errgroup.WithContext(ctx)

	for i := range allTests {
		testNum := i
		testName := allTests[i]
		g.Go(func() error {
			return e.runTest(t, testBin, outputDir, testName, testNum, includeSubtests, tests, subs)
		})
	}
	if err := g.Wait(); err != nil {
		return []string{}, err
	}

	if includeSubtests {
		coveringSubs, err := e.coveringSubs(ctx, t, testBin, outputDir, tests, subs)
		if err != nil {
			return []string{}, err
		}
		for i := range tests {
			if tests[i] != "" {
				coveredBy = append(coveredBy, tests[i])
			}
			if len(coveringSubs[i]) != 0 {
				coveredBy = append(coveredBy, coveringSubs[i]...)
			}
		}
	} else {
		for i := range tests {
			if tests[i] != "" {
				coveredBy = append(coveredBy, tests[i])
			}
		}
	}

	return coveredBy, nil
}

func (e errGroupFinder) runTest(t *Tester, testBin, outputDir, testName string, testNum int, includeSubtests bool, tests []string, subs [][]string) error {
	coverout, stdout, err := t.runCompiledTest(testName, testBin, outputDir)
	if err != nil {
		return fmt.Errorf("error running test '%s': %s", testName, err)
	}

	prof, err := cover.New(coverout)
	if err != nil {
		return fmt.Errorf("error parsing coverage output: %s", err)
	}
	if err := coverout.Close(); err != nil {
		return err
	}

	if prof.Covers(t.testPos.file, t.testPos.line, t.testPos.col) {
		tests[testNum] = testName
		if includeSubtests {
			subs[testNum], err = subtests(stdout)
			if err != nil {
				return fmt.Errorf("error finding subtests: %s", err)
			}
		}
	}
	return nil
}

func (e errGroupFinder) coveringSubs(ctx context.Context, t *Tester, testBin, outputDir string, tests []string, subs [][]string) ([][]string, error) {
	errGroup, _ := errgroup.WithContext(ctx)
	coveringSubs := make([][]string, len(tests))
	for i := range tests {
		if tests[i] == "" {
			continue
		}
		testNum := i
		errGroup.Go(func() error {
			if len(subs[testNum]) != 0 {
				if len(subs[testNum]) == 1 {
					coveringSubs[testNum] = subs[testNum]
					return nil
				}
				coveringSubTests, err := e.coveringTests(t, testBin, outputDir, subs[testNum], false)
				if err != nil {
					return err
				}
				coveringSubs[testNum] = coveringSubTests
			}
			return nil
		})
	}
	if err := errGroup.Wait(); err != nil {
		return [][]string{}, err
	}
	return coveringSubs, nil
}
