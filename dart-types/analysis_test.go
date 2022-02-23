package darttypes

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/loader"
)

func TestInterfaces(t *testing.T) {
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

	h.processInterfaces()

	itf := h.interfaces["itfName"]
	if len(itf.members) != 2 {
		t.Fatalf("unexpected members %v", itf.members)
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

func TestMain(t *testing.T) {
	// fn := "../../goACVE/server/core/rawdata/models.go"
	fn := "../../intendance/server/models/models.go"
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

	out, err := os.Create("test.dart")
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()

	err = decls.Generate(out, h)
	if err != nil {
		t.Fatal(err)
	}
}
