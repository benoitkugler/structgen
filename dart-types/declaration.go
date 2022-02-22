package darttypes

import (
	"fmt"

	"github.com/benoitkugler/structgen/loader"
)

var _ loader.Declaration = declaration{}

// declaration is as top-level type declaration
type declaration struct {
	name   string
	type_  dartType
	origin string
}

func (decl declaration) Id() string {
	return decl.name
}

func (decl declaration) Render() string {
	out := "// " + decl.origin + "\n"

	if object, isClass := decl.type_.(class); isClass {
		out += fmt.Sprintf(`class %s %s 
		%s
		`, decl.name, object.render(), object.renderJSONconvertors())
	} else if _, isEnum := decl.type_.(enum); isEnum { // named already defined
		out += decl.type_.render()
	} else {
		out += fmt.Sprintf("typedef %s = %s;", decl.name, decl.type_.render())
	}
	return out
}

type jsonFunction struct {
	name    string
	content string
}

func (jf jsonFunction) Id() string {
	return "_json" + jf.name
}

func (jf jsonFunction) Render() string { return jf.content }
