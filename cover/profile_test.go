package cover

import (
	"bytes"
	"testing"
)

const coverOut = `mode: set
fmt/errors.go:17.52,23.25 6 1
fmt/errors.go:28.2,29.12 2 1
fmt/errors.go:23.25,25.3 1 1
fmt/errors.go:25.8,27.3 1 1
fmt/errors.go:37.36,39.2 1 1
fmt/errors.go:41.36,43.2 1 1
fmt/format.go:54.28,56.2 1 1
fmt/format.go:58.33,61.2 2 1
fmt/format.go:64.35,65.12 1 1
fmt/format.go:68.2,72.23 4 1
fmt/format.go:77.2,78.12 2 1
fmt/format.go:82.2,83.25 2 1
fmt/format.go:86.2,86.23 1 1
fmt/format.go:65.12,67.3 1 1
fmt/format.go:72.23,75.3 2 0
fmt/format.go:78.12,80.3 1 1
fmt/format.go:83.25,85.3 1 1`

var coverTests = map[string]struct {
	cover         string
	file          string
	line, col     int
	expectCovered bool
}{
	"start_of_coverage": {
		cover: coverOut,
		file:  "errors.go",
		line:  17, col: 52,
		expectCovered: true,
	},
	"end_of_coverage": {
		cover: coverOut,
		file:  "errors.go",
		line:  23, col: 25,
		expectCovered: true,
	},
	"in_uncovered_block": {
		cover: coverOut,
		file:  "format.go",
		line:  73, col: 0,
		expectCovered: false,
	},
	"middle_of_covered_line": {
		cover: coverOut,
		file:  "format.go",
		line:  86, col: 4,
		expectCovered: true,
	},
	"uncovered_file": {
		cover: coverOut,
		file:  "fake_file.go",
		line:  86, col: 4,
		expectCovered: false,
	},
}

func TestCovers(t *testing.T) {
	for testName, test := range coverTests {
		t.Run(testName, func(t *testing.T) {
			var b bytes.Buffer
			b.WriteString(test.cover)

			profile, err := New(&b)
			if err != nil {
				t.Fatalf("Error creating profile: %s", err)
			}

			covered := profile.Covers(test.file, test.line, test.col)
			if covered != test.expectCovered {
				t.Errorf("Unexpected coverage result (expected = %v, actual = %v)", test.expectCovered, covered)
			}
		})
	}
}

var parseLineTests = map[string]struct {
	line           string
	expectErr      bool
	expectedParsed coverLine
}{
	"covered_line": {
		line:      "fmt/scan.go:737.2,737.22 1 1",
		expectErr: false,
		expectedParsed: coverLine{
			pkg:  "fmt",
			file: "scan.go",
			coverBlock: coverBlock{
				startLine: 737,
				startCol:  2,
				endLine:   737,
				endCol:    22,
				numStmt:   1,
				count:     1,
			},
		},
	},
	"uncovered_line": {
		line:      "fmt/format.go:72.23,75.3 2 0",
		expectErr: false,
		expectedParsed: coverLine{
			pkg:  "fmt",
			file: "format.go",
			coverBlock: coverBlock{
				startLine: 72,
				startCol:  23,
				endLine:   75,
				endCol:    3,
				numStmt:   2,
				count:     0,
			},
		},
	},
	"mode_line": {
		line:      "mode: set",
		expectErr: true,
	},
}

func TestParseLine(t *testing.T) {
	for testName, test := range parseLineTests {
		t.Run(testName, func(t *testing.T) {
			parsed, err := parseLine(test.line)
			if test.expectErr && err == nil {
				t.Errorf("Unexpectedly no error")
			}
			if !test.expectErr && err != nil {
				t.Errorf("Unexpected error: %s", err)
			}

			if parsed != test.expectedParsed {
				t.Errorf("Unexpected parsed line (expected = %#v, actual = %#v)", test.expectedParsed, parsed)
			}
		})
	}
}
