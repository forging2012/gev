package swagger

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type ISwagRouter interface {
	Params(params ...*Param) ISwagRouter
	Body(body interface{}) ISwagRouter
	Data(data interface{}) ISwagRouter
	Info(info ...string) ISwagRouter
	GET(route string, handler gin.HandlerFunc)
	POST(route string, handler gin.HandlerFunc)
	Handle(ms, route string, handler gin.HandlerFunc)

	QueryParam(name, desc string) *Param
	PathParam(name, desc string) *Param
	FileParam(name, desc string) *Param
}

type SwagRouter struct {
	engine  *Swagger
	group   *gin.RouterGroup
	params  []*Param
	body    interface{}
	data    interface{}
	summary string
	desc    string
}

func (this *SwagRouter) clear() {
	this.params = nil
	this.body = nil
	this.data = nil
	this.summary = ""
	this.desc = ""
}

func (this *SwagRouter) Params(ps ...*Param) ISwagRouter {
	this.params = ps
	return this
}

func (this *SwagRouter) Body(body interface{}) ISwagRouter {
	this.body = body
	return this
}

func (this *SwagRouter) Data(data interface{}) ISwagRouter {
	this.data = data
	return this
}

func (this *SwagRouter) Info(info ...string) ISwagRouter {
	if len(info) < 1 {
		return this
	}
	this.summary = info[0]
	this.desc = strings.Join(info[1:], "\n")
	return this
}

func (this *SwagRouter) Handle(ms, route string, handler gin.HandlerFunc) {
	this.engine.AddPath(this.group.BasePath(), route, ms, this.summary, this.desc, this.params, this.body, this.data)
	this.group.Handle(ms, route, handler)
	this.clear()
}

func (this *SwagRouter) GET(route string, handler gin.HandlerFunc) {
	this.Handle("GET", route, handler)
}

func (this *SwagRouter) POST(route string, handler gin.HandlerFunc) {
	// if this.body != nil {
	// 	New := this.body.New
	// 	this.Handle("POST", route, func(c *gin.Context) {
	// 		body := New()
	// 		if err := c.BindJSON(body); err != nil {
	// 			c.Set("body", body)
	// 			handler(c)
	// 		} else {
	// 			c.Set("error", err)
	// 		}
	// 	})
	// } else {
	// 	this.Handle("POST", route, handler)
	// }
	this.Handle("POST", route, handler)
}

func (this *SwagRouter) QueryParam(name, desc string) *Param {
	return &Param{"query", name, desc, "string", false, "", false}
}

func (this *SwagRouter) PathParam(name, desc string) *Param {
	return &Param{"path", name, desc, "string", true, "", false}
}

func (this *SwagRouter) FileParam(name, desc string) *Param {
	return &Param{"formData", name, desc, "file", false, "form", true}
}
