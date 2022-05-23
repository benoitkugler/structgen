// Defines type as understood by SQL
package sqltypes

import (
	"fmt"
	"go/types"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/orm/jsonsql"
)

type SQLType struct {
	Go         types.Type
	IsNullable bool
	Type       sqlType
	JSON       jsonsql.TypeJSON // might be null
	GoName     string
}

func (s SQLType) Declaration(field string) string {
	ct := s.Type.Constraint(field)
	if !s.IsNullable {
		ct += " NOT NULL"
	}
	// add the eventual JSON validation function
	if s.JSON != nil {
		funcName := jsonsql.FunctionName(s.JSON)
		ct += fmt.Sprintf(" CONSTRAINT %s_%s CHECK (%s(%s))", field, funcName, funcName, field)
	}
	return s.Type.string() + " " + ct
}

type sqlType interface {
	// Constraint is an optionnal constraint to add to the create statement
	Constraint(field string) string
	string() string
}

type Builtin string

func (Builtin) Constraint(string) string { return "" }
func (b Builtin) string() string         { return string(b) }

type Enum struct {
	underlying Builtin
	enums.Type
}

func (e Enum) Constraint(field string) string {
	return fmt.Sprintf(" CHECK (%s IN %s)", field, e.AsTuple())
}
func (e Enum) string() string { return e.underlying.string() }

// Array is a one-dimensionnal SQL array
type Array struct {
	element Builtin
	length  int64 // -1 for a slice
}

func (a Array) Constraint(field string) string {
	if a.length == -1 {
		return ""
	}
	return fmt.Sprintf(" CHECK (array_length(%s, 1) = %d)", field, a.length)
}

func (a Array) string() string {
	return fmt.Sprintf("%s[]", a.element.string())
}
