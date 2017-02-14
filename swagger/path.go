package swagger

import (
	"strings"
)

type Path map[string]*Method

func (m Path) SetMethod(ms string, method *Method) {
	m[strings.ToLower(ms)] = method
}
