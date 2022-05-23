package crud

import (
	"bytes"
	"fmt"
	"go/types"

	"github.com/benoitkugler/structgen/loader"
	"github.com/benoitkugler/structgen/orm"
)

var _ loader.Handler = &handler{} // interface conformity

type structSQL struct {
	packageName string
	orm.GoSQLTable
}

type structSQLTest struct {
	packageName string
	orm.GoSQLTable
}

// return `true` is typ package name is the current package
func (st structSQL) isTypeLocal(typ *types.Named) bool {
	return typ.Obj().Pkg().Name() == st.packageName
}

func (m structSQL) Render() []loader.Declaration {
	args := m.GoSQLTable

	tmpl := templateStructLink
	if args.HasID() {
		tmpl = templateStructWithID
	}

	var out bytes.Buffer
	if err := templateScan.Execute(&out, args); err != nil {
		panic(err)
	}
	if err := tmpl.Execute(&out, args); err != nil {
		panic(err)
	}

	decls := []loader.Declaration{{Id: m.Id(), Content: out.String()}}

	// generate the JSON Value interface method
	for _, field := range m.Fields {
		if field.Type.JSON != nil {
			goTypeName := field.Type.GoName
			if goTypeName == "" {
				panic(fmt.Sprintf("JSON field %s is not named: SQL Value interface can't be implemented", field.GoName))
			}

			// check if the type is in the same package
			if m.isTypeLocal(field.Type.Go.(*types.Named)) {
				decls = append(decls, loader.Declaration{
					Id: "json_value" + goTypeName,
					Content: fmt.Sprintf(`
					func (s *%s) Scan(src interface{}) error { return loadJSON(s, src) }
					func (s %s) Value() (driver.Value, error) { return dumpJSON(s) }
					`, goTypeName, goTypeName),
				})
			}
		}
	}

	return decls
}

func (m structSQLTest) Render() []loader.Declaration {
	args := m.GoSQLTable
	var out bytes.Buffer
	if err := templateTest.Execute(&out, args); err != nil {
		panic(err)
	}

	return []loader.Declaration{{Id: m.Id(), Content: out.String()}}
}

type handler struct {
	PackageName string

	uniqueConstraints map[string][]string // table name -> unique cols
	tables            []structSQL

	IsTest bool
}

func NewHandler(packageName string, isTest bool) *handler {
	return &handler{PackageName: packageName, IsTest: isTest, uniqueConstraints: make(map[string][]string)}
}

func (l handler) Header() string {
	var dbInterface string
	if !l.IsTest {
		dbInterface = utils + `
		type scanner interface {
			Scan(...interface{}) error
		}

		// DB groups transaction like objects
		type DB interface {
			Exec(query string, args ...interface{}) (sql.Result, error)
			Query(query string, args ...interface{}) (*sql.Rows, error)
			QueryRow(query string, args ...interface{}) *sql.Row 
			Prepare(query string) (*sql.Stmt, error)
		}`
	}

	return fmt.Sprintf(`
	package %s

	// Code generated by structgen. DO NOT EDIT.

	import "database/sql"

	%s 

	`, l.PackageName, dbInterface)
}

func (l handler) Footer() string {
	var out bytes.Buffer
	for _, table := range l.tables {
		table.SetUniqueColumns(l.uniqueConstraints)
		if err := templateSelectBy.Execute(&out, table); err != nil {
			panic(err)
		}
		if table.HasID() { // the lookup methods are only valid for link tables
			continue
		}
		if err := templateStructLinkToLookup.Execute(&out, table); err != nil {
			panic(err)
		}
	}
	return out.String()
}

func (l *handler) HandleType(typ types.Type) loader.Type {
	item, isTable := orm.TypeToSQLStruct(typ, nil)
	if !isTable {
		return nil
	}
	var decl loader.Type
	if l.IsTest {
		decl = structSQLTest{l.PackageName, item}
	} else {
		table := structSQL{l.PackageName, item}
		l.tables = append(l.tables, table)
		decl = table
	}
	return decl
}

func (l handler) HandleComment(comment loader.Comment) error {
	if comment.Tag != "sql" { // ignored
		return nil
	}
	column := orm.IsUniqueConstraint(comment)
	if column != "" { // we have a unique field
		l.uniqueConstraints[comment.TypeName] = append(l.uniqueConstraints[comment.TypeName], column)
	}
	return nil
}
