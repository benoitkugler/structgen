// Defines type as understood by SQL
package sqltypes

import (
	"fmt"

	"github.com/benoitkugler/structgen/enums"
)

type SQLType struct {
	IsNullable bool
	Type       sqlType
}

func (s SQLType) Declaration(field string) string {
	ct := s.Type.Constraint(field)
	if !s.IsNullable {
		ct += " NOT NULL"
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
	enums.EnumType
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
