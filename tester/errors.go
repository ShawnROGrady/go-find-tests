package tester

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// testErr represents an error that occured while running a test
type testErr struct {
	testName string
	output   string
}

func (t *testErr) Error() string {
	return fmt.Sprintf("FAIL: %s - %s", t.testName, strings.TrimSpace(t.output))
}

// parseTestError reads the output of a test run and attempts to construct a human-readable testErr
func parseTestError(origErr error, output io.Reader) error {
	var (
		scanner      = bufio.NewScanner(output)
		lastOutEvent TestEvent
	)

	if _, ok := origErr.(*exec.ExitError); !ok {
		// errors running a test will result in an exec.ExitError
		return origErr
	}

	for scanner.Scan() {
		event := TestEvent{}
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			// swallow unmarshal error since it won't provide any additional context as to why there was an error initially
			return origErr
		}

		switch event.Action {
		case "output":
			lastOutEvent = event
		case "fail":
			return &testErr{testName: event.Test, output: lastOutEvent.Output}
		}
	}

	return parseCommandErr(origErr)
}

// commandErr is an alternate to exec.ExitError which includes stderr if present
type commandErr struct {
	err    error
	stderr string
}

func (c *commandErr) Error() string {
	return fmt.Sprintf("%s - %s", c.err, c.stderr)
}

func parseCommandErr(origErr error) error {
	if exitErr, ok := origErr.(*exec.ExitError); ok && len(exitErr.Stderr) != 0 {
		return &commandErr{err: origErr, stderr: string(exitErr.Stderr)}
	}

	return origErr
}
