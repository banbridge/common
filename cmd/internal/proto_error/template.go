package proto_error

import (
	"bytes"
	"text/template"
)

var errorsTemplate = `
{{ range .Errors }}

{{ if .HasComment }}{{ .Comment }}{{ end -}}
func Is{{.CamelValue}}(ctx context.Context, err error) bool {
	if err == nil {
		return false
	}
	e := biz_err.FromError(ctx, err)
	return e.Code() == "{{ .BizCode }}"
}

{{ if .HasComment }}{{ .Comment }}{{ end -}}
func Error{{ .CamelValue }}(ctx context.Context, format string, args ...any) *errors.BizError {
	return biz_err.NewError(ctx, "{{ .BizCode }}", fmt.Sprintf(format, args...),
		biz_err.WithHttpStatus({{ .HTTPCode }}), biz_err.WithBizMsg("{{ .BizMsg }}"), biz_err.WithReason({{ .Name }}_{{ .Value }}.String()), biz_err.WithDepth(3))	
}

{{- end }}
`

type errorInfo struct {
	Name       string
	Value      string
	HTTPCode   int
	BizCode    string
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
