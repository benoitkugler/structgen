package crud

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/benoitkugler/structgen/loader"
)

func TestMain(t *testing.T) {
	fn := "../../../goACVE/server/core/rawdata/models.go"
	// fn := "../../../intendance/server/models/models.go"
	pkg, err := loader.LoadSource(fn)
	if err != nil {
		t.Fatal(err)
	}
	fullPath, err := filepath.Abs(fn)
	if err != nil {
		t.Fatal(err)
	}
	typeHandler := NewHandler("skldl", false)
	decls, err := loader.WalkFile(fullPath, pkg, typeHandler)
	if err != nil {
		t.Fatal(err)
	}

	err = decls.Generate(os.Stdout, typeHandler)
	if err != nil {
		t.Fatal(err)
	}
}
