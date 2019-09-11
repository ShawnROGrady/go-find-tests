package cover

import (
	"fmt"
	"regexp"
	"strconv"
)

// Re-implementation of https://github.com/golang/tools/blob/master/cover/profile.go
// found the current API a bit awkward for testing + dealing w/ different profiles

func init() {
	outLineReg = regexp.MustCompile(outLineFmt)
}

var outLineReg *regexp.Regexp

const (
	outLineFmt   = `^([a-zA-Z0-9/\.]+)/([a-zA-Z0-9]+.go):([0-9]+)\.([0-9]+),([0-9]+)\.([0-9]+) ([0-9]+) ([0-9]+)$`
	packageIndex = iota
	fileIndex
	startLineIndex
	startColIndex
	endLineIndex
	endColIndex
	numStmtIndex
	countIndex
)

type coverLine struct {
	pkg       string
	file      string
	startLine int
	startCol  int
	endLine   int
	endCol    int
	numStmt   int
	count     int
}

func parseLine(line string) (coverLine, error) {
	subexps := outLineReg.FindStringSubmatch(line)
	if len(subexps) == 0 {
		return coverLine{}, fmt.Errorf("Unmatched line: %s", line)
	}

	startL, err := strconv.Atoi(subexps[startLineIndex])
	if err != nil {
		return coverLine{}, err
	}
	startC, err := strconv.Atoi(subexps[startColIndex])
	if err != nil {
		return coverLine{}, err
	}
	endL, err := strconv.Atoi(subexps[endLineIndex])
	if err != nil {
		return coverLine{}, err
	}
	endC, err := strconv.Atoi(subexps[endColIndex])
	if err != nil {
		return coverLine{}, err
	}
	stmts, err := strconv.Atoi(subexps[numStmtIndex])
	if err != nil {
		return coverLine{}, err
	}
	count, err := strconv.Atoi(subexps[countIndex])
	if err != nil {
		return coverLine{}, err
	}

	return coverLine{
		pkg:       subexps[packageIndex],
		file:      subexps[fileIndex],
		startLine: startL,
		startCol:  startC,
		endLine:   endL,
		endCol:    endC,
		numStmt:   stmts,
		count:     count,
	}, nil
}
