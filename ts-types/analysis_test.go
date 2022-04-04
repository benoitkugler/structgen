package tstypes

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"testing"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/loader"
	"github.com/benoitkugler/structgen/utils"
)

func TestMain(t *testing.T) {
	fn := "../../goACVE/server/core/rawdata/models.go"
	// fn := "../../intendance/server/models/models.go"
	pkg, err := loader.LoadSource(fn)
	if err != nil {
		t.Fatal(err)
	}

	et, err := enums.FetchEnums(pkg)
	if err != nil {
		t.Fatal(err)
	}

	fullPath, err := filepath.Abs(fn)
	if err != nil {
		t.Fatal(err)
	}

	h := NewHandler(et, pkg.Types.Scope())
	decls, err := loader.WalkFile(fullPath, pkg, h)
	if err != nil {
		t.Fatal(err)
	}

	if err := decls.Generate(os.Stdout, h); err != nil {
		t.Fatal(err)
	}
}

func TestTime(t *testing.T) {
	const source = `package main

	import "time"

	var t = time.Now()

	func main() {
	}
	`

	fset := token.NewFileSet()

	// Parse the input string, []byte, or io.Reader,
	// recording position information in fset.
	// ParseFile returns an *ast.File, a syntax tree.
	f, err := parser.ParseFile(fset, "hello.go", source, 0)
	if err != nil {
		t.Fatal(err) // parse error
	}

	// A Config controls various options of the type checker.
	// The defaults work fine except for one setting:
	// we must specify how to deal with imports.
	conf := types.Config{Importer: importer.Default()}

	// Type-check the package containing only file f.
	// Check returns a *types.Package.
	pkg, err := conf.Check("cmd/hello", fset, []*ast.File{f}, nil)
	if err != nil {
		t.Fatal(err) // type error
	}

	objTime := pkg.Scope().Lookup(pkg.Scope().Names()[1])
	fmt.Println(objTime.Type().String())

	if utils.IsUnderlyingTime(objTime.Type()) != true {
		t.Fatal()
	}
}
