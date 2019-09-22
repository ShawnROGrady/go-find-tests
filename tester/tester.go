package tester

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

// Tester performs the main testing logic
type Tester struct {
	testPos position
	finder  coverFinder
}

// New constructs a new tester
func New(path string, line, col int) (*Tester, error) {
	pos := position{
		line: line,
		col:  col,
	}

	if err := pos.setFilePkg(path); err != nil {
		return nil, err
	}

	return &Tester{
		testPos: pos,
		finder:  errGroupFinder{},
	}, nil
}

// CoveredBy returns the tests which cover the provided position
func (t *Tester) CoveredBy() ([]string, error) {
	outputDir, err := ioutil.TempDir("", "test_finder")
	if err != nil {
		return []string{}, err
	}
	defer os.RemoveAll(outputDir)

	allTests, err := findTests(t.testPos.pkg)
	if err != nil {
		return []string{}, err
	}

	return t.finder.coveringTests(t, outputDir, allTests)
}

func (t *Tester) runTest(testName, outputDest string) (io.ReadCloser, error) {
	cmd := exec.Command("go", "test", t.testPos.pkg, "-run", testName, "-coverprofile", outputDest)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return os.Open(outputDest)
}
