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

func NewHandler(enumsTable enums.EnumTable) *handler {
	return &handler{
		enumsTable:  enumsTable,
		itfs:        interfaces.NewAnalyser(),
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
	out := d.analyseType(typ, nil)
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

type importMap struct {
	goPackage      string
	dartImportPath string
}

// check for tag with the form <name>:<path>
func parseDartExternTag(tag string) (importMap, bool) {
	de := reflect.StructTag(tag).Get("dart-extern")
	if i := strings.Index(de, ":"); i != -1 {
		return importMap{goPackage: de[:i], dartImportPath: de[i+1:]}, true
	}
	return importMap{}, false
}

func (d handler) convertFields(structType *types.Struct) []classField {
	var out []classField
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		tag := structType.Tag(i)
		finalName, exported := utils.GetFieldName(field, tag, "dart")
		if finalName == "" || !exported { // ignored field
			continue
		}

		var imported *importMap
		if extern, ok := parseDartExternTag(tag); ok {
			// do not generate type definition for imported type, rely on external import
			imported = &extern
		}

		// special case for embedded structs : merge the struct fields
		if field.Embedded() {
			// recursion on embeded struct
			dartField := d.analyseType(field.Type(), imported)
			if st, isStruct := dartField.(*class); isStruct {
				out = append(out, st.fields...)
			}
			continue
		}

		dartField := d.analyseType(field.Type(), imported)
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
// if externImport is not nil, external named types (such as client.Answer) are converted to `imported` types
func (d *handler) analyseType(typ types.Type, externImport *importMap) dartType {
	if dt, ok := d.types[typ]; ok {
		return dt
	}
	out := d.createType(typ, externImport)
	// do not register imported as interface
	if _, isImported := out.(imported); !isImported {
		d.itfs.NewInterface(typ)
	}
	d.types[typ] = out
	return out
}

func (d handler) createType(typ types.Type, externImport *importMap) dartType {
	if typ == nil {
		return dartAny
	}
	na, isNamed := typ.(*types.Named)

	// special case for time.Time
	if utils.IsUnderlyingTime(typ) {
		return dartTime
	}

	if isNamed {
		finalName := na.Obj().Name()
		origin := typ.String()

		if externImport != nil { // check for external refs
			if na.Obj().Pkg().Name() == externImport.goPackage {
				return imported{name_: na.Obj().Name(), importPath: externImport.dartImportPath}
			}
		}

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
		underlyingDartType := d.analyseType(typ.Underlying(), externImport)

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
		return d.analyseType(under.Elem(), externImport)
	case *types.Struct:
		out := &class{renderCache: d.renderCache}
		// register the struct before calling convertFields
		// to properly handle recursive types
		d.types[typ] = out
		out.fields = d.convertFields(under)
		return out
	case *types.Array:
		valueDart := d.analyseType(under.Elem(), externImport)
		return list{element: valueDart}
	case *types.Slice:
		valueDart := d.analyseType(under.Elem(), externImport)
		return list{element: valueDart}
	case *types.Map:
		return dict{
			key:     d.analyseType(under.Key(), externImport),
			element: d.analyseType(under.Elem(), externImport),
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

			if dartMember == nil {
				panic("missing member for " + member.String())
			}

			if cl, isClass := dartMember.(*class); isClass {
				cl.interfaces = append(cl.interfaces, dartITF.name_)
			}

			dartITF.members = append(dartITF.members, typeWithTag{
				type_: dartMember,
				tag:   member.Obj().Name(),
			})
		}
	}
}
