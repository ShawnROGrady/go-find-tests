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
		fmt.Fprintf(dst, "%s\n", tests[i])
	}
	return nil
}

func printCoveringPostions(dst io.Writer, positions map[string]finder.TestPosition, jsonFmt bool, lineFmt string) error {
	if jsonFmt {
		b, err := json.Marshal(positions)
		if err != nil {
			return err
		}
		_, err = dst.Write(b)
		return err
	}

	for k, v := range positions {
		fmt.Fprintf(dst, "%s\n", fmtPosition(v, k, lineFmt))
	}
	return nil
}

func fmtPosition(pos finder.TestPosition, testName, fmt string) string {
	line := strings.ReplaceAll(fmt, "%t", testName)
	line = strings.ReplaceAll(line, "%f", pos.File)
	line = strings.ReplaceAll(line, "%l", strconv.Itoa(pos.Line))
	line = strings.ReplaceAll(line, "%c", strconv.Itoa(pos.Col))
	line = strings.ReplaceAll(line, "%o", strconv.Itoa(pos.Offset))

	return line
}
