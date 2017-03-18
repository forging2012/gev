package gev

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

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

type IBind interface {
	Bind(g ISwagRouter, model IModel)
}

func Today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

var (
	Db, _ = xorm.NewEngine("mysql", "root:199337@/youyue?parseTime=true&loc=Asia%2FShanghai")
	App   = gin.New()
	// Db, _        = xorm.NewEngine("sqlite3", "./test.db")
	token_expire = 86400
	UserVerify   IVerifyModel
	Swag         = swagger.NewSwagger()
	_gev_path    = ""
	Log          = log.New(os.Stdout, "[ gev ]\t", log.Ltime|log.Lshortfile)
)

type RouterGroup gin.RouterGroup

func (r *RouterGroup) Bind(model IBind) {
	model.Bind(Swag.Bind((*gin.RouterGroup)(r)), nil)
}

func Bind(prefix string, model IBind, summary ...string) {
	pbd := Swag.Bind(App.Group(prefix), summary...)
	model.Bind(pbd, nil)
}

func init() {
	gin.SetMode(gin.ReleaseMode)
	// if out, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY, 0664); err == nil {
	// 	App.Use(gin.LoggerWithWriter(out))
	// 	Log = log.New(out, "[ gev ]\t", log.Ltime|log.Lshortfile)
	// 	Db.SetLogger(xorm.NewSimpleLogger(out))
	// } else {
	// 	Log.Println(out)
	// }
	App.Use(gin.Logger())
	App.Use(gin.Recovery())
	_, file, _, _ := runtime.Caller(0)
	if index := strings.LastIndex(file, "/"); index > 0 {
		_gev_path = file[:index]
	}
	CopySwagger()
}

func CopySwagger() {
	if info, err := os.Stat("api"); err != nil || !info.IsDir() {
		cmd := exec.Command("cp", "-R", _gev_path+"/api", "api")
		err := cmd.Start()
		if err != nil {
			Log.Println(err)
		}
	} else {
		Log.Println("swagger文件夹已经存在")
	}
}

func SetDb(driverName, dataSourceName string) {
	Db, _ = xorm.NewEngine(driverName, dataSourceName)
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
	if host == "" {
		host = ":8017"
	}
	Swag.WriteJson("api/swagger.json")

	Db.ShowSQL(true)
	AutoRestart()
	Server := endless.NewServer(host, App)
	Server.ListenAndServe()
}
