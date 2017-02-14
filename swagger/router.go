package swagger

import "github.com/gin-gonic/gin"

type ISwagRouter interface {
	Params(params ...*Param) ISwagRouter
	Body(body interface{}) ISwagRouter
	Data(data interface{}) ISwagRouter
	Info(summary, desc string) ISwagRouter
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

func (swr *SwagRouter) clear() {
	swr.params = nil
	swr.body = nil
	swr.data = nil
	swr.summary = ""
	swr.desc = ""
}

func (swr *SwagRouter) Params(ps ...*Param) ISwagRouter {
	swr.params = ps
	return swr
}

func (swr *SwagRouter) Body(body interface{}) ISwagRouter {
	swr.body = body
	return swr
}

func (swr *SwagRouter) Data(data interface{}) ISwagRouter {
	swr.data = data
	return swr
}

func (swr *SwagRouter) Info(summary, desc string) ISwagRouter {
	swr.summary = summary
	swr.desc = desc
	return swr
}

func (swr *SwagRouter) Handle(ms, route string, handler gin.HandlerFunc) {
	swr.engine.AddPath(swr.group.BasePath(), route, ms, swr.summary, swr.desc, swr.params, swr.body, swr.data)
	swr.group.Handle(ms, route, handler)
	swr.clear()
}

func (swr *SwagRouter) GET(route string, handler gin.HandlerFunc) {
	swr.Handle("GET", route, handler)
}

func (swr *SwagRouter) POST(route string, handler gin.HandlerFunc) {
	swr.Handle("POST", route, handler)
}

func (swr *SwagRouter) QueryParam(name, desc string) *Param {
	return &Param{"query", name, desc, "string", false, "", false}
}

func (swr *SwagRouter) PathParam(name, desc string) *Param {
	return &Param{"path", name, desc, "string", true, "", false}
}

func (swr *SwagRouter) FileParam(name, desc string) *Param {
	return &Param{"file", name, desc, "file", false, "form", true}
}
