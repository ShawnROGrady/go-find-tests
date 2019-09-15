package tester

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// position represents the position we want to test
type position struct {
	file, pkg string
	line, col int
}

// setFilePkg sets the file and package from the provided path
func (p *position) setFilePkg(path string) error {
	dir, file := filepath.Split(path)
	pkg, err := packageName(dir)
	if err != nil {
		return err
	}
	p.file = file
	p.pkg = pkg
	return nil
}

// packageName returns the go package name associated with the provided directory
func packageName(dir string) (string, error) {
	output, err := exec.Command("go", "list", dir).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(output), "\n"), nil
}
