package enums

import (
	"bytes"
	"fmt"
	"go/types"
	"log"
	"sort"
	"text/template"

	"github.com/benoitkugler/structgen/loader"
)

var (
	tmplInt = template.Must(template.New("labels_int").Parse(`
		{{ .Name }}Labels = [...]string{
			{{ range .Values -}}
				{{ .VarName }}: "{{ .Label }}",
			{{ end -}}
		}`))

	tmplDefault = template.Must(template.New("labels").Parse(`
		{{ .Name }}Labels = map[{{ .Name }}]string{
			{{ range .Values -}}
				{{ .VarName }}: "{{ .Label }}",
			{{ end -}}
		}`))
)

// labels returns the code for
// mapping variable
func (e Type) labels() string {
	var out bytes.Buffer
	tmpl := tmplDefault
	if e.IsInt {
		tmpl = tmplInt
	}
	if err := tmpl.Execute(&out, e); err != nil {
		log.Fatal(err)
	}
	return out.String()
}

var _ loader.Handler = Handler{}

type Handler struct {
	Enums       EnumTable
	PackageName string
}

func (d Handler) HandleType(typ types.Type) loader.Type      { return nil }
func (d Handler) HandleComment(comment loader.Comment) error { return nil }

func (d Handler) Header() string {
	out := fmt.Sprintf("package %s \n // DO NOT EDIT - autogenerated from enums.go \n\n  var ( \n ", d.PackageName)

	// print in deterministic order
	var keys []string
	for key := range d.Enums {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		out += fmt.Sprintln(d.Enums[key].labels())
	}
	out += ")"
	return out
}

func (d Handler) Footer() string { return "" }
