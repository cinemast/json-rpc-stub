package main

import (
	"github.com/cinemast/json-rpc-stub/codegen"
	"github.com/cinemast/json-rpc-stub/specification"
	"os/exec"
	"fmt"
)

func main() {
	fmt.Println("json-rpc-stub tool")

	_, err := exec.LookPath("quicktype")
	if err != nil {
		fmt.Println("quicktype was not found on your PATH, please install it using `npm install -g quicktype`")
		panic(err)
	}

	spec,err := specification.NewSpecification("examples/warehouse.json")

	cxx := codegen.NewJsonRpcCxx(spec, "warehouse", "WarehouseApp", "gen")
	err = cxx.GenerateTypes()
	if err != nil {
		panic(err)
	}
	err = cxx.GenerateServer()
	if err != nil {
		panic(err)
	}
	
	fmt.Println()
}