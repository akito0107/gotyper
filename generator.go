package gotyper

import (
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"strconv"
)

type Generator struct {
	w   io.Writer
	pak string
}

func NewGenerator(w io.Writer) *Generator {
	return &Generator{w, "main"}
}

func (g *Generator) Generate(types []*ast.TypeSpec) error {

	f := &ast.File{
		Name: ast.NewIdent(g.pak),
	}

	var decls []ast.Decl

	decls = append(decls, &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: strconv.Quote("reflect"),
				},
			},
		},
	})

	decls = append(decls, &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{
					ast.NewIdent("TypesMapper"),
				},
				Type: &ast.MapType{
					Key: ast.NewIdent("string"),
					Value: &ast.SelectorExpr{
						X:   ast.NewIdent("reflect"),
						Sel: ast.NewIdent("Type"),
					},
				},
			},
		},
	})
	assign := &ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{
			ast.NewIdent("TypesMapper"),
		},
		Rhs: []ast.Expr{
			&ast.CompositeLit{
				Type: &ast.MapType{
					Key: ast.NewIdent("string"),
					Value: &ast.SelectorExpr{
						X:   ast.NewIdent("reflect"),
						Sel: ast.NewIdent("Type"),
					},
				},
				Elts: []ast.Expr{},
			},
		},
	}
	stmt := []ast.Stmt{assign}
	for _, t := range types {
		stmt = append(stmt, &ast.AssignStmt{
			Tok: token.ASSIGN,
			Lhs: []ast.Expr{
				&ast.IndexExpr{
					X: ast.NewIdent("TypesMapper"),
					Index: &ast.BasicLit{
						Kind:  token.STRING,
						Value: strconv.Quote(t.Name.Name),
					},
				},
			},
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("reflect"),
						Sel: ast.NewIdent("TypeOf"),
					},
					Args: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.Ident{
								Name: "Test",
							},
							Elts: []ast.Expr{},
						},
					},
				},
			},
		})
	}
	decls = append(decls, &ast.FuncDecl{
		Name: ast.NewIdent("init"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
		},
		Body: &ast.BlockStmt{
			List: stmt,
		},
	})

	f.Decls = decls

	format.Node(g.w, token.NewFileSet(), f)
	return nil
}
