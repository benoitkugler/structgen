package creation

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/loader"
	"github.com/benoitkugler/structgen/orm"
	"github.com/benoitkugler/structgen/orm/jsonsql"
)

func NewGenHandler(enumsTable enums.EnumTable, eraseJSONDecl bool) loader.Handler {
	return &sqlGenHandler{enumsTable: enumsTable, lookupEnumTable: enumsTable.AsLookupTable(), eraseJSONDecl: eraseJSONDecl}
}

type TableGen struct {
	orm.GoSQLTable
}

// func (t TableGen) Id() string {
// 	return t.Name
// }

func (t TableGen) Render() []loader.Declaration {
	fieldsDecl := make([]string, len(t.Fields))
	fieldsName := make([]string, len(t.Fields))
	for i, f := range t.Fields {
		fieldsDecl[i] = "\t" + f.CreateStmt()
		fieldsName[i] = f.SQLName
	}

	// json validation first
	out := t.jsonValidations()

	decl := loader.Declaration{
		Id: t.Name,
		Content: fmt.Sprintf(`
CREATE TABLE %s (
%s
);`, t.TableName(), strings.Join(fieldsDecl, ",\n")),
	}

	return append(out, decl)
}

// add the json validation functions
func (t TableGen) jsonValidations() []loader.Declaration {
	var out []loader.Declaration
	for _, f := range t.Fields {
		if f.Type.JSON != nil {
			out = append(out, f.Type.JSON.Validations()...)
		}
	}
	return out
}

// encode constraints we want to defer
type constraint interface {
	Render() string
}

type sqlGenHandler struct {
	lookupEnumTable map[string]string // cached from `enumsTable`
	enumsTable      enums.EnumTable
	constraints     []constraint
	eraseJSONDecl   bool
}

func (l sqlGenHandler) Header() string {
	out := `
	-- DO NOT EDIT - autogenerated by structgen 
		   
	`
	if l.eraseJSONDecl {
		out += jsonsql.SetupSQLCode
	}
	return out
}

func (l sqlGenHandler) Footer() string {
	chunks := make([]string, 0, len(l.constraints))
	for _, c := range l.constraints {
		chunks = append(chunks, c.Render())
	}
	return strings.Join(chunks, "\n")
}

func (l *sqlGenHandler) HandleType(typ types.Type) loader.Type {
	table, isTable := orm.TypeToSQLStruct(typ, l.enumsTable)
	if !isTable {
		return nil
	}
	decl := TableGen{GoSQLTable: table}

	// register the constraints
	for _, f := range table.Fields {
		foreignConstraint, has := f.ForeignConstraint(decl.Name)
		if has {
			l.constraints = append(l.constraints, foreignConstraint)
		}
	}

	return decl
}

func (l *sqlGenHandler) HandleComment(comment loader.Comment) error {
	switch comment.Tag {
	case "sql":
	case "noTableSql":
		comment.TypeName = ""
	default:
		return nil
	}
	ct, err := orm.NewConstraint(comment, l.lookupEnumTable)
	l.constraints = append(l.constraints, ct)
	return err
}
