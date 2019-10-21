package tester

import (
	"context"
	"io"
	"runtime"
	"sync"

	"github.com/ShawnROGrady/go-find-tests/cover"
	"golang.org/x/sync/errgroup"
)

type coverFinder interface {
	coveringTests(t *Tester, testBin, outputDir string, allTests []string, includeSubtests bool) ([]string, error)
}

// synchronousFinder runs reach test synchronously
// primarily useful for establishing a baseline
type synchronousFinder struct{}

func (s synchronousFinder) coveringTests(t *Tester, testBin, outputDir string, allTests []string, includeSubtests bool) ([]string, error) {
	coveredBy := []string{}
	for i := range allTests {
		coverout, stdout, err := t.runCompiledTest(allTests[i], testBin, outputDir)
		if err != nil {
			return []string{}, err
		}

		prof, err := cover.New(coverout)
		if err != nil {
			return []string{}, err
		}
		if err := coverout.Close(); err != nil {
			return []string{}, err
		}

		if prof.Covers(t.testPos.file, t.testPos.line, t.testPos.col) {
			coveredBy = append(coveredBy, allTests[i])
			if includeSubtests {
				subTests := subtests(stdout)
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

func (s errGroupFinder) coveringTests(t *Tester, testBin, outputDir string, allTests []string, includeSubtests bool) ([]string, error) {
	var (
		ctx       = context.Background()
		coveredBy = []string{}
		tests     = make([]string, len(allTests))
		subs      = make([][]string, len(allTests))
	)
	g, _ := errgroup.WithContext(ctx)

	for i := range allTests {
		testNum := i
		testName := allTests[i]
		g.Go(func() error {
			coverout, stdout, err := t.runCompiledTest(testName, testBin, outputDir)
			if err != nil {
				return err
			}

			prof, err := cover.New(coverout)
			if err != nil {
				return err
			}
			if err := coverout.Close(); err != nil {
				return err
			}

			if prof.Covers(t.testPos.file, t.testPos.line, t.testPos.col) {
				tests[testNum] = testName
				if includeSubtests {
					subs[testNum] = subtests(stdout)
				}
			}
			return nil
		})
	}
	err := g.Wait()
	for i := range tests {
		if tests[i] != "" {
			coveredBy = append(coveredBy, tests[i])
			if includeSubtests {
				if len(subs[i]) != 0 {
					coveringSubs, err := s.coveringTests(t, testBin, outputDir, subs[i], false)
					if err != nil {
						return coveredBy, err
					}
					coveredBy = append(coveredBy, coveringSubs...)
				}
			}
		}
	}

	return coveredBy, err
}

type testRun struct {
	testName  string
	testBin   string
	outputDir string
}

func testGen(testBin, outputDir string, testNames ...string) <-chan testRun {
	out := make(chan testRun, len(testNames))

	for i := range testNames {
		out <- testRun{testName: testNames[i], testBin: testBin, outputDir: outputDir}
	}
	close(out)
	return out
}

type testRes struct {
	testName string
	cover    io.ReadCloser
	output   io.Reader
	err      error
}

func runTests(done <-chan struct{}, t *Tester, testRuns <-chan testRun, res chan<- testRes) {
	for run := range testRuns {
		coverOut, stdout, err := t.runCompiledTest(run.testName, run.testBin, run.outputDir)
		select {
		case res <- testRes{cover: coverOut, output: stdout, testName: run.testName, err: err}:
		case <-done:
			return
		}
	}
}

type coverRes struct {
	testName string
	err      error
	testOut  io.Reader
}

func findCoveringTests(done <-chan struct{}, t *Tester, testResults <-chan testRes, res chan<- coverRes) {
	for testRes := range testResults {
		if testRes.err != nil {
			select {
			case res <- coverRes{err: testRes.err}:
			case <-done:
			}
			return
		}
		prof, err := cover.New(testRes.cover)
		if err != nil {
			select {
			case res <- coverRes{err: testRes.err}:
			case <-done:
			}
			return
		}
		if err := testRes.cover.Close(); err != nil {
			select {
			case res <- coverRes{err: testRes.err}:
			case <-done:
			}
			return
		}

		tName := ""
		if prof.Covers(t.testPos.file, t.testPos.line, t.testPos.col) {
			tName = testRes.testName
		}
		select {
		case <-done:
			return
		case res <- coverRes{testName: tName, testOut: testRes.output}:
		}
	}
}

// pipelineFinder separates each component of test in to a pipeline ran by a specified number of workers
// inspired by: https://blog.golang.org/pipelines
type pipelineFinder struct {
	maxWorkers int
}

func (p pipelineFinder) coveringTests(t *Tester, testBin, outputDir string, allTests []string, includeSubtests bool) ([]string, error) {
	done := make(chan struct{})
	defer close(done)

	in := testGen(testBin, outputDir, allTests...)

	maxWorkers := p.maxWorkers
	if maxWorkers == 0 {
		maxWorkers = runtime.NumCPU() * 4
	}
	numWorkers := maxWorkers
	if len(allTests) < maxWorkers {
		numWorkers = len(allTests)
	}

	var runWg, findWg sync.WaitGroup
	runWg.Add(numWorkers)
	findWg.Add(numWorkers)

	testResults := make(chan testRes, numWorkers)
	coverResults := make(chan coverRes, numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			runTests(done, t, in, testResults)
			runWg.Done()
		}()
		go func() {
			findCoveringTests(done, t, testResults, coverResults)
			findWg.Done()
		}()
	}

	go func() {
		runWg.Wait()
		close(testResults)
	}()

	go func() {
		findWg.Wait()
		close(coverResults)
	}()

	coveredBy := []string{}

	for res := range coverResults {
		if res.err != nil {
			return []string{}, res.err
		}
		if res.testName != "" {
			coveredBy = append(coveredBy, res.testName)
			if includeSubtests {
				subTests := subtests(res.testOut)
				coveringSubs, err := p.coveringTests(t, testBin, outputDir, subTests, false)
				if err != nil {
					return coveredBy, err
				}
				coveredBy = append(coveredBy, coveringSubs...)
			}
		}
	}

	return coveredBy, nil
}
