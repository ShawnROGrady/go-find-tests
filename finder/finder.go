package finder

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// TestPosition represents the location of a tests declaration
type TestPosition struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Col    int    `json:"col"`
	Offset int    `json:"offset"`
}

// PackageTests returns the positions of all tests within a package
func PackageTests(dir string) (map[string]TestPosition, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(f os.FileInfo) bool {
		return strings.HasSuffix(f.Name(), "_test.go")
	}, 0)
	if err != nil {
		return nil, err
	}

	testFuncs := make(map[string]TestPosition)

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			funcFinder := &testFuncFinder{
				fset:      fset,
				testFuncs: make(map[string]TestPosition),
			}
			ast.Walk(funcFinder, file)
			for k, v := range funcFinder.testFuncs {
				testFuncs[k] = v
			}
		}
	}
	return testFuncs, nil
}
