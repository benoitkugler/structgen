// Package interfaces provides a way to serialize/deserialize
// interface values into JSON.
package interfaces

import (
	"bytes"
	"fmt"
	"go/format"
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

// Interfaces exposes the relation between types and interfaces
// in the parsed code.
type Interfaces []Interface

// Implements returns the interface names implemented by `typ`
func (itfs Interfaces) Implements(typ *types.Named) (interfaces []string) {
	for _, itf := range itfs {
		if itf.hasMember(typ) {
			interfaces = append(interfaces, itf.Name.Obj().Name())
		}
	}
	return interfaces
}

var _ loader.Handler = (*handler)(nil)

// Analyzer may be used to handle interface types.
type Analyzer struct {
	interfaces map[*types.Named]bool // with underlying type *types.Interface
	types      map[*types.Named]bool
}

func NewAnalyser() Analyzer {
	return Analyzer{
		interfaces: make(map[*types.Named]bool),
		types:      make(map[*types.Named]bool),
	}
}

type handler struct {
	analyzer Analyzer

	packageName string
}

func NewHandler(packageName string) loader.Handler {
	return handler{packageName: packageName, analyzer: NewAnalyser()}
}

func (handler) HandleComment(loader.Comment) error { return nil }

func (h handler) Header() string {
	return fmt.Sprintf(`package %s
	
	import "encoding/json"

	// Code generated by structgen/interfaces. DO NOT EDIT
	`, h.packageName)
}

func (h handler) Footer() string {
	var out bytes.Buffer
	for _, itf := range h.analyzer.Process() {
		out.WriteString(itf.json())
		out.WriteByte('\n')
	}

	b, err := format.Source(out.Bytes())
	if err != nil {
		return out.String()
	}
	return string(b)
}

// HandleType implements loader.Handler, but always return a nil value.
func (h handler) HandleType(typ types.Type) loader.Type {
	h.analyzer.HandleType(typ)
	return nil
}

// HandleType adds `typ` to types to analyse.
// It is only useful for Named and Interfaces.
func (an Analyzer) HandleType(typ types.Type) {
	named, ok := typ.(*types.Named)
	if !ok { // we do not support anonymous interfaces
		return
	}

	// do not add the interface as member of itself
	if _, isItf := typ.Underlying().(*types.Interface); isItf {
		an.interfaces[named] = true
	} else {
		an.types[named] = true
	}
}

// Process uses the accumulated types to find
// the interfaces and their members.
func (an Analyzer) Process() Interfaces {
	var out Interfaces
	for namedITF := range an.interfaces {
		itf := namedITF.Underlying().(*types.Interface)

		item := Interface{Name: namedITF}
		for t := range an.types {
			if types.Implements(t, itf) {
				item.Members = append(item.Members, t)
			}
		}

		sort.Slice(item.Members, func(i, j int) bool {
			return item.Members[i].Obj().Name() < item.Members[j].Obj().Name()
		})

		out = append(out, item)
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Name.Obj().Name() < out[j].Name.Obj().Name() })

	return out
}
