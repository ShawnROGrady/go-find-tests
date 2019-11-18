package main

import (
	"bytes"
	"testing"

	"github.com/ShawnROGrady/go-find-tests/tester"
)

/*
NOTE: since tester.New() doesn't play well with testdata directory
these tests all check coverage of actual implementation. Since this
is likely to change these test will need to be updated accordingly
*/
var runTests = map[string]struct {
	conf           runConfig
	path           string
	line           int
	col            int
	expectErr      bool
	expectedOutput string
}{
	"default_options": {
		conf: runConfig{
			lineFmt: defaultLineFmt,
		},
		path: "../../cover/profile.go",
		line: 154, col: 12, // success case parseLine()
		expectErr:      false,
		expectedOutput: "TestCovers\nTestParseLine\n",
	},
	"run_filter_set": {
		conf: runConfig{
			lineFmt: defaultLineFmt,
			testerConf: tester.Config{
				Run: "TestCovers",
			},
		},
		path: "../../cover/profile.go",
		line: 154, col: 12, // success case parseLine()
		expectErr:      false,
		expectedOutput: "TestCovers\n",
	},
	"json_printing_no_subs": {
		conf: runConfig{
			lineFmt: defaultLineFmt,
			jsonFmt: true,
		},
		path: "../../cover/profile.go",
		line: 154, col: 12, // success case parseLine()
		expectErr:      false,
		expectedOutput: `["TestCovers","TestParseLine"]`,
	},
	"subs_enabled": {
		conf: runConfig{
			lineFmt: defaultLineFmt,
			testerConf: tester.Config{
				IncludeSubtests: true,
			},
		},
		path: "../../cover/profile.go",
		line: 154, col: 12, // success case parseLine()
		expectErr:      false,
		expectedOutput: "TestCovers\nTestCovers/end_of_coverage\nTestCovers/in_uncovered_block\nTestCovers/middle_of_covered_line\nTestCovers/start_of_coverage\nTestCovers/uncovered_file\nTestParseLine\nTestParseLine/covered_line\nTestParseLine/uncovered_line\n",
	},
	"json_printing_subs_enabled": {
		conf: runConfig{
			lineFmt: defaultLineFmt,
			testerConf: tester.Config{
				IncludeSubtests: true,
			},
			jsonFmt: true,
		},
		path: "../../cover/profile.go",
		line: 154, col: 12, // success case parseLine()
		expectErr:      false,
		expectedOutput: `["TestCovers","TestCovers/end_of_coverage","TestCovers/in_uncovered_block","TestCovers/middle_of_covered_line","TestCovers/start_of_coverage","TestCovers/uncovered_file","TestParseLine","TestParseLine/covered_line","TestParseLine/uncovered_line"]`,
	},
	"with_positions": {
		conf: runConfig{
			lineFmt:        defaultLineFmt,
			printPositions: true,
		},
		path: "../../cover/profile.go",
		line: 154, col: 12, // success case parseLine()
		expectErr:      false,
		expectedOutput: "TestCovers:../../cover/profile_test.go:65:1:\nTestParseLine:../../cover/profile_test.go:127:1:\n",
	},
	"json_printing_with_positions": {
		conf: runConfig{
			lineFmt:        defaultLineFmt,
			printPositions: true,
			jsonFmt:        true,
		},
		path: "../../cover/profile.go",
		line: 154, col: 12, // success case parseLine()
		expectErr:      false,
		expectedOutput: `{"TestCovers":{"file":"../../cover/profile_test.go","line":65,"col":1,"offset":1270},"TestParseLine":{"file":"../../cover/profile_test.go","line":127,"col":1,"offset":2543}}`,
	},
	"json_printing_with_positions_and_subs": {
		conf: runConfig{
			lineFmt:        defaultLineFmt,
			printPositions: true,
			jsonFmt:        true,
			testerConf: tester.Config{
				IncludeSubtests: true,
			},
		},
		path: "../../cover/profile.go",
		line: 154, col: 12, // success case parseLine()
		expectErr:      false,
		expectedOutput: `{"TestCovers":{"file":"../../cover/profile_test.go","line":65,"col":1,"offset":1270,"subtests":["TestCovers/end_of_coverage","TestCovers/in_uncovered_block","TestCovers/middle_of_covered_line","TestCovers/start_of_coverage","TestCovers/uncovered_file"]},"TestParseLine":{"file":"../../cover/profile_test.go","line":127,"col":1,"offset":2543,"subtests":["TestParseLine/covered_line","TestParseLine/uncovered_line"]}}`,
	},
	"with_positions_subs_enabled": {
		conf: runConfig{
			lineFmt:        defaultLineFmt,
			printPositions: true,
			testerConf: tester.Config{
				IncludeSubtests: true,
			},
		},
		path: "../../cover/profile.go",
		line: 154, col: 12, // success case parseLine()
		expectErr:      false,
		expectedOutput: "TestCovers:../../cover/profile_test.go:65:1:TestCovers/end_of_coverage,TestCovers/in_uncovered_block,TestCovers/middle_of_covered_line,TestCovers/start_of_coverage,TestCovers/uncovered_file\nTestParseLine:../../cover/profile_test.go:127:1:TestParseLine/covered_line,TestParseLine/uncovered_line\n",
	},
}

func TestRun(t *testing.T) {
	for testName, testCase := range runTests {
		t.Run(testName, func(t *testing.T) {
			var b bytes.Buffer

			err := run(testCase.conf, testCase.path, testCase.line, testCase.col, &b)
			if err != nil {
				if !testCase.expectErr {
					t.Errorf("Unexpected error: %s", err)
				}
				return
			}

			if testCase.expectErr {
				t.Error("Unexpectedly no error")
				return
			}

			actual := b.String()
			if actual != testCase.expectedOutput {
				t.Errorf("Unexpected output (expected = '%s', actual = '%s')", testCase.expectedOutput, actual)
			}
		})
	}
}
