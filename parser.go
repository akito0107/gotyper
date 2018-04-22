package gotyper

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func Parse(file string) ([]*ast.TypeSpec, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", file, 0)
	if err != nil {
		return nil, err
	}
	var types []*ast.TypeSpec
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			types = append(types, x)
		}
		return true
	})

	return types, nil
}
