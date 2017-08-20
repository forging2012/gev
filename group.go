package gev

import (
	"github.com/gin-gonic/gin"
	"github.com/inu1255/go-swagger/core"
)

type IRouter interface {
	DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	HEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	OPTIONS(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	PATCH(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	Handle(httpMethod string, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	Group(string, ...gin.HandlerFunc) *RouterGroup
	Body(body interface{}) IRouter
	Data(data interface{}) IRouter
	Info(info ...string) IRouter
	QueryParam(name, desc string) *core.Param
	PathParam(name, desc string) *core.Param
	FileParam(name, desc string) *core.Param
}

type RouterGroup struct {
	*gin.RouterGroup
	swag *core.SwagRouter
}

func (this *RouterGroup) Any(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	this.swag.AddPath(this.BasePath(), relativePath, "GET")
	this.swag.AddPath(this.BasePath(), relativePath, "POST")
	this.swag.AddPath(this.BasePath(), relativePath, "PUT")
	this.swag.AddPath(this.BasePath(), relativePath, "PATCH")
	this.swag.AddPath(this.BasePath(), relativePath, "HEAD")
	this.swag.AddPath(this.BasePath(), relativePath, "OPTIONS")
	this.swag.AddPath(this.BasePath(), relativePath, "DELETE")
	this.swag.AddPath(this.BasePath(), relativePath, "CONNECT")
	this.swag.AddPath(this.BasePath(), relativePath, "TRACE")
	this.swag.Clear()
	return this.RouterGroup.Any(relativePath, handlers...)
}

func (this *RouterGroup) DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return this.Handle("DELETE", relativePath, handlers...)
}

func (this *RouterGroup) GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return this.Handle("GET", relativePath, handlers...)
}

func (this *RouterGroup) HEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return this.Handle("HEAD", relativePath, handlers...)
}

func (this *RouterGroup) OPTIONS(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return this.Handle("OPTIONS", relativePath, handlers...)
}

func (this *RouterGroup) PATCH(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return this.Handle("PATCH", relativePath, handlers...)
}

func (this *RouterGroup) POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return this.Handle("POST", relativePath, handlers...)
}

func (this *RouterGroup) PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return this.Handle("PUT", relativePath, handlers...)
}

func (this *RouterGroup) Handle(httpMethod string, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	this.swag.AddPath(this.BasePath(), relativePath, httpMethod)
	this.swag.Clear()
	return this.RouterGroup.Handle(httpMethod, relativePath, handlers...)
}

func (this *RouterGroup) Group(relativePath string, handlers ...gin.HandlerFunc) *RouterGroup {
	group := new(RouterGroup)
	group.RouterGroup = this.RouterGroup.Group(relativePath, handlers...)
	group.swag = this.swag
	return group
}

func (this *RouterGroup) Body(body interface{}) IRouter {
	this.swag.Body(body)
	return this
}

func (this *RouterGroup) Data(data interface{}) IRouter {
	this.swag.Data(data)
	return this
}

func (this *RouterGroup) Info(info ...string) IRouter {
	this.swag.Info(info...)
	return this
}

func (this *RouterGroup) QueryParam(name, desc string) *core.Param {
	return this.swag.QueryParam(name, desc)
}

func (this *RouterGroup) PathParam(name, desc string) *core.Param {
	return this.swag.PathParam(name, desc)
}

func (this *RouterGroup) FileParam(name, desc string) *core.Param {
	return this.swag.FileParam(name, desc)
}
