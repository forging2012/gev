package core

import "strings"

type Info struct {
	Version     string `json:"version,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

func (i *Info) Add(desc ...string) {
	i.Description = i.Description + "\n<br/>" + strings.Join(desc, "\n<br/>")
}
