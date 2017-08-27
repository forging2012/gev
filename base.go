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
	// Class `json:"-" xorm:"-"`
	ctx *gin.Context `json:"-" xorm:"-"`
}

func (this *BaseService) Before(ctx *gin.Context) bool {
	this.ctx = ctx
	return true
}

func (this *BaseService) After(data interface{}, err error) {
	err_hander.Api(this.ctx, data, err)
}

func (this *BaseService) Finish(err interface{}) {
	if err != nil {
		Log.Printf("%v\n\033[31m%s\033[0m", err, string(stack()))
		err_hander.Err(this.ctx, 500, errors.New(fmt.Sprintf("系统错误 : %v", err)))
	}
}
