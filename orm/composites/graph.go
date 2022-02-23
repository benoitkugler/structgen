package composites

import (
	"io"
	"sort"
	"strings"
	"text/template"

	"github.com/benoitkugler/structgen/orm"
)

var templateComposite = template.Must(template.New("").Funcs(orm.FnMap).Parse(`
type {{ .Name }} struct {
	{{ with $array:= .Composition -}}
		{{ index $array 0 }} ` + "`" + `json:"-"` + "`" + `
		{{ index $array 1 }} ` + "`" + `json:"-"` + "`" + `
	{{ end -}}
}

type {{.Name}}s []{{.Name}}

func scanOne{{ .Name }}(row scanner) ({{ .Name }}, error) {
	var s {{ .Name }}
	{{ if .HasPrivateFields }}var dummy interface{}{{ end }}

	err := row.Scan(
		{{ range $i, $table := .Tables }}
				{{- range .Fields -}}
					{{- if .Exported -}}
						&s.{{ $table.Name }}.{{ .GoName }},
					{{- else -}}
						&dummy, 
					{{- end -}}
				{{- end }}
		{{ end }}
	)
	return s, err
}

func Scan{{ .Name }}(row *sql.Row) ({{ .Name }}, error) {
	return scanOne{{ .Name }}(row)
}

func Scan{{ .Name}}s(rs *sql.Rows) ({{.Name}}s , error) {
	var (
		s {{ .Name }}
		err error
	)
	defer func() {
		errClose := rs.Close()
		if err == nil {
			err = errClose
		}
	}()
	structs := make({{.Name}}s , 0, 16)
	for rs.Next() {
		s, err = scanOne{{ .Name }}(rs)
		if err != nil {
			return nil, err
		}
		structs = append(structs, s)
	}
	if err = rs.Err(); err != nil {
		return nil, err
	}
	return structs, nil
}
`))

type compositeTable struct {
	Origin string
	Tables []orm.GoSQLTable
}

func (c compositeTable) Name() string {
	return compositeTableName(c.Tables)
}

// aggregate the names
func compositeTableName(tables []orm.GoSQLTable) string {
	var name string
	for _, table := range tables {
		name += table.Name
	}
	return name
}

// Composition returns :
// 	for 2 tables : the two names
//  for more : the first name, and the composition of the other
func (c compositeTable) Composition() [2]string {
	if len(c.Tables) == 2 {
		return [2]string{c.Tables[0].QualifiedGoName(c.Origin), c.Tables[1].QualifiedGoName(c.Origin)}
	}
	return [2]string{c.Tables[0].QualifiedGoName(c.Origin), compositeTableName(c.Tables[1:])}
}

func (c compositeTable) HasPrivateFields() bool {
	for _, table := range c.Tables {
		for _, field := range table.Fields {
			if !field.Exported {
				return true
			}
		}
	}
	return false
}

type lien struct {
	table   string
	foreign string
}

type graph struct {
	tables map[string]orm.GoSQLTable
	liens  map[lien]bool
}

func newGraph(sts []orm.GoSQLTable) graph {
	out := graph{tables: map[string]orm.GoSQLTable{}, liens: map[lien]bool{}}
	for _, s := range sts {
		sqlName := s.TableName()
		out.tables[sqlName] = s
		fs := s.Fields
		for _, key := range fs.ForeignKeys() {
			out.liens[lien{table: sqlName, foreign: key.ForeignKey()}] = true
		}
	}
	return out
}

// return the next tables
func (g graph) voisins(table string) []string {
	var out sort.StringSlice
	for l := range g.liens {
		if l.table == table {
			out = append(out, l.foreign)
		}
	}
	sort.Sort(out)
	return out
}

func (g graph) addNext(origin string, paths [][]string) [][]string {
	vs := g.voisins(origin)
	var out [][]string
	for _, next := range vs {
		var nextPaths [][]string
		// update paths
		for _, path := range paths {
			nextPaths = append(nextPaths, append(path, next))
		}
		// recurse on next neighbour
		nextPathsRec := g.addNext(next, nextPaths)
		// add to global
		out = append(out, nextPathsRec...)
	}
	// we return subpaths as well
	return append(paths, out...)
}

// we want unicity of set of strings
func hash(s []string) string {
	sort.Strings(append([]string{}, s...)) // copy to avoid mutating
	return strings.Join(s, "")
}

func (g graph) sortedOrigines() []string {
	var keys sort.StringSlice
	for origin := range g.tables {
		keys = append(keys, origin)
	}
	sort.Sort(keys)
	return keys
}

// returns all paths, excepted singletons
// path are trimmed of link tables
func (g graph) extractPaths() [][]string {
	var (
		out  [][]string
		seen = map[string]bool{}
	)
	for _, origin := range g.sortedOrigines() {
		paths := g.addNext(origin, [][]string{{origin}})
		for _, path := range paths {
			var trimmed []string
			for _, table := range path {
				if g.tables[table].HasID() {
					trimmed = append(trimmed, table)
				}
			}
			ha := hash(trimmed)
			if len(trimmed) >= 2 && !seen[ha] {
				seen[ha] = true
				out = append(out, trimmed)
			}
		}
	}
	return out
}

func (g graph) render(origin string, out io.Writer) error {
	paths := g.extractPaths()
	for _, path := range paths {
		args := compositeTable{Origin: origin}
		for _, table := range path {
			args.Tables = append(args.Tables, g.tables[table])
		}
		if err := templateComposite.Execute(out, args); err != nil {
			return err
		}
	}
	return nil
}
