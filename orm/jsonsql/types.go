// Analyse types to generate compile-time
// json validation functions for PSQL
package jsonsql

import (
	"fmt"
	"go/types"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/loader"
	"github.com/benoitkugler/structgen/utils"
)

const (
	Dynamic basic = "" // undefined constraint
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

func NewTypeJSON(t types.Type, enums enums.EnumTable) TypeJSON {
	switch t := t.(type) {
	case *types.Basic:
		return newBasic(t)
	case *types.Map:
		return newMap(t, enums)
	case *types.Slice:
		return newArrayFromSlice(t, enums)
	case *types.Array:
		return newArrayFromArray(t, enums)
	case *types.Struct:
		panic(fmt.Sprintf("anonymous struct not supported: %s", t))
	case *types.Named:
		if utils.IsUnderlyingTime(t) {
			// special case for time, JSONed as a string
			return String
		}
		if enum, basic, ok := enums.Lookup(t); ok {
			under := newBasic(basic)
			return enumValue{basic: under, enumType: enum}
		} else if st, isStruct := t.Underlying().(*types.Struct); isStruct {
			return newStruct(st, enums, t)
		}
		return NewTypeJSON(t.Underlying(), enums)
	default:
		return Dynamic
	}
}

type field struct {
	type_ TypeJSON
	key   string
}

// Struct is a fixed field struct
type Struct struct {
	name   *types.Named
	fields []field
}

func newStruct(t *types.Struct, enums enums.EnumTable, name *types.Named) Struct {
	var fields []field
	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)
		key, isExported := utils.GetFieldName(f, t.Tag(i), "json")
		if !isExported {
			continue
		}
		fields = append(fields, field{key: key, type_: NewTypeJSON(f.Type(), enums)})
	}
	return Struct{fields: fields, name: name}
}

func (b Struct) Id() string {
	pkg := b.name.Obj().Pkg().Name()[:3]
	return pkg + "_" + b.name.Obj().Name()
}

type Map struct {
	elem TypeJSON
}

func newMap(t *types.Map, enums enums.EnumTable) Map {
	return Map{elem: NewTypeJSON(t.Elem(), enums)}
}

// Array encode a slice or an array (homogenous)
type Array struct {
	elem   TypeJSON
	length int64 // -1 for slice
}

func newArrayFromArray(t *types.Array, enums enums.EnumTable) Array {
	return Array{length: t.Len(), elem: NewTypeJSON(t.Elem(), enums)}
}

func newArrayFromSlice(t *types.Slice, enums enums.EnumTable) Array {
	return Array{length: -1, elem: NewTypeJSON(t.Elem(), enums)}
}

type basic string

func newBasic(t *types.Basic) basic {
	switch t.Info() {
	case types.IsString:
		return String
	case types.IsBoolean:
		return Boolean
	case types.IsFloat, types.IsInteger:
		return Number
	default:
		return Dynamic
	}
}

type enumValue struct {
	basic
	enumType enums.Type
}
