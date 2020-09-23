package tstypes

import (
	"fmt"

	"github.com/benoitkugler/structgen/loader"
)

var _ loader.Declaration = Declaration{}

// Declaration is as top-level type declaration
type Declaration struct {
	Name   string
	Type   TsType
	Origin string
}

func (decl Declaration) Id() string {
	return decl.Name
}

func (decl Declaration) Render() string {
	out := "// " + decl.Origin + "\n"
	object, isInterface := decl.Type.(TsObject)
	if isInterface && len(object.Embeded) == 0 {
		out += fmt.Sprintf("export interface %s %s", decl.Name, decl.Type.Render())
	} else if enum, isEnum := decl.Type.(TsEnum); isEnum {
		out += fmt.Sprintf("export enum %s %s", decl.Name, enum.Render())
	} else {
		out += fmt.Sprintf("export type %s = %s", decl.Name, decl.Type.Render())
	}
	return out
}

type timesStringDefinition struct{}

func (timesStringDefinition) Id() string {
	return "__times_string_def"
}

func (timesStringDefinition) Render() string {
	return `
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
	`
}
