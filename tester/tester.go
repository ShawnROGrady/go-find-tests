package tester

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ShawnROGrady/go-find-tests/cover"
	"golang.org/x/sync/errgroup"
)

// Tester performs the main testing logic
type Tester struct {
	testPos         position
	includeSubtests bool
	dir             string // directory of test
}

// New constructs a new tester
func New(path string, line, col int, includeSubtests bool) (*Tester, error) {
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

	return &Tester{
		testPos:         pos,
		includeSubtests: includeSubtests,
		dir:             dir,
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

	allTests, err := findTests(t.testPos.pkg)
	if err != nil {
		return []string{}, fmt.Errorf("error finding tests in go pkg %s: %s", t.testPos.pkg, err)
	}

	return t.coveringTests(testBin, outputDir, allTests, t.includeSubtests)
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

	var cmd *exec.Cmd
	if t.includeSubtests {
		cmd = exec.Command("go", "tool", "test2json", testBin, "-test.run", testName, "-test.coverprofile", coverOut.String(), "-test.outputdir", outputDir, "-test.v")
	} else {
		cmd = exec.Command("go", "tool", "test2json", testBin, "-test.run", testName, "-test.coverprofile", coverOut.String(), "-test.outputdir", outputDir)
	}

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

func (t *Tester) coveringTests(testBin, outputDir string, allTests []string, includeSubtests bool) ([]string, error) {
	var (
		ctx       = context.Background()
		coveredBy = []string{}
		tests     = make([]string, len(allTests))
		subs      = make([][]string, len(allTests))
	)
	g, ctx := errgroup.WithContext(ctx)

	for i := range allTests {
		testNum := i
		testName := allTests[i]
		g.Go(func() error {
			coverout, stdout, err := t.runCompiledTest(testName, testBin, outputDir)
			if err != nil {
				return fmt.Errorf("error running test '%s': %s", testName, err)
			}

			prof, err := cover.New(coverout)
			if err != nil {
				return fmt.Errorf("error parsing coverage output: %s", err)
			}
			if err := coverout.Close(); err != nil {
				return err
			}

			if prof.Covers(t.testPos.file, t.testPos.line, t.testPos.col) {
				tests[testNum] = testName
				if includeSubtests {
					subs[testNum], err = subtests(stdout)
					if err != nil {
						return fmt.Errorf("error finding subtests: %s", err)
					}
				}
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return []string{}, err
	}

	if includeSubtests {
		errGroup, _ := errgroup.WithContext(ctx)
		coveringSubs := make([][]string, len(tests))
		for i := range tests {
			if tests[i] == "" {
				continue
			}
			testNum := i
			errGroup.Go(func() error {
				if len(subs[testNum]) != 0 {
					if len(subs[testNum]) == 1 {
						coveringSubs[testNum] = subs[testNum]
						return nil
					}
					coveringSubTests, err := t.coveringTests(testBin, outputDir, subs[testNum], false)
					if err != nil {
						return err
					}
					coveringSubs[testNum] = coveringSubTests
				}
				return nil
			})
		}
		if err := errGroup.Wait(); err != nil {
			return []string{}, err
		}
		for i := range tests {
			if tests[i] != "" {
				coveredBy = append(coveredBy, tests[i])
			}
			if len(coveringSubs[i]) != 0 {
				coveredBy = append(coveredBy, coveringSubs[i]...)
			}
		}
	} else {
		for i := range tests {
			if tests[i] != "" {
				coveredBy = append(coveredBy, tests[i])
			}
		}
	}

	return coveredBy, nil
}
