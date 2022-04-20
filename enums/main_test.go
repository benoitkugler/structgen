package enums

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"testing"
	"time"

	"github.com/benoitkugler/structgen/loader"
)

func TestParse(t *testing.T) {
	// fn := "../../goACVE/server/directeurs/types.go"
	fn := "testenums.go"
	pa, err := loader.LoadSource(fn)
	if err != nil {
		t.Fatal(err)
	}
	ti := time.Now()
	l, err := FetchEnums(pa)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("recursively fetching enums :", time.Since(ti))
	if len(l) != 2 {
		t.Errorf("expected 2 enums, got %d", len(l))
	}
}

func commentValueSpecAt(pos token.Pos, file *ast.File) (out string) {
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		if n.Pos() == pos {
			if spec, is := n.(*ast.ValueSpec); is {
				out = strings.TrimSpace(spec.Comment.Text())
				return false
			}
		}
		return true
	})
	return
}

func TestFetch(t *testing.T) {
	fn := "testenums.go"
	pa, err := loader.LoadSource(fn)
	if err != nil {
		t.Fatal(err)
	}

	scope := pa.Types.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		_, ok := obj.(*types.Const)
		if ok {
			declFile := pa.Fset.File(obj.Pos()).Pos(0)
			for _, file := range pa.Syntax { // select the right file
				if file.Pos() == declFile {
					fmt.Println(commentValueSpecAt(obj.Pos(), file))
				}
			}
		}

		// pa.Fset.Position(obj.Pos())
		// pa.Fset.File(obj.Pos()).Name()
		// obj.Pos()
	}
}
