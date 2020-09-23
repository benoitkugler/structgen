package orm

import (
	"go/types"

	"github.com/benoitkugler/structgen/utils"
)

func TypeToSQLStruct(typ types.Type) (GoSQLTable, bool) {
	// we only keep named structs
	if named, isNamed := typ.(*types.Named); isNamed {
		if str, isStruct := named.Underlying().(*types.Struct); isStruct {
			table := NewGoSQLTable(named.Obj().Name(), str)
			return table, true
		}
	}
	return GoSQLTable{}, false
}

type GoSQLTable struct {
	Name   string // local type of the struct
	Fields fields

	uniqueColumns map[string]bool // sql names for unique columns of the table
}

func NewGoSQLTable(name string, type_ *types.Struct) GoSQLTable {
	args := GoSQLTable{Name: name}
	args.Fields = extractStructFields(type_)
	return args
}

func (t GoSQLTable) QualifiedGoName(package_ string) string {
	return package_ + "." + t.Name
}

func (t GoSQLTable) HasID() bool {
	for _, f := range t.Fields {
		if f.IsPrimary() {
			return true
		}
	}
	return false
}

func (m GoSQLTable) Id() string {
	return m.Name
}

func tableName(goName string) string {
	return toSnakeCase(goName) + "s"
}

// TableName returns the sql table name
func (m GoSQLTable) TableName() string {
	return tableName(m.Name)
}

func extractStructFields(type_ *types.Struct) []SQLField {
	var out []SQLField
	for i := 0; i < type_.NumFields(); i++ {
		field := type_.Field(i)
		sqlFieldName, exported := utils.GetFieldName(field, type_.Tag(i), "sql")
		if sqlFieldName == "" { // field ignored
			continue
		}

		// for embedded structs, we flatten the fields
		if underlyingType, isStruct := field.Type().Underlying().(*types.Struct); field.Embedded() && isStruct {
			// extract the fields ...
			sqlFields := extractStructFields(underlyingType)
			// ... and merge them to the outer struct
			out = append(out, sqlFields...)
			continue
		}
		constraint := parseForeignKeyConstraint(type_.Tag(i))
		goFieldName := field.Name()
		sf := SQLField{GoName: goFieldName, SQLName: sqlFieldName, Type: field.Type(), Exported: exported, onDelete: constraint}
		out = append(out, sf)
	}
	return out
}

// SetUniqueColumns use the constraints to add information
// of the column with unique constraint in the table
func (m *GoSQLTable) SetUniqueColumns(constraints map[string][]string) {
	set := map[string]bool{}
	for _, col := range constraints[m.Name] {
		set[col] = true
	}
	m.uniqueColumns = set
}

func (m GoSQLTable) IsColumnUnique(sqlName string) bool {
	return m.uniqueColumns[sqlName]
}
