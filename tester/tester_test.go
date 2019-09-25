package tester

import (
	"fmt"
	"os/exec"
	"runtime"
	"sort"
	"testing"
)

const (
	sequentialName    = "sequential"
	errGroupName      = "err_group"
	pipelineCPUName   = "pipeline_num_cpu"
	pipelineCPU2Name  = "pipeline_num_cpu*2"
	pipelineCPU4Name  = "pipeline_num_cpu*4"
	pipelineCPU6Name  = "pipeline_num_cpu*6"
	pipelineCPU8Name  = "pipeline_num_cpu*8"
	pipelineCPU10Name = "pipeline_num_cpu*10"
)

var shortSkip = map[string]bool{
	pipelineCPU2Name:  true,
	pipelineCPU4Name:  true,
	pipelineCPU6Name:  true,
	pipelineCPU8Name:  true,
	pipelineCPU10Name: true,
}

var allFinders = []struct {
	newFinder func() coverFinder
	name      string
}{
	{
		name:      sequentialName,
		newFinder: func() coverFinder { return synchronousFinder{} },
	},
	{
		name:      errGroupName,
		newFinder: func() coverFinder { return errGroupFinder{} },
	},
	{
		name:      pipelineCPUName,
		newFinder: func() coverFinder { return pipelineFinder{maxWorkers: runtime.NumCPU()} },
	},
	{
		name:      pipelineCPU2Name,
		newFinder: func() coverFinder { return pipelineFinder{maxWorkers: runtime.NumCPU() * 2} },
	},
	{
		name:      pipelineCPU4Name,
		newFinder: func() coverFinder { return pipelineFinder{maxWorkers: runtime.NumCPU() * 4} },
	},
	{
		name:      pipelineCPU6Name,
		newFinder: func() coverFinder { return pipelineFinder{maxWorkers: runtime.NumCPU() * 6} },
	},
	{
		name:      pipelineCPU8Name,
		newFinder: func() coverFinder { return pipelineFinder{maxWorkers: runtime.NumCPU() * 8} },
	},
	{
		name:      pipelineCPU10Name,
		newFinder: func() coverFinder { return pipelineFinder{maxWorkers: runtime.NumCPU() * 10} },
	},
}

var coveredByTests = map[string]struct {
	fileDir         string
	fileName        string
	line, col       int
	expectCoveredBy []string
	expectErr       bool
}{
	"covered_by_1_of_4_tests": {
		fileDir:  "size",
		fileName: "size.go",
		line:     22, col: 0, // body of isEnormous()
		expectCoveredBy: []string{"TestIsEnormous"},
	},
	"covered_by_3_of_4_tests": {
		fileDir:  "size",
		fileName: "size.go",
		line:     8, col: 0, // negative case of size()
		expectCoveredBy: []string{"TestSize", "TestNegativeSize", "TestIsNegative"},
	},
	"covered_by_0_of_4_tests": {
		fileDir:  "size",
		fileName: "size.go",
		line:     10, col: 0, // zero case of size()
		expectCoveredBy: []string{},
	},
	"invalid_path": {
		fileDir:   "bad_path",
		fileName:  "size.go",
		expectErr: true,
	},
	"failing_test": {
		fileDir:  "failing",
		fileName: "fail.go",
		line:     5, col: 0,
		expectErr: true,
	},
	"covered_by_5_of_25_tests": {
		fileDir:  "len25",
		fileName: "len.go",
		line:     10, col: 0, // zero case of length()
		expectCoveredBy: []string{
			"TestEmptyStringIsEmpty",
			"TestEmptyStringIsShort",
			"TestEmptyStringIsLong",
			"TestEmptyStringIsVeryLong",
			"TestEmptyStringIsNovel",
		},
	},
	"covered_by_4_of_20_tests": {
		fileDir:  "len20",
		fileName: "len.go",
		line:     12, col: 0, // short case of length()
		expectCoveredBy: []string{
			"TestShortStringIsEmpty",
			"TestShortStringIsShort",
			"TestShortStringIsLong",
			"TestShortStringIsVeryLong",
		},
	},
	"covered_by_5_of_15_tests": {
		fileDir:  "len15",
		fileName: "len.go",
		line:     30, col: 0, // body of isLong()
		expectCoveredBy: []string{
			"TestEmptyStringIsLong",
			"TestShortStringIsLong",
			"TestLongStringIsLong",
			"TestVeryLongStringIsLong",
			"TestNovelIsLong",
		},
	},
	"covered_by_5_of_10_tests": {
		fileDir:  "len10",
		fileName: "len.go",
		line:     22, col: 0, // body of isEmpty()
		expectCoveredBy: []string{
			"TestEmptyStringIsEmpty",
			"TestShortStringIsEmpty",
			"TestLongStringIsEmpty",
			"TestVeryLongStringIsEmpty",
			"TestNovelIsEmpty",
		},
	},
}

func TestCoveredBy(t *testing.T) {
	for testName, test := range coveredByTests {
		for i := range allFinders {
			t.Run(fmt.Sprintf("%s_%s", testName, allFinders[i].name), func(t *testing.T) {
				if testing.Short() {
					if shouldSkip, ok := shortSkip[allFinders[i].name]; ok && shouldSkip {
						t.SkipNow()
					}
				}
				t.Parallel()
				// TODO: figure out better solution to handling testdata
				// currently getting the package associate testdata returns a string
				// beginning with '_', which throughs of the later 'go test' calls
				tester := &Tester{
					testPos: position{
						file: test.fileName,
						pkg:  fmt.Sprintf("../testdata/%s", test.fileDir),
						line: test.line,
						col:  test.col,
					},
					finder: allFinders[i].newFinder(),
				}

				coveredBy, err := tester.CoveredBy()
				if test.expectErr {
					if err == nil {
						t.Errorf("Unexpectedly no error")
					}
				} else {
					if err != nil {
						if exitErr, ok := err.(*exec.ExitError); ok {
							t.Errorf("Unexpected error checking for covering tests: %s", exitErr.Stderr)
						} else {
							t.Errorf("Unexpected error checking for covering tests: %#v", err)
						}
						return
					}
				}

				if len(coveredBy) != len(test.expectCoveredBy) {
					t.Errorf("Unexpected CoveredBy (expected = %v, actual = %v)", test.expectCoveredBy, coveredBy)
				}

				sort.Slice(coveredBy, func(i, j int) bool { return coveredBy[i] < coveredBy[j] })
				expectCoveredBy := []string{}
				expectCoveredBy = append(expectCoveredBy, test.expectCoveredBy...)
				sort.Slice(expectCoveredBy, func(i, j int) bool { return expectCoveredBy[i] < expectCoveredBy[j] })
				for i := range coveredBy {
					if coveredBy[i] != expectCoveredBy[i] {
						t.Errorf("Unexpected CoveredBy[%d] (expected = %s, actual = %s)", i, expectCoveredBy[i], coveredBy[i])
					}
				}
			})
		}
	}
}

var coveringTests []string

func BenchmarkCoveredBy(b *testing.B) {
	for testName, test := range coveredByTests {
		for i := range allFinders {
			// TODO: figure out better solution to handling testdata
			// currently getting the package associate testdata returns a string
			// beginning with '_', which throughs of the later 'go test' calls
			tester := &Tester{
				testPos: position{
					file: test.fileName,
					pkg:  fmt.Sprintf("../testdata/%s", test.fileDir),
					line: test.line,
					col:  test.col,
				},
				finder: allFinders[i].newFinder(),
			}
			b.Run(fmt.Sprintf("%s_%s", testName, allFinders[i].name), func(b *testing.B) {
				if testing.Short() {
					if shouldSkip, ok := shortSkip[allFinders[i].name]; ok && shouldSkip {
						b.SkipNow()
					}
				}
				var (
					covered []string
					err     error
				)
				for n := 0; n < b.N; n++ {
					covered, err = tester.CoveredBy()
					if err != nil {
						if !test.expectErr {
							b.Errorf("Unexpected error: %s", err)
						}
					}
				}
				coveringTests = covered
			})
		}
	}
}
