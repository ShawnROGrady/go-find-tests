package finder

import (
	"go/ast"
	"go/token"
	"strings"
)

type testFuncFinder struct {
	fset      *token.FileSet
	testFuncs map[string]TestPosition
}

func (f *testFuncFinder) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		fun := n.Type
		if !strings.HasPrefix(n.Name.Name, "Test") {
			return nil
		}
		currentFile := f.fset.File(fun.Func)
		pos := currentFile.Position(fun.Func)
		f.testFuncs[n.Name.Name] = TestPosition{
			File:   currentFile.Name(),
			Line:   pos.Line,
			Col:    pos.Column,
			Offset: pos.Offset,
		}
		// TODO: should search for subs
		return nil
	}
	return f
}
