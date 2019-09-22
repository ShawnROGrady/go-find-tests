package tester

import (
	"context"
	"strings"

	"github.com/ShawnROGrady/go-find-tests/cover"
	"golang.org/x/sync/errgroup"
)

type coverFinder interface {
	coveringTests(t *Tester, outputDir string, allTests []string) ([]string, error)
}

// synchronousFinder runs reach test synchronously
// primarily useful for establishing a baseline
type synchronousFinder struct{}

func (s synchronousFinder) coveringTests(t *Tester, outputDir string, allTests []string) ([]string, error) {
	coveredBy := []string{}
	for i := range allTests {
		var dst strings.Builder
		dst.WriteString(outputDir)
		dst.WriteString(allTests[i])
		dst.WriteString(".out")
		output, err := t.runTest(allTests[i], dst.String())
		if err != nil {
			return []string{}, err
		}

		prof, err := cover.New(output)
		if err != nil {
			return []string{}, err
		}
		if err := output.Close(); err != nil {
			return []string{}, err
		}

		if prof.Covers(t.testPos.file, t.testPos.line, t.testPos.col) {
			coveredBy = append(coveredBy, allTests[i])
		}
	}
	return coveredBy, nil
}

// errGroupFinder runs each test in a separate go routine managed by an error group
type errGroupFinder struct{}

func (s errGroupFinder) coveringTests(t *Tester, outputDir string, allTests []string) ([]string, error) {
	var (
		ctx       = context.Background()
		coveredBy = []string{}
		tests     = make([]string, len(allTests))
	)
	g, _ := errgroup.WithContext(ctx)

	for i := range allTests {
		testNum := i
		testName := allTests[i]
		g.Go(func() error {
			var dst strings.Builder
			dst.WriteString(outputDir)
			dst.WriteString(testName)
			dst.WriteString(".out")
			output, err := t.runTest(testName, dst.String())
			if err != nil {
				return err
			}

			prof, err := cover.New(output)
			if err != nil {
				return err
			}
			if err := output.Close(); err != nil {
				return err
			}

			if prof.Covers(t.testPos.file, t.testPos.line, t.testPos.col) {
				tests[testNum] = testName
			}
			return nil
		})
	}
	err := g.Wait()
	for i := range tests {
		if tests[i] != "" {
			coveredBy = append(coveredBy, tests[i])
		}
	}
	return coveredBy, err
}
