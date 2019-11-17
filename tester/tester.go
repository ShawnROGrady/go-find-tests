package tester

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Tester performs the main testing logic
type Tester struct {
	testPos         position
	includeSubtests bool
	short           bool
	run             string
	dir             string // directory of test
	coverFinder     coverFinder
}

// Config represents configuration options for the Tester
type Config struct {
	IncludeSubtests bool
	Short           bool   // sets '-short' when running tests
	Run             string // which tests should be run, if empty defaults to '.' (sets '-list' flag)
	Seq             bool   // all tests should be run sequentially
}

// New constructs a new tester
func New(path string, line, col int, conf Config) (*Tester, error) {
	pos := position{
		line: line,
		col:  col,
	}

	if err := pos.setFilePkg(path); err != nil {
		return nil, err
	}

	var (
		dir string
		err error
	)
	if strings.HasPrefix(path, ".") {
		rel, _ := filepath.Split(path)
		dir, err = filepath.Abs(rel)
	}
	if err != nil {
		return nil, err
	}

	runExp := "." // should default to running all
	if conf.Run != "" {
		runExp = "."
	}

	var finder coverFinder
	if conf.Seq {
		finder = sequentialFinder{}
	} else {
		finder = errGroupFinder{}
	}

	return &Tester{
		testPos:         pos,
		includeSubtests: conf.IncludeSubtests,
		short:           conf.Short,
		run:             runExp,
		dir:             dir,
		coverFinder:     finder,
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
		return []string{}, fmt.Errorf("error compiling test for go pkg %s: %s", t.testPos.pkg, err)
	}

	allTests, err := findTests(t.testPos.pkg, t.run)
	if err != nil {
		return []string{}, fmt.Errorf("error finding tests in go pkg %s: %s", t.testPos.pkg, err)
	}

	if len(allTests) == 0 {
		return []string{}, nil
	}

	return t.coverFinder.coveringTests(t, testBin, outputDir, allTests, t.includeSubtests)
}

func (t *Tester) compileTest(outputDir string) (string, error) {
	var binName strings.Builder
	s := strings.Split(t.testPos.pkg, "/")
	binName.WriteString(s[len(s)-1])
	binName.WriteString(".test")

	testBin := filepath.Join(outputDir, binName.String())

	cmd := exec.Command("go", "test", t.testPos.pkg, "-cover", "-c", "-o", testBin)
	err := cmd.Run()
	if err != nil {
		return "", parseCommandErr(err)
	}
	return testBin, nil
}

func (t *Tester) runCompiledTest(testName, testBin, outputDir string) (io.ReadCloser, io.Reader, error) {
	var coverOut strings.Builder
	coverOut.WriteString(strings.Replace(testName, "/", "", -1))
	coverOut.WriteString(".out")

	pathToCover := filepath.Join(outputDir, coverOut.String())

	cmdArgs := []string{"tool", "test2json", testBin, "-test.run", testName, "-test.coverprofile", coverOut.String(), "-test.outputdir", outputDir}
	if t.includeSubtests {
		cmdArgs = append(cmdArgs, "-test.v")
	}
	if t.short {
		cmdArgs = append(cmdArgs, "-test.short")
	}

	cmd := exec.Command("go", cmdArgs...)

	if t.dir != "" {
		// run test in same dir as file to prevent issues due to dependency on file structure
		cmd.Dir = t.dir
	}

	var buf bytes.Buffer
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		return nil, nil, parseTestError(err, &buf)
	}

	coverProf, err := os.Open(pathToCover)
	return coverProf, &buf, err
}
