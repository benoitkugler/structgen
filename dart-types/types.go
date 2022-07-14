package darttypes

import (
	"fmt"
	"go/types"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/benoitkugler/structgen/enums"
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

func lowerFirst(s string) string {
	return strings.ToLower(s[0:1]) + s[1:]
}

type basic string

func (b basic) name() string       { return string(b) }
func (b basic) functionId() string { return lowerFirst(string(b)) }

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

// convert to dart convention
func (cl classField) dartName() string {
	return lowerFirst(cl.name)
}

type class struct {
	origin     string
	name_      string // needed for constructors
	fields     []classField
	interfaces []string // interfaces implemented

	renderCache map[dartType]bool
}

func (cl class) name() string       { return cl.name_ }
func (cl class) functionId() string { return lowerFirst(cl.name_) }

func (cl *class) Render() (out []loader.Declaration) {
	if cl.renderCache[cl] {
		return nil
	}
	cl.renderCache[cl] = true
	var fields, initFields, interpolatedFields []string
	for _, field := range cl.fields {
		if field.type_.name() != cl.name_ {
			out = append(out, field.type_.Render()...)
		}
		fields = append(fields, fmt.Sprintf("final %s %s;", field.type_.name(), field.dartName()))
		initFields = append(initFields, fmt.Sprintf("this.%s", field.dartName()))
		interpolatedFields = append(interpolatedFields, fmt.Sprintf("$%s", field.dartName()))
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

		const %s(%s);

		@override
		String toString() {
			return "%s(%s)";
		}
		}
		
		%s
	`, cl.origin, cl.name_, implements,
			strings.Join(fields, "\n"), cl.name_, strings.Join(initFields, ", "),
			cl.name_, strings.Join(interpolatedFields, ", "),
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
	return "list" + strings.Title(l.element.functionId())
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
	return "dict" + strings.Title(d.key.functionId()) + strings.Title(d.element.functionId())
}

type enum struct {
	origin string
	enum   enums.Type
}

func (e enum) name() string       { return strings.Title(e.enum.Name) }
func (e enum) functionId() string { return lowerFirst(e.enum.Name) }

func (e enum) Render() []loader.Declaration {
	if e.enum.IsInt {
		// we have to sort by values, which must be ints
		sort.Slice(e.enum.Values, func(i, j int) bool {
			vi, err := strconv.Atoi(e.enum.Values[i].Value)
			if err != nil {
				panic(err)
			}
			vj, err := strconv.Atoi(e.enum.Values[j].Value)
			if err != nil {
				panic(err)
			}
			return vi < vj
		})
	}

	var names, values, labels []string
	for _, v := range e.enum.Values {
		if unicode.IsLower(rune(v.VarName[0])) {
			continue
		}
		names = append(names, lowerFirst(v.VarName))
		labels = append(labels, fmt.Sprintf("%q", v.Label))
		values = append(values, v.Value)
	}

	var fromValue string
	if e.enum.IsInt { // we can just use Dart builtin enums

		fromValue = fmt.Sprintf(`static %s fromValue(int i) {
			return %s.values[i];
		}
		
		int toValue() {
			return index;
		}
		`, e.name(), e.name())
	} else { // add lookup array
		fromValue = fmt.Sprintf(`
		static const _values = [
			%s
		];
		static %s fromValue(String s) {
			return _values.indexOf(s) as %s;
		}
	
		String toValue() {
			return _values[index];
		}
		`, strings.Join(values, ", "), e.name(), e.name())
	}

	// labels are not used for now
	// static const _labels = [
	// 		%s
	// 	];

	// 	String label() { return _labels[index]; }

	enumDecl := fmt.Sprintf(`enum  %s {
		%s
	}
	
	extension _%sExt on %s {
		

		%s
	}
	`, e.name(), strings.Join(names, ", "), e.name(), e.name(), fromValue)

	content := "// " + e.origin + "\n" + enumDecl
	content += "\n" + e.json()

	return []loader.Declaration{{Id: e.enum.Name, Content: content}}
}

type imported struct {
	name_      string
	importPath string
}

func (n imported) name() string       { return n.name_ }
func (n imported) functionId() string { return lowerFirst(n.name_) }

// special cased in Header(); since dart imports must come first
func (n imported) Render() []loader.Declaration { return nil }

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

type typeWithTag struct {
	type_ dartType
	tag   string
}

// interface type, handled as union type
type union struct {
	origin  string
	name_   string
	type_   *types.Interface
	members []typeWithTag // completed after analysis
}

func (u *union) name() string       { return u.name_ }
func (u *union) functionId() string { return lowerFirst(u.name_) }

func (u *union) Render() (out []loader.Declaration) {
	// ensure order
	for _, member := range u.members {
		out = append(out, member.type_.Render()...)
	}

	content := fmt.Sprintf(`abstract class %s {}
	`, u.name_)

	content += u.json()

	out = append(out, loader.Declaration{Id: u.name_, Content: content})
	return out
}
