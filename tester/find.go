package tester

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
)

func findTests(pkg, runExpr string) ([]string, error) {
	output, err := exec.Command("go", "test", "-list", runExpr, pkg).Output()
	if err != nil {
		return []string{}, parseCommandErr(err)
	}

	var (
		b       = bytes.NewBuffer(output)
		scanner = bufio.NewScanner(b)
		tests   = []string{}
	)

	for scanner.Scan() {
		txt := scanner.Text()
		if strings.HasPrefix(txt, "Test") {
			tests = append(tests, txt)
		}
	}

	return tests, nil
}
