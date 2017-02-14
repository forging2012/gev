package swagger

type Security struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	In   string `json:"in,omitempty"`
}
