package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ShawnROGrady/go-find-tests/finder"
)

func printTests(dst io.Writer, tests []string, jsonFmt bool) error {
	if jsonFmt {
		b, err := json.Marshal(tests)
		if err != nil {
			return err
		}
		_, err = dst.Write(b)
		return err
	}
	for i := range tests {
		if _, err := fmt.Fprintf(dst, "%s\n", tests[i]); err != nil {
			return err
		}
	}
	return nil
}

type testPosition struct {
	finder.TestPosition
	SubTests []string `json:"subtests,omitempty"`
}

func printCoveringPostions(dst io.Writer, positions map[string]*testPosition, positionTests []string, jsonFmt bool, lineFmt string) error {
	if jsonFmt {
		b, err := json.Marshal(positions)
		if err != nil {
			return err
		}
		_, err = dst.Write(b)
		return err
	}

	for i := range positionTests {
		if _, err := fmt.Fprintf(dst, "%s\n", fmtPosition(*positions[positionTests[i]], positionTests[i], lineFmt)); err != nil {
			return err
		}
	}
	return nil
}

func fmtPosition(pos testPosition, testName, fmt string) string {
	line := strings.ReplaceAll(fmt, "%t", testName)
	line = strings.ReplaceAll(line, "%f", pos.File)
	line = strings.ReplaceAll(line, "%l", strconv.Itoa(pos.Line))
	line = strings.ReplaceAll(line, "%c", strconv.Itoa(pos.Col))
	line = strings.ReplaceAll(line, "%o", strconv.Itoa(pos.Offset))
	line = strings.ReplaceAll(line, "%s", strings.Join(pos.SubTests, ","))

	return line
}
