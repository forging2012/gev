package gev

import (
	"github.com/inu1255/gohelper"

	"github.com/gin-gonic/gin"
)

func Ok(c *gin.Context, data interface{}) {
	c.IndentedJSON(200, gin.H{"code": 0, "data": data})
}
func Err(c *gin.Context, code int, err error) {
	msg := err.Error()
	if code == 0 {
		table := gohelper.Str2bytes(msg)
		count := len(table)
		if count > 32 {
			count = 32
		}
		for i := 0; i < count; i++ {
			code += int(table[i])
		}
	}
	c.IndentedJSON(200, gin.H{"code": code, "msg": msg})
}
func Api(c *gin.Context, data interface{}, err error) {
	if err != nil {
		Err(c, 0, err)
		return
	}
	Ok(c, data)
}
