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
		testOutput: `{"Time":"2019-11-10T18:10:33.350927-06:00","Action":"run","Package":"github.com/ShawnROGrady/go-find-tests/testdata/subtests","Test":"TestIsEmpty"}
{"Time":"2019-11-10T18:10:33.351288-06:00","Action":"output","Package":"github.com/ShawnROGrady/go-find-tests/testdata/subtests","Test":"TestIsEmpty","Output":"=== RUN   TestIsEmpty\n"}
{"Time":"2019-11-10T18:10:33.351315-06:00","Action":"run","Package":"github.com/ShawnROGrady/go-find-tests/testdata/subtests","Test":"TestIsEmpty/empty_input"}
{"Time":"2019-11-10T18:10:33.351324-06:00","Action":"output","Package":"github.com/ShawnROGrady/go-find-tests/testdata/subtests","Test":"TestIsEmpty/empty_input","Output":"=== RUN TestIsEmpty/empty_input\n"}
{"Time":"2019-11-10T18:10:33.351334-06:00","Action":"run","Package":"github.com/ShawnROGrady/go-find-tests/testdata/subtests","Test":"TestIsEmpty/short_input"}
{"Time":"2019-11-10T18:10:33.351341-06:00","Action":"output","Package":"github.com/ShawnROGrady/go-find-tests/testdata/subtests","Test":"TestIsEmpty/short_input","Output":"=== RUN TestIsEmpty/short_input\n"}
{"Time":"2019-11-10T18:10:33.351348-06:00","Action":"run","Package":"github.com/ShawnROGrady/go-find-tests/testdata/subtests","Test":"TestIsEmpty/long_input"}
{"Time":"2019-11-10T18:10:33.351354-06:00","Action":"output","Package":"github.com/ShawnROGrady/go-find-tests/testdata/subtests","Test":"TestIsEmpty/long_input","Output":"=== RUN  TestIsEmpty/long_input\n"}
{"Time":"2019-11-10T18:10:33.351375-06:00","Action":"output","Package":"github.com/ShawnROGrady/go-find-tests/testdata/subtests","Test":"TestIsEmpty","Output":"--- PASS: TestIsEmpty (0.00s)\n"}
{"Time":"2019-11-10T18:10:33.351385-06:00","Action":"output","Package":"github.com/ShawnROGrady/go-find-tests/testdata/subtests","Test":"TestIsEmpty/empty_input","Output":"    --- PASS: TestIsEmpty/empty_input (0.00s)\n"}
{"Time":"2019-11-10T18:10:33.351523-06:00","Action":"pass","Package":"github.com/ShawnROGrady/go-find-tests/testdata/subtests","Test":"TestIsEmpty/empty_input","Elapsed":0}
{"Time":"2019-11-10T18:10:33.35155-06:00","Action":"output","Package":"github.com/ShawnROGrady/go-find-tests/testdata/subtests","Test":"TestIsEmpty/short_input","Output":"    --- PASS: TestIsEmpty/short_input (0.00s)\n"}
{"Time":"2019-11-10T18:10:33.351559-06:00","Action":"pass","Package":"github.com/ShawnROGrady/go-find-tests/testdata/subtests","Test":"TestIsEmpty/short_input","Elapsed":0}
{"Time":"2019-11-10T18:10:33.351566-06:00","Action":"output","Package":"github.com/ShawnROGrady/go-find-tests/testdata/subtests","Test":"TestIsEmpty/long_input","Output":"    --- PASS: TestIsEmpty/long_input (0.00s)\n"}
{"Time":"2019-11-10T18:10:33.351572-06:00","Action":"pass","Package":"github.com/ShawnROGrady/go-find-tests/testdata/subtests","Test":"TestIsEmpty/long_input","Elapsed":0}
{"Time":"2019-11-10T18:10:33.353468-06:00","Action":"pass","Package":"github.com/ShawnROGrady/go-find-tests/testdata/subtests","Test":"TestIsEmpty","Elapsed":0}
{"Time":"2019-11-10T18:10:33.353526-06:00","Action":"output","Package":"github.com/ShawnROGrady/go-find-tests/testdata/subtests","Output":"PASS\n"}
{"Time":"2019-11-10T18:10:33.353601-06:00","Action":"output","Package":"github.com/ShawnROGrady/go-find-tests/testdata/subtests","Output":"ok  \tgithub.com/ShawnROGrady/go-find-tests/testdata/subtests\t0.008s\n"}
{"Time":"2019-11-10T18:10:33.353614-06:00","Action":"pass","Package":"github.com/ShawnROGrady/go-find-tests/testdata/subtests","Elapsed":0.008}
`,
		expectedSubtests: []string{"TestIsEmpty/empty_input", "TestIsEmpty/short_input", "TestIsEmpty/long_input"},
	},
	"no_subtests": {
		testOutput: `{"Time":"2019-11-10T18:09:29.637515-06:00","Action":"run","Package":"github.com/ShawnROGrady/go-find-tests/testdata/len10","Test":"TestEmptyStringIsEmpty"}
{"Time":"2019-11-10T18:09:29.637824-06:00","Action":"output","Package":"github.com/ShawnROGrady/go-find-tests/testdata/len10","Test":"TestEmptyStringIsEmpty","Output":"=== RUN   TestEmptyStringIsEmpty\n"}
{"Time":"2019-11-10T18:09:29.637856-06:00","Action":"output","Package":"github.com/ShawnROGrady/go-find-tests/testdata/len10","Test":"TestEmptyStringIsEmpty","Output":"--- PASS: TestEmptyStringIsEmpty (0.00s)\n"}
{"Time":"2019-11-10T18:09:29.637868-06:00","Action":"pass","Package":"github.com/ShawnROGrady/go-find-tests/testdata/len10","Test":"TestEmptyStringIsEmpty","Elapsed":0}
{"Time":"2019-11-10T18:09:29.637891-06:00","Action":"output","Package":"github.com/ShawnROGrady/go-find-tests/testdata/len10","Output":"PASS\n"}
{"Time":"2019-11-10T18:09:29.637954-06:00","Action":"output","Package":"github.com/ShawnROGrady/go-find-tests/testdata/len10","Output":"ok  \tgithub.com/ShawnROGrady/go-find-tests/testdata/len10\t0.006s\n"}
{"Time":"2019-11-10T18:09:29.637968-06:00","Action":"pass","Package":"github.com/ShawnROGrady/go-find-tests/testdata/len10","Elapsed":0.006}
`,
		expectedSubtests: []string{},
	},
}

func TestSubtests(t *testing.T) {
	for testName, testCase := range subtestsTests {
		t.Run(testName, func(t *testing.T) {
			var b bytes.Buffer
			b.WriteString(testCase.testOutput)

			subtests, err := subtests(&b)
			if err != nil {
				t.Fatalf("Unexpected error parsing subtests: %s", err)
			}

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
