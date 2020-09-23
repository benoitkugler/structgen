// This file create helpers functions
// to generate random data for tests
package data

import (
	"fmt"
	"go/types"
	"io"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/loader"
	"github.com/benoitkugler/structgen/utils"
)

type Handler struct {
	PackageName string
	EnumsTable  enums.EnumTable
}

func (d Handler) HandleType(topLevelDecl *loader.Declarations, typ types.Type) {
	d.analyseType(topLevelDecl, typ)
}
func (d Handler) HandleComment(topLevelDecl *loader.Declarations, comment loader.Comment) error {
	return nil
}

func (d Handler) WriteHeader(w io.Writer) error {
	_, err := fmt.Fprintf(w, "package %s", d.PackageName)
	return err
}
func (d Handler) WriteFooter(w io.Writer) error { return nil }

func (d Handler) convertFields(topLevelDecl *loader.Declarations, structType *types.Struct) []structField {
	var out []structField
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		dataFn := d.analyseType(topLevelDecl, field.Type())
		out = append(out, structField{Name: field.Name(), Id: dataFn.Id()})
	}
	return out
}

func isTime(typ *types.Named) bool {
	return typ.Obj().Pkg().Name() == "time" && typ.Obj().Name() == "Time"
}

func (d Handler) analyseType(topLevelDecl *loader.Declarations, typ types.Type) DataFunction {
	var decl DataFunction
	if named, isNamed := typ.(*types.Named); isNamed {
		// special case for structs :
		// we dont generate a random function for underlying type
		if st, isStruct := typ.Underlying().(*types.Struct); isStruct {
			// special case for time.Time, we use a shortcut
			if utils.IsUnderlyingTime(typ) {
				underFn := FnTime{type_: named}
				topLevelDecl.Add(underFn)
				if isTime(named) {
					decl = underFn
				} else {
					decl = FnNamed{TargetPackage: d.PackageName, Type_: named, Underlying: underFn}
				}
			} else {
				fields := d.convertFields(topLevelDecl, st)
				decl = FnStruct{TargetPackage: d.PackageName, Type_: named, Fields: fields}
			}
		} else if enum, isEnum := d.EnumsTable[named.Obj().Name()]; isEnum {
			decl = FnEnum{TargetPackage: d.PackageName, Type_: named, Underlying: enum}
		} else {
			// extract underlying type
			underFn := d.analyseType(topLevelDecl, typ.Underlying())
			decl = FnNamed{TargetPackage: d.PackageName, Type_: named, Underlying: underFn}
		}

		// add top level declaration
		topLevelDecl.Add(decl)
		return decl
	}

	switch typU := typ.Underlying().(type) {
	case *types.Basic:
		decl = FnBasic{type_: typU}
	case *types.Pointer:
		// indirection
		elem := d.analyseType(topLevelDecl, typU.Elem())
		decl = FnPointer{TargetPackage: d.PackageName, Elem: elem}
	case *types.Struct:
		panic("annonymous struct are not supported")
	case *types.Array:
		valueFn := d.analyseType(topLevelDecl, typU.Elem())
		decl = FnArray{TargetPackage: d.PackageName, Length: typU.Len(), Elem: valueFn}
	case *types.Slice:
		valueFn := d.analyseType(topLevelDecl, typU.Elem())
		decl = FnSlice{TargetPackage: d.PackageName, Elem: valueFn}
	case *types.Map:
		decl = FnMap{TargetPackage: d.PackageName,
			Key:  d.analyseType(topLevelDecl, typU.Key()),
			Elem: d.analyseType(topLevelDecl, typU.Elem()),
		}
	default:
		panic(fmt.Sprintf("type %v not supported", typ.Underlying()))
	}
	topLevelDecl.Add(decl)
	return decl
}
