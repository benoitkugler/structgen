// Analyse types to generated compile-time
// json validation functions for PSQL
package jsonsql

import (
	"go/types"

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

func funcName(t TypeJSON) string {
	return "structgen_validate_json_" + t.Id()
}

func NewTypeJSON(t types.Type) TypeJSON {
	switch t := t.(type) {
	case *types.Basic:
		return newBasic(t)
	case *types.Map:
		return newMap(t)
	case *types.Slice:
		return newArrayFromSlice(t)
	case *types.Array:
		return newArrayFromArray(t)
	case *types.Struct:
		return newStruct(t)
	case *types.Named:
		return NewTypeJSON(t.Underlying())
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

func newStruct(t *types.Struct) Struct {
	var fields []field
	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)
		key, isExported := utils.GetFieldName(f, t.Tag(i), "json")
		if !isExported {
			continue
		}
		fields = append(fields, field{key: key, type_: NewTypeJSON(f.Type())})
	}
	return Struct{fields: fields}
}

type Map struct {
	elem TypeJSON
}

func newMap(t *types.Map) Map {
	return Map{elem: NewTypeJSON(t.Elem())}
}

// Array encode a slice or an array (homogenous)
type Array struct {
	length int64 // -1 for slice
	elem   TypeJSON
}

func newArrayFromArray(t *types.Array) Array {
	return Array{length: t.Len(), elem: NewTypeJSON(t.Elem())}
}

func newArrayFromSlice(t *types.Slice) Array {
	return Array{length: -1, elem: NewTypeJSON(t.Elem())}
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
