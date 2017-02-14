package gev

import (
	// _ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/inu1255/gev/swagger"
	_ "github.com/mattn/go-sqlite3"
)

type ISwagRouter interface {
	swagger.ISwagRouter
}

var (
	// Db, _ = xorm.NewEngine("mysql", "root:199337@/youyue")
	App          = gin.Default()
	Db, _        = xorm.NewEngine("sqlite3", "./test.db")
	token_expire = 86400
	UserVerify   IVerifyModel
	Host         = "localhost:8017"
	Swag         = swagger.NewSwagger()
)

type RouterGroup gin.RouterGroup

func (r *RouterGroup) Bind(model IModel) {
	model.Bind(Swag.Bind((*gin.RouterGroup)(r)), nil)
}

func Bind(prefix string, model IModel, handlers ...gin.HandlerFunc) {
	pbd := Swag.Bind(App.Group(prefix, handlers...))
	model.Bind(pbd, nil)
}

func init() {
}

func Run() {
	Swag.Host = Host
	Swag.WriteJson("api/swagger.json")

	Db.ShowSQL(true)
	App.Run(Host)
}