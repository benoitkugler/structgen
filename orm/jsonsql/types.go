// Analyse types to generate compile-time
// json validation functions for PSQL
package jsonsql

import (
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

	// AddValidation write the needed declaration for the
	// function name
	AddValidation(*loader.Declarations)
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
		if utils.IsUnderlyingTime(t) {
			// special case for time, JSONed as a string
			return String
		}
		return newStruct(t, enums)
	case *types.Named:
		if enum, basic, ok := enums.Lookup(t); ok {
			under := newBasic(basic)
			return enumValue{basic: under, enumType: enum}
		}
		return NewTypeJSON(t.Underlying(), enums)
	default:
		return Dynamic
	}
}

type field struct {
	key   string
	type_ TypeJSON
}

// Struct is a fixed field struct
type Struct struct {
	fields []field
}

func newStruct(t *types.Struct, enums enums.EnumTable) Struct {
	var fields []field
	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)
		key, isExported := utils.GetFieldName(f, t.Tag(i), "json")
		if !isExported {
			continue
		}
		fields = append(fields, field{key: key, type_: NewTypeJSON(f.Type(), enums)})
	}
	return Struct{fields: fields}
}

type Map struct {
	elem TypeJSON
}

func newMap(t *types.Map, enums enums.EnumTable) Map {
	return Map{elem: NewTypeJSON(t.Elem(), enums)}
}

// Array encode a slice or an array (homogenous)
type Array struct {
	length int64 // -1 for slice
	elem   TypeJSON
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
	enumType enums.EnumType
}
