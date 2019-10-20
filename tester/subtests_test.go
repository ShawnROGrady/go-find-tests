package tester

import (
	"bytes"
	"testing"
)

var subtestsTests = map[string]struct {
	testOutput       string
	expectedSubtests []string
}{
	"3_subtests": {
		testOutput: `=== RUN   TestIsEmpty
=== RUN   TestIsEmpty/empty_input
=== RUN   TestIsEmpty/short_input
=== RUN   TestIsEmpty/long_input
--- PASS: TestIsEmpty (0.00s)
    --- PASS: TestIsEmpty/empty_input (0.00s)
    --- PASS: TestIsEmpty/short_input (0.00s)
    --- PASS: TestIsEmpty/long_input (0.00s)
PASS
coverage: 66.7% of statements
ok      github.com/ShawnROGrady/go-find-tests/testdata/subtests 0.007s  coverage: 66.7% of statements
`,
		expectedSubtests: []string{"TestIsEmpty/empty_input", "TestIsEmpty/short_input", "TestIsEmpty/long_input"},
	},
	"no_subtests": {
		testOutput: `=== RUN   TestEmptyStringIsEmpty
--- PASS: TestEmptyStringIsEmpty (0.00s)
PASS
coverage: 44.4% of statements
ok      github.com/ShawnROGrady/go-find-tests/testdata/len10    0.014s  coverage: 44.4% of statements
`,
		expectedSubtests: []string{},
	},
}

func TestSubtests(t *testing.T) {
	for testName, testCase := range subtestsTests {
		t.Run(testName, func(t *testing.T) {
			var b bytes.Buffer
			b.WriteString(testCase.testOutput)

			subtests := subtests(&b)

			if len(subtests) != len(testCase.expectedSubtests) {
				t.Fatalf("Unexpected subtests [expected = %q, actual = %q]", testCase.expectedSubtests, subtests)
			}

			for i := range subtests {
				if subtests[i] != testCase.expectedSubtests[i] {
					t.Errorf("Unexpected subtests[%d] [expected = '%s', actual ='%s']", i, testCase.expectedSubtests[i], subtests[i])
				}
			}
		})
	}
}
