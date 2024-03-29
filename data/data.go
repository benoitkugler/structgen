package data

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/loader"
	"github.com/benoitkugler/structgen/utils"
)

var _ loader.Type = dataFunction(nil)

// dataFunction describe a function
// able to generate random data.
type dataFunction interface {
	// Identify the type generated
	Id() string
	Type() types.Type
	Render() []loader.Declaration
}

func typeName(target string, typ types.Type) string {
	n, _ := utils.TypeName(target, typ)
	return n
}

type FnBasic struct {
	type_ *types.Basic
}

func (f FnBasic) Id() string {
	if f.type_.Name() == "byte" {
		return "uint8"
	}
	return f.type_.Name()
}

func (f FnBasic) Type() types.Type {
	return f.type_
}

func (f FnBasic) Render() []loader.Declaration {
	var code string
	switch f.type_.Kind() {
	case types.Bool:
		code = fnBool()
	case types.Int:
		code = fnInt("int")
	case types.Rune:
		code = fnInt("rune")
	case types.Int64:
		code = fnInt("int64")
	case types.Uint8:
		code = fnInt("uint8")
	case types.Int8:
		code = fnInt("int8")
	case types.Float64:
		code = fnFloat64()
	case types.String:
		code = fnString()
	default:
		panic(fmt.Sprintf("basic type %v not supported", f.type_))
	}
	return []loader.Declaration{{Id: f.Id(), Content: code}}
}

func fnBool() string {
	return `
	func randbool() bool {
		i := rand.Int31n(2)
		return i == 1
	}`
}

func fnInt(intType string) string {
	return fmt.Sprintf(`
	func rand%s() %s {
		return %s(rand.Intn(1000000))
	}`, intType, intType, intType)
}

func fnFloat64() string {
	return `
	func randfloat64() float64 {
		return rand.Float64() * float64(rand.Int31())
	}`
}

func fnString() string {
	return `
	var letterRunes2  = []rune("azertyuiopqsdfghjklmwxcvbn123456789é@!?&èïab ")

	func randstring() string {
		b := make([]rune, 50)
		maxLength := len(letterRunes2)		
		for i := range b {
			b[i] = letterRunes2[rand.Intn(maxLength)]
		}
		return string(b)
	}`
}

type fnTime struct {
	type_ *types.Named
}

func (f fnTime) Id() string {
	return "tTime"
}

func (f fnTime) Type() types.Type {
	return types.NewNamed(types.NewTypeName(0, nil, "Time", nil), &types.Struct{}, nil)
}

func (f fnTime) Render() []loader.Declaration {
	return []loader.Declaration{{Id: f.Id(), Content: `
	func randtTime() time.Time {
		return time.Unix(int64(rand.Int31()), 5)
	}
	`}}
}

type fnArray struct {
	Elem          dataFunction
	TargetPackage string
	Length        int64
}

func (f fnArray) Id() string {
	return fmt.Sprintf("Array%d%s", f.Length, f.Elem.Id())
}

func (f fnArray) Type() types.Type {
	return types.NewArray(f.Elem.Type(), f.Length)
}

func (f fnArray) Render() []loader.Declaration {
	decls := f.Elem.Render() // recurse for deps

	elemString := typeName(f.TargetPackage, f.Elem.Type())
	decl := loader.Declaration{
		Id: f.Id(), Content: fmt.Sprintf(`
	func rand%s() [%d]%s {
		var out [%d]%s
		for i := range out {
			out[i] = rand%s()
		}
		return out
	}`, f.Id(), f.Length, elemString, f.Length, elemString, f.Elem.Id()),
	}

	decls = append(decls, decl)
	return decls
}

type fnSlice struct {
	Elem          dataFunction
	TargetPackage string
}

func (f fnSlice) Id() string {
	return fmt.Sprintf("Slice%s", f.Elem.Id())
}

func (f fnSlice) Type() types.Type {
	return types.NewSlice(f.Elem.Type())
}

func (f fnSlice) Render() []loader.Declaration {
	decls := f.Elem.Render() // recurse for deps

	elemString := typeName(f.TargetPackage, f.Elem.Type())
	decl := loader.Declaration{
		Id: f.Id(), Content: fmt.Sprintf(`
	func rand%s() []%s {
		l := 40 + rand.Intn(10)
		out := make([]%s, l)
		for i := range out {
			out[i] = rand%s()
		}
		return out
	}`, f.Id(), elemString, elemString, f.Elem.Id()),
	}

	decls = append(decls, decl)
	return decls
}

type fnMap struct {
	Key           dataFunction
	Elem          dataFunction
	TargetPackage string
}

func (f fnMap) Id() string {
	return fmt.Sprintf("Map%s%s", f.Key.Id(), f.Elem.Id())
}

func (f fnMap) Type() types.Type {
	return types.NewMap(f.Key.Type(), f.Elem.Type())
}

func (f fnMap) Render() []loader.Declaration {
	decls := f.Key.Render() // recurse for deps
	decls = append(decls, f.Elem.Render()...)

	keyString := typeName(f.TargetPackage, f.Key.Type())
	elemString := typeName(f.TargetPackage, f.Elem.Type())
	decl := loader.Declaration{
		Id: f.Id(), Content: fmt.Sprintf(`
	func rand%s() map[%s]%s {
		l := 40 + rand.Intn(10)
		out := make(map[%s]%s, l)
		for i := 0; i < l; i++ {
			out[rand%s()] = rand%s()
		}
		return out
	}`, f.Id(), keyString, elemString, keyString, elemString, f.Key.Id(), f.Elem.Id()),
	}

	decls = append(decls, decl)
	return decls
}

// structField stores one property of an object
type structField struct {
	type_ dataFunction
	Name  string
	Id    string
}

// fnStruct generate a named random struct
type fnStruct struct {
	TargetPackage string
	Type_         *types.Named
	Fields        []structField
}

func (f fnStruct) Id() string {
	packageName := f.Type_.Obj().Pkg().Name()
	localName := f.Type_.Obj().Name()
	if packageName == f.TargetPackage {
		return localName
	}
	return packageName[:3] + "_" + localName
}

func (f fnStruct) Type() types.Type {
	return f.Type_
}

func (f fnStruct) Render() (decls []loader.Declaration) {
	fieldsCode := ""
	for _, field := range f.Fields {
		decls = append(decls, field.type_.Render()...)
		fieldsCode += fmt.Sprintf("\t%s: rand%s(),\n", field.Name, field.Id)
	}
	typeN := typeName(f.TargetPackage, f.Type_)
	fCode := loader.Declaration{
		Id: f.Id(), Content: fmt.Sprintf(`
	func rand%s() %s {
		return %s{
			%s
		}
	}`, f.Id(), typeN, typeN, fieldsCode),
	}

	decls = append(decls, fCode)
	return decls
}

type fnPointer struct {
	Elem          dataFunction
	TargetPackage string
}

func (f fnPointer) Id() string {
	return fmt.Sprintf("Pointer%s", f.Elem.Id())
}

func (f fnPointer) Type() types.Type {
	return types.NewPointer(f.Elem.Type())
}

func (f fnPointer) Render() []loader.Declaration {
	decls := f.Elem.Render()
	elemString := typeName(f.TargetPackage, f.Elem.Type())
	decl := loader.Declaration{
		Id: f.Id(), Content: fmt.Sprintf(`
	func rand%s() *%s {
		data := rand%s()
		return &data
	}`, f.Id(), elemString, f.Elem.Id()),
	}

	decls = append(decls, decl)
	return decls
}

type fnNamed struct {
	Type_         *types.Named
	Underlying    dataFunction
	TargetPackage string
}

func (f fnNamed) Id() string {
	return f.Type_.Obj().Name()
}

func (f fnNamed) Type() types.Type {
	return f.Type_
}

func (f fnNamed) Render() []loader.Declaration {
	decls := f.Underlying.Render()

	tn := typeName(f.TargetPackage, f.Type_)
	decl := loader.Declaration{
		Id: f.Id(),
		Content: fmt.Sprintf(`
	func rand%s() %s {
		return %s(rand%s())
	}`, f.Id(), tn, tn, f.Underlying.Id()),
	}

	decls = append(decls, decl)
	return decls
}

// fnEnum is also a NamedType but has it's own rand()
// function
type fnEnum struct {
	TargetPackage string
	Type_         *types.Named
	Underlying    enums.Type
}

func (f fnEnum) Id() string {
	return f.Type_.Obj().Name()
}

func (f fnEnum) Type() types.Type {
	return f.Type_
}

func (f fnEnum) Render() []loader.Declaration {
	tn, origin := utils.TypeName(f.TargetPackage, f.Type_)
	if origin == f.TargetPackage {
		origin = ""
	}
	return []loader.Declaration{{
		Id: f.Id(), Content: fmt.Sprintf(`
	func rand%s() %s {
		choix := %s
		i := rand.Intn(len(choix))
		return choix[i]
	}`, f.Id(), tn, f.Underlying.AsArray(origin)),
	}}
}

// fnInterface is a named union type
type fnInterface struct {
	TargetPackage string
	typ_          *types.Named
	members       []dataFunction
}

func (f fnInterface) Id() string {
	return f.typ_.Obj().Name()
}

func (f fnInterface) Type() types.Type {
	return f.typ_
}

func (f fnInterface) Render() []loader.Declaration {
	var (
		choix       []string
		membersDecl []loader.Declaration
	)
	for _, member := range f.members {
		choix = append(choix, fmt.Sprintf("rand%s(),\n", member.Id()))
		membersDecl = append(membersDecl, member.Render()...)
	}

	qualifiedName := typeName(f.TargetPackage, f.typ_)
	return append([]loader.Declaration{{
		Id: f.Id(),
		Content: fmt.Sprintf(`
	func rand%s() %s {
		choix := [...]%s{
			%s
		}
		i := rand.Intn(%d)
		return choix[i]
	}`, f.typ_.Obj().Name(), qualifiedName, qualifiedName,
			strings.Join(choix, ""),
			len(f.members)),
	}}, membersDecl...)
}
