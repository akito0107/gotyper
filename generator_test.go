package gotyper_test

import (
	"bytes"
	"testing"

	"fmt"

	"github.com/akito0107/gotyper"
	"github.com/andreyvit/diff"
)

func TestGenerator_Generate(t *testing.T) {
	var buf bytes.Buffer
	g := gotyper.NewGenerator(&buf)

	in := `
package main
type Test struct {
  name string
}
`
	expect := `package main

import "reflect"

var TypesMapper map[string]reflect.Type

func init() {
    TypesMapper = new(map[string]reflect.Type)
    TypesMapper["Test"] = reflect.TypeOf(Test)
}
`
	types, _ := gotyper.Parse(in)
	if err := g.Generate(types); err != nil {
		t.Fatal(err)
	}

	if act := buf.String(); act != expect {
		t.Errorf("expect: \n%s\n but got: \n%s\n", expect, act)
		fmt.Println(diff.CharacterDiff(expect, act))
	}

}
