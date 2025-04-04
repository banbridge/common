GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
API_PROTO_FILES=$(shell find api -name *.proto)
ERROR_PROTO_FILES=$(shell find errorx -name *.proto)
API_PB_GO_FILES=$(shell find api -name *.pb.go)
CONF_PROTO_FILES=$(shell find internal -name conf*.proto)
REPO_PROTO_FILES=$(shell find repository -name *.proto)
PLATFORM=$(shell uname)
PROTOCVER=$(shell protoc --version | awk '{print $1}')

.PHONY: test
# local build
test:
# 	go build -o $(GOPATH)/bin/protoc-gen-go-kitex cmd/petal/protoc-gen-go-kitex/main.go
# 	go build -o $(GOPATH)/bin/protoc-gen-go-hertz cmd/petal/protoc-gen-go-hertz/main.go
	go build -o $(GOPATH)/bin/protoc-gen-go-error cmd/protoc-gen-go-error/main.go
# 	go build -o $(GOPATH)/bin/protoc-gen-go-fastpb cmd/petal/protoc-gen-go-fastpb/main.go
