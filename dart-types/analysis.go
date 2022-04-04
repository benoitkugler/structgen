package darttypes

import (
	"fmt"
	"go/types"
	"reflect"
	"sort"
	"strings"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/interfaces"
	"github.com/benoitkugler/structgen/loader"
	"github.com/benoitkugler/structgen/utils"
)

var _ loader.Handler = (*handler)(nil)

func NewHandler(enumsTable enums.EnumTable, pkg *types.Package) *handler {
	return &handler{
		enumsTable:  enumsTable,
		itfs:        interfaces.NewAnalyser(pkg.Scope()),
		types:       make(map[types.Type]dartType),
		renderCache: make(map[dartType]bool),
	}
}

// stored a map of enum fields
type handler struct {
	enumsTable enums.EnumTable

	itfs *interfaces.Analyzer
	// mapping from go types to the one generated by the analysis,
	// used in processInterfaces()
	types map[types.Type]dartType

	renderCache map[dartType]bool
}

func (d *handler) HandleType(typ types.Type) loader.Type {
	out := d.analyseType(typ)
	return out
}

func (d handler) HandleComment(comment loader.Comment) error {
	return nil
}

func (d *handler) Header() string {
	d.processInterfaces()

	imports := d.processImported()

	return fmt.Sprintf(`// Code generated by structgen. DO NOT EDIT
	
	%s 

	typedef JSON = Map<String, dynamic>; // alias to shorten JSON convertors

	`, imports)
}

func (d *handler) processImported() string {
	paths := map[string]bool{}
	for _, dart := range d.types {
		if imp, isImported := dart.(imported); isImported {
			importLine := fmt.Sprintf("import '%s';", imp.importPath)
			paths[importLine] = true
		}
	}
	var sorted []string
	for p := range paths {
		sorted = append(sorted, p)
	}
	sort.Strings(sorted)
	return strings.Join(sorted, "\n")
}

func (d handler) Footer() string { return "" }

func (d handler) convertFields(structType *types.Struct) []classField {
	var out []classField
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		tag := structType.Tag(i)
		finalName, exported := utils.GetFieldName(field, tag, "dart")
		if finalName == "" || !exported { // ignored field
			continue
		}

		if dartExtern := reflect.StructTag(tag).Get("dart-extern"); dartExtern != "" {
			// do not generate type definition for the field type, rely on external import
			na, ok := field.Type().(*types.Named)
			if !ok {
				panic("dart-extern only works with named types")
			}

			ty := d.analyseImportedType(na, dartExtern)

			out = append(out, classField{name: finalName, type_: ty})
			continue
		}

		// special case for embedded structs : merge the struct fields
		if field.Embedded() {
			// recursion on embeded struct
			dartField := d.analyseType(field.Type())
			if st, isStruct := dartField.(*class); isStruct {
				out = append(out, st.fields...)
			}
			continue
		}

		dartField := d.analyseType(field.Type())
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
func (d *handler) analyseType(typ types.Type) dartType {
	if dt, ok := d.types[typ]; ok {
		return dt
	}
	out := d.createType(typ)
	d.itfs.NewInterface(typ)
	d.types[typ] = out
	return out
}

func (d handler) analyseImportedType(na *types.Named, externPath string) dartType {
	if dt, ok := d.types[na]; ok {
		return dt
	}
	out := imported{name_: na.Obj().Name(), importPath: externPath}
	d.types[na] = out
	return out
}

func (d handler) createType(typ types.Type) dartType {
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
		return dartTime
	}

	if isNamed {
		finalName := na.Obj().Name()
		origin := typ.String()

		// first we look for enums type (which usually have underlying basic types)
		e, isEnum := d.enumsTable[finalName]
		if isEnum {
			return enum{origin: origin, enum: e}
		}

		finalName = strings.Title(finalName)

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

		// otherwise, extract underlying type
		underlyingDartType := d.analyseType(typ.Underlying())

		// type name is required for classes
		cl, isClass := underlyingDartType.(*class)
		if isClass {
			cl.origin = origin
			cl.name_ = finalName
			return cl
		}

		return named{origin: origin, name_: finalName, underlying: underlyingDartType}
	}

	switch under := typ.Underlying().(type) {
	case *types.Basic:
		return analyseBasicType(under)
	case *types.Pointer:
		// indirection
		return d.analyseType(under.Elem())
	case *types.Struct:
		out := &class{renderCache: d.renderCache}
		// register the struct before calling convertFields
		// to properly handle recursive types
		d.types[typ] = out
		out.fields = d.convertFields(under)
		return out
	case *types.Array:
		valueDart := d.analyseType(under.Elem())
		return list{element: valueDart}
	case *types.Slice:
		valueDart := d.analyseType(under.Elem())
		return list{element: valueDart}
	case *types.Map:
		return dict{
			key:     d.analyseType(under.Key()),
			element: d.analyseType(under.Elem()),
		}
	}
	// unhandled type:
	return dartAny
}

func (h *handler) processInterfaces() {
	for _, itf := range h.itfs.Itfs() {
		dartITF := h.types[itf.Name].(*union)

		for _, member := range itf.Members {
			dartMember := h.types[member]
			if cl, isClass := dartMember.(*class); isClass {
				cl.interfaces = append(cl.interfaces, dartITF.name_)
			}

			dartITF.members = append(dartITF.members, dartMember)
		}
	}
}
