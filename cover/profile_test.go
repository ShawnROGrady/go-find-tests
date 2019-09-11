package cover

import "testing"

var parseLineTests = map[string]struct {
	line           string
	expectErr      bool
	expectedParsed coverLine
}{
	"covered_line": {
		line:      "fmt/scan.go:737.2,737.22 1 1",
		expectErr: false,
		expectedParsed: coverLine{
			pkg:       "fmt",
			file:      "scan.go",
			startLine: 737,
			startCol:  2,
			endLine:   737,
			endCol:    22,
			numStmt:   1,
			count:     1,
		},
	},
	"uncovered_line": {
		line:      "fmt/format.go:72.23,75.3 2 0",
		expectErr: false,
		expectedParsed: coverLine{
			pkg:       "fmt",
			file:      "format.go",
			startLine: 72,
			startCol:  23,
			endLine:   75,
			endCol:    3,
			numStmt:   2,
			count:     0,
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
