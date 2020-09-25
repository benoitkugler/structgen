package orm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/benoitkugler/structgen/orm/sqltypes"
)

type SQLField struct {
	GoName   string
	SQLName  string
	Type     sqltypes.SQLType
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

func (s SQLField) CreateStmt() string {
	var typeDecl string
	if s.IsPrimary() {
		typeDecl = "serial PRIMARY KEY"
	} else {
		typeDecl = s.Type.Declaration(s.SQLName)
	}
	// we defer foreign contraints in separate declaration
	return fmt.Sprintf("%s %s", s.SQLName, typeDecl)
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
