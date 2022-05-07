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
	an := NewAnalyser(en)
	for _, name := range pkg.Types.Scope().Names() {
		if ty, ok := pkg.Types.Scope().Lookup(name).(*types.TypeName); ok {
			jsonT := an.Convert(ty.Type())

			l = append(l, jsonT.Validations()...)
		}
	}
	_ = ioutil.WriteFile("test/valid.sql", []byte(loader.ToString(l)), os.ModePerm)
}

func TestItfs(t *testing.T) {
	fn := "test/models.go"
	pkg, err := loader.LoadSource(fn)
	if err != nil {
		t.Fatal(err)
	}
	en, err := enums.FetchEnums(pkg)
	if err != nil {
		t.Fatal(err)
	}

	an := NewAnalyser(en)
	for _, name := range pkg.Types.Scope().Names() {
		if ty, ok := pkg.Types.Scope().Lookup(name).(*types.TypeName); ok {
			an.Convert(ty.Type())
		}
	}
}
