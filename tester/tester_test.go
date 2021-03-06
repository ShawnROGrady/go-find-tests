package tester

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"testing"
)

var allFinders = map[string]func() coverFinder{
	"sequential": func() coverFinder { return sequentialFinder{} },
	"err_group":  func() coverFinder { return errGroupFinder{} },
}

var coveredByTests = map[string]struct {
	fileDir         string
	fileName        string
	includeSubtests bool
	short           bool
	runExpr         string
	line, col       int
	expectCoveredBy []string
	expectErr       bool
	expectedErr     error
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
	"test_skipped_due_to_short_flag": {
		fileDir:  "size",
		fileName: "size.go",
		line:     8, col: 0, // negative case of size()
		short:           true,
		expectCoveredBy: []string{"TestSize", "TestIsNegative"},
	},
	"runexp_matches_covering_test": {
		fileDir:  "size",
		fileName: "size.go",
		line:     8, col: 0, // negative case of size()
		runExpr:         "TestIs",
		expectCoveredBy: []string{"TestIsNegative"},
	},
	"runexp_only_match_not_covering_test": {
		fileDir:  "size",
		fileName: "size.go",
		line:     8, col: 0, // negative case of size()
		runExpr:         "TestIsEnormous",
		expectCoveredBy: []string{},
	},
	"runexp_no_matches": {
		fileDir:  "size",
		fileName: "size.go",
		line:     8, col: 0, // negative case of size()
		runExpr:         "SomethingInvalid",
		expectCoveredBy: []string{},
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
		expectedErr: fmt.Errorf("error running test 'TestSum': %s", &testErr{
			testName: "TestSum",
			output:   "fail_test.go:12: Unexpected sum(1, 2) (expected = 3, actual = 2)",
		}),
	},
	"subtests_enabled_covered_by_subtests": {
		fileDir:  "subtests",
		fileName: "len.go",
		line:     9, col: 0, // "empty" case of length()
		includeSubtests: true,
		expectCoveredBy: []string{
			"TestIsEmpty",
			"TestIsEmpty/empty_input",
			"TestIsShort",
			"TestIsShort/empty_input",
		},
	},
	"subtests_disabled_covered_by_subtests": {
		fileDir:  "subtests",
		fileName: "len.go",
		line:     9, col: 0, // "empty" case of length()
		includeSubtests: false,
		expectCoveredBy: []string{
			"TestIsEmpty",
			"TestIsShort",
		},
	},
	"subtests_enabled_no_subtests": {
		fileDir:         "size",
		fileName:        "size.go",
		includeSubtests: true,
		line:            22, col: 0, // body of isEnormous()
		expectCoveredBy: []string{"TestIsEnormous"},
	},
}

func TestCoveredBy(t *testing.T) {
	for finderName, newFinder := range allFinders {
		t.Run(fmt.Sprintf("finder=%s", finderName), func(t *testing.T) {
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
						includeSubtests: test.includeSubtests,
						short:           test.short,
						run:             test.runExpr,
						coverFinder:     newFinder(),
					}

					// This logic is normally handled in the constructor
					if tester.run == "" {
						tester.run = "."
					}

					coveredBy, err := tester.CoveredBy()
					if test.expectErr {
						if err == nil {
							t.Errorf("Unexpectedly no error")
						}
						if test.expectedErr != nil {
							if test.expectedErr.Error() != err.Error() {
								t.Errorf("Unexpected error message (expected = '%s', actual = '%s')", test.expectedErr, err)
							}
						}
					} else {
						if err != nil {
							t.Errorf("Unexpected error checking for covering tests: %s", err)
							return
						}
					}

					if len(coveredBy) != len(test.expectCoveredBy) {
						t.Errorf("Unexpected CoveredBy (expected = %v, actual = %v)", test.expectCoveredBy, coveredBy)
						return
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
		})
	}
}

var coveringTests []string

var coveredByBenchmarks = map[string]struct {
	fileDir         string
	fileName        string
	includeSubtests bool
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
	"subtests_enabled_covered_by_subtests": {
		fileDir:  "subtests",
		fileName: "len.go",
		line:     9, col: 0, // "empty" case of length()
		includeSubtests: true,
		expectCoveredBy: []string{
			"TestIsEmpty",
			"TestIsEmpty/empty_input",
			"TestIsShort",
			"TestIsShort/empty_input",
		},
	},
	"subtests_disabled_covered_by_subtests": {
		fileDir:  "subtests",
		fileName: "len.go",
		line:     9, col: 0, // "empty" case of length()
		includeSubtests: false,
		expectCoveredBy: []string{
			"TestIsEmpty",
			"TestIsShort",
		},
	},
	"subtests_enabled_no_subtests": {
		fileDir:         "size",
		fileName:        "size.go",
		includeSubtests: true,
		line:            22, col: 0, // body of isEnormous()
		expectCoveredBy: []string{"TestIsEnormous"},
	},
	"subtests_enabled_covered_by_3testsX2subtests": {
		fileDir:  "subtests",
		fileName: "len.go",
		line:     6, col: 0, // first line of length()
		includeSubtests: true,
		expectCoveredBy: []string{
			"TestIsEmpty", "TestIsEmpty/empty_input", "TestIsEmpty/short_input", "TestIsEmpty/long_input",
			"TestIsShort", "TestIsShort/empty_input", "TestIsShort/short_input", "TestIsShort/long_input",
		},
	},
	"subtests_enabled_covered_by_5testsX5subtests": {
		fileDir:  "subtests_5_5",
		fileName: "len.go",
		line:     7, col: 0, // first line of length()
		includeSubtests: true,
		expectCoveredBy: []string{
			"TestIsEmpty", "TestIsEmpty/empty_input", "TestIsEmpty/short_input", "TestIsEmpty/long_input", "TestIsEmpty/very_long_input", "TestIsEmpty/novel_input",
			"TestIsShort", "TestIsShort/empty_input", "TestIsShort/short_input", "TestIsShort/long_input", "TestIsShort/very_long_input", "TestIsShort/novel_input",
			"TestIsLong", "TestIsLong/empty_input", "TestIsLong/short_input", "TestIsLong/long_input", "TestIsLong/very_long_input", "TestIsLong/novel_input",
			"TestIsVeryLong", "TestIsVeryLong/empty_input", "TestIsVeryLong/short_input", "TestIsVeryLong/long_input", "TestIsVeryLong/very_long_input", "TestIsVeryLong/novel_input",
			"TestIsNovel", "TestIsNovel/empty_input", "TestIsNovel/short_input", "TestIsNovel/long_input", "TestIsNovel/very_long_input", "TestIsNovel/novel_input",
		},
	},
}

func BenchmarkCoveredBy(b *testing.B) {
	for finderName, newFinder := range allFinders {
		b.Run(fmt.Sprintf("finder=%s", finderName), func(b *testing.B) {
			for testName, test := range coveredByBenchmarks {
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
					includeSubtests: test.includeSubtests,
					coverFinder:     newFinder(),
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
		})
	}
}

func BenchmarkCoveringTests(b *testing.B) {
	for finderName, newFinder := range allFinders {
		b.Run(fmt.Sprintf("finder=%s", finderName), func(b *testing.B) {
			for testName, test := range coveredByBenchmarks {
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
					coverFinder:     newFinder(),
					includeSubtests: test.includeSubtests,
				}
				outputDir, err := ioutil.TempDir("", "test_finder")
				if err != nil {
					b.Fatalf("Error creating tmp dir: %s", err)
				}
				testBin, err := tester.compileTest(outputDir)
				if err != nil {
					if !test.expectErr {
						b.Errorf("Error compiling test: %s", err)
					}
					os.RemoveAll(outputDir)
					continue
				}
				allTests, err := findTests(tester.testPos.pkg, ".")
				if err != nil {
					if !test.expectErr {
						b.Errorf("Error finding tests: %s", err)
					}
					os.RemoveAll(outputDir)
					continue
				}
				b.Run(testName, func(b *testing.B) {
					var (
						covered []string
						err     error
					)
					for n := 0; n < b.N; n++ {
						covered, err = tester.coverFinder.coveringTests(tester, testBin, outputDir, allTests, tester.includeSubtests)
						if err != nil {
							if !test.expectErr {
								b.Errorf("Unexpected error: %s", err)
							}
						}
					}
					coveringTests = covered
				})

				os.RemoveAll(outputDir)
			}
		})
	}
}
