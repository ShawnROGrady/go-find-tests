package main

import (
	"testing"

	"github.com/ShawnROGrady/go-find-tests/finder"
)

var (
	coveringTestName = "TestPackageTests"
	coveringTestPos  = finder.TestPosition{
		File:   "finder/finder_test.go",
		Line:   79,
		Col:    1,
		Offset: 1580,
	}
)

var fmtPositionTests = map[string]string{
	"%t:%f:%l:%c": "TestPackageTests:finder/finder_test.go:79:1",
	"%f:%t":       "finder/finder_test.go:TestPackageTests",
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
