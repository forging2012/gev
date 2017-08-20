package gev

import (
	"os"
	"strings"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/inu1255/go-swagger/core"
)

type Engine struct {
	*gin.Engine
	swag *core.SwagRouter
}

func (this *Engine) Any(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
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
	return this.Engine.Any(relativePath, handlers...)
}

func (this *Engine) GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return this.Handle("GET", relativePath, handlers...)
}

func (this *Engine) POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return this.Handle("POST", relativePath, handlers...)
}

func (this *Engine) PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return this.Handle("PUT", relativePath, handlers...)
}

func (this *Engine) PATCH(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return this.Handle("PATCH", relativePath, handlers...)
}

func (this *Engine) HEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return this.Handle("HEAD", relativePath, handlers...)
}

func (this *Engine) OPTIONS(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return this.Handle("OPTIONS", relativePath, handlers...)
}

func (this *Engine) DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return this.Handle("DELETE", relativePath, handlers...)
}

func (this *Engine) Group(relativePath string, handlers ...gin.HandlerFunc) *RouterGroup {
	group := new(RouterGroup)
	group.RouterGroup = this.Engine.Group(relativePath, handlers...)
	group.swag = this.swag
	return group
}

func (this *Engine) Handle(httpMethod string, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	this.swag.AddPath(this.BasePath(), relativePath, httpMethod)
	this.swag.Clear()
	return this.Engine.Handle(httpMethod, relativePath, handlers...)
}

func (this *Engine) Body(body interface{}) IRouter {
	this.swag.Body(body)
	return this
}

func (this *Engine) Data(data interface{}) IRouter {
	this.swag.Data(data)
	return this
}

func (this *Engine) Info(info ...string) IRouter {
	this.swag.Info(info...)
	return this
}

func (this *Engine) QueryParam(name, desc string) *core.Param {
	return this.swag.QueryParam(name, desc)
}

func (this *Engine) PathParam(name, desc string) *core.Param {
	return this.swag.PathParam(name, desc)
}

func (this *Engine) FileParam(name, desc string) *core.Param {
	return this.swag.FileParam(name, desc)
}

func (this *Engine) Swagger(relativePath string) {
	core.CopySwagger()
	this.swag.WriteJson("api/swagger.json")
	this.Static(relativePath, "api")
}

func (this *Engine) Run(addr ...string) (err error) {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		Log.Println("go run 不启动热更新", os.Args[0])
		return this.Engine.Run(addr...)
	}
	AutoRestart()
	Server := endless.NewServer(resolveAddress(addr), this.Engine)
	return Server.ListenAndServe()
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); len(port) > 0 {
			return ":" + port
		}
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too much parameters")
	}
}

func New() *Engine {
	e := gin.New()
	engine := new(Engine)
	engine.Engine = e
	engine.swag = core.NewSwagRouter()
	return engine
}
