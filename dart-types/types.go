package darttypes

import (
	"fmt"
	"go/types"
	"sort"
	"strings"

	"github.com/benoitkugler/structgen/enums"
	dartEnums "github.com/benoitkugler/structgen/enums/dart"
	"github.com/benoitkugler/structgen/loader"
)

// defines Dart types and how to render them as Dart code

var _ loader.Type = dartType(nil)

type dartType interface {
	// return the dart code to define the type
	Render() []loader.Declaration

	name() string       // how to refer to this type
	functionId() string // how to call function associated with this type
}

type basic string

func (b basic) name() string       { return string(b) }
func (b basic) functionId() string { return string(b) }

func (b basic) Render() []loader.Declaration {
	return []loader.Declaration{{Id: b.name() + "json", Content: b.json()}}
}

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
	origin     string
	name_      string // needed for constructors
	fields     []classField
	interfaces []string // interfaces implemented
}

func (cl class) name() string       { return cl.name_ }
func (cl class) functionId() string { return cl.name_ }

func (cl *class) Render() (out []loader.Declaration) {
	var fields, initFields []string
	for _, field := range cl.fields {
		out = append(out, field.type_.Render()...)
		fields = append(fields, fmt.Sprintf("final %s %s;", field.type_.name(), field.name))
		initFields = append(initFields, fmt.Sprintf("this.%s", field.name))
	}

	var implements string
	if len(cl.interfaces) != 0 {
		implements = "implements " + strings.Join(cl.interfaces, ", ")
	}

	decl := loader.Declaration{
		Id: cl.name_, Content: fmt.Sprintf(`
		// %s
		class %s %s {
		%s

		%s(%s);
		}
		
		%s
	`, cl.origin, cl.name_, implements,
			strings.Join(fields, "\n"), cl.name_, strings.Join(initFields, ", "),
			cl.json(),
		),
	}
	out = append(out, decl)

	return out
}

type list struct {
	element dartType
}

func (l list) name() string {
	return fmt.Sprintf("List<%s>", l.element.name())
}

func (l list) functionId() string {
	return "List_" + l.element.functionId()
}

func (l list) Render() []loader.Declaration {
	out := l.element.Render()
	out = append(out, loader.Declaration{
		Id:      l.functionId(),
		Content: l.json(),
	})
	return out
}

type dict struct {
	key     dartType
	element dartType
}

func (d dict) name() string {
	return fmt.Sprintf("Map<%s,%s>", d.key.name(), d.element.name())
}

func (d dict) Render() []loader.Declaration {
	out := append(d.key.Render(), d.element.Render()...)

	out = append(out, loader.Declaration{
		Id:      d.functionId(),
		Content: d.json(),
	})
	return out
}

func (d dict) functionId() string {
	return "Dict_" + d.key.functionId() + "_" + d.element.functionId()
}

type enum struct {
	origin string
	enum   enums.Type
}

func (e enum) name() string       { return e.enum.Name }
func (e enum) functionId() string { return e.enum.Name }

func (e enum) Render() []loader.Declaration {
	content := "// " + e.origin + "\n" + dartEnums.EnumAsDart(e.enum)
	content += "\n" + e.json()

	return []loader.Declaration{{Id: e.enum.Name, Content: content}}
}

// named refers to a previously declared type
type named struct {
	underlying dartType
	origin     string
	name_      string
}

func (n named) name() string       { return string(n.name_) }
func (n named) functionId() string { return n.underlying.functionId() }

func (n named) Render() []loader.Declaration {
	out := n.underlying.Render()

	content := "// " + n.origin + "\n"
	content += fmt.Sprintf("typedef %s = %s;\n", n.name_, n.underlying.name())
	content += n.json()

	out = append(out, loader.Declaration{Id: n.name_, Content: content})
	return out
}

// interface type, handled as union type
type union struct {
	origin  string
	name_   string
	type_   *types.Interface
	members []dartType // completed after analysis
}

func (u *union) name() string       { return u.name_ }
func (u *union) functionId() string { return u.name_ }

func (u *union) Render() (out []loader.Declaration) {
	// ensure order
	sort.Slice(u.members, func(i, j int) bool { return u.members[i].name() < u.members[j].name() })
	for _, member := range u.members {
		out = append(out, member.Render()...)
	}

	content := fmt.Sprintf(`// Corresponding Go code
	/*
	%s
	*/ 
	abstract class %s {}
	`, u.goJSON(), u.name_)

	content += u.json()

	out = append(out, loader.Declaration{Id: u.name_, Content: content})
	return out
}
