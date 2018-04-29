package main

import (
	"flag"
	"go/ast"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/akito0107/gotyper"
	"github.com/pkg/errors"
)

var dryrun = flag.Bool("dryrun", false, "dryrun mode")
var path = flag.String("path", "./", "input package path")
var filename = flag.String("out", "typer.go", "output file name")

func main() {
	flag.Parse()

	fls, err := ioutil.ReadDir(*path)
	if err != nil {
		log.Fatal(err)
	}

	var files []os.FileInfo

	for _, f := range fls {
		if strings.HasSuffix(f.Name(), ".go") && !strings.HasSuffix(f.Name(), "_test.go") {
			files = append(files, f)
		}
	}

	var spec []*ast.TypeSpec
	var pkgName string

	for _, f := range files {
		file, err := os.Open(*path + f.Name())
		if err != nil {
			log.Fatal(err)
		}
		b, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
		s, pkg, err := gotyper.Parse(string(b))
		if err != nil {
			log.Fatal(err)
		}
		if pkgName != "" && pkgName != pkg {
			log.Fatal(errors.Errorf("multiple package: %s, %s", pkgName, pkg))
		}
		pkgName = pkg
		spec = append(spec, s...)
	}
	var out io.Writer
	if *dryrun {
		out = os.Stdout
	} else {
		f, err := os.Create(*path + *filename)
		if err != nil {
			log.Fatal(err)
		}
		out = f
	}

	gen := gotyper.NewGenerator(out, pkgName)
	if err := gen.Generate(spec); err != nil {
		log.Fatal(err)
	}
}
