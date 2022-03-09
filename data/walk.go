// This file create helpers functions
// to generate random data for tests
package data

import (
	"fmt"
	"go/types"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/interfaces"
	"github.com/benoitkugler/structgen/loader"
	"github.com/benoitkugler/structgen/utils"
)

var _ loader.Handler = handler{}

type handler struct {
	EnumsTable  enums.EnumTable
	PackageName string

	itfs interfaces.Analyzer
	// mapping from go types to the one generated by the analysis,
	// used in processInterfaces()
	types map[types.Type]dataFunction
}

func NewHandler(packageName string, enums enums.EnumTable) loader.Handler {
	return handler{
		PackageName: packageName,
		EnumsTable:  enums,
		itfs:        interfaces.NewAnalyser(),
		types:       make(map[types.Type]dataFunction),
	}
}

func (d handler) HandleType(typ types.Type) loader.Type {
	d.itfs.HandleType(typ)
	return d.analyseType(typ)
}

func (d handler) analyseType(typ types.Type) dataFunction {
	if dt, ok := d.types[typ]; ok {
		return dt
	}
	out := d.createType(typ)
	d.itfs.HandleType(typ)
	d.types[typ] = out
	return out
}

func (d handler) HandleComment(comment loader.Comment) error { return nil }

func (d handler) Header() string {
	d.processInterfaces()

	return fmt.Sprintf(`package %s
	
	// Code generated by structgen. DO NOT EDIT.

	`, d.PackageName)
}
func (d handler) Footer() string { return "" }

func (d handler) convertFields(structType *types.Struct) (fields []structField) {
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		dataFn := d.analyseType(field.Type())
		fields = append(fields, structField{Name: field.Name(), Id: dataFn.Id(), type_: dataFn})
	}
	return fields
}

func isTime(typ *types.Named) bool {
	return typ.Obj().Pkg().Name() == "time" && typ.Obj().Name() == "Time"
}

// return the corresponding function, as well as all its dependencies.
// deps already contain decl.
func (d handler) createType(typ types.Type) (decl dataFunction) {
	if named, isNamed := typ.(*types.Named); isNamed {
		// special case for structs :
		// we dont generate a random function for underlying type
		if st, isStruct := typ.Underlying().(*types.Struct); isStruct {
			// special case for time.Time, we use a shortcut
			if utils.IsUnderlyingTime(typ) {
				underFn := fnTime{type_: named}
				if isTime(named) {
					decl = underFn
				} else {
					decl = fnNamed{TargetPackage: d.PackageName, Type_: named, Underlying: underFn}
				}
			} else {
				fields := d.convertFields(st)
				decl = fnStruct{TargetPackage: d.PackageName, Type_: named, Fields: fields}
			}
		} else if _, isInterface := typ.Underlying().(*types.Interface); isInterface {
			decl = &fnInterface{TargetPackage: d.PackageName, typ_: named}
		} else if enum, isEnum := d.EnumsTable[named.Obj().Name()]; isEnum {
			decl = fnEnum{TargetPackage: d.PackageName, Type_: named, Underlying: enum}
		} else {
			// extract underlying type
			underFn := d.analyseType(typ.Underlying())
			decl = fnNamed{TargetPackage: d.PackageName, Type_: named, Underlying: underFn}
		}

		// add top level declaration
		return decl
	}

	switch typU := typ.Underlying().(type) {
	case *types.Basic:
		decl = FnBasic{type_: typU}
	case *types.Interface:

	case *types.Pointer:
		// indirection
		elem := d.analyseType(typU.Elem())
		decl = fnPointer{TargetPackage: d.PackageName, Elem: elem}
	case *types.Struct:
		panic("annonymous struct are not supported")
	case *types.Array:
		valueFn := d.analyseType(typU.Elem())
		decl = fnArray{TargetPackage: d.PackageName, Length: typU.Len(), Elem: valueFn}
	case *types.Slice:
		valueFn := d.analyseType(typU.Elem())
		decl = fnSlice{TargetPackage: d.PackageName, Elem: valueFn}
	case *types.Map:
		decl = fnMap{
			TargetPackage: d.PackageName,
			Key:           d.analyseType(typU.Key()),
			Elem:          d.analyseType(typU.Elem()),
		}
	default:
		panic(fmt.Sprintf("type %v not supported", typ.Underlying()))
	}

	return decl
}

func (h handler) processInterfaces() {
	itfs := h.itfs.Process()
	for _, itf := range itfs {
		dartITF := h.types[itf.Name].(*fnInterface)

		for _, member := range itf.Members {
			fnMember := h.types[member]
			dartITF.members = append(dartITF.members, fnMember)
		}
	}
}
