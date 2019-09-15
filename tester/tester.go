package tester

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/ShawnROGrady/go-find-tests/cover"
)

// Tester performs the main testing logic
type Tester struct {
	testPos position
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

	coveredBy := []string{}
	// TODO: each loop iteration should be handled by separate go routine
	for i := range allTests {
		var dst strings.Builder
		dst.WriteString(outputDir)
		dst.WriteString(allTests[i])
		dst.WriteString(".out")
		output, err := t.runTest(allTests[i], dst.String())
		if err != nil {
			return []string{}, err
		}

		prof, err := cover.New(output)
		if err != nil {
			return []string{}, err
		}
		if err := output.Close(); err != nil {
			return []string{}, err
		}

		if prof.Covers(t.testPos.file, t.testPos.line, t.testPos.col) {
			coveredBy = append(coveredBy, allTests[i])
		}
	}
	return coveredBy, nil
}

func (t *Tester) runTest(testName, outputDest string) (io.ReadCloser, error) {
	cmd := exec.Command("go", "test", t.testPos.pkg, "-run", testName, "-coverprofile", outputDest)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return os.Open(outputDest)
}
