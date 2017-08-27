package gev

import "github.com/gin-gonic/gin"

var err_hander IErrorHander = new(ErrorHander)

// 处理错误接口
type IErrorHander interface {
	Ok(c *gin.Context, data interface{})
	Err(c *gin.Context, code int, err error)
	Api(c *gin.Context, data interface{}, err error)
}

type ErrorHander int

func (this *ErrorHander) Ok(c *gin.Context, data interface{}) {
	if data != nil {
		c.IndentedJSON(200, gin.H{"code": 0, "data": data})
	}
}

func (this *ErrorHander) Err(c *gin.Context, code int, err error) {
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

func (this *ErrorHander) Api(c *gin.Context, data interface{}, err error) {
	if err != nil {
		if v, ok := err.(iApiError); ok {
			this.Err(c, v.Code(), err)
		} else {
			this.Err(c, 0, err)
		}
		return
	}
	this.Ok(c, data)
}

// 错误码接口
type iApiError interface {
	error
	Code() int
}

// 错误结构
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

func Error(code int, msg string) error {
	e := new(apiError)
	e.code = code
	e.msg = msg
	return e
}

func SetErrorHander(eh IErrorHander) {
	err_hander = eh
}
