package tstypes

import (
	"fmt"

	"github.com/benoitkugler/structgen/enums"
	tsEnums "github.com/benoitkugler/structgen/enums/ts"
)

// This file defines a representation of ts types.
// These types are built from go/types,
// and know how to render themselves to .ts code.

// tsType is the common interface of all ts types.
type tsType interface {
	Render() string
}

// NullableTsType wraps a type, making him nullable.
type NullableTsType struct {
	tsType
}

func (t NullableTsType) Render() string {
	return t.tsType.Render() + " | null"
}

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

// one of string, number, boolean
type tsBasic string

func (t tsBasic) Render() string {
	return string(t)
}

// TsNamedType represents a defined user type (such as an interface)
type TsNamedType tsBasic

func (t TsNamedType) Render() string {
	return string(t)
}

// TsMap represents a mapping object
type TsMap struct {
	Key  tsType
	Elem tsType
}

func (t TsMap) Render() string {
	return fmt.Sprintf("{ [key: %s]: %s }", t.Key.Render(), t.Elem.Render())
}

// TsArray represents an array
type TsArray struct {
	Elem tsType
}

func (t TsArray) Render() string {
	return t.Elem.Render() + "[]"
}

// TsEnum represents an enum type
type TsEnum enums.Type

func (t TsEnum) Render() string {
	return tsEnums.EnumAsTypeScript(enums.Type(t))
}

// StructField stores one propery of an object
type StructField struct {
	Type tsType
	Name string
}

// TsObject represents an annonymous interface
type TsObject struct {
	Fields  []StructField
	Embeded []tsType
}

func (t TsObject) Render() string {
	out := "{\n"
	for _, field := range t.Fields {
		out += fmt.Sprintf("\t%s: %s,\n", field.Name, field.Type.Render())
	}
	out += "}"
	for _, embeded := range t.Embeded {
		out += " & " + embeded.Render()
	}
	return out
}
