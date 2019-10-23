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

func (s errGroupFinder) coveringTests(t *Tester, testBin, outputDir string, allTests []string, includeSubtests bool) ([]string, error) {
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
	if err := g.Wait(); err != nil {
		return []string{}, err
	}

	if includeSubtests {
		errGroup, _ := errgroup.WithContext(ctx)
		coveringSubs := make([][]string, len(tests))
		for i := range tests {
			if tests[i] != "" {
				coveredBy = append(coveredBy, tests[i])
			}
			testNum := i
			errGroup.Go(func() error {
				if len(subs[testNum]) != 0 {
					if len(subs[testNum]) == 1 {
						coveringSubs[testNum] = subs[testNum]
						return nil
					}
					coveringSubTests, err := s.coveringTests(t, testBin, outputDir, subs[testNum], false)
					if err != nil {
						return err
					}
					coveringSubs[testNum] = coveringSubTests
				}
				return nil
			})
		}
		if err := errGroup.Wait(); err != nil {
			return []string{}, err
		}
		for i := range coveredBy {
			if len(coveringSubs[i]) != 0 {
				coveredBy = append(coveredBy[:i], append(coveringSubs[i], coveredBy[i:]...)...)
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

	if includeSubtests {
		var allCoveringTests = make([]subRes, len(allTests))
		if err := p.findCoveringSubs(t, testBin, outputDir, coverResults, allCoveringTests); err != nil {
			return []string{}, err
		}
		for i, res := range allCoveringTests {
			if res.err != nil {
				return []string{}, res.err
			}
			if res.testName != "" {
				coveredBy = append(coveredBy, res.testName)
				if coveringSubs := allCoveringTests[i].subtests; len(coveringSubs) != 0 {
					coveredBy = append(coveredBy, coveringSubs...)
				}
			}
		}
	} else {
		for res := range coverResults {
			if res.err != nil {
				return []string{}, res.err
			}
			if res.testName != "" {
				coveredBy = append(coveredBy, res.testName)
			}
		}
	}

	return coveredBy, nil
}

type subRes struct {
	subtests []string
	coverRes
}

func (p pipelineFinder) findCoveringSubs(t *Tester, testBin, outputDir string, coverResults chan coverRes, dst []subRes) error {
	var wg sync.WaitGroup
	i := 0

	for res := range coverResults {
		if res.err != nil {
			return res.err
		}

		dst[i] = subRes{coverRes: res, subtests: []string{}}
		if res.testName != "" {
			wg.Add(1)
			go func(i int, res coverRes) {
				subTests := subtests(res.testOut)
				if len(subTests) == 0 {
					wg.Done()
					return
				}
				if len(subTests) == 1 {
					dst[i].subtests = subTests
					wg.Done()
					return
				}
				coveringSubs, err := p.coveringTests(t, testBin, outputDir, subTests, false)
				if err != nil {
					dst[i].err = err
					wg.Done()
					return
				}

				dst[i].subtests = coveringSubs
				wg.Done()
			}(i, res)
			i++
		}
	}
	wg.Wait()
	return nil
}
