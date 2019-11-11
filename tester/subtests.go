package tester

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
	"time"
)

// TestEvent represents a single output event of a test run
// see: go doc test2json
type TestEvent struct {
	Time    time.Time // encodes as an RFC3339-format string
	Action  string
	Package string
	Test    string
	Elapsed float64 // seconds
	Output  string
}

func subtests(r io.Reader) ([]string, error) {
	var (
		subtests = []string{}
		scanner  = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		event := TestEvent{}
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return nil, err
		}

		// TODO: should be able to distinguish different levels of sub tests
		if event.Action == "pass" && strings.Contains(event.Test, "/") {
			subtests = append(subtests, event.Test)
		}
	}
	return subtests, nil
}
