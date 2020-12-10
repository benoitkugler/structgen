package crud

import (
	"text/template"

	"github.com/benoitkugler/structgen/orm"
)

var (
	templateScan = template.Must(template.New("").Funcs(orm.FnMap).Parse(`
func scanOne{{ .Name }}(row scanner) ({{ .Name }}, error) {
	var s {{.Name}}
	err := row.Scan({{range .Fields}}
		&s.{{.GoName}},{{end}}
	)
	return s, err
}

func Scan{{ .Name }}(row *sql.Row) ({{.Name}}, error) {
	return scanOne{{ .Name }}(row)
}

func SelectAll{{ .Name }}s(tx DB) ({{ .Name }}s, error) {
	rows, err := tx.Query("SELECT * FROM {{snake .Name}}s")
	if err != nil {
		return nil, err
	}
	return Scan{{ .Name }}s(rows)
}
`))

	templateStructWithID = template.Must(template.New("").Funcs(orm.FnMap).Parse(`

// Select{{ .Name }} returns the entry matching id.
func Select{{ .Name }}(tx DB, id int64) ({{ .Name }}, error) {
	row := tx.QueryRow("SELECT * FROM {{snake .Name}}s WHERE id = $1", id)
	return Scan{{ .Name }}(row)
}

// Select{{ .Name }}s returns the entry matching the given ids.
func Select{{ .Name }}s(tx DB, ids ...int64) ({{ .Name }}s, error) {
	rows, err := tx.Query("SELECT * FROM {{snake .Name}}s WHERE id = ANY($1)", pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	return Scan{{ .Name }}s(rows)
}

type {{.Name}}s map[int64]{{.Name}}

func (m {{.Name}}s) Ids() Ids {
	out := make(Ids, 0, len(m))
	for i := range m {
		out = append(out, i)
	}
	return out
}

func Scan{{ .Name }}s(rs *sql.Rows) ({{.Name}}s, error) {
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
	structs := make({{.Name}}s,  16)
	for rs.Next() {
		s, err = scanOne{{ .Name }}(rs)
		if err != nil {
			return nil, err
		}
		structs[s.Id] = s
	}
	if err = rs.Err(); err != nil {
		return nil, err
	}
	return structs, nil
}

// Insert {{ .Name }} in the database and returns the item with id filled.
func (item {{ .Name }}) Insert(tx DB) (out {{.Name}}, err error) {
	row := tx.QueryRow(` + "`" + `INSERT INTO {{snake .Name}}s (
		{{range $i, $e :=  .Fields.Exported.NoId }}{{if $i}},{{end}}{{ $e.SQLName }}{{end}}
		) VALUES (
		{{range $i, $e :=  .Fields.Exported.NoId }}{{if $i}},{{end}}${{inc $i}}{{end}}
		) RETURNING 
		{{range $i, $e := .Fields}}{{if $i}},{{end}}{{ $e.SQLName }}{{end}};
		` + "`" + `{{range  .Fields.Exported.NoId }},item.{{.GoName}}{{end}})
	return Scan{{ .Name }}(row)
}

// Update {{ .Name }} in the database and returns the new version.
func (item {{ .Name }}) Update(tx DB) (out {{.Name}}, err error) {
	row := tx.QueryRow(` + "`" + `UPDATE {{snake .Name}}s SET (
		{{range $i, $e := .Fields.Exported.NoId }}{{if $i}},{{end}}{{ $e.SQLName }}{{end}}
		) = (
		{{range $i, $e := .Fields.Exported.NoId }}{{if $i}},{{end}}${{inc (inc $i)}}{{end}}
		) WHERE id = $1 RETURNING 
		{{range $i, $e := .Fields }}{{if $i}},{{end}}{{ $e.SQLName }}{{end}};
		` + "`" + `{{range .Fields.Exported }},item.{{.GoName}}{{end}})
	return Scan{{ .Name }}(row)
}

// Deletes the {{ .Name }} and returns the item
func Delete{{ .Name }}ById(tx DB, id int64) ({{ .Name }}, error) {
	row := tx.QueryRow("DELETE FROM {{snake .Name}}s WHERE id = $1 RETURNING *;", id)
	return Scan{{ .Name }}(row)
}

// Deletes the {{ .Name }} in the database and returns the ids.
func Delete{{ .Name }}sByIds(tx DB, ids ...int64) (Ids, error) {
	rows, err := tx.Query("DELETE FROM {{ snake .Name }}s WHERE id = ANY($1) RETURNING id", pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	return ScanIds(rows)
}	
`))

	templateStructLink = template.Must(template.New("").Funcs(orm.FnMap).Parse(`
type {{.Name}}s []{{.Name}}

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

// Insert the links {{ .Name}} in the database.
func InsertMany{{ .Name}}s(tx *sql.Tx, items ...{{ .Name}}) error {
	if len(items) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("{{snake .Name}}s", 
		{{range .Fields.Exported }}"{{ .SQLName }}",{{end}}
	))
	if err != nil {
		return err
	}

	for _, item := range items {
		_, err = stmt.Exec({{range $i, $e := .Fields.Exported }}{{if $i}},{{end}}item.{{.GoName}}{{end}})
		if err != nil {
			return err
		}
	}

	if _, err = stmt.Exec(); err != nil {
		return err
	}
	
	if err = stmt.Close(); err != nil {
		return err
	}
	return nil
}

// Delete the link {{ .Name }} in the database.
// Only the {{range .Fields.ForeignKeys }}'{{ .GoName }}' {{end}}fields are used.
func (item {{ .Name }}) Delete(tx DB) error {
	_, err := tx.Exec(` + "`" + `DELETE FROM {{snake .Name}}s WHERE 
	{{range $i, $e := .Fields.ForeignKeys }}{{if $i}} AND {{end}}
	{{- if $e.Type.IsNullable -}}
		( {{ $e.SQLName }} IS NULL OR {{ $e.SQLName }} = ${{inc $i}})
	{{- else -}}
		{{ $e.SQLName }} = ${{inc $i}}
	{{- end -}}	
	{{end}};` +
		"`" + ` {{range .Fields.ForeignKeys }},item.{{.GoName}}{{end}})
	return err
}

`))

	templateStructLinkToLookup = template.Must(template.New("").Funcs(orm.FnMap).Parse(`
{{range .Fields.ForeignKeys }}
	{{ if .Type.IsNullable }}
	{{ else }}
		// By{{ .GoName }} returns a map with '{{ .GoName }}' as keys.
		{{- if $.IsColumnUnique .SQLName }}
		func (items {{$.Name}}s) By{{ .GoName }}() map[int64]{{ $.Name }} {
			out := make(map[int64]{{ $.Name }}, len(items))
			for _, target := range items {
				out[target.{{ .GoName }}] = target
			}
			return out
		}	
		{{ else }}
		func (items {{$.Name}}s) By{{ .GoName }}() map[int64]{{ $.Name }}s {
			out := make(map[int64]{{ $.Name }}s)
			for _, target := range items {
				out[target.{{ .GoName }}] = append(out[target.{{ .GoName }}], target)
			}
			return out
		}	
		{{ end}}
	{{ end }}
{{end}}`))

	templateSelectBy = template.Must(template.New("").Funcs(orm.FnMap).Parse(`
{{range .Fields.ForeignKeys }}
{{- if $.IsColumnUnique .SQLName }}
// Select{{ $.Name }}By{{ .GoName }} return zero or one item, thanks to a UNIQUE constraint
func Select{{ $.Name }}By{{ .GoName }}(tx DB, {{ varname .GoName }} int64) (item {{ $.Name }}, found bool, err error) {
	row := tx.QueryRow("SELECT * FROM {{ snake $.Name }}s WHERE {{ .SQLName }} = $1", {{ varname .GoName}})
	item, err = Scan{{ $.Name }}(row)
	if err == sql.ErrNoRows {
		return item, false, nil
	}
	return item, true, err
}	
{{ end }}

func Select{{ $.Name }}sBy{{ .GoName }}s(tx DB, {{ varname .GoName }}s ...int64) ({{ $.Name }}s, error) {
	rows, err := tx.Query("SELECT * FROM {{ snake $.Name }}s WHERE {{ .SQLName }} = ANY($1)", pq.Int64Array({{ varname .GoName}}s))
	if err != nil {
		return nil, err
	}
	return Scan{{ $.Name }}s(rows)
}	

{{ if $.HasID }}
func Delete{{ $.Name }}sBy{{ .GoName }}s(tx DB, {{ varname .GoName }}s ...int64) (Ids, error) {
	rows, err := tx.Query("DELETE FROM {{ snake $.Name }}s WHERE {{ .SQLName }} = ANY($1) RETURNING id", pq.Int64Array({{ varname .GoName}}s))
	if err != nil {
		return nil, err
	}
	return ScanIds(rows)
}	
{{ else }}
func Delete{{ $.Name }}sBy{{ .GoName }}s(tx DB, {{ varname .GoName }}s ...int64) ({{ $.Name }}s, error)  {
	rows, err := tx.Query("DELETE FROM {{ snake $.Name }}s WHERE {{ .SQLName }} = ANY($1) RETURNING *", pq.Int64Array({{ varname .GoName}}s))
	if err != nil {
		return nil, err
	}
	return Scan{{ $.Name }}s(rows)
}	
{{ end }}

{{end}}`))

	templateTest = template.Must(template.New("").Funcs(orm.FnMap).Parse(`
func queries{{.Name}}(tx *sql.Tx, item {{.Name}}) ({{.Name}}, error) {
	{{ if .HasID }} item, err := item.Insert(tx)
	{{ else }} err := InsertMany{{ .Name }}s(tx, item) {{end}}
	if err != nil {
		return item, err
	}
	rows, err := tx.Query("SELECT * FROM {{snake .Name}}s")
	if err != nil {
		return item, err
	}
	items, err := Scan{{ .Name }}s(rows)
	if err != nil {
		return item, err
	}
	{{ if .HasID  }}
		_ = items.Ids()
	{{ else }}
		_ = len(items)
	{{ end }}

	{{ if .HasID  }}
	item, err = item.Update(tx)
	if err != nil {
		return item, err
	}
	_, err = Select{{ .Name }}(tx, item.Id)
	{{ else }} 
	row := tx.QueryRow(` + "`" + `SELECT * FROM {{snake .Name}}s WHERE 
		{{range $i, $e := .Fields.ForeignKeys }}{{if $i}} AND {{end}}
		{{- if $e.Type.IsNullable -}}
			( {{ $e.SQLName }} IS NULL OR {{ $e.SQLName }} = ${{inc $i}})
		{{- else -}}
			{{ $e.SQLName }} = ${{inc $i}}
		{{- end -}}	
		{{end}};` + "`" +
		`{{range .Fields.ForeignKeys }},item.{{.GoName}}{{end}})
	_, err = Scan{{ .Name }}(row)
	{{ end }}

	return item, err
}

`))
)
