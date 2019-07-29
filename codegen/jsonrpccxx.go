package codegen

import (
	"strings"
	"text/template"
	"github.com/cinemast/json-rpc-stub/specification"
	"io"
)

type JsonRpcCxx struct {
	Writer io.Writer
	Namespace string
	Class string
	Spec *specification.Specification
}

const (
	server = `#pragma once
#include <jsonrpccxx/server.hpp>
#include "types.h"
namespace {{.Namespace}} {
	class {{.Class}} {
	public:
		{{ range $key, $value := .Spec.Procedures }}
		// {{$value.Description}}
		{{ ToReturnType $value.ReturnType }} {{$key}}({{ range $index,$param := $value.Params }}{{if $index}}, {{end}}{{ToCppType $param.Type}} {{$param.Name}}{{end}});{{end}}

		bool Bind(jsonrpccxx::JsonRpcServer &server) {
			bool result = true;{{ range $key, $proc := .Spec.Procedures }}
			result &= server.Add("{{$key}}", GetHandle({{$.Class}}::{{$key}}, *this), { {{ range $index,$param := $proc.Params }}{{if $index}}, {{end}}"{{$param.Name}}"{{end}} });{{end}}
			return result;
		}
		//TODO: Add binding method
	};
}`
	client = `#pragma once
#include <jsonrpccxx/client.hpp>
#include "types.h"
namespace {{.Namespace}} {
	class {{.Class}} {
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

func NewJsonRpcCxx(writer io.Writer, spec *specification.Specification) *JsonRpcCxx {
	return &JsonRpcCxx{Writer: writer, Spec: spec, Namespace: "warehouse::foo", Class: "WarehouseClient"}
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
	tmpl.Execute(cxx.Writer, cxx)
	return nil
}

func (cxx *JsonRpcCxx) GenerateClient() error {
	return cxx.GenerateTemplate(client)
}

func (cxx *JsonRpcCxx) GenerateServer() error {
	return cxx.GenerateTemplate(server)
}