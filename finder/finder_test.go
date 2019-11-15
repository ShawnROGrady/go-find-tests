package finder

import (
	"testing"
)

var packageTestsTests = map[string]struct {
	dir               string
	expectedPositions map[string]TestPosition
	expectErr         bool
}{
	"10_tests_1_file": {
		dir: "../testdata/len10",
		expectedPositions: map[string]TestPosition{
			"TestEmptyStringIsEmpty": {
				File:   "../testdata/len10/len_test.go",
				Line:   8,
				Col:    1,
				Offset: 52,
			},
			"TestEmptyStringIsShort": {
				File:   "../testdata/len10/len_test.go",
				Line:   43,
				Col:    1,
				Offset: 775,
			},
			"TestLongStringIsEmpty": {
				File:   "../testdata/len10/len_test.go",
				Line:   22,
				Col:    1,
				Offset: 309,
			},
			"TestLongStringIsShort": {
				File:   "../testdata/len10/len_test.go",
				Line:   57,
				Col:    1,
				Offset: 1032,
			},
			"TestNovelIsEmpty": {
				File:   "../testdata/len10/len_test.go",
				Line:   36,
				Col:    1,
				Offset: 614,
			},
			"TestNovelIsShort": {
				File:   "../testdata/len10/len_test.go",
				Line:   71,
				Col:    1,
				Offset: 1337,
			},
			"TestShortStringIsEmpty": {
				File:   "../testdata/len10/len_test.go",
				Line:   15,
				Col:    1,
				Offset: 179,
			},
			"TestShortStringIsShort": {
				File:   "../testdata/len10/len_test.go",
				Line:   50,
				Col:    1,
				Offset: 900,
			},
			"TestVeryLongStringIsEmpty": {
				File:   "../testdata/len10/len_test.go",
				Line:   29,
				Col:    1,
				Offset: 445,
			},
			"TestVeryLongStringIsShort": {
				File:   "../testdata/len10/len_test.go",
				Line:   64,
				Col:    1,
				Offset: 1168,
			},
		},
	},
}

func TestPackageTests(t *testing.T) {
	for testName, testCase := range packageTestsTests {
		t.Run(testName, func(t *testing.T) {
			pkgTests, err := PackageTests(testCase.dir)
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

			if len(testCase.expectedPositions) != len(pkgTests) {
				t.Errorf("Unexpected tests (expected = %v, actual = %v)", testCase.expectedPositions, pkgTests)
				return
			}

			for k, v := range testCase.expectedPositions {
				if pkgTests[k] != v {
					t.Errorf("Unexpected positions[%s] (expected = %v, actual = %v)", k, v, pkgTests[k])
				}
			}
		})
	}
}
