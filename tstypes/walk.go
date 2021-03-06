package tstypes

import (
	"fmt"
	"go/types"
	"io"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/loader"
	"github.com/benoitkugler/structgen/utils"
)

var _ loader.Handler = Driver{}

func NewHandler(enumsTable enums.EnumTable) Driver {
	return Driver{enumsTable: enumsTable}
}

// stored a map of enum fields
type Driver struct {
	enumsTable enums.EnumTable
}

func (d Driver) HandleType(topLevelDecl *loader.Declarations, typ types.Type) {
	d.AnalyseType(topLevelDecl, typ)
}
func (d Driver) HandleComment(topLevelDecl *loader.Declarations, comment loader.Comment) error {
	return nil
}

func (d Driver) WriteHeader(w io.Writer) error {
	_, err := fmt.Fprintln(w, "// DO NOT EDIT -- autogenerated by structgen")
	return err
}
func (d Driver) WriteFooter(w io.Writer) error { return nil }

func (d Driver) convertFields(topLevelDecl *loader.Declarations, structType *types.Struct) ([]StructField, []TsType) {
	var (
		out      []StructField
		embedded []TsType
	)
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		tag := structType.Tag(i)
		finalName, exported := utils.GetFieldName(field, tag, "ts")
		if finalName == "" || !exported { // ignored field
			continue
		}

		// special case for embedded
		if field.Embedded() {
			// recursion on embeded struct
			fieldTsType := d.AnalyseType(topLevelDecl, field.Type())
			// add it to embeded
			embedded = append(embedded, fieldTsType)
			continue
		}

		tsFieldType := d.AnalyseType(topLevelDecl, field.Type())
		out = append(out, StructField{Name: finalName, Type: tsFieldType})
	}
	return out, embedded
}

func analyseBasicType(typ *types.Basic) TsType {
	info := typ.Info()
	if info&types.IsBoolean != 0 {
		return TsBoolean
	} else if info&(types.IsInteger|types.IsFloat) != 0 {
		return TsNumber
	} else if info&types.IsString != 0 {
		return TsString
	} else {
		return TsAny
	}
}

// AnalyseType converts a go type into a ts equivalent
// Named types (such as non-anonymous structs or enums) are extracted into new top levels declarations
func (d Driver) AnalyseType(topLevelDecl *loader.Declarations, typ types.Type) TsType {
	if typ == nil {
		return TsAny
	}
	named, isNamed := typ.(*types.Named)

	// special case for Date
	if isNamed && named.Obj().Name() == "Date" {
		topLevelDecl.Add(timesStringDefinition{})
		return TsDate
	}

	// special case for time.Time
	if utils.IsUnderlyingTime(typ) {
		topLevelDecl.Add(timesStringDefinition{})
		return TsTime
	}

	if isNamed {
		finalName := named.Obj().Name()
		// first we look for enums type (which usually have underlying basic types)
		var underlyingTsType TsType
		if enum, isEnum := d.enumsTable[finalName]; isEnum {
			underlyingTsType = TsEnum(enum)
		}

		// otherwise, extract underlying type
		if underlyingTsType == nil {
			underlyingTsType = d.AnalyseType(topLevelDecl, typ.Underlying())
		}
		// caveat for String, which is reserved
		if finalName == "String" {
			finalName = "String_"
		}
		decl := Declaration{Name: finalName, Type: underlyingTsType, Origin: typ.String()}
		// add top level declaration
		topLevelDecl.Add(decl)
		// return named type
		return TsNamedType(finalName)
	}

	switch typ := typ.Underlying().(type) {
	case *types.Basic:
		return analyseBasicType(typ)
	case *types.Pointer:
		// indirection
		return d.AnalyseType(topLevelDecl, typ.Elem())
	case *types.Struct:
		fields, embedded := d.convertFields(topLevelDecl, typ)
		return TsObject{Fields: fields, Embeded: embedded}
	case *types.Array:
		valueTsType := d.AnalyseType(topLevelDecl, typ.Elem())
		return TsArray{Elem: valueTsType}
	case *types.Slice:
		valueTsType := d.AnalyseType(topLevelDecl, typ.Elem())
		return NullableTsType{TsType: TsArray{Elem: valueTsType}}
	case *types.Map:
		return NullableTsType{TsType: TsMap{
			Key:  d.AnalyseType(topLevelDecl, typ.Key()),
			Elem: d.AnalyseType(topLevelDecl, typ.Elem()),
		}}
	}
	// unhandled type:
	return TsAny
}
