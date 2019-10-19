package subtests

import "testing"

var isEmptyTests = map[string]struct {
	input       string
	expectEmpty bool
}{
	"empty_input": {
		input:       "",
		expectEmpty: true,
	},
	"short_input": {
		input:       "hello",
		expectEmpty: false,
	},
	"long_input": {
		input:       "hello world!",
		expectEmpty: false,
	},
}

func TestIsEmpty(t *testing.T) {
	for testName, testCase := range isEmptyTests {
		t.Run(testName, func(t *testing.T) {
			empty := isEmpty(testCase.input)
			if empty != testCase.expectEmpty {
				t.Errorf("unexpected isEmpty('%s') [expected = %v, actual = %v]", testCase.input, testCase.expectEmpty, empty)
			}
		})
	}
}

var isShortTests = map[string]struct {
	input       string
	expectShort bool
}{
	"empty_input": {
		input:       "",
		expectShort: false,
	},
	"short_input": {
		input:       "hello",
		expectShort: true,
	},
	"long_input": {
		input:       "hello world!",
		expectShort: false,
	},
}

func TestIsShort(t *testing.T) {
	for testName, testCase := range isShortTests {
		t.Run(testName, func(t *testing.T) {
			empty := isShort(testCase.input)
			if empty != testCase.expectShort {
				t.Errorf("unexpected isShort('%s') [expected = %v, actual = %v]", testCase.input, testCase.expectShort, empty)
			}
		})
	}
}
