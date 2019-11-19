package main

import "testing"

var parsePositionTests = map[string]struct {
	providedArg string
	expectedPos *pos
	expectErr   bool
}{
	"valid_arg_w_col": {
		providedArg: "./cover/profile.go:154.11",
		expectedPos: &pos{
			file: "./cover/profile.go",
			line: 154,
			col:  11,
		},
	},
	"valid_arg_no_col": {
		providedArg: "./cover/profile.go:154",
		expectedPos: &pos{
			file: "./cover/profile.go",
			line: 154,
			col:  0,
		},
	},
	"no_line_or_col": {
		providedArg: "./cover/profile.go",
		expectedPos: &pos{
			file: "./cover/profile.go",
			line: 154,
			col:  0,
		},
	},
	"std_lib_file": {
		providedArg: "fmt/errors.go:17.52",
		expectedPos: &pos{
			file: "fmt/errors.go",
			line: 17,
			col:  52,
		},
	},
}

func TestParsePosition(t *testing.T) {
	for testName, testCase := range parsePositionTests {
		t.Run(testName, func(t *testing.T) {
			pos, err := parsePosition(testCase.providedArg)
			if err != nil {
				if testCase.expectErr {
					t.Errorf("Unexpected error: %s", err)
				}
				return
			}

			if testCase.expectErr {
				t.Error("Unexpectedly no error")
				return
			}

			if testCase.expectedPos == nil && pos != nil {
				t.Errorf("parsed position unexpectedly not nil: %v", pos)
				return
			}

			if testCase.expectedPos != nil && pos == nil {
				t.Errorf("parsed position unexpectedly nil")
				return
			}

			if *testCase.expectedPos != *pos {
				t.Errorf("Unexpected parsed postion (expected = %#v, actual = %#v)", testCase.expectedPos, pos)
			}
		})
	}
}
