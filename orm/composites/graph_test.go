package composites

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/benoitkugler/structgen/loader"
)

func TestGraph(t *testing.T) {
	fn := "../../../goACVE/server/core/rawdata/models.go"
	pkg, err := loader.LoadSource(fn)
	if err != nil {
		t.Fatal(err)
	}
	fullPath, err := filepath.Abs(fn)
	if err != nil {
		t.Fatal(err)
	}
	handler := Composites{OriginPackageName: "skldl"}
	loader.WalkFile(fullPath, pkg, &handler)
	tables := handler.tables

	g := newGraph(tables)
	err = g.render(handler.OriginPackageName, os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
}
