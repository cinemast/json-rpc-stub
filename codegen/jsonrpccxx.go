package codegen

import (
	"github.com/cinemast/json-rpc-stub/specification"
	"io"
)

type JsonRpcCxx struct {
	Name string
	Writer io.Writer
	Spec *specification.Specification
}

func NewJsonRpcCxx() *JsonRpcCxx {
	return &JsonRpcCxx{}
}

func (cxx *JsonRpcCxx) GenerateClient() {

}

func (cxx *JsonRpcCxx) GenerateServer() {

}