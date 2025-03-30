package server

import (
	"bytes"
	"html/template"
)

const serviceTemplate = `
{{- /* delete empty line */ -}}
package service

import (
	{{- if .UseContext }}
	"context"
	{{- end }}
	pb "{{ .Package }}"
	"code.byted.org/kite/kitex/server"
	{{- if .GoogleEmpty }}
	"google.golang.org/protobuf/types/known/emptypb"
	{{- end }}
)

type {{ .Service }}ServiceImpl struct {
	server.Server
}

func New{{ .Service }}ServiceImpl() *{{ .Service }}ServiceImpl {
	return &{{ .Service }}ServiceImpl{}
}

{{- $s1 := "google.protobuf.Empty" }}
{{ range .Methods }}
{{- if eq .Type 1 }}
func (s *{{ .Service}}ServiceImpl) {{ .Name }}(ctx context.Context, req {{ if eq .Request $s1 }}*emptypb.Empty{{ else }}*pb.{{ .Request }}{{ end }}) ({{ if eq .Reply $s1 }}*emptypb.Empty{{ else }}*pb.{{ .Reply }}{{ end}}, error) {
	return {{ if eq .Reply $s1 }}&emptypb.Empty{}{{ else }}&pb.{{ .Reply }}{}{{ end }}, nil
}

{{- else if eq .Type 2 }}

{{- else if eq .Type 3 }}

{{- else if eq .Type 4 }}

{{- end }}
{{- end }}
`

// Service is a proto service.
type Service struct {
	Package     string
	Service     string
	Methods     []*Method
	GoogleEmpty bool

	UseIO      bool
	UseContext bool
}

type MethodType uint8

const (
	unaryType          MethodType = 1
	twoWayStreamsType  MethodType = 2
	requestStreamsType MethodType = 3
	returnsStreamsType MethodType = 4
)

type Method struct {
	Service string
	Name    string
	Request string
	Reply   string

	// type: unary or stream
	Type MethodType
}

func (s *Service) excute() ([]byte, error) {
	const empty = "google.protobuf.Empty"
	buf := new(bytes.Buffer)
	for _, method := range s.Methods {
		if (method.Type == unaryType && (method.Request == empty || method.Reply == empty)) ||
			(method.Type == returnsStreamsType && method.Request == empty) {
			s.GoogleEmpty = true
		}
		if method.Type == twoWayStreamsType || method.Type == requestStreamsType {
			s.UseIO = true
		}
		if method.Type == unaryType {
			s.UseContext = true
		}
	}
	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
