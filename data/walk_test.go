package data

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/loader"
)

func TestMain(t *testing.T) {
	fn := "../../goACVE/server/core/rawdata/models.go"
	// fn := "../../intendance/server/models/models.go"
	pkg, err := loader.LoadSource(fn)
	if err != nil {
		t.Fatal(err)
	}
	en, err := enums.FetchEnums(pkg)
	if err != nil {
		t.Fatal(err)
	}
	fullPath, err := filepath.Abs(fn)
	if err != nil {
		t.Fatal(err)
	}
	decls, err := loader.WalkFile(fullPath, pkg, Handler{PackageName: "models", EnumsTable: en})
	if err != nil {
		t.Fatal(err)
	}
	if err := decls.Render(os.Stdout); err != nil {
		t.Fatal(err)
	}
}
