package add

import (
	"bytes"
	"strings"
	"text/template"
)

const protoTemplate = `
syntax = "proto3";

package {{.Package}};

option go_package = "{{.GoPackage}}";
option java_multiple_files = true;
option java_package = "{{.JavaPackage}}";

service {{.Service}} {
	rpc SayHello (HelloRequest) returns (HelloReply);
}

message HelloRequest {}
message HelloReply {}

`

func (p *Proto) execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("proto").Parse(strings.TrimSpace(protoTemplate))
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
