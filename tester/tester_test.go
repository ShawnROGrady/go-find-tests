package tester

import (
	"fmt"
	"os/exec"
	"sort"
	"testing"
)

var coveredByTests = map[string]struct {
	fileDir         string
	fileName        string
	line, col       int
	expectCoveredBy []string
	expectErr       bool
}{
	"covered_by_one_test": {
		fileDir:  "size",
		fileName: "size.go",
		line:     22, col: 0, // body of isEnormous()
		expectCoveredBy: []string{"TestIsEnormous"},
	},
	"covered_by_three_tests": {
		fileDir:  "size",
		fileName: "size.go",
		line:     8, col: 0, // negative case of size()
		expectCoveredBy: []string{"TestSize", "TestNegativeSize", "TestIsNegative"},
	},
	"not_covered": {
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
}

func TestCoveredBy(t *testing.T) {
	for testName, test := range coveredByTests {
		t.Run(testName, func(t *testing.T) {
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
				finder: errGroupFinder{},
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
			sort.Slice(test.expectCoveredBy, func(i, j int) bool { return test.expectCoveredBy[i] < test.expectCoveredBy[j] })
			for i := range coveredBy {
				if coveredBy[i] != test.expectCoveredBy[i] {
					t.Errorf("Unexpected CoveredBy[%d] (expected = %s, actual = %s)", i, test.expectCoveredBy[i], coveredBy[i])
				}
			}
		})
	}
}

var coveringTests []string

func BenchmarkCoveredBy(b *testing.B) {
	for testName, test := range coveredByTests {
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
			finder: errGroupFinder{},
		}
		b.Run(testName, func(b *testing.B) {
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
