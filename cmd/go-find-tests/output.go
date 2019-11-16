package main

import (
	"encoding/json"
	"fmt"
	"io"

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

func printCoveringPostions(dst io.Writer, positions map[string]finder.TestPosition, jsonFmt bool) error {
	if jsonFmt {
		b, err := json.Marshal(positions)
		if err != nil {
			return err
		}
		_, err = dst.Write(b)
		return err
	}

	for k, v := range positions {
		// TODO: allow output fmt to be changed
		fmt.Fprintf(dst, "%s:%s\n", k, v)
	}
	return nil
}
