package gev

import "github.com/gin-gonic/gin"

func Ok(c *gin.Context, data interface{}) {
	c.IndentedJSON(200, gin.H{"code": 0, "data": data})
}
func Err(c *gin.Context, code int, err error) {
	c.IndentedJSON(200, gin.H{"code": code, "msg": err.Error()})
}
func Api(c *gin.Context, data interface{}, err error) {
	if err != nil {
		Err(c, 52, err)
		return
	}
	Ok(c, data)
}
