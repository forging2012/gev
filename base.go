package gev

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type IService interface {
	Before(ctx *gin.Context) bool
	Finish(err interface{})
	After(data interface{}, err error)
}

type BaseService struct {
	Class `json:"-" xorm:"-"`
	Ctx   *gin.Context `json:"-" xorm:"-"`
}

func (this *BaseService) Before(ctx *gin.Context) bool {
	this.Ctx = ctx
	return true
}
func (this *BaseService) After(data interface{}, err error) {
	respApi(this.Ctx, data, err)
}
func (this *BaseService) Finish(err interface{}) {
	if err != nil {
		Log.Printf("%v\n\033[31m%s\033[0m", err, string(stack()))
		respErr(this.Ctx, 500, errors.New(fmt.Sprintf("系统错误 : %v", err)))
	}
}
