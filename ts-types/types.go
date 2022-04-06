package tstypes

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/benoitkugler/structgen/enums"
	tsEnums "github.com/benoitkugler/structgen/enums/ts"
	"github.com/benoitkugler/structgen/loader"
)

// This file defines a representation of ts types.
// These types are built from go/types,
// and know how to render themselves to .ts code.

var _ loader.Type = Type(nil)

// Type is the common interface of all types
// used in the generated TypeScript code.
type Type interface {
	Render() []loader.Declaration

	// return the Name referencing the type
	Name() string
}

// nullableTsType wraps a type, making him nullable.
type nullableTsType struct {
	Type
}

func (t nullableTsType) Name() string {
	return t.Type.Name() + " | null"
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

func (t tsBasic) Name() string { return string(t) }

// namedType represents a defined user type,
// appart from enums and structs.
type namedType struct {
	underlying Type
	origin     string
	name_      string
}

func (named namedType) Render() []loader.Declaration {
	deps := named.underlying.Render()

	code := fmt.Sprintf(`// %s
	export type %s = %s`, named.origin, named.name_, named.underlying.Name())

	deps = append(deps, loader.Declaration{Id: named.name_, Content: code})
	return deps
}

func (t namedType) Name() string { return t.name_ }

// dict represents a mapping object
type dict struct {
	key  Type
	elem Type
}

func (t dict) Render() []loader.Declaration {
	// the map itself has no additional declarations
	return append(t.key.Render(), t.elem.Render()...)
}

func (t dict) Name() string {
	return fmt.Sprintf("{ [key: %s]: %s }", t.key.Name(), t.elem.Name())
}

// array represents an array
type array struct {
	elem Type
}

func (t array) Render() []loader.Declaration {
	// the array itself has no additional declarations
	return t.elem.Render()
}

func (t array) Name() string {
	return t.elem.Name() + "[]"
}

// enumT represents an enum type
type enumT struct {
	origin string
	enum   enums.Type
}

func (t enumT) Render() []loader.Declaration {
	return []loader.Declaration{{
		Id: t.enum.Name,
		Content: "// " + t.origin + "\n" +
			"export " + tsEnums.EnumAsTypeScript(t.enum),
	}}
}

func (t enumT) Name() string { return t.enum.Name }

// structField stores one propery of an object
type structField struct {
	Type Type
	Name string
}

// class represents an interface
type class struct {
	origin  string
	name_   string
	fields  []structField
	embeded []Type
}

func (t class) Name() string { return t.name_ }

func (t class) Render() (decls []loader.Declaration) {
	out := "// " + t.origin + "\n"

	if len(t.embeded) == 0 { // prefer interface syntax
		out += fmt.Sprintf("export interface %s {\n", t.name_)
	} else {
		out += fmt.Sprintf("export type %s = {\n", t.name_)
	}

	for _, field := range t.fields {
		decls = append(decls, field.Type.Render()...)
		out += fmt.Sprintf("\t%s: %s,\n", field.Name, field.Type.Name())
	}
	out += "}"
	for _, embeded := range t.embeded {
		decls = append(decls, embeded.Render()...)
		out += " & " + embeded.Name()
	}

	decls = append(decls, loader.Declaration{Id: t.name_, Content: out})
	return decls
}

type union struct {
	origin  string
	name_   string
	type_   *types.Interface
	members []Type // completed after analysis
}

func (u *union) Render() []loader.Declaration {
	var (
		members     []string
		kindEnum    []string
		membersDecl []loader.Declaration
	)
	enumKindName := u.name_ + "Kind"
	for i, m := range u.members {
		members = append(members, m.Name())
		kindEnum = append(kindEnum, fmt.Sprintf("%s = %d", m.Name(), i))
		membersDecl = append(membersDecl, m.Render()...)
	}
	code := fmt.Sprintf(`
	export enum %s {
		%s
	}
	
	export interface %s {
		Kind: %s
		Data: %s
	}`, enumKindName, strings.Join(kindEnum, ",\n"), u.name_, enumKindName, strings.Join(members, " | "))

	return append([]loader.Declaration{
		{Id: u.name_, Content: code},
	}, membersDecl...)
}

func (u *union) Name() string { return u.name_ }
