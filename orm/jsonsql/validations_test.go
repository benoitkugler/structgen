package jsonsql

import (
	"go/types"
	"io/ioutil"
	"os"
	"testing"

	"github.com/benoitkugler/structgen/loader"
)

func TestValidations(t *testing.T) {
	// fn := "test/models.go"
	fn := "../../../intendance/server/models/models.go"
	pkg, err := loader.LoadSource(fn)
	if err != nil {
		t.Fatal(err)
	}
	l := loader.NewDeclarations()
	for _, name := range pkg.Types.Scope().Names() {
		if ty, ok := pkg.Types.Scope().Lookup(name).(*types.TypeName); ok {
			jsonT := NewTypeJSON(ty.Type())
			jsonT.AddValidation(l)
		}
	}
	_ = ioutil.WriteFile("test/valid.sql", []byte(l.ToString()), os.ModePerm)
}
