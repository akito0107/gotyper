package gotyper_test

import (
	"testing"

	"github.com/akito0107/gotyper"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name string
		in   string
		out  []string
	}{
		{
			name: "Simple Struct",
			in: `
package main
type Test struct {
  label string
}
`,
			out: []string{"Test"},
		},
		{
			name: "Named Types",
			in: `package main
type Int int
type Uint uint
type Map map[string]string
`,
			out: []string{"Int", "Uint", "Map"},
		},
		{
			name: "with Bracket",
			in: `package main
type (
	// The Spec type stands for any of *ImportSpec, *ValueSpec, and *TypeSpec.
	Spec interface {
		Node
		specNode()
	}

	// An ImportSpec node represents a single package import.
	ImportSpec struct {
		Doc     *CommentGroup // associated documentation; or nil
		Name    *Ident        // local package name (including "."); or nil
		Path    *BasicLit     // import path
		Comment *CommentGroup // line comments; or nil
		EndPos  token.Pos     // end of spec (overrides Path.Pos if nonzero)
	}
)
`,
			out: []string{"Spec", "ImportSpec"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			types, err := gotyper.Parse(c.in)
			if err != nil {
				t.Fatal(err)
			}

			if len(types) != len(c.out) {
				t.Errorf("should be same length: %v", types)
			}

			for i := 0; i < len(types); i++ {
				if types[i].Name.Name != c.out[i] {
					t.Errorf("must be %s but %s", c.out[i], types[i].Name.Name)
				}
			}
		})
	}
}

//func TestParse2(t *testing.T) {
//	expect := `
//package main
//
//import "reflect"
//
//var TypesMapper map[string]reflect.Type
//
//func init() {
//  typeMapper := map[string]reflect.Type{}
//  typeMapper["Test"] = reflect.TypeOf(Test{})
//}
//`
//
//	fset := token.NewFileSet()
//	f, _ := parser.ParseFile(fset, "", expect, 0)
//	fmt.Println(token.VAR)
//	pp.Print(f)
//}
