package swagger

type Param struct {
	In            string `json:"in,omitempty"`
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	Type          string `json:"type,omitempty"`
	Required      bool   `json:"required,omitempty"`
	ParamType     string `json:"paramType,omitempty"`
	AllowMultiple bool   `json:"allowMultiple,omitempty"`
}
