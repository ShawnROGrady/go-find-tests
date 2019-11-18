package main

import (
	"testing"

	"github.com/ShawnROGrady/go-find-tests/finder"
)

var (
	coveringTestName = "TestPackageTests"
	coveringTestPos  = testPosition{
		TestPosition: finder.TestPosition{
			File:   "finder/finder_test.go",
			Line:   79,
			Col:    1,
			Offset: 1580,
		},
		SubTests: []string{"TestPackageTests/10_tests_1_file", "TestPackageTests/20_tests_2_files"},
	}
)

var fmtPositionTests = map[string]string{
	"%t:%f:%l:%c":    "TestPackageTests:finder/finder_test.go:79:1",
	"%f:%t":          "finder/finder_test.go:TestPackageTests",
	"%t:%f:%l:%c:%s": "TestPackageTests:finder/finder_test.go:79:1:TestPackageTests/10_tests_1_file,TestPackageTests/20_tests_2_files",
}

func TestFmtPosition(t *testing.T) {
	for outputFmt, expectedOutput := range fmtPositionTests {
		t.Run(outputFmt, func(t *testing.T) {
			output := fmtPosition(coveringTestPos, coveringTestName, outputFmt)
			if output != expectedOutput {
				t.Errorf("Unexpected output (expected = '%s', actual = '%s')", expectedOutput, output)
			}
		})
	}
}
