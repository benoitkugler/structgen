package sqltypes

import (
	"go/types"
	"log"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/orm/jsonsql"
	"github.com/benoitkugler/structgen/utils"
)

const JSONB Builtin = "jsonb"

var sqlBasicTypes = map[types.BasicKind]Builtin{
	types.Bool:    "boolean",
	types.Int:     "integer",
	types.Int8:    "integer",
	types.Int16:   "integer",
	types.Int32:   "integer",
	types.Int64:   "integer",
	types.Uint:    "integer",
	types.Uint8:   "integer",
	types.Uint16:  "integer",
	types.Uint32:  "integer",
	types.Uint64:  "integer",
	types.Uintptr: "integer",
	types.Float32: "real",
	types.Float64: "real",
	types.String:  "varchar",
}

func newBuiltin(typ *types.Basic) Builtin {
	kind := typ.Kind()
	sqlType, in := sqlBasicTypes[kind]
	if !in {
		log.Printf("warning : unknow basic type %s, jsonb used", typ)
		sqlType = JSONB
	}
	return sqlType
}

type arrayLike interface {
	Elem() types.Type
}

// only basic array are supported, jsonb is used as fallback
// length = -1 for a slice
func newTypeFromArray(typ arrayLike, length int64) sqlType {
	elem := typ.Elem()
	switch elemTyp := elem.Underlying().(type) {
	case *types.Basic:
		sqlElemType := newBuiltin(elemTyp)
		return Array{element: sqlElemType, length: length}
	default:
		log.Printf("unknow array element type %s, jsonb used for the whole array", elem)
		return JSONB
	}
}

// return the type of a sql.NullXXX struct
// or nil
func isNullable(typ *types.Named) types.Type {
	isFieldValid := func(field *types.Var) bool {
		typ, ok := field.Type().Underlying().(*types.Basic)
		return ok && typ.Info() == types.IsBoolean && field.Name() == "Valid"
	}

	str, ok := typ.Underlying().(*types.Struct)
	if !ok || str.NumFields() != 2 { // not a possible struct
		return nil
	}
	if isFieldValid(str.Field(0)) {
		return str.Field(1).Type()
	} else if isFieldValid(str.Field(1)) {
		return str.Field(0).Type()
	}
	return nil
}

// NewSQLType returns the equivalent SQL type
// Special case for serial (ID) must be handled by the caller
// An enum table is needed to detect the types which should be an enum.
func NewSQLType(typ types.Type, enums enums.EnumTable) SQLType {
	var (
		out               SQLType
		isAnonymousStruct bool
	)

	switch typ := typ.(type) {
	case *types.Basic:
		out = SQLType{Go: typ, Type: newBuiltin(typ), IsNullable: false}
	case *types.Array:
		out = SQLType{Go: typ, Type: newTypeFromArray(typ, typ.Len()), IsNullable: false}
	case *types.Slice:
		// since pq lib convert nil slice to null
		// we have to make this types nullable

		// special case for []byte
		if basic, ok := typ.Elem().(*types.Basic); ok && basic.Kind() == types.Byte {
			out = SQLType{Go: typ, Type: Builtin("bytea"), IsNullable: true}
		} else {
			out = SQLType{Go: typ, Type: newTypeFromArray(typ, -1), IsNullable: true}
		}
	case *types.Named:
		if typ.Obj().Name() == "Date" {
			// special case for Date type
			out = SQLType{Go: typ, Type: Builtin("date"), IsNullable: false}
		} else if utils.IsUnderlyingTime(typ) {
			out = SQLType{Go: typ, Type: Builtin("timestamp (0) with time zone"), IsNullable: true}
		} else if nullableType := isNullable(typ); nullableType != nil {
			out = NewSQLType(nullableType, enums) // convert associated type
			out.IsNullable = true                 // mark as nullable
		} else if enum, basic, isEnum := enums.Lookup(typ); isEnum {
			under := newBuiltin(basic)
			out = SQLType{Go: typ, Type: Enum{underlying: under, Type: enum}}
		} else {
			out = NewSQLType(typ.Underlying(), enums)
			out.GoName = typ.Obj().Name()
		}
		out.Go = typ
	default:
		out = SQLType{Go: typ, Type: JSONB, IsNullable: false}
		_, isAnonymousStruct = typ.(*types.Struct)
	}

	if out.Type == JSONB && !isAnonymousStruct {
		// add the additional validation information
		an := jsonsql.NewAnalyser(enums)
		out.JSON = an.Convert(typ)
	}
	return out
}
