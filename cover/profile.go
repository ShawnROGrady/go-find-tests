package cover

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Re-implementation of https://github.com/golang/tools/blob/master/cover/profile.go
// found the current API a bit awkward for testing + dealing w/ different profiles

func init() {
	outLineReg = regexp.MustCompile(outLineFmt)
}

var outLineReg *regexp.Regexp

const (
	outLineFmt   = `^([a-zA-Z0-9\/\.\-_]+)\/([a-zA-Z0-9\-_]+.go):([0-9]+)\.([0-9]+),([0-9]+)\.([0-9]+) ([0-9]+) ([0-9]+)$`
	packageIndex = iota
	fileIndex
	startLineIndex
	startColIndex
	endLineIndex
	endColIndex
	numStmtIndex
	countIndex
)

// Profile represents the output of a cover profile
type Profile map[string]coverBlocks

// New contructs a new Profile
func New(r io.Reader) (*Profile, error) {
	var (
		prof    = make(Profile)
		scanner = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		s := scanner.Text()
		if strings.HasPrefix(s, "mode: ") {
			continue
		}
		line, err := parseLine(s)
		if err != nil {
			return nil, err
		}
		prof[line.file] = append(prof[line.file], line.coverBlock)
	}
	// sort blocks for easy traversal
	for k := range prof {
		sort.Sort(prof[k])
	}
	return &prof, nil
}

// Covers returns whether or not the statement at the given position is covered by the profile
func (p *Profile) Covers(file string, line, col int) bool {
	if prof, ok := (*p)[file]; ok {
		for i := range prof {
			if prof[i].inBlock(line, col) {
				return prof[i].count != 0
			}
		}
	}
	return false
}

// alias to implement sort.Interface
type coverBlocks []coverBlock

func (c coverBlocks) Len() int {
	return len(c)
}

func (c coverBlocks) Less(i, j int) bool {
	if c[i].startLine == c[j].startLine {
		return c[i].startCol < c[j].endCol
	}
	return c[i].startLine < c[j].startLine
}

func (c coverBlocks) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type coverBlock struct {
	startLine int
	startCol  int
	endLine   int
	endCol    int
	numStmt   int
	count     int
}

func (c coverBlock) inBlock(line, col int) bool {
	if c.startLine <= line && c.endLine >= line {
		if c.startLine == line && c.endLine == line {
			return c.startCol <= col && c.endCol >= col
		}
		if c.startLine == line {
			return c.startCol <= col
		}
		if c.endLine == line {
			return c.endCol >= col
		}
		return true
	}
	return false
}

type coverLine struct {
	pkg  string
	file string
	coverBlock
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
		pkg:  subexps[packageIndex],
		file: subexps[fileIndex],
		coverBlock: coverBlock{
			startLine: startL,
			startCol:  startC,
			endLine:   endL,
			endCol:    endC,
			numStmt:   stmts,
			count:     count,
		},
	}, nil
}
