package composites

import (
	"fmt"
	"go/types"
	"io"

	"github.com/benoitkugler/structgen/loader"
	"github.com/benoitkugler/structgen/orm"
)

const PackageComposites = "composites"

var _ loader.Handler = &Composites{} // interface conformity

// Composites create composite types
// and scans function, and write it in a separate package
type Composites struct {
	// package types where extracted from
	OriginPackageName string

	tables []orm.GoSQLTable
}

func (l Composites) WriteHeader(w io.Writer) error {
	_, err := fmt.Fprintf(w, `
	// DON'T EDIT - automatically generated by structgen //

	package %s

	import "database/sql"

	type scanner interface {
		Scan(...interface{}) error
	}
	
	`, PackageComposites)
	return err
}

func (l Composites) WriteFooter(w io.Writer) error {
	g := newGraph(l.tables)
	return g.Render(l.OriginPackageName, w)
}

func (l *Composites) HandleType(topLevelDecl *loader.Declarations, typ types.Type) {
	item, isTable := orm.TypeToSQLStruct(typ)
	if !isTable {
		return
	}
	l.tables = append(l.tables, item)
}

func (l Composites) HandleComment(topLevelDecl *loader.Declarations, comment loader.Comment) error {
	return nil
}