package tstypes

import (
	"fmt"
	"go/types"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/interfaces"
	"github.com/benoitkugler/structgen/loader"
	"github.com/benoitkugler/structgen/utils"
)

var _ loader.Handler = handler{}

func NewHandler(enumsTable enums.EnumTable, pkg *types.Scope) handler {
	return handler{
		enumsTable:  enumsTable,
		itfs:        interfaces.NewAnalyser(pkg),
		types:       make(map[types.Type]Type),
		renderCache: make(map[Type]bool),
	}
}

// stored a map of enum fields
type handler struct {
	enumsTable enums.EnumTable

	itfs *interfaces.Analyzer
	// mapping from go types to the one generated by the analysis,
	// used in processInterfaces()
	types map[types.Type]Type

	renderCache map[Type]bool
}

func (d handler) HandleType(typ types.Type) loader.Type {
	return d.AnalyseType(typ)
}

// AnalyseType converts a go type into a ts equivalent
// Named types (such as non-anonymous structs or enums) are extracted into new top levels declarations
func (d handler) AnalyseType(typ types.Type) Type {
	if dt, ok := d.types[typ]; ok {
		return dt
	}
	out := d.createType(typ)
	itf, ok := d.itfs.NewInterface(typ)
	if ok {
		// also analyse members
		for _, member := range itf.Members {
			d.AnalyseType(member)
		}
	}
	d.types[typ] = out
	return out
}

func (h handler) processInterfaces() {
	for _, itf := range h.itfs.Itfs() {
		tsITF := h.types[itf.Name].(*union)

		for _, member := range itf.Members {
			tsMember := h.types[member]
			if tsMember == nil {
				panic(fmt.Sprintf("interface member %s not analyzed", member.Obj().Name()))
			}
			tsITF.members = append(tsITF.members, tsMember)
		}
	}
}

func (d handler) HandleComment(comment loader.Comment) error { return nil }

func (d handler) Header() string {
	d.processInterfaces()

	return "// Code generated by structgen DO NOT EDIT"
}
func (d handler) Footer() string { return "" }

func (d handler) convertFields(structType *types.Struct) ([]structField, []Type) {
	var (
		out      []structField
		embedded []Type
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
			fieldTsType := d.createType(field.Type())
			// add it to embeded
			embedded = append(embedded, fieldTsType)
			continue
		}

		tsFieldType := d.createType(field.Type())
		out = append(out, structField{Name: finalName, Type: tsFieldType})
	}
	return out, embedded
}

func analyseBasicType(typ *types.Basic) Type {
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

// createType converts a Go type into a TypeScript equivalent.
// Named types (such as non-anonymous structs or enums) are extracted into new top levels declarations
func (d handler) createType(typ types.Type) Type {
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
			return enumT{origin: origin, enum: enum}
		}

		// handle interface after ending the walk
		if inter, isInter := typ.Underlying().(*types.Interface); isInter {
			// interface are handled after walking
			return &union{
				origin: origin,
				name_:  finalName,
				type_:  inter,
				// members are completed after walking the file
			}
		}

		// caveat for String, which is reserved
		if finalName == "String" {
			finalName = "String_"
		}

		// otherwise, extract underlying type and look for structs
		underlyingTsType := d.AnalyseType(typ.Underlying())
		if st, isObject := underlyingTsType.(*class); isObject {
			st.origin = origin
			st.name_ = finalName
			// return st
		}

		return namedType{origin: origin, name_: finalName, underlying: underlyingTsType}
	}

	switch typ := typ.Underlying().(type) {
	case *types.Basic:
		return analyseBasicType(typ)
	case *types.Pointer:
		// indirection
		return d.AnalyseType(typ.Elem())
	case *types.Struct:
		out := &class{renderCache: d.renderCache}
		// register the struct before calling convertFields
		// to properly handle recursive types
		d.types[typ] = out
		out.fields, out.embeded = d.convertFields(typ)
		return out
	case *types.Array:
		valueTsType := d.AnalyseType(typ.Elem())
		return array{elem: valueTsType}
	case *types.Slice:
		valueTsType := d.AnalyseType(typ.Elem())
		return nullableTsType{Type: array{elem: valueTsType}}
	case *types.Map:
		return nullableTsType{Type: dict{
			key:  d.AnalyseType(typ.Key()),
			elem: d.AnalyseType(typ.Elem()),
		}}
	}
	// unhandled type:
	return TsAny
}
