package tester

import (
	"bytes"
	"context"
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
		return []string{}, err
	}

	allTests, err := findTests(t.testPos.pkg)
	if err != nil {
		return []string{}, err
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
	return testBin, err
}

func (t *Tester) runCompiledTest(testName, testBin, outputDir string) (io.ReadCloser, io.Reader, error) {
	var coverOut strings.Builder
	coverOut.WriteString(strings.Replace(testName, "/", "", -1))
	coverOut.WriteString(".out")

	pathToCover := filepath.Join(outputDir, coverOut.String())

	var cmd *exec.Cmd
	if t.includeSubtests {
		cmd = exec.Command(testBin, "-test.run", testName, "-test.coverprofile", coverOut.String(), "-test.outputdir", outputDir, "-test.v")
	} else {
		cmd = exec.Command(testBin, "-test.run", testName, "-test.coverprofile", coverOut.String(), "-test.outputdir", outputDir)
	}

	if t.dir != "" {
		// run test in same dir as file to prevent issues due to dependency on file structure
		cmd.Dir = t.dir
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, nil, err
	}

	coverProf, err := os.Open(pathToCover)
	return coverProf, bytes.NewBuffer(output), err
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
				return err
			}

			prof, err := cover.New(coverout)
			if err != nil {
				return err
			}
			if err := coverout.Close(); err != nil {
				return err
			}

			if prof.Covers(t.testPos.file, t.testPos.line, t.testPos.col) {
				tests[testNum] = testName
				if includeSubtests {
					subs[testNum] = subtests(stdout)
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
			if tests[i] != "" {
				coveredBy = append(coveredBy, tests[i])
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
		for i := range coveredBy {
			if len(coveringSubs[i]) != 0 {
				coveredBy = append(coveredBy[:i], append(coveringSubs[i], coveredBy[i:]...)...)
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
