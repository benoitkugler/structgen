package creation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/loader"
)

func TestSQL(t *testing.T) {
	// fn := "../../../goACVE/core/rawdata/rawdata.go"
	fn := "../../../intendance/server/models/models.go"
	pkg, err := loader.LoadSource(fn)
	if err != nil {
		t.Fatal(err)
	}
	fullPath, err := filepath.Abs(fn)
	if err != nil {
		t.Fatal(err)
	}
	en, err := enums.FetchEnums(pkg)
	if err != nil {
		t.Fatal(err)
	}
	handler := NewGenHandler(en)
	decls, err := loader.WalkFile(fullPath, pkg, handler)
	if err != nil {
		t.Fatal(err)
	}

	err = decls.Generate(os.Stdout, handler)
	if err != nil {
		t.Fatal(err)
	}
}
