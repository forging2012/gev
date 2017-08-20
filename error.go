package gev

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type iApiError interface {
	error
	Code() int
}

type apiError struct {
	code int
	msg  string
}

func (this *apiError) Error() string {
	return this.msg
}

func (this *apiError) Code() int {
	return this.code
}

func Err(code int, msg string) error {
	e := new(apiError)
	e.code = code
	e.msg = msg
	return e
}

func respOk(c *gin.Context, data interface{}) {
	if data != nil {
		c.IndentedJSON(200, gin.H{"code": 0, "data": data})
	}
}

func respErr(c *gin.Context, code int, err error) {
	msg := err.Error()
	if code == 0 {
		table := str2bytes(msg)
		count := len(table)
		if count > 32 {
			count = 32
		}
		for i := 0; i < count; i++ {
			code += int(table[i])
		}
	}
	Log.Println("code:\033[41;37m", code, "\033[0m msg:", msg)
	c.IndentedJSON(200, gin.H{"code": code, "msg": msg})
}

func respApi(c *gin.Context, data interface{}, err error) {
	if err != nil {
		if v, ok := err.(iApiError); ok {
			respErr(c, v.Code(), err)
		} else {
			respErr(c, 0, err)
		}
		return
	}
	respOk(c, data)
}

func NeedAuth(c *gin.Context) (interface{}, bool) {
	if user, ok := c.Get("user"); ok {
		return user, true
	} else {
		respErr(c, 1255, errors.New("需要登录"))
		return nil, false
	}
}
