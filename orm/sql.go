package orm

import (
	"fmt"
	"go/types"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/utils"
)

type SQLField struct {
	GoName   string
	SQLName  string
	Type     types.Type
	Exported bool
	onDelete string
}

func (s SQLField) IsPrimary() bool {
	return s.GoName == "Id"
}

// ForeignKey returns the name to the table this field references
// or ""
func (s SQLField) ForeignKey() string {
	if !s.IsPrimary() && strings.HasPrefix(s.GoName, "Id") {
		goTableName := strings.TrimPrefix(s.GoName, "Id")
		return tableName(goTableName)
	}
	return ""
}

func convertBasicType(typ *types.Basic) string {
	kind := typ.Kind()
	sqlType, in := sqlBasicTypes[kind]
	if !in {
		log.Printf("warning : unknow basic type %s, jsonb used", typ)
		sqlType = "jsonb"
	}
	return sqlType
}

type arrayLike interface {
	Elem() types.Type
}

// only basic array are supported, jsonb is used as fallback
// length = -1 for a slice
func convertArrayType(typ arrayLike, length int64) string {
	elem := typ.Elem()
	switch elemTyp := elem.Underlying().(type) {
	case *types.Basic:
		lengthS := ""
		if length > 0 {
			lengthS = strconv.Itoa(int(length))
		}
		sqlElemType := convertBasicType(elemTyp)
		return fmt.Sprintf("%s[%s]", sqlElemType, lengthS)
	default:
		log.Printf("unknow array element type %s, jsonb used for the whole array", elem)
		return "jsonb"
	}
}

func isFieldValid(field *types.Var) bool {
	typ, ok := field.Type().Underlying().(*types.Basic)
	return ok && typ.Info() == types.IsBoolean && field.Name() == "Valid"
}

// return the type of a sql.NullXXX struct
// or nil
func isNullable(typ *types.Named) types.Type {
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

func (s SQLField) SQLType() string {
	if s.IsPrimary() {
		return "serial"
	}
	sqlType, isNullable := convertType(s.Type)
	if !isNullable {
		sqlType = sqlType + " NOT NULL"
	}
	return sqlType
}

func (s SQLField) IsNullable() bool {
	_, ok := convertType(s.Type)
	return ok
}

// sql type, is nullable
func convertType(typ types.Type) (string, bool) {
	switch typ := typ.(type) {
	case *types.Basic:
		return convertBasicType(typ), false
	case *types.Array:
		return convertArrayType(typ, typ.Len()), false
	case *types.Slice:
		// since pq lib convert nil slice to null
		// we have to make this types nullable

		// special case for []byte
		if basic, ok := typ.Elem().(*types.Basic); ok && basic.Kind() == types.Byte {
			return "bytea", true
		}
		return convertArrayType(typ, -1), true
	case *types.Named:
		if typ.Obj().Name() == "Date" {
			// special case for Date type
			return "date", false
		} else if utils.IsUnderlyingTime(typ) {
			return "timestamp (0) with time zone", true
		} else if nullableType := isNullable(typ); nullableType != nil {
			out, _ := convertType(nullableType)
			return out, true // mark as nullable
		} else {
			return convertType(typ.Underlying())
		}
	}
	return "jsonb", false
}

func (s SQLField) CreateStmt(enums enums.EnumTable) string {
	var constraint string
	if s.IsPrimary() {
		constraint = " PRIMARY KEY"
	}
	if named, isNamed := s.Type.(*types.Named); isNamed {
		if enum, isEnum := enums[named.Obj().Name()]; isEnum {
			constraint = fmt.Sprintf(" CHECK (%s IN %s)", s.SQLName, enum.AsTuple())
		}
	}
	if array, isArray := s.Type.Underlying().(*types.Array); isArray {
		constraint = fmt.Sprintf(" CHECK (array_length(%s, 1) = %d)", s.SQLName, array.Len())
	}
	// we defer foreign contraints in separate declaration
	return fmt.Sprintf("%s %s%s", s.SQLName, s.SQLType(), constraint)
}

func parseForeignKeyConstraint(fullTag string) string {
	sTag := reflect.StructTag(fullTag)
	return sTag.Get("sql_foreign")
}

type fields []SQLField

// excludes primary key
func (fs fields) NoId() fields {
	var out fields
	for _, f := range fs {
		if !f.IsPrimary() {
			out = append(out, f)
		}
	}
	return out
}

// select foreign keys
func (fs fields) ForeignKeys() fields {
	var out fields
	for _, f := range fs {
		if !f.IsPrimary() && f.ForeignKey() != "" {
			// we found a foreign key
			out = append(out, f)
		}
	}
	return out
}

// exclude private (non exported) fields
func (fs fields) Exported() fields {
	var out fields
	for _, f := range fs {
		if f.Exported {
			out = append(out, f)
		}
	}
	return out
}
