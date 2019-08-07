package specification

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type ParamType string

const (
	Array ParamType = "array"
	Object ParamType = "object"
	Integer ParamType = "integer"
	Number ParamType = "number"
	String ParamType = "string"
	Boolean ParamType = "boolean"
	Null ParamType = "null"
)

type Specification struct {
	Schema string `json:"$schema"`
	Version string `json:"version"`
	Path string
	Procedures map[string]Procedure `json:"procedures"`
}

type Procedure struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Params []Parameter `json:"params"`
	ReturnType Type `json:"returns"`
}

type Type struct {
	Type ParamType `json:"type"`
	RefType string `json:"$ref"`
}

type Parameter struct {
	Type
	Name string `json:"name"`
}

func NewSpecification(path string) (*Specification, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result Specification
	result.Path = path
	if err := json.Unmarshal([]byte(byteValue), &result); err != nil {
		return nil, err
	}
	return &result, nil
}