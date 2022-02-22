package darttypes

import (
	"fmt"
	"go/types"
	"io"

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

func (d handler) HandleType(topLevelDecl *loader.Declarations, typ types.Type) {
	d.analyseType(topLevelDecl, typ)
}

func (d handler) HandleComment(topLevelDecl *loader.Declarations, comment loader.Comment) error {
	return nil
}

func (d handler) WriteHeader(w io.Writer) error {
	_, err := fmt.Fprintln(w, `// Code generated by structgen. DO NOT EDIT
	
	typedef JSON = Map<String, dynamic>; // alias to shorten JSON convertors

	`)
	return err
}
func (d handler) WriteFooter(w io.Writer) error { return nil }

func (d handler) convertFields(topLevelDecl *loader.Declarations, structType *types.Struct) []classField {
	var out []classField
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		tag := structType.Tag(i)
		finalName, exported := utils.GetFieldName(field, tag, "dart")
		if finalName == "" || !exported { // ignored field
			continue
		}

		// special case for embedded structs : merge the struct fields
		if field.Embedded() {
			// recursion on embeded struct
			dartField := d.analyseType(topLevelDecl, field.Type())
			if st, isStruct := dartField.(class); isStruct {
				out = append(out, st.fields...)
			}
			continue
		}

		dartField := d.analyseType(topLevelDecl, field.Type())
		out = append(out, classField{name: finalName, type_: dartField})
	}
	return out
}

func analyseBasicType(typ *types.Basic) basic {
	info := typ.Info()
	if info&types.IsBoolean != 0 {
		return dartBool
	} else if info&types.IsInteger != 0 {
		return dartInt
	} else if info&types.IsFloat != 0 {
		return dartFloat
	} else if info&types.IsString != 0 {
		return dartString
	} else {
		return dartAny
	}
}

// analyseType converts a go type into a ts equivalent
// Named types (such as non-anonymous structs or enums) are extracted into new top levels declarations
func (d handler) analyseType(topLevelDecl *loader.Declarations, typ types.Type) dartType {
	if typ == nil {
		return dartAny
	}
	na, isNamed := typ.(*types.Named)

	// special case for Date
	// if isNamed && named.Obj().Name() == "Date" {
	// 	topLevelDecl.Add(timesStringDefinition{})
	// 	return TsDate
	// }

	// special case for time.Time
	if utils.IsUnderlyingTime(typ) {
		topLevelDecl.Add(jsonFunction{
			name:    string(dartTime),
			content: dartTime.renderJSONconvertors(),
		})
		return dartTime
	}

	if isNamed {
		finalName := na.Obj().Name()
		// first we look for enums type (which usually have underlying basic types)
		var (
			underlyingDartType dartType
			jsonDecl           jsonFunction
		)
		e, isEnum := d.enumsTable[finalName]
		if isEnum {
			underlyingDartType = enum(e)

			jsonDecl = jsonFunction{
				name:    finalName,
				content: enum(e).renderJSONconvertors(),
			}
		}

		// otherwise, extract underlying type
		if !isEnum {
			underlyingDartType = d.analyseType(topLevelDecl, typ.Underlying())
		}

		// type name is required for classes
		cl, isClass := underlyingDartType.(class)
		if isClass {
			cl.name = finalName
			underlyingDartType = cl
		}

		if !isEnum && !isClass {
			jsonDecl = jsonFunction{
				name: finalName,
				content: named(finalName).renderJSONconvertors(
					underlyingDartType.fromJSONBody(),
					underlyingDartType.toJSONBody(),
				),
			}
		}

		decl := declaration{name: finalName, type_: underlyingDartType, origin: typ.String()}
		// add top level declaration
		topLevelDecl.Add(decl)
		topLevelDecl.Add(jsonDecl)
		// return named type
		return named(finalName)
	}

	switch typ := typ.Underlying().(type) {
	case *types.Basic:
		out := analyseBasicType(typ)
		topLevelDecl.Add(jsonFunction{
			name:    out.render(),
			content: out.renderJSONconvertors(),
		})
		return out
	case *types.Pointer:
		// indirection
		return d.analyseType(topLevelDecl, typ.Elem())
	case *types.Struct:
		fields := d.convertFields(topLevelDecl, typ)
		return class{fields: fields}
	case *types.Array:
		valueDart := d.analyseType(topLevelDecl, typ.Elem())
		return list{element: valueDart}
	case *types.Slice:
		valueDart := d.analyseType(topLevelDecl, typ.Elem())
		return list{element: valueDart}
	case *types.Map:
		return dict{
			key:     d.analyseType(topLevelDecl, typ.Key()),
			element: d.analyseType(topLevelDecl, typ.Elem()),
		}
	}
	// unhandled type:
	return dartAny
}
