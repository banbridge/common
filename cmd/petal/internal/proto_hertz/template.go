package proto_hertz

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	_ "embed"
)

//go:embed template.go.tpl
var tpl string

type serviceDesc struct {
	Name      string // Greeter
	FullName  string // helloworld.Greeter
	FilePath  string // api/helloword/helloworld.proto
	Methods   []*methodDesc
	MethodSet map[string]*methodDesc
}

func (s *serviceDesc) execute() string {
	s.MethodSet = make(map[string]*methodDesc)
	for _, m := range s.Methods {
		s.MethodSet[m.Name] = m
	}

	buf := new(bytes.Buffer)
	tmpl, err := template.New("http").Parse(strings.TrimSpace(tpl))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return strings.Trim(buf.String(), "\r\n")
}

// InterfaceName service interface name
func (s *serviceDesc) InterfaceName() string {
	return s.Name + "HTTPServer"
}

func (s *serviceDesc) NameHttp() string {
	return s.Name + "Http"
}

type methodDesc struct {
	// method
	Name    string // SayHello
	Num     int    // 一个rpc方法可以对应多个http请求
	Request string // SayHelloReq
	Reply   string // SayHelloResp
	// http_rule
	Path         string // 路由
	Method       string // HTTP Method
	Body         string
	ResponseBody string
}

// HandlerName for hertz handler name
func (m *methodDesc) HandlerName() string {
	return fmt.Sprintf("%s_%d", m.Name, m.Num)
}

// Metadata Key for method
func (m *methodDesc) MdMethodKey() string {
	return fmt.Sprintf("Md_%s_%d", m.Name, m.Num)
}

// ResponseBody type: raw、openapi、custom
// default: raw
func (m *methodDesc) ResponseBodyRaw() bool {
	return m.ResponseBody == "raw"
}

func (m *methodDesc) ResponseBodyOpenAPI() bool {
	return m.ResponseBody == "openapi"
}

func (m *methodDesc) ResponseBodyCustom() bool {
	return m.ResponseBody == "custom"
}

// initPathParams {xx} -> :xx
func (m *methodDesc) initPathParams() {
	paths := strings.Split(m.Path, "/")
	for i, p := range paths {
		if len(p) > 0 && (p[0] == '{' && p[len(p)-1] == '}') {
			paths[i] = ":" + p[1:len(p)-1]
		}
	}
	m.Path = strings.Join(paths, "/")
}
