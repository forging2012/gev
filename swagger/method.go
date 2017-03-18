package core

import (
	"strings"
)

type Method struct {
	Summary     string               `json:"summary,omitempty"`
	Description string               `json:"description,omitempty"`
	Tags        []string             `json:"tags,omitempty"`
	Parameters  []interface{}        `json:"parameters,omitempty"`
	Responses   map[string]*Response `json:"responses,omitempty"`
}

func NewMethod(summary, desc string) *Method {
	method := &Method{
		Summary:     summary,
		Description: desc,
		Tags:        make([]string, 0, 1),
		Parameters:  make([]interface{}, 0),
	}
	// method.SetResponse(MapSchema)
	return method
}

// 添加路由中的第一段作为 tag
func (m *Method) AddTag(ps string) {
	if len(ps) == 0 {
		return
	}
	index := strings.Index(ps[1:], "/")
	if index < 0 {
		m.Tags = append(m.Tags, ps[1:])
	} else {
		m.Tags = append(m.Tags, ps[1:index+1])
	}
}

func (m *Method) SetResponse(data interface{}) {
	if m.Responses == nil {
		m.Responses = make(map[string]*Response)
	}
	m.Responses["200"] = &Response{Schema: data}
}
