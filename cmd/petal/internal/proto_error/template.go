package proto_error

import (
	"bytes"
	"text/template"
)

var errorsTemplate = `
{{ range .Errors }}

{{ if .HasComment }}{{ .Comment }}{{ end -}}
func Is{{.CamelValue}}(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == {{ .Name }}_{{ .Value }}.String() && e.Code == {{ .HTTPCode }} 
}

{{ if .HasComment }}{{ .Comment }}{{ end -}}
func Error{{ .CamelValue }}(ctx context.Context, format string, args ...interface{}) *errors.Error {
	 return errors.New(ctx, {{ .HTTPCode }}, int({{ .Name }}_{{ .Value }}.Number()), {{ .Name }}_{{ .Value }}.String(), "{{ .BizMsg }}", fmt.Sprintf(format, args...))
}

{{- end }}
`

type errorInfo struct {
	Name       string
	Value      string
	HTTPCode   int
	BizCode    int
	CamelValue string
	Comment    string
	HasComment bool
	BizMsg     string
}

type errorWrapper struct {
	Errors []*errorInfo
}

func (e *errorWrapper) execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("errorx").Parse(errorsTemplate)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, e); err != nil {
		panic(err)
	}
	return buf.String()
}
