package interfaces

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/benoitkugler/structgen/loader"
)

func TestAnalyze(t *testing.T) {
	fn := "test/test.go"
	pkg, err := loader.LoadSource(fn)
	if err != nil {
		t.Fatal(err)
	}
	fullPath, err := filepath.Abs(fn)
	if err != nil {
		t.Fatal(err)
	}

	h := NewHandler("test")
	decls, err := loader.WalkFile(fullPath, pkg, h)
	if err != nil {
		t.Fatal(err)
	}

	itfs := h.(handler).analyzer.Process()
	if len(itfs) != 2 {
		t.Fatal(itfs)
	}
	if itfs[0].Name.Obj().Name() != "union1" {
		t.Fatal(itfs[0])
	}
	if len(itfs[0].Members) != 3 || len(itfs[1].Members) != 1 {
		t.Fatal()
	}

	out, err := os.Create("test/gen.go")
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()

	err = decls.Generate(out, h)
	if err != nil {
		t.Fatal(err)
	}

	err = exec.Command("goimports", "-w", "test/gen.go").Run()
	if err != nil {
		t.Fatal(err)
	}
}
