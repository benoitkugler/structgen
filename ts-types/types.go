package tstypes

import (
	"fmt"

	"github.com/benoitkugler/structgen/enums"
	tsEnums "github.com/benoitkugler/structgen/enums/ts"
	"github.com/benoitkugler/structgen/loader"
)

// This file defines a representation of ts types.
// These types are built from go/types,
// and know how to render themselves to .ts code.

var _ loader.Type = tsType(nil)

// tsType is the common interface of all ts types.
type tsType interface {
	Render() []loader.Declaration

	// return the name referencing the type
	name() string
}

// NullableTsType wraps a type, making him nullable.
type NullableTsType struct {
	tsType
}

func (t NullableTsType) name() string {
	return t.tsType.name() + " | null"
}

var timesStringDefinition = loader.Declaration{
	Id: "__times_string_def",
	Content: `
	class DateTag {
		private _ :"D" = "D"
	}
	
	class TimeTag {
		private _ :"T" = "T"
	}
	
	// AAAA-MM-YY date format
	export type Date_ = string & DateTag
	
	// ISO date-time string
	export type Time = string & TimeTag
	`,
}

// one of string, number, boolean
type tsBasic string

const (
	TsString  tsBasic = "string"
	TsNumber  tsBasic = "number"
	TsBoolean tsBasic = "boolean"
	// Represent a go time.Time.
	// An alias will be added
	TsTime tsBasic = "Time"
	// Represent a go time.Time.
	// An alias will be added
	TsDate tsBasic = "Date_"
	TsAny  tsBasic = "any"
)

func (t tsBasic) Render() []loader.Declaration {
	// special case for date and time
	switch t {
	case TsTime, TsDate:
		return []loader.Declaration{timesStringDefinition}
	default:
		return nil
	}
}

func (t tsBasic) name() string { return string(t) }

// TsNamedType represents a defined user type,
// appart from enums and structs.
type TsNamedType struct {
	origin     string
	name_      string
	underlying tsType
}

func (named TsNamedType) Render() []loader.Declaration {
	deps := named.underlying.Render()

	code := fmt.Sprintf(`// %s
	export type %s = %s`, named.origin, named.name_, named.underlying.name())

	deps = append(deps, loader.Declaration{Id: named.name_, Content: code})
	return deps
}

func (t TsNamedType) name() string { return t.name_ }

// TsMap represents a mapping object
type TsMap struct {
	key  tsType
	elem tsType
}

func (t TsMap) Render() []loader.Declaration {
	// the map itself has no additional declarations
	return append(t.key.Render(), t.elem.Render()...)
}

func (t TsMap) name() string {
	return fmt.Sprintf("{ [key: %s]: %s }", t.key.Render(), t.elem.Render())
}

// TsArray represents an array
type TsArray struct {
	elem tsType
}

func (t TsArray) Render() []loader.Declaration {
	// the array itself has no additional declarations
	return t.elem.Render()
}

func (t TsArray) name() string {
	return t.elem.name() + "[]"
}

// TsEnum represents an enum type
type TsEnum struct {
	enum   enums.Type
	origin string
}

func (t TsEnum) Render() []loader.Declaration {
	return []loader.Declaration{{
		Id: t.enum.Name,
		Content: "// " + t.origin + "\n" +
			tsEnums.EnumAsTypeScript(t.enum),
	}}
}

func (t TsEnum) name() string { return t.enum.Name }

// StructField stores one propery of an object
type StructField struct {
	Type tsType
	Name string
}

// TsObject represents an interface
type TsObject struct {
	origin  string
	name_   string
	fields  []StructField
	embeded []tsType
}

func (t TsObject) name() string { return t.name_ }

func (t TsObject) Render() (decls []loader.Declaration) {
	out := "// " + t.origin + "\n"

	if len(t.embeded) == 0 { // prefer interface syntax
		out += fmt.Sprintf("export interface %s {\n", t.name_)
	} else {
		out += fmt.Sprintf("export type %s = {\n", t.name_)
	}

	for _, field := range t.fields {
		decls = append(decls, field.Type.Render()...)
		out += fmt.Sprintf("\t%s: %s,\n", field.Name, field.Type.name())
	}
	out += "}"
	for _, embeded := range t.embeded {
		decls = append(decls, embeded.Render()...)
		out += " & " + embeded.name()
	}

	decls = append(decls, loader.Declaration{Id: t.name_, Content: out})
	return decls
}
