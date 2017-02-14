package swagger

type Definition struct {
	Required   []string               `json:"required,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Type       string                 `json:"type,omitempty"`
}

func NewDefinition() *Definition {
	return &Definition{
		Required:   make([]string, 0),
		Type:       "object",
		Properties: make(map[string]interface{}),
	}
}
