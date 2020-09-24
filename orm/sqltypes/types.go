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

func (s SQLType) Constraint(field string) string {
	ct := s.Type.Constraint(field)
	if s.IsNullable {
		return ct
	}
	return ct + " NOT NULL"
}

type sqlType interface {
	// Constraint is an optionnal constraint to add to the create statement
	Constraint(field string) string
}

type Builtin string

func (Builtin) Constraint(string) string { return "" }

type Enum struct {
	enums.EnumType
}

func (e Enum) Constraint(field string) string {
	return fmt.Sprintf(" CHECK (%s IN %s)", field, e.AsTuple())
}

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
