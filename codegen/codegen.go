package codegen

import (
)

type CodeGenerator interface {
	GenerateClient()
	GenerateServer()
}