package gev

import (
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/inu1255/gev/swagger"
	// _ "github.com/mattn/go-sqlite3"
)

type ISwagRouter interface {
	swagger.ISwagRouter
}

var (
	Db, _ = xorm.NewEngine("mysql", "root:199337@/youyue")
	App   = gin.Default()
	// Db, _        = xorm.NewEngine("sqlite3", "./test.db")
	token_expire = 86400
	UserVerify   IVerifyModel
	Host         = ""
	Swag         = swagger.NewSwagger()
	PkgPath      = ""
	Log          = log.New(os.Stdout, "[ gev ]\t", log.Ltime|log.Lshortfile)
)

type RouterGroup gin.RouterGroup

func (r *RouterGroup) Bind(model IModel) {
	model.Bind(Swag.Bind((*gin.RouterGroup)(r)), nil)
}

func Bind(prefix string, model IModel, summary ...string) {
	pbd := Swag.Bind(App.Group(prefix), summary...)
	model.Bind(pbd, nil)
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	if index := strings.LastIndex(file, "/"); index > 0 {
		PkgPath = file[:index]
	}
}

func Cross() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "x-auth-token,x-device,x-uuid,content-type")
			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(200)
			}
		}
	}
}

func Description(info ...string) {
	Swag.Info.Add(info...)
}

func Run(host string) {
	if host != "" {
		Host = host
	}
	Swag.Host = Host
	Swag.WriteJson("api/swagger.json")

	Db.ShowSQL(true)
	AutoRestart()
	endless.ListenAndServe(":8017", App)
}
