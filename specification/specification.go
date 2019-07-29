package specification

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