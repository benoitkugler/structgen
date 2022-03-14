package darttypes

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/loader"
)

func TestGenerate(t *testing.T) {
	fn := "test/test.go"
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

	h := NewHandler(et)
	decls, err := loader.WalkFile(fullPath, pkg, h)
	if err != nil {
		t.Fatal(err)
	}

	out, err := os.Create("test/gen.dart")
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()

	err = decls.Generate(out, h)
	if err != nil {
		t.Fatal(err)
	}
}
