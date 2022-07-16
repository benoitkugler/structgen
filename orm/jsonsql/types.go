// Analyse types to generate compile-time
// json validation functions for PSQL
package jsonsql

import (
	"fmt"
	"go/types"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/interfaces"
	"github.com/benoitkugler/structgen/loader"
	"github.com/benoitkugler/structgen/utils"
)

const (
	Dynamic basic = "" // type not supported -> no constraint applied
	String  basic = "string"
	Number  basic = "number"
	Boolean basic = "boolean"
)

type TypeJSON interface {
	// Id is used to construct the function name
	Id() string

	// Validation returns the needed declaration for the
	// function name
	Validations() []loader.Declaration
}

// FunctionName returns the name of the validation function
// associated with `t`
func FunctionName(t TypeJSON) string {
	return "structgen_validate_json_" + t.Id()
}

// Analyzer converts a Go type to its json checks.
type Analyzer struct {
	enums       enums.EnumTable
	cache       map[types.Type]TypeJSON
	renderCache map[TypeJSON]bool
}

func NewAnalyser(enums enums.EnumTable) *Analyzer {
	return &Analyzer{
		enums:       enums,
		cache:       make(map[types.Type]TypeJSON),
		renderCache: make(map[TypeJSON]bool),
	}
}

// Convert returns the json-SQL type for `t`.
func (an *Analyzer) Convert(t types.Type) TypeJSON {
	// check for the cache
	if out, ok := an.cache[t]; ok {
		return out
	}
	out := an.create(t)
	an.cache[t] = out
	return out
}

func (an *Analyzer) create(t types.Type) TypeJSON {
	switch t := t.(type) {
	case *types.Basic:
		return newBasic(t)
	case *types.Map:
		return an.newMap(t)
	case *types.Slice:
		return an.newArrayFromSlice(t)
	case *types.Array:
		return an.newArrayFromArray(t)
	case *types.Struct:
		panic(fmt.Sprintf("anonymous struct not supported: %s", t))
	case *types.Named:
		if utils.IsUnderlyingTime(t) {
			// special case for time, JSONed as a string
			return String
		}
		if enum, basic, ok := an.enums.Lookup(t); ok {
			under := newBasic(basic)
			return enumValue{basic: under, enumType: enum}
		} else if st, isStruct := t.Underlying().(*types.Struct); isStruct {
			return an.newStruct(st, t)
		} else if _, isItf := t.Underlying().(*types.Interface); isItf {
			return an.newUnion(t)
		}
		return an.Convert(t.Underlying())
	default:
		return Dynamic
	}
}

type field struct {
	type_ TypeJSON
	key   string
}

// class is a fixed field struct
type class struct {
	name   *types.Named
	fields []field

	renderCache map[TypeJSON]bool // to handle recursive types
}

func (an *Analyzer) newStruct(t *types.Struct, name *types.Named) *class {
	// register the output struct before recursing, to properly handle
	// recursive types
	out := &class{name: name, renderCache: an.renderCache}
	an.cache[name] = out

	var fields []field
	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)
		key, isExported := utils.GetFieldName(f, t.Tag(i), "json")
		if !isExported {
			continue
		}
		fields = append(fields, field{key: key, type_: an.Convert(f.Type())})
	}
	out.fields = fields
	return out
}

func idFromNamed(typ *types.Named) string {
	pkg := typ.Obj().Pkg().Name()[:3]
	return pkg + "_" + typ.Obj().Name()
}

func (b *class) Id() string {
	return idFromNamed(b.name)
}

type Map struct {
	elem TypeJSON
}

func (an *Analyzer) newMap(t *types.Map) Map {
	return Map{elem: an.Convert(t.Elem())}
}

// Array encode a slice or an array (homogenous)
type Array struct {
	elem   TypeJSON
	length int64 // -1 for slice
}

func (an *Analyzer) newArrayFromArray(t *types.Array) Array {
	return Array{length: t.Len(), elem: an.Convert(t.Elem())}
}

func (an *Analyzer) newArrayFromSlice(t *types.Slice) Array {
	return Array{length: -1, elem: an.Convert(t.Elem())}
}

type basic string

func newBasic(t *types.Basic) basic {
	info := t.Info()
	switch {
	case info&types.IsString != 0:
		return String
	case info&types.IsBoolean != 0:
		return Boolean
	case info&types.IsNumeric != 0:
		return Number
	default:
		return Dynamic
	}
}

type enumValue struct {
	basic
	enumType enums.Type
}

type typeWithTag struct {
	type_ TypeJSON
	tag   string
}

// sum type, defined in Go by a closed interface
type union struct {
	itf     interfaces.Interface
	members []typeWithTag // associated with itf.Members
}

func (u union) Id() string {
	return idFromNamed(u.itf.Name)
}

func (an *Analyzer) newUnion(t *types.Named) union {
	ana := interfaces.NewAnalyser()
	itf, _ := ana.NewInterface(t)
	out := union{itf: itf, members: make([]typeWithTag, len(itf.Members))}
	for i, m := range itf.Members {
		out.members[i] = typeWithTag{
			type_: an.Convert(m),
			tag:   m.Obj().Name(),
		}
	}
	return out
}
