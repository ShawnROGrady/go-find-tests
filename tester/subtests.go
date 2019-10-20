package tester

import (
	"bufio"
	"io"
	"regexp"
)

func init() {
	outSubtestReg = regexp.MustCompile(outSubtestFmt)
}

var outSubtestReg *regexp.Regexp

const (
	outSubtestFmt = `^=== RUN   ([a-zA-Z0-9\-\_]+\/[a-zA-Z0-9\-\_]+)`
)

func subtests(r io.Reader) []string {
	var (
		subtests = []string{}
		scanner  = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		s := scanner.Text()
		subexps := outSubtestReg.FindStringSubmatch(s)
		if len(subexps) == 0 {
			continue
		}

		subTest := subexps[1]
		subtests = append(subtests, subTest)
	}
	return subtests
}
