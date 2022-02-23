package jsonsql

import (
	"go/types"
	"io/ioutil"
	"os"
	"testing"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/loader"
)

func TestValidations(t *testing.T) {
	fn := "test/models.go"
	// fn := "../../../intendance/server/models/models.go"
	pkg, err := loader.LoadSource(fn)
	if err != nil {
		t.Fatal(err)
	}
	en, err := enums.FetchEnums(pkg)
	if err != nil {
		t.Fatal(err)
	}

	var l []loader.Declaration
	for _, name := range pkg.Types.Scope().Names() {
		if ty, ok := pkg.Types.Scope().Lookup(name).(*types.TypeName); ok {
			jsonT := NewTypeJSON(ty.Type(), en)

			l = append(l, jsonT.Validations()...)
		}
	}
	_ = ioutil.WriteFile("test/valid.sql", []byte(loader.Render(l)), os.ModePerm)
}
