package codegen

import (
	"errors"
	"os/exec"
	"strings"
	"path"
	"text/template"
	"github.com/cinemast/json-rpc-stub/specification"
	"os"	
)

type JsonRpcCxx struct {
	Namespace string
	Path string
	Class string
	Spec *specification.Specification
}

const (
	server = `#pragma once
#include <jsonrpccxx/server.hpp>
#include "types.hpp"
namespace {{.Namespace}} {
	class {{.Class}}Server {
	public:
		{{ range $key, $value := .Spec.Procedures }}
		// {{$value.Description}}
		{{ ToReturnType $value.ReturnType }} {{$key}}({{ range $index,$param := $value.Params }}{{if $index}}, {{end}}{{ToCppType $param.Type}} {{$param.Name}}{{end}});{{end}}

		bool Bind(jsonrpccxx::JsonRpcServer &server) {
			bool result = true;{{ range $key, $proc := .Spec.Procedures }}
			result &= server.Add("{{$key}}", GetHandle({{$.Class}}::{{$key}}, *this), { {{ range $index,$param := $proc.Params }}{{if $index}}, {{end}}"{{$param.Name}}"{{end}} });{{end}}
			return result;
		}
	};
}`
	client = `#pragma once
#include <jsonrpccxx/client.hpp>
#include "types.hpp"
namespace {{.Namespace}} {
	class {{.Class}}Client {
	public:
		explicit {{.Class}}(jsonrpccxx::JsonRpcClient &client) : client(client) {}
		{{ range $key, $value := .Spec.Procedures }}
		// {{$value.Description}}
		{{ ToReturnType $value.ReturnType }} {{$key}}({{ range $index,$param := $value.Params }}{{if $index}}, {{end}}{{ToCppType $param.Type}} {{$param.Name}}{{end}}) { return client.CallMethod<{{ ToReturnType $value.ReturnType }}>(1, "{{$key}}", { {{range $index,$param := $value.Params}}{{if $index}}, {{end}}{{$param.Name}}{{end}} }); }{{end}}

	private:
		jsonrpccxx::JsonRpcClient &client;
	};
}`
)

func NewJsonRpcCxx(spec *specification.Specification, namespace string, class string, path string) *JsonRpcCxx {
	os.MkdirAll(path, os.ModePerm)
	return &JsonRpcCxx{Spec: spec, Namespace: namespace, Path: path, Class: class}
}

func ToReturnType(t specification.Type) string {
	if t.RefType != "" {
		parts := strings.Split(t.RefType, "/")
		return parts[len(parts)-1]
	}
	switch t.Type {
	case specification.Boolean: return "bool"
	case specification.String: return "std::string"
	case specification.Integer: return "int"
	case specification.Number: return "double"
	case specification.Object:
	case specification.Array: 
	}
	return "void"
}

func ToCppType(t specification.Type) string {
	if t.RefType != "" {
		parts := strings.Split(t.RefType, "/")
		return "const " + parts[len(parts)-1]+"&"
	}
	switch t.Type {
	case specification.Boolean: return "bool"
	case specification.String: return "const std::string&"
	case specification.Integer: return "int"
	case specification.Number: return "double"
	case specification.Object:
	case specification.Array: 
	}
	return "void"
}

func (cxx *JsonRpcCxx) GenerateTemplate(tpl string) error {
	tmpl, err := template.New("template").Funcs(template.FuncMap{
		"ToCppType": ToCppType,
		"ToReturnType": ToReturnType,
	}).Parse(tpl)
	if err != nil {
		return err
	}

	f, err := os.Create(path.Join(cxx.Path, "server.hpp"))
	
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl.Execute(f, cxx)
	return nil
}

func (cxx *JsonRpcCxx) GenerateTypes() error {
	filePath := path.Join(cxx.Path, "types.hpp")
	cmd := exec.Command("quicktype", "--src-lang", "schema", cxx.Spec.Path, "-o", filePath, "--lang", "c++", "--code-format", "with-struct", "--no-boost", "--include-location", "global-include", "--namespace", cxx.Namespace)
	output,err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(err.Error() + ": " + string(output))
	}
	return nil
}

func (cxx *JsonRpcCxx) GenerateClient() error {
	return cxx.GenerateTemplate(client)
}

func (cxx *JsonRpcCxx) GenerateServer() error {
	return cxx.GenerateTemplate(server)
}