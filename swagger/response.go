package swagger

type Response struct {
	Description string      `json:"description,omitempty"`
	Schema      interface{} `json:"schema,omitempty"`
}

func (r *Response) SetSchema(define string) {
	if r.Schema == nil {
		r.Schema = make(map[string]string)
	} else if _, ok := r.Schema.(map[string]string); !ok {
		r.Schema = make(map[string]string)
	}
	r.Schema.(map[string]string)["$ref"] = "#/definitions/errorModel"
}
