package tester

import "testing"

var setFilePkgTests = map[string]struct {
	path                      string
	expectedFile, expectedPkg string
	expectErr                 bool
}{
	"std_lib_file": {
		path:         "fmt/format.go",
		expectedFile: "format.go",
		expectedPkg:  "fmt",
		expectErr:    false,
	},
	"this_file": {
		path:         "./package_test.go",
		expectedFile: "package_test.go",
		expectedPkg:  "github.com/ShawnROGrady/go-find-tests/tester",
		expectErr:    false,
	},
	"invalid_file": {
		path:      "./bad/path/to/file.go",
		expectErr: true,
	},
}

func TestSetFilePkg(t *testing.T) {
	for testName, test := range setFilePkgTests {
		t.Run(testName, func(t *testing.T) {
			p := &Position{}

			err := p.SetFilePkg(test.path)
			if err != nil {
				if test.expectErr {
					return
				}
				t.Errorf("Unexpected err: %s", err)
				return
			}
			if test.expectErr {
				t.Errorf("Unexpectedly no error")
			}

			if test.expectedFile != p.file {
				t.Errorf("Unexpected file (expected = '%s', actual = '%s')", test.expectedFile, p.file)
			}
			if test.expectedPkg != p.pkg {
				t.Errorf("Unexpected pkg (expected = '%s', actual = '%s')", test.expectedPkg, p.pkg)
			}
		})
	}
}
