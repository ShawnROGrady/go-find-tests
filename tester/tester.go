package tester

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	testBin, err := t.compileTest(outputDir)
	if err != nil {
		return []string{}, err
	}

	allTests, err := findTests(t.testPos.pkg)
	if err != nil {
		return []string{}, err
	}

	return t.finder.coveringTests(t, testBin, outputDir, allTests)
}

func (t *Tester) compileTest(outputDir string) (string, error) {
	var binName strings.Builder
	s := strings.Split(t.testPos.pkg, "/")
	binName.WriteString(s[len(s)-1])
	binName.WriteString(".test")

	testBin := filepath.Join(outputDir, binName.String())

	cmd := exec.Command("go", "test", t.testPos.pkg, "-cover", "-c", "-o", testBin)
	err := cmd.Run()
	return testBin, err
}

func (t *Tester) runCompiledTest(testName, testBin, outputDir string) (io.ReadCloser, error) {
	var coverOut strings.Builder
	coverOut.WriteString(testName)
	coverOut.WriteString(".out")

	cmd := exec.Command(testBin, "-test.run", testName, "-test.coverprofile", coverOut.String(), "-test.outputdir", outputDir)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return os.Open(filepath.Join(outputDir, coverOut.String()))
}
