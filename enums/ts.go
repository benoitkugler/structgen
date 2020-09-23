package enums

import (
	"bytes"
	"log"
	"text/template"
)

var (
	tplEnumDef = template.Must(template.New("enums_def").Parse(`{
	{{ range .Values -}}
	{{ .VarName }} = {{ .Value }},
	{{ end -}}
};
`))

	tpltEnumValues = template.Must(template.New("enums").Parse(`
	export const {{ .Name }}Labels: { [key in {{ .Name }}]: string } = {
		{{ range .Values -}}
			[{{ $.Name }}.{{ .VarName }}]: "{{ .Label }}",
		{{ end }}
	}
`))
)

func (e EnumType) TsDefinition() string {
	var out bytes.Buffer
	if err := tplEnumDef.Execute(&out, e); err != nil {
		log.Fatal(err)
	}
	if err := tpltEnumValues.Execute(&out, e); err != nil {
		log.Fatal(err)
	}
	return out.String()
}
