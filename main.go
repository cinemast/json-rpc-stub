package main

import (
	"io/ioutil"
	"github.com/cinemast/json-rpc-stub/specification"
	"encoding/json"
	"os"
	"fmt"
)

func main() {
	fmt.Println("json-rpc-stub tool")
	jsonFile, err := os.Open("examples/warehouse.json")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result specification.Specification
	
	if err := json.Unmarshal([]byte(byteValue), &result); err != nil {
		panic(err)
	}

	fmt.Println(result.Procedures["AddProduct"].Params[0].Type.RefType)
	
	fmt.Println()
}