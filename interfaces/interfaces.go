// Package interfaces provides a way to serialize/deserialize
// interface values into JSON.
package interfaces

import (
	"fmt"
	"go/types"
	"sort"

	"github.com/benoitkugler/structgen/loader"
)

// Interface represents a (named) Go interface.
type Interface struct {
	Name    *types.Named   // Go type
	Members []*types.Named // the types implementing this interface, sorted by name
}

func (itf Interface) hasMember(typ *types.Named) bool {
	for _, t := range itf.Members {
		if t == typ {
			return true
		}
	}
	return false
}

// Render returns the JSON routine functions
func (itf Interface) Render() []loader.Declaration {
	return []loader.Declaration{{
		Id:      itf.Name.Obj().Name(),
		Content: itf.json(),
	}}
}

// Implements returns the interface names implemented by `typ`
func (an *Analyzer) Implements(typ *types.Named) (interfaces []string) {
	for _, itf := range an.Itfs() {
		if itf.hasMember(typ) {
			interfaces = append(interfaces, itf.Name.Obj().Name())
		}
	}
	return interfaces
}

var _ loader.Handler = (*handler)(nil)

// Analyzer may be used to handle interface types.
type Analyzer struct {
	pkgNamedTypes []*types.Named
	itfs          map[*types.Named]Interface
}

func NewAnalyser(pkg *types.Scope) *Analyzer {
	return &Analyzer{
		pkgNamedTypes: allNamedTypes(pkg),
		itfs:          make(map[*types.Named]Interface),
	}
}

func (an *Analyzer) Itfs() (out []Interface) {
	for _, v := range an.itfs {
		out = append(out, v)
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Name.Obj().Name() < out[j].Name.Obj().Name() })

	return out
}

type handler struct {
	analyzer    *Analyzer
	packageName string
}

func NewHandler(packageName string, pkg *types.Package) loader.Handler {
	return &handler{packageName: packageName, analyzer: NewAnalyser(pkg.Scope())}
}

func (handler) HandleComment(loader.Comment) error { return nil }

func (h handler) Header() string {
	return fmt.Sprintf(`package %s
	
	import "encoding/json"

	// Code generated by structgen/interfaces. DO NOT EDIT
	`, h.packageName)
}

func (h handler) Footer() string {
	return ""
}

// HandleType implements loader.Handler, but always return a nil value.
func (h *handler) HandleType(typ types.Type) loader.Type {
	switch under := typ.Underlying().(type) {
	case *types.Struct:
		var out class
		for i := 0; i < under.NumFields(); i++ {
			if itf := h.HandleType(under.Field(i).Type()); itf != nil {
				out.fields = append(out.fields, itf)
			}
		}
		return out
	case *types.Slice:
		itf, isItf := h.analyzer.NewInterface(under.Elem())
		if isItf {
			named, ok := typ.(*types.Named)
			if !ok {
				panic(fmt.Sprintf("type []%s is not name", itf.Name.Obj().Name()))
			}
			return itfSlice{name: named.Obj().Name(), elem: itf}
		}
	}
	itf, isItf := h.analyzer.NewInterface(typ)
	if isItf {
		return itf
	}

	return nil
}

// NewInterface adds `typ` to types to analyse.
// It is only useful for Interfaces.
func (an *Analyzer) NewInterface(typ types.Type) (Interface, bool) {
	named, ok := typ.(*types.Named)
	if !ok { // we do not support anonymous interfaces
		return Interface{}, false
	}

	if _, isItf := typ.Underlying().(*types.Interface); !isItf {
		return Interface{}, false
	}

	if itf, has := an.itfs[named]; has {
		return itf, true
	}
	itf := an.processITF(named)
	an.itfs[named] = itf
	return itf, true
}

// // Process uses the accumulated types to find
// // the interfaces and their members.
// func (an Analyzer) Process(pkg *types.Package) Interfaces {
// 	types_ := allNamedTypes(pkg)

// 	var out Interfaces
// 	for namedITF := range an.interfaces {
// 		itf := namedITF.Underlying().(*types.Interface)

// 		item := Interface{Name: namedITF}
// 		for _, t := range types_ {
// 			if types.Implements(t, itf) {
// 				item.Members = append(item.Members, t)
// 			}
// 		}

// 		sort.Slice(item.Members, func(i, j int) bool {
// 			return item.Members[i].Obj().Name() < item.Members[j].Obj().Name()
// 		})

// 		out = append(out, item)
// 	}

// 	sort.Slice(out, func(i, j int) bool { return out[i].Name.Obj().Name() < out[j].Name.Obj().Name() })

// 	return out
// }

// processITF uses the accumulated types to find
// the interfaces and their members.
func (an *Analyzer) processITF(namedITF *types.Named) Interface {
	itf := namedITF.Underlying().(*types.Interface)

	item := Interface{Name: namedITF}
	for _, t := range an.pkgNamedTypes {
		// do not add the interface as member of itself
		if _, isItf := t.Underlying().(*types.Interface); isItf {
			continue
		}

		if types.Implements(t, itf) {
			item.Members = append(item.Members, t)
		}
	}

	sort.Slice(item.Members, func(i, j int) bool {
		return item.Members[i].Obj().Name() < item.Members[j].Obj().Name()
	})
	return item
}

func allNamedTypes(scope *types.Scope) (out []*types.Named) {
	for _, name := range scope.Names() {
		obj, ok := scope.Lookup(name).(*types.TypeName)
		if !ok {
			continue
		}

		if named, isNamed := obj.Type().(*types.Named); isNamed {
			out = append(out, named)
		}
	}
	return out
}
