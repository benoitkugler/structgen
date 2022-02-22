package darttypes

import (
	"fmt"
	"strings"

	"github.com/benoitkugler/structgen/enums"
	dartEnums "github.com/benoitkugler/structgen/enums/dart"
)

// defines Dart types and how to render them as Dart code

type dartType interface {
	// return the dart code to define the type
	render() string

	fromJSONBody() string
	toJSONBody() string
}

type basic string

func (b basic) render() string { return string(b) }

const (
	dartString basic = "String"
	dartInt    basic = "int"
	dartFloat  basic = "double"
	dartBool   basic = "bool"
	// Represent a go time.Time.
	// An alias will be added
	dartTime basic = "DateTime"
	// // Represent a go time.Time.
	// // An alias will be added
	// dartDate basic = "Date_"
	dartAny basic = "dynamic"
)

type classField struct {
	type_ dartType
	name  string
}

type class struct {
	name   string // needed for constructors
	fields []classField
}

func (cl class) render() string {
	var fields, initFields []string
	for _, f := range cl.fields {
		fields = append(fields, fmt.Sprintf("final %s %s;", f.type_.render(), f.name))
		initFields = append(initFields, fmt.Sprintf("this.%s", f.name))
	}

	return fmt.Sprintf(`{
		%s

		%s(%s);
	} 
	`,
		strings.Join(fields, "\n"), cl.name, strings.Join(initFields, ", "))
}

type list struct {
	element dartType
}

func (l list) render() string {
	return fmt.Sprintf("List<%s>", l.element.render())
}

type dict struct {
	key     dartType
	element dartType
}

func (d dict) render() string {
	return fmt.Sprintf("Map<%s,%s>", d.key.render(), d.element.render())
}

type enum enums.Type

func (e enum) render() string {
	return dartEnums.EnumAsDart(enums.Type(e))
}

// named refers to a previously declared type
type named string

func (n named) render() string { return string(n) }
