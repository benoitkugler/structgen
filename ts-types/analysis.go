package tstypes

import (
	"go/types"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/loader"
	"github.com/benoitkugler/structgen/utils"
)

var _ loader.Handler = handler{}

func NewHandler(enumsTable enums.EnumTable) handler {
	return handler{enumsTable: enumsTable}
}

// stored a map of enum fields
type handler struct {
	enumsTable enums.EnumTable
}

func (d handler) HandleType(typ types.Type) loader.Type {
	return d.analyseType(typ)
}

func (d handler) HandleComment(comment loader.Comment) error { return nil }

func (d handler) Header() string {
	return "// Code generated by structgen DO NOT EDIT"
}
func (d handler) Footer() string { return "" }

func (d handler) convertFields(structType *types.Struct) ([]StructField, []tsType) {
	var (
		out      []StructField
		embedded []tsType
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
			fieldTsType := d.analyseType(field.Type())
			// add it to embeded
			embedded = append(embedded, fieldTsType)
			continue
		}

		tsFieldType := d.analyseType(field.Type())
		out = append(out, StructField{Name: finalName, Type: tsFieldType})
	}
	return out, embedded
}

func analyseBasicType(typ *types.Basic) tsType {
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

// analyseType converts a go type into a ts equivalent
// Named types (such as non-anonymous structs or enums) are extracted into new top levels declarations
func (d handler) analyseType(typ types.Type) tsType {
	if typ == nil {
		return TsAny
	}
	named, isNamed := typ.(*types.Named)

	// special case for Date
	if isNamed && named.Obj().Name() == "Date" {
		return TsDate
	}

	// special case for time.Time
	if utils.IsUnderlyingTime(typ) {
		return TsTime
	}

	if isNamed {
		finalName := named.Obj().Name()
		origin := typ.String()
		// first we look for enums type (which usually have underlying basic types)
		if enum, isEnum := d.enumsTable[finalName]; isEnum {
			return TsEnum{origin: origin, enum: enum}
		}

		// caveat for String, which is reserved
		if finalName == "String" {
			finalName = "String_"
		}

		// otherwise, extract underlying type and look for structs
		underlyingTsType := d.analyseType(typ.Underlying())
		if st, isObject := underlyingTsType.(TsObject); isObject {
			st.origin = origin
			st.name_ = finalName
			return st
		}

		return TsNamedType{origin: origin, name_: finalName, underlying: underlyingTsType}
	}

	switch typ := typ.Underlying().(type) {
	case *types.Basic:
		return analyseBasicType(typ)
	case *types.Pointer:
		// indirection
		return d.analyseType(typ.Elem())
	case *types.Struct:
		fields, embedded := d.convertFields(typ)
		return TsObject{fields: fields, embeded: embedded}
	case *types.Array:
		valueTsType := d.analyseType(typ.Elem())
		return TsArray{elem: valueTsType}
	case *types.Slice:
		valueTsType := d.analyseType(typ.Elem())
		return NullableTsType{tsType: TsArray{elem: valueTsType}}
	case *types.Map:
		return NullableTsType{tsType: TsMap{
			key:  d.analyseType(typ.Key()),
			elem: d.analyseType(typ.Elem()),
		}}
	}
	// unhandled type:
	return TsAny
}
