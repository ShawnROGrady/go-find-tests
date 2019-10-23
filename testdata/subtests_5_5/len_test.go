package subtests

import (
	"strings"
	"testing"
)

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
	"very_long_input": {
		input:       strings.Repeat("a", 500),
		expectEmpty: false,
	},
	"novel_input": {
		input:       strings.Repeat("a", 5000),
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
	"very_long_input": {
		input:       strings.Repeat("a", 500),
		expectShort: false,
	},
	"novel_input": {
		input:       strings.Repeat("a", 5000),
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

var isLongTests = map[string]struct {
	input      string
	expectLong bool
}{
	"empty_input": {
		input:      "",
		expectLong: false,
	},
	"short_input": {
		input:      "hello",
		expectLong: false,
	},
	"long_input": {
		input:      "hello world!",
		expectLong: true,
	},
	"very_long_input": {
		input:      strings.Repeat("a", 500),
		expectLong: false,
	},
	"novel_input": {
		input:      strings.Repeat("a", 5000),
		expectLong: false,
	},
}

func TestIsLong(t *testing.T) {
	for testName, testCase := range isLongTests {
		t.Run(testName, func(t *testing.T) {
			empty := isLong(testCase.input)
			if empty != testCase.expectLong {
				t.Errorf("unexpected isLong('%s') [expected = %v, actual = %v]", testCase.input, testCase.expectLong, empty)
			}
		})
	}
}

var isVeryLongTests = map[string]struct {
	input          string
	expectVeryLong bool
}{
	"empty_input": {
		input:          "",
		expectVeryLong: false,
	},
	"short_input": {
		input:          "hello",
		expectVeryLong: false,
	},
	"long_input": {
		input:          "hello world!",
		expectVeryLong: false,
	},
	"very_long_input": {
		input:          strings.Repeat("a", 500),
		expectVeryLong: true,
	},
	"novel_input": {
		input:          strings.Repeat("a", 5000),
		expectVeryLong: false,
	},
}

func TestIsVeryLong(t *testing.T) {
	for testName, testCase := range isVeryLongTests {
		t.Run(testName, func(t *testing.T) {
			empty := isVeryLong(testCase.input)
			if empty != testCase.expectVeryLong {
				t.Errorf("unexpected isVeryLong('%s') [expected = %v, actual = %v]", testCase.input, testCase.expectVeryLong, empty)
			}
		})
	}
}

var isNovelTests = map[string]struct {
	input       string
	expectNovel bool
}{
	"empty_input": {
		input:       "",
		expectNovel: false,
	},
	"short_input": {
		input:       "hello",
		expectNovel: false,
	},
	"long_input": {
		input:       "hello world!",
		expectNovel: false,
	},
	"very_long_input": {
		input:       strings.Repeat("a", 500),
		expectNovel: false,
	},
	"novel_input": {
		input:       strings.Repeat("a", 5000),
		expectNovel: true,
	},
}

func TestIsNovel(t *testing.T) {
	for testName, testCase := range isNovelTests {
		t.Run(testName, func(t *testing.T) {
			empty := isNovel(testCase.input)
			if empty != testCase.expectNovel {
				t.Errorf("unexpected isNovel('%s') [expected = %v, actual = %v]", testCase.input, testCase.expectNovel, empty)
			}
		})
	}
}
