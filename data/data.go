package data

import (
	"fmt"
	"go/types"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/loader"
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

type FnTime struct {
	type_ *types.Named
}

func (f FnTime) Id() string {
	return "tTime"
}

func (f FnTime) Type() types.Type {
	return types.NewNamed(types.NewTypeName(0, nil, "Time", nil), &types.Struct{}, nil)
}

func (f FnTime) Render() []loader.Declaration {
	return []loader.Declaration{{Id: f.Id(), Content: `
	func randtTime() time.Time {
		return time.Unix(int64(rand.Int31()), 5)
	}
	`}}
}

func typeName(targetPackage string, type_ types.Type) string {
	if named, isNamed := type_.(*types.Named); isNamed {
		localName := named.Obj().Name()
		packageName := named.Obj().Pkg().Name()
		if packageName == targetPackage {
			return localName
		}
		return packageName + "." + localName
	}
	return type_.String()
}

type FnArray struct {
	Elem          dataFunction
	TargetPackage string
	Length        int64
}

func (f FnArray) Id() string {
	return fmt.Sprintf("Array%d%s", f.Length, f.Elem.Id())
}

func (f FnArray) Type() types.Type {
	return types.NewArray(f.Elem.Type(), f.Length)
}

func (f FnArray) Render() []loader.Declaration {
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

type FnSlice struct {
	Elem          dataFunction
	TargetPackage string
}

func (f FnSlice) Id() string {
	return fmt.Sprintf("Slice%s", f.Elem.Id())
}

func (f FnSlice) Type() types.Type {
	return types.NewSlice(f.Elem.Type())
}

func (f FnSlice) Render() []loader.Declaration {
	decls := f.Elem.Render() // recurse for deps

	elemString := typeName(f.TargetPackage, f.Elem.Type())
	decl := loader.Declaration{
		Id: f.Id(), Content: fmt.Sprintf(`
	func rand%s() []%s {
		l := rand.Intn(10)
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

type FnMap struct {
	Key           dataFunction
	Elem          dataFunction
	TargetPackage string
}

func (f FnMap) Id() string {
	return fmt.Sprintf("Map%s%s", f.Key.Id(), f.Elem.Id())
}

func (f FnMap) Type() types.Type {
	return types.NewMap(f.Key.Type(), f.Elem.Type())
}

func (f FnMap) Render() []loader.Declaration {
	decls := f.Key.Render() // recurse for deps
	decls = append(decls, f.Elem.Render()...)

	keyString := typeName(f.TargetPackage, f.Key.Type())
	elemString := typeName(f.TargetPackage, f.Elem.Type())
	decl := loader.Declaration{
		Id: f.Id(), Content: fmt.Sprintf(`
	func rand%s() map[%s]%s {
		l := rand.Intn(10)
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

// FnStruct generate a named random struct
type FnStruct struct {
	TargetPackage string
	Type_         *types.Named
	Fields        []structField
}

func (f FnStruct) Id() string {
	packageName := f.Type_.Obj().Pkg().Name()
	localName := f.Type_.Obj().Name()
	if packageName == f.TargetPackage {
		return localName
	}
	return packageName + localName
}

func (f FnStruct) Type() types.Type {
	return f.Type_
}

func (f FnStruct) Render() (decls []loader.Declaration) {
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

type FnPointer struct {
	Elem          dataFunction
	TargetPackage string
}

func (f FnPointer) Id() string {
	return fmt.Sprintf("Pointer%s", f.Elem.Id())
}

func (f FnPointer) Type() types.Type {
	return types.NewPointer(f.Elem.Type())
}

func (f FnPointer) Render() []loader.Declaration {
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

type FnNamed struct {
	Type_         *types.Named
	Underlying    dataFunction
	TargetPackage string
}

func (f FnNamed) Id() string {
	return f.Type_.Obj().Name()
}

func (f FnNamed) Type() types.Type {
	return f.Type_
}

func (f FnNamed) Render() []loader.Declaration {
	decls := f.Underlying.Render()

	tn := typeName(f.TargetPackage, f.Type_)
	decl := loader.Declaration{
		Id: f.Id(), Content: fmt.Sprintf(`
	func rand%s() %s {
		return %s(rand%s())
	}`, f.Id(), tn, tn, f.Underlying.Id()),
	}

	decls = append(decls, decl)
	return decls
}

// FnEnum is also a NamedType but has it's own rand()
// function
type FnEnum struct {
	TargetPackage string
	Type_         *types.Named
	Underlying    enums.Type
}

func (f FnEnum) Id() string {
	return f.Type_.Obj().Name()
}

func (f FnEnum) Type() types.Type {
	return f.Type_
}

func (f FnEnum) Render() []loader.Declaration {
	tn := typeName(f.TargetPackage, f.Type_)
	return []loader.Declaration{{
		Id: f.Id(), Content: fmt.Sprintf(`
	func rand%s() %s {
		choix := %s
		i := rand.Intn(len(choix))
		return choix[i]
	}`, f.Id(), tn, f.Underlying.AsArray()),
	}}
}
